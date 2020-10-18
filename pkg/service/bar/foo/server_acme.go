package foo

import (
	"context"

	"diegomarangoni.dev/go/pkg/lib/cito"
	pb "diegomarangoni.dev/go/pkg/pb/service/bar/foo/v1"
	pb_name_namespace "diegomarangoni.dev/go/pkg/pb/service/namespace/name/v1"
	name_namespace "diegomarangoni.dev/go/pkg/service/namespace/name"
	"go.uber.org/zap"
)

type Acme struct {
	pb.UnimplementedAcmeServer

	client *cito.Client
	logger *zap.Logger
}

func (h Acme) Ping(ctx context.Context, r *pb.PingRequest) (*pb.PongResponse, error) {
	nameNamespaceConn, err := h.client.Conn(name_namespace.Service{})
	if nil != err {
		h.logger.Error("failed to stablish connection", zap.Error(err))
		return nil, err
	}

	nameNamespaceClient := pb_name_namespace.NewExampleClient(nameNamespaceConn)

	req := &pb_name_namespace.HelloWorldRequest{
		Name: "Jhon Doe",
	}

	res, err := nameNamespaceClient.HelloWorld(ctx, req)
	if nil != err {
		h.logger.Error("failed to say hello", zap.Error(err))
		return nil, err
	}

	h.logger.Info("saying hello", zap.String("HelloWorldResponse", res.Content))

	return &pb.PongResponse{}, nil
}
