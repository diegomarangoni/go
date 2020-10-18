package cito

import "go.uber.org/zap/zapcore"

// Kubernetes add additional information to logging
type Kubernetes struct {
	// Pod represents the k8s pod name
	Pod string
	// Node represents the k8s node name where current service pod is running
	Node string
}

// MarshalLogObject tells zap how to handle struct encoding
func (k Kubernetes) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("pod", k.Pod)
	enc.AddString("node", k.Node)

	return nil
}
