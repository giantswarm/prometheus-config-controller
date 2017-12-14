package prometheus

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test_Prometheus_getJobName tests the getJobName function.
func Test_Prometheus_getJobName(t *testing.T) {
	tests := []struct {
		service         v1.Service
		expectedJobName string
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedJobName: "guest-cluster-bar",
		},
	}

	for index, test := range tests {
		jobName := getJobName(test.service)

		if test.expectedJobName != jobName {
			t.Fatalf(
				"%d: expected job name does not match returned job name.\nexpected: %s\nreturned: %s\n",
				index,
				test.expectedJobName,
				jobName,
			)
		}
	}
}

// Test_Prometheus_getTargetHost tests the getTargetHost function.
func Test_Prometheus_getTargetHost(t *testing.T) {
	tests := []struct {
		service            v1.Service
		expectedTargetHost string
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedTargetHost: "foo.bar",
		},
	}

	for index, test := range tests {
		targetHost := getTargetHost(test.service)

		if test.expectedTargetHost != targetHost {
			t.Fatalf(
				"%d: expected target host does not match returned target host.\nexpected: %s\nreturned: %s\n",
				index,
				test.expectedTargetHost,
				targetHost,
			)
		}
	}
}

// Test_Prometheus_getTarget tests the getTarget function.
func Test_Prometheus_getTarget(t *testing.T) {
	tests := []struct {
		service        v1.Service
		expectedTarget model.LabelSet
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedTarget: model.LabelSet{model.AddressLabel: "foo.bar"},
		},
	}

	for index, test := range tests {
		target := getTarget(test.service)

		if !reflect.DeepEqual(test.expectedTarget, target) {
			t.Fatalf(
				"%d: expected target does not match returned target.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedTarget),
				spew.Sdump(target),
			)
		}
	}
}

// Test_Prometheus_GetScrapeConfigs tests the GetScrapeConfigs function.
func Test_Prometheus_GetScrapeConfigs(t *testing.T) {
	tests := []struct {
		services             []v1.Service
		certificateDirectory string

		expectedScrapeConfigs []config.ScrapeConfig
	}{
		// Test that when there are no services available,
		// no scrape configs are returned.
		{
			services:             nil,
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a non-annotated service does not create a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a service that specifies the cluster annotation creates a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{
				{
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
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.xa5ly",
								}},
								Role: config.KubernetesRoleEndpoint,
								TLSConfig: config.TLSConfig{
									CAFile:             "/certs/xa5ly-ca.pem",
									CertFile:           "/certs/xa5ly-crt.pem",
									KeyFile:            "/certs/xa5ly-key.pem",
									InsecureSkipVerify: false,
								},
							},
							{
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.xa5ly",
								}},
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
							TargetLabel: ClusterLabel,
							Replacement: ClusterLabel,
						},
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
			},
		},

		// Test that two services that specify different clusters create separate configs.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "0ba9v",
						Annotations: map[string]string{
							ClusterAnnotation: "0ba9v",
						},
					},
				},
			},
			certificateDirectory: "/certs",

			expectedScrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "guest-cluster-0ba9v",
					Scheme:  "https",
					HTTPClientConfig: config.HTTPClientConfig{
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/0ba9v-ca.pem",
							CertFile:           "/certs/0ba9v-crt.pem",
							KeyFile:            "/certs/0ba9v-key.pem",
							InsecureSkipVerify: true,
						},
					},
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						KubernetesSDConfigs: []*config.KubernetesSDConfig{
							{
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.0ba9v",
								}},
								Role: config.KubernetesRoleEndpoint,
								TLSConfig: config.TLSConfig{
									CAFile:             "/certs/0ba9v-ca.pem",
									CertFile:           "/certs/0ba9v-crt.pem",
									KeyFile:            "/certs/0ba9v-key.pem",
									InsecureSkipVerify: false,
								},
							},
							{
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.0ba9v",
								}},
								Role: config.KubernetesRoleNode,
								TLSConfig: config.TLSConfig{
									CAFile:             "/certs/0ba9v-ca.pem",
									CertFile:           "/certs/0ba9v-crt.pem",
									KeyFile:            "/certs/0ba9v-key.pem",
									InsecureSkipVerify: false,
								},
							},
						},
					},
					RelabelConfigs: []*config.RelabelConfig{
						{
							TargetLabel: ClusterLabel,
							Replacement: ClusterLabel,
						},
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "0ba9v",
						},
					},
				},
				{
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
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.xa5ly",
								}},
								Role: config.KubernetesRoleEndpoint,
								TLSConfig: config.TLSConfig{
									CAFile:             "/certs/xa5ly-ca.pem",
									CertFile:           "/certs/xa5ly-crt.pem",
									KeyFile:            "/certs/xa5ly-key.pem",
									InsecureSkipVerify: false,
								},
							},
							{
								APIServer: config.URL{&url.URL{
									Scheme: "https",
									Host:   "apiserver.xa5ly",
								}},
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
							TargetLabel: ClusterLabel,
							Replacement: ClusterLabel,
						},
						{
							TargetLabel: ClusterIDLabel,
							Replacement: "xa5ly",
						},
					},
				},
			},
		},
	}

	for index, test := range tests {
		scrapeConfigs, err := GetScrapeConfigs(test.services, test.certificateDirectory)
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !reflect.DeepEqual(test.expectedScrapeConfigs, scrapeConfigs) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedScrapeConfigs),
				spew.Sdump(scrapeConfigs),
			)
		}
	}
}

