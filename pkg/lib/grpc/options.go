package grpc

type Options struct {
	// ListenPort of the gRPC server
	ListenPort int64
	// Kubernetes contains additional information of pod/node name
	Kubernetes *Kubernetes
	// LogAnyRequest forces to log all requests and not only on error
	LogAnyRequest bool
	// ServerReflection enables gRPC server reflection
	ServerReflection bool
}
