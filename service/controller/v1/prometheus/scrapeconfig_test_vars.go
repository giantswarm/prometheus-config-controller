package prometheus

import (
	"net/url"
	"strings"

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
					APIServer: config.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
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
		MetricRelabelConfigs: []*config.RelabelConfig{
			// drop several bucket latency metric
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricDropBucketLatencies,
			},
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricsDropReflectorRegexp,
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
					APIServer: config.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
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
				Regex:        NSRegexp,
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
					APIServer: config.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
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
		MetricRelabelConfigs: []*config.RelabelConfig{
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricsDropReflectorRegexp,
			},
		},
	}
	TestConfigOneNodeExporter = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-node-exporter",
		Scheme:  "http",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
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
					APIServer: config.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
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
				Regex:        ClusterAutoscalerPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ClusterAutoscalerNamespace, key.ClusterAutoscalerMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        CoreDNSPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CoreDNSNamespace, key.CoreDNSMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        ElasticLoggingPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ElasticLoggingNamespace, key.ElasticLoggingMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        NetExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NetExporterNamespace, key.NetExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        NicExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NicExporterNamespace, key.NicExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        VaultExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.VaultExporterNamespace, key.VaultExporterMetricPort),
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
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
			// drop useless IC metrics
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricDropICRegexp,
			},
		},
	}
)

var (
	TestConfigTwoApiserver    config.ScrapeConfig
	TestConfigTwoCadvisor     config.ScrapeConfig
	TestConfigTwoKubelet      config.ScrapeConfig
	TestConfigTwoNodeExporter config.ScrapeConfig
	TestConfigTwoWorkload     config.ScrapeConfig
)

func init() {

	TestConfigTwoApiserver = TestConfigOneApiserver
	TestConfigTwoCadvisor = TestConfigOneCadvisor
	TestConfigTwoKubelet = TestConfigOneKubelet
	TestConfigTwoNodeExporter = TestConfigOneNodeExporter
	TestConfigTwoWorkload = TestConfigOneWorkload

	apiServer := config.URL{&url.URL{
		Scheme: "https",
		Host:   "apiserver.0ba9v",
	}}

	clusterID := "0ba9v"

	tlsConfig := config.TLSConfig{
		CAFile:             "/certs/0ba9v-ca.pem",
		CertFile:           "/certs/0ba9v-crt.pem",
		KeyFile:            "/certs/0ba9v-key.pem",
		InsecureSkipVerify: false,
	}

	{
		TestConfigTwoApiserver.JobName = "guest-cluster-0ba9v-apiserver"
		TestConfigTwoApiserver.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoApiserver.HTTPClientConfig.TLSConfig.InsecureSkipVerify = true
		TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRoleEndpoint,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoApiserver.RelabelConfigs = nil
		for _, r := range TestConfigOneApiserver.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoApiserver.RelabelConfigs = append(TestConfigTwoApiserver.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoCadvisor.JobName = "guest-cluster-0ba9v-cadvisor"
		TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRoleNode,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoCadvisor.RelabelConfigs = nil
		for _, r := range TestConfigOneCadvisor.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoCadvisor.RelabelConfigs = append(TestConfigTwoCadvisor.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoKubelet.JobName = "guest-cluster-0ba9v-kubelet"
		TestConfigTwoKubelet.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoKubelet.HTTPClientConfig.TLSConfig.InsecureSkipVerify = true
		TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRoleNode,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoKubelet.RelabelConfigs = nil
		for _, r := range TestConfigOneKubelet.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoKubelet.RelabelConfigs = append(TestConfigTwoKubelet.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoNodeExporter.JobName = "guest-cluster-0ba9v-node-exporter"
		TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRoleEndpoint,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoNodeExporter.RelabelConfigs = nil
		for _, r := range TestConfigOneNodeExporter.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoNodeExporter.RelabelConfigs = append(TestConfigTwoNodeExporter.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoWorkload.JobName = "guest-cluster-0ba9v-workload"
		TestConfigTwoWorkload.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRoleEndpoint,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoWorkload.RelabelConfigs = nil
		for _, r := range TestConfigOneWorkload.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoWorkload.RelabelConfigs = append(TestConfigTwoWorkload.RelabelConfigs, &newRelabelConfig)
		}
	}
}
