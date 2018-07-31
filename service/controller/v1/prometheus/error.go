package prometheus

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

var reloadError = &microerror.Error{
	Kind: "reloadError",
}

// IsReloadError asserts reloadError.
func IsReloadError(err error) bool {
	return microerror.Cause(err) == reloadError
}
