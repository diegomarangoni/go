package cito

import (
	"crypto/tls"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	logger  *zap.Logger
	address string
	logAll  bool
}

func newClient(instance Instance, logger *zap.Logger, opts ClientOptions) (*Client, error) {
	return &Client{
		logger:  logger,
		address: opts.Address,
		logAll:  opts.LogAllRequests,
	}, nil
}

func (c *Client) Conn(svc Service) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			grpc_zap.UnaryClientInterceptor(c.logger),
		),
		grpc.WithChainStreamInterceptor(
			grpc_zap.StreamClientInterceptor(c.logger),
		),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}

	return grpc.Dial(c.address, opts...)
}
