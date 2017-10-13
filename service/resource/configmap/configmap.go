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

	CertificateDirectory string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,

		CertificateDirectory: "",
		ConfigMapKey:         "",
		ConfigMapName:        "",
		ConfigMapNamespace:   "",
	}
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	certificateDirectory string
	configMapKey         string
	configMapName        string
	configMapNamespace   string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.CertificateDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertificateDirectory must not be empty")
	}
	if config.ConfigMapKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapKey must not be empty")
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

		certificateDirectory: config.CertificateDirectory,
		configMapKey:         config.ConfigMapKey,
		configMapName:        config.ConfigMapName,
		configMapNamespace:   config.ConfigMapNamespace,
	}

	return resource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
