package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

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

	// address is the address of the Prometheus instance we manage.
	address            string
	configMapKey       string
	configMapName      string
	configMapNamespace string
	minimumReloadTime  time.Duration

	isReloadRequested      bool
	isReloadRequestedMutex sync.Mutex
	lastReloadTime         time.Time
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

	service := &Service{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		address:            config.Address,
		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
		minimumReloadTime:  config.MinimumReloadTime,

		isReloadRequested:      false,
		isReloadRequestedMutex: sync.Mutex{},
		lastReloadTime:         time.Time{},
	}

	return service, nil
}

func (s *Service) Reload(ctx context.Context) error {
	reloadRequired, err := s.isReloadRequired(ctx)
	if err != nil {
		return microerror.Maskf(reloadError, err.Error())
	}

	if reloadRequired {
		if err := s.reload(ctx); err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *Service) isReloadRateLimited(ctx context.Context) bool {
	timeSinceLastReload := time.Since(s.lastReloadTime)

	if timeSinceLastReload < s.minimumReloadTime {
		s.logger.LogCtx(ctx, "debug", fmt.Sprintf("ignoring reload request, only %s since last reload, minimum time between is %s", timeSinceLastReload, s.minimumReloadTime))
		configurationReloadIgnoredCount.Inc()

		return true
	}

	return false
}

func (s *Service) RequestReload(ctx context.Context) {
	s.logger.LogCtx(ctx, "debug", "reload requested")

	s.isReloadRequestedMutex.Lock()
	defer s.isReloadRequestedMutex.Unlock()

	s.isReloadRequested = true
}

func (s *Service) IsReloadRequested() bool {
	s.isReloadRequestedMutex.Lock()
	defer s.isReloadRequestedMutex.Unlock()

	return s.isReloadRequested
}

func (s *Service) isReloadRequired(ctx context.Context) (bool, error) {
	s.logger.LogCtx(ctx, "debug", "checking if reload is required")

	configurationReloadCheckCount.Inc()

	if s.isReloadRateLimited(ctx) {
		return false, nil
	}

	if s.IsReloadRequested() {
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

// getConfigFromKubernetes returns the configuration that is in the configmap.
func (s *Service) getConfigFromKubernetes(ctx context.Context) (string, error) {
	s.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching configmap: %s/%s", s.configMapNamespace, s.configMapName))

	configMap, err := s.k8sClient.CoreV1().ConfigMaps(s.configMapNamespace).Get(
		s.configMapName, metav1.GetOptions{},
	)
	if err != nil {
		return "", microerror.Maskf(reloadError, err.Error())
	}

	val, ok := configMap.Data[s.configMapKey]
	if !ok {
		return "", microerror.Maskf(reloadError, "configmap key not present")
	}

	return val, nil
}

// getConfigFromPrometheus returns the configuration that is currently loaded in Prometheus.
func (s *Service) getConfigFromPrometheus(ctx context.Context) (string, error) {
	s.logger.LogCtx(ctx, "debug", "fetching current prometheus config")

	configUrl, err := s.configUrl()
	if err != nil {
		return "", microerror.Mask(err)
	}

	res, err := http.Get(configUrl)
	if err != nil {
		return "", microerror.Mask(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", microerror.Maskf(reloadError, "a non-200 HTTP status code was returned: %d", res.StatusCode)
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
			return "", microerror.Mask(err)
		}

		if resp.Status != "success" {
			return "", microerror.Maskf(reloadError, "prometheus returned non-success response status when reloding config: %s", resp.Status)
		}

		config = resp.Data.YAML
	}

	return config, nil
}

func (s *Service) reload(ctx context.Context) error {
	s.logger.LogCtx(ctx, "debug", "reloading prometheus config")

	reloadUrl, err := s.reloadUrl()
	if err != nil {
		return microerror.Mask(err)
	}

	res, err := http.Post(reloadUrl, "", nil)
	if err != nil {
		return microerror.Mask(err)
	}
	if res.StatusCode != http.StatusOK {
		return microerror.Maskf(reloadError, "a non-200 status code was returned: %d", res.StatusCode)
	}

	configurationReloadCount.Inc()

	s.lastReloadTime = time.Now()

	return nil
}

// configUrl returns the url to fetch the current Prometheus configuration.
func (s *Service) configUrl() (string, error) {
	return s.getUrl(s.address, ConfigPath)
}

// reloadUrl returns the url to reload the Prometheus configuration.
func (s *Service) reloadUrl() (string, error) {
	return s.getUrl(s.address, ReloadPath)
}

// getUrl appends the given route to the address.
func (s *Service) getUrl(address, route string) (string, error) {
	u, err := url.ParseRequestURI(s.address)
	if err != nil {
		return "", microerror.Mask(err)
	}
	u.Path = path.Join(u.Path, route)

	return u.String(), nil
}
