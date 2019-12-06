package controller

import (
	"time"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1"
)

type PrometheusConfig struct {
	K8sClient k8sclient.Interface
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
	ResyncPeriod       time.Duration
}

type Prometheus struct {
	*controller.Controller
}

func NewPrometheus(config PrometheusConfig) (*Prometheus, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	var err error

	var resourceSetV1 *controller.ResourceSet
	{
		c := v1.ResourceSetConfig{
			K8sClient: config.K8sClient.K8sClient(),
			Logger:    config.Logger,

			ConfigMapKey:       config.ConfigMapKey,
			ConfigMapName:      config.ConfigMapName,
			ConfigMapNamespace: config.ConfigMapNamespace,
			CertComponentName:  config.CertComponentName,
			CertDirectory:      config.CertDirectory,
			CertNamespace:      config.CertNamespace,
			CertPermission:     config.CertPermission,
			MinReloadTime:      config.MinReloadTime,
			ProjectName:        config.ProjectName,
			PrometheusAddress:  config.PrometheusAddress,
		}

		resourceSetV1, err = v1.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      config.ProjectName,
			ResourceSets: []*controller.ResourceSet{
				resourceSetV1,
			},
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(corev1.Service)
			},
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
