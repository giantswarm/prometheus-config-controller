package prometheus

import (
	"github.com/prometheus/prometheus/config"
	"k8s.io/api/core/v1"
)

const (
	// ClusterAnnotation is the Kubernetes annotation that identifies Services
	// that the prometheus-config-controller should scrape.
	ClusterAnnotation = "giantswarm.io/prometheus-cluster"

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

	// PrometheusServicePortLabel is the Prometheus label added by the Kubernetes
	// service discovery to hold an endpoints port.
	PrometheusServicePortLabel = "__meta_kubernetes_pod_container_port_number"

	// PrometheusServicePortLabel is the Prometheus label added by the Kubernetes
	// service discovery to hold a node's name.
	PrometheusKubernetesNodeNameLabel = "__meta_kubernetes_node_name"

	// CadvisorMetricsPath is the path under which cadvisor metrics can be scraped.
	CadvisorMetricsPath = "/api/v1/nodes/${1}:4194/proxy/metrics"
)

var (
	// EndpointRegexp is the regular expression against which endpoint service names
	// must match to be scraped.
	// The empty string is also matched, so that nodes (which have no service name),
	// are also matched.
	EndpointRegexp = config.MustNewRegexp(`(\s*|kube-state-metrics|kubernetes|node-exporter)`)

	// EndpointPortRegexp is the regular expression against which endpoint service ports
	// must match to be scraped.
	// We specify ports in the cases where users are running kube-state-metrics or node-exporters,
	// to ensure we only scrape the correct instances.
	EndpointPortRegexp = config.MustNewRegexp(`(\s*|443|10300|10301)`)

	// HTTPEndpointRegexp is the regular expression against which endpoint service
	// names that we want to scrape via HTTP need to match.
	HTTPEndpointRegexp = config.MustNewRegexp(`(kube-state-metrics|node-exporter)`)

	// GroupRegex is the regular expression with which we can group strings.
	GroupRegex = config.MustNewRegexp(`(.+)`)
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
