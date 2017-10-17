package prometheus

import (
	"strings"

	"k8s.io/client-go/pkg/api/v1"
)

const (
	CertificateAnnotation = "giantswarm.io/prometheus-config-controller/certificate"
	ClusterAnnotation     = "giantswarm.io/prometheus-config-controller/cluster"
)

// GetCertificateName returns the value of the certificate annotation, splitted by '/'.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetCertificateName(service v1.Service) (string, string) {
	reference := service.ObjectMeta.Annotations[CertificateAnnotation]
	parts := strings.Split(reference, "/")

	return parts[0], parts[1]
}

// GetGroupName returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetGroupName(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
