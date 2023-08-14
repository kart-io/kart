package kart_grpc

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/costa92/logger"
	openmetrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/openmetrics/v2"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/providers/zap/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/kart-io/kart/internal/host"
)

const HEALTHCHECK_SERVICE = "grpc.health.v1.Health"

type GrpcServer struct {
	*grpc.Server
	ctx                        context.Context
	listener                   net.Listener
	once                       sync.Once
	err                        error
	config                     *GrpcConfig
	endpoint                   *url.URL
	timeout                    time.Duration
	customInterceptorOuterMost bool // default is false, so custom interceptor will be inner most by default
	unaryInterceptor           []grpc.UnaryServerInterceptor
	streamInterceptor          []grpc.StreamServerInterceptor
	grpcOpts                   []grpc.ServerOption
	health                     *health.Server

	metric             bool
	metricSubSystem    string
	metricHandlingTime bool
	serverMetrics      *openmetrics.ServerMetrics // internal state

	tracing  bool
	logging  bool
	recovery bool
}

func NewGrpcServer(opts ...Option) *GrpcServer {
	g := &GrpcServer{
		timeout: 1 * time.Second,
		health:  health.NewServer(),
	}
	for _, opt := range opts {
		opt(g)
	}
	g.initGrpcServer()
	return g
}

func (s *GrpcServer) initGrpcServer() {
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	if s.customInterceptorOuterMost {
		logger.Infow("init custom interceptors to be the outer most wrapper around the real call")
		if len(s.unaryInterceptor) > 0 {
			unaryInterceptors = append(unaryInterceptors, s.unaryInterceptor...)
		}
		if len(s.streamInterceptor) > 0 {
			streamInterceptors = append(streamInterceptors, s.streamInterceptor...)
		}
	}

	if s.recovery {
		grpcRecoveryOpts := []grpcRecovery.Option{
			grpcRecovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
				md, ok := metadata.FromIncomingContext(ctx)
				if ok {
					logger.Errorw("grpcRecovery panic triggered", "grpc_metadata", md, "panic_info", p)
				} else {
					logger.Errorw("grpcRecovery panic triggered", "panic_info", p)
				}
				return fmt.Errorf("[grpcRecovery] context panic triggered: %v", p)
			}),
		}
		unaryInterceptors = append(unaryInterceptors, grpcRecovery.UnaryServerInterceptor(grpcRecoveryOpts...))
		streamInterceptors = append(streamInterceptors, grpcRecovery.StreamServerInterceptor(grpcRecoveryOpts...))
	}

	if s.tracing {
		unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, otelgrpc.StreamServerInterceptor())
		logger.Infow("grpc server tracing middleware enabled")
	}

	if s.metric {
		serverMetricsOptions := []openmetrics.ServerMetricsOption{
			openmetrics.WithServerCounterOptions(func(opts *prometheus.CounterOpts) {
				opts.Subsystem = s.metricSubSystem
			}),
		}
		if s.metricHandlingTime {
			serverMetricsOptions = append(serverMetricsOptions,
				openmetrics.WithServerHandlingTimeHistogram(func(opts *prometheus.HistogramOpts) {
					opts.Subsystem = s.metricSubSystem
				}))
		}
		prom := openmetrics.NewServerMetrics(serverMetricsOptions...)
		s.serverMetrics = prom
		unaryInterceptors = append(unaryInterceptors, openmetrics.UnaryServerInterceptor(prom))
		streamInterceptors = append(streamInterceptors, openmetrics.StreamServerInterceptor(prom))
		logger.Infow("grpc server metric middleware enabled")
	}

	if s.logging {
		codeToLevel := func(c codes.Code) logging.Level {
			if c == codes.InvalidArgument {
				return logging.WARNING
			}
			level := logging.DefaultServerCodeToLevel(c)
			return level
		}
		healthCheckMethod := fmt.Sprintf("/%s/Check", HEALTHCHECK_SERVICE)
		decider := logging.WithDecider(func(fullMethodName string, err error) logging.Decision {
			if err == nil && fullMethodName == healthCheckMethod {
				return logging.NoLogCall
			}
			return logging.LogStartAndFinishCall
		})
		opts := []logging.Option{
			decider,
			logging.WithLevels(codeToLevel),
		}
		unaryInterceptors = append(unaryInterceptors, logging.UnaryServerInterceptor(grpczap.InterceptorLogger(logger.GetZapLogger()), opts...))
		streamInterceptors = append(streamInterceptors, logging.StreamServerInterceptor(grpczap.InterceptorLogger(logger.GetZapLogger()), opts...))
		logger.Infow("grpc server logging middleware enabled")
	}

	// The first interceptor will be the outer most, while the last interceptor will be the inner most wrapper around the real call.
	// All unary interceptors added by this method will be chained.
	// we make the user interceptor the last appended, so it will be the inner most wrapper around the real call.
	if !s.customInterceptorOuterMost {
		logger.Infow("append custom interceptors to be the inner most wrapper around the real call")
		if len(s.unaryInterceptor) > 0 {
			unaryInterceptors = append(unaryInterceptors, s.unaryInterceptor...)
		}
		if len(s.streamInterceptor) > 0 {
			streamInterceptors = append(streamInterceptors, s.streamInterceptor...)
		}
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
	if len(s.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, s.grpcOpts...)
	}

	gs := grpc.NewServer(grpcOpts...)
	s.Server = gs

	s.health.SetServingStatus(HEALTHCHECK_SERVICE, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s.Server, s.health)
	reflection.Register(s.Server)
}

// Endpoint return a real address to registry endpoint.
// examples:
// grpc://127.0.0.1:9000?isSecure=false
func (s *GrpcServer) Endpoint() (*url.URL, error) {
	config := s.config
	address := fmt.Sprintf("%s:%s", config.Addr, config.Port)
	s.once.Do(func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(address, s.listener)
		if err != nil {
			lis.Close()
			s.err = err
			return
		}
		s.listener = lis
		s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

func (s *GrpcServer) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	logger.Infow("[gRPC] server started", "listen_addr", s.listener.Addr().String())

	if s.metric && s.serverMetrics != nil {
		s.serverMetrics.InitializeMetrics(s.Server)
		if err := s.serverMetrics.Register(prometheus.DefaultRegisterer); err != nil {
			logger.Errorw("[gRPC] server register prometheus metrics failed", "err", err)
		}
	}

	s.health.Resume()
	return s.Server.Serve(s.listener)
}

func (s *GrpcServer) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.health.Shutdown()
	logger.Info("[gRPC] server stopping")
	return nil
}
