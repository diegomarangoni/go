package discovery

import (
	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/clientv3/naming"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPC interface {
	Conn(string, *naming.GRPCResolver) (*grpc.ClientConn, error)
}

type ClientDiscovery struct {
	logger   *zap.Logger
	grpc     GRPC
	etcd     *clientv3.Client
	resolver *naming.GRPCResolver
}

func (d *ClientDiscovery) Conn(service Service) *grpc.ClientConn {
	conn, err := d.grpc.Conn(service.String(), d.resolver)
	if nil != err {
		d.logger.Error("failed to create connection", zap.Error(err))
		return nil
	}

	return conn
}

func NewClient(grpc GRPC, opts *Options) (*ClientDiscovery, error) {
	if nil == opts {
		opts = &Options{}
	}

	if nil == opts.Logger {
		var err error
		opts.Logger, err = zap.NewProduction()
		if nil != err {
			return nil, err
		}
	}

	if nil == opts.Etcd {
		opts.Etcd = &Etcd{
			Endpoints: []string{"http://localhost:2379"},
		}
	}

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints: opts.Etcd.Endpoints,
	})
	if nil != err {
		return nil, err
	}

	return &ClientDiscovery{
		logger:   opts.Logger,
		grpc:     grpc,
		etcd:     etcd,
		resolver: &naming.GRPCResolver{Client: etcd},
	}, nil
}
