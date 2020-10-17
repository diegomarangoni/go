package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	context      context.Context
	cancel       context.CancelFunc
	signal       chan os.Signal
	listen       net.Listener
	logger       *zap.Logger
	registerFunc RegisterFunc
	logAll       bool
	reflection   bool
	server       *grpc.Server
}

func NewServer(instance Instance, opts ServerOptions) (*Server, error) {
	chansig := make(chan os.Signal, 1)
	signal.Notify(chansig, os.Interrupt, syscall.SIGTERM)

	var listenPort int64 = 20020
	if opts.ListenPort > 0 {
		listenPort = opts.ListenPort
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if nil != err {
		return nil, err
	}

	fields := []zap.Field{
		zap.String("go", runtime.Version()),
		zap.Object("instance", instance),
	}
	if opts.Kubernetes != nil {
		fields = append(fields, zap.Object("k8s", opts.Kubernetes))
	}
	logger := opts.Logger
	if opts.Logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if nil != err {
			return nil, err
		}
	}
	logger = logger.WithOptions(zap.Fields(fields...))

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		context:      ctx,
		cancel:       cancel,
		signal:       chansig,
		listen:       listen,
		logger:       logger,
		registerFunc: opts.RegisterServerFunc,
		logAll:       opts.LogAllRequests,
		reflection:   opts.ServerReflection,
	}, nil
}

type RegisterFunc func(s *grpc.Server)

func (s *Server) Addr() net.Addr {
	if nil == s.listen {
		return nil
	}
	return s.listen.Addr()
}

func (s *Server) Context() context.Context {
	return s.context
}

func (s *Server) ListenAndServe() error {
	eg, egctx := errgroup.WithContext(s.context)

	eg.Go(func() error {
		return s.run()
	})

	select {
	case <-s.signal:
		s.logger.Warn("Shutdown signal received")
		break
	case <-egctx.Done():
		s.logger.Warn("Service stopped")
		break
	}

	s.cancel()
	signal.Stop(s.signal)
	s.server.GracefulStop()
	s.logger.Sync()

	return eg.Wait()
}

func (s *Server) run() error {
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(ep string, err error) bool {
			if err != nil {
				return true
			}

			if false == s.logAll || ep == "/grpc.health.v1.Health/Check" {
				return false
			}

			if s.logAll {
				return true
			}

			return false
		}),
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(s.logger, zapOpts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(s.logger, zapOpts...),
		)),
	}

	s.server = grpc.NewServer(opts...)

	if nil != s.registerFunc {
		s.registerFunc(s.server)
	}

	if s.reflection {
		reflection.Register(s.server)
	}

	hc := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.server, hc)
	hc.SetServingStatus("server", grpc_health_v1.HealthCheckResponse_SERVING)
	defer hc.SetServingStatus("server", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	s.logger.Info(
		"gRPC server listening",
		zap.String("listen", s.listen.Addr().String()),
	)

	err := s.server.Serve(s.listen)
	if nil != err && grpc.ErrServerStopped != err {
		return err
	}

	return nil
}
