package cito

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service interface {
	Namespace() string
	Name() string
}

type ServiceWantsLogger interface {
	SetLogger(*zap.Logger)
}

type ServiceWantsClient interface {
	SetClient(*Client)
}

type ServiceWantsServer interface {
	SetServer(*grpc.Server)
}
