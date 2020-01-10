package configmap

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"

	"github.com/giantswarm/operatorkit/controller/context/finalizerskeptcontext"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()

	if update != nil {
		patch.SetUpdateChange(update)
	}

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
		r.logger.LogCtx(ctx, "debug", "current configmap does not match desired configmap, need to update")
		return desiredConfigMap, nil
	}

	r.logger.LogCtx(ctx, "debug", "current configmap matches desired configmap, no update needed")

	return nil, nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	configMapToUpdate, err := toConfigMap(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if configMapToUpdate != nil {
		r.logger.LogCtx(ctx, "debug", "updating configmap")

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

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "reloading prometheus configuration")
		// We attempt to reload Prometheus even if the configmap hasn't updated,
		// as the PrometheusReloader takes care that we don't reload too often.
		err := r.prometheusReloader.Reload(ctx)
		if prometheus.IsReloadThrottle(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not reload prometheus configuration")

			r.logger.LogCtx(ctx, "level", "debug", "message", err.Error())
			r.logger.LogCtx(ctx, "level", "debug", "message", "keeping finalizers")
			finalizerskeptcontext.SetKept(ctx)
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "reloaded prometheus configuration")
		}
	}

	return nil
}

func toConfigMap(v interface{}) (*v1.ConfigMap, error) {
	if v == nil {
		return nil, nil
	}

	configMap, ok := v.(*v1.ConfigMap)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", configMap, v)
	}

	return configMap, nil
}
