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
