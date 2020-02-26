package prometheus

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	v1 "k8s.io/api/core/v1"
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

	// KubernetesSDPodNameLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes pod name.
	KubernetesSDPodNameLabel = model.LabelName("__meta_kubernetes_pod_name")

	// KubernetesSDPodNodeNameLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes pod node name.
	KubernetesSDPodNodeNameLabel = model.LabelName("__meta_kubernetes_pod_node_name")

	// KubernetesSDServiceNameLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service.
	KubernetesSDServiceNameLabel = model.LabelName("__meta_kubernetes_service_name")

	// KubernetesSDServiceGiantSwarmMonitoringPresentLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service presenting the annotation giantswarm_io_monitoring.
	KubernetesSDServiceGiantSwarmMonitoringPresentLabel = model.LabelName("__meta_kubernetes_service_annotationpresent_giantswarm_io_monitoring")

	// KubernetesSDServiceGiantSwarmMonitoringLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service presenting the annotation giantswarm_io_monitoring as true.
	KubernetesSDServiceGiantSwarmMonitoringLabel = model.LabelName("__meta_kubernetes_service_annotation_giantswarm_io_monitoring")

	// KubernetesSDServiceGiantSwarmMonitoringAppTypeLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service type of managed application (default, optional).
	KubernetesSDServiceGiantSwarmMonitoringAppTypeLabel = model.LabelName("__meta_kubernetes_service_annotation_giantswarm_io_monitoring_app_type")

	// KubernetesSDServiceGiantSwarmMonitoringPathLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service path.
	KubernetesSDServiceGiantSwarmMonitoringPathLabel = model.LabelName("__meta_kubernetes_service_annotation_giantswarm_io_monitoring_path")

	// KubernetesSDServiceGiantSwarmMonitoringPathPresentLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service presenting the annotation giantswarm_io_monitoring_path.
	KubernetesSDServiceGiantSwarmMonitoringPathPresentLabel = model.LabelName("__meta_kubernetes_service_annotationpresent_giantswarm_io_monitoring_path")

	// KubernetesSDServiceGiantSwarmMonitoringPortLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service port number.
	KubernetesSDServiceGiantSwarmMonitoringPortLabel = model.LabelName("__meta_kubernetes_service_annotation_giantswarm_io_monitoring_port")

	// KubernetesSDServiceGiantSwarmMonitoringPortPresentLabel is the label applied to the target
	// by Prometheus Kubernetes service discovery that holds the target's Kubernetes service presenting the annotation giantswarm_io_monitoring_port.
	KubernetesSDServiceGiantSwarmMonitoringPortPresentLabel = model.LabelName("__meta_kubernetes_service_annotationpresent_giantswarm_io_monitoring_port")
)

// Prometheus Kubernetes metrics labels.
var (
	// MetricExportedNamespaceLabel is label for filtering by k8s namespace in kube-state-metric.
	MetricExportedNamespaceLabel = model.LabelName("exported_namespace")
	// MetricNamespaceLabel is label for filtering by k8s namespace
	MetricNamespaceLabel = model.LabelName("namespace")
	// MetricNameLabel is label for filtering by metric name.
	MetricNameLabel = model.LabelName("__name__")
	// MetricSystemdNameLabel is a label for filtering by systemd unit name.
	MetricSystemdNameLabel = model.LabelName("name")
	// MetricSystemdStateLabel is a label for filtering by systemd unit state.
	MetricSystemdStateLabel = model.LabelName("state")
	// MetricFSTypeLabel is a label for filtering by mount filesystem type.
	MetricFSTypeLabel = model.LabelName("fstype")
)

// Prometheus POD service discovery labels.
var (
	// PodSDContainerNameLabel is a label applied to the target by Prometheus
	// POD service discovery that holds the target's POD container name.
	PodSDContainerNameLabel = model.LabelName("__meta_kubernetes_pod_container_name")

	// PodSDPodNameLabel is the label applied to the target by Prometheus POD
	// service discovery that holds the target's Kubernetes POD name.
	PodSDPodNameLabel = model.LabelName("__meta_kubernetes_pod_name")

	// PodSDNamespaceLabel is the label applied to the target by Prometheus POD
	// service discovery that holds the target's Kubernetes namespace.
	PodSDNamespaceLabel = model.LabelName("__meta_kubernetes_namespace")
)

