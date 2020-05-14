package grpc

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"
)

type Service struct {
	Name    string
	Port    int64
	Version *Version
}

func New(svc Service, opts ...Option) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	log, err := zap.NewProduction(zap.Fields(
		zap.String("go", runtime.Version()),
		zap.String("service", svc.Name),
		zap.Object("version", svc.Version),
	))
	if nil != err {
		return nil, err
	}

	return &Server{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
	}, nil
}
