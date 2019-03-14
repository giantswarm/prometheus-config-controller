package service

import (
	"context"
	"sync"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/service/controller"
	"github.com/giantswarm/prometheus-config-controller/service/healthz"
)

type Config struct {
	Logger micrologger.Logger

	Description string
	Flag        *flag.Flag
	GitCommit   string
	Name        string
	Source      string
	Viper       *viper.Viper
}

type Service struct {
	Healthz *healthz.Service
	Version *version.Service

	bootOnce             sync.Once
	prometheusController *controller.Prometheus
}

func New(config Config) (*Service, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:    config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster:  config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			KubeConfig: config.Viper.GetString(config.Flag.Service.Kubernetes.KubeConfig),
			TLS: k8srestconfig.ConfigTLS{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
				CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
				KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var healthzService *healthz.Service
	{
		c := healthz.Config{
			K8sClient: k8sClient,
			Logger:    config.Logger,
		}

		healthzService, err = healthz.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusController *controller.Prometheus
	{
		c := controller.PrometheusConfig{
			K8sClient: k8sClient,
			Logger:    config.Logger,

			ConfigMapKey:       config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Key),
			ConfigMapName:      config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Name),
			ConfigMapNamespace: config.Viper.GetString(config.Flag.Service.Resource.ConfigMap.Namespace),
			CertComponentName:  config.Viper.GetString(config.Flag.Service.Resource.Certificate.ComponentName),
			CertDirectory:      config.Viper.GetString(config.Flag.Service.Resource.Certificate.Directory),
			CertNamespace:      config.Viper.GetString(config.Flag.Service.Resource.Certificate.Namespace),
			CertPermission:     config.Viper.GetInt(config.Flag.Service.Resource.Certificate.Permission),
			MinReloadTime:      config.Viper.GetDuration(config.Flag.Service.Resource.ConfigMap.MinimumReloadTime),
			ProjectName:        config.Name,
			PrometheusAddress:  config.Viper.GetString(config.Flag.Service.Prometheus.Address),
			ResyncPeriod:       config.Viper.GetDuration(config.Flag.Service.Controller.ResyncPeriod),
		}

		prometheusController, err = controller.NewPrometheus(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		c := version.Config{
			Description:    config.Description,
			GitCommit:      config.GitCommit,
			Name:           config.Name,
			Source:         config.Source,
			VersionBundles: NewVersionBundles(),
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Healthz: healthzService,
		Version: versionService,

		bootOnce:             sync.Once{},
		prometheusController: prometheusController,
	}

	return s, nil
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		go s.prometheusController.Boot(context.Background())
	})
}
