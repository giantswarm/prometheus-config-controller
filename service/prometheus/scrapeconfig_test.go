package prometheus

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Test_Prometheus_GetTargets tests the GetTargets function.
func Test_Prometheus_GetTargets(t *testing.T) {
	tests := []struct {
		service         v1.Service
		expectedTargets []model.LabelSet
	}{
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			expectedTargets: []model.LabelSet{
				{
					model.AddressLabel: "foo.bar",
				},
			},
		},

		// Test that if a single port is annotated, that port is part of the target.
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						PortAnnotation: "30101",
					},
				},
			},
			expectedTargets: []model.LabelSet{
				{
					model.AddressLabel: "foo.bar:30101",
				},
			},
		},

		// Test that if multiple ports are annotated,
		// each port is a separate target.
		{
			service: v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						PortAnnotation: "30101,30102",
					},
				},
			},
			expectedTargets: []model.LabelSet{
				{
					model.AddressLabel: "foo.bar:30101",
				},
				{
					model.AddressLabel: "foo.bar:30102",
				},
			},
		},
	}

	for index, test := range tests {
		targets := GetTargets(test.service)

		if !reflect.DeepEqual(test.expectedTargets, targets) {
			t.Fatalf(
				"%d: expected targets do not match returned targets.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedTargets),
				spew.Sdump(targets),
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
							InsecureSkipVerify: false,
						},
					},
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
								},
								Labels: model.LabelSet{
									ClusterLabel:   "",
									ClusterIDLabel: "xa5ly",
								},
							},
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
							InsecureSkipVerify: false,
						},
					},
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.0ba9v"},
								},
								Labels: model.LabelSet{
									ClusterLabel:   "",
									ClusterIDLabel: "0ba9v",
								},
							},
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
							InsecureSkipVerify: false,
						},
					},
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
								},
								Labels: model.LabelSet{
									ClusterLabel:   "",
									ClusterIDLabel: "xa5ly",
								},
							},
						},
					},
				},
			},
		},

		// Test that two services that specify the same cluster annotation create a scrape config together.
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
						Name:      "kubelet",
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
							InsecureSkipVerify: false,
						},
					},
					ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
						StaticConfigs: []*config.TargetGroup{
							{
								Targets: []model.LabelSet{
									model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
									model.LabelSet{model.AddressLabel: "kubelet.xa5ly"},
								},
								Labels: model.LabelSet{
									ClusterLabel:   "",
									ClusterIDLabel: "xa5ly",
								},
							},
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
					InsecureSkipVerify: false,
				},
			},
			ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
				StaticConfigs: []*config.TargetGroup{
					{
						Targets: []model.LabelSet{
							model.LabelSet{model.AddressLabel: "apiserver.0ba9v"},
						},
						Labels: model.LabelSet{
							ClusterLabel:   "",
							ClusterIDLabel: "0ba9v",
						},
					},
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
					InsecureSkipVerify: false,
				},
			},
			ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
				StaticConfigs: []*config.TargetGroup{
					{
						Targets: []model.LabelSet{
							model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
						},
						Labels: model.LabelSet{
							ClusterLabel:   "",
							ClusterIDLabel: "xa5ly",
						},
					},
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
						InsecureSkipVerify: false,
					},
				},
				ServiceDiscoveryConfig: config.ServiceDiscoveryConfig{
					StaticConfigs: []*config.TargetGroup{
						{
							Targets: []model.LabelSet{
								model.LabelSet{model.AddressLabel: "apiserver.xa5ly"},
							},
							Labels: model.LabelSet{
								ClusterLabel:   "",
								ClusterIDLabel: "xa5ly",
							},
						},
					},
				},
			},

			expectedConfig: `job_name: guest-cluster-xa5ly
scheme: https
static_configs:
- targets:
  - apiserver.xa5ly
  labels:
    cluster_id: xa5ly
    prometheus_config_controller: ""
tls_config:
  ca_file: /certs/xa5ly-ca.pem
  cert_file: /certs/xa5ly-crt.pem
  key_file: /certs/xa5ly-key.pem
  insecure_skip_verify: false
`,
		},
	}

	for index, test := range tests {
		data, err := yaml.Marshal(test.scrapeConfig)
		if err != nil {
			t.Fatalf("%d: error occurred marshaling yaml: %s\n", index, err)
		}

		if test.expectedConfig != string(data) {
			t.Fatalf(
				"%d: expected scrape configs do not match returned scrape configs.\nexpected: %s\nreturned: %s\n",
				index,
				test.expectedConfig,
				string(data),
			)
		}
	}
}
