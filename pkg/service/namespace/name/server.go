package name

import (
	pb "diegomarangoni.dev/go/pkg/pb/service/namespace/name/v1"
	"google.golang.org/grpc"
)

func RegisterServer(s *grpc.Server) {
	pb.RegisterExampleServer(s, &Example{})
}
