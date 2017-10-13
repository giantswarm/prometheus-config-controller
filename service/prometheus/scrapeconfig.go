package prometheus

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/prometheus-config-controller/service/key"
)

const (
	// httpsScheme is the scheme for https connections.
	httpsScheme = "https"
)

// GetTarget takes a Kubernetes Service, and returns a LabelSet,
// suitable for use as a target.
func GetTarget(service v1.Service) model.LabelSet {
	targetName := fmt.Sprintf("%s.%s", service.Name, service.Namespace)
	target := model.LabelSet{model.LabelName(targetName): ""}

	return target
}

// GetScrapeConfigs takes a list of Kubernetes Services,
// and returns a list of Prometheus ScrapeConfigs.
func GetScrapeConfigs(services []v1.Service, certificateDirectory string) ([]config.ScrapeConfig, error) {
	filteredServices := FilterInvalidServices(services)
	groupedServices := GroupServices(filteredServices)

	scrapeConfigs := []config.ScrapeConfig{}
	for groupName, services := range groupedServices {
		targets := []model.LabelSet{}
		for _, service := range services {
			targets = append(targets, GetTarget(service))
		}

		scrapeConfig := config.ScrapeConfig{
			JobName: groupName,
			Scheme:  httpsScheme,
			HTTPClientConfig: config.HTTPClientConfig{
				TLSConfig: config.TLSConfig{
					CAFile:   key.CAPath(certificateDirectory, groupName),
					CertFile: key.CrtPath(certificateDirectory, groupName),
					KeyFile:  key.KeyPath(certificateDirectory, groupName),
				},
			},
			ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
				StaticConfigs: []*config.TargetGroup{
					{
						Targets: targets,
					},
				},
			},
		}

		scrapeConfigs = append(scrapeConfigs, scrapeConfig)
	}

	return scrapeConfigs, nil
}
