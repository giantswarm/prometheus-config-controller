package endpoint

import (
	"github.com/giantswarm/microendpoint/endpoint/healthz"
	"github.com/giantswarm/microendpoint/endpoint/version"
	healthzservice "github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/prometheus-config-controller/service"
)

type Config struct {
	Logger  micrologger.Logger
	Service *service.Service
}

func DefaultConfig() Config {
	return Config{
		Logger:  nil,
		Service: nil,
	}
}

type Endpoint struct {
	Healthz *healthz.Endpoint
	Version *version.Endpoint
}

func New(config Config) (*Endpoint, error) {
	var err error

	var newHealthzEndpoint *healthz.Endpoint
	{
		healthzConfig := healthz.DefaultConfig()

		healthzConfig.Logger = config.Logger
		healthzConfig.Services = []healthzservice.Service{
			config.Service.Healthz.K8s,
		}

		newHealthzEndpoint, err = healthz.New(healthzConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newVersionEndpoint *version.Endpoint
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Logger = config.Logger
		versionConfig.Service = config.Service.Version

		newVersionEndpoint, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newEndpoint := &Endpoint{
		Healthz: newHealthzEndpoint,
		Version: newVersionEndpoint,
	}

	return newEndpoint, nil
}
