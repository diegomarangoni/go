package foo

import (
	"context"

	pb "diegomarangoni.dev/go/pkg/pb/service/bar/foo/v1"
)

type Acme struct {
}

func (h Acme) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{}, nil
}
