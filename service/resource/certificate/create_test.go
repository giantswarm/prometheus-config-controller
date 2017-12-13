package certificate

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus/prometheustest"
)

// Test_Resource_Certificate_ApplyCreateChange tests the ApplyCreateChange method.
func Test_Resource_Certificate_ApplyCreateChange(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheustest.New()

	resourceConfig.CertificateComponentName = "prometheus"
	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.CertificateNamespace = "default"
	resourceConfig.CertificatePermission = 0644

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ApplyCreateChange(context.TODO(), v1.Service{}, []certificateFile{}); err != nil {
		t.Fatalf("error returned applying create change: %s\n", err)
	}
}
