package configmap

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var configMapNotFoundError = &microerror.Error{
	Kind: "configMapNotFoundError",
}

// IsConfigMapNotFound asserts configMapNotFoundError.
func IsConfigMapNotFound(err error) bool {
	return microerror.Cause(err) == configMapNotFoundError
}

var configMapKeyNotFoundError = &microerror.Error{
	Kind: "configMapKeyNotFoundError",
}

// IsConfigMapKeyNotFound asserts configMapKeyNotFoundError.
func IsConfigMapKeyNotFound(err error) bool {
	return microerror.Cause(err) == configMapKeyNotFoundError
}

var invalidConfigMapError = &microerror.Error{
	Kind: "invalidConfigMapError",
}

// IsInvalidConfigMap asserts invalidConfigMapError.
func IsInvalidConfigMap(err error) bool {
	return microerror.Cause(err) == invalidConfigMapError
}

var wrongNameError = &microerror.Error{
	Kind: "wrongNameError",
}

// IsWrongName asserts wrongNameError.
func IsWrongName(err error) bool {
	return microerror.Cause(err) == wrongNameError
}

var wrongNamespaceError = &microerror.Error{
	Kind: "wrongNamespaceError",
}

// IsWrongNamespace asserts wrongNamespaceError.
func IsWrongNamespace(err error) bool {
	return microerror.Cause(err) == wrongNamespaceError
}

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
