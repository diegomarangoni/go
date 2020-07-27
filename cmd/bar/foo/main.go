package main

import (
	"log"

	"diegomarangoni.dev/go/pkg/lib/discovery"
	"diegomarangoni.dev/go/pkg/lib/grpc"
	pb "diegomarangoni.dev/go/pkg/service/bar/foo"
	"diegomarangoni.dev/typenv"
	"go.uber.org/zap"
)

var (
	GitCommit    string
	GitBranch    string
	BuildDate    string
	BuildVersion string
)

func init() {
	typenv.SetGlobalDefault(
		typenv.E(typenv.Bool, "DEBUG", true),
	)
}

func main() {
	logger, err := zap.NewProduction()
	if nil != err {
		log.Panicf("failed to create a logger instance: %v", err)
	}

	svc := grpc.Service{
		Name: "foo.bar",
		Version: &grpc.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}

	opts := grpc.Options{
		Logger:           logger,
		ListenPort:       typenv.Int64("LISTEN_PORT", 20020),
		LogAllRequests:   typenv.Bool("LOG_ALL_REQUESTS", typenv.Bool("DEBUG")),
		ServerReflection: typenv.Bool("SERVER_REFLECTION", typenv.Bool("DEBUG")),
	}

	if "" != typenv.String("POD_NAME") && "" != typenv.String("NODE_NAME") {
		opts.Kubernetes = &grpc.Kubernetes{
			Pod:  typenv.String("POD_NAME"),
			Node: typenv.String("NODE_NAME"),
		}
	}

	srv, err := grpc.NewServer(svc, opts)
	if nil != err {
		logger.Panic("failed to create grpc server", zap.Error(err))
	}

	srv.RegisterServersFunc(pb.RegisterServers)

	etcd := &discovery.Etcd{
		Endpoints: []string{typenv.String("ETCD_SERVICE", "http://cluster1.etcd.svc.cluster.local:2379")},
	}
	discv, err := discovery.NewService(srv.Context(),
		discovery.Service{Name: svc.Name, Address: srv.Addr()},
		discovery.Options{Logger: logger, Etcd: etcd})
	if nil != err {
		logger.Panic("service discovery failed", zap.Error(err))
	}
	go discv.BestEffortRun()

	err = srv.ListenAndServe()
	if nil != err {
		logger.Panic("grpc server failed", zap.Error(err))
	}
}
