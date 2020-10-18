package cito

import (
	"runtime"

	"go.uber.org/zap"
)

type Cito struct {
	service Service
	logger  *zap.Logger
	client  *Client
	server  *Server
}

func New(instance Instance, k8s *Kubernetes, serverOpts *ServerOptions, clientOpts *ClientOptions) (*Cito, error) {
	logger, err := zap.NewProduction()
	if nil != err {
		return nil, err
	}
	fields := []zap.Field{
		zap.String("go", runtime.Version()),
		zap.Object("instance", instance),
	}
	if k8s != nil {
		fields = append(fields, zap.Object("k8s", k8s))
	}
	logger = logger.WithOptions(zap.Fields(fields...))

	var client *Client
	if nil != clientOpts {
		client, err = newClient(instance, logger, *clientOpts)
		if nil != err {
			return nil, err
		}
	}

	var server *Server
	if nil != serverOpts {
		server, err = newServer(instance, logger, *serverOpts)
		if nil != err {
			return nil, err
		}
	}

	return &Cito{
		service: instance.Service,
		logger:  logger,
		client:  client,
		server:  server,
	}, nil
}

func (eg *Cito) Run() error {
	serviceLogger, wantsLogger := eg.service.(ServiceWantsLogger)
	if wantsLogger && eg.client != nil {
		serviceLogger.SetLogger(eg.logger)
	}

	serviceClient, wantsClient := eg.service.(ServiceWantsClient)
	if wantsClient && eg.client != nil {
		serviceClient.SetClient(eg.client)
	}

	serviceServer, wantsServer := eg.service.(ServiceWantsServer)
	if wantsServer && eg.server != nil {
		serviceServer.SetServer(eg.server.server)
	}

	if nil == eg.server {
		return nil
	}

	return eg.server.listenAndServe()
}
