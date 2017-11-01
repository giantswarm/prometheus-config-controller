package configmap

import (
	"context"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_GetDeleteState tests the GetDeleteState method.
func Test_Resource_ConfigMap_GetDeleteState(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	deleteState, err := resource.GetDeleteState(context.TODO(), v1.Service{}, v1.ConfigMap{}, v1.ConfigMap{})
	if err != nil {
		t.Fatalf("error returned getting delete state: %s\n", err)
	}

	if deleteState != nil {
		t.Fatalf("delete state should be nil, was: %#v", deleteState)
	}
}

// Test_Resource_ConfigMap_ProcessDeleteState tests the ProcessDeleteState method.
func Test_Resource_ConfigMap_ProcessDeleteState(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ProcessDeleteState(context.TODO(), v1.Service{}, v1.ConfigMap{}); err != nil {
		t.Fatalf("error returned processing delete state: %s\n", err)
	}
}