// Test_Prometheus_GetScrapeConfigs_Deterministic tests that the GetScrapeConfigs function is deterministic,
// and that scrape configs are returned in alphabetical order by the job name.
func Test_Prometheus_GetScrapeConfigs_Deterministic(t *testing.T) {
	services := []v1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apiserver",
				Namespace: "xa5ly",
				Annotations: map[string]string{
					ClusterAnnotation: "xa5ly",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apiserver",
				Namespace: "0ba9v",
				Annotations: map[string]string{
					ClusterAnnotation: "0ba9v",
				},
			},
		},
	}

	expectedScrapeConfigs := []config.ScrapeConfig{
		{
			JobName: "guest-cluster-0ba9v",
			Scheme:  "https",
			HTTPClientConfig: config.HTTPClientConfig{
				TLSConfig: config.TLSConfig{
					CAFile:             "/certs/0ba9v-ca.pem",
					CertFile:           "/certs/0ba9v-crt.pem",
					KeyFile:            "/certs/0ba9v-key.pem",
					InsecureSkipVerify: true,
				},
			},
			ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
				KubernetesSDConfigs: []*config.KubernetesSDConfig{
					{
						APIServer: config.URL{&url.URL{
							Scheme: "https",
							Host:   "apiserver.0ba9v",
						}},
						Role: config.KubernetesRoleEndpoint,
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/0ba9v-ca.pem",
							CertFile:           "/certs/0ba9v-crt.pem",
							KeyFile:            "/certs/0ba9v-key.pem",
							InsecureSkipVerify: false,
						},
					},
					{
						APIServer: config.URL{&url.URL{
							Scheme: "https",
							Host:   "apiserver.0ba9v",
						}},
						Role: config.KubernetesRoleNode,
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/0ba9v-ca.pem",
							CertFile:           "/certs/0ba9v-crt.pem",
							KeyFile:            "/certs/0ba9v-key.pem",
							InsecureSkipVerify: false,
						},
					},
				},
			},
			RelabelConfigs: []*config.RelabelConfig{
				{
					TargetLabel: ClusterLabel,
					Replacement: ClusterLabel,
				},
				{
					TargetLabel: ClusterIDLabel,
					Replacement: "0ba9v",
				},
			},
		},
		{
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
						APIServer: config.URL{&url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						}},
						Role: config.KubernetesRoleEndpoint,
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/xa5ly-ca.pem",
							CertFile:           "/certs/xa5ly-crt.pem",
							KeyFile:            "/certs/xa5ly-key.pem",
							InsecureSkipVerify: false,
						},
					},
					{
						APIServer: config.URL{&url.URL{
							Scheme: "https",
							Host:   "apiserver.xa5ly",
						}},
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
					TargetLabel: ClusterLabel,
					Replacement: ClusterLabel,
				},
				{
					TargetLabel: ClusterIDLabel,
					Replacement: "xa5ly",
				},
			},
		},
	}

	for index := 0; index < 50; index++ {
		scrapeConfigs, err := GetScrapeConfigs(services, "/certs")
		if err != nil {
			t.Fatalf("%d: error returned creating scrape configs: %s\n", index, err)
		}

		if !reflect.DeepEqual(expectedScrapeConfigs, scrapeConfigs) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs. GetScrapeConfigs not deterministic.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(expectedScrapeConfigs),
				spew.Sdump(scrapeConfigs),
			)
		}
	}
}

