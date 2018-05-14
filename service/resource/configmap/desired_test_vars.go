package configmap

import (
	"net/url"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
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
					prometheus.KubernetesSDNamespaceLabel,
					prometheus.KubernetesSDServiceNameLabel,
				},
				Regex:  prometheus.APIServerRegexp,
				Action: config.RelabelKeep,
			},
			{
				TargetLabel: prometheus.AppLabel,
				Replacement: prometheus.KubernetesAppName,
			},
			{
				TargetLabel: prometheus.ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: prometheus.ClusterTypeLabel,
				Replacement: prometheus.GuestClusterType,
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
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeNameLabel},
				Replacement:  prometheus.CadvisorMetricsPath,
				TargetLabel:  model.MetricsPathLabel,
			},
			{
				TargetLabel: prometheus.AppLabel,
				Replacement: prometheus.CadvisorAppName,
			},
			{
				TargetLabel: prometheus.ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: prometheus.ClusterTypeLabel,
				Replacement: prometheus.GuestClusterType,
			},
			{
				TargetLabel:  prometheus.IPLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeAddressInternalIPLabel},
			},
			{
				TargetLabel:  prometheus.RoleLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeLabelRole},
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeLabelRole},
				Regex:        prometheus.EmptyRegexp,
				Replacement:  prometheus.WorkerRole,
				TargetLabel:  prometheus.RoleLabel,
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
			// keep only kube-system cadvisor metrics
			{
				Action:       prometheus.ActionKeep,
				SourceLabels: model.LabelNames{prometheus.MetricNamespaceLabel},
				Regex:        prometheus.KubeSystemGiantswarmNSRegexp,
			},
			// drop cadvisor metrics about container network statistics
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricDropContainerNetworkRegexp,
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
				TargetLabel: prometheus.AppLabel,
				Replacement: prometheus.KubeletAppName,
			},
			{
				TargetLabel: prometheus.ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: prometheus.ClusterTypeLabel,
				Replacement: prometheus.GuestClusterType,
			},
			{
				TargetLabel:  prometheus.IPLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeAddressInternalIPLabel},
			},
			{
				TargetLabel:  prometheus.RoleLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeLabelRole},
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNodeLabelRole},
				Regex:        prometheus.EmptyRegexp,
				Replacement:  prometheus.WorkerRole,
				TargetLabel:  prometheus.RoleLabel,
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
					prometheus.KubernetesSDNamespaceLabel,
					prometheus.KubernetesSDServiceNameLabel,
				},
				Regex:  prometheus.NodeExporterRegexp,
				Action: config.RelabelKeep,
			},
			{
				SourceLabels: model.LabelNames{model.AddressLabel},
				Regex:        prometheus.KubeletPortRegexp,
				Replacement:  prometheus.NodeExporterPort,
				TargetLabel:  model.AddressLabel,
			},
			{
				TargetLabel: prometheus.AppLabel,
				Replacement: prometheus.NodeExporterAppName,
			},
			{
				TargetLabel: prometheus.ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: prometheus.ClusterTypeLabel,
				Replacement: prometheus.GuestClusterType,
			},
			{
				SourceLabels: model.LabelNames{model.AddressLabel},
				Regex:        prometheus.NodeExporterPortRegexp,
				Replacement:  prometheus.GroupCapture,
				TargetLabel:  prometheus.IPLabel,
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
			// Drop many mounts that are not interesting based on fstype.
			{
				Action:       prometheus.ActionKeep,
				SourceLabels: model.LabelNames{prometheus.MetricFSTypeLabel},
				Regex:        prometheus.MetricDropFStypeRegexp,
			},
			// We care only about systemd state failed, we can drop the rest.
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel, prometheus.MetricSystemdStateLabel},
				Regex:        prometheus.MetricDropSystemdStateRegexp,
			},
			// Drop all systemd units that are connected to docker mounts or networking, we don't care about them.
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel, prometheus.MetricSystemdNameLabel},
				Regex:        prometheus.MetricDropSystemdNameRegexp,
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
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNamespaceLabel, prometheus.KubernetesSDServiceNameLabel},
				Regex:        prometheus.WhitelistRegexp,
				Action:       config.RelabelKeep,
			},
			{
				TargetLabel:  prometheus.AppLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDServiceNameLabel},
			},
			{
				TargetLabel:  prometheus.NamespaceLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDNamespaceLabel},
			},
			{
				TargetLabel: prometheus.ClusterIDLabel,
				Replacement: "xa5ly",
			},
			{
				TargetLabel: prometheus.ClusterTypeLabel,
				Replacement: prometheus.GuestClusterType,
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
			// keep only kube-system cadvisor metrics
			{
				Action:       prometheus.ActionKeep,
				SourceLabels: model.LabelNames{prometheus.MetricExportedNamespaceLabel},
				Regex:        prometheus.KubeSystemGiantswarmNSRegexp,
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
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].APIServer.Host = apiserver
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CAFile = caFile
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.CertFile = crtFile
	TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs[0].TLSConfig.KeyFile = keyFile
	TestConfigTwoWorkload.RelabelConfigs[3].Replacement = clusterID
}
