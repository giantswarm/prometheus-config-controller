package certificate

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var missingError = microerror.New("missing")

// IsMissing asserts missingError.
func IsMissing(err error) bool {
	return microerror.Cause(err) == missingError
}
