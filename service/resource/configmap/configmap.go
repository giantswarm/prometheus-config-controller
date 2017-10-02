package configmap

import (
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

const (
	Name = "configmap"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	ConfigMapName      string
	ConfigMapNamespace string
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,

		ConfigMapName:      "",
		ConfigMapNamespace: "",
	}
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	configMapName      string
	configMapNamespace string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.ConfigMapName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapName must not be empty")
	}
	if config.ConfigMapNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapNamespace must not be empty")
	}

	resource := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
	}

	return resource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
