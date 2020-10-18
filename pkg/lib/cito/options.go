package cito

type ClientOptions struct {
	// Address contains host:port combination of your project
	Address string
	// LogAllRequests and not only the ones that returned error
	LogAllRequests bool
}

type ServerOptions struct {
	// ListenPort of the gRPC server
	ListenPort int64
	// LogAllRequests and not only the ones that returned error
	LogAllRequests bool
	// Reflection enables gRPC server reflection
	Reflection bool
}
