package v1

import (
	"os"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/resource/metricsresource"
	"github.com/giantswarm/operatorkit/controller/resource/retryresource"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/certificate"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/configmap"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/resource/patchfinalizer"
)

type ResourceSetConfig struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	CertComponentName  string
	CertDirectory      string
	CertNamespace      string
	CertPermission     int
	MinReloadTime      time.Duration
	ProjectName        string
	PrometheusAddress  string
}

func NewResourceSet(config ResourceSetConfig) (*controller.ResourceSet, error) {
	var err error

	var prometheusReloader prometheus.PrometheusReloader
	{
		c := prometheus.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			Address:            config.PrometheusAddress,
			ConfigMapKey:       config.ConfigMapKey,
			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,
			MinimumReloadTime:  config.MinReloadTime,
		}

		prometheusReloader, err = prometheus.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var certificateResource controller.Resource
	{
		c := certificate.Config{
			Fs:                 afero.NewOsFs(),
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			PrometheusReloader: prometheusReloader,

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

	var configMapResource controller.Resource
	{
		c := configmap.Config{
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			PrometheusReloader: prometheusReloader,

			CertDirectory:      config.CertDirectory,
			ConfigMapKey:       config.ConfigMapKey,
			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,
		}

		ops, err := configmap.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		configMapResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var patchFinalizerResource controller.Resource
	{
		c := patchfinalizer.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		patchFinalizerResource, err = patchfinalizer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []controller.Resource{
		certificateResource,
		configMapResource,
		patchFinalizerResource,
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

	handlesFunc := func(obj interface{}) bool {
		return true
	}

	var resourceSet *controller.ResourceSet
	{
		c := controller.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = controller.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}

func toCRUDResource(logger micrologger.Logger, ops controller.CRUDResourceOps) (*controller.CRUDResource, error) {
	c := controller.CRUDResourceConfig{
		Logger: logger,
		Ops:    ops,
	}

	r, err := controller.NewCRUDResource(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}
