package configmap

import (
	"context"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_GetCurrentState tests the GetCurrentState method.
func Test_Resource_ConfigMap_GetCurrentState(t *testing.T) {
	configMapName := "prometheus-config"
	configMapNamespace := "monitoring"

	tests := []struct {
		setUp             func(kubernetes.Interface) error
		expectedConfigMap *v1.ConfigMap
	}{
		// Test that when the configmap does not exist, GetCurrentState returns nil.
		{
			setUp: func(k8sClient kubernetes.Interface) error {
				return nil
			},
			expectedConfigMap: nil,
		},

		// Test that when the configmap does exist, GetCurrentState returns the configmap.
		{
			setUp: func(k8sClient kubernetes.Interface) error {
				_, err := k8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      configMapName,
						Namespace: configMapNamespace,
					},
					Data: map[string]string{
						"foo": "bar",
					},
				})

				return err
			},
			expectedConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					"foo": "bar",
				},
			},
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.ConfigMapKey = "prometheus.yml"
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%v: error returned creating resource: %v\n", index, err)
		}

		if err := test.setUp(fakeK8sClient); err != nil {
			t.Fatalf("%v: error returned during setup: %v\n", index, err)
		}

		currentState, err := resource.GetCurrentState(context.TODO(), v1.Service{})
		if err != nil {
			t.Fatalf("%v: error returned getting current state: %v\n", index, err)
		}

		if !(test.expectedConfigMap == nil && currentState == nil) && !reflect.DeepEqual(test.expectedConfigMap, currentState) {
			t.Fatalf(
				"%v: expected configmap does not match returned current state.\nexpected: %v\nreturned: %v\n",
				index,
				test.expectedConfigMap,
				currentState,
			)
		}
	}
}
