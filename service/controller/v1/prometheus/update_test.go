package prometheus

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	sd_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/pkg/relabel"
)

// Test_Prometheus_isManaged tests the isManaged function.
func Test_Prometheus_isManaged(t *testing.T) {
	tests := []struct {
		scrapeConfig config.ScrapeConfig
		isManaged    bool
	}{
		{
			scrapeConfig: config.ScrapeConfig{},
			isManaged:    false,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "xa5ly",
				ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
					StaticConfigs: []*targetgroup.Group{
						{
							Targets: []model.LabelSet{
								{"apiserver.xa5ly": ""},
							},
						},
					},
				},
			},
			isManaged: false,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "workload-cluster-xa5ly",
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
						TargetLabel: ClusterIDLabel,
						Replacement: "xa5ly",
					},
				},
			},
			isManaged: true,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "workload-cluster-xa5ly-cadvisor",
				RelabelConfigs: []*relabel.Config{
					{
						TargetLabel: ClusterIDLabel,
						Replacement: "xa5ly",
					},
				},
			},
			isManaged: true,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "workload-cluster-xa5ly-cadvisor",
			},
			isManaged: true,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "management-cluster-gauss",
			},
			isManaged: false,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "management-cluster-gauss-cadvisor",
			},
			isManaged: false,
		},
	}

	for index, test := range tests {
		returnedIsManaged := isManaged(test.scrapeConfig)

		if test.isManaged != returnedIsManaged {
			t.Fatalf(
				"%d: incorrect managed: expected: %t, received: %t, for: \n%s",
				index,
				test.isManaged,
				returnedIsManaged,
				spew.Sdump(test.scrapeConfig),
			)
		}
	}
}

// Test_Prometheus_UpdateConfig tests the UpdateConfig function.
func Test_Prometheus_UpdateConfig(t *testing.T) {
	tests := []struct {
		config        config.Config
		scrapeConfigs []config.ScrapeConfig

		expectedConfig config.Config
	}{
		// Test an empty config, and one scrape config,
		// returns a config containing the scrape config.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "workload-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*relabel.Config{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
		},

		// Test a config containing one scrape config,
		// and given the same scrape config,
		// returns a config containing said scrape config only once.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "workload-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*relabel.Config{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
		},

		// Test a config containing one scrape config,
		// and given two scrape configs - including the old one,
		// returns a config containing both scrape configs.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "workload-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*relabel.Config{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
				{
					JobName: "workload-cluster-jf0sj",
					ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
						KubernetesSDConfigs: []*kubernetes.SDConfig{
							{
								APIServer: config_util.URL{
									URL: &url.URL{
										Scheme: "https",
										Host:   "apiserver.jf02j",
									},
								},
								Role: kubernetes.RoleEndpoint,
							},
						},
					},
					RelabelConfigs: []*relabel.Config{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "jf02j",
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
					{
						JobName: "workload-cluster-jf0sj",
						ServiceDiscoveryConfig: sd_config.ServiceDiscoveryConfig{
							KubernetesSDConfigs: []*kubernetes.SDConfig{
								{
									APIServer: config_util.URL{
										URL: &url.URL{
											Scheme: "https",
											Host:   "apiserver.jf02j",
										},
									},
									Role: kubernetes.RoleEndpoint,
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "jf02j",
							},
						},
					},
				},
			},
		},

		// Test that adding a scrape config does not affect other existing scrape configs.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "workload-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*relabel.Config{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
		},

		// Test that an existing config that contains a cluster scrape job,
		// and an empty list of scrape configs,
		// returns a config that does not include the old cluster scrape job.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "workload-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*relabel.Config{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{},
			},
		},
	}

	for index, test := range tests {
		newConfig, err := UpdateConfig(test.config, test.scrapeConfigs)
		if err != nil {
			t.Fatalf("%d: error returned merging config: %s\n", index, err)
		}

		if !reflect.DeepEqual(test.expectedConfig, newConfig) {
			t.Fatalf(
				"%d: expected config does not match returned config.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedConfig),
				spew.Sdump(newConfig),
			)
		}
	}
}
