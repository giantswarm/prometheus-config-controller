package prometheus

import (
	"net/url"
	"strings"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/key"
	"github.com/prometheus/common/model"

	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/config"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/pkg/relabel"
)

var (
	TestConfigOneApiserver = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-apiserver",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: true,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{
					KubernetesSDNamespaceLabel,
					KubernetesSDServiceNameLabel,
				},
				Regex:  APIServerRegexp,
				Action: relabel.Keep,
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
		MetricRelabelConfigs: []*relabel.Config{
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
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneAWSNode = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-aws-node",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RolePod,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{PodSDNamespaceLabel, PodSDPodNameLabel},
				Regex:        AWSNodePodRegexp,
				Action:       relabel.Keep,
			},
			{
				TargetLabel:  AppLabel,
				SourceLabels: model.LabelNames{PodSDContainerNameLabel},
			},
			{
				TargetLabel:  NamespaceLabel,
				SourceLabels: model.LabelNames{PodSDNamespaceLabel},
			},
			{
				TargetLabel:  PodNameLabel,
				SourceLabels: model.LabelNames{PodSDPodNameLabel},
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
				SourceLabels: model.LabelNames{PodSDPodNameLabel},
				Regex:        AWSNodePodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.AWSNodeNamespace, key.AWSNodeMetricPort),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneCadvisor = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-cadvisor",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleNode,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
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
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneCalicoNode = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-calico-node",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RolePod,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{PodSDNamespaceLabel, PodSDPodNameLabel},
				Regex:        CalicoNodePodRegexp,
				Action:       relabel.Keep,
			},
			{
				TargetLabel:  AppLabel,
				SourceLabels: model.LabelNames{PodSDContainerNameLabel},
			},
			{
				TargetLabel:  NamespaceLabel,
				SourceLabels: model.LabelNames{PodSDNamespaceLabel},
			},
			{
				TargetLabel:  PodNameLabel,
				SourceLabels: model.LabelNames{PodSDPodNameLabel},
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
				SourceLabels: model.LabelNames{PodSDPodNameLabel},
				Regex:        CalicoNodePodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CalicoNodeNamespace, key.CalicoNodeMetricPort),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneDocker = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-docker-daemon",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleNode,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				TargetLabel: model.AddressLabel,
				Replacement: "apiserver.xa5ly",
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDNodeNameLabel},
				Replacement:  DockerMetricsPath,
				TargetLabel:  model.MetricsPathLabel,
			},
			{
				TargetLabel: AppLabel,
				Replacement: DockerAppName,
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
		MetricRelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        DockerMetricsNameRegexp,
				Action:       ActionKeep,
			},
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneKubelet = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-kubelet",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: true,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleNode,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
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
		MetricRelabelConfigs: []*relabel.Config{
			{
				Action:       ActionDrop,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricsDropReflectorRegexp,
			},
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneNodeExporter = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-node-exporter",
		Scheme:  "http",
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{
					KubernetesSDNamespaceLabel,
					KubernetesSDServiceNameLabel,
				},
				Regex:  NodeExporterRegexp,
				Action: relabel.Keep,
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
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}

	TestConfigOneKubeProxy = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-kube-proxy",
		Scheme:  "https",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RolePod,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
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
				Regex:        KubeProxyPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.KubeProxyNamespace, key.KubeProxyMetricPort),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
			// keep only kube-proxy iptables restore errors metrics
			{
				Action:       ActionKeep,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricsKeepKubeProxyIptableRegexp,
			},
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}

	TestConfigOneWorkload = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-workload",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		Scheme: "https",
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
				Regex:        ServiceWhitelistRegexp,
				Action:       relabel.Keep,
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel, PodSDGiantswarmServiceTypeLabel},
				Regex:        KiamPodNameRegexpNonManaged,
				Action:       relabel.Drop,
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
				TargetLabel:  NodeLabel,
				SourceLabels: model.LabelNames{KubernetesSDPodNodeNameLabel},
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
				Regex:        CalicoNodePodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.CalicoNodeNamespace, key.CalicoNodeMetricPort),
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
				Regex:        KiamPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.KiamNamespace, key.KiamMetricPort),
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDPodNameLabel},
				Regex:        VaultExporterPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.VaultExporterNamespace, key.VaultExporterMetricPort),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
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
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneIngress = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-ingress",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		Scheme: "https",
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
				Regex:        IngressWhitelistRegexp,
				Action:       relabel.Keep,
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
				TargetLabel:  NodeLabel,
				SourceLabels: model.LabelNames{KubernetesSDPodNodeNameLabel},
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
				Regex:        NginxIngressControllerPodNameRegexp,
				TargetLabel:  MetricPathLabel,
				Replacement:  key.APIProxyPodMetricsPath(key.NginxIngressControllerNamespace, key.NginxIngressControllerMetricPort),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
			{
				Action:       ActionRelabel,
				SourceLabels: model.LabelNames{MetricExportedNamespaceLabel, MetricNamespaceLabel},
				Regex:        RelabelNamespaceRegexp,
				Replacement:  GroupCapture,
				TargetLabel:  ExportedNamespaceLabel,
			},
			{
				Action:       ActionKeep,
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        MetricKeepICRegexp,
			},
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneManagedApp = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-managed-app",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		Scheme: "https",
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPresentLabel},
				Regex:        relabel.MustNewRegexp(`(true)`),
				Action:       relabel.Keep,
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringLabel},
				Regex:        relabel.MustNewRegexp(`(true)`),
				Action:       relabel.Keep,
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPortPresentLabel},
				Regex:        relabel.MustNewRegexp(`(true)`),
				Action:       relabel.Keep,
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPathPresentLabel},
				Regex:        relabel.MustNewRegexp(`(true)`),
				Action:       relabel.Keep,
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
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringAppTypeLabel},
				Regex:        relabel.MustNewRegexp(`(optional|default)`),
				TargetLabel:  AppTypeLabel,
			},
			{
				SourceLabels: model.LabelNames{KubernetesSDServiceGiantSwarmMonitoringPresentLabel},
				Regex:        relabel.MustNewRegexp(`(true)`),
				TargetLabel:  AppIsManaged,
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
				SourceLabels: model.LabelNames{
					model.LabelName(NamespaceLabel),
					model.LabelName(PodNameLabel),
					KubernetesSDServiceGiantSwarmMonitoringPortLabel,
					KubernetesSDServiceGiantSwarmMonitoringPathLabel},
				Regex:       ManagedAppSourceRegexp,
				TargetLabel: MetricPathLabel,
				Replacement: key.ManagedAppPodMetricsPath(),
			},
		},
		MetricRelabelConfigs: []*relabel.Config{
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
	TestConfigOneKubeStateManagedApp = config.ScrapeConfig{
		JobName: "guest-cluster-xa5ly-kube-state-managed-app",
		HTTPClientConfig: config_util.HTTPClientConfig{
			TLSConfig: config_util.TLSConfig{
				CAFile:             "/certs/xa5ly-ca.pem",
				CertFile:           "/certs/xa5ly-crt.pem",
				KeyFile:            "/certs/xa5ly-key.pem",
				InsecureSkipVerify: false,
			},
		},
		Scheme: "https",
		ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*kubernetes.SDConfig{
				{
					APIServer: config_util.URL{
						URL: &url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						},
					},
					Role: kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: config_util.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
		},
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{KubernetesSDNamespaceLabel, KubernetesSDServiceNameLabel},
				Regex:        KubeStateMetricsServiceNameRegexp,
				Action:       relabel.Keep,
			},
			{
				TargetLabel: KubeStateMetricsForManagedApps,
				Replacement: "true",
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
		},
		MetricRelabelConfigs: []*relabel.Config{
			{
				SourceLabels: model.LabelNames{MetricNameLabel},
				Regex:        KubeStateMetricsManagedAppMetricsNameRegexp,
				Action:       ActionKeep,
			},
			{
				SourceLabels: model.LabelNames{MetricExportedNamespaceLabel},
				TargetLabel:  NamespaceLabel,
			},
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
			{
				TargetLabel: ProviderLabel,
				Replacement: "aws-test",
			},
		},
	}
)

