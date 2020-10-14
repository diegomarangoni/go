package discovery

import (
	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/clientv3/naming"
	"google.golang.org/grpc"
)

type GRPC interface {
	Dial(*naming.GRPCResolver) *grpc.ClientConn
}

type ClientDiscovery struct {
	grpc     GRPC
	etcd     *clientv3.Client
	resolver *naming.GRPCResolver
}

func (d *ClientDiscovery) Conn(service Service) *grpc.ClientConn {
	return d.grpc.Dial(d.resolver)
}

func NewClient(grpc GRPC) (*ClientDiscovery, error) {
	return &ClientDiscovery{
		grpc:     grpc,
		etcd:     nil,
		resolver: nil,
	}, nil
}
