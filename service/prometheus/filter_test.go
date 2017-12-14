package prometheus

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test_Prometheus_FilterInvalidServices tests the FilterInvalidServices function.
func Test_Prometheus_FilterInvalidServices(t *testing.T) {
	tests := []struct {
		services         []v1.Service
		expectedServices []v1.Service
	}{
		// Test a service without a cluster annotation is filtered.
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

		// Test that two services without cluster annotations are filtered.
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

		// Test that a service with a cluster annotation is not filtered.
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
			expectedServices: []v1.Service{
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
