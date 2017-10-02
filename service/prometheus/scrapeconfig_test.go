package prometheus

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Test_Prometheus_GetTarget tests the GetTarget function.
func Test_Prometheus_GetTarget(t *testing.T) {
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
			expectedTarget: model.LabelSet{model.LabelName("foo.bar"): ""},
		},
	}

	for index, test := range tests {
		target := GetTarget(test.service)

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
		services              []v1.Service
		expectedScrapeConfigs []config.ScrapeConfig
	}{
		// Test that when there are no services available,
		// no scrape configs are returned.
		{
			services:              nil,
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
			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a service with the cluster annotation,
		// but without a certificate annotation,
		// does not create a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			expectedScrapeConfigs: []config.ScrapeConfig{},
		},

		// Test that a service that specifies both the cluster and certificate
		// annotations creates a scrape config.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation:     "xa5ly",
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
			},
			expectedScrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "xa5ly",
					Scheme:  "https",
					HTTPClientConfig: config.HTTPClientConfig{
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/xa5ly/ca.pem",
							CertFile:           "/certs/xa5ly/crt.pem",
							KeyFile:            "/certs/xa5ly/key.pem",
							InsecureSkipVerify: false,
						},
					},
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

		// Test that two services that specify the same cluster and certificate
		// annotations create a scrape config together.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation:     "xa5ly",
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubelet",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							ClusterAnnotation:     "xa5ly",
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
			},
			expectedScrapeConfigs: []config.ScrapeConfig{
				{
					JobName: "xa5ly",
					Scheme:  "https",
					HTTPClientConfig: config.HTTPClientConfig{
						TLSConfig: config.TLSConfig{
							CAFile:             "/certs/xa5ly/ca.pem",
							CertFile:           "/certs/xa5ly/crt.pem",
							KeyFile:            "/certs/xa5ly/key.pem",
							InsecureSkipVerify: false,
						},
					},
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
	}

	for index, test := range tests {
		scrapeConfigs, err := GetScrapeConfigs(test.services)
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
