package grpc

import (
	"runtime"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.etcd.io/etcd/v3/clientv3/naming"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	logger *zap.Logger
	logAll bool
}

func (c *Client) Conn(service string, resolver *naming.GRPCResolver) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			grpc_zap.UnaryClientInterceptor(c.logger),
		),
		grpc.WithChainStreamInterceptor(
			grpc_zap.StreamClientInterceptor(c.logger),
		),
	}

	opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(resolver)))

	return grpc.Dial(service, opts...)
}

func NewClient(inst Instance, opts ClientOptions) (*Client, error) {
	fields := []zap.Field{
		zap.String("go", runtime.Version()),
		zap.Object("instance", inst),
	}
	if opts.Kubernetes != nil {
		fields = append(fields, zap.Object("k8s", opts.Kubernetes))
	}
	logger := opts.Logger
	if opts.Logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if nil != err {
			return nil, err
		}
	}
	logger = logger.WithOptions(zap.Fields(fields...))

	return &Client{
		logger: logger,
		logAll: opts.LogAllRequests,
	}, nil
}
