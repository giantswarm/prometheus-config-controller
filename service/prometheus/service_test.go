package prometheus

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Prometheus_New tests the New function.
func Test_Prometheus_New(t *testing.T) {
	tests := []struct {
		config func() Config

		expectedErrorHandler func(error) bool
	}{
		// Test that the default config returns an error.
		{
			config: DefaultConfig,

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the kubernetes client must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: nil,
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the logger must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    nil,

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the prometheus address must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap key must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap name must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the configmap namespace must not be empty.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the minimum reload time must not be zero.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "",
					MinimumReloadTime:  time.Duration(0),
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that a valid config produces a service.
		{
			config: func() Config {
				return Config{
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					Address:            "http://127.0.0.1:8080",
					ConfigMapKey:       "prometheus.yml",
					ConfigMapName:      "prometheus",
					ConfigMapNamespace: "monitoring",
					MinimumReloadTime:  2 * time.Minute,
				}
			},

			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		config := test.config()

		service, err := New(config)
		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned creating service: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned creating service: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned creating service\n", index)
		}

		if test.expectedErrorHandler == nil && service == nil {
			t.Fatalf("%d: returned service was nil", index)
		}
	}
}

// Test_Prometheus_Reload tests the Reload method.
func Test_Prometheus_Reload(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		configMap *v1.ConfigMap
		handler   func(w http.ResponseWriter, r *http.Request)

		expectedErrorHandler func(error) bool
	}{
		// Test that an error is returned if the configmap does not exist.
		{
			configMap: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				t.Fatalf("unexpected http request, configmap does not exist")
			},

			expectedErrorHandler: IsReloadError,
		},

		// Test that an error is returned if the configmap does not contain
		// the required key.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				t.Fatalf("unexpected http request, configmap key does not exist")
			},

			expectedErrorHandler: IsReloadError,
		},

		// Test that if the current Prometheus configuration matches the configmap,
		// a reload is not executed.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `foobar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					io.WriteString(w, "<html><pre>foobar</pre></html>")
					return
				}
				t.Fatalf("unexpected http request, reload is not required")
			},

			expectedErrorHandler: nil,
		},

		// Test that if the current Prometheus configuration does not match the configmap,
		// a reload is executed correctly.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `bar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					io.WriteString(w, "<html><pre>foo</pre></html>")
					return
				}
				if r.URL.Path != prometheusReloadPath {
					t.Fatalf("unexpected http request, reload is required")
				}
			},

			expectedErrorHandler: nil,
		},

		// Test that an error is returned if the config route returns an error.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `bar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					http.Error(w, fmt.Sprintf("error getting prometheus config"), http.StatusInternalServerError)
					return
				}
				t.Fatalf("unexpected http request, should only access erroring config route")
			},

			expectedErrorHandler: IsReloadError,
		},

		// Test that an error is returned if the config route returns garbage.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `bar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					io.WriteString(w, "lwnefknfiefnpeijfpqofjqpwofjqpwofjqpwofjpofjwpofjwpeofj")
					return
				}
			},

			expectedErrorHandler: IsReloadError,
		},

		// Test that an error is returned if the config route returns an empty string.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `bar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					io.WriteString(w, "")
					return
				}
			},

			expectedErrorHandler: IsReloadError,
		},

		// Test that an error is returned if the reload route returns an error.
		{
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"prometheus.yml": `bar`,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == prometheusConfigPath {
					io.WriteString(w, "<html><pre>foo</pre></html>")
					return
				}
				if r.URL.Path == prometheusReloadPath {
					http.Error(w, fmt.Sprintf("error reloading prometheus"), http.StatusInternalServerError)
					return
				}
			},

			expectedErrorHandler: IsReloadError,
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		if test.configMap != nil {
			if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(test.configMap); err != nil {
				t.Fatalf("%d: error returned creating configmap: %s\n", index, err)
			}
		}

		testServer := httptest.NewServer(http.HandlerFunc(test.handler))
		defer testServer.Close()

		prometheusReloaderConfig := DefaultConfig()

		prometheusReloaderConfig.K8sClient = fakeK8sClient
		prometheusReloaderConfig.Logger = microloggertest.New()

		prometheusReloaderConfig.Address = testServer.URL
		prometheusReloaderConfig.ConfigMapKey = configMapKey
		prometheusReloaderConfig.ConfigMapName = configMapName
		prometheusReloaderConfig.ConfigMapNamespace = configMapNamespace
		prometheusReloaderConfig.MinimumReloadTime = 1 * time.Second

		service, err := New(prometheusReloaderConfig)
		if err != nil {
			t.Fatalf("error returned creating service: %s\n", err)
		}

		reloadErr := service.Reload()

		if reloadErr != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned reloading prometheus: %s\n", index, reloadErr)
		}
		if reloadErr != nil && !test.expectedErrorHandler(reloadErr) {
			t.Fatalf("%d: incorrect error returned reloading prometheus: %s\n", index, reloadErr)
		}
		if reloadErr == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned reloading prometheus\n", index)
		}
	}
}
