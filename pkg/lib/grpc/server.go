package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

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
	log          *zap.Logger
	sigint       chan os.Signal
	server       *grpc.Server
	listenPort   int64
	registerFunc RegisterFunc
	logAny       bool
	reflection   bool
}

type RegisterFunc func(s *grpc.Server)

func (s *Server) RegisterServers(fn RegisterFunc) {
	s.registerFunc = fn
}

func (s *Server) ListenAndServe() error {
	g, gctx := errgroup.WithContext(s.context)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.listenPort))
	if nil != err {
		return err
	}

	g.Go(func() error {
		return s.run(l)
	})

	select {
	case <-s.sigint:
		s.log.Warn("Shutdown signal received")
		break
	case <-gctx.Done():
		s.log.Warn("Service stopped")
		break
	}

	signal.Stop(s.sigint)

	s.cancel()
	s.log.Sync()
	s.server.GracefulStop()

	return g.Wait()
}

func (s *Server) run(l net.Listener) error {
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(ep string, err error) bool {
			if err != nil {
				return true
			}

			if false == s.logAny || ep == "/grpc.health.v1.Health/Check" {
				return false
			}

			if s.logAny {
				return true
			}

			return false
		}),
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(s.log, zapOpts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(s.log, zapOpts...),
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

	s.log.Info(
		"gRPC server listening",
		zap.String("listen", l.Addr().String()),
	)

	err := s.server.Serve(l)
	if nil != err && grpc.ErrServerStopped != err {
		return err
	}

	return nil
}
