package configmap

import (
	"net/url"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
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
					APIServer: config.URL{URL: &url.URL{
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
		MetricRelabelConfigs: []*config.RelabelConfig{
			// drop several bucket latency metric
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricDropBucketLatencies,
			},
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricsDropReflectorRegexp,
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
					APIServer: config.URL{URL: &url.URL{
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
				Regex:        prometheus.NSRegexp,
			},
			// drop cadvisor metrics about container network statistics
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricDropContainerNetworkRegexp,
			},
		},
	}
	TestConfigOneCalicoNode = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-calico-node",
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
					Role: config.KubernetesRolePod,
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
				SourceLabels: model.LabelNames{prometheus.PodSDNamespaceLabel, prometheus.PodSDPodNameLabel},
				Regex:        prometheus.CalicoNodePodRegexp,
				Action:       config.RelabelKeep,
			},
			{
				TargetLabel:  prometheus.AppLabel,
				SourceLabels: model.LabelNames{prometheus.PodSDContainerNameLabel},
			},
			{
				TargetLabel:  prometheus.NamespaceLabel,
				SourceLabels: model.LabelNames{prometheus.PodSDNamespaceLabel},
			},
			{
				TargetLabel:  prometheus.PodNameLabel,
				SourceLabels: model.LabelNames{prometheus.PodSDPodNameLabel},
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
				TargetLabel: prometheus.AddressLabel,
				Replacement: key.APIServiceHost(key.PrefixMaster, "xa5ly"),
			},
			{
				SourceLabels: model.LabelNames{prometheus.PodSDPodNameLabel},
				Regex:        prometheus.CalicoNodePodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CalicoNodeNamespace, key.CalicoNodeMetricPort),
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{},
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
					APIServer: config.URL{URL: &url.URL{
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
		MetricRelabelConfigs: []*config.RelabelConfig{
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricsDropReflectorRegexp,
			},
		},
	}
	TestConfigOneNodeExporter = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-node-exporter",
		Scheme:  "http",
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: config.URL{URL: &url.URL{
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
				Action:       prometheus.ActionDrop,
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
					APIServer: config.URL{URL: &url.URL{
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
				Regex:        prometheus.ServiceWhitelistRegexp,
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
				TargetLabel:  prometheus.PodNameLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
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
				TargetLabel: prometheus.AddressLabel,
				Replacement: key.APIServiceHost(key.PrefixMaster, "xa5ly"),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.KubeStateMetricsPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.KubeStateMetricsNamespace, key.KubeStateMetricsPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.NginxIngressControllerPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NginxIngressControllerNamespace, key.NginxIngressControllerMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.CalicoNodePodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CalicoNodeNamespace, key.CalicoNodeMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.ChartOperatorPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ChartOperatorNamespace, key.ChartOperatorMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.CertExporterPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CertExporterNamespace, key.CertExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.ClusterAutoscalerPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ClusterAutoscalerNamespace, key.ClusterAutoscalerMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.CoreDNSPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CoreDNSNamespace, key.CoreDNSMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.ElasticLoggingPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.ElasticLoggingNamespace, key.ElasticLoggingMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.NetExporterPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NetExporterNamespace, key.NetExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.NicExporterPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NicExporterNamespace, key.NicExporterMetricPort),
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
				Regex:        prometheus.VaultExporterPodNameRegexp,
				TargetLabel:  prometheus.MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.VaultExporterNamespace, key.VaultExporterMetricPort),
			},
		},
		MetricRelabelConfigs: []*config.RelabelConfig{
			{
				Action:       prometheus.ActionRelabel,
				SourceLabels: model.LabelNames{prometheus.MetricExportedNamespaceLabel, prometheus.MetricNamespaceLabel},
				Regex:        prometheus.RelabelNamespaceRegexp,
				Replacement:  prometheus.GroupCapture,
				TargetLabel:  prometheus.ExportedNamespaceLabel,
			},
			// keep only kube-system cadvisor metrics
			{
				Action:       prometheus.ActionKeep,
				SourceLabels: model.LabelNames{prometheus.MetricExportedNamespaceLabel},
				Regex:        prometheus.NSRegexp,
			},
			// drop useless IC metrics
			{
				Action:       prometheus.ActionDrop,
				SourceLabels: model.LabelNames{prometheus.MetricNameLabel},
				Regex:        prometheus.MetricDropICRegexp,
			},
		},
	}
	TestConfigOneManagedApp = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-managed-app",
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
				SourceLabels: model.LabelNames{prometheus.KubernetesSDServiceGiantSwarmMonitoringPresentLabel},
				Regex:        config.MustNewRegexp(`(true)`),
				Action:       config.RelabelKeep,
			},
			{
				SourceLabels: model.LabelNames{prometheus.KubernetesSDServiceGiantSwarmMonitoringLabel},
				Regex:        config.MustNewRegexp(`(true)`),
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
				TargetLabel:  prometheus.PodNameLabel,
				SourceLabels: model.LabelNames{prometheus.KubernetesSDPodNameLabel},
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
				SourceLabels: model.LabelNames{
					model.LabelName(prometheus.NamespaceLabel),
					model.LabelName(prometheus.PodNameLabel),
					prometheus.KubernetesSDServiceGiantSwarmMonitoringPortLabel,
				},
				Regex:       prometheus.ManagedAppSourceRegexp,
				TargetLabel: prometheus.AddressLabel,
				Replacement: key.ManagedAppPodMetricsPath("xa5ly"),
			},
		},
	}
)

