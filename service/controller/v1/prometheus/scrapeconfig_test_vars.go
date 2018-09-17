package prometheus

import (
	"net/url"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
)

var (
	TestConfigOneApiserver = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-apiserver",
		Scheme:  "https",
		HTTPClientConfig: config.HTTPClientConfig{
			TLSConfig: config.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: true,
			},
		},
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{&url.URL{
						Scheme: "https",
						Host:   "apiserver.xa5ly",
					}},
					Role: config.KubernetesRoleEndpoint,
					TLSConfig: config.TLSConfig{
						CAFile:             "/certs/xa5ly-ca.pem",
						CertFile:           "/certs/xa5ly-crt.pem",
						KeyFile:            "/certs/xa5ly-key.pem",
						InsecureSkipVerify: false,
					},
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			{
				SourceLabels: model.LabelNames{
					KubernetesSDNamespaceLabel,
					KubernetesSDServiceNameLabel,
				},
				Regex:  APIServerRegexp,
				Action: config.RelabelKeep,
			},
			{
				TargetLabel: AppLabel,
				Replacement: KubernetesAppName,
			},
			{
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
			},
		},
	}
	TestConfigOneCadvisor = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-cadvisor",
		Scheme:  "https",
		HTTPClientConfig: config.HTTPClientConfig{
			TLSConfig: config.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{&url.URL{
						Scheme: "https",
						Host:   "apiserver.xa5ly",
					}},
					Role: config.KubernetesRoleNode,
					TLSConfig: config.TLSConfig{
						CAFile:             "/certs/xa5ly-ca.pem",
						CertFile:           "/certs/xa5ly-crt.pem",
						KeyFile:            "/certs/xa5ly-key.pem",
						InsecureSkipVerify: false,
					},
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			{
				TargetLabel: model.AddressLabel,
				Replacement: "apiserver.xa5ly",
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDNodeNameLabel},
				Replacement:  CadvisorMetricsPath,
				TargetLabel:  model.MetricsPathLabel,
			},
			{
				TargetLabel: AppLabel,
				Replacement: CadvisorAppName,
			},
			{
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
			},
			{
				TargetLabel:  IPLabel,
				SourceLabels: model.LabelNames{KubernetesSDNodeAddressInternalIPLabel},
			},
			{
				TargetLabel:  RoleLabel,
				SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
				Regex:        EmptyRegexp,
				Replacement:  WorkerRole,
				TargetLabel:  RoleLabel,
			},
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
	}
	TestConfigOneKubelet = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-kubelet",
		Scheme:  "https",
		HTTPClientConfig: config.HTTPClientConfig{
			TLSConfig: config.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: true,
			},
		},
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{&url.URL{
						Scheme: "https",
						Host:   "apiserver.xa5ly",
					}},
					Role: config.KubernetesRoleNode,
					TLSConfig: config.TLSConfig{
						CAFile:             "/certs/xa5ly-ca.pem",
						CertFile:           "/certs/xa5ly-crt.pem",
						KeyFile:            "/certs/xa5ly-key.pem",
						InsecureSkipVerify: false,
					},
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			{
				TargetLabel: AppLabel,
				Replacement: KubeletAppName,
			},
			{
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
			},
			{
				TargetLabel:  IPLabel,
				SourceLabels: model.LabelNames{KubernetesSDNodeAddressInternalIPLabel},
			},
			{
				TargetLabel:  RoleLabel,
				SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDNodeLabelRole},
				Regex:        EmptyRegexp,
				Replacement:  WorkerRole,
				TargetLabel:  RoleLabel,
			},
		},
	}
	TestConfigOneNodeExporter = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-node-exporter",
		Scheme:  "http",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{&url.URL{
						Scheme: "https",
						Host:   "apiserver.xa5ly",
					}},
					Role: config.KubernetesRoleEndpoint,
					TLSConfig: config.TLSConfig{
						CAFile:             "/certs/xa5ly-ca.pem",
						CertFile:           "/certs/xa5ly-crt.pem",
						KeyFile:            "/certs/xa5ly-key.pem",
						InsecureSkipVerify: false,
					},
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			{
				SourceLabels: model.LabelNames{
					KubernetesSDNamespaceLabel,
					KubernetesSDServiceNameLabel,
				},
				Regex:  NodeExporterRegexp,
				Action: config.RelabelKeep,
			},
			{
				SourceLabels: model.LabelNames{model.AddressLabel},
				Regex:        KubeletPortRegexp,
				Replacement:  NodeExporterPort,
				TargetLabel:  model.AddressLabel,
			},
			{
				TargetLabel: AppLabel,
				Replacement: NodeExporterAppName,
			},
			{
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
			},
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
	}
	TestConfigOneWorkload = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-workload",
		HTTPClientConfig: config.HTTPClientConfig{
			TLSConfig: config.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		Scheme: "https",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{&url.URL{
						Scheme: "https",
						Host:   "apiserver.xa5ly",
					}},
					Role: config.KubernetesRoleEndpoint,
					TLSConfig: config.TLSConfig{
						CAFile:             "/certs/xa5ly-ca.pem",
						CertFile:           "/certs/xa5ly-crt.pem",
						KeyFile:            "/certs/xa5ly-key.pem",
						InsecureSkipVerify: false,
					},
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			{
				SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
				Regex:        WhitelistRegexp,
				Action:       config.RelabelKeep,
			},
			{
				TargetLabel:  AppLabel,
				SourceLabels: model.LabelNames{KubernetesSDServiceNameLabel},
			},
			{
				TargetLabel:  NamespaceLabel,
				SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel},
			},
			{
				TargetLabel:  PodNameLabel,
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
			},
			{
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
			},
			{
				TargetLabel: AddressLabel,
				Replacement: key.APIServiceHost(key.PrefixMaster, "xa5ly"),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        KubeStateMetricsPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.KubeStateMetricsNamespace, key.KubeStateMetricsPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        NginxIngressControllerPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NginxIngressControllerNamespace, key.NginxIngressControllerMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        ChartOperatorPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ChartOperatorNamespace, key.ChartOperatorMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        CertExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CertExporterNamespace, key.CertExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        NetExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NetExporterNamespace, key.NetExporterMetricPort),
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
			{
				Action:       ActionRelabel,
				SourceLabels: model.LabelNames{MetricExportedNamespaceLabel, MetricNamespaceLabel},
				Regex:        KubeSystemRelabelNamespaceRegexp,
				Replacement:  NamespaceKubeSystemLabel,
				TargetLabel:  ExportedNamespaceLabel,
			},
			// keep only kube-system cadvisor metrics
			{
				Action:       ActionKeep,
				SourceLabels: model.LabelNames{MetricExportedNamespaceLabel},
				Regex:        KubeSystemGiantswarmNSRegexp,
			},
			// drop useless IC metrics
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricDropICRegexp,
			},
			// drop several bucket latency metric
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricDropBucketLatencies,
			},
		},
	}
)

