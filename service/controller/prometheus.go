package controller

import (
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-config-controller/pkg/project"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

type PrometheusConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	ConfigMapPath      string
	CertComponentName  string
	CertDirectory      string
	CertNamespace      string
	CertPermission     int
	PrometheusAddress  string
}

type Prometheus struct {
	*controller.Controller
}

func NewPrometheus(config PrometheusConfig) (*Prometheus, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.ConfigMapKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapKey must not be empty", config)
	}
	if config.ConfigMapName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapName must not be empty", config)
	}
	if config.ConfigMapNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapNamespace must not be empty", config)
	}
	if config.ConfigMapPath == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapPath must not be empty", config)
	}
	if config.CertComponentName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.CertComponentName must not be empty", config)
	}
	if config.CertDirectory == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.CertDirectory must not be empty", config)
	}
	if config.CertNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.CertNamespace must not be empty", config)
	}
	if config.CertPermission == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.CertPermission must not be empty", config)
	}
	if config.PrometheusAddress == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusAddress must not be empty", config)
	}

	var err error

	var resourceSet *controller.ResourceSet
	{
		c := prometheusResourceSetConfig{
			K8sClient: config.K8sClient.K8sClient(),
			Logger:    config.Logger,

			ConfigMapKey:       config.ConfigMapKey,
			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,
			ConfigMapPath:      config.ConfigMapPath,
			CertComponentName:  config.CertComponentName,
			CertDirectory:      config.CertDirectory,
			CertNamespace:      config.CertNamespace,
			CertPermission:     config.CertPermission,
			PrometheusAddress:  config.PrometheusAddress,
		}

		resourceSet, err = newPrometheusResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(corev1.Service)
			},
			Logger: config.Logger,
			ResourceSets: []*controller.ResourceSet{
				resourceSet,
			},
			Selector: key.LabelSelectorService(),

			Name: project.Name(),
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Prometheus{
		Controller: operatorkitController,
	}

	return c, nil
}
