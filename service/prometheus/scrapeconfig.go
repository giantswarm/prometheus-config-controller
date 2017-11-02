package prometheus

import (
	"fmt"
	"sort"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/prometheus-config-controller/service/key"
)

const (
	// httpsScheme is the scheme for https connections.
	httpsScheme = "https"

	// jobNamePrefix is the prefix for job names.
	jobNamePrefix = "guest-cluster"
)

// GetTarget takes a Kubernetes Service, and returns a LabelSet,
// suitable for use as a target.
func GetTarget(service v1.Service) model.LabelSet {
	targetName := fmt.Sprintf("%s.%s", service.Name, service.Namespace)
	target := model.LabelSet{model.AddressLabel: model.LabelValue(targetName)}

	return target
}

// GetScrapeConfigs takes a list of Kubernetes Services,
// and returns a list of Prometheus ScrapeConfigs.
func GetScrapeConfigs(services []v1.Service, certificateDirectory string) ([]config.ScrapeConfig, error) {
	filteredServices := FilterInvalidServices(services)
	groupedServices := GroupServices(filteredServices)

	scrapeConfigs := []config.ScrapeConfig{}
	for clusterID, services := range groupedServices {
		targets := []model.LabelSet{}
		for _, service := range services {
			targets = append(targets, GetTarget(service))
		}

		scrapeConfig := config.ScrapeConfig{
			JobName: fmt.Sprintf("%s-%s", jobNamePrefix, clusterID),
			Scheme:  httpsScheme,
			HTTPClientConfig: config.HTTPClientConfig{
				TLSConfig: config.TLSConfig{
					CAFile:   key.CAPath(certificateDirectory, clusterID),
					CertFile: key.CrtPath(certificateDirectory, clusterID),
					KeyFile:  key.KeyPath(certificateDirectory, clusterID),
				},
			},
			ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
				StaticConfigs: []*config.TargetGroup{
					{
						Targets: targets,
						Labels: model.LabelSet{
							ClusterLabel:   "",
							ClusterIDLabel: model.LabelValue(clusterID),
						},
					},
				},
			},
		}

		scrapeConfigs = append(scrapeConfigs, scrapeConfig)
	}

	sort.Slice(scrapeConfigs, func(i, j int) bool {
		return scrapeConfigs[i].JobName < scrapeConfigs[j].JobName
	})

	return scrapeConfigs, nil
}
