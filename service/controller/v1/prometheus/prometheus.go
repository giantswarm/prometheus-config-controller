package prometheus

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/relabel"
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
	// DeploymentTypeLabel is a label added by kube-state-metrics to Deployment related metrics.
	DeploymentTypeLabel = model.LabelName("deployment")
	// DaemonSetTypeLabel is a label added by kube-state-metrics to DaemonSet related metrics.
	DaemonSetTypeLabel = model.LabelName("daemonset")
	// StatefulSetTypeLabel is a label added by kube-state-metrics to StatefulSet related metrics.
	StatefulSetTypeLabel = model.LabelName("statefulset")
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

	// AppIsManaged is the label used to mark metrics coming from managed apps marked with "giantswarm.io/monitoring"
	// k8s annotation.
	AppIsManaged = "is_managed_app"

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

	// KubeStateMetricsForManagedApps is the label used to mark metrics coming from kube-state-metrics
	// for use with managed apps. They are used to show the basic "is deployment OK" metric.
	KubeStateMetricsForManagedApps = "kube_state_metrics_for_managed_app"

	// NamespaceLabel is the label used to hold the application's namespace.
	NamespaceLabel = "namespace"

	// NamespaceKubeSystemLabel is the label for kube-system namespace.
	NamespaceKubeSystemLabel = "kube-system"

	// ManagedAppWorkloadTypeLabel is the label for showing the workload type (deployment, statefulset, daemonset)
	// of a managed app.
	ManagedAppWorkloadTypeLabel = "workload_type"

	// ManagedAppWorkloadNameLabel is the label for storing any workload's name.
	ManagedAppWorkloadNameLabel = "workload_name"

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

	// ManagedAppsDeployment is the value used to indicate a managed app workload of type Deployment.
	ManagedAppsDeployment = "deployment"

	// ManagedAppsDaemonSet is the value used to indicate a managed app workload of type DaemonSet.
	ManagedAppsDaemonSet = "daemonset"

	// ManagedAppsStatefulSet is the value used to indicate a managed app workload of type StatefulSet.
	ManagedAppsStatefulSet = "statefulset"
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
	APIServerRegexp = relabel.MustNewRegexp(`default;kubernetes`)

	// CalicoNodePodRegexp is the regular expression to match calico-node pod name and namespace.
	CalicoNodePodRegexp = relabel.MustNewRegexp(`kube-system;calico-node.*`)

	// CalicoNodePodNameRegexp is the regular expression to match calico-node pod name.
	CalicoNodePodNameRegexp = relabel.MustNewRegexp(`(calico-node.*)`)

	// EmptyRegexp is the regular expression to match against the empty string.
	EmptyRegexp = relabel.MustNewRegexp(``)

	// NonEmptyRegexp is the regular expression to match against the non-empty string.
	NonEmptyRegexp = relabel.MustNewRegexp(`(.+)`)

	// KubeletPortRegexp is the regular expression to match against the
	// Kubelet IP (including port), and capture the IP.
	KubeletPortRegexp = relabel.MustNewRegexp(`(.*):10250`)

	// NSRegexp is the regular expression to match against the specified namespaces.
	NSRegexp = relabel.MustNewRegexp(`(kube-system|giantswarm.*|vault-exporter)`)

	// MetricDropBucketLatencies is the regular expression to match against the several bucket latencies metrics.
	MetricDropBucketLatencies = relabel.MustNewRegexp(`(apiserver_admission_controller_admission_latencies_seconds_.*|apiserver_admission_step_admission_latencies_seconds_.*|apiserver_request_count|apiserver_request_duration_seconds_.*|apiserver_request_latencies_.*|apiserver_request_total|apiserver_response_sizes_.*|rest_client_request_latency_seconds_.*)`)

	// MetricDropContainerNetworkRegexp is the regular expression to match againts cadvisor container network metrics.
	MetricDropContainerNetworkRegexp = relabel.MustNewRegexp(`container_network_.*`)

	// MetricDropFStypeRegexp is the regular expression to match againts not interesting filesystem (for node exporter metrics).
	MetricDropFStypeRegexp = relabel.MustNewRegexp(`(cgroup|devpts|mqueue|nsfs|overlay|tmpfs)`)

	// MetricDropICRegexp is the regular expression to match against useless metric exposed by IC.
	MetricDropICRegexp = relabel.MustNewRegexp(`(ingress_controller_ssl_expire_time_seconds|nginx.*)`)

	// MetricDropSystemdStateRegexp is the regular expression to match against not interesting systemd unit (for node exporter metrics).
	MetricDropSystemdStateRegexp = relabel.MustNewRegexp(`node_systemd_unit_state;(active|activating|deactivating|inactive)`)

	// MetricDropSystemdNameRegexp is the regular expression to match against not interesting systemd units(docker mounts and calico network devices).
	MetricDropSystemdNameRegexp = relabel.MustNewRegexp(`node_systemd_unit_state;(dev-disk-by|run-docker-netns|sys-devices|sys-subsystem-net|var-lib-docker-overlay2|var-lib-docker-containers|var-lib-kubelet-pods).*`)

	// MetricsDropReflectorRegexp is the regular expression to match against spammy reflector metrics returned by the Kubelet.
	MetricsDropReflectorRegexp = relabel.MustNewRegexp(`(reflector.*)`)

	// ElasticLoggingPodNameRegexp is the regular expression to match elastic-logging-elasticsearch-exporter pod name.
	ElasticLoggingPodNameRegexp = relabel.MustNewRegexp(`(elastic-logging-elasticsearch-exporter.*)`)

	// NginxIngressControllerPodNameRegexp is the regular expression to match nginx ic pod name.
	NginxIngressControllerPodNameRegexp = relabel.MustNewRegexp(`(nginx-ingress-controller.*)`)

	// KubeStateMetricsPodNameRegexp is the regular expression to match kube-state-metrics pod name.
	KubeStateMetricsPodNameRegexp = relabel.MustNewRegexp(`(kube-state-metrics.*)`)

	// KubeStateMetricsServiceNameRegexpis the regular expression to match kube-state-metrics service name.
	KubeStateMetricsServiceNameRegexp = relabel.MustNewRegexp(`(kube-system;kube-state-metrics)`)

	// KubeStateMetricsManagedAppMetricsNameRegexp is the regular expression to keep only KSM metrics realted to SLI of managed apps.
	KubeStateMetricsManagedAppMetricsNameRegexp = relabel.MustNewRegexp(`(kube_deployment_status_replicas_unavailable|kube_deployment_labels|kube_daemonset_status_number_unavailable|kube_daemonset_labels|kube_statefulset_status_replicas|kube_statefulset_status_replicas_current|kube_statefulset_labels)`)

	// ChartOperatorPodNameRegexp is the regular expression to match chart-operator pod name.
	ChartOperatorPodNameRegexp = relabel.MustNewRegexp(`(chart-operator.*)`)

	// CertExporterPodNameRegexp is the regular expression to match cert-exporter pod name.
	CertExporterPodNameRegexp = relabel.MustNewRegexp(`(cert-exporter.*)`)

	// ClusterAutoscalerPodNameRegexp is the regular expression to match cluster-autoscaler pod name.
	ClusterAutoscalerPodNameRegexp = relabel.MustNewRegexp(`(cluster-autoscaler.*)`)

	// CoreDNSPodNameRegexp is the regular expression to match coredns pod name.
	CoreDNSPodNameRegexp = relabel.MustNewRegexp(`(coredns.*)`)

	// NetExporterPodNameRegexp is the regular expression to match net-exporter pod name.
	NetExporterPodNameRegexp = relabel.MustNewRegexp(`(net-exporter.*)`)

	// NicExporterPodNameRegexp is the regular expression to match nic-exporter pod name.
	NicExporterPodNameRegexp = relabel.MustNewRegexp(`(nic-exporter.*)`)

	// VaultExporterPodNameRegexp is the regular expression to match against the
	// vault-exporter name.
	VaultExporterPodNameRegexp = relabel.MustNewRegexp(`(vault-exporter.*)`)

	// RelabelNamespaceRegexp is the regular expression to match against metrics with empty exported_namespace and namespace kube-system.
	RelabelNamespaceRegexp = relabel.MustNewRegexp(`;(kube-system|giantswarm.*|vault-exporter)`)

	ManagedAppSourceRegexp = relabel.MustNewRegexp(`(.*);(.*);(.*);(.*)`)

	// NodeExporterRegexp is the regular expression to match against the
	// node-exporter name.
	NodeExporterRegexp = relabel.MustNewRegexp(`kube-system;node-exporter`)

	// NodeExporterPortRegexp is the regular expression to match against the
	// node-exporter IP (including port), and capture the IP.
	NodeExporterPortRegexp = relabel.MustNewRegexp(`(.*):10300`)

	// ServiceWhitelistRegexp is the regular expression to match workload targets to scrape.
	ServiceWhitelistRegexp = relabel.MustNewRegexp(`(kube-system;(cert-exporter|cluster-autoscaler|coredns|kube-state-metrics|net-exporter|nic-exporter|nginx-ingress-controller))|(giantswarm;chart-operator)|(giantswarm-elastic-logging;elastic-logging-elasticsearch-exporter)|(vault-exporter;vault-exporter)`)
)

// GetClusterID returns the value of the cluster annotation.
// Assumed that the service contains this annotation, see `FilterInvalidServices`.
func GetClusterID(service v1.Service) string {
	return service.ObjectMeta.Annotations[ClusterAnnotation]
}
