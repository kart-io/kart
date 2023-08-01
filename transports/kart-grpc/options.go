package kart_grpc

type Option func(grpcServer *GrpcServer)

func WithConfig(config *GrpcConfig) Option {
	return func(grpcServer *GrpcServer) {
		grpcServer.config = config
	}
}
