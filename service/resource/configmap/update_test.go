package configmap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/prometheus-config-controller/service/prometheus"
	"github.com/giantswarm/prometheus-config-controller/service/prometheus/prometheustest"
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

		// Test that if the current state and desired state are the same,
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
			expectedErrorHandler:         IsWrongName,
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
			expectedErrorHandler:         IsWrongNamespace,
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace
		resourceConfig.ReloadWaitTime = 1 * time.Second

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

// Test_Resource_ConfigMap_ProcessUpdateState tests the ProcessUpdateState method.
func Test_Resource_ConfigMap_ProcessUpdateState(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		setUpConfigMap *v1.ConfigMap
		updateState    *v1.ConfigMap

		expectedConfigMap    *v1.ConfigMap
		expectedErrorHandler func(error) bool
	}{
		// Test if the initial configmap is nil, and the update state is some
		// configmap, an error occurs.
		{
			setUpConfigMap: nil,
			updateState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},

			expectedConfigMap:    nil,
			expectedErrorHandler: IsConfigMapNotFound,
		},

		// Test if the initial configmap exists, and the update state is nil,
		// the expected configmap matches the initial configmap.
		{
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			updateState: nil,

			expectedConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test if the initial configmap exists, and the update state is the same,
		// the expected configmap matches the initial configmap.
		{
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			updateState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},

			expectedConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			expectedErrorHandler: nil,
		},

		// Test if the initial configmap exists, and a different update state exists,
		// the expected configmap matches the update state.
		{
			setUpConfigMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "foo",
				},
			},
			updateState: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
					Namespace: configMapNamespace,
				},
				Data: map[string]string{
					configMapKey: "bar",
				},
			},

			expectedConfigMap: &v1.ConfigMap{
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
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertificateDirectory = "/certs"
		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace
		resourceConfig.ReloadWaitTime = 1 * time.Second

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		if test.setUpConfigMap != nil {
			if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(test.setUpConfigMap); err != nil {
				t.Fatalf("%d: error returned setting up configmap: %s\n", index, err)
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

		if test.expectedConfigMap == nil {
			_, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Get(
				configMapName, metav1.GetOptions{},
			)
			if !errors.IsNotFound(err) {
				t.Fatalf("%d: unexpectedly found configmap", index)
			}
		}

		if test.expectedConfigMap != nil {
			configMap, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Get(
				configMapName, metav1.GetOptions{},
			)
			if err != nil {
				t.Fatalf("%d: error returned getting configmap: %s\n", index, err)
			}

			if !reflect.DeepEqual(*test.expectedConfigMap, *configMap) {
				t.Fatalf(
					"%d: expected configmap does not match returned desired state.\nexpected:\n%s\nreturned:\n%s\n",
					index,
					spew.Sdump(*test.expectedConfigMap),
					spew.Sdump(*configMap),
				)
			}
		}
	}
}

// Test_Resource_ConfigMap_Reload tests that the configmap is reloaded correctly.
func Test_Resource_ConfigMap_Reload(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	var receivedMessage *http.Request = nil

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMessage = r
	}))
	defer testServer.Close()

	prometheusConfig := prometheus.DefaultConfig()
	prometheusConfig.Logger = microloggertest.New()
	prometheusConfig.Address = testServer.URL

	prometheusReloader, err := prometheus.New(prometheusConfig)
	if err != nil {
		t.Fatalf("error returned creating prometheus service: %s\n", err)
	}

	fakeK8sClient := fake.NewSimpleClientset()

	if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "foo",
		},
	}); err != nil {
		t.Fatalf("error returned creating configmap: %s\n", err)
	}

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheusReloader

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = configMapKey
	resourceConfig.ConfigMapName = configMapName
	resourceConfig.ConfigMapNamespace = configMapNamespace
	resourceConfig.ReloadWaitTime = 1 * time.Second

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if err := resource.ProcessUpdateState(context.TODO(), v1.Service{}, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "bar",
		},
	}); err != nil {
		t.Fatalf("error returned processing update state: %s\n", err)
	}

	if receivedMessage == nil {
		t.Fatalf("handler did not receive message")
	}

	if receivedMessage.Method != "POST" {
		t.Fatalf("incorrect method used: %s\n", receivedMessage.Method)
	}

	if receivedMessage.URL.Path != "/-/reload" {
		t.Fatalf("incorrect path used: %s\n", receivedMessage.URL.Path)
	}
}

// Test_Resource_ConfigMap_Reload_Count tests that the configmap is reloaded the correct number of times.
func Test_Resource_ConfigMap_Reload_Count(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	var reloadCount int = 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reloadCount++
	}))
	defer testServer.Close()

	prometheusConfig := prometheus.DefaultConfig()
	prometheusConfig.Logger = microloggertest.New()
	prometheusConfig.Address = testServer.URL

	prometheusReloader, err := prometheus.New(prometheusConfig)
	if err != nil {
		t.Fatalf("error returned creating prometheus service: %s\n", err)
	}

	fakeK8sClient := fake.NewSimpleClientset()

	if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "foo",
		},
	}); err != nil {
		t.Fatalf("error returned creating configmap: %s\n", err)
	}

	resourceConfig := DefaultConfig()

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheusReloader

	resourceConfig.CertificateDirectory = "/certs"
	resourceConfig.ConfigMapKey = configMapKey
	resourceConfig.ConfigMapName = configMapName
	resourceConfig.ConfigMapNamespace = configMapNamespace
	resourceConfig.ReloadWaitTime = 1 * time.Second

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if reloadCount != 0 {
		t.Fatalf("incorrect reload count - should be 0, was: %d", reloadCount)
	}

	if err := resource.ProcessUpdateState(context.TODO(), v1.Service{}, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "bar",
		},
	}); err != nil {
		t.Fatalf("error returned processing update state: %s\n", err)
	}

	if reloadCount != 1 {
		t.Fatalf("incorrect reload count - should be 1, was: %d", reloadCount)
	}

	// Check that a nil processing does not cause a reload.
	if err := resource.ProcessUpdateState(context.TODO(), v1.Service{}, nil); err != nil {
		t.Fatalf("error returned processing update state: %s\n", err)
	}

	if reloadCount != 1 {
		t.Fatalf("incorrect reload count - should be 1, was: %d", reloadCount)
	}
}
