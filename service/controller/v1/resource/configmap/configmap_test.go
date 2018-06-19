package configmap

import (
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
)

// Test_Resource_Configmap_New tests the New function.
func Test_Resource_Configmap_New(t *testing.T) {
	prometheusReloaderConfig := prometheus.DefaultConfig()

	prometheusReloaderConfig.K8sClient = fake.NewSimpleClientset()
	prometheusReloaderConfig.Logger = microloggertest.New()

	prometheusReloaderConfig.Address = "http://127.0.0.1:9090"
	prometheusReloaderConfig.ConfigMapKey = "prometheus.yml"
	prometheusReloaderConfig.ConfigMapName = "prometheus"
	prometheusReloaderConfig.ConfigMapNamespace = "monitoring"
	prometheusReloaderConfig.MinimumReloadTime = 2 * time.Minute

	prometheusReloader, _ := prometheus.New(prometheusReloaderConfig)

	tests := []struct {
		config func() Config

		expectedErrorHandler func(error) bool
	}{
		// Test that the default config returns an error.
		{
			config: func() Config { return Config{} },

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the kubernetes client must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          nil,
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the logger must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             nil,
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the prometheus reloader must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: nil,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the certificate directory must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap key must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap name must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap namespace must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that a valid config produces a configmap resource.
		{
			config: func() Config {
				return Config{
					K8sClient:          fake.NewSimpleClientset(),
					Logger:             microloggertest.New(),
					PrometheusReloader: prometheusReloader,

					CertDirectory:      "/certs",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
				}
			},

			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		config := test.config()

		service, err := New(config)
		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned creating configmap resource: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned creating configmap resource: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned creating configmap resource\n", index)
		}

		if test.expectedErrorHandler == nil && service == nil {
			t.Fatalf("%d: returned configmap resource was nil", index)
		}
	}
}
