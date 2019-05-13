package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/service/controller"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
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

	logger micrologger.Logger

	bootOnce             sync.Once
	prometheusAddress    string
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

		logger: config.Logger,

		bootOnce: sync.Once{},

		prometheusAddress:    config.Viper.GetString(config.Flag.Service.Prometheus.Address),
		prometheusController: prometheusController,
	}

	return s, nil
}

func (s *Service) Boot() {
	ctx := context.TODO()

	err := s.boot(ctx)
	if err != nil {
		s.logger.LogCtx(ctx, "level", "error", "message", "failed to boot the service", "stack", fmt.Sprintf("%#v", err))
		panic(fmt.Sprintf("failed to boot the service, please see the logs"))
	}
}

func (s *Service) boot(ctx context.Context) error {
	// Wait for Prometheus to be ready before booting the controller.
	// Otherwise it will fail to (re)load the configuration.
	{
		s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for Prometheus to be up"))

		url := key.PrometheusURLConfig(s.prometheusAddress)

		o := func() error {
			res, err := http.Get(url)
			if err != nil {
				return microerror.Maskf(waitError, "failed request URL %#q with error %#q", url, err)
			}

			if res.StatusCode < 200 || res.StatusCode > 299 {
				return microerror.Maskf(waitError, "expected 2xx response for URL %#q but got %d", url, res.StatusCode)
			}

			return nil
		}
		b := backoff.NewMaxRetries(10, 60*time.Second)
		n := backoff.NewNotifier(s.logger, ctx)

		// Prometheus won't start in 90 seconds anyway so let's not
		// spam with logs and wait for it.
		time.Sleep(90 * time.Second)

		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		s.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waited for Prometheus to be up"))
	}

	s.bootOnce.Do(func() {
		go s.prometheusController.Boot(ctx)
	})

	return nil
}
