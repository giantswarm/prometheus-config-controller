package certificate

import (
	"github.com/giantswarm/operatorkit/framework"
)

const (
	Name = "certificate"
)

type Config struct{}

func DefaultConfig() Config {
	return Config{}
}

type Resource struct{}

func New(config Config) (*Resource, error) {
	resource := &Resource{}

	return resource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
