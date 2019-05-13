package prometheus

import "context"

// PrometheusReloader represents a service that can reload Prometheus configuration.
type PrometheusReloader interface {
	// Reload should reload the Prometheus configuration, possibly taking
	// rate limiting into account.
	Reload(ctx context.Context) error

	// RequestReload should specify that the next call to Reload should force
	// the reload to happen. Rate-limiting is allowed, but a reload must happen
	// eventually.
	RequestReload(ctx context.Context)
}
