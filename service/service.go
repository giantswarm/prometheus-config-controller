package service

import (
	"os"
	"sync"
	"time"

	"github.com/cenk/backoff"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8sclient"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/logresource"
	"github.com/giantswarm/operatorkit/framework/resource/metricsresource"
	"github.com/giantswarm/operatorkit/framework/resource/retryresource"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/service/controller"
	"github.com/giantswarm/prometheus-config-controller/service/healthz"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
	"github.com/giantswarm/prometheus-config-controller/service/resource/certificate"
	"github.com/giantswarm/prometheus-config-controller/service/resource/configmap"
)

type Config struct {
	Flag   *flag.Flag
	Logger micrologger.Logger
	Viper  *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string

	ControllerBackOffDuration time.Duration
	FrameworkBackOffDuration  time.Duration
	ResourceRetries           int
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

		ControllerBackOffDuration: time.Duration(0),
		FrameworkBackOffDuration:  time.Duration(0),
		ResourceRetries:           0,
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

	if config.ControllerBackOffDuration == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ControllerBackOffDuration must not be zero")
	}
	if config.FrameworkBackOffDuration == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.FrameworkBackOffDuration must not be zero")
	}
	if config.ResourceRetries == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ResourceRetries must not be zero")
	}

	var err error

	var newFs afero.Fs
	{
		newFs = afero.NewOsFs()
	}

	var newK8sClient kubernetes.Interface
	{
		k8sConfig := k8sclient.DefaultConfig()

		k8sConfig.Logger = config.Logger

		k8sConfig.Address = config.Viper.GetString(config.Flag.Service.Kubernetes.Address)
		k8sConfig.InCluster = config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster)
		k8sConfig.TLS.CAFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile)
		k8sConfig.TLS.CrtFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile)
		k8sConfig.TLS.KeyFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile)

		newK8sClient, err = k8sclient.New(k8sConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newPrometheusReloader prometheus.PrometheusReloader
	{
		prometheusConfig := prometheus.DefaultConfig()

		prometheusConfig.K8sClient = newK8sClient
		prometheusConfig.Logger = config.Logger

		prometheusConfig.Address = config.Viper.GetString(config.Flag.Service.Prometheus.Address)
		prometheusConfig.ConfigMapKey = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Key)
		prometheusConfig.ConfigMapName = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Name)
		prometheusConfig.ConfigMapNamespace = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Namespace)
		prometheusConfig.MinimumReloadTime = config.Viper.GetDuration(config.Flag.Service.Resource.ConfigMap.MinimumReloadTime)

		newPrometheusReloader, err = prometheus.New(prometheusConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newCertificateResource framework.Resource
	{
		certificateConfig := certificate.DefaultConfig()

		certificateConfig.Fs = newFs
		certificateConfig.K8sClient = newK8sClient
		certificateConfig.Logger = config.Logger
		certificateConfig.PrometheusReloader = newPrometheusReloader

		certificateConfig.CertificateComponentName = config.Viper.GetString(config.Flag.Service.Resource.Certificate.ComponentName)
		certificateConfig.CertificateDirectory = config.Viper.GetString(config.Flag.Service.Resource.Certificate.Directory)
		certificateConfig.CertificateNamespace = config.Viper.GetString(config.Flag.Service.Resource.Certificate.Namespace)
		certificateConfig.CertificatePermission = os.FileMode(config.Viper.GetInt(config.Flag.Service.Resource.Certificate.Permission))

		newCertificateResource, err = certificate.New(certificateConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newConfigMapResource framework.Resource
	{
		configMapConfig := configmap.DefaultConfig()

		configMapConfig.K8sClient = newK8sClient
		configMapConfig.Logger = config.Logger
		configMapConfig.PrometheusReloader = newPrometheusReloader

		configMapConfig.CertificateDirectory = config.Viper.GetString(config.Flag.Service.Resource.Certificate.Directory)
		configMapConfig.ConfigMapKey = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Key)
		configMapConfig.ConfigMapName = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Name)
		configMapConfig.ConfigMapNamespace = config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Namespace)

		newConfigMapResource, err = configmap.New(configMapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var resources []framework.Resource
	{
		resources = []framework.Resource{
			newCertificateResource,
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
		metricsWrapConfig.Name = config.Name
		resources, err = metricsresource.Wrap(resources, metricsWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newOperatorFramework *framework.Framework
	{
		backOff := backoff.NewExponentialBackOff()
		backOff.MaxElapsedTime = config.FrameworkBackOffDuration

		frameworkConfig := framework.DefaultConfig()

		frameworkConfig.BackOff = backOff
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
		backOff := backoff.NewExponentialBackOff()
		backOff.MaxElapsedTime = config.ControllerBackOffDuration

		controllerConfig := controller.DefaultConfig()

		controllerConfig.BackOff = backOff
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
