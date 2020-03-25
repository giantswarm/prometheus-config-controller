package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/k8sclient/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/prometheus-config-controller/flag"
	"github.com/giantswarm/prometheus-config-controller/service/controller"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

type Config struct {
	Flag   *flag.Flag
	Logger micrologger.Logger
	Viper  *viper.Viper

	Description string
	GitCommit   string
	ProjectName string
	Source      string
	Version     string
}

type Service struct {
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

	var k8sClient *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1alpha1.AddToScheme,
			},
			Logger:     config.Logger,
			RestConfig: restConfig,
		}
		k8sClient, err = k8sclient.NewClients(c)
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
			PrometheusAddress:  config.Viper.GetString(config.Flag.Service.Prometheus.Address),
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
			Name:           config.ProjectName,
			Source:         config.Source,
			Version:        config.Version,
			VersionBundles: NewVersionBundles(),
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
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

		// Prometheus won't start in 90 seconds anyway so let's not
		// spam with logs and wait for it.
		time.Sleep(90 * time.Second)

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
