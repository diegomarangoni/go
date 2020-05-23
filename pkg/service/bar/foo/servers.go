package foo

import (
	pb "diegomarangoni.dev/go/pkg/pb/service/bar/foo/v1"
	"google.golang.org/grpc"
)

func RegisterServers(s *grpc.Server) {
	pb.RegisterAcmeServer(s, &Acme{})
}
