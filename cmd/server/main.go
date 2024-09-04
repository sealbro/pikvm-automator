package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/config"
	"github.com/sealbro/pikvm-automator/internal/grpc_ext"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/server"
	"github.com/sealbro/pikvm-automator/internal/services"
	"github.com/sealbro/pikvm-automator/internal/storage"
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
	sent := player.Start(ctx)

	clientCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	err, receive := piKvmClient.Start(clientCtx, sent)
	if err != nil {
		logger.ErrorContext(ctx, "pikvm client start", slog.Any("err", err))
		return
	}

	trigger := queue.NewExpressionTrigger(logger, player)

	commandRepository := storage.NewCommandRepository(conf.CommandsPath)
	templateReplacer := services.NewTemplateReplacer(logger, commandRepository, conf)
	automatorServer := server.NewPiKvmAutomatorServer(logger, player, commandRepository, templateReplacer, trigger, conf)

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
			bytes, ok := <-receive
			if !ok {
				logger.InfoContext(ctx, "receive closed")
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
