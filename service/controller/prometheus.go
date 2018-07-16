package controller

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/informer"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1"
)

type PrometheusConfig struct {
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

	var newInformer *informer.Informer
	{
		c := informer.Config{
			Logger:  config.Logger,
			Watcher: config.K8sClient.CoreV1().Services(""),

			RateWait:     informer.DefaultRateWait,
			ResyncPeriod: config.ResyncPeriod,
		}

		newInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var resourceSetV1 *controller.ResourceSet
	{
		c := v1.ResourceSetConfig{
			K8sClient: config.K8sClient,
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
			Informer: newInformer,
			Logger:   config.Logger,
			ResourceSets: []*controller.ResourceSet{
				resourceSetV1,
			},
			RESTClient: config.K8sClient.CoreV1().RESTClient(),

			Name: config.ProjectName,
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
