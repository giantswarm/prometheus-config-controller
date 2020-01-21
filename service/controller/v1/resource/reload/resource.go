package reload

import (
	"context"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

const (
	Name = "reloadv1"
)

type Config struct {
	Logger micrologger.Logger

	PrometheusAddress string
}

type Resource struct {
	logger micrologger.Logger

	prometheusAddress string
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.PrometheusAddress == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusAddress must not be empty", config)
	}

	r := &Resource{
		logger: config.Logger,

		prometheusAddress: config.PrometheusAddress,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ensure(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "reloading prometheus configuration")

	res, err := http.Post(key.PrometheusURLReload(r.prometheusAddress), "", nil)
	if err != nil {
		return microerror.Mask(err)
	}
	if res.StatusCode != http.StatusOK {
		return microerror.Maskf(executionFailedError, "non-200 status code = %d was returned", res.StatusCode)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "reloaded prometheus configuration")
	return nil
}
