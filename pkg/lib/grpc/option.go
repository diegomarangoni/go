package grpc

import (
	"google.golang.org/grpc"
)

type Option func(*Server) error

type RegisterServersFunc func(s *grpc.Server)

func SetRegisterServersFunc(fn RegisterServersFunc) Option {
	return func(srv *Server) error {
		fn(srv.grpc)

		return nil
	}
}
