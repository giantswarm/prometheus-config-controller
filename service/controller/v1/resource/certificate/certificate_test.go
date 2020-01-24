package certificate

import (
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes/fake"
)

// Test_Resource_Certificate_New tests the New function.
func Test_Resource_Certificate_New(t *testing.T) {
	tests := []struct {
		config func() Config

		expectedErrorHandler func(error) bool
	}{
		// Test that the default config returns an error.
		{
			config: func() Config { return Config{} },

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the fs must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        nil,
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the kubernetes clientset must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: nil,
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the logger must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    nil,

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the certificate component name must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the certificate directory must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the certificate namespace must not be empty.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the certificate permission must not be zero.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0,
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that a valid config produces a configmap resource.
		{
			config: func() Config {
				return Config{
					Fs:        afero.NewMemMapFs(),
					K8sClient: fake.NewSimpleClientset(),
					Logger:    microloggertest.New(),

					CertComponentName: "prometheus",
					CertDirectory:     "/certs",
					CertNamespace:     "default",
					CertPermission:    0600,
				}
			},

			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		config := test.config()

		service, err := New(config)
		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned creating certificate resource: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned creating certificate resource: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned creating certificate resource\n", index)
		}

		if test.expectedErrorHandler == nil && service == nil {
			t.Fatalf("%d: returned certificate resource was nil", index)
		}
	}
}
