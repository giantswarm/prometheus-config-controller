package prometheus

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
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
				ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
					StaticConfigs: []*config.TargetGroup{
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
				JobName: "guest-cluster-xa5ly",
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
						TargetLabel: ClusterIDLabel,
						Replacement: "xa5ly",
					},
				},
			},
			isManaged: true,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "guest-cluster-xa5ly-cadvisor",
				RelabelConfigs: []*config.RelabelConfig{
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
				JobName: "guest-cluster-xa5ly-cadvisor",
			},
			isManaged: true,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "host-cluster-gauss",
			},
			isManaged: false,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "host-cluster-gauss-cadvisor",
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
					JobName: "guest-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
					JobName: "guest-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
					JobName: "guest-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
				{
					JobName: "guest-cluster-jf0sj",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						KubernetesSDConfigs: []*config.KubernetesSDConfig{
							{
								APIServer: config.URL{
									URL: &url.URL{
										Scheme: "https",
										Host:   "apiserver.jf02j",
									},
								},
								Role: config.KubernetesRoleEndpoint,
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
							{
								TargetLabel: ClusterIDLabel,
								Replacement: "xa5ly",
							},
						},
					},
					{
						JobName: "guest-cluster-jf0sj",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							KubernetesSDConfigs: []*config.KubernetesSDConfig{
								{
									APIServer: config.URL{
										URL: &url.URL{
											Scheme: "https",
											Host:   "apiserver.jf02j",
										},
									},
									Role: config.KubernetesRoleEndpoint,
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
					JobName: "guest-cluster-xa5ly",
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
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
						JobName: "guest-cluster-xa5ly",
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
								},
							},
						},
						RelabelConfigs: []*config.RelabelConfig{
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
