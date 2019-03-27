package prometheus

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var reloadThrottleError = &microerror.Error{
	Kind: "reloadThrottleError",
}

// IsReloadThrottle asserts reloadThrottleError.
func IsReloadThrottle(err error) bool {
	return microerror.Cause(err) == reloadThrottleError
}
