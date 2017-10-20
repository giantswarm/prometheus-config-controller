package configmap

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
)

// Test_Resource_ConfigMap_GetDesiredState tests the GetDesiredState method.
func Test_Resource_ConfigMap_GetDesiredState(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		setUpPrometheusConfiguration *config.Config
		setUpConfigMap               *v1.ConfigMap
		setUpServices                []*v1.Service

		expectedPrometheusConfiguration *config.Config
		expectedErrorHandler            func(error) bool
	}{
		// Test that if the configmap does not exist,
		// an error is returned.
		{
			setUpPrometheusConfiguration: nil,
			setUpConfigMap:               nil,
			setUpServices:                nil,

			expectedPrometheusConfiguration: nil,
			expectedErrorHandler:            IsConfigMapNotFound,
		},

		// Test that if the configmap does exist, but is empty,
		// an error is returned.
		{
			setUpPrometheusConfiguration: nil,
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{},
			},
			setUpServices: nil,

			expectedPrometheusConfiguration: nil,
			expectedErrorHandler:            IsConfigMapKeyNotFound,
		},

		// Test that if the configmap does exist, with an invalid config,
		// an error is returned.
		{
			setUpPrometheusConfiguration: nil,
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "lkjwhndwlkjndlkwdnwdn", // Intentionally garbage.
				},
			},
			setUpServices: nil,

			expectedPrometheusConfiguration: nil,
			expectedErrorHandler:            IsInvalidConfigMap,
		},

		// Test that if the configmap does exist, with a valid config,
		// and no services exist,
		// the configmap is returned without modifications.
		// Note - the returned config is marshalled, so some defaults appear.
		{
			setUpPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
				},
			},
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
			},
			setUpServices: nil,

			expectedPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval:     model.Duration(1 * time.Minute),
					ScrapeTimeout:      model.Duration(10 * time.Second),
					EvaluationInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName:        "kubernetes-nodes",
						ScrapeInterval: model.Duration(1 * time.Minute),
						ScrapeTimeout:  model.Duration(10 * time.Second),
						MetricsPath:    "/metrics",
						Scheme:         "http",
					},
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that if the configmap does exist, with a valid config,
		// and a non-annotated service exists,
		// the configmap is returned without modification.
		{
			setUpPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
				},
			},
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
			},
			setUpServices: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
				},
			},

			expectedPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval:     model.Duration(1 * time.Minute),
					ScrapeTimeout:      model.Duration(10 * time.Second),
					EvaluationInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName:        "kubernetes-nodes",
						ScrapeInterval: model.Duration(1 * time.Minute),
						ScrapeTimeout:  model.Duration(10 * time.Second),
						MetricsPath:    "/metrics",
						Scheme:         "http",
					},
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that if the configmap does exist, with a valid config,
		// and an annotated service exists,
		// the configmap is returned with the new service.
		{
			setUpPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName: "kubernetes-nodes",
					},
				},
			},
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
			},
			setUpServices: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
			},

			expectedPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval:     model.Duration(1 * time.Minute),
					ScrapeTimeout:      model.Duration(10 * time.Second),
					EvaluationInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName:        "kubernetes-nodes",
						ScrapeInterval: model.Duration(1 * time.Minute),
						ScrapeTimeout:  model.Duration(10 * time.Second),
						MetricsPath:    "/metrics",
						Scheme:         "http",
					},
					{
						JobName: "xa5ly",
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
								},
							},
						},
					},
				},
			},
			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		if test.setUpPrometheusConfiguration != nil {
			prometheusConfig, err := yaml.Marshal(test.setUpPrometheusConfiguration)
			if err != nil {
				t.Fatalf("%d: error returned marshaling prometheus configuration: %s\n", index, err)
			}

			if test.setUpConfigMap.Data == nil {
				test.setUpConfigMap.Data = map[string]string{}
			}
			test.setUpConfigMap.Data[configMapKey] = string(prometheusConfig)
		}

		if test.setUpConfigMap != nil {
			if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(test.setUpConfigMap); err != nil {
				t.Fatalf("%d: error returned setting up configmap: %s\n", index, err)
			}
		}

		for _, service := range test.setUpServices {
			if _, err := fakeK8sClient.CoreV1().Services(service.Namespace).Create(service); err != nil {
				t.Fatalf("%d: error returned setting up service: %s\n", index, err)
			}
		}

		desiredState, err := resource.GetDesiredState(context.TODO(), v1.Service{})

		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned getting desired state: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned getting desired state: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned getting desired state\n", index)
		}

		if test.expectedPrometheusConfiguration == nil && desiredState != nil {
			t.Fatalf("%d: unexpected configmap returned getting desired state: %s\n", index, spew.Sdump(desiredState))
		}

		if test.expectedPrometheusConfiguration != nil {
			desiredStateConfigMap, ok := desiredState.(*v1.ConfigMap)
			if !ok {
				t.Fatalf("%d: could not cast desired state to configmap: %s\n", index, spew.Sdump(desiredState))
			}

			expectedPrometheusConfigurationBytes, err := yaml.Marshal(test.expectedPrometheusConfiguration)
			if err != nil {
				t.Fatalf("%d: could not marshal expected prometheus configuration: %s\n", index, err)
			}
			expectedPrometheusConfiguration := string(expectedPrometheusConfigurationBytes)

			returnedPrometheusConfiguration, ok := desiredStateConfigMap.Data[configMapKey]
			if !ok {
				t.Fatalf("%d: configuration key not found in desired state configmap: %s\n", index, spew.Sdump(desiredState))
			}

			if expectedPrometheusConfiguration != returnedPrometheusConfiguration {
				t.Fatalf(
					"%d: expected configmap does not match returned desired state.\nexpected:\n%s\nreturned:\n%s\n",
					index,
					expectedPrometheusConfiguration,
					returnedPrometheusConfiguration,
				)
			}
		}
	}
}
