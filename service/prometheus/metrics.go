package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "prometheus_config_controller"
	prometheusSubsystem = "prometheus_reloader"
)

var (
	configurationReloadCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "configuration_reload_count",
			Help:      "Count of the times we have reloaded the prometheus configuration.",
		},
	)
)

func init() {
	prometheus.MustRegister(configurationReloadCount)
}
