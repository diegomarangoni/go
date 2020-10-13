package discovery

import (
	"context"
	"net"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/clientv3/naming"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	grpc_naming "google.golang.org/grpc/naming"
)

type Lease struct {
	clientv3.Lease
	ID       clientv3.LeaseID
	Revision int64
}

type Discovery struct {
	context  context.Context
	etcd     *clientv3.Client
	address  string
	logger   *zap.Logger
	service  string
	resolver *naming.GRPCResolver
	lease    *Lease
}

var (
	defaultTimeout time.Duration = 5 * time.Second
	defaultTTL     int64         = 60
)

func (d *Discovery) Run() error {
	ctx, cancel := context.WithTimeout(d.context, defaultTimeout)
	err := d.newLease(ctx)
	cancel()
	if nil != err {
		return err
	}

	defer d.lease.Close()

	ctx, cancel = context.WithTimeout(d.context, defaultTimeout)
	err = d.register(ctx)
	cancel()
	if nil != err {
		return err
	}

	ctx, cancel = context.WithCancel(d.context)
	err = d.keepAlive(ctx)
	cancel()
	if nil != err {
		return err
	}

	ctx, cancel = context.WithTimeout(d.context, defaultTimeout)
	err = d.deregister(context.Background())
	cancel()
	if nil != err {
		return err
	}

	return nil
}

func (d *Discovery) BestEffortRun() {
	eg, egctx := errgroup.WithContext(d.context)

	eg.Go(func() error {
		err := d.Run()
		d.logger.Error("service discovery failed", zap.Error(err))
		return err
	})

	select {
	case <-d.context.Done():
		break
	case <-egctx.Done():
		time.Sleep(defaultTimeout)
		d.BestEffortRun()
	}
}

func (d *Discovery) Resolver() *naming.GRPCResolver {
	return d.resolver
}

func (d *Discovery) keepAlive(ctx context.Context) error {
	keepAlive, err := d.lease.KeepAlive(ctx, d.lease.ID)
	if nil != err {
		return err
	}

	for {
		select {
		case resp := <-keepAlive:
			if nil == resp {
				return nil
			}

			if resp.Revision > d.lease.Revision {
				// something changed, register again to override any changes
				err = d.register(ctx)
				if nil != err {
					return err
				}
				d.lease.Revision = resp.Revision + 1
			}
		}
	}
}

func (d *Discovery) newLease(ctx context.Context) error {
	lease := clientv3.NewLease(d.etcd)

	leaseGrant, err := lease.Grant(ctx, defaultTTL)
	if nil != err {
		d.logger.Error("failed to obtain lease", zap.Error(err))
		return err
	}

	d.lease = &Lease{
		Lease:    lease,
		ID:       leaseGrant.ID,
		Revision: leaseGrant.Revision,
	}

	return nil
}

func (d *Discovery) register(ctx context.Context) error {
	op := grpc_naming.Update{Op: grpc_naming.Add, Addr: d.address}

	err := d.resolver.Update(ctx, d.service, op, clientv3.WithLease(d.lease.ID))
	if nil != err {
		d.logger.Error("failed to register", zap.Error(err))
		return err
	}

	return nil
}

func (d *Discovery) deregister(ctx context.Context) error {
	op := grpc_naming.Update{Op: grpc_naming.Delete, Addr: d.address}

	err := d.resolver.Update(ctx, d.service, op)
	if nil != err {
		d.logger.Error("failed to deregister", zap.Error(err))
		return err
	}

	return nil
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

func NewService(ctx context.Context, svc Service, opts *Options) (*Discovery, error) {
	if nil == ctx {
		ctx = context.Background()
	}

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

	ip, err := guessHostIP(svc.Address)
	if nil != err {
		return nil, err
	}

	return &Discovery{
		context:  ctx,
		etcd:     etcd,
		address:  ip,
		logger:   opts.Logger,
		service:  svc.Name,
		resolver: &naming.GRPCResolver{Client: etcd},
	}, nil
}
