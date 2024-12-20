package grpc_ext

type GatewayConfig struct {
	GrpcPassthroughAuth bool   `env:"GRPC_PASSTHROUGH_AUTH, default=true"`
	GrpcProtocol        string `env:"GRPC_PROTOCOL, default=tcp"`
	GrpcAddress         string `env:"GRPC_ADDRESS, default=:32023"`
	GrpcGatewayAddress  string `env:"GRPC_GATEWAY_ADDRESS, default=:8032"`
}
