package kart_http

import (
	"context"
	"net/http"

	"github.com/costa92/errors"
	"github.com/costa92/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/kart-io/kart/transports/kart-http/middlewares"
)

type Server struct {
	GinEngin            *gin.Engine
	httpServer          *http.Server
	InsecureServingInfo *InsecureServingInfo
	middlewares         []string
	healthz             bool
	enableMetrics       bool
	enableProfiling     bool
}

func NewServer(opts ...Option) *Server {
	srv := &Server{}
	for _, o := range opts {
		o(srv)
	}
	initAPIServer(srv)
	return srv
}

func initAPIServer(s *Server) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

func (s *Server) Setup() {
	gin.ForceConsoleColor()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Infow("gin endpoint setup ", "httpMethod", httpMethod, "absolutePath",
			absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}
}

func (s *Server) InstallMiddlewares() {
	for _, m := range s.middlewares {
		mw, ok := middlewares.Middlewares[m]
		if !ok {
			continue
		}
		s.GinEngin.Use(mw)
	}
}

func (s *Server) InstallAPIs() {
	// Healthz 检测健康
	if s.healthz {
		s.GinEngin.GET("/healthz", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, http.StatusText(http.StatusOK))
		})
	}
	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.GinEngin)
	}

	// install pprof handler
	if s.enableProfiling {
		pprof.Register(s.GinEngin)
	}
}

func (s *Server) Start(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorw("appService recover err", "err", err)
		}
	}()

	s.httpServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s.GinEngin,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	logger.Infow("start run http server", "address", s.InsecureServingInfo.Address)
	if err := s.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) { // 如果是关闭状态，不当异常处理
			logger.Errorw("start run failed server:", "address", s.InsecureServingInfo.Address)
			return err
		}
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Infow("[HTTP] server stopping")
	return s.httpServer.Shutdown(ctx)
}
