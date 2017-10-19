package certificate

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Resource_Certificate_GetCreateState tests the GetCreateState method.
func Test_Resource_Certificate_GetCreateState(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateComponentName = "prometheus"
	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.CertificateNamespace = "default"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	createState, err := resource.GetCreateState(context.TODO(), v1.Service{}, []certificateFile{}, []certificateFile{})
	if err != nil {
		t.Fatalf("error returned getting create state: %s\n", err)
	}

	if createState != nil {
		t.Fatalf("create state should be nil, was: %#v", createState)
	}
}

// Test_Resource_Certificate_ProcessCreateState tests the ProcessCreateState method.
func Test_Resource_Certificate_ProcessCreateState(t *testing.T) {
	fs := afero.NewMemMapFs()
	fakeK8sClient := fake.NewSimpleClientset()

	resourceConfig := DefaultConfig()

	resourceConfig.Fs = fs
	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateComponentName = "prometheus"
	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.CertificateNamespace = "default"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ProcessCreateState(context.TODO(), v1.Service{}, []certificateFile{}); err != nil {
		t.Fatalf("error returned processing create state: %s\n", err)
	}
}
