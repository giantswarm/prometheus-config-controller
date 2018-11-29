package prometheus

import (
	"k8s.io/api/core/v1"
)

// FilterInvalidServices takes a list of Kubernetes Services,
// and returns a list of valid Services.
func FilterInvalidServices(services []v1.Service) []v1.Service {
	filteredServices := []v1.Service{}

	for _, service := range services {
		{
			isDeleted := service.GetDeletionTimestamp() != nil
			if isDeleted {
				continue
			}
		}

		{
			_, hasAnnotation := service.ObjectMeta.Annotations[ClusterAnnotation]
			if !hasAnnotation {
				continue
			}
		}

		filteredServices = append(filteredServices, service)
	}

	return filteredServices
}
