package certificate

import (
	"context"
	"path"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus/prometheustest"
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
		// and the desired state contains the same certificate, with different data,
		// the update state is the new certificate.
		{
			currentState: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
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
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertificateComponentName = "prometheus"
		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.CertificateNamespace = "default"
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
			updateStateCertificateFiles, err := toCertificateFiles(updateState)
			if err != nil {
				t.Fatalf("%d: could not cast update state to certificate files: %s\n", index, spew.Sdump(updateState))
			}

			if !reflect.DeepEqual(test.expectedUpdateStateCertificateFiles, updateStateCertificateFiles) {
				t.Fatalf(
					"%d: expected update state does not match returned update state.\nexpected: %s\nreturned: %s\n",
					index,
					spew.Sdump(test.expectedUpdateStateCertificateFiles),
					spew.Sdump(updateStateCertificateFiles),
				)
			}
		}
	}
}

// Test_Resource_Certificate_ProcessUpdateState tests the ProcessUpdateState method.
func Test_Resource_Certificate_ProcessUpdateState(t *testing.T) {
	tests := []struct {
		currentCertificateFiles []certificateFile
		updateState             []certificateFile

		expectedCertificateFiles []certificateFile
		expectedErrorHandler     func(error) bool
	}{
		// Test that when the updateState is nil and no certificates are on disk,
		// no certificates are written, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{},
			updateState:             nil,

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that when the updateState is empty and no certificates are on disk,
		//  no certificates are written, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{},
			updateState:             []certificateFile{},

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that when the updateState contains one certificate and no certificates are on disk,
		// one certificate is written, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{},
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

		// Test that when the updateState contains two certificates and no certificates are on disk,
		// two certificates are written, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{},
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

		// Test that when the updateState contains one certificate,
		// and the same certificate is on disk,
		// the certificate is not updated, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
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

		// Test that when the updateState contains no certificates,
		// and there is one certificate on disk,
		// the certificate is removed, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
			updateState: []certificateFile{},

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that when the updateState is nil,
		// and there is one certificate on disk,
		// the certificate is not removed, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
			updateState: nil,

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that when the updateState contains one certificate,
		// and there is one certificate on disk with different data,
		// the certificate is updated, and no error is returned.
		{
			currentCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
			updateState: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "bar",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that when the updateState contains two certificates,
		// and one of the certificates is on disk,
		// the other certificate is added,
		// and no error is returned.
		{
			currentCertificateFiles: []certificateFile{
				{
					path: "/certs/foo",
					data: "foo",
				},
			},
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
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertificateComponentName = "prometheus"
		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.CertificateNamespace = "default"
		resourceConfig.CertificatePermission = 0644

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		if err := fs.Mkdir(resourceConfig.CertificateDirectory, 0600); err != nil {
			t.Fatalf("%d: error returned creating certificate directory: %s\n", index, err)
		}

		for _, currentCertificateFile := range test.currentCertificateFiles {
			if err := afero.WriteFile(fs, currentCertificateFile.path, []byte(currentCertificateFile.data), 0600); err != nil {
				t.Fatalf("%d: error returned writing current certificate file: %s\n", index, err)
			}
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

		fileInfos, err := afero.ReadDir(fs, resourceConfig.CertificateDirectory)
		if err != nil {
			t.Fatalf("%d: error returned reading directory: %s\n", index, err)
		}

		if len(fileInfos) == 0 && len(test.expectedCertificateFiles) > 0 {
			t.Fatalf("%d: expected certificates not found: %#v\n", index, test.expectedCertificateFiles)
		}

		for _, fileInfo := range fileInfos {
			foundFile := false

			path := path.Join(resourceConfig.CertificateDirectory, fileInfo.Name())
			data, err := afero.ReadFile(fs, path)
			if err != nil {
				t.Fatalf("%d: could not read expected certificate file: %s\n", index, err)
			}

			for _, expectedCertificateFile := range test.expectedCertificateFiles {
				if path == expectedCertificateFile.path && string(data) == expectedCertificateFile.data {
					foundFile = true
				}
			}

			if !foundFile {
				t.Fatalf("%d: unexpected certificate found: %s, %s", index, path, string(data))
			}
		}
	}
}

// Test_Resource_Certificate_toCertificateFiles tests the Test_Resource_Certificate_toCertificateFiles function.
func Test_Resource_Certificate_toCertificateFiles(t *testing.T) {
	tests := []struct {
		v interface{}

		expectedCertificateFiles []certificateFile
		expectedErrorHandler     func(error) bool
	}{
		// Test that a nil interface returns nil.
		{
			v: nil,

			expectedCertificateFiles: nil,
			expectedErrorHandler:     nil,
		},

		// Test that a pointer to a slice of certificate files
		// returns an error.
		{
			v: &[]certificateFile{},

			expectedCertificateFiles: nil,
			expectedErrorHandler:     IsWrongTypeError,
		},

		// Test that a slice of certificate files
		// returns a slice of certificate files.
		{
			v: []certificateFile{},

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},
	}

	for index, test := range tests {
		returnedCertificateFiles, err := toCertificateFiles(test.v)

		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned processing update state: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned processing update state: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned processing update state\n", index)
		}

		if !reflect.DeepEqual(test.expectedCertificateFiles, returnedCertificateFiles) {
			t.Fatalf(
				"%d: expected certificate files do not match returned certificate files.\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedCertificateFiles),
				spew.Sdump(returnedCertificateFiles),
			)
		}
	}
}
