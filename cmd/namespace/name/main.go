package main

import (
	"log"

	"diegomarangoni.dev/go/pkg/lib/cito"
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

	instance := cito.Instance{
		Service: &name_namespace.Service{},
		Version: &cito.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}

	var k8s *cito.Kubernetes
	if "" != typenv.String("POD_NAME") && "" != typenv.String("NODE_NAME") {
		k8s = &cito.Kubernetes{
			Pod:  typenv.String("POD_NAME"),
			Node: typenv.String("NODE_NAME"),
		}
	}

	clientOpts := &cito.ClientOptions{
		Address:        typenv.String("CITO_ENDPOINT", "myproject.cito.dev:443"),
		LogAllRequests: typenv.Bool("LOG_ALL_REQUESTS", typenv.Bool("DEBUG")),
	}

	serverOpts := &cito.ServerOptions{
		ListenPort:     typenv.Int64("LISTEN_PORT", 20020),
		LogAllRequests: typenv.Bool("LOG_ALL_REQUESTS", typenv.Bool("DEBUG")),
		Reflection:     typenv.Bool("SERVER_REFLECTION", typenv.Bool("DEBUG")),
	}

	service, err := cito.New(instance, k8s, serverOpts, clientOpts)
	if nil != err {
		logger.Panic("unable to create service", zap.Error(err))
	}

	err = service.Run()
	if nil != err {
		logger.Panic("service server failed", zap.Error(err))
	}
}
