package certificate

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/spf13/afero"
)

const (
	Name = "certificate"
)

type Config struct {
	Fs     afero.Fs
	Logger micrologger.Logger

	CertificateDirectory string
}

func DefaultConfig() Config {
	return Config{
		Fs:     nil,
		Logger: nil,

		CertificateDirectory: "",
	}
}

type Resource struct {
	fs     afero.Fs
	logger micrologger.Logger

	certificateDirectory string
}

func New(config Config) (*Resource, error) {
	if config.Fs == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Fs must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.CertificateDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CertificateDirectory must not be empty")
	}

	resource := &Resource{
		fs:     config.Fs,
		logger: config.Logger,

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
