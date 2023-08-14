package transports

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/costa92/errors"
	"github.com/costa92/logger"
	"github.com/segmentio/ksuid"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"

	"github.com/kart-io/kart"
)

// ID returns app instance id.
func (gs *GenericAPIServer) ID() string { return "" }

// Name returns service name.
func (gs *GenericAPIServer) Name() string { return "" }

// Version returns app version.
func (gs *GenericAPIServer) Version() string { return "" }

type GenericAPIServer struct {
	opts   options
	ctx    context.Context
	cancel func()
}

// NewGenericAPIServer 实例化
func NewGenericAPIServer(opts ...Option) *GenericAPIServer {
	opt := options{
		ctx:           context.Background(),
		sigs:          []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		handleSignals: true,
	}
	opt.id = ksuid.New().String()
	for _, o := range opts {
		o(&opt)
	}
	ctx, cancel := context.WithCancel(opt.ctx)
	return &GenericAPIServer{
		ctx:    ctx,
		cancel: cancel,
		opts:   opt,
	}
}

// Run 开始运行
func (gs *GenericAPIServer) Run() error {
	//nolint: forbidigo
	if _, err := maxprocs.Set(maxprocs.Logger(logger.Infof)); err != nil {
		return err
	}
	ctx := NewContext(gs.ctx, gs)
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range gs.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()
	c := make(chan os.Signal, 1)
	// warning: you need manually call App.Stop() to stop the application if you set handleSignals to false
	if gs.opts.handleSignals {
		signal.Notify(c, gs.opts.sigs...)
	}

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return gs.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (gs *GenericAPIServer) Stop() error {
	if gs.cancel != nil {
		gs.cancel()
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s kart.AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s kart.AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(kart.AppInfo)
	return
}
