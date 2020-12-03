package prometheus

import (
	"fmt"
	"net/url"
	"sort"
	"time"

	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/relabel"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
)

const (
	// jobNamePrefix is the prefix for all guest cluster jobs.
	jobNamePrefix = "guest-cluster"

	// HttpScheme is the scheme for http connections.
	HttpScheme = "http"
	// HttpsScheme is the scheme for https connections.
	HttpsScheme = "https"

	// APIServerJobType is the job type for scraping Kubernetes API servers.
	APIServerJobType = "apiserver"
	// AWSNodeJobType is the job type for scraping aws-node Pods.
	AWSNodeJobType = "aws-node"
	// CadvisorJobType is the job type for scraping Cadvisor.
	CadvisorJobType = "cadvisor"
	// CalicoNodeJobType is the job type for scraping calico-node pods.
	CalicoNodeJobType = "calico-node"
	// DockerDaemonJobType is the job type for scraping docker daemon.
	DockerDaemonJobType = "docker-daemon"
	// EtcdJobType is the job type for scraping etcd.
	EtcdJobType = "etcd"
	// KubeletJobType is the job type for scraping kubelets.
	KubeletJobType = "kubelet"
	// ManagedAppJobType is the job type for scraping managed app metrics.
	ManagedAppJobType = "managed-app"
	// NodeExporterJobType is the job type for scraping node-exporters
	NodeExporterJobType = "node-exporter"
	// WorkloadJobType is the job type for scraping general workloads.
	WorkloadJobType = "workload"
	// IngressJobType is the job type for scraping the ingress controller
	IngressJobType = "ingress"
	// KubeStateManagedAppJobType is the job type for scraping kube-state-metrics-provided endpoints for managed apps.
	KubeStateManagedAppJobType = "kube-state-managed-app"
	// KubeProxyJobType is the job type for scraping node-exporters
	KubeProxyJobType = "kube-proxy"

	// ActionKeep is action type that keeps only matching metrics.
	ActionKeep = "keep"
	// ActionDrop is action type that drops matching metrics.
	ActionDrop = "drop"
	// ActionRelabel is action type that relabel metrics.
	ActionRelabel = "replace"
)

// getJobName takes a cluster ID, and returns a suitable job name.
func getJobName(service v1.Service, name string) string {
	return fmt.Sprintf("%s-%s-%s", jobNamePrefix, service.Namespace, name)
}

// getTargetHost takes a Kubernetes Service, and returns a suitable host.
func getTargetHost(service v1.Service) string {
	return fmt.Sprintf("%s.%s", service.Name, service.Namespace)
}

// getTarget takes a Kubernetes Service, and returns a LabelSet,
// suitable for use as a target.
func getTarget(service v1.Service) model.LabelSet {
	return model.LabelSet{
		model.AddressLabel: model.LabelValue(getTargetHost(service)),
	}
}

// getEtcdTarget takes a etcd url, and returns a LabelSet,
// suitable for use as a target.
func getEtcdTarget(etcdUrl string) model.LabelSet {
	return model.LabelSet{
		model.AddressLabel: model.LabelValue(etcdUrl),
	}
}

