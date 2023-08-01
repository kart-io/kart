package app

import (
	"github.com/gin-gonic/gin"
	"github.com/kart-io/kart/example/wire-example/config"
	"github.com/kart-io/kart/example/wire-example/controller"
	"github.com/kart-io/kart/example/wire-example/routers"
	"github.com/kart-io/kart/example/wire-example/service"
	"github.com/kart-io/kart/transports"
	kartHttp "github.com/kart-io/kart/transports/kart-http"
)

// Run 通过 Complete 函数完成
func Run(config *config.Config) error {
	newServerConfig, _ := buildGenericConfig(config)
	httpServer := newServerConfig.Complete().New()

	// 实例化
	_ = service.NewBaseService()
	// 实例化控制器
	apiCtr := controller.ProvideApiController()

	routers.InitRoute(httpServer.GinEngin, apiCtr)
	// 实例化服务
	gs := transports.NewGenericAPIServer(
		transports.Server(
			httpServer,
		),
	)
	// 运行服务
	return gs.Run()
}

func buildGenericConfig(cfg *config.Config) (newServerConfig *kartHttp.ServerConfig, lastErr error) {
	// 实例化参数
	newServerConfig = kartHttp.NewServerConfig()
	// 重新给 Config 赋值
	if lastErr = cfg.ServerRunOption.ApplyTo(newServerConfig); lastErr != nil {
		return
	}
	if lastErr = cfg.InsecureServingOptions.ApplyTo(newServerConfig); lastErr != nil {
		return
	}
	return
}

// RunV1 通过 With 方式处理
func RunV1(config *config.Config) error {
	var lastErr error
	// 实例化参数
	newServerConfig := kartHttp.NewServerConfig()
	// 重新给 Config 赋值
	if lastErr = config.ServerRunOption.ApplyTo(newServerConfig); lastErr != nil {
		return nil
	}
	if lastErr = config.InsecureServingOptions.ApplyTo(newServerConfig); lastErr != nil {
		return nil
	}

	handler := gin.Default()
	// 实例化 http 服务
	httpServer := kartHttp.NewServer(
		kartHttp.WithGinEngin(handler),
		kartHttp.WithInsecureServingInfo(newServerConfig.InsecureServing),
	)
	// 添加路由
	routers.InitRoute(httpServer.GinEngin, nil)
	// 实例化服务
	gs := transports.NewGenericAPIServer(
		transports.Server(
			httpServer,
		),
	)
	// 运行服务
	return gs.Run()
}
