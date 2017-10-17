package certificate

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Resource_Certificate_GetDeleteState tests the GetDeleteState method.
func Test_Resource_Certificate_GetDeleteState(t *testing.T) {
	fs := afero.NewMemMapFs()

	resourceConfig := DefaultConfig()

	resourceConfig.Fs = fs
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateDirectory = "/certs"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	deleteState, err := resource.GetDeleteState(context.TODO(), v1.Service{}, []certificateFile{}, []certificateFile{})
	if err != nil {
		t.Fatalf("error returned getting delete state: %s\n", err)
	}

	if deleteState != nil {
		t.Fatalf("delete state should be nil, was: %#v", deleteState)
	}
}

// Test_Resource_Certificate_ProcessDeleteState tests the ProcessDeleteState method.
func Test_Resource_Certificate_ProcessDeleteState(t *testing.T) {
	fs := afero.NewMemMapFs()

	resourceConfig := DefaultConfig()

	resourceConfig.Fs = fs
	resourceConfig.Logger = microloggertest.New()

	resourceConfig.CertificateDirectory = "/certs"

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ProcessDeleteState(context.TODO(), v1.Service{}, []certificateFile{}); err != nil {
		t.Fatalf("error returned processing delete state: %s\n", err)
	}
}
