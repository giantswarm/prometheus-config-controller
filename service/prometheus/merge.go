package prometheus

import (
	"github.com/prometheus/prometheus/config"
)

// ConfigMerge takes an existing Prometheus configuration,
// and a list of Prometheus scrape configurations.
// A new configuration is returned, that includes both the scrape configurations
// in the prior configuration, as well as the new scrape configs.
func ConfigMerge(promcfg config.Config, scrapeConfigs []config.ScrapeConfig) (config.Config, error) {
	for _, scrapeConfig := range scrapeConfigs {
		presentAlready := false

		// Update existing jobs first.
		for index, existingScrapeConfig := range promcfg.ScrapeConfigs {
			if scrapeConfig.JobName == existingScrapeConfig.JobName {
				presentAlready = true
				promcfg.ScrapeConfigs[index] = &scrapeConfig
			}
		}

		// If the job does not exist, add it.
		if !presentAlready {
			promcfg.ScrapeConfigs = append(promcfg.ScrapeConfigs, &scrapeConfig)
		}
	}

	return promcfg, nil
}
