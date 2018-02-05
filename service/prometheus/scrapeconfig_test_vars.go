package prometheus

import (
	"net/url"

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
					KubernetesSDEndpointPortNameLabel,
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
	TestConfigOneWorkload = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-workload",
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
				SourceLabels: model.LabelNames{KubernetesSDServiceNameLabel},
				Regex:        NodeExporterRegexp,
				Action:       config.RelabelDrop,
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
				TargetLabel: ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: ClusterTypeLabel,
				Replacement: GuestClusterType,
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
	tlsConfig := config.TLSConfig{
		CAFile:             "/certs/0ba9v-ca.pem",
		CertFile:           "/certs/0ba9v-crt.pem",
		KeyFile:            "/certs/0ba9v-key.pem",
		InsecureSkipVerify: false,
	}

	TestConfigTwoApiserver.JobName = "guest-cluster-0ba9v-apiserver"
	TestConfigTwoApiserver.HTTPClientConfig.TLSConfig = tlsConfig
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig = tlsConfig
	TestConfigTwoApiserver.RelabelConfigs[2].Replacement = clusterID

	TestConfigTwoCadvisor.JobName = "guest-cluster-0ba9v-cadvisor"
	TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig = tlsConfig
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig = tlsConfig
	TestConfigTwoCadvisor.RelabelConfigs[0].Replacement = apiserver
	TestConfigTwoCadvisor.RelabelConfigs[3].Replacement = clusterID

	TestConfigTwoKubelet.JobName = "guest-cluster-0ba9v-kubelet"
	TestConfigTwoKubelet.HTTPClientConfig.TLSConfig = tlsConfig
	TestConfigTwoKubelet.HTTPClientConfig.TLSConfig.InsecureSkipVerify = true
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig = tlsConfig
	TestConfigTwoKubelet.RelabelConfigs[1].Replacement = clusterID

	TestConfigTwoNodeExporter.JobName = "guest-cluster-0ba9v-node-exporter"
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig = tlsConfig
	TestConfigTwoNodeExporter.RelabelConfigs[2].Replacement = clusterID

	TestConfigTwoWorkload.JobName = "guest-cluster-0ba9v-workload"
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig = tlsConfig
	TestConfigTwoWorkload.RelabelConfigs[3].Replacement = clusterID
}
