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

func (h Example) HelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	if "world" == strings.ToLower(r.Name) {
		return nil, errors.New("You can't say hello to entire world")
	}

	return &pb.HelloWorldResponse{
		Content: fmt.Sprintf("Hello %s", r.Name),
	}, nil
}
