package kart_http

import (
	"github.com/gin-gonic/gin"
)

type Option func(server *Server)

func WithGinEngin(engin *gin.Engine) Option {
	return func(s *Server) {
		s.GinEngin = engin
	}
}

func WithMiddlewares(middlewares []string) Option {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithHealthz(healthz bool) Option {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithEnableMetrics(enableMetrics bool) Option {
	return func(s *Server) {
		s.enableMetrics = enableMetrics
	}
}

func WithEnableProfiling(enableProfiling bool) Option {
	return func(s *Server) {
		s.enableProfiling = enableProfiling
	}
}

func WithInsecureServingInfo(insecureServingInfo *InsecureServingInfo) Option {
	return func(s *Server) {
		s.InsecureServingInfo = insecureServingInfo
	}
}
