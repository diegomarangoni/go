package bar

import (
	pb "diegomarangoni.dev/go/pkg/pb/service/foo/bar/v1"
	"google.golang.org/grpc"
)

func RegisterServers(s *grpc.Server) {
	pb.RegisterAcmeServer(s, &Acme{})
}
