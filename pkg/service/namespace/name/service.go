package name

import (
	"diegomarangoni.dev/go/pkg/lib/cito"
	pb "diegomarangoni.dev/go/pkg/pb/service/namespace/name/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	client *cito.Client
	logger *zap.Logger
}

func (Service) Name() string {
	return "name"
}

func (Service) Namespace() string {
	return "namespace"
}

func (svc *Service) SetLogger(l *zap.Logger) {
	svc.logger = l
}

func (svc *Service) SetServer(s *grpc.Server) {
	pb.RegisterExampleServer(s, &Example{})
}

func (svc *Service) SetClient(c *cito.Client) {
	svc.client = c
}
