package grpc

import (
	"go.uber.org/zap/zapcore"
)

// Service is the identification of the gRPC service
// Is recommended to use the dotted naming pattern. eg.: {name}.{namespace}
type Service struct {
	// Name uniquely identifies the service
	Name string
	// Version is a stamp of service current version
	Version *Version
}

// MarshalLogObject tells zap how to handle struct encoding
func (s Service) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", s.Name)
	enc.AddObject("version", s.Version)

	return nil
}
