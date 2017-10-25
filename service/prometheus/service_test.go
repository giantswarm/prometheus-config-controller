package prometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
)

// Test_Prometheus_New tests the New function.
func Test_Prometheus_New(t *testing.T) {
	tests := []struct {
		config func() Config

		expectedErrorHandler func(error) bool
	}{
		// Test that the default config returns an error.
		{
			config: DefaultConfig,

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that a logger must not be empty.
		{
			config: func() Config {
				return Config{
					Logger:  nil,
					Address: "http://127.0.0.1:8080",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the prometheus address must not be empty.
		{
			config: func() Config {
				return Config{
					Logger:  microloggertest.New(),
					Address: "",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that the prometheus address must be valid.
		{
			config: func() Config {
				return Config{
					Logger:  microloggertest.New(),
					Address: "jabberwocky",
				}
			},

			expectedErrorHandler: IsInvalidConfig,
		},

		// Test that a valid config produces a service.
		{
			config: func() Config {
				return Config{
					Logger:  microloggertest.New(),
					Address: "http://127.0.0.1:8080",
				}
			},

			expectedErrorHandler: nil,
		},
	}

	for index, test := range tests {
		config := test.config()

		service, err := New(config)
		if err != nil && test.expectedErrorHandler == nil {
			t.Fatalf("%d: unexpected error returned creating service: %s\n", index, err)
		}
		if err != nil && !test.expectedErrorHandler(err) {
			t.Fatalf("%d: incorrect error returned creating service: %s\n", index, err)
		}
		if err == nil && test.expectedErrorHandler != nil {
			t.Fatalf("%d: expected error not returned creating service\n", index)
		}

		if test.expectedErrorHandler == nil && service == nil {
			t.Fatalf("%d: returned service was nil", index)
		}
	}
}

// Test_Prometheus_Reload tests the Reload method.
func Test_Prometheus_Reload(t *testing.T) {
	var receivedMessage *http.Request = nil

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMessage = r
	}))
	defer testServer.Close()

	defaultConfig := DefaultConfig()

	defaultConfig.Logger = microloggertest.New()

	defaultConfig.Address = testServer.URL

	service, err := New(defaultConfig)
	if err != nil {
		t.Fatalf("error returned creating service: %s\n", err)
	}

	if err := service.Reload(); err != nil {
		t.Fatalf("an error returned reloading prometheus: %s\n", err)
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

// Test_Prometheus_Reload_Failure tests the Reload method if the reload fails.
// See https://github.com/prometheus/prometheus/blob/099df0c5f00c45c007a9779a2e4ab51cf4d076bf/web/web.go#L598
// Prometheus returns http.StatusInternalServerError on error.
func Test_Prometheus_Reload_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, fmt.Sprintf("beep boop i failed! D:"), http.StatusInternalServerError)
	}))
	defer testServer.Close()

	defaultConfig := DefaultConfig()

	defaultConfig.Logger = microloggertest.New()

	defaultConfig.Address = testServer.URL

	service, err := New(defaultConfig)
	if err != nil {
		t.Fatalf("error returned creating service: %s\n", err)
	}

	reloadErr := service.Reload()

	if reloadErr == nil {
		t.Fatalf("a nil error was returned\n")
	}
	if !IsReloadError(reloadErr) {
		t.Fatalf("incorrect error returned: %s\n", reloadError)
	}
}
