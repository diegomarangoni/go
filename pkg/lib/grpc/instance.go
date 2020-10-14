package grpc

import (
	"go.uber.org/zap/zapcore"
)

// Instance is the identification of the gRPC service
// Is recommended to use the dotted naming pattern. eg.: {name}.{namespace}
type Instance struct {
	// Service uniquely identifies the service
	Service Service
	// Version is a stamp of service current version
	Version *Version
}

// MarshalLogObject tells zap how to handle struct encoding
func (s Instance) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("service", s.Service.Name())
	enc.AddObject("version", s.Version)

	return nil
}