// Giant Swarm metrics schema labels.
var (
	// AddressLabel is the label used to hold target ip and port.
	AddressLabel = "__address__"

	// AppLabel is the label used to hold the application's name.
	AppLabel = "app"

	// AppTypeLabel is the label used to hold the type of managed application (optional, default), if applicable.
	AppTypeLabel = "app_type"

	// ClusterIDLabel is the label used to hold the cluster's ID.
	ClusterIDLabel = "cluster_id"

	// ClusterTypeLabel is the label used to hold the cluster's type.
	ClusterTypeLabel = "cluster_type"

	// ExportedNamespaceLabel is the label used to hold the application's namespace.
	ExportedNamespaceLabel = "exported_namespace"

	// IPLabel is the label used to hold the machine's IP.
	IPLabel = "ip"

	// NamespaceLabel is the label used to hold the application's namespace.
	NamespaceLabel = "namespace"

	// NamespaceKubeSystemLabel is the label for kube-system namespace.
	NamespaceKubeSystemLabel = "kube-system"

	// MetricPathLabel is the label used to hold the scrape metrics path.
	MetricPathLabel = "__metrics_path__"

	// PodNameLabel is the label used to hold the pod name.
	PodNameLabel = "pod_name"

	// NodeLabel is the label used to hold the node name.
	NodeLabel = "node"

	// RoleLabel is the label used to hold the machine's role.
	RoleLabel = "role"
)

// Giant Swarm metrics schema values.
const (
	// CadvisorAppName is the label value for Cadvisor targets.
	CadvisorAppName = "cadvisor"

	// CalicoNodeAppName is the label value for calico-node targets.
	CalicoNodeAppName = "calico-node"

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
	CadvisorMetricsPath = "/api/v1/nodes/${1}:10250/proxy/metrics/cadvisor"

	// NodeExporterPort is the path under which node-exporter metrics can be scraped.
	NodeExporterPort = "${1}:10300"

	// GroupCapture is the regular expression to match against the first capture group.
	GroupCapture = "${1}"
)

