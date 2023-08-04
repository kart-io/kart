package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/kart-io/kart/transports"
	kart_grpc "github.com/kart-io/kart/transports/kart-grpc"
	kartHttp "github.com/kart-io/kart/transports/kart-http"
)

var ProviderHttpSeverSet = wire.NewSet(NewConfig, NewHttpSever)

type App struct {
	GenericAPIServer *transports.GenericAPIServer
}

func NewConfig() *kartHttp.ServerConfig {
	return &kartHttp.ServerConfig{
		InsecureServing: &kartHttp.InsecureServingInfo{
			Address: "0.0.0.0:8080",
		},
		Healthz:       true,
		EnableMetrics: true,
	}
}

func NewHttpSever(config *kartHttp.ServerConfig, engine *gin.Engine) (*App, error) {
	gs, err := initSever(config, engine)
	if err != nil {
		return nil, err
	}
	return &App{
		GenericAPIServer: gs,
	}, nil
}

func initSever(config *kartHttp.ServerConfig, handler *gin.Engine) (*transports.GenericAPIServer, error) {
	// 实例化 http 服务
	httpServer := kartHttp.NewServer(
		kartHttp.WithGinEngin(handler),
		kartHttp.WithInsecureServingInfo(config.InsecureServing),
	)

	grcConfig := &kart_grpc.GrpcConfig{
		Port: "8081",
		Addr: "0.0.0.0",
	}
	grpcServer := kart_grpc.NewGrpcServer(kart_grpc.WithConfig(grcConfig))
	// 运行 http 与 rpc
	gs := transports.NewGenericAPIServer(
		transports.Server(
			httpServer,
			grpcServer,
		),
		transports.Name(config.Name),
	)
	return gs, nil
}

func (a *App) Run() error {
	return a.GenericAPIServer.Run()
}
