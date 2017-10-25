package prometheus

const (
	// prometheusReloadPath is the Prometheus API route that reloads the configuration
	// when POSTed to.
	prometheusReloadPath = "/-/reload"
)

// PrometheusReloader represents a service that can reload Prometheus configuration.
type PrometheusReloader interface {
	Reload() error
}
