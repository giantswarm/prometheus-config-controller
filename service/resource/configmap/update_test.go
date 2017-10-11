package configmap

import (
	"context"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Resource_ConfigMap_GetUpdateState tests the GetUpdateState method.
func Test_Resource_ConfigMap_GetUpdateState(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		currentState *v1.ConfigMap
		desiredState *v1.ConfigMap

		expectedUpdateStateConfigMap *v1.ConfigMap
		expectedErrorHandler         func(error) bool
	}{
		// Test that if current state and desired state are both nil,
		// the update state is nil.
		{
			currentState: nil,
			desiredState: nil,

			expectedUpdateStateConfigMap: nil,
			expectedErrorHandler:         nil,
		},

		// Test that that if the current state and desired state are the same,
		// and the configmap is empty,
		// the update state is nil.
		{
			currentState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{},
			},
			desiredState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{},
			},

			expectedUpdateStateConfigMap: nil,
			expectedErrorHandler:         nil,
		},

		// Test that if the current state and desired state are the same,
		// the update state is nil.
		{
			currentState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			desiredState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},

			expectedUpdateStateConfigMap: nil,
			expectedErrorHandler:         nil,
		},

		// Test that if the current state and desired state are different,
		// the update state matches the desired state.
		{
			currentState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			desiredState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "bar",
				},
			},

			expectedUpdateStateConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "bar",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test that if the current state and desired state configmaps have
		// different names, an error is returned.
		{
			currentState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			desiredState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "another-name",
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},

			expectedUpdateStateConfigMap: nil,
			expectedErrorHandler:         IsConfigMapWrongName,
		},

		// Test that if the current state and desired state configmaps are in
		// different namespaces, an error is returned.
		{
			currentState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			desiredState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: "another-namespace",
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},

			expectedUpdateStateConfigMap: nil,
			expectedErrorHandler:         IsConfigMapWrongNamespace,
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()

		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		createState, updateState, deleteState, err := resource.GetUpdateState(
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

		// We don't create or delete the configmap,
		// so create and delete state should always be nil.
		if createState != nil {
			t.Fatalf("%d: createState should be nil, returned: %#v\n", index, createState)
		}
		if deleteState != nil {
			t.Fatalf("%d: deleteState should be nil, returned: %#v\n", index, deleteState)
		}

		if updateState == nil && test.expectedUpdateStateConfigMap != nil {
			t.Fatalf("%d: updateState was nil, should be: %s\n", index, spew.Sdump(test.expectedUpdateStateConfigMap))
		}

		if updateState != nil {
			updateStateConfigMap, ok := updateState.(*v1.ConfigMap)
			if !ok {
				t.Fatalf("%d: could not cast update state to configmap: %s\n", index, spew.Sdump(updateState))
			}

			if !reflect.DeepEqual(*updateStateConfigMap, *test.expectedUpdateStateConfigMap) {
				t.Fatalf(
					"%d: expected update state does not match returned update state.\nexpected: %s\nreturned: %s\n",
					index,
					spew.Sdump(test.expectedUpdateStateConfigMap),
					spew.Sdump(updateStateConfigMap),
				)
			}
		}
	}
}
