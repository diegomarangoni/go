package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *zap.Logger
	grpc   *grpc.Server
}

func (s *Server) ListenAndServe() error {
	return nil
}
