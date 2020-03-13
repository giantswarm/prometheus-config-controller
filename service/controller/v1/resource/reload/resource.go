package reload

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

const (
	Name = "reloadv1"

	// minReloadInterval is the minimum time that has to pass between
	// Prometheus reload calls unless ConfigMap resource version changes.
	minReloadInterval = 2 * time.Minute
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	ConfigMapKey       string
	ConfigMapName      string
	ConfigMapNamespace string
	ConfigMapPath      string
	PrometheusAddress  string
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	lastReloadTime                   time.Time
	lastSeenConfigMapResourceVersion string

	configMapKey       string
	configMapName      string
	configMapNamespace string
	configMapPath      string
	prometheusAddress  string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.ConfigMapKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapKey must not be empty", config)
	}
	if config.ConfigMapName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapName must not be empty", config)
	}
	if config.ConfigMapNamespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapNamespace must not be empty", config)
	}
	if config.ConfigMapPath == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ConfigMapPath must not be empty", config)
	}
	if config.PrometheusAddress == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusAddress must not be empty", config)
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		lastReloadTime: time.Now().Add(-minReloadInterval),

		configMapKey:       config.ConfigMapKey,
		configMapName:      config.ConfigMapName,
		configMapNamespace: config.ConfigMapNamespace,
		configMapPath:      config.ConfigMapPath,
		prometheusAddress:  config.PrometheusAddress,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ensure(ctx context.Context, obj interface{}) error {
	var cm *corev1.ConfigMap
	{
		filePath := path.Join(r.configMapPath, r.configMapKey)
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding ConfigMap at %#q", filePath))

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return microerror.Mask(err)
		}

		var configMap corev1.ConfigMap
		err = yaml.Unmarshal(content, &configMap)
		if err != nil {
			return microerror.Mask(err)
		}
		cm = &configMap

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found %#q ConfigMap in namespace %#q", r.configMapName, r.configMapNamespace))
	}

	var reloadRequired bool
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "finding out if prometheus configuration needs to be reloaded")

		switch {
		case cm.ResourceVersion != r.lastSeenConfigMapResourceVersion:
			reloadRequired = true
		}

		if reloadRequired {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found out that prometheus configuration needs to be reloaded")
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found out that prometheus configuration does not need to be reloaded")

			r.logger.LogCtx(ctx, "level", "debug", "message", "canceling resource")
			return nil
		}
	}

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "reloading prometheus configuration")

		res, err := http.Post(key.PrometheusURLReload(r.prometheusAddress), "", nil)
		if err != nil {
			return microerror.Mask(err)
		}
		if res.StatusCode != http.StatusOK {
			return microerror.Maskf(executionFailedError, "non-200 status code = %d was returned", res.StatusCode)
		}

		r.lastSeenConfigMapResourceVersion = cm.ResourceVersion
		r.lastReloadTime = time.Now()

		r.logger.LogCtx(ctx, "level", "debug", "message", "reloaded prometheus configuration")
	}

	return nil
}
