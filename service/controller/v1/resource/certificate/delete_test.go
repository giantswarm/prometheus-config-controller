package certificate

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/spf13/afero"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus/prometheustest"
)

// Test_Resource_Certificate_NewDeletePatch tests the NewDeletePatch method.
func Test_Resource_Certificate_NewDeletePatch(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := Config{}

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertComponentName = "prometheus"
	resourceConfig.CertDirectory = "/certs"
	resourceConfig.CertNamespace = "default"
	resourceConfig.CertPermission = 0644

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	deletePatch, err := resource.NewDeletePatch(context.TODO(), v1.Service{}, []certificateFile{}, []certificateFile{})
	if err != nil {
		t.Fatalf("error returned getting delete patch: %s\n", err)
	}

	if deletePatch != nil {
		t.Fatalf("delete patch should be nil, was: %#v", deletePatch)
	}
}

// Test_Resource_Certificate_ApplyDeleteChange tests the ApplyDeleteChange method.
func Test_Resource_Certificate_ApplyDeleteChange(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := Config{}

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertComponentName = "prometheus"
	resourceConfig.CertDirectory = "/certs"
	resourceConfig.CertNamespace = "default"
	resourceConfig.CertPermission = 0644

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ApplyDeleteChange(context.TODO(), v1.Service{}, []certificateFile{}); err != nil {
		t.Fatalf("error returned applying delete change: %s\n", err)
	}
}