var (
	TestConfigTwoApiserver    config.ScrapeConfig
	TestConfigTwoCadvisor     config.ScrapeConfig
	TestConfigTwoCalicoNode   config.ScrapeConfig
	TestConfigTwoKubelet      config.ScrapeConfig
	TestConfigTwoNodeExporter config.ScrapeConfig
	TestConfigTwoWorkload     config.ScrapeConfig
	TestConfigTwoManagedApp   config.ScrapeConfig
)

func init() {
	// Copy base of test data structures. Deep copying of required fields is
	// done further below.
	TestConfigTwoApiserver = TestConfigOneApiserver
	TestConfigTwoCadvisor = TestConfigOneCadvisor
	TestConfigTwoCalicoNode = TestConfigOneCalicoNode
	TestConfigTwoKubelet = TestConfigOneKubelet
	TestConfigTwoNodeExporter = TestConfigOneNodeExporter
	TestConfigTwoWorkload = TestConfigOneWorkload
	TestConfigTwoManagedApp = TestConfigOneManagedApp

	apiServer := config.URL{URL: &url.URL{
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
		TestConfigTwoCalicoNode.JobName = "guest-cluster-0ba9v-calico-node"
		TestConfigTwoCalicoNode.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoCalicoNode.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			&config.KubernetesSDConfig{
				APIServer: apiServer,
				Role:      config.KubernetesRolePod,
				TLSConfig: tlsConfig,
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoCalicoNode.RelabelConfigs = nil
		for _, r := range TestConfigOneCalicoNode.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoCalicoNode.RelabelConfigs = append(TestConfigTwoCalicoNode.RelabelConfigs, &newRelabelConfig)
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

	{
		TestConfigTwoManagedApp.JobName = "guest-cluster-0ba9v-managed-app"
		TestConfigTwoManagedApp.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoManagedApp.ServiceDiscoveryConfig.KubernetesSDConfigs = []*config.KubernetesSDConfig{
			{
				APIServer: apiServer,
				Role:      config.KubernetesRoleEndpoint,
				TLSConfig: tlsConfig,
			},
		}
		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoManagedApp.RelabelConfigs = nil
		for _, r := range TestConfigOneManagedApp.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoManagedApp.RelabelConfigs = append(TestConfigTwoManagedApp.RelabelConfigs, &newRelabelConfig)
		}
	}
}
