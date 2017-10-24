// Package prometheustest provides a test implementation of the PrometheusReloader interface.
package prometheustest

type TestService struct{}

func New() *TestService {
	return &TestService{}
}

func (s *TestService) Reload() error {
	return nil
}
