package prometheus

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
)

// Test_Prometheus_ConfigMerge tests the ConfigMerge function.
func Test_Prometheus_ConfigMerge(t *testing.T) {
	tests := []struct {
		config         config.Config
		scrapeConfigs  []config.ScrapeConfig
		expectedConfig config.Config
	}{
		// Test an empty config and no scrape configs,
		// returns an empty config.
		{
			config:         config.Config{},
			scrapeConfigs:  []config.ScrapeConfig{},
			expectedConfig: config.Config{},
		},

		// Test an empty config, and one scrape config,
		// returns a config containing the scrape config.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "foo",
				},
			},
			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "foo",
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
						JobName: "foo",
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "foo",
				},
			},
			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "foo",
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
								},
							},
						},
					},
				},
			},
		},

		// Test that adding a scrape config does not affect existing scrape configs.
		{
			config: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "foo",
					},
				},
			},
			scrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "bar",
				},
			},
			expectedConfig: config.Config{
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "foo",
					},
					{
						JobName: "bar",
					},
				},
			},
		},
	}

	for index, test := range tests {
		newConfig, err := ConfigMerge(test.config, test.scrapeConfigs)
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
