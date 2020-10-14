package main

import (
	"log"

	"diegomarangoni.dev/go/pkg/lib/discovery"
	"diegomarangoni.dev/go/pkg/lib/grpc"
	name_namespace "diegomarangoni.dev/go/pkg/service/namespace/name"
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

	instance := grpc.Instance{
		Service: name_namespace.Service{},
		Version: &grpc.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}

	serverOpts := grpc.ServerOptions{
		Logger:             logger,
		ListenPort:         typenv.Int64("LISTEN_PORT", 20020),
		LogAllRequests:     typenv.Bool("LOG_ALL_REQUESTS", typenv.Bool("DEBUG")),
		ServerReflection:   typenv.Bool("SERVER_REFLECTION", typenv.Bool("DEBUG")),
		RegisterServerFunc: name_namespace.RegisterServer,
	}
	if "" != typenv.String("POD_NAME") && "" != typenv.String("NODE_NAME") {
		serverOpts.Kubernetes = &grpc.Kubernetes{
			Pod:  typenv.String("POD_NAME"),
			Node: typenv.String("NODE_NAME"),
		}
	}
	server, err := grpc.NewServer(instance, serverOpts)
	if nil != err {
		logger.Panic("failed to create grpc server", zap.Error(err))
	}

	clientOpts := grpc.ClientOptions{
		Logger:         logger,
		LogAllRequests: typenv.Bool("LOG_ALL_REQUESTS", typenv.Bool("DEBUG")),
	}
	if "" != typenv.String("POD_NAME") && "" != typenv.String("NODE_NAME") {
		clientOpts.Kubernetes = &grpc.Kubernetes{
			Pod:  typenv.String("POD_NAME"),
			Node: typenv.String("NODE_NAME"),
		}
	}
	client, err := grpc.NewClient(instance, clientOpts)
	if nil != err {
		logger.Panic("failed to create grpc client", zap.Error(err))
	}
	clientDiscovery, err := discovery.NewClient(client)
	if nil != err {
		logger.Panic("failed to create discovery client", zap.Error(err))
	}
	name_namespace.RegisterClient(clientDiscovery)

	go func() {
		etcd := &discovery.Etcd{
			Endpoints: []string{typenv.String("ETCD_SERVICE", "http://localhost:2379")},
		}
		service := discovery.Instance{
			Service: instance.Service,
			Address: server.Addr(),
		}
		options := &discovery.Options{
			Logger: logger,
			Etcd:   etcd,
		}
		serviceDiscovery, err := discovery.NewService(server.Context(), service, options)
		if nil != err {
			logger.Panic("service discovery failed", zap.Error(err))
		}
		go serviceDiscovery.BestEffortRun()
	}()

	err = server.ListenAndServe()
	if nil != err {
		logger.Panic("grpc server failed", zap.Error(err))
	}
}
