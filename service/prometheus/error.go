package prometheus

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var reloadError = microerror.New("reload")

// IsReloadError asserts reloadError.
func IsReloadError(err error) bool {
	return microerror.Cause(err) == reloadError
}
