package grpc

import (
	"runtime"

	"go.etcd.io/etcd/v3/clientv3/naming"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	logger *zap.Logger
	logAll bool
}

func (*Client) Dial(resolver *naming.GRPCResolver) *grpc.ClientConn {
	return nil
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
