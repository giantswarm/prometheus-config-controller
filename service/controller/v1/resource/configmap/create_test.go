package configmap

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_ApplyCreateChange tests the ApplyCreateChange method.
func Test_Resource_ConfigMap_ApplyCreateChange(t *testing.T) {
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

	if err := resource.ApplyCreateChange(context.TODO(), v1.Service{}, v1.ConfigMap{}); err != nil {
		t.Fatalf("error returned applying create change: %s\n", err)
	}
}
