package configmap

import (
	"context"

	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentConfigMap, err := toConfigMap(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desiredConfigMap, err := toConfigMap(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// If the current or desired state configmaps are empty, perform no action.
	if currentConfigMap == nil {
		return nil, nil
	}
	if desiredConfigMap == nil {
		return nil, nil
	}

	// If the current and desired state configmaps have different names or namespaces,
	// something bad is going on, so error out.
	if currentConfigMap.Name != desiredConfigMap.Name {
		return nil, microerror.Mask(wrongNameError)
	}
	if currentConfigMap.Namespace != desiredConfigMap.Namespace {
		return nil, microerror.Mask(wrongNamespaceError)
	}

	// If the current state does not match the desired state,
	// set the desired state as update.
	if currentConfigMap.Data[r.configMapKey] != desiredConfigMap.Data[r.configMapKey] {
		r.logger.Log("debug", "current configmap does not match desired configmap, need to update")
		return desiredConfigMap, nil
	}

	r.logger.Log("debug", "current configmap matches desired configmap, no update needed")

	return nil, nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	configMapToUpdate, err := toConfigMap(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if configMapToUpdate != nil {
		r.logger.Log("debug", "updating configmap")

		timer := prometheusclient.NewTimer(kubernetesResource.WithLabelValues("configmap", "update"))
		_, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Update(configMapToUpdate)
		timer.ObserveDuration()

		if errors.IsConflict(err) {
			// fall through, we'll update it on the next reconciliation loop.
		} else if errors.IsNotFound(err) {
			return microerror.Mask(configMapNotFoundError)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	// We attempt to reload Prometheus even if the configmap hasn't updated,
	// as the PrometheusReloader takes care that we don't reload too often.
	if err := r.prometheusReloader.Reload(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func toConfigMap(v interface{}) (*v1.ConfigMap, error) {
	if v == nil {
		return nil, nil
	}

	configMap, ok := v.(*v1.ConfigMap)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", v1.ConfigMap{}, v)
	}

	return configMap, nil
}