// Regular expressions.
var (
	// APIServerRegexp is the regular expression to match against Kubernetes API servers.
	APIServerRegexp = config.MustNewRegexp(`default;kubernetes`)

	// CalicoNodePodRegexp is the regular expression to match calico-node pod name and namespace.
	CalicoNodePodRegexp = config.MustNewRegexp(`kube-system;calico-node.*`)

	// CalicoNodePodNameRegexp is the regular expression to match calico-node pod name.
	CalicoNodePodNameRegexp = config.MustNewRegexp(`(calico-node.*)`)

	// EmptyRegexp is the regular expression to match against the empty string.
	EmptyRegexp = config.MustNewRegexp(``)

	// KubeletPortRegexp is the regular expression to match against the
	// Kubelet IP (including port), and capture the IP.
	KubeletPortRegexp = config.MustNewRegexp(`(.*):10250`)

	// NSRegexp is the regular expression to match against the specified namespaces.
	NSRegexp = config.MustNewRegexp(`(kube-system|giantswarm.*|vault-exporter)`)

	// MetricDropBucketLatencies is the regular expression to match against the several bucket latencies metrics.
	MetricDropBucketLatencies = config.MustNewRegexp(`(apiserver_admission_controller_admission_latencies_seconds_.*|apiserver_admission_step_admission_latencies_seconds_.*|apiserver_request_count|apiserver_request_duration_seconds_.*|apiserver_request_latencies_.*|apiserver_request_total|apiserver_response_sizes_.*|rest_client_request_latency_seconds_.*)`)

	// MetricDropContainerNetworkRegexp is the regular expression to match againts cadvisor container network metrics.
	MetricDropContainerNetworkRegexp = config.MustNewRegexp(`container_network_.*`)

	// MetricDropFStypeRegexp is the regular expression to match againts not interesting filesystem (for node exporter metrics).
	MetricDropFStypeRegexp = config.MustNewRegexp(`(cgroup|devpts|mqueue|nsfs|overlay|tmpfs)`)

	// MetricDropICRegexp is the regular expression to match against useless metric exposed by IC.
	MetricDropICRegexp = config.MustNewRegexp(`(ingress_controller_ssl_expire_time_seconds|nginx.*)`)

	// MetricDropSystemdStateRegexp is the regular expression to match against not interesting systemd unit (for node exporter metrics).
	MetricDropSystemdStateRegexp = config.MustNewRegexp(`node_systemd_unit_state;(active|activating|deactivating|inactive)`)

	// MetricDropSystemdNameRegexp is the regular expression to match against not interesting systemd units(docker mounts and calico network devices).
	MetricDropSystemdNameRegexp = config.MustNewRegexp(`node_systemd_unit_state;(dev-disk-by|run-docker-netns|sys-devices|sys-subsystem-net|var-lib-docker-overlay2|var-lib-docker-containers|var-lib-kubelet-pods).*`)

	// MetricsDropReflectorRegexp is the regular expression to match against spammy reflector metrics returned by the Kubelet.
	MetricsDropReflectorRegexp = config.MustNewRegexp(`(reflector.*)`)

	// ElasticLoggingPodNameRegexp is the regular expression to match elastic-logging-elasticsearch-exporter pod name.
	ElasticLoggingPodNameRegexp = config.MustNewRegexp(`(elastic-logging-elasticsearch-exporter.*)`)

	// NginxIngressControllerPodNameRegexp is the regular expression to match nginx ic pod name.
	NginxIngressControllerPodNameRegexp = config.MustNewRegexp(`(nginx-ingress-controller.*)`)

	// KubeStateMetricsPodNameRegexp is the regular expression to match kube-state-metrics pod name.
	KubeStateMetricsPodNameRegexp = config.MustNewRegexp(`(kube-state-metrics.*)`)

	// KubeStateMetricsServiceNameRegexpis the regular expression to match kube-state-metrics service name.
	KubeStateMetricsServiceNameRegexp = config.MustNewRegexp(`(kube-system;kube-state-metrics)`)

	// ChartOperatorPodNameRegexp is the regular expression to match chart-operator pod name.
	ChartOperatorPodNameRegexp = config.MustNewRegexp(`(chart-operator.*)`)

	// CertExporterPodNameRegexp is the regular expression to match cert-exporter pod name.
	CertExporterPodNameRegexp = config.MustNewRegexp(`(cert-exporter.*)`)

	// ClusterAutoscalerPodNameRegexp is the regular expression to match cluster-autoscaler pod name.
	ClusterAutoscalerPodNameRegexp = config.MustNewRegexp(`(cluster-autoscaler.*)`)

	// CoreDNSPodNameRegexp is the regular expression to match coredns pod name.
	CoreDNSPodNameRegexp = config.MustNewRegexp(`(coredns.*)`)

	// NetExporterPodNameRegexp is the regular expression to match net-exporter pod name.
	NetExporterPodNameRegexp = config.MustNewRegexp(`(net-exporter.*)`)

	// NicExporterPodNameRegexp is the regular expression to match nic-exporter pod name.
	NicExporterPodNameRegexp = config.MustNewRegexp(`(nic-exporter.*)`)

	// VaultExporterPodNameRegexp is the regular expression to match against the
	// vault-exporter name.
	VaultExporterPodNameRegexp = config.MustNewRegexp(`(vault-exporter.*)`)

	// RelabelNamespaceRegexp is the regular expression to match against metrics with empty exported_namespace and namespace kube-system.
	RelabelNamespaceRegexp = config.MustNewRegexp(`;(kube-system|giantswarm.*|vault-exporter)`)

	ManagedAppSourceRegexp = config.MustNewRegexp(`(.*);(.*);(.*);(.*)`)

	// NodeExporterRegexp is the regular expression to match against the
	// node-exporter name.
	NodeExporterRegexp = config.MustNewRegexp(`kube-system;node-exporter`)

	// NodeExporterPortRegexp is the regular expression to match against the
	// node-exporter IP (including port), and capture the IP.
	NodeExporterPortRegexp = config.MustNewRegexp(`(.*):10300`)

	// ServiceWhitelistRegexp is the regular expression to match workload targets to scrape.
	ServiceWhitelistRegexp = config.MustNewRegexp(`(kube-system;(cert-exporter|cluster-autoscaler|coredns|kube-state-metrics|net-exporter|nic-exporter|nginx-ingress-controller))|(giantswarm;chart-operator)|(giantswarm-elastic-logging;elastic-logging-elasticsearch-exporter)|(vault-exporter;vault-exporter)`)
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
