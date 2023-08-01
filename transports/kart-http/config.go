package kart_http

import "github.com/gin-gonic/gin"

type ServerConfig struct {
	Name            string
	InsecureServing *InsecureServingInfo
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
	}
}

// CompletedConfig is the completed configuration for HttpServer.
type CompletedConfig struct {
	*ServerConfig
}

func (c *ServerConfig) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) New() *Server {
	gin.SetMode(c.Mode)
	srv := &Server{
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		GinEngin:            gin.New(),
	}
	initAPIServer(srv)
	return srv
}

// InsecureServingInfo  http 服务
type InsecureServingInfo struct {
	Address string
}
