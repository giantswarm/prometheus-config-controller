package prometheus

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	Address string
	// ConfigMapKey is the key in the configmap under which the prometheus configuration is held.
	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,

		Address:            "",
		ConfigMapKey:       "",
		ConfigMapName:      "",
		ConfigMapNamespace: "",
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

	service := &Service{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		address:            config.Address,
		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
	}

	return service, nil
}

func (s *Service) Reload() error {
	reloadRequired, err := s.isReloadRequired()
	if err != nil {
		return microerror.Maskf(reloadError, err.Error())
	}

	if reloadRequired {
		if err := s.reload(); err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *Service) isReloadRequired() (bool, error) {
	s.logger.Log("debug", "checking if reload is required")

	configurationReloadCheckCount.Inc()

	kubernetesConfiguration, err := s.getConfigFromKubernetes()
	if err != nil {
		return false, microerror.Mask(err)
	}

	prometheusConfiguration, err := s.getConfigFromPrometheus()
	if err != nil {
		return false, microerror.Mask(err)
	}

	if kubernetesConfiguration != prometheusConfiguration {
		configurationReloadRequiredCount.Inc()

		s.logger.Log("debug", "kubernetes and prometheus configuration do not match, reload required")
		return true, nil
	}

	s.logger.Log("debug", "kubernetes and prometheus configuration match, reload not required")
	return false, nil
}

// getConfigFromKubernetes returns the configuration that is in the configmap.
func (s *Service) getConfigFromKubernetes() (string, error) {
	s.logger.Log("debug", fmt.Sprintf("fetching configmap: %s/%s", s.configMapNamespace, s.configMapName))

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
func (s *Service) getConfigFromPrometheus() (string, error) {
	s.logger.Log("debug", "fetching current prometheus config")

	configUrl, err := s.configUrl()
	if err != nil {
		return "", microerror.Mask(err)
	}

	res, err := http.Get(configUrl)
	if err != nil {
		return "", microerror.Mask(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", microerror.Maskf(reloadError, "a non-200 status code was returned: %d", res.StatusCode)
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", microerror.Mask(err)
	}
	defer res.Body.Close()

	configPage := html.UnescapeString(string(buf))

	startAnchor := "<pre>"
	endAnchor := "</pre>"

	if !strings.Contains(configPage, startAnchor) && !strings.Contains(configPage, endAnchor) {
		return "", microerror.Maskf(reloadError, "required start and end anchors not found in configpage")
	}

	i := strings.Split(configPage, startAnchor)
	j := strings.Split(i[1], endAnchor)

	config := j[0]

	return config, nil
}

func (s *Service) reload() error {
	s.logger.Log("debug", "reloading prometheus config")

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

	return nil
}

// configUrl returns the url to fetch the current Prometheus configuration.
func (s *Service) configUrl() (string, error) {
	return s.getUrl(s.address, prometheusConfigPath)
}

// reloadUrl returns the url to reload the Prometheus configuration.
func (s *Service) reloadUrl() (string, error) {
	return s.getUrl(s.address, prometheusReloadPath)
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
