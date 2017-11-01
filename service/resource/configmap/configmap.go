package configmap

import (
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
)

const (
	Name = "configmap"
)

type Config struct {
	K8sClient          kubernetes.Interface
	Logger             micrologger.Logger
	PrometheusReloader prometheus.PrometheusReloader

	CertificateDirectory string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	ReloadWaitTime     time.Duration
}

func DefaultConfig() Config {
	return Config{
		K8sClient:          nil,
		Logger:             nil,
		PrometheusReloader: nil,

		CertificateDirectory: "",
		ConfigMapKey:         "",
		ConfigMapName:        "",
		ConfigMapNamespace:   "",
		ReloadWaitTime:       time.Duration(0),
	}
}

type Resource struct {
	k8sClient          kubernetes.Interface
	logger             micrologger.Logger
	prometheusReloader prometheus.PrometheusReloader

	certificateDirectory string
	configMapKey         string
	configMapName        string
	configMapNamespace   string
	reloadWaitTime       time.Duration
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
	if config.ReloadWaitTime == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ReloadWaitTime must not be zero")
	}

	resource := &Resource{
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		prometheusReloader: config.PrometheusReloader,

		certificateDirectory: config.CertificateDirectory,
		configMapKey:         config.ConfigMapKey,
		configMapName:        config.ConfigMapName,
		configMapNamespace:   config.ConfigMapNamespace,
		reloadWaitTime:       config.ReloadWaitTime,
	}

	return resource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