var (
	TestConfigTwoApiserver           config.ScrapeConfig
	TestConfigTwoAWSNode             config.ScrapeConfig
	TestConfigTwoCadvisor            config.ScrapeConfig
	TestConfigTwoCalicoNode          config.ScrapeConfig
	TestConfigTwoDocker              config.ScrapeConfig
	TestConfigTwoKubelet             config.ScrapeConfig
	TestConfigTwoNodeExporter        config.ScrapeConfig
	TestConfigTwoWorkload            config.ScrapeConfig
	TestConfigTwoIngress             config.ScrapeConfig
	TestConfigTwoManagedApp          config.ScrapeConfig
	TestConfigTwoKubeStateManagedApp config.ScrapeConfig
	TestConfigTwoKubeProxy           config.ScrapeConfig
)

func init() {
	// Copy base of test data structures. Deep copying of required fields is
	// done further below.
	TestConfigTwoApiserver = TestConfigOneApiserver
	TestConfigTwoAWSNode = TestConfigOneAWSNode
	TestConfigTwoCadvisor = TestConfigOneCadvisor
	TestConfigTwoCalicoNode = TestConfigOneCalicoNode
	TestConfigTwoDocker = TestConfigOneDocker
	TestConfigTwoKubelet = TestConfigOneKubelet
	TestConfigTwoNodeExporter = TestConfigOneNodeExporter
	TestConfigTwoWorkload = TestConfigOneWorkload
	TestConfigTwoIngress = TestConfigOneIngress
	TestConfigTwoManagedApp = TestConfigOneManagedApp
	TestConfigTwoKubeStateManagedApp = TestConfigOneKubeStateManagedApp
	TestConfigTwoKubeProxy = TestConfigOneKubeProxy

	apiServer := config_util.URL{URL: &url.URL{
		Scheme: "https",
		Host:   "apiserver.0ba9v",
	}}

	clusterID := "0ba9v"

	tlsConfig := config_util.TLSConfig{
		CAFile:             "/certs/0ba9v-ca.pem",
		CertFile:           "/certs/0ba9v-crt.pem",
		KeyFile:            "/certs/0ba9v-key.pem",
		InsecureSkipVerify: false,
	}

	{
		TestConfigTwoApiserver.JobName = "guest-cluster-0ba9v-apiserver"
		TestConfigTwoApiserver.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoApiserver.HTTPClientConfig.TLSConfig.InsecureSkipVerify = true
		TestConfigTwoApiserver.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleEndpoint,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoAWSNode.JobName = "guest-cluster-0ba9v-aws-node"
		TestConfigTwoAWSNode.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoAWSNode.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RolePod,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoAWSNode.RelabelConfigs = nil
		for _, r := range TestConfigOneAWSNode.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoAWSNode.RelabelConfigs = append(TestConfigTwoAWSNode.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoCadvisor.JobName = "guest-cluster-0ba9v-cadvisor"
		TestConfigTwoCadvisor.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoCadvisor.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleNode,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoDocker.JobName = "guest-cluster-0ba9v-docker-daemon"
		TestConfigTwoDocker.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoDocker.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleNode,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoDocker.RelabelConfigs = nil
		for _, r := range TestConfigOneDocker.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoDocker.RelabelConfigs = append(TestConfigTwoDocker.RelabelConfigs, &newRelabelConfig)
		}
	}

	{
		TestConfigTwoCalicoNode.JobName = "guest-cluster-0ba9v-calico-node"
		TestConfigTwoCalicoNode.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoCalicoNode.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RolePod,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoKubelet.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleNode,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoNodeExporter.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleEndpoint,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoWorkload.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleEndpoint,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
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
		TestConfigTwoIngress.JobName = "guest-cluster-0ba9v-ingress"
		TestConfigTwoIngress.HTTPClientConfig.TLSConfig = tlsConfig
		TestConfigTwoIngress.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
			{
				APIServer: apiServer,
				Role:      kubernetes.RoleEndpoint,
				HTTPClientConfig: config_util.HTTPClientConfig{
					TLSConfig: tlsConfig,
				},
			},
		}

		// Deepcopy relabel configs and change clusterID to match second test config.
		TestConfigTwoIngress.RelabelConfigs = nil
		for _, r := range TestConfigOneIngress.RelabelConfigs {
			newRelabelConfig := *r
			newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
			TestConfigTwoIngress.RelabelConfigs = append(TestConfigTwoIngress.RelabelConfigs, &newRelabelConfig)
		}
	}
	{
		{
			TestConfigTwoManagedApp.JobName = "guest-cluster-0ba9v-managed-app"
			TestConfigTwoManagedApp.HTTPClientConfig.TLSConfig = tlsConfig
			TestConfigTwoManagedApp.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
				{
					APIServer: apiServer,
					Role:      kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: tlsConfig,
					},
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

	{
		{
			TestConfigTwoKubeStateManagedApp.JobName = "guest-cluster-0ba9v-kube-state-managed-app"
			TestConfigTwoKubeStateManagedApp.HTTPClientConfig.TLSConfig = tlsConfig
			TestConfigTwoKubeStateManagedApp.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
				{
					APIServer: apiServer,
					Role:      kubernetes.RoleEndpoint,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: tlsConfig,
					},
				},
			}

			// Deepcopy relabel configs and change clusterID to match second test config.
			TestConfigTwoKubeStateManagedApp.RelabelConfigs = nil
			for _, r := range TestConfigOneKubeStateManagedApp.RelabelConfigs {
				newRelabelConfig := *r
				newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
				TestConfigTwoKubeStateManagedApp.RelabelConfigs = append(TestConfigTwoKubeStateManagedApp.RelabelConfigs, &newRelabelConfig)
			}
		}
	}

	{
		{
			TestConfigTwoKubeProxy.JobName = "guest-cluster-0ba9v-kube-proxy"
			TestConfigTwoKubeProxy.HTTPClientConfig.TLSConfig = tlsConfig
			TestConfigTwoKubeProxy.ServiceDiscoveryConfig.KubernetesSDConfigs = []*kubernetes.SDConfig{
				{
					APIServer: apiServer,
					Role:      kubernetes.RolePod,
					HTTPClientConfig: config_util.HTTPClientConfig{
						TLSConfig: tlsConfig,
					},
				},
			}

			// Deepcopy relabel configs and change clusterID to match second test config.
			TestConfigTwoKubeProxy.RelabelConfigs = nil
			for _, r := range TestConfigOneKubeProxy.RelabelConfigs {
				newRelabelConfig := *r
				newRelabelConfig.Replacement = strings.ReplaceAll(r.Replacement, "xa5ly", clusterID)
				TestConfigTwoKubeProxy.RelabelConfigs = append(TestConfigTwoKubeProxy.RelabelConfigs, &newRelabelConfig)
			}
		}
	}
}
