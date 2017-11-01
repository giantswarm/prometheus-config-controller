package prometheus

const (
	// prometheusConfigPath is the Prometheus route that returns the current
	// configuration webpage.
	prometheusConfigPath = "/config"
	// prometheusReloadPath is the Prometheus API route that reloads the configuration
	// when POSTed to.
	prometheusReloadPath = "/-/reload"
)

// PrometheusReloader represents a service that can reload Prometheus configuration.
type PrometheusReloader interface {
	Reload() error
}
