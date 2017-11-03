package configmap

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "prometheus_config_controller"
	prometheusSubsystem = "configmap_resource"
)

var (
	configmapSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "configmap_size",
			Help:      "Size of the prometheus configmap.",
		},
	)

	kubernetesResource = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prometheusNamespace,
			Name:      "kubernetes_resource",
			Help:      "Histogram for actions against Kubernetes resources.",
		},
		[]string{"resource", "action"},
	)

	scrapeConfigCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Subsystem: prometheusSubsystem,
			Name:      "scrape_config_count",
			Help:      "Count of scrape configs computed.",
		},
	)
)

func init() {
	prometheus.MustRegister(configmapSize)
	prometheus.MustRegister(kubernetesResource)
	prometheus.MustRegister(scrapeConfigCount)
}
