package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pikvm_automator "github.com/sealbro/pikvm-automator"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/config"
	"github.com/sealbro/pikvm-automator/internal/grpc_ext"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/repository"
	"github.com/sealbro/pikvm-automator/internal/server"
	"github.com/sealbro/pikvm-automator/internal/services"
	"github.com/sealbro/pikvm-automator/internal/sqlite"
	"github.com/sealbro/pikvm-automator/pkg/pikvm"
	"github.com/sethvargo/go-envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	var conf config.PiKvmAutomatorConfig
	if err := envconfig.Process(ctx, &conf); err != nil {
		logger.ErrorContext(ctx, "failed to process config", slog.Any("err", err))
		return
	}

	piKvmClient := pikvm.NewPiKvmClient(logger, conf.PiKvmConfig)

	player := queue.NewExpressionPlayer(logger)
	sender := player.Start(ctx)

	clientCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	receiverChan := make(chan []byte, 20)
	defer func() {
		close(receiverChan)
		receiverChan = nil
	}()
	receiver := func(data []byte) {
		if receiverChan == nil {
			return
		}
		receiverChan <- data
	}

	err := piKvmClient.Start(clientCtx, sender, receiver)
	if err != nil {
		logger.ErrorContext(ctx, "pikvm client start", slog.Any("err", err))
		return
	}

	trigger := queue.NewExpressionTrigger(logger, player)

	db, err := sqlite.New(conf.DatabasePath)
	if err != nil {
		logger.ErrorContext(ctx, "sqlite new", slog.Any("err", err))
		return
	}
	err = sqlite.ApplySchema(ctx, db, pikvm_automator.SchemaSql)
	if err != nil {
		logger.WarnContext(ctx, "sqlite apply schema", slog.Any("err", err))
	}

	queries := repository.New(db)

	templateReplacer := services.NewTemplateReplacer(logger, queries, conf)
	automatorServer := server.NewPiKvmAutomatorServer(logger, player, queries, templateReplacer, trigger, conf)

	grpc_ext.NewGRPC(logger, conf.GatewayConfig).
		AddHTTPGateway(conf.GrpcGatewayAddress).
		AddServerImplementation(func(registrar grpc.ServiceRegistrar, mux *runtime.ServeMux) error {
			//frontend.AddFrontend(mux)
			gen.RegisterPiKvmAutomatorServer(registrar, automatorServer)
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			return gen.RegisterPiKvmAutomatorHandlerFromEndpoint(ctx, mux, conf.GrpcAddress, opts)
		}).Run()

	go func() {
		for {
			bytes, ok := <-receiverChan
			if !ok {
				logger.InfoContext(ctx, "receiver closed")
				return
			}

			logger.DebugContext(ctx, "receive", slog.Any("event", string(bytes)))
			//fmt.Println(string(bytes))

			recvEvent := &pikvm.PiKVMRecvEvent{}
			recvErr := recvEvent.UnmarshalJSON(bytes)
			if recvErr != nil {
				logger.WarnContext(ctx, "unmarshal event", slog.Any("err", recvErr))
			} else {
				if recvEvent.EventType == pikvm.HidState {
					event := recvEvent.Event.(pikvm.HIDStateEvent)
					if event.Online {
						trigger.Rise(queue.PiKvmHidOnline)
					}
				}
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
