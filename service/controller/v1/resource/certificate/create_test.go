package certificate

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/spf13/afero"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// Test_Resource_Certificate_ApplyCreateChange tests the ApplyCreateChange method.
func Test_Resource_Certificate_ApplyCreateChange(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := Config{}

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertComponentName = "prometheus"
	resourceConfig.CertDirectory = "/certs"
	resourceConfig.CertNamespace = "default"
	resourceConfig.CertPermission = 0644

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ApplyCreateChange(context.TODO(), v1.Service{}, []certificateFile{}); err != nil {
		t.Fatalf("error returned applying create change: %s\n", err)
	}
}
