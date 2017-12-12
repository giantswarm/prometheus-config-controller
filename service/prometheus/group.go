package prometheus

import (
	"k8s.io/api/core/v1"
)

// FilterInvalidServices takes a list of Kubernetes Services,
// and returns a list of valid Services.
func FilterInvalidServices(services []v1.Service) []v1.Service {
	filteredServices := []v1.Service{}

	for _, service := range services {
		if _, ok := service.ObjectMeta.Annotations[ClusterAnnotation]; !ok {
			continue
		}

		filteredServices = append(filteredServices, service)
	}

	return filteredServices
}

// GroupServices takes a list of Kubernetes Services,
// and returns a map of the Kubernetes Services, with the cluster annotation
// value acting as the key.
// Services that do not specify a cluster annotation are dropped.
func GroupServices(services []v1.Service) map[string][]v1.Service {
	groupedServices := map[string][]v1.Service{}

	for _, service := range services {
		// If the service does not specify a cluster annotation, drop it.
		if _, ok := service.ObjectMeta.Annotations[ClusterAnnotation]; !ok {
			continue
		}

		// If the group doesn't exist yet in the map, create the empty list.
		cluster := service.ObjectMeta.Annotations[ClusterAnnotation]
		if _, ok := groupedServices[cluster]; !ok {
			groupedServices[cluster] = []v1.Service{}
		}

		// And add the new service to the group.
		groupedServices[cluster] = append(groupedServices[cluster], service)
	}

	return groupedServices
}
