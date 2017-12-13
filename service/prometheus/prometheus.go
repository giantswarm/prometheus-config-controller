package prometheus

import (
	"k8s.io/api/core/v1"
)

const (
	// ClusterAnnotation is the Kubernetes annotation that identifies Services
	// that the prometheus-config-controller should scrape.
	ClusterAnnotation = "giantswarm.io/prometheus-cluster"

	// ClusterLabel is the Prometheus label used to identify jobs
	// managed by the prometheus-config-controller.
	ClusterLabel = "prometheus_config_controller"

	// ClusterIDLabel is the Prometheus label used to identify guest cluster
	// metrics by external clients.
	ClusterIDLabel = "cluster_id"
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
