package configmap

import (
	"context"
	"fmt"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/prometheus/config"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

func (r *Resource) getDesiredState(ctx context.Context) (*corev1.ConfigMap, error) {
	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching configmap: %s/%s", r.configMapNamespace, r.configMapName))

	configMap, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Get(
		r.configMapName, metav1.GetOptions{},
	)

	if errors.IsNotFound(err) {
		return nil, microerror.Maskf(configMapNotFoundError, "%s/%s", r.configMapNamespace, r.configMapName)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	configMapData, ok := configMap.Data[r.configMapKey]
	if !ok {
		return nil, microerror.Maskf(configMapKeyNotFoundError, "%s/%s - %s", r.configMapNamespace, r.configMapName, r.configMapKey)
	}

	prometheusConfig, err := config.Load(configMapData)
	if err != nil {
		return nil, microerror.Maskf(invalidConfigMapError, err.Error())
	}

	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching all services"))

	services, err := r.k8sClient.CoreV1().Services("").List(metav1.ListOptions{
		LabelSelector: key.LabelSelectorService().String(),
	})

	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("computing desired state of configmap"))
	scrapeConfigs, err := prometheus.GetScrapeConfigs(services.Items, r.certDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	newPrometheusConfig, err := prometheus.UpdateConfig(*prometheusConfig, scrapeConfigs)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Prometheus YAML marshalling obscures remote write passwords for security,
	// so we store the Cortex password, and write it back after marshalling.
	var remoteWritePassword string
	if len(newPrometheusConfig.RemoteWriteConfigs) == 1 {
		remoteWritePassword = string(newPrometheusConfig.RemoteWriteConfigs[0].HTTPClientConfig.BasicAuth.Password)
	}

	newConfigMapData, err := yaml.Marshal(newPrometheusConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if remoteWritePassword != "" && strings.Count(string(newConfigMapData), "password: <secret>") == 1 {
		newConfigMapData = []byte(
			strings.Replace(string(newConfigMapData), "password: <secret>", fmt.Sprintf("password: %s", remoteWritePassword), 1),
		)
	}

	configMap.Data[r.configMapKey] = string(newConfigMapData)

	return configMap, nil
}
