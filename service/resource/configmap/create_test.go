package configmap

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

// Test_Resource_ConfigMap_GetCreateState tests the GetCreateState method.
func Test_Resource_ConfigMap_GetCreateState(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	createState, err := resource.GetCreateState(context.TODO(), v1.Service{}, v1.ConfigMap{}, v1.ConfigMap{})
	if err != nil {
		t.Fatalf("error returned getting create state: %s\n", err)
	}

	if createState != nil {
		t.Fatalf("create state should be nil, was: %#v", createState)
	}
}

// Test_Resource_ConfigMap_ProcessCreateState tests the ProcessCreateState method.
func Test_Resource_ConfigMap_ProcessCreateState(t *testing.T) {
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = "prometheus.yml"
	resourceConfig.ConfigMapName = "prometheus"
	resourceConfig.ConfigMapNamespace = "monitoring"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ProcessCreateState(context.TODO(), v1.Service{}, v1.ConfigMap{}); err != nil {
		t.Fatalf("error returned processing create state: %s\n", err)
	}
}
