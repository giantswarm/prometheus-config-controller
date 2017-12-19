package prometheus

import (
	"github.com/prometheus/prometheus/config"
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

	// NamespaceLabel is the Prometheus label we use internally for an
	// endpoints namespace.
	NameLabel = "kubernetes_name"

	// NamespaceLabel is the Prometheus label we use internally for
	// an endpoints namespace.
	NamespaceLabel = "kubernetes_namespace"

	// PrometheusNamespaceLabel is the Prometheus label added by the Kubernetes
	// service discovery to hold an endpoint targets namespace.
	PrometheusNamespaceLabel = "__meta_kubernetes_namespace"

	// PrometheusServiceNameLabel is the Prometheus label added by the Kubernetes
	// service discovery to hold an endpoints service name.
	PrometheusServiceNameLabel = "__meta_kubernetes_service_name"
)

var (
	// EndpointRegexp is the regular expression against which endpoint service names
	// must match to be scraped.
	// The empty string is also matched, so that nodes (which have no service name),
	// are also matched.
	EndpointRegexp = config.MustNewRegexp(`(\s*|kubernetes|node-exporter)`)

	// HTTPEndpointRegexp is the regular expression against which endpoint service
	// names that we want to scrape via HTTP need to match.
	HTTPEndpointRegexp = config.MustNewRegexp(`(node-exporter)`)
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
