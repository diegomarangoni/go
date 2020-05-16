package grpc

import "go.uber.org/zap/zapcore"

// Version contains the binary stamping
type Version struct {
	GitCommit    string
	GitBranch    string
	BuildDate    string
	BuildVersion string
}

// MarshalLogObject tells zap how to handle struct encoding
func (v Version) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GitCommit", v.GitCommit)
	enc.AddString("GitBranch", v.GitBranch)
	enc.AddString("BuildDate", v.BuildDate)
	enc.AddString("BuildVersion", v.BuildVersion)

	return nil
}
