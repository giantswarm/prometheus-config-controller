package server

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/giantswarm/prometheus-config-controller/server/endpoint"
	"github.com/giantswarm/prometheus-config-controller/service"
)

type Config struct {
	Service *service.Service

	MicroServerConfig microserver.Config
}

func DefaultConfig() Config {
	return Config{
		Service: nil,

		MicroServerConfig: microserver.DefaultConfig(),
	}
}

func New(config Config) (microserver.Server, error) {
	var err error

	var newEndpointCollection *endpoint.Endpoint
	{
		endpointConfig := endpoint.DefaultConfig()

		endpointConfig.Logger = config.MicroServerConfig.Logger
		endpointConfig.Service = config.Service

		newEndpointCollection, err = endpoint.New(endpointConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newServer := &server{
		logger: config.MicroServerConfig.Logger,

		config: config.MicroServerConfig,
	}

	newServer.config.Endpoints = []microserver.Endpoint{
		newEndpointCollection.Healthz,
		newEndpointCollection.Version,
	}

	newServer.config.ErrorEncoder = newServer.newErrorEncoder()

	return newServer, nil
}

type server struct {
	logger micrologger.Logger

	config microserver.Config
}

func (s *server) Boot() {}

func (s *server) Config() microserver.Config {
	return s.config
}

func (s *server) Shutdown() {}

func (s *server) newErrorEncoder() kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		rErr := err.(microserver.ResponseError)
		uErr := rErr.Underlying()

		rErr.SetCode(microserver.CodeInternalError)
		rErr.SetMessage(uErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