// Test_Prometheus_YamlMarshal tests that Prometheus marshals yaml correctly.
func Test_Prometheus_YamlMarshal(t *testing.T) {
	tests := []struct {
		scrapeConfig config.ScrapeConfig

		expectedConfig string
	}{
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
							APIServer: config.URL{&url.URL{
								Scheme: "https",
								Host:   "apiserver.xa5ly",
							}},
							Role: config.KubernetesRoleEndpoint,
							TLSConfig: config.TLSConfig{
								CAFile:             "/certs/xa5ly-ca.pem",
								CertFile:           "/certs/xa5ly-crt.pem",
								KeyFile:            "/certs/xa5ly-key.pem",
								InsecureSkipVerify: false,
							},
						},
						{
							APIServer: config.URL{&url.URL{
								Scheme: "https",
								Host:   "apiserver.xa5ly",
							}},
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
						TargetLabel: ClusterLabel,
						Replacement: ClusterLabel,
					},
					{
						TargetLabel: ClusterIDLabel,
						Replacement: "xa5ly",
					},
				},
			},

			expectedConfig: `job_name: guest-cluster-xa5ly
scheme: https
kubernetes_sd_configs:
- api_server: https://apiserver.xa5ly
  role: endpoints
  tls_config:
    ca_file: /certs/xa5ly-ca.pem
    cert_file: /certs/xa5ly-crt.pem
    key_file: /certs/xa5ly-key.pem
    insecure_skip_verify: false
- api_server: https://apiserver.xa5ly
  role: node
  tls_config:
    ca_file: /certs/xa5ly-ca.pem
    cert_file: /certs/xa5ly-crt.pem
    key_file: /certs/xa5ly-key.pem
    insecure_skip_verify: false
tls_config:
  ca_file: /certs/xa5ly-ca.pem
  cert_file: /certs/xa5ly-crt.pem
  key_file: /certs/xa5ly-key.pem
  insecure_skip_verify: true
relabel_configs:
- source_labels: []
  target_label: prometheus_config_controller
  replacement: prometheus_config_controller
- source_labels: []
  target_label: cluster_id
  replacement: xa5ly
`,
		},
	}

	for index, test := range tests {
		data, err := yaml.Marshal(test.scrapeConfig)
		if err != nil {
			t.Fatalf("%d: error occurred marshaling yaml: %s\n", index, err)
		}

		expectedLines := strings.Split(test.expectedConfig, "\n")
		returnedLines := strings.Split(string(data), "\n")

		if len(expectedLines) == len(returnedLines) {
			for i := 0; i < len(expectedLines); i++ {
				if expectedLines[i] != returnedLines[i] {
					t.Logf("\nexpected line:\n'%s'\nreturned line:\n'%s'\n", expectedLines[i], returnedLines[i])
				}
			}
		}

		if test.expectedConfig != string(data) {
			t.Fatalf(
				"%d: expected scrape config does not match returned scrape config.\nexpected:\n%s\nreturned:\n%s\n",
				index,
				test.expectedConfig,
				string(data),
			)
		}
	}
}
