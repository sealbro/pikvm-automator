package grpc_ext

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log/slog"
	"math"
	"net"
	"net/http"
	"time"
)

var maxSendMsgSize = grpc.MaxSendMsgSize(math.MaxInt32)

type ServerRegistrationFunc func(registrar grpc.ServiceRegistrar, mux *runtime.ServeMux) error

type Server struct {
	grpcSever            *grpc.Server
	httpServer           *http.Server
	httpMux              *runtime.ServeMux
	grpcClientConnection *grpc.ClientConn
	logger               *slog.Logger
	config               GatewayConfig
}

// NewGRPC new grpc server
func NewGRPC(logger *slog.Logger, config GatewayConfig) *Server {
	return &Server{
		logger:    logger,
		config:    config,
		grpcSever: grpc.NewServer(maxSendMsgSize),
	}
}

// AddHTTPGateway adds HTTP grpc-gateway
func (s *Server) AddHTTPGateway(address string) *Server {
	s.httpMux = runtime.NewServeMux()
	s.httpServer = &http.Server{
		Addr:              address,
		Handler:           s.httpMux,
		ReadHeaderTimeout: time.Minute * 5,
	}
	return s
}

// AddServerImplementation adds server implementation
func (s *Server) AddServerImplementation(regFunc ServerRegistrationFunc) *Server {
	if err := regFunc(s.grpcSever, s.httpMux); err != nil {
		s.logger.Error("failed to add server implementation", slog.Any("err", err))
		return nil
	}
	return s
}

// AddGrpcHealthCheck adds grpc health check
func (s *Server) AddGrpcHealthCheck() *Server {
	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s.grpcSever, hs)
	return s
}

func (s *Server) Run() {
	go func() {
		if err := s.ListenAndServe(); err != nil {
			s.logger.Error("can't run listener", slog.Any("err", err))
		}
	}()
}

// ListenAndServe start grpc-gateway
func (s *Server) ListenAndServe() error {
	network := s.config.GrpcProtocol
	address := s.config.GrpcAddress
	if s.httpServer != nil {
		s.logger.Info("Starting gRPC and REST servers", slog.String("grpc", fmt.Sprintf("%v:%v", network, address)), slog.String("http", s.httpServer.Addr))
	} else {
		s.logger.Info("Starting gRPC server", slog.String("grpc", fmt.Sprintf("%v:%v", network, address)))
	}

	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	if s.httpServer != nil {
		group.Go(func() error {
			return s.httpServer.ListenAndServe()
		})
	}

	group.Go(func() error {
		return s.grpcSever.Serve(l)
	})

	return group.Wait()
}

// Shutdown to run on app shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	s.grpcSever.GracefulStop()
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
