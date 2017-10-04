package service

import (
	"sync"
	"time"

	"github.com/cenk/backoff"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8s"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/logresource"
	"github.com/giantswarm/operatorkit/framework/metricsresource"
	"github.com/giantswarm/operatorkit/framework/retryresource"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/service/controller"
	"github.com/giantswarm/prometheus-config-controller/service/healthz"
	configmapresource "github.com/giantswarm/prometheus-config-controller/service/resource/configmap"
)

type Config struct {
	Flag   *flag.Flag
	Logger micrologger.Logger
	Viper  *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string

	ResourceRetries           int
	ControllerBackOffDuration time.Duration
}

func DefaultConfig() Config {
	return Config{
		Flag:   nil,
		Logger: nil,
		Viper:  nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",

		ResourceRetries:           0,
		ControllerBackOffDuration: time.Duration(0),
	}
}

type Service struct {
	Controller *controller.Controller
	Healthz    *healthz.Service
	Version    *version.Service

	bootOnce sync.Once
}

func New(config Config) (*Service, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Viper must not be empty")
	}

	if config.ResourceRetries == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResourceRetries must not be zero")
	}
	if config.ControllerBackOffDuration == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ControllerBackOffDuration must not be zero")
	}

	var err error

	var newK8sClient kubernetes.Interface
	{
		k8sConfig := k8s.DefaultConfig()

		k8sConfig.Logger = config.Logger

		k8sConfig.Address = config.Viper.GetString(config.Flag.Service.Kubernetes.Address)
		k8sConfig.InCluster = config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster)
		k8sConfig.TLS.CAFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile)
		k8sConfig.TLS.CrtFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile)
		k8sConfig.TLS.KeyFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile)

		newK8sClient, err = k8s.NewClient(k8sConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newConfigMapResource framework.Resource
	{
		configMapConfig := configmapresource.DefaultConfig()

		configMapConfig.K8sClient = newK8sClient
		configMapConfig.Logger = config.Logger

		configMapConfig.ConfigMapKey = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Key)
		configMapConfig.ConfigMapName = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Name)
		configMapConfig.ConfigMapNamespace = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Namespace)

		newConfigMapResource, err = configmapresource.New(configMapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var resources []framework.Resource
	{
		resources = []framework.Resource{
			newConfigMapResource,
		}

		logWrapConfig := logresource.DefaultWrapConfig()
		logWrapConfig.Logger = config.Logger
		resources, err = logresource.Wrap(resources, logWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		retryWrapConfig := retryresource.DefaultWrapConfig()
		retryWrapConfig.BackOffFactory = func() backoff.BackOff {
			return backoff.WithMaxTries(backoff.NewExponentialBackOff(), uint64(config.ResourceRetries))
		}
		retryWrapConfig.Logger = config.Logger
		resources, err = retryresource.Wrap(resources, retryWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		metricsWrapConfig := metricsresource.DefaultWrapConfig()
		metricsWrapConfig.Namespace = config.Name
		resources, err = metricsresource.Wrap(resources, metricsWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newControllerBackOff *backoff.ExponentialBackOff
	{
		newControllerBackOff = backoff.NewExponentialBackOff()
		newControllerBackOff.MaxElapsedTime = config.ControllerBackOffDuration
	}

	var newOperatorFramework *framework.Framework
	{
		frameworkConfig := framework.DefaultConfig()

		frameworkConfig.Logger = config.Logger
		frameworkConfig.Resources = resources

		newOperatorFramework, err = framework.New(frameworkConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newHealthzService *healthz.Service
	{
		healthzConfig := healthz.DefaultConfig()

		healthzConfig.K8sClient = newK8sClient
		healthzConfig.Logger = config.Logger

		newHealthzService, err = healthz.New(healthzConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newController *controller.Controller
	{
		controllerConfig := controller.DefaultConfig()

		controllerConfig.BackOff = newControllerBackOff
		controllerConfig.K8sClient = newK8sClient
		controllerConfig.Logger = config.Logger
		controllerConfig.OperatorFramework = newOperatorFramework

		controllerConfig.ResyncPeriod = config.Viper.GetDuration(config.Flag.Service.Controller.ResyncPeriod)

		newController, err = controller.New(controllerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newVersionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		newVersionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		Controller: newController,
		Healthz:    newHealthzService,
		Version:    newVersionService,

		bootOnce: sync.Once{},
	}

	return newService, nil
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		s.Controller.Boot()
	})
}
