package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// Address is the address of the Prometheus instance to manage.
	Address string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	MinimumReloadTime  time.Duration
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,

		Address:            "",
		ConfigMapKey:       "",
		ConfigMapName:      "",
		ConfigMapNamespace: "",
		MinimumReloadTime:  0,
	}
}

type Service struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	configMapKey       string
	configMapName      string
	configMapNamespace string
	minimumReloadTime  time.Duration

	isReloadRequested      bool
	isReloadRequestedMutex sync.Mutex
	lastReloadTime         time.Time
	urlConfig              *url.URL
	urlReload              *url.URL
}

func New(config Config) (*Service, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.Address == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Address must not be empty")
	}
	if config.ConfigMapKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapKey must not be empty")
	}
	if config.ConfigMapName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapName must not be empty")
	}
	if config.ConfigMapNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ConfigMapNamespace must not be empty")
	}
	if config.MinimumReloadTime == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.MinimumReloadTime must not be zero")
	}

	urlBase, err := url.ParseRequestURI(config.Address)
	if err != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Address must be a valid URI but got %#q", config.Address)
	}

	urlConfigRelative, err := url.Parse(ConfigPath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	urlReloadRelative, err := url.Parse(ReloadPath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	service := &Service{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
		minimumReloadTime:  config.MinimumReloadTime,

		isReloadRequested:      false,
		isReloadRequestedMutex: sync.Mutex{},
		lastReloadTime:         time.Time{},
		urlConfig:              urlBase.ResolveReference(urlConfigRelative),
		urlReload:              urlBase.ResolveReference(urlReloadRelative),
	}

	return service, nil
}

func (s *Service) Reload(ctx context.Context) error {
	err := s.throttleReload(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	reloadRequired, err := s.isReloadRequired(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	if reloadRequired {
		if err := s.reload(ctx); err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *Service) RequestReload(ctx context.Context) {
	s.logger.LogCtx(ctx, "debug", "reload requested")

	s.isReloadRequestedMutex.Lock()
	defer s.isReloadRequestedMutex.Unlock()

	s.isReloadRequested = true
}

// getConfigFromKubernetes returns the configuration that is in the configmap.
func (s *Service) getConfigFromKubernetes(ctx context.Context) (string, error) {
	s.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching configmap: %s/%s", s.configMapNamespace, s.configMapName))

	configMap, err := s.k8sClient.CoreV1().ConfigMaps(s.configMapNamespace).Get(s.configMapName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return "", microerror.Maskf(executionFailedError, "configmap %#q in namespace %#q does not exist", s.configMapName, s.configMapNamespace)
	} else if err != nil {
		return "", microerror.Mask(err)
	}

	val, ok := configMap.Data[s.configMapKey]
	if !ok {
		return "", microerror.Maskf(executionFailedError, "configmap key not present")
	}

	return val, nil
}

// getConfigFromPrometheus returns the configuration that is currently loaded in Prometheus.
func (s *Service) getConfigFromPrometheus(ctx context.Context) (string, error) {
	s.logger.LogCtx(ctx, "debug", "fetching current prometheus config")

	res, err := http.Get(s.urlConfig.String())
	if err != nil {
		return "", microerror.Mask(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", microerror.Maskf(executionFailedError, "non-200 HTTP status code was returned: %d", res.StatusCode)
	}

	var config string
	{
		decoder := json.NewDecoder(res.Body)
		defer res.Body.Close()

		resp := struct {
			Status string `json:"status"`
			Data   struct {
				YAML string `json:"yaml"`
			} `json:"data"`
		}{}

		err = decoder.Decode(&resp)
		if err != nil {
			return "", microerror.Maskf(executionFailedError, "decode prometheus config response with error: %s", err.Error())
		}

		if resp.Status != "success" {
			return "", microerror.Maskf(executionFailedError, "prometheus returned non-success response status when reloding config: %s", resp.Status)
		}

		config = resp.Data.YAML
	}

	return config, nil
}

func (s *Service) isReloadRequired(ctx context.Context) (bool, error) {
	s.logger.LogCtx(ctx, "debug", "checking if reload is required")

	configurationReloadCheckCount.Inc()

	s.isReloadRequestedMutex.Lock()
	isReloadRequested := s.isReloadRequested
	s.isReloadRequestedMutex.Unlock()

	if isReloadRequested {
		s.logger.LogCtx(ctx, "debug", "reload was requested previously")

		configurationReloadRequiredCount.Inc()

		s.isReloadRequestedMutex.Lock()
		s.isReloadRequested = false
		s.isReloadRequestedMutex.Unlock()

		return true, nil
	}

	kubernetesConfiguration, err := s.getConfigFromKubernetes(ctx)
	if err != nil {
		return false, microerror.Mask(err)
	}

	prometheusConfiguration, err := s.getConfigFromPrometheus(ctx)
	if err != nil {
		return false, microerror.Mask(err)
	}

	if kubernetesConfiguration != prometheusConfiguration {
		configurationReloadRequiredCount.Inc()

		s.logger.LogCtx(ctx, "debug", "kubernetes and prometheus configuration do not match, reload required")
		return true, nil
	}

	s.logger.LogCtx(ctx, "debug", "kubernetes and prometheus configuration match, reload not required")
	return false, nil
}

func (s *Service) reload(ctx context.Context) error {
	s.logger.LogCtx(ctx, "debug", "reloading prometheus config")

	res, err := http.Post(s.urlReload.String(), "", nil)
	if err != nil {
		return microerror.Mask(err)
	}
	if res.StatusCode != http.StatusOK {
		return microerror.Maskf(executionFailedError, "non-200 status code was returned: %d", res.StatusCode)
	}

	configurationReloadCount.Inc()

	s.lastReloadTime = time.Now()

	return nil
}

func (s *Service) throttleReload(ctx context.Context) error {
	timeSinceLastReload := time.Since(s.lastReloadTime)

	if timeSinceLastReload < s.minimumReloadTime {
		configurationReloadIgnoredCount.Inc()

		return microerror.Maskf(reloadThrottleError, "%s since last reload, minimum time between is %s", timeSinceLastReload, s.minimumReloadTime)
	}

	return nil
}
