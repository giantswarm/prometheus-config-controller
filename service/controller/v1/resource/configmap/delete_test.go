package configmap

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_NewDeletePatch tests the NewDeletePatch method.
func Test_Resource_ConfigMap_NewDeletePatch(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := Config{}

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	deletePatch, err := resource.NewDeletePatch(context.TODO(), v1.Service{}, v1.ConfigMap{}, v1.ConfigMap{})
	if err != nil {
		t.Fatalf("error returned getting delete patch: %s\n", err)
	}

	if deletePatch != nil {
		t.Fatalf("delete patch should be nil, was: %#v", deletePatch)
	}
}

// Test_Resource_ConfigMap_ApplyDeleteChange tests the ApplyDeleteChange method.
func Test_Resource_ConfigMap_ApplyDeleteChange(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := Config{}

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ApplyDeleteChange(context.TODO(), v1.Service{}, v1.ConfigMap{}); err != nil {
		t.Fatalf("error returned applying delete patch: %s\n", err)
	}
}
