package configmap

import (
	"context"
	"fmt"

	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"
)

// GetCurrentState returns the current state of the prometheus config configmap.
// If the configmap exists, it is returned.
// If the configmap does not exist, nil is returned.
func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("debug", fmt.Sprintf("fetching configmap: %s/%s", r.configMapNamespace, r.configMapName))

	timer := prometheusclient.NewTimer(kubernetesResource.WithLabelValues("configmap", "get"))
	configMap, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Get(
		r.configMapName, metav1.GetOptions{},
	)
	timer.ObserveDuration()

	if errors.IsNotFound(err) {
		r.logger.Log("debug", "configmap does not exist")
		return nil, nil
	} else if err != nil {
		return nil, microerror.Maskf(err, "an error occurred fetching the configmap")
	}

	r.logger.Log("debug", "found configmap")

	configmapSize.Set(float64(configMap.Size()))

	return configMap, nil
}
