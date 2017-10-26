package prometheus

import (
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
								model.LabelSet{"apiserver.xa5ly": ""},
							},
						},
					},
				},
			},
			isManaged: false,
		},

		{
			scrapeConfig: config.ScrapeConfig{
				JobName: "xa5ly",
				ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
					StaticConfigs: []*config.TargetGroup{
						{
							Targets: []model.LabelSet{
								model.LabelSet{"apiserver.xa5ly": ""},
							},
							Labels: model.LabelSet{ClusterLabel: ""},
						},
					},
				},
			},
			isManaged: true,
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
					JobName: "xa5ly",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{"apiserver.xa5ly": ""},
								},
								Labels: model.LabelSet{ClusterLabel: ""},
							},
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{"apiserver.xa5ly": ""},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
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
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{"apiserver.xa5ly": ""},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "xa5ly",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{"apiserver.xa5ly": ""},
								},
								Labels: model.LabelSet{ClusterLabel: ""},
							},
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{"apiserver.xa5ly": ""},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
				},
			},
		},

		// Test a config containing one scrape config,
		// and given a scrape config with the same name but different values,
		// returns a config containing the new scrape config.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{"apiserver.xa5ly": ""},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "xa5ly",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{"apiserver.xa5ly": ""},
									model.LabelSet{"kubelet.xa5ly": ""},
								},
								Labels: model.LabelSet{ClusterLabel: ""},
							},
						},
					},
				},
			},

			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{"apiserver.xa5ly": ""},
										model.LabelSet{"kubelet.xa5ly": ""},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
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
					JobName: "xa5ly",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
								},
								Labels: model.LabelSet{ClusterLabel: ""},
							},
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
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
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
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
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

		// Test that an existing config that contains a non-managed job,
		// and two existing cluster scrape jobs,
		// and one of the scrape jobs is removed, and another updated,
		// returns a config that returns the non-managed job, and the updated job.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
					{
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
					{
						JobName: "ru85y",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{model.AddressLabel: "apiserver.ru85y"},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "xa5ly",
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
									model.LabelSet{model.AddressLabel: "kube-state-metrics.xa5ly"},
								},
								Labels: model.LabelSet{ClusterLabel: ""},
							},
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
						JobName: "xa5ly",
						ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
							StaticConfigs: []*config.TargetGroup{
								{
									Targets: []model.LabelSet{
										model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
										model.LabelSet{model.AddressLabel: "kube-state-metrics.xa5ly"},
									},
									Labels: model.LabelSet{ClusterLabel: ""},
								},
							},
						},
					},
				},
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
