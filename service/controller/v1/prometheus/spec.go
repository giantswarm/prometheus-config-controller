package prometheus

const (
	// prometheusConfigPath is the Prometheus route that returns the current
	// configuration webpage.
	prometheusConfigPath = "/api/v1/status/config"
	// prometheusReloadPath is the Prometheus API route that reloads the configuration
	// when POSTed to.
	prometheusReloadPath = "/-/reload"
)

// PrometheusReloader represents a service that can reload Prometheus configuration.
type PrometheusReloader interface {
	// Reload should reload the Prometheus configuration, possibly taking
	// rate limiting into account.
	Reload() error

	// RequestReload should specify that the next call to Reload should force
	// the reload to happen. Rate-limiting is allowed, but a reload must happen
	// eventually.
	RequestReload()
}
