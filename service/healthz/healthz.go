package healthz

import (
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/k8shealthz"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,
	}
}

type Service struct {
	K8s healthz.Service
}

func New(config Config) (*Service, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	var err error

	var newK8sService healthz.Service
	{
		k8sConfig := k8shealthz.DefaultConfig()

		k8sConfig.K8sClient = config.K8sClient
		k8sConfig.Logger = config.Logger

		newK8sService, err = k8shealthz.New(k8sConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		K8s: newK8sService,
	}

	return newService, nil
}
