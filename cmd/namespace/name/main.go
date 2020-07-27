package main

import (
	"log"

	"diegomarangoni.dev/go/pkg/lib/discovery"
	"diegomarangoni.dev/go/pkg/lib/grpc"
	pb "diegomarangoni.dev/go/pkg/service/namespace/name"
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
		typenv.E(typenv.Int64, "LISTEN_PORT", 20020),
	)
}

func main() {
	logger, err := zap.NewProduction()
	if nil != err {
		log.Panicf("failed to create a logger instance: %v", err)
	}

	debug := typenv.Bool("DEBUG", true)

	svc := grpc.Service{
		Name: "name.namespace",
		Version: &grpc.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}
	opts := grpc.Options{
		Logger:           logger,
		ListenPort:       typenv.Int64("LISTEN_PORT"),
		LogAllRequests:   typenv.Bool("LOG_ALL_REQUESTS", debug),
		ServerReflection: typenv.Bool("SERVER_REFLECTION", debug),
	}
	if "" != typenv.String("POD_NAME") && "" != typenv.String("NODE_NAME") {
		opts.Kubernetes = &grpc.Kubernetes{
			Pod:  typenv.String("POD_NAME"),
			Node: typenv.String("NODE_NAME"),
		}
	}
	srv, err := grpc.New(svc, opts)
	if nil != err {
		logger.Panic("failed to create grpc server", zap.Error(err))
	}
	srv.RegisterServers(pb.RegisterServers)

	svcdisc := discovery.NewService(srv)
	err = svcdisc.Register()
	if nil != err {
		logger.Panic("unable to register", zap.Error(err))
	}
	defer func() {
		err := svcdisc.Unregister()
		if nil != err {
			logger.Panic("unable to deregister", zap.Error(err))
		}
	}()

	err = srv.ListenAndServe()
	if nil != err {
		logger.Panic("grpc server failed", zap.Error(err))
	}
}
