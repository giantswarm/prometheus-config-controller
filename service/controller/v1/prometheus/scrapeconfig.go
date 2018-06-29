package prometheus

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"k8s.io/api/core/v1"

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
	// CadvisorJobType is the job type for scraping Cadvisor.
	CadvisorJobType = "cadvisor"
	// KubeletJobType is the job type for scraping kubelets.
	KubeletJobType = "kubelet"
	// NodeExporterJobType is the job type for scraping node-exporters
	NodeExporterJobType = "node-exporter"
	// WorkloadJobType is the job type for scraping general workloads.
	WorkloadJobType = "workload"

	// ActionKeep is action type that keeps only matching metrics.
	ActionKeep = "keep"
	// ActionDrop is action type that drops matching metrics.
	ActionDrop = "drop"
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

// getScrapeConfigs takes a Service, and returns a list of ScrapeConfigs.
// It is assumed that filtering has already taken place, and the cluster annotation exists.
func getScrapeConfigs(service v1.Service, certificateDirectory string) []config.ScrapeConfig {
	clusterID := GetClusterID(service)

	secureTLSConfig := config.TLSConfig{
		CAFile:             key.CAPath(certificateDirectory, clusterID),
		CertFile:           key.CrtPath(certificateDirectory, clusterID),
		KeyFile:            key.KeyPath(certificateDirectory, clusterID),
		InsecureSkipVerify: false,
	}
	secureHTTPClientConfig := config.HTTPClientConfig{
		TLSConfig: secureTLSConfig,
	}
	insecureTLSConfig := config.TLSConfig{
		CAFile:             key.CAPath(certificateDirectory, clusterID),
		CertFile:           key.CrtPath(certificateDirectory, clusterID),
		KeyFile:            key.KeyPath(certificateDirectory, clusterID),
		InsecureSkipVerify: true,
	}
	insecureHTTPClientConfig := config.HTTPClientConfig{
		TLSConfig: insecureTLSConfig,
	}

	endpointSDConfig := config.ServiceDiscoveryConfig{
		KubernetesSDConfigs: []*config.KubernetesSDConfig{
			{
				APIServer: config.URL{&url.URL{
					Scheme: HttpsScheme,
					Host:   getTargetHost(service),
				}},
				Role:      config.KubernetesRoleEndpoint,
				TLSConfig: secureTLSConfig,
			},
		},
	}
	nodeSDConfig := config.ServiceDiscoveryConfig{
		KubernetesSDConfigs: []*config.KubernetesSDConfig{
			{
				APIServer: config.URL{&url.URL{
					Scheme: HttpsScheme,
					Host:   getTargetHost(service),
				}},
				Role:      config.KubernetesRoleNode,
				TLSConfig: secureTLSConfig,
			},
		},
	}

	clusterIDLabelRelabelConfig := &config.RelabelConfig{
		TargetLabel: ClusterIDLabel,
		Replacement: clusterID,
	}
	clusterTypeLabelRelabelConfig := &config.RelabelConfig{
		TargetLabel: ClusterTypeLabel,
		Replacement: GuestClusterType,
	}
	rewriteAddress := &config.RelabelConfig{
		TargetLabel: AddressLabel,
		Replacement: key.APIServiceHost(key.PrefixMaster, clusterID),
	}
	rewriteKubeStateMetricPath := &config.RelabelConfig{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        KubeStateMetricsPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.KubeStaeMetricsPort),
	}
	rewriteICMetricPath := &config.RelabelConfig{
		SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
		Regex:        NginxICPodNameRegexp,
		TargetLabel:  MetricPathLabel,
		Replacement:  key.APIProxyPodMetricsPath(key.NginxICMetricPort),
	}

	ipLabelRelabelConfig := &config.RelabelConfig{
		TargetLabel:  IPLabel,
		SourceLabels: model.LabelNames{KubernetesSDNodeAddressInternalIPLabel},
	}
	roleLabelRelabelConfig := &config.RelabelConfig{
		TargetLabel:  RoleLabel,
		SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
	}
	missingRoleLabelRelabelConfig := &config.RelabelConfig{
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
			RelabelConfigs: []*config.RelabelConfig{
				// Only keep api server endpoints.
				{
					SourceLabels: model.LabelNames{
						KubernetesSDNamespaceLabel,
						KubernetesSDServiceNameLabel,
					},
					Regex:  APIServerRegexp,
					Action: config.RelabelKeep,
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
		},

		{
			JobName:                getJobName(service, CadvisorJobType),
			Scheme:                 HttpsScheme,
			HTTPClientConfig:       secureHTTPClientConfig,
			ServiceDiscoveryConfig: nodeSDConfig,
			RelabelConfigs: []*config.RelabelConfig{
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
			MetricRelabelConfigs: []*config.RelabelConfig{
				// keep only kube-system cadvisor metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricNamespaceLabel},
					Regex:        KubeSystemGiantswarmNSRegexp,
				},
				// drop cadvisor metrics about container network statistics
				{
					Action:       ActionDrop,
					SourceLabels: model.LabelNames{MetricNameLabel},
					Regex:        MetricDropContainerNetworkRegexp,
				},
			},
		},

		{
			JobName:                getJobName(service, KubeletJobType),
			Scheme:                 HttpsScheme,
			HTTPClientConfig:       insecureHTTPClientConfig,
			ServiceDiscoveryConfig: nodeSDConfig,
			RelabelConfigs: []*config.RelabelConfig{
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
		},

		{
			JobName:                getJobName(service, NodeExporterJobType),
			Scheme:                 HttpScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*config.RelabelConfig{
				// Only keep node-exporter endpoints.
				{
					SourceLabels: model.LabelNames{
						KubernetesSDNamespaceLabel,
						KubernetesSDServiceNameLabel,
					},
					Regex:  NodeExporterRegexp,
					Action: config.RelabelKeep,
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
			MetricRelabelConfigs: []*config.RelabelConfig{
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
			},
		},

		{
			JobName:                getJobName(service, WorkloadJobType),
			HTTPClientConfig:       secureHTTPClientConfig,
			Scheme:                 HttpsScheme,
			ServiceDiscoveryConfig: endpointSDConfig,
			RelabelConfigs: []*config.RelabelConfig{
				// Only keep kube-state-metrics targets.
				{
					SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
					Regex:        WhitelistRegexp,
					Action:       config.RelabelKeep,
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
				// Add cluster_id label.
				clusterIDLabelRelabelConfig,
				// Add cluster_type label.
				clusterTypeLabelRelabelConfig,
				// rewrite host to api proxy
				rewriteAddress,
				// rewrite metrics scrape path to connect pods
				rewriteKubeStateMetricPath,
				rewriteICMetricPath,
			},
			MetricRelabelConfigs: []*config.RelabelConfig{
				// keep only kube-system cadvisor metrics
				{
					Action:       ActionKeep,
					SourceLabels: model.LabelNames{MetricExportedNamespaceLabel},
					Regex:        KubeSystemGiantswarmNSRegexp,
				},
			},
		},
	}

	return scrapeConfigs
}

// GetScrapeConfigs takes a list of Kubernetes Services,
// and returns a list of Prometheus ScrapeConfigs.
func GetScrapeConfigs(services []v1.Service, certificateDirectory string) ([]config.ScrapeConfig, error) {
	filteredServices := FilterInvalidServices(services)

	scrapeConfigs := []config.ScrapeConfig{}
	for _, service := range filteredServices {
		scrapeConfigs = append(scrapeConfigs, getScrapeConfigs(service, certificateDirectory)...)
	}

	sort.Slice(scrapeConfigs, func(i, j int) bool {
		return scrapeConfigs[i].JobName < scrapeConfigs[j].JobName
	})

	return scrapeConfigs, nil
}
