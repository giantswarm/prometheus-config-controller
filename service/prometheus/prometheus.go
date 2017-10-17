package prometheus

import (
	"k8s.io/client-go/pkg/api/v1"
)

const (
	ClusterAnnotation = "giantswarm.io/prometheus-config-controller/cluster"
)

// GetGroupName returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetGroupName(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
