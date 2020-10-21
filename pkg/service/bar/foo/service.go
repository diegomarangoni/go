package foo

import (
	"diegomarangoni.dev/go/pkg/lib/cito"
	pb "diegomarangoni.dev/go/pkg/pb/service/bar/foo/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewCitoService() cito.Service {
	return &Service{}
}

type Service struct {
	client *cito.Client
	logger *zap.Logger
}

func (Service) Namespace() string {
	return "bar"
}

func (Service) Name() string {
	return "foo"
}

func (svc *Service) SetLogger(l *zap.Logger) {
	svc.logger = l
}

func (svc *Service) SetClient(c *cito.Client) {
	svc.client = c
}

func (svc *Service) SetServer(s *grpc.Server) {
	pb.RegisterAcmeServer(s, &Acme{
		client: svc.client,
		logger: svc.logger,
	})
}
