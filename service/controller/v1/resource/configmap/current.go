package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCurrentState returns the current state of the prometheus config configmap.
// If the configmap exists, it is returned.
// If the configmap does not exist, nil is returned.
func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.LogCtx(ctx, "debug", fmt.Sprintf("fetching configmap: %s/%s", r.configMapNamespace, r.configMapName))

	configMap, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Get(
		r.configMapName, metav1.GetOptions{},
	)

	if errors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "debug", "configmap does not exist")
		return nil, nil
	} else if err != nil {
		return nil, microerror.Maskf(err, "an error occurred fetching the configmap")
	}

	r.logger.LogCtx(ctx, "debug", "found configmap")

	return configMap, nil
}
