package grpc

import (
	"go.uber.org/zap"
)

type Options struct {
	Logger *zap.Logger
	// ListenPort of the gRPC server
	ListenPort int64
	// Kubernetes contains additional information of pod/node name
	Kubernetes *Kubernetes
	// LogAllRequests and not only the ones that returned error
	LogAllRequests bool
	// ServerReflection enables gRPC server reflection
	ServerReflection bool
}
