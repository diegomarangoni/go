package bar

import (
	"context"

	pb "diegomarangoni.dev/go/pkg/pb/service/foo/bar/v1"
)

type Acme struct {
}

func (h Acme) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{}, nil
}
