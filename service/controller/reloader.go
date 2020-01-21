package controller

import (
	"time"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/prometheus-config-controller/pkg/project"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

type ReloaderConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	PrometheusAddress string
}

type Reloader struct {
	*controller.Controller
}

func NewReloader(config ReloaderConfig) (*Reloader, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.PrometheusAddress == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusAddress must not be empty", config)
	}

	var err error

	var resourceSet *controller.ResourceSet
	{
		c := reloaderResourceSetConfig{
			Logger: config.Logger,

			PrometheusAddress: config.PrometheusAddress,
		}

		resourceSet, err = newReloaderResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(corev1.ConfigMap)
			},
			Logger: config.Logger,
			ResourceSets: []*controller.ResourceSet{
				resourceSet,
			},
			Selector: key.LabelSelectorConfigMap(),

			Name:         project.Name() + "-reloader",
			ResyncPeriod: 3 * time.Minute,
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Reloader{
		Controller: operatorkitController,
	}

	return c, nil
}
