package certificate

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

const (
	Name = "certificatev1"
)

type Config struct {
	Fs                 afero.Fs
	K8sClient          kubernetes.Interface
	Logger             micrologger.Logger
	PrometheusReloader prometheus.PrometheusReloader

	CertComponentName string
	CertDirectory     string
	CertNamespace     string
	CertPermission    os.FileMode
}

type Resource struct {
	fs                 afero.Fs
	k8sClient          kubernetes.Interface
	logger             micrologger.Logger
	prometheusReloader prometheus.PrometheusReloader

	certComponentName string
	certDirectory     string
	certNamespace     string
	certPermission    os.FileMode
}

func New(config Config) (*Resource, error) {
	if config.Fs == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Fs must not be empty")
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.PrometheusReloader == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.PrometheusReloader must not be empty")
	}

	if config.CertComponentName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertComponentName must not be empty")
	}
	if config.CertDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertDirectory must not be empty")
	}
	if config.CertNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertNamespace must not be empty")
	}
	if config.CertPermission == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.CertPermission must not be zero")
	}

	r := &Resource{
		fs:                 config.Fs,
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		prometheusReloader: config.PrometheusReloader,

		certComponentName: config.CertComponentName,
		certDirectory:     config.CertDirectory,
		certNamespace:     config.CertNamespace,
		certPermission:    config.CertPermission,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
