package foo

import (
	"context"
	"fmt"

	pb "diegomarangoni.dev/go/pkg/pb/service/bar/foo/v1"
	pb_name_namespace "diegomarangoni.dev/go/pkg/pb/service/namespace/name/v1"
	name_namespace "diegomarangoni.dev/go/pkg/service/namespace/name"
)

type Acme struct {
}

func (h Acme) Ping(ctx context.Context, r *pb.PingRequest) (*pb.PongResponse, error) {
	nameNamespaceClient := pb_name_namespace.NewExampleClient(client.Conn(name_namespace.Service{}))

	req := &pb_name_namespace.HelloWorldRequest{
		Name: "Jhon Doe",
	}

	res, err := nameNamespaceClient.HelloWorld(ctx, req)
	if nil != err {
		fmt.Printf("failed to say hello: %s", err.Error())
		return nil, err
	}

	fmt.Printf("saying hello: %s", res.Content)

	return &pb.PongResponse{}, nil
}
