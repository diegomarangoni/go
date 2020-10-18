package cito

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
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
	signal chan os.Signal
	listen net.Listener
	logger *zap.Logger
	server *grpc.Server
}

func newServer(instance Instance, logger *zap.Logger, opts ServerOptions) (*Server, error) {
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

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(ep string, err error) bool {
			if err != nil {
				return true
			}

			if false == opts.LogAllRequests || ep == "/grpc.health.v1.Health/Check" {
				return false
			}

			if opts.LogAllRequests {
				return true
			}

			return false
		}),
	}

	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger, zapOpts...),
		)),
	}

	server := grpc.NewServer(grpcOpts...)

	if opts.Reflection {
		reflection.Register(server)
	}

	return &Server{
		signal: chansig,
		listen: listen,
		logger: logger,
		server: server,
	}, nil
}

func (s *Server) listenAndServe() error {
	ctx, cancel := context.WithCancel(context.Background())

	eg, egctx := errgroup.WithContext(ctx)

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

	cancel()
	signal.Stop(s.signal)
	s.server.GracefulStop()
	s.logger.Sync()

	return eg.Wait()
}

func (s *Server) run() error {
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
