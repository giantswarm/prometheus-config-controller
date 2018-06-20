package configmap

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

const (
	Name = "configmapv1"
)

type Config struct {
	K8sClient          kubernetes.Interface
	Logger             micrologger.Logger
	PrometheusReloader prometheus.PrometheusReloader

	CertDirectory string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
}

type Resource struct {
	k8sClient          kubernetes.Interface
	logger             micrologger.Logger
	prometheusReloader prometheus.PrometheusReloader

	certDirectory      string
	configMapKey       string
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
	if config.PrometheusReloader == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.PrometheusReloader must not be empty")
	}

	if config.CertDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertDirectory must not be empty")
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

	r := &Resource{
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		prometheusReloader: config.PrometheusReloader,

		certDirectory:      config.CertDirectory,
		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
