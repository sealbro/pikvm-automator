package grpc_ext

type GatewayConfig struct {
	GrpcProtocol       string `env:"GRPC_PROTOCOL, default=tcp"`
	GrpcAddress        string `env:"GRPC_ADDRESS, default=:32000"`
	GrpcGatewayAddress string `env:"GRPC_GATEWAY_ADDRESS, default=:8032"`
}
