package configmap

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/giantswarm/micrologger/microloggertest"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus"
	"github.com/giantswarm/prometheus-config-controller/service/controller/v1/prometheus/prometheustest"
)

// Test_Resource_ConfigMap_newUpdateChange tests the newUpdateChange method.
func Test_Resource_ConfigMap_newUpdateChange(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		currentState *v1.ConfigMap
		desiredState *v1.ConfigMap

		expectedUpdateChangeConfigMap *v1.ConfigMap
		expectedErrorHandler          func(error) bool
	}{
		// Test that if current state and desired state are both nil,
		// the update change is nil.
		{
			currentState: nil,
			desiredState: nil,

			expectedUpdateChangeConfigMap: nil,
			expectedErrorHandler:          nil,
		},

		// Test that if the current state and desired state are the same,
		// and the configmap is empty,
		// the update change is nil.
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

			expectedUpdateChangeConfigMap: nil,
			expectedErrorHandler:          nil,
		},

		// Test that if the current state and desired state are the same,
		// the update change is nil.
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

			expectedUpdateChangeConfigMap: nil,
			expectedErrorHandler:          nil,
		},

		// Test that if the current state and desired state are different,
		// the update change matches the desired state.
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

			expectedUpdateChangeConfigMap: &v1.ConfigMap{
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

			expectedUpdateChangeConfigMap: nil,
			expectedErrorHandler:          IsWrongName,
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

			expectedUpdateChangeConfigMap: nil,
			expectedErrorHandler:          IsWrongNamespace,
		},
	}

	for index, test := range tests {
		fakeK8sClient := fake.NewSimpleClientset()

		resourceConfig := Config{}

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertDirectory = "/certs"
		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		updateChange, err := resource.newUpdateChange(
			context.TODO(), v1.Service{}, test.currentState, test.desiredState,
		)

		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned getting update change: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned getting update change: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned getting update change\n", index)
		}

		if updateChange == nil && test.expectedUpdateChangeConfigMap != nil {
			t.Fatalf("%d: updateChange was nil, should be: %s\n", index, spew.Sdump(test.expectedUpdateChangeConfigMap))
		}

		if updateChange != nil {
			updateChangeConfigMap, ok := updateChange.(*v1.ConfigMap)
			if !ok {
				t.Fatalf("%d: could not cast update state to configmap: %s\n", index, spew.Sdump(updateChange))
			}

			if !reflect.DeepEqual(*updateChangeConfigMap, *test.expectedUpdateChangeConfigMap) {
				t.Fatalf(
					"%d: expected update change does not match returned update change.\nexpected: %s\nreturned: %s\n",
					index,
					spew.Sdump(test.expectedUpdateChangeConfigMap),
					spew.Sdump(updateChangeConfigMap),
				)
			}
		}
	}
}

