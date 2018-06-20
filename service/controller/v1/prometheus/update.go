package prometheus

import (
	"strings"

	"github.com/prometheus/prometheus/config"
)

// UpdateConfig takes an existing Prometheus configuration,
// and a list of Prometheus scrape configurations.
// A new configuration is returned, that includes both the scrape configurations
// in the prior configuration, as well as the new scrape configs.
func UpdateConfig(promcfg config.Config, scrapeConfigs []config.ScrapeConfig) (config.Config, error) {
	desiredScrapeConfigs := []*config.ScrapeConfig{}

	// Make sure to preserve all scrape configs that the prometheus-config-controller does not manage.
	for _, config := range promcfg.ScrapeConfigs {
		if !isManaged(*config) {
			desiredScrapeConfigs = append(desiredScrapeConfigs, config)
		}
	}

	// And append the supplied, desired scrape configs.
	for i, _ := range scrapeConfigs {
		desiredScrapeConfigs = append(desiredScrapeConfigs, &scrapeConfigs[i])
	}

	promcfg.ScrapeConfigs = desiredScrapeConfigs

	return promcfg, nil
}

// isManaged returns true if the given scrape config is managed by the prometheus-config-controller,
// false otherwise.
func isManaged(scrapeConfig config.ScrapeConfig) bool {
	return strings.HasPrefix(scrapeConfig.JobName, jobNamePrefix)
}
