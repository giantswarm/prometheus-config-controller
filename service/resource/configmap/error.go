package configmap

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var configMapNotFoundError = microerror.New("configmap not found")

// IsConfigMapNotFound asserts configMapNotFoundError.
func IsConfigMapNotFound(err error) bool {
	return microerror.Cause(err) == configMapNotFoundError
}

var configMapKeyNotFoundError = microerror.New("configmap key not found")

// IsConfigMapKeyNotFound asserts configMapKeyNotFoundError.
func IsConfigMapKeyNotFound(err error) bool {
	return microerror.Cause(err) == configMapKeyNotFoundError
}

var invalidConfigMapError = microerror.New("invalid config map")

// IsInvalidConfigMap asserts invalidConfigMapError.
func IsInvalidConfigMap(err error) bool {
	return microerror.Cause(err) == invalidConfigMapError
}

var configMapAssertionError = microerror.New("configmap assertion")

// IsConfigMapAssertion asserts configMapAssertionError.
func IsConfigMapAssertion(err error) bool {
	return microerror.Cause(err) == configMapAssertionError
}

var configMapWrongNameError = microerror.New("configmap wrong name")

// IsConfigMapWrongName asserts configMapWrongNameError.
func IsConfigMapWrongName(err error) bool {
	return microerror.Cause(err) == configMapWrongNameError
}

var configMapWrongNamespaceError = microerror.New("configmap wrong namespace")

// IsConfigMapWrongNamespace asserts configMapWrongNamespaceError.
func IsConfigMapWrongNamespace(err error) bool {
	return microerror.Cause(err) == configMapWrongNamespaceError
}