// Test_Resource_ConfigMap_ApplyUpdateChange tests the ApplyUpdateChange method.
func Test_Resource_ConfigMap_ApplyUpdateChange(t *testing.T) {
	configMapKey := "prometheus.yml"
	configMapName := "prometheus"
	configMapNamespace := "monitoring"

	tests := []struct {
		setUpConfigMap *v1.ConfigMap
		updateChange   *v1.ConfigMap

		expectedConfigMap    *v1.ConfigMap
		expectedErrorHandler func(error) bool
	}{
		// Test if the initial configmap is nil, and the update change is some
		// configmap, an error occurs.
		{
			setUpConfigMap: nil,
			updateChange: &v1.ConfigMap{
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

		// Test if the initial configmap exists, and the update change is nil,
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
			updateChange: nil,

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

		// Test if the initial configmap exists, and the update change is the same,
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
			updateChange: &v1.ConfigMap{
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

		// Test if the initial configmap exists, and a different update change exists,
		// the expected configmap matches the update change.
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
			updateChange: &v1.ConfigMap{
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

		resourceConfig := Config{}

		resourceConfig.K8sClient = fakeK8sClient
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.PrometheusReloader = prometheustest.New()

		resourceConfig.CertDirectory = "/certs"
		resourceConfig.ConfigMapKey = configMapKey
		resourceConfig.ConfigMapName = configMapName
		resourceConfig.ConfigMapNamespace = configMapNamespace

		resource, err := New(resourceConfig)
		if err != nil {
			t.Fatalf("%d: error returned creating resource: %s\n", index, err)
		}

		if test.setUpConfigMap != nil {
			if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(test.setUpConfigMap); err != nil {
				t.Fatalf("%d: error returned setting up configmap: %s\n", index, err)
			}
		}

		updateErr := resource.ApplyUpdateChange(context.TODO(), v1.Service{}, test.updateChange)

		if updateErr != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned applying update change: %s\n", index, updateErr)
		}
		if updateErr != nil && !test.expectedErrorHandler(updateErr) {
			t.Fatalf("%d: incorrect error returned applying update change: %s\n", index, updateErr)
		}
		if updateErr == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned applying update change\n", index)
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

	fakeK8sClient := fake.NewSimpleClientset()

	if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "bar",
		},
	}); err != nil {
		t.Fatalf("error returned creating configmap: %s\n", err)
	}

	var receivedReloadMessage *http.Request = nil
	configRequestCount := 0
	reloadRequestCount := 0

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == prometheus.ConfigPath {
			configRequestCount++

			io.WriteString(w, "{ \"status\": \"success\", \"data\": { \"yaml\": \"foo\" }}")
			return
		}
		if r.URL.Path == prometheus.ReloadPath {
			receivedReloadMessage = r
			reloadRequestCount++

			io.WriteString(w, "")
			return
		}
	}))
	defer testServer.Close()

	prometheusConfig := prometheus.DefaultConfig()

	prometheusConfig.K8sClient = fakeK8sClient
	prometheusConfig.Logger = microloggertest.New()

	prometheusConfig.Address = testServer.URL
	prometheusConfig.ConfigMapKey = configMapKey
	prometheusConfig.ConfigMapName = configMapName
	prometheusConfig.ConfigMapNamespace = configMapNamespace
	prometheusConfig.MinimumReloadTime = 10 * time.Millisecond

	prometheusReloader, err := prometheus.New(prometheusConfig)
	if err != nil {
		t.Fatalf("error returned creating prometheus service: %s\n", err)
	}

	resourceConfig := Config{}

	resourceConfig.K8sClient = fakeK8sClient
	resourceConfig.Logger = microloggertest.New()
	resourceConfig.PrometheusReloader = prometheusReloader

	resourceConfig.CertDirectory = "/certs"
	resourceConfig.ConfigMapKey = configMapKey
	resourceConfig.ConfigMapName = configMapName
	resourceConfig.ConfigMapNamespace = configMapNamespace

	resource, err := New(resourceConfig)
	if err != nil {
		t.Fatalf("error returned creating resource: %s\n", err)
	}

	if configRequestCount != 0 {
		t.Fatalf("incorrect config request count before update - should be 0, was: %d", configRequestCount)
	}
	if reloadRequestCount != 0 {
		t.Fatalf("incorrect reload request count before update - should be 0, was: %d", reloadRequestCount)
	}

	if err := resource.ApplyUpdateChange(context.TODO(), v1.Service{}, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "bar",
		},
	}); err != nil {
		t.Fatalf("error returned applying update change: %s\n", err)
	}

	if receivedReloadMessage == nil {
		t.Fatalf("handler did not receive reload message")
	}

	if receivedReloadMessage.Method != "POST" {
		t.Fatalf("incorrect method used for reload: %s\n", receivedReloadMessage.Method)
	}

	if receivedReloadMessage.URL.Path != "/-/reload" {
		t.Fatalf("incorrect path used for reload: %s\n", receivedReloadMessage.URL.Path)
	}

	if configRequestCount != 1 {
		t.Fatalf("incorrect config request count after update - should be 1, was: %d", configRequestCount)
	}
	if reloadRequestCount != 1 {
		t.Fatalf("incorrect reload request count after update - should be 1, was: %d", reloadRequestCount)
	}

	// Update the configmap to match the prometheus config.
	if _, err := fakeK8sClient.CoreV1().ConfigMaps(configMapNamespace).Update(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
		Data: map[string]string{
			configMapKey: "foo",
		},
	}); err != nil {
		t.Fatalf("error returned updating configmap: %s\n", err)
	}

	// Wait out the rate limit
	time.Sleep(20 * time.Millisecond)

	// Check that a nil processing does not cause a reload.
	if err := resource.ApplyUpdateChange(context.TODO(), v1.Service{}, nil); err != nil {
		t.Fatalf("error returned applying update change: %s\n", err)
	}

	if configRequestCount != 2 {
		t.Fatalf("incorrect config request count after nil update - should be 2, was: %d", configRequestCount)
	}
	if reloadRequestCount != 1 {
		t.Fatalf("incorrect reload request count after nil update - should be 1, was: %d", reloadRequestCount)
	}
}
