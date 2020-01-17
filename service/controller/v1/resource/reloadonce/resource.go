package reloadonce

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller/context/finalizerskeptcontext"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

const (
	Name = "reloadoncev1"
)

type Config struct {
	Logger             micrologger.Logger
	PrometheusReloader prometheus.PrometheusReloader
}

type Resource struct {
	logger             micrologger.Logger
	prometheusReloader prometheus.PrometheusReloader

	reloaded bool
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.PrometheusReloader == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.PrometheusReloader must not be empty")
	}

	r := &Resource{
		logger:             config.Logger,
		prometheusReloader: config.PrometheusReloader,

		reloaded: false,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

// ensure reloads Prometheus configuration once. This is required at startup
// when configuration ConfigMap doesn't change but it isn't loaded to
// Prometheus. It can't be done at the boot time because prometheus may be not
// ready yet.
func (r *Resource) ensure(ctx context.Context, obj interface{}) error {
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "reloading prometheus configuration")

		if r.reloaded {
			r.logger.LogCtx(ctx, "level", "debug", "message", "prometheus configuration already reloaded")

			r.logger.LogCtx(ctx, "level", "debug", "message", "cancelling resource")
			return nil
		}

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
			r.reloaded = true
			r.logger.LogCtx(ctx, "level", "debug", "message", "reloaded prometheus configuration")
		}
	}

	return nil
}
