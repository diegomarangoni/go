package main

import (
	"log"

	"diegomarangoni.dev/go/pkg/lib/env"
	"diegomarangoni.dev/go/pkg/lib/grpc"
	pb "diegomarangoni.dev/go/pkg/service/namespace/name"
)

var (
	GitCommit    string
	GitBranch    string
	BuildDate    string
	BuildVersion string
)

func init() {
	env.SetGlobalDefault(
		env.E(env.Int64, "LISTEN_PORT", 20020),
	)
}

func main() {
	svc := grpc.Service{
		Name: "example.namespace",
		Version: &grpc.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}

	debug := env.Bool("DEBUG", true)

	opts := grpc.Options{
		ListenPort:       env.Int64("LISTEN_PORT"),
		LogAllRequests:   env.Bool("LOG_ALL_REQUESTS", debug),
		ServerReflection: env.Bool("SERVER_REFLECTION", debug),
	}

	if "" != env.String("POD_NAME") && "" != env.String("NODE_NAME") {
		opts.Kubernetes = &grpc.Kubernetes{
			Pod:  env.String("POD_NAME"),
			Node: env.String("NODE_NAME"),
		}
	}

	srv, err := grpc.New(svc, opts)
	if nil != err {
		log.Fatalf("failed to create grpc server: %v", err)
	}

	srv.RegisterServers(pb.RegisterServers)

	err = srv.ListenAndServe()
	if nil != err {
		log.Fatalf("server failed: %v", err)
	}
}
