package main

import (
	"log"

	"github.com/diegomarangoni/gomonorepo/pkg/lib/grpc"
	"github.com/diegomarangoni/gomonorepo/pkg/service/namespace/name"
)

var (
	GitCommit    string
	GitBranch    string
	BuildDate    string
	BuildVersion string
)

func main() {
	svc := grpc.Service{
		Name: "namespace.example",
		Port: 20021,
		Version: &grpc.Version{
			GitCommit:    GitCommit,
			GitBranch:    GitBranch,
			BuildDate:    BuildDate,
			BuildVersion: BuildVersion,
		},
	}

	opts := []grpc.Option{
		grpc.SetRegisterServersFunc(name.RegisterServers),
	}

	srv, err := grpc.New(svc, opts...)
	if nil != err {
		log.Fatal(err)
	}

	err = srv.ListenAndServe()
	if nil != err {
		log.Fatal(err)
	}
}
