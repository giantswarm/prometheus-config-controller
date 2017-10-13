package certificate

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Resource_Certificate_GetCurrentState tests the GetCurrentState method.
func Test_Resource_Certificate_GetCurrentState(t *testing.T) {
	var fileMode os.FileMode = 0644

	tests := []struct {
		setUp                func(afero.Fs) error
		certificateDirectory string

		expectedCertificateFiles []certificateFile
	}{
		// Test when the filesystem is empty, there are no certificates returned.
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			certificateDirectory: "/certs",

			expectedCertificateFiles: []certificateFile{},
		},

		// Test when one certificate exists on the filesystem,
		// one certificate is returned.
		{
			setUp: func(fs afero.Fs) error {
				return afero.WriteFile(fs, "/certs/kf83j-ca.pem", []byte("foo"), fileMode)
			},
			certificateDirectory: "/certs",

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/kf83j-ca.pem",
					data: "foo",
				},
			},
		},

		// Test when two certificates exist on the filesystem,
		// two certificates are returned.
		{
			setUp: func(fs afero.Fs) error {
				if err := afero.WriteFile(fs, "/certs/kf83j-ca.pem", []byte("foo"), fileMode); err != nil {
					return microerror.Mask(err)
				}
				if err := afero.WriteFile(fs, "/certs/kf83j-crt.pem", []byte("bar"), fileMode); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			certificateDirectory: "/certs",

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/kf83j-ca.pem",
					data: "foo",
				},
				{
					path: "/certs/kf83j-crt.pem",
					data: "bar",
				},
			},
		},

		// Test that files in another directory are not returned.
		{
			setUp: func(fs afero.Fs) error {
				return afero.WriteFile(fs, "/somewhere-else/jf92j-ca.pem", []byte("foo"), fileMode)
			},
			certificateDirectory: "/certs",

			expectedCertificateFiles: []certificateFile{},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()

		resourceConfig := DefaultConfig()

		resourceConfig.Fs = fs
		resourceConfig.Logger = microloggertest.New()

		resourceConfig.CertificateDirectory = test.certificateDirectory

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		if err := fs.Mkdir(test.certificateDirectory, fileMode); err != nil {
			t.Fatalf("%d: error returned creating certificate directory: %s\n", index, err)
		}

		if err := test.setUp(fs); err != nil {
			t.Fatalf("%d: error returned during setup: %s\n", index, err)
		}

		currentState, err := resource.GetCurrentState(context.TODO(), v1.Service{})
		if err != nil {
			t.Fatalf("%d: error returned getting current state: %v\n", index, err)
		}

		if !(test.expectedCertificateFiles == nil && currentState == nil) &&
			!reflect.DeepEqual(test.expectedCertificateFiles, currentState) {
			t.Fatalf(
				"%d: expected configmap does not match returned current state.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedCertificateFiles),
				spew.Sdump(currentState),
			)
		}
	}
}
