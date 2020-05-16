package grpc

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Service is the identification of the gRPC service
// Is recommended to use the dotted naming pattern. eg.: {name}.{namespace}
type Service struct {
	// Name uniquely identifies the service
	Name string
	// Version is a stamp of service current version
	Version *Version
}

// MarshalLogObject tells zap how to handle struct encoding
func (s Service) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", s.Name)
	enc.AddObject("version", s.Version)

	return nil
}

func New(svc Service, opts Options) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	fields := []zap.Field{
		zap.String("go", runtime.Version()),
		zap.Object("svc", svc),
	}

	if opts.Kubernetes != nil {
		fields = append(fields, zap.Object("k8s", opts.Kubernetes))
	}

	log, err := zap.NewProduction(zap.Fields(fields...))
	if nil != err {
		return nil, err
	}

	var listenPort int64 = 20020
	if opts.ListenPort > 0 {
		listenPort = opts.ListenPort
	}

	return &Server{
		context:    ctx,
		cancel:     cancel,
		log:        log,
		sigint:     sigint,
		listenPort: listenPort,
		logAny:     opts.LogAnyRequest,
		reflection: opts.ServerReflection,
	}, nil
}
