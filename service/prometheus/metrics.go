package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "prometheus_config_controller"
	prometheusSubsystem = "prometheus_reloader"
)

var (
	configurationReloadCheckCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "configuration_reload_check_count",
			Help:      "Count of the times we have checked if a reload of the prometheus configuration is necessary.",
		},
	)

	configurationReloadIgnoredCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "configuration_reload_ignored_count",
			Help:      "Count of the times we have ignored a reload request due to rate limiting.",
		},
	)

	configurationReloadRequiredCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "configuration_reload_required_count",
			Help:      "Count of the times we need to reload the prometheus configuration.",
		},
	)

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
	prometheus.MustRegister(configurationReloadCheckCount)
	prometheus.MustRegister(configurationReloadIgnoredCount)
	prometheus.MustRegister(configurationReloadRequiredCount)
	prometheus.MustRegister(configurationReloadCount)
}
