package prometheus

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Test_Prometheus_FilterInvalidServices tests the FilterInvalidServices function.
func Test_Prometheus_FilterInvalidServices(t *testing.T) {
	tests := []struct {
		services         []v1.Service
		expectedServices []v1.Service
	}{
		// Test a service without either cluster or certificate annotations is filtered.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
				},
			},
			expectedServices: []v1.Service{},
		},

		// Test a service with just a cluster annotation is filtered.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			expectedServices: []v1.Service{},
		},

		// Test a service with just a certificate annotation is filtered.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
						Annotations: map[string]string{
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
			},
			expectedServices: []v1.Service{},
		},

		// Test that a service with both cluster and certificate annotations
		// is not filtered.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
						Annotations: map[string]string{
							ClusterAnnotation:     "xa5ly",
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
			},
			expectedServices: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
						Annotations: map[string]string{
							ClusterAnnotation:     "xa5ly",
							CertificateAnnotation: "default/xa5ly-prometheus",
						},
					},
				},
			},
		},

		// Test that two services without either annotations are filtered.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
				},

				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "bar",
					},
				},
			},
			expectedServices: []v1.Service{},
		},
	}

	for index, test := range tests {
		filteredServices := FilterInvalidServices(test.services)

		if !reflect.DeepEqual(test.expectedServices, filteredServices) {
			t.Fatalf(
				"%d: expected filtered services do not match returned filtered services\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedServices),
				spew.Sdump(filteredServices),
			)
		}
	}
}

// Test_Prometheus_GroupServices tests the GroupServices function.
func Test_Prometheus_GroupServices(t *testing.T) {
	tests := []struct {
		services                []v1.Service
		expectedGroupedServices map[string][]v1.Service
	}{
		// Test that an empty services list leads to an empty grouped services map.
		{
			services:                []v1.Service{},
			expectedGroupedServices: map[string][]v1.Service{},
		},

		// Test that a service that does not specify the cluster annotation
		// is dropped.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foo",
					},
				},
			},
			expectedGroupedServices: map[string][]v1.Service{},
		},

		// Test that a single service, specifying the cluster annotation,
		// is grouped on its own.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foo",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			expectedGroupedServices: map[string][]v1.Service{
				"xa5ly": []v1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
							Annotations: map[string]string{
								ClusterAnnotation: "xa5ly",
							},
						},
					},
				},
			},
		},

		// Test that two services, specifying the same cluster annotation,
		// are grouped together.
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foo",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
			},
			expectedGroupedServices: map[string][]v1.Service{
				"xa5ly": []v1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
							Annotations: map[string]string{
								ClusterAnnotation: "xa5ly",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
							Annotations: map[string]string{
								ClusterAnnotation: "xa5ly",
							},
						},
					},
				},
			},
		},

		// Test that two services, specifying the different cluster annotation,
		// are grouped separately
		{
			services: []v1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foo",
						Annotations: map[string]string{
							ClusterAnnotation: "xa5ly",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar",
						Annotations: map[string]string{
							ClusterAnnotation: "jeo0d",
						},
					},
				},
			},
			expectedGroupedServices: map[string][]v1.Service{
				"xa5ly": []v1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo",
							Annotations: map[string]string{
								ClusterAnnotation: "xa5ly",
							},
						},
					},
				},
				"jeo0d": []v1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
							Annotations: map[string]string{
								ClusterAnnotation: "jeo0d",
							},
						},
					},
				},
			},
		},
	}

	for index, test := range tests {
		groupedServices := GroupServices(test.services)

		if !reflect.DeepEqual(test.expectedGroupedServices, groupedServices) {
			t.Fatalf(
				"%d: expected grouped services do not match returned grouped services\nexpected: %s\nreturned: %s\n",
				index,
				spew.Sdump(test.expectedGroupedServices),
				spew.Sdump(groupedServices),
			)
		}
	}
}
