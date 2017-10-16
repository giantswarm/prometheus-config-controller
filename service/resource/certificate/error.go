package certificate

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var missingCertificateError = microerror.New("missing certificate")

// IsMissingCertificate asserts missingCertificateError.
func IsMissingCertificate(err error) bool {
	return microerror.Cause(err) == missingCertificateError
}
