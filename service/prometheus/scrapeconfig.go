package prometheus

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"k8s.io/api/core/v1"

	"github.com/giantswarm/prometheus-config-controller/service/key"
)

const (
	// httpsScheme is the scheme for https connections.
	httpsScheme = "https"

	// jobNamePrefix is the prefix for job names.
	jobNamePrefix = "guest-cluster"
)

// getJobName takes a cluster ID, and returns a suitable job name.
func getJobName(service v1.Service) string {
	return fmt.Sprintf("%s-%s", jobNamePrefix, service.Namespace)
}

// getTargetHost takes a Kubernetes Service, and returns a suitable host.
func getTargetHost(service v1.Service) string {
	return fmt.Sprintf("%s.%s", service.Name, service.Namespace)
}

// getTarget takes a Kubernetes Service, and returns a LabelSet,
// suitable for use as a target.
func getTarget(service v1.Service) model.LabelSet {
	return model.LabelSet{
		model.AddressLabel: model.LabelValue(getTargetHost(service)),
	}
}

// getScrapeConfig takes a Service, and returns a ScrapeConfig.
// It is assumed that filtering has already taken place, and the cluster annotation exists.
func getScrapeConfig(service v1.Service, certificateDirectory string) config.ScrapeConfig {
	clusterID := GetClusterID(service)

	apiServer := config.URL{&url.URL{
		Scheme: httpsScheme,
		Host:   getTargetHost(service),
	}}

	tlsConfig := config.TLSConfig{
		CAFile:   key.CAPath(certificateDirectory, clusterID),
		CertFile: key.CrtPath(certificateDirectory, clusterID),
		KeyFile:  key.KeyPath(certificateDirectory, clusterID),
	}

	clientTlsConfig := tlsConfig
	clientTlsConfig.InsecureSkipVerify = true

	kubernetesTlsConfig := tlsConfig
	kubernetesTlsConfig.InsecureSkipVerify = false

	scrapeConfig := config.ScrapeConfig{
		JobName: getJobName(service),
		Scheme:  httpsScheme,
		HTTPClientConfig: config.HTTPClientConfig{
			TLSConfig: clientTlsConfig,
		},
		ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*config.KubernetesSDConfig{
				{
					APIServer: apiServer,
					Role:      config.KubernetesRoleEndpoint,
					TLSConfig: kubernetesTlsConfig,
				},
				{
					APIServer: apiServer,
					Role:      config.KubernetesRoleNode,
					TLSConfig: kubernetesTlsConfig,
				},
			},
		},
		RelabelConfigs: []*config.RelabelConfig{
			// Add the cluster id label, so we can identify the specific
			// guest cluster.
			{
				TargetLabel: ClusterIDLabel,
				Replacement: clusterID,
				Action:      config.RelabelReplace,
			},
			// Copy the meta service name label to a named label.
			{
				SourceLabels: model.LabelNames{PrometheusServiceNameLabel},
				TargetLabel:  NameLabel,
				Action:       config.RelabelReplace,
			},
			// Copy the meta namespace name label to a named label.
			{
				SourceLabels: model.LabelNames{PrometheusNamespaceLabel},
				TargetLabel:  NamespaceLabel,
				Action:       config.RelabelReplace,
			},
			// Drop any targets that don't match the regexp.
			{
				SourceLabels: model.LabelNames{PrometheusServiceNameLabel},
				Regex:        EndpointRegexp,
				Action:       config.RelabelKeep,
			},
		},
	}

	return scrapeConfig
}

// GetScrapeConfigs takes a list of Kubernetes Services,
// and returns a list of Prometheus ScrapeConfigs.
func GetScrapeConfigs(services []v1.Service, certificateDirectory string) ([]config.ScrapeConfig, error) {
	filteredServices := FilterInvalidServices(services)

	scrapeConfigs := []config.ScrapeConfig{}
	for _, service := range filteredServices {
		scrapeConfigs = append(scrapeConfigs, getScrapeConfig(service, certificateDirectory))
	}

	sort.Slice(scrapeConfigs, func(i, j int) bool {
		return scrapeConfigs[i].JobName < scrapeConfigs[j].JobName
	})

	return scrapeConfigs, nil
}
