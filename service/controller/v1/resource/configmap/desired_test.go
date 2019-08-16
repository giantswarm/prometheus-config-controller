package configmap

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_GetDesiredState tests the GetDesiredState method.
func Test_Resource_ConfigMap_GetDesiredState(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		name                         string
		setUpPrometheusConfiguration *config.Config
		setUpConfigMap               *v1.ConfigMap
		setUpServices                []*v1.Service

		expectedPrometheusConfiguration *config.Config
		expectedErrorHandler            func(error) bool
	}{
		// Test that if the configmap does not exist,
		// an error is returned.
		{
			name:                         "return error when configmap doesn't exist",
			setUpPrometheusConfiguration: nil,
			setUpConfigMap:               nil,
			setUpServices:                nil,

			expectedPrometheusConfiguration: nil,
			expectedErrorHandler:            IsConfigMapNotFound,
		},

		// Test that if the configmap does exist, but is empty,
		// an error is returned.
		{
			name:                         "return error when configmap exists but is empty",
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
			name:                         "return error when configmap exists but is invalid config",
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
			name: "return unchanged configmap when configuration is valid but no services exist",
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
			name: "return unchanged configmap when configuration is valid but no annotated services exist",
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
		// and an annotated service exists, without the app=master label,
		// the configmap is returned without modification.
		{
			name: "return unchanged configmap when configuration is valid but no annotated services exist with app=master label",
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
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that if the configmap does exist, with a valid config,
		// and an annotated service exists,
		// the configmap is returned with the new service.
		{
			name: "return configmap with new service when configuration is valid and annotated service is found",
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
						Labels: map[string]string{
							"app": "master",
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

					&TestConfigOneApiserver,
					&TestConfigOneCadvisor,
					&TestConfigOneKubelet,
					&TestConfigOneNodeExporter,
					&TestConfigOneWorkload,
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that if the configmap exists, with a service already,
		// and the service no longer exists,
		// the configmap is returned without the service.
		{
			name: "return updated configmap without the service when previously configured service doesn't exist anymore",
			setUpPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					{
						JobName:        "kubernetes-nodes",
						ScrapeInterval: model.Duration(1 * time.Minute),
						ScrapeTimeout:  model.Duration(10 * time.Second),
						MetricsPath:    "/metrics",
						Scheme:         "http",
					},

					&TestConfigOneApiserver,
					&TestConfigOneCadvisor,
					&TestConfigOneKubelet,
					&TestConfigOneNodeExporter,
					&TestConfigOneWorkload,
				},
			},
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
			},
			setUpServices: []*v1.Service{},

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

		// Test that if the configmap exists, with a service already,
		// and a new service is added,
		// the configmap is returned with both services.
		{
			name: "return configmap with two services when it previously had only one and another new service is found",
			setUpPrometheusConfiguration: &config.Config{
				GlobalConfig: config.GlobalConfig{
					ScrapeInterval: model.Duration(1 * time.Minute),
				},
				ScrapeConfigs: []*config.ScrapeConfig{
					&TestConfigOneApiserver,
					&TestConfigOneCadvisor,
					&TestConfigOneKubelet,
					&TestConfigOneNodeExporter,
					&TestConfigOneWorkload,
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
						Labels: map[string]string{
							"app": "master",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "0ba9v",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "0ba9v",
						},
						Labels: map[string]string{
							"app": "master",
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
					&TestConfigTwoApiserver,
					&TestConfigTwoCadvisor,
					&TestConfigTwoKubelet,
					&TestConfigTwoNodeExporter,
					&TestConfigTwoWorkload,

					&TestConfigOneApiserver,
					&TestConfigOneCadvisor,
					&TestConfigOneKubelet,
					&TestConfigOneNodeExporter,
					&TestConfigOneWorkload,
				},
			},
			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("case %d: %s", index, test.name), func(t *testing.T) {
			fakeK8sClient := fake.NewSimpleClientset()

			resourceConfig := Config{}

			resourceConfig.K8sClient = fakeK8sClient
			resourceConfig.Logger = microloggertest.New()
			resourceConfig.PrometheusReloader = prometheustest.New()

			resourceConfig.CertDirectory = "/certs"
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
						"%d: expected configmap does not match returned desired state.\ndiff:\n%s\n",
						index,
						cmp.Diff(expectedPrometheusConfiguration, returnedPrometheusConfiguration),
					)
				}
			}
		})
	}
}
