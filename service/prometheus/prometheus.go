package prometheus

import (
	"k8s.io/client-go/pkg/api/v1"
)

const (
	ClusterAnnotation = "giantswarm.io/prometheus-cluster"

	// ClusterLabel is the Prometheus label used to identify jobs
	// managed by the prometheus-config-controller.
	ClusterLabel = "prometheus_config_controller"
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
