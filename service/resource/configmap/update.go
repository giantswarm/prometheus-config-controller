package configmap

import (
	"context"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
)

func (r *Resource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	currentStateConfigMap, ok := currentState.(*v1.ConfigMap)
	if !ok {
		return nil, nil, nil, configMapAssertionError
	}

	desiredStateConfigMap, ok := desiredState.(*v1.ConfigMap)
	if !ok {
		return nil, nil, nil, configMapAssertionError
	}

	// If the current or desired state configmaps are empty,
	// perform no action.
	if currentStateConfigMap == nil {
		return nil, nil, nil, nil
	}
	if desiredStateConfigMap == nil {
		return nil, nil, nil, nil
	}

	// If the current and desired state configmaps have different names or namespaces,
	// something bad is going on, so error out.
	if currentStateConfigMap.Name != desiredStateConfigMap.Name {
		return nil, nil, nil, configMapWrongNameError
	}
	if currentStateConfigMap.Namespace != desiredStateConfigMap.Namespace {
		return nil, nil, nil, configMapWrongNamespaceError
	}

	// If the current state does not match the desired state,
	// return the desired state as update.
	if currentStateConfigMap.Data[r.configMapKey] != desiredStateConfigMap.Data[r.configMapKey] {
		return nil, desiredStateConfigMap, nil, nil
	}

	return nil, nil, nil, nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	updateStateConfigmap, ok := updateState.(*v1.ConfigMap)
	if !ok {
		return configMapAssertionError
	}

	if updateStateConfigmap != nil {
		_, err := r.k8sClient.CoreV1().ConfigMaps(r.configMapNamespace).Update(updateStateConfigmap)
		if errors.IsNotFound(err) {
			return configMapNotFoundError
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
