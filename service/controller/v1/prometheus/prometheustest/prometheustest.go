// Package prometheustest provides a test implementation of the PrometheusReloader interface.
package prometheustest

import "context"

type TestService struct{}

func New() *TestService {
	return &TestService{}
}

func (s *TestService) Reload(ctx context.Context) error {
	return nil
}

func (s *TestService) RequestReload(ctx context.Context) {}
