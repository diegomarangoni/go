package cito

import (
	"go.uber.org/zap/zapcore"
)

// Instance is the identification of the gRPC service
type Instance struct {
	// Service uniquely identifies the service
	Service Service
	// Version is a stamp of service current version
	Version *Version
}

// MarshalLogObject tells zap how to handle struct encoding
func (s Instance) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", s.Service.Name())
	enc.AddString("namespace", s.Service.Namespace())
	enc.AddObject("version", s.Version)

	return nil
}