// getScrapeConfigs takes a Service, and returns a list of ScrapeConfigs.
// It is assumed that filtering has already taken place, and the cluster annotation exists.
func getScrapeConfigs(service v1.Service, metaConfig Config) []config.ScrapeConfig {
	certificateDirectory := metaConfig.CertDirectory
	clusterID := GetClusterID(service)
	provider := metaConfig.Provider

	secureTLSConfig := config_util.TLSConfig{
		CAFile:             key.CAPath(certificateDirectory, clusterID),
		CertFile:           key.CrtPath(certificateDirectory, clusterID),
		KeyFile:            key.KeyPath(certificateDirectory, clusterID),
		InsecureSkipVerify: false,
	}
	secureHTTPClientConfig := config_util.HTTPClientConfig{
		TLSConfig: secureTLSConfig,
	}
	insecureTLSConfig := config_util.TLSConfig{
		CAFile:             key.CAPath(certificateDirectory, clusterID),
		CertFile:           key.CrtPath(certificateDirectory, clusterID),
		KeyFile:            key.KeyPath(certificateDirectory, clusterID),
		InsecureSkipVerify: true,
	}
	insecureHTTPClientConfig := config_util.HTTPClientConfig{
		TLSConfig: insecureTLSConfig,
	}

	endpointSDConfig := sd_config.ServiceDiscoveryConfig{
		KubernetesSDConfigs: []*kubernetes.SDConfig{
			{
				APIServer: config_util.URL{
					URL: &url.URL{
						Scheme: HttpsScheme,
						Host:   getTargetHost(service),
					},
				},
				Role: kubernetes.RoleEndpoint,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: secureTLSConfig,
				},
			},
		},
	}
	nodeSDConfig := sd_config.ServiceDiscoveryConfig{
		KubernetesSDConfigs: []*kubernetes.SDConfig{
			{
				APIServer: config_util.URL{
					URL: &url.URL{
						Scheme: HttpsScheme,
						Host:   getTargetHost(service),
					},
				},
				Role: kubernetes.RoleNode,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: secureTLSConfig,
				},
			},
		},
	}
	podSDConfig := sd_config.ServiceDiscoveryConfig{
		KubernetesSDConfigs: []*kubernetes.SDConfig{
			{
				APIServer: config_util.URL{
					URL: &url.URL{
						Scheme: HttpsScheme,
						Host:   getTargetHost(service),
					},
				},
				Role: kubernetes.RolePod,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: secureTLSConfig,
				},
			},
		},
	}

	clusterIDLabelRelabelConfig := &relabel.Config{
		TargetLabel: ClusterIDLabel,
		Replacement: clusterID,
	}
	clusterTypeLabelRelabelConfig := &relabel.Config{
		TargetLabel: ClusterTypeLabel,
		Replacement: GuestClusterType,
	}
	providerLabelRelabelConfig := &relabel.Config{
		TargetLabel: ProviderLabel,
		Replacement: provider,
	}
	reflectorRelabelConfig := &relabel.Config{
		Action:       ActionDrop,
		SourceLabels: model.LabelNames{MetricNameLabel},
		Regex:        MetricsDropReflectorRegexp,
	}
	rewriteAddress := &relabel.Config{
		TargetLabel: AddressLabel,
		Replacement: key.APIServiceHost(key.PrefixMaster, clusterID),
	}
	rewriteManagedAppMetricPath := &relabel.Config{
		SourceLabels: model.LabelNames{model.LabelName(NamespaceLabel), model.LabelName(PodNameLabel), KubernetesSDServiceGiantSwarmMonitoringPortLabel, KubernetesSDServiceGiantSwarmMonitoringPathLabel},
		Regex:        ManagedAppSourceRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.ManagedAppPodMetricsPath(),
	}
	rewriteKubeStateMetricPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        KubeStateMetricsPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.KubeStateMetricsNamespace, key.KubeStateMetricsPort),
	}
	rewriteICMetricPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        NginxIngressControllerPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.NginxIngressControllerNamespace, key.NginxIngressControllerMetricPort),
	}
	rewriteAWSNodePath := &relabel.Config{
		SourceLabels: model.LabelNames{PodSDPodNameLabel},
		Regex:        AWSNodePodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.AWSNodeNamespace, key.AWSNodeMetricPort),
	}
	rewriteCalicoNodePath := &relabel.Config{
		SourceLabels: model.LabelNames{PodSDPodNameLabel},
		Regex:        CalicoNodePodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.CalicoNodeNamespace, key.CalicoNodeMetricPort),
	}
	rewriteChartOperatorPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        ChartOperatorPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.ChartOperatorNamespace, key.ChartOperatorMetricPort),
	}
	rewriteCertExporterPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        CertExporterPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.CertExporterNamespace, key.CertExporterMetricPort),
	}
	rewriteClusterAutoscalerPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        ClusterAutoscalerPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.ClusterAutoscalerNamespace, key.ClusterAutoscalerMetricPort),
	}
	rewriteCoreDNSPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        CoreDNSPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.CoreDNSNamespace, key.CoreDNSMetricPort),
	}
	rewriteElasticLoggingMetricPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        ElasticLoggingPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.ElasticLoggingNamespace, key.ElasticLoggingMetricPort),
	}
	rewriteNetExporterPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        NetExporterPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.NetExporterNamespace, key.NetExporterMetricPort),
	}
	rewriteNicExporterPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        NicExporterPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.NicExporterNamespace, key.NicExporterMetricPort),
	}
	rewriteKiamPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        KiamPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.KiamNamespace, key.KiamMetricPort),
	}
	rewriteKubeProxyPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        KubeProxyPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.KubeProxyNamespace, key.KubeProxyMetricPort),
	}
	rewriteVaultExporterPath := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        VaultExporterPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.VaultExporterNamespace, key.VaultExporterMetricPort),
	}

	ipLabelRelabelConfig := &relabel.Config{
		TargetLabel:  IPLabel,
		SourceLabels: model.LabelNames{KubernetesSDNodeAddressInternalIPLabel},
	}
	roleLabelRelabelConfig := &relabel.Config{
		TargetLabel:  RoleLabel,
		SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
	}
	missingRoleLabelRelabelConfig := &relabel.Config{
		SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
		Regex:        EmptyRegexp,
		Replacement:  WorkerRole,
		TargetLabel:  RoleLabel,
	}

	scrapeConfigs := []config.ScrapeConfig{
		{
			JobName:                getJobName(service, APIServerJobType),
			Scheme:                 HttpsScheme,
			HTTPClientConfig:       insecureHTTPClientConfig,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep api server endpoints.
				{
					SourceLabels: model.LabelNames{
						KubernetesSDNamespaceLabel,
						KubernetesSDServiceNameLabel,
					},
					Regex:  APIServerRegexp,
					Action: relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: KubernetesAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// drop several bucket latency metric
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        MetricDropBucketLatencies,
				},
				reflectorRelabelConfig,
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, CadvisorJobType),
			Scheme:                 HttpsScheme,
			HTTPClientConfig:       secureHTTPClientConfig,
			ServiceDiscoveryConfig: nodeSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Relabel address to kubernetes service.
				{
					TargetLabel: model.AddressLabel,
					Replacement: getTargetHost(service),
				},
				// Relabel metrics path to cadvisor proxy.
				{
					SourceLabels: model.LabelNames{KubernetesSDNodeNameLabel},
					Replacement:  CadvisorMetricsPath,
					TargetLabel:  model.MetricsPathLabel,
				},
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: CadvisorAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// Add ip label.
				ipLabelRelabelConfig,
				// Add role label.
				roleLabelRelabelConfig,
				missingRoleLabelRelabelConfig,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// keep only kube-system cadvisor metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricNamespaceLabel},
					Regex:        NSRegexp,
				},
				// drop cadvisor metrics about container network statistics
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        MetricDropContainerNetworkRegexp,
				},
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, KubeletJobType),
			Scheme:                 HttpsScheme,
			HTTPClientConfig:       insecureHTTPClientConfig,
			ServiceDiscoveryConfig: nodeSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: KubeletAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// Add ip label.
				ipLabelRelabelConfig,
				// Add role label.
				roleLabelRelabelConfig,
				missingRoleLabelRelabelConfig,
			},
			MetricRelabelConfigs: []*relabel.Config{
				reflectorRelabelConfig,
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, AWSNodeJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: podSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep kube-state-metrics targets.
				{
					SourceLabels: model.LabelNames{PodSDNamespaceLabel, PodSDPodNameLabel},
					Regex:        AWSNodePodRegexp,
					Action:       relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel:  AppLabel,
					SourceLabels: model.LabelNames{PodSDContainerNameLabel},
				},
				// Add namespace label.
				{
					TargetLabel:  NamespaceLabel,
					SourceLabels: model.LabelNames{PodSDNamespaceLabel},
				},
				// Add pod_name label.
				{
					TargetLabel:  PodNameLabel,
					SourceLabels: model.LabelNames{PodSDPodNameLabel},
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteAWSNodePath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, CalicoNodeJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: podSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep calico node targets.
				{
					SourceLabels: model.LabelNames{PodSDNamespaceLabel, PodSDPodNameLabel},
					Regex:        CalicoNodePodRegexp,
					Action:       relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel:  AppLabel,
					SourceLabels: model.LabelNames{PodSDContainerNameLabel},
				},
				// Add namespace label.
				{
					TargetLabel:  NamespaceLabel,
					SourceLabels: model.LabelNames{PodSDNamespaceLabel},
				},
				// Add pod_name label.
				{
					TargetLabel:  PodNameLabel,
					SourceLabels: model.LabelNames{PodSDPodNameLabel},
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteCalicoNodePath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, DockerDaemonJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: nodeSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Relabel address to kubernetes service.
				{
					TargetLabel: model.AddressLabel,
					Replacement: getTargetHost(service),
				},
				// Relabel metrics path to cadvisor proxy.
				{
					SourceLabels: model.LabelNames{KubernetesSDNodeNameLabel},
					Replacement:  DockerMetricsPath,
					TargetLabel:  model.MetricsPathLabel,
				},
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: DockerAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// Add ip label.
				ipLabelRelabelConfig,
				// Add role label.
				roleLabelRelabelConfig,
				missingRoleLabelRelabelConfig,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// Keep only metrics with names listed in DockerMetricsNameRegexp.
				{
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        DockerMetricsNameRegexp,
					Action:       ActionKeep,
				},
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, KubeStateManagedAppJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep kube-state-metrics targets.
				{
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
					Regex:        KubeStateMetricsServiceNameRegexp,
					Action:       relabel.Keep,
				},
				// Add kube_state_metrics_for_managed_apps label.
				{
					TargetLabel: KubeStateMetricsForManagedApps,
					Replacement: "true",
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteKubeStateMetricPath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// keep only metrics with names listed in KubeStateMetricsManagedAppMetricsNameRegexp
				{
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        KubeStateMetricsManagedAppMetricsNameRegexp,
					Action:       ActionKeep,
				},
				// copy exported_namespace as namespace
				{
					SourceLabels: model.LabelNames{MetricExportedNamespaceLabel},
					TargetLabel:  NamespaceLabel,
				},
				// apply correct workload type label
				{
					SourceLabels: model.LabelNames{DeploymentTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadTypeLabel,
					Replacement:  ManagedAppsDeployment,
				},
				{
					SourceLabels: model.LabelNames{DaemonSetTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadTypeLabel,
					Replacement:  ManagedAppsDaemonSet,
				},
				{
					SourceLabels: model.LabelNames{StatefulSetTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadTypeLabel,
					Replacement:  ManagedAppsStatefulSet,
				},
				// copy type-specific workload name label into generic "workload_name"
				{
					SourceLabels: model.LabelNames{DeploymentTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadNameLabel,
					Replacement:  GroupCapture,
				},
				{
					SourceLabels: model.LabelNames{DaemonSetTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadNameLabel,
					Replacement:  GroupCapture,
				},
				{
					SourceLabels: model.LabelNames{StatefulSetTypeLabel},
					Regex:        NonEmptyRegexp,
					TargetLabel:  ManagedAppWorkloadNameLabel,
					Replacement:  GroupCapture,
				},
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, NodeExporterJobType),
			Scheme:                 HttpScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep node-exporter endpoints.
				{
					SourceLabels: model.LabelNames{
						KubernetesSDNamespaceLabel,
						KubernetesSDServiceNameLabel,
					},
					Regex:  NodeExporterRegexp,
					Action: relabel.Keep,
				},
				// Relabel address to node-exporter port.
				{
					SourceLabels: model.LabelNames{model.AddressLabel},
					Regex:        KubeletPortRegexp,
					Replacement:  NodeExporterPort,
					TargetLabel:  model.AddressLabel,
				},
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: NodeExporterAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// Add ip label.
				{
					SourceLabels: model.LabelNames{model.AddressLabel},
					Regex:        NodeExporterPortRegexp,
					Replacement:  GroupCapture,
					TargetLabel:  IPLabel,
				},
			},
			MetricRelabelConfigs: []*relabel.Config{
				// Drop many mounts that are not interesting based on fstype.
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricFSTypeLabel},
					Regex:        MetricDropFStypeRegexp,
				},
				// We care only about systemd state failed, we can drop the rest.
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricNameLabel, MetricSystemdStateLabel},
					Regex:        MetricDropSystemdStateRegexp,
				},
				// Drop all systemd units that are connected to docker mounts or networking, we don't care about them.
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricNameLabel, MetricSystemdNameLabel},
					Regex:        MetricDropSystemdNameRegexp,
				},
				providerLabelRelabelConfig,
			},
		},
		{
			JobName:                getJobName(service, WorkloadJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep kube-state-metrics targets.
				{
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
					Regex:        ServiceWhitelistRegexp,
					Action:       relabel.Keep,
				},
				// Drop non-managed kiam pods and keep only Giantswarm managed kiam pods
				{
					SourceLabels: model.LabelNames{KubernetesSDPodNameLabel, PodSDGiantswarmServiceTypeLabel},
					Regex:        KiamPodNameRegexpNonManaged,
					Action:       relabel.Drop,
				},
				// Add app label.
				{
					TargetLabel:  AppLabel,
					SourceLabels: model.LabelNames{KubernetesSDServiceNameLabel},
				},
				// Add namespace label.
				{
					TargetLabel:  NamespaceLabel,
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel},
				},
				// Add pod_name label.
				{
					TargetLabel:  PodNameLabel,
					SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				},
				// Add node label.
				{
					TargetLabel:  NodeLabel,
					SourceLabels: model.LabelNames{KubernetesSDPodNodeNameLabel},
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteKubeStateMetricPath,
				rewriteCalicoNodePath,
				rewriteChartOperatorPath,
				rewriteCertExporterPath,
				rewriteClusterAutoscalerPath,
				rewriteCoreDNSPath,
				rewriteElasticLoggingMetricPath,
				rewriteNetExporterPath,
				rewriteNicExporterPath,
				rewriteKiamPath,
				rewriteVaultExporterPath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// relabel namespace to exported_namespace for endpoints in kube-system namespace.
				// this keeps metrics from nginx ingress controller from being dropped by filter below
				{
					Action:       ActionRelabel,
					SourceLabels: model.LabelNames{MetricExportedNamespaceLabel, MetricNamespaceLabel},
					Regex:        RelabelNamespaceRegexp,
					Replacement:  GroupCapture,
					TargetLabel:  ExportedNamespaceLabel,
				},
				// keep only kube-system cadvisor metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricExportedNamespaceLabel},
					Regex:        NSRegexp,
				},
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, IngressJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep ingress controller targets.
				{
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
					Regex:        IngressWhitelistRegexp,
					Action:       relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel:  AppLabel,
					SourceLabels: model.LabelNames{KubernetesSDServiceNameLabel},
				},
				// Add namespace label.
				{
					TargetLabel:  NamespaceLabel,
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel},
				},
				// Add pod_name label.
				{
					TargetLabel:  PodNameLabel,
					SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				},
				// Add node label.
				{
					TargetLabel:  NodeLabel,
					SourceLabels: model.LabelNames{KubernetesSDPodNodeNameLabel},
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteICMetricPath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// relabel namespace to exported_namespace for endpoints in kube-system namespace.
				// this keeps metrics from nginx ingress controller from being dropped by filter below
				{
					Action:       ActionRelabel,
					SourceLabels: model.LabelNames{MetricExportedNamespaceLabel, MetricNamespaceLabel},
					Regex:        RelabelNamespaceRegexp,
					Replacement:  GroupCapture,
					TargetLabel:  ExportedNamespaceLabel,
				},
				// keep useful IC metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        MetricKeepICRegexp,
				},
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, ManagedAppJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep monitoring label presents
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPresentLabel},
					Regex:        relabel.MustNewRegexp(`(true)`),
					Action:       relabel.Keep,
				},
				// Only keep monitoring label as true
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringLabel},
					Regex:        relabel.MustNewRegexp(`(true)`),
					Action:       relabel.Keep,
				},
				// Only keep when monitoring port presents in annotation.
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPortPresentLabel},
					Regex:        relabel.MustNewRegexp(`(true)`),
					Action:       relabel.Keep,
				},
				// Only keep when monitoring path presents in annotation.
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPathPresentLabel},
					Regex:        relabel.MustNewRegexp(`(true)`),
					Action:       relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel:  AppLabel,
					SourceLabels: model.LabelNames{KubernetesSDServiceNameLabel},
				},
				// Add namespace label.
				{
					TargetLabel:  NamespaceLabel,
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel},
				},
				// Add pod_name label.
				{
					TargetLabel:  PodNameLabel,
					SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				},
				// Add application type label.
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringAppTypeLabel},
					Regex:        relabel.MustNewRegexp(`(optional|default)`),
					TargetLabel:  AppTypeLabel,
				},
				// Add is_managed_app label.
				{
					SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPresentLabel},
					Regex:        relabel.MustNewRegexp(`(true)`),
					TargetLabel:  AppIsManaged,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// Relabel metrics path to specific managed app proxy.
				rewriteManagedAppMetricPath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				providerLabelRelabelConfig,
			},
		},

		{
			JobName:                getJobName(service, KubeProxyJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: podSDConfig,
			RelabelConfigs: []*relabel.Config{
				// Only keep node-exporter endpoints.
				{
					SourceLabels: model.LabelNames{
						KubernetesSDPodNameLabel,
					},
					Regex:  KubeProxyPodNameRegexp,
					Action: relabel.Keep,
				},
				// Add app label.
				{
					TargetLabel: AppLabel,
					Replacement: KubeProxyAppName,
				},
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				rewriteKubeProxyPath,
			},
			MetricRelabelConfigs: []*relabel.Config{
				// keep only kube-proxy iptables restore errors metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        MetricsKeepKubeProxyIptableRegexp,
				},
				providerLabelRelabelConfig,
			},
		},
	}
	// check if we can add etcd monitoring

	//  to ensure all components in cloud are ready we delay creation of etcd scrape config by 30 minutes
	etcdScrapeDelay := metav1.Time{Time: time.Now().Add(-time.Minute * 30)}

	if service.CreationTimestamp.Before(&etcdScrapeDelay) {
		if _, ok := service.Annotations[key.AnnotationEtcdDomain]; ok {
			// prepare etcd static discovery config
			etcdStaticConfig := sd_config.ServiceDiscoveryConfig{
				StaticConfigs: []*targetgroup.Group{
					{
						Targets: []model.LabelSet{
							getEtcdTarget(service.Annotations[key.AnnotationEtcdDomain]),
						},
						Labels: model.LabelSet{
							model.LabelName(ClusterTypeLabel): model.LabelValue(GuestClusterType),
							model.LabelName(ClusterIDLabel):   model.LabelValue(clusterID),
							model.LabelName(ProviderLabel):    model.LabelValue(provider),
						},
					},
				},
			}

			etcdScrapeConfig := config.ScrapeConfig{
				JobName:                getJobName(service, EtcdJobType),
				Scheme:                 HttpsScheme,
				HTTPClientConfig:       secureHTTPClientConfig,
				ServiceDiscoveryConfig: etcdStaticConfig,
				MetricRelabelConfigs: []*relabel.Config{
					providerLabelRelabelConfig,
				},
			}
			// append etcd scrape config
			scrapeConfigs = append(scrapeConfigs, etcdScrapeConfig)
		}
	}

	return scrapeConfigs
}

// GetScrapeConfigs takes a list of Kubernetes Services,
// and returns a list of Prometheus ScrapeConfigs.
func GetScrapeConfigs(services []v1.Service, metaConfig Config) ([]config.ScrapeConfig, error) {
	filteredServices := FilterInvalidServices(services)

	scrapeConfigs := []config.ScrapeConfig{}
	for _, service := range filteredServices {
		scrapeConfigs = append(scrapeConfigs, getScrapeConfigs(service, metaConfig)...)
	}

	sort.Slice(scrapeConfigs, func(i, j int) bool {
		return scrapeConfigs[i].JobName < scrapeConfigs[j].JobName
	})

	return scrapeConfigs, nil
}
