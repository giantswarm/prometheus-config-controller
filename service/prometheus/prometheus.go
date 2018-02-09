package prometheus

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"k8s.io/api/core/v1"
)

const (
	// ClusterAnnotation is the Kubernetes annotation that identifies Services
	// that the prometheus-config-controller should scrape.
	ClusterAnnotation = "giantswarm.io/prometheus-cluster"
)

// Prometheus Kubernetes service discovery labels.
var (
	// KubernetesSDNamespaceLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes namespace.
	KubernetesSDNamespaceLabel = model.LabelName("__meta_kubernetes_namespace")

	// KubernetesSDNodeNameLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes node name.
	KubernetesSDNodeNameLabel = model.LabelName("__meta_kubernetes_node_name")

	// KubernetesSDNodeAddressInternalIPLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes node addres internal IP.
	KubernetesSDNodeAddressInternalIPLabel = model.LabelName("__meta_kubernetes_node_address_InternalIP")

	// KubernetesSDNodeLabelRole is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes node role label.
	KubernetesSDNodeLabelRole = model.LabelName("__meta_kubernetes_node_label_role")

	// KubernetesSDServiceNameLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service.
	KubernetesSDServiceNameLabel = model.LabelName("__meta_kubernetes_service_name")
)

// Giant Swarm metrics schema labels.
var (
	// AppLabel is the label used to hold the application's name.
	AppLabel = "app"

	// ClusterIDLabel is the label used to hold the cluster's ID.
	ClusterIDLabel = "cluster_id"

	// ClusterTypeLabel is the label used to hold the cluster's type.
	ClusterTypeLabel = "cluster_type"

	// IPLabel is the label used to hold the machine's IP.
	IPLabel = "ip"

	// NamespaceLabel is the label used to hold the application's namespace.
	NamespaceLabel = "namespace"

	// RoleLabel is the label used to hold the machine's role.
	RoleLabel = "role"
)

// Giant Swarm metrics schema values.
const (
	// CadvisorAppName is the label value for Cadvisor targets.
	CadvisorAppName = "cadvisor"

	// GuestClusterType is the cluster type for guest clusters.
	GuestClusterType = "guest"

	// KubeletAppName is the label value for kubelets.
	KubeletAppName = "kubelet"

	// KubernetesAppName is the label value for Kubernetes API servers.
	KubernetesAppName = "kubernetes"

	// NodeExporterAppName is the label value for node-exporters.
	NodeExporterAppName = "node-exporter"

	// WorkerRole is the label value used for Kubernetes workers.
	WorkerRole = "worker"
)

// Path replacements.
const (
	// CadvisorMetricsPath is the path under which cadvisor metrics can be scraped.
	CadvisorMetricsPath = "/api/v1/nodes/${1}:4194/proxy/metrics"

	// NodeExporterPort is the path under which node-exporter metrics can be scraped.
	NodeExporterPort = "${1}:10300"

	// GroupCapture is the regular expression to match against the first capture group.
	GroupCapture = "${1}"
)

// Regular expressions.
var (
	// APIServerRegexp is the regular expression to match against Kubernetes API servers.
	APIServerRegexp = config.MustNewRegexp(`default;kubernetes`)

	// EmptyRegexp is the regular expression to match against the empty string.
	EmptyRegexp = config.MustNewRegexp(``)

	// KubeletPortRegexp is the regular expression to match against the
	// Kubelet IP (including port), and capture the IP.
	KubeletPortRegexp = config.MustNewRegexp(`(.*):10250`)

	// NodeExporterRegexp is the regular expression to match against the
	// node-exporter name.
	NodeExporterRegexp = config.MustNewRegexp(`kube-system;node-exporter`)

	// NodeExporterPortRegexp is the regular expression to match against the
	// node-exporter IP (including port), and capture the IP.
	NodeExporterPortRegexp = config.MustNewRegexp(`(.*):10300`)

	// WhitelistRegexp is the regular expression to match workload targets to scrape.
	WhitelistRegexp = config.MustNewRegexp(`kube-system;kube-state-metrics|node-exporter`)
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
