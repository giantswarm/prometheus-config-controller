package certificate

import (
	"context"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/prometheus-config-controller/service/key"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
)

// Test_Resource_Certificate_GetDesiredState tests the GetDesiredState method.
func Test_Resource_Certificate_GetDesiredState(t *testing.T) {
	defaultCertificateDirectory := "/certs"

	tests := []struct {
		certificateDirectory string
		services             []*v1.Service
		secrets              []*v1.Secret

		expectedCertificateFiles []certificateFile
		expectedErrorHandler     func(error) bool
	}{
		// Test that an empty list of services leads to an empty list of certificate files.
		{
			certificateDirectory: defaultCertificateDirectory,
			services:             nil,
			secrets:              nil,

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that a service without any annotation does not create a certificate file.
		{
			certificateDirectory: defaultCertificateDirectory,
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
				},
			},
			secrets: nil,

			expectedCertificateFiles: []certificateFile{},
			expectedErrorHandler:     nil,
		},

		// Test that a service with a cluster annotation,
		// but the certificate is missing, produces an error.
		{
			certificateDirectory: defaultCertificateDirectory,
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			secrets: nil,

			expectedCertificateFiles: nil,
			expectedErrorHandler:     IsMissing,
		},

		// Test that a service with a cluster annotation,
		// and the certificate (containing just a ca) being present, returns the certificate.
		{
			certificateDirectory: defaultCertificateDirectory,
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "xa5ly-prometheus",
						Namespace: "default",
						Labels: map[string]string{
							"clusterComponent": "prometheus",
							"clusterID":        "xa5ly",
						},
					},
					Data: map[string][]byte{
						"ca": []byte("foo"),
					},
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					// Note: we use the key method here, as opposed to a hard-coded string,
					// to make the test more specific - we care that the path matches `CAPath` string,
					// not the exact string.
					path: key.CAPath(defaultCertificateDirectory, "xa5ly"),
					data: "foo",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that a service with a cluster annotation,
		// and a certificate with ca, crt, and key fields, returns three certificates.
		{
			certificateDirectory: defaultCertificateDirectory,
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "xa5ly-prometheus",
						Namespace: "default",
						Labels: map[string]string{
							"clusterComponent": "prometheus",
							"clusterID":        "xa5ly",
						},
					},
					Data: map[string][]byte{
						"ca":  []byte("foo"),
						"crt": []byte("bar"),
						"key": []byte("baz"),
					},
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: key.CAPath(defaultCertificateDirectory, "xa5ly"),
					data: "foo",
				},
				{
					path: key.CrtPath(defaultCertificateDirectory, "xa5ly"),
					data: "bar",
				},
				{
					path: key.KeyPath(defaultCertificateDirectory, "xa5ly"),
					data: "baz",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that two services, both with cluster annotations,
		// and certificates that have only ca field, return two certificates.
		{
			certificateDirectory: defaultCertificateDirectory,
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "al9qy",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "al9qy",
						},
					},
				},
			},
			secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "xa5ly-prometheus",
						Namespace: "default",
						Labels: map[string]string{
							"clusterComponent": "prometheus",
							"clusterID":        "xa5ly",
						},
					},
					Data: map[string][]byte{
						"ca": []byte("foo"),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "al9qy-prometheus",
						Namespace: "default",
						Labels: map[string]string{
							"clusterComponent": "prometheus",
							"clusterID":        "al9qy",
						},
					},
					Data: map[string][]byte{
						"ca": []byte("bar"),
					},
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: key.CAPath(defaultCertificateDirectory, "xa5ly"),
					data: "foo",
				},
				{
					path: key.CAPath(defaultCertificateDirectory, "al9qy"),
					data: "bar",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that a service with a cluster annotation,
		// and a certificate with a ca field, returns one certificate,
		// with the correct certificate directory.
		{
			certificateDirectory: "/foo/bar",
			services: []*v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "apiserver",
						Namespace: "xa5ly",
						Annotations: map[string]string{
							prometheus.ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "xa5ly-prometheus",
						Namespace: "default",
						Labels: map[string]string{
							"clusterComponent": "prometheus",
							"clusterID":        "xa5ly",
						},
					},
					Data: map[string][]byte{
						"ca": []byte("foo"),
					},
				},
			},

			expectedCertificateFiles: []certificateFile{
				{
					path: key.CAPath("/foo/bar", "xa5ly"),
					data: "foo",
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

		resourceConfig.CertificateComponentName = "prometheus"
		resourceConfig.CertificateDirectory = test.certificateDirectory
		resourceConfig.CertificateNamespace = "default"

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		for _, service := range test.services {
			if _, err := fakeK8sClient.CoreV1().Services(service.Namespace).Create(service); err != nil {
				t.Fatalf("%d: error returned creating service: %s\n", index, err)
			}
		}

		for _, secret := range test.secrets {
			if _, err := fakeK8sClient.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
				t.Fatalf("%d: error returned creating secret: %s\n", index, err)
			}
		}

		desiredState, err := resource.GetDesiredState(context.TODO(), v1.Service{})

		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned getting desired state: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned getting desired state: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned getting desired state\n", index)
		}

		if test.expectedCertificateFiles != nil {
			if !reflect.DeepEqual(test.expectedCertificateFiles, desiredState) {
				t.Fatalf(
					"%d: expected certificate files do not match returned certificate files.\nexpected:\n%s\nreturned:\n%s\n",
					index,
					spew.Sdump(test.expectedCertificateFiles),
					spew.Sdump(desiredState),
				)
			}
		}
	}
}
