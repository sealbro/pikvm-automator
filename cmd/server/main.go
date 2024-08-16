package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gen "github.com/sealbro/pikvm-automator/generated/go"
	"github.com/sealbro/pikvm-automator/internal/config"
	"github.com/sealbro/pikvm-automator/internal/grpc_ext"
	"github.com/sealbro/pikvm-automator/internal/queue"
	"github.com/sealbro/pikvm-automator/internal/server"
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

	// https://docs.pikvm.org/api/#websocket-events
	// https://github.com/pikvm/kvmd/blob/d7675cdf231896851c492e1e0e719d2ccdc28082/hid/arduino/lib/drivers/usb-keymap.h

	piKvmClient := pikvm.NewPiKvmClient(logger, conf.PiKvmConfig)

	player := queue.NewExpressionPlayer(logger)
	sent := player.Play(ctx)

	clientCtx, clientCancel := context.WithCancel(ctx)
	defer clientCancel()

	err, receive := piKvmClient.Start(clientCtx, sent)
	if err != nil {
		logger.ErrorContext(ctx, "pikvm client start", slog.Any("err", err))
		return
	}

	go func() {
		for {
			bytes, ok := <-receive
			if !ok {
				logger.InfoContext(ctx, "receive closed")
				return
			}
			logger.DebugContext(ctx, "receive", slog.String("data", string(bytes)))
		}
	}()

	automatorServer := server.NewPiKvmAutomatorServer(player)

	grpc_ext.NewGRPC(logger, conf.GatewayConfig).
		AddHTTPGateway(conf.GrpcGatewayAddress).
		AddServerImplementation(func(registrar grpc.ServiceRegistrar, mux *runtime.ServeMux) error {
			gen.RegisterPiKvmAutomatorServer(registrar, automatorServer)
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			return gen.RegisterPiKvmAutomatorHandlerFromEndpoint(ctx, mux, conf.GrpcAddress, opts)
		}).Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
