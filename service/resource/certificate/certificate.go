package certificate

import (
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

const (
	Name = "certificate"
)

type Config struct {
	Fs        afero.Fs
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	CertificateDirectory string
}

func DefaultConfig() Config {
	return Config{
		Fs:        nil,
		K8sClient: nil,
		Logger:    nil,

		CertificateDirectory: "",
	}
}

type Resource struct {
	certificatetprService *certificatetpr.Service
	fs                    afero.Fs
	k8sClient             kubernetes.Interface
	logger                micrologger.Logger

	certificateDirectory string
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

	if config.CertificateDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertificateDirectory must not be empty")
	}

	certificatetprConfig := certificatetpr.DefaultServiceConfig()

	certificatetprConfig.K8sClient = config.K8sClient
	certificatetprConfig.Logger = config.Logger

	certificatetprService, err := certificatetpr.NewService(certificatetprConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	resource := &Resource{
		certificatetprService: certificatetprService,
		fs:        config.Fs,
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		certificateDirectory: config.CertificateDirectory,
	}

	return resource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
