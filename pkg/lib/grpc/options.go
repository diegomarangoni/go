package grpc

import (
	"go.uber.org/zap"
)

type ClientOptions struct {
	// Logger intance
	Logger *zap.Logger
	// Kubernetes contains additional information of pod/node name
	Kubernetes *Kubernetes
	// LogAllRequests and not only the ones that returned error
	LogAllRequests bool
}

type ServerOptions struct {
	// Logger intance
	Logger *zap.Logger
	// ListenPort of the gRPC server
	ListenPort int64
	// Kubernetes contains additional information of pod/node name
	Kubernetes *Kubernetes
	// LogAllRequests and not only the ones that returned error
	LogAllRequests bool
	// ServerReflection enables gRPC server reflection
	ServerReflection bool
	// RegisterServerFunc is called passing grpc server instance
	RegisterServerFunc RegisterFunc
}
