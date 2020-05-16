package name

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pb "diegomarangoni.dev/go/pkg/pb/service/namespace/name/v1"
)

type Example struct {
}

func (h Example) HelloWorld(ctx context.Context, req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	if "world" == strings.ToLower(req.Name) {
		return nil, errors.New("You can't say hello to entire world")
	}

	return &pb.HelloWorldResponse{
		Content: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
