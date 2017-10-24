package certificate

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "prometheus_config_controller"
	prometheusSubsystem = "certificate_resource"
)

var (
	certificateCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "certificate_count",
			Help:      "Number of certificates on disk.",
		},
	)
)

func init() {
	prometheus.MustRegister(certificateCount)
}
