package resource

import (
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/crud"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/retryresource"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/certificate"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/configmap"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/reload"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	CertComponentName  string
	CertDirectory      string
	CertNamespace      string
	CertPermission     int
	PrometheusAddress  string
	Provider           string
}

func New(config Config) ([]resource.Interface, error) {
	var err error

	var certificateResource resource.Interface
	{
		c := certificate.Config{
			Fs:        afero.NewOsFs(),
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			CertComponentName: config.CertComponentName,
			CertDirectory:     config.CertDirectory,
			CertNamespace:     config.CertNamespace,
			CertPermission:    os.FileMode(config.CertPermission),
		}

		ops, err := certificate.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		certificateResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var configMapResource resource.Interface
	{
		c := configmap.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			CertDirectory:      config.CertDirectory,
			ConfigMapKey:       config.ConfigMapKey,
			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,

			Provider: config.Provider,
		}

		configMapResource, err = configmap.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var reloadResource resource.Interface
	{
		c := reload.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,
			PrometheusAddress:  config.PrometheusAddress,
		}

		reloadResource, err = reload.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		certificateResource,
		configMapResource,
		reloadResource,
	}

	{
		c := retryresource.WrapConfig{
			Logger: config.Logger,
		}

		resources, err = retryresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	{
		c := metricsresource.WrapConfig{}

		resources, err = metricsresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resources, nil
}

func toCRUDResource(logger micrologger.Logger, ops crud.Interface) (*crud.Resource, error) {
	c := crud.ResourceConfig{
		CRUD:   ops,
		Logger: logger,
	}

	r, err := crud.NewResource(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}
