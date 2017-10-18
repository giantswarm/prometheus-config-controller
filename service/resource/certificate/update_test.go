package certificate

import (
	"context"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Resource_Certificate_GetUpdateState tests the GetUpdateState method.
func Test_Resource_Certificate_GetUpdateState(t *testing.T) {
	tests := []struct {
		currentState []certificateFile
		desiredState []certificateFile

		expectedUpdateStateCertificateFiles []certificateFile
		expectedErrorHandler                func(error) bool
	}{
		// Test that when the current state and desired state are both nil,
		// a nil update state is returned.
		{
			currentState: nil,
			desiredState: nil,

			expectedUpdateStateCertificateFiles: nil,
			expectedErrorHandler:                nil,
		},

		// Test that when the current state is nil, and desired state is empty,
		// a nil update state is returned.
		{
			currentState: nil,
			desiredState: []certificateFile{},

			expectedUpdateStateCertificateFiles: nil,
			expectedErrorHandler:                nil,
		},

		// Test that when the current state is empty, and desired state is nil,
		// a nil update state is returned.
		{
			currentState: []certificateFile{},
			desiredState: nil,

			expectedUpdateStateCertificateFiles: nil,
			expectedErrorHandler:                nil,
		},

		// Test that when the current and desired state are both empty,
		// a nil update state is returned.
		{
			currentState: []certificateFile{},
			desiredState: []certificateFile{},

			expectedUpdateStateCertificateFiles: nil,
			expectedErrorHandler:                nil,
		},

		// Test that when the current state is empty,
		// and the desired state contains a certificate file,
		// the update state contains the certificate file.
		{
			currentState: []certificateFile{},
			desiredState: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},

			expectedUpdateStateCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that when the current state contains a certificate,
		// and the desired state contains the same certificate,
		// the update state is nil.
		{
			currentState: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},
			desiredState: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},

			expectedUpdateStateCertificateFiles: nil,
			expectedErrorHandler:                nil,
		},

		// Test that when the current state contains a certificate,
		// and the desired state is empty,
		// the update state is empty.
		{
			currentState: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},
			desiredState: []certificateFile{},

			expectedUpdateStateCertificateFiles: []certificateFile{},
			expectedErrorHandler:                nil,
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.Fs = fs
		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.CertificatePermission = 0644

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		createState, deleteState, updateState, err := resource.GetUpdateState(
			context.TODO(), v1.Service{}, test.currentState, test.desiredState,
		)

		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned getting update state: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned getting update state: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned getting update state\n", index)
		}

		if createState != nil {
			t.Fatalf("%d: createState should be nil, returned: %#v\n", index, createState)
		}
		if deleteState != nil {
			t.Fatalf("%d: deleteState should be nil, returned: %#v\n", index, deleteState)
		}

		if updateState == nil && test.expectedUpdateStateCertificateFiles != nil {
			t.Fatalf("%d: updateState was nil, should be: %s\n", index, spew.Sdump(test.expectedUpdateStateCertificateFiles))
		}

		if updateState != nil {
			updateStateCertificateFiles, ok := updateState.(*[]certificateFile)
			if !ok {
				t.Fatalf("%d: could not cast update state to certificate files: %s\n", index, spew.Sdump(updateState))
			}

			if !reflect.DeepEqual(test.expectedUpdateStateCertificateFiles, *updateStateCertificateFiles) {
				t.Fatalf(
					"%d: expected update state does not match returned update state.\nexpected: %s\nreturned: %s\n",
					index,
					spew.Sdump(test.expectedUpdateStateCertificateFiles),
					spew.Sdump(*updateStateCertificateFiles),
				)
			}
		}
	}
}

// Test_Resource_Certificate_ProcessUpdateState tests the ProcessUpdateState method.
func Test_Resource_Certificate_ProcessUpdateState(t *testing.T) {
	tests := []struct {
		updateState []certificateFile

		expectedCertificateFiles []certificateFile
		expectedErrorHandler     func(error) bool
	}{
		// Test that when the updateState is nil, no certificates are written,
		// and no error is returned.
		{
			updateState: nil,

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that when the updateState is empty, no certificates are written,
		// and no error is returned.
		{
			updateState: []certificateFile{},

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that when the updateState contains one certificate,
		// one certificate is written, and no error is returned.
		{
			updateState: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that when the updateState contains two certificates,
		// two certificates are written, and no error is returned.
		{
			updateState: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
				{
					path: "/certs/bar",
					data: "bar",
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
				{
					path: "/certs/bar",
					data: "bar",
				},
			},
			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.Fs = fs
		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.CertificatePermission = 0644

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		updateErr := resource.ProcessUpdateState(context.TODO(), v1.Service{}, test.updateState)

		if updateErr != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned processing update state: %s\n", index, updateErr)
		}
		if updateErr != nil && !test.expectedErrorHandler(updateErr) {
			t.Fatalf("%d: incorrect error returned processing update state: %s\n", index, updateErr)
		}
		if updateErr == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned processing update state\n", index)
		}

		for _, expectedCertificateFile := range test.expectedCertificateFiles {
			data, err := afero.ReadFile(fs, expectedCertificateFile.path)
			if err != nil {
				t.Fatalf("%d: could not read expected certificate file: %s\n", index, err)
			}

			if string(data) != expectedCertificateFile.data {
				t.Fatalf(
					"%d: expected certificate does not match written certificate.\nexpected: %s\nreturned: %s\n",
					index,
					expectedCertificateFile.data,
					string(data),
				)
			}
		}
	}
}
