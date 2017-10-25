package configmap

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/microerror"
)

func (r *Resource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	currentConfigMap, err := toConfigMap(currentState)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	desiredConfigMap, err := toConfigMap(desiredState)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	// If the current or desired state configmaps are empty,
	// perform no action.
	if currentConfigMap == nil {
		return nil, nil, nil, nil
	}
	if desiredConfigMap == nil {
		return nil, nil, nil, nil
	}

	// If the current and desired state configmaps have different names or namespaces,
	// something bad is going on, so error out.
	if currentConfigMap.Name != desiredConfigMap.Name {
		return nil, nil, nil, microerror.Mask(wrongNameError)
	}
	if currentConfigMap.Namespace != desiredConfigMap.Namespace {
		return nil, nil, nil, microerror.Mask(wrongNamespaceError)
	}

	// If the current state does not match the desired state,
	// return the desired state as update.
	if currentConfigMap.Data[r.configMapKey] != desiredConfigMap.Data[r.configMapKey] {
		r.logger.Log("debug", "current configmap does not match desired configmap, need to update")
		return nil, nil, desiredConfigMap, nil
	}

	r.logger.Log("debug", "current configmap matches desired configmap, no update needed")

	return nil, nil, nil, nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	configMapToUpdate, err := toConfigMap(updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	if configMapToUpdate != nil {
		r.logger.Log("debug", "updating configmap")
		_, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Update(configMapToUpdate)
		if errors.IsNotFound(err) {
			return microerror.Mask(configMapNotFoundError)
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

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
