package discovery

import (
	"context"
	"net"

	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/clientv3/naming"
	"go.uber.org/zap"
)

type Discovery struct {
	etcd *clientv3.Client
	ip   string
}

func (d *Discovery) Run() error {
	return nil
}

func (d *Discovery) BestEffortRun() {
	_ = &naming.GRPCResolver{Client: d.etcd}

	// var lease *Lease

	// func() {
	// 	ctx, end := context.WithTimeout(s.context, 5*time.Second)
	// 	defer end()

	// 	lease, err = s.register(ctx, advertise)
	// 	if nil != err {
	// 		s.log.Error("Failed to register service", zap.Error(err))
	// 	}
	// }()

	// g.Go(func() error {
	// 	return s.runMetricsServer(metrics)
	// })

	// g.Go(func() error {
	// 	return s.runServiceServer(service)
	// })

	// g.Go(func() error {
	// 	return s.keepAlive(lease)
	// })
}

type Service struct {
	Name    string
	Address net.Addr
}

type Options struct {
	Logger *zap.Logger
	Etcd   *Etcd
}

type Etcd struct {
	Endpoints []string
}

func NewService(ctx context.Context, svc Service, opts Options) (*Discovery, error) {
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

	ip, err := guessHostIP(svc.Address.(*net.TCPAddr).IP.String())
	if nil != err {
		return nil, err
	}

	return &Discovery{
		etcd: etcd,
		ip:   ip,
	}, nil
}