var (
	TestConfigTwoApiserver    config.ScrapeConfig = TestConfigOneApiserver
	TestConfigTwoCadvisor     config.ScrapeConfig = TestConfigOneCadvisor
	TestConfigTwoKubelet      config.ScrapeConfig = TestConfigOneKubelet
	TestConfigTwoNodeExporter config.ScrapeConfig = TestConfigOneNodeExporter
	TestConfigTwoWorkload     config.ScrapeConfig = TestConfigOneWorkload
)

func init() {
	apiserver := "apiserver.0ba9v"
	clusterID := "0ba9v"
	caFile := "/certs/0ba9v-ca.pem"
	crtFile := "/certs/0ba9v-crt.pem"
	keyFile := "/certs/0ba9v-key.pem"

	TestConfigTwoApiserver.JobName = "guest-cluster-0ba9v-apiserver"
	TestConfigTwoApiserver.HTTPClientConfig.TLSConfig.CAFile = caFile
	TestConfigTwoApiserver.HTTPClientConfig.TLSConfig.CertFile = crtFile
	TestConfigTwoApiserver.HTTPClientConfig.TLSConfig.KeyFile = keyFile
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoApiserver.RelabelConfigs[2].Replacement = clusterID

	TestConfigTwoCadvisor.JobName = "guest-cluster-0ba9v-cadvisor"
	TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig.CAFile = caFile
	TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig.CertFile = crtFile
	TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig.KeyFile = keyFile
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoCadvisor.RelabelConfigs[0].Replacement = apiserver
	TestConfigTwoCadvisor.RelabelConfigs[3].Replacement = clusterID

	TestConfigTwoKubelet.JobName = "guest-cluster-0ba9v-kubelet"
	TestConfigTwoKubelet.HTTPClientConfig.TLSConfig.CAFile = caFile
	TestConfigTwoKubelet.HTTPClientConfig.TLSConfig.CertFile = crtFile
	TestConfigTwoKubelet.HTTPClientConfig.TLSConfig.KeyFile = keyFile
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoKubelet.RelabelConfigs[1].Replacement = clusterID

	TestConfigTwoNodeExporter.JobName = "guest-cluster-0ba9v-node-exporter"
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoNodeExporter.RelabelConfigs[3].Replacement = clusterID

	TestConfigTwoWorkload.JobName = "guest-cluster-0ba9v-workload"
	TestConfigTwoWorkload.HTTPClientConfig.TLSConfig.CAFile = caFile
	TestConfigTwoWorkload.HTTPClientConfig.TLSConfig.CertFile = crtFile
	TestConfigTwoWorkload.HTTPClientConfig.TLSConfig.KeyFile = keyFile
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoWorkload.RelabelConfigs[4].Replacement = clusterID
	TestConfigTwoWorkload.RelabelConfigs[6].Replacement = key.APIServiceHost(key.PrefixMaster, clusterID)
}
