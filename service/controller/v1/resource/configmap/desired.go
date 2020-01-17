package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/prometheus/config"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching configmap: %s/%s", r.configMapNamespace, r.configMapName))

	configMap, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Get(
		r.configMapName, metav1.GetOptions{},
	)

	if errors.IsNotFound(err) {
		return nil, microerror.Maskf(configMapNotFoundError, "%s/%s", r.configMapNamespace, r.configMapName)
	} else if err != nil {
		return nil, microerror.Maskf(err, "an error occurred fetching the configmap %s/%s", r.configMapNamespace, r.configMapName)
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
		LabelSelector: key.ServiceLabelSelector().String(),
	})

	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred listing all services")
	}

	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("computing desired state of configmap"))
	scrapeConfigs, err := prometheus.GetScrapeConfigs(services.Items, r.certDirectory)
	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred creating scrape configs")
	}

	newPrometheusConfig, err := prometheus.UpdateConfig(*prometheusConfig, scrapeConfigs)
	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred merging prometheus config")
	}

	newConfigMapData, err := yaml.Marshal(newPrometheusConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "an error occurred marshaling yaml")
	}

	configMap.Data[r.configMapKey] = string(newConfigMapData)

	return configMap, nil
}
