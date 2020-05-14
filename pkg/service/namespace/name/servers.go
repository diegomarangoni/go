package name

import (
	name "github.com/diegomarangoni/gomonorepo/pkg/pb/service/namespace/name/v1"
	"google.golang.org/grpc"
)

func RegisterServers(s *grpc.Server) {
	name.RegisterExampleServer(s, &Example{})
}
