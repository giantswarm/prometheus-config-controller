package prometheus

import (
	"fmt"
	"sort"
	"strings"

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
func GetTargets(service v1.Service) []model.LabelSet {
	baseTargetName := fmt.Sprintf("%s.%s", service.Name, service.Namespace)

	// Check if the service is annotated with any ports.
	ports := []string{}
	if val, ok := service.ObjectMeta.Annotations[PortAnnotation]; ok {
		ports = strings.Split(val, ",")
	}

	targetNames := []string{}

	// If we have ports specified, append them.
	if len(ports) > 0 {
		for _, port := range ports {
			targetName := fmt.Sprintf("%s:%s", baseTargetName, port)
			targetNames = append(targetNames, targetName)
		}
	} else {
		targetNames = append(targetNames, baseTargetName)
	}

	// And then construct a proper set of labelsets.
	targets := []model.LabelSet{}
	for _, targetName := range targetNames {
		target := model.LabelSet{model.AddressLabel: model.LabelValue(targetName)}
		targets = append(targets, target)
	}

	return targets
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
			targets = append(targets, GetTargets(service)...)
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
