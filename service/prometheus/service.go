package prometheus

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Logger micrologger.Logger

	Address string
}

func DefaultConfig() Config {
	return Config{
		Logger: nil,

		Address: "",
	}
}

type Service struct {
	logger micrologger.Logger

	// address is concatenated with the reload path in New.
	address string
}

func New(config Config) (*Service, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.Address == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Address must not be empty")
	}

	u, err := url.ParseRequestURI(config.Address)
	if err != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Address is invalid: %s", err)
	}
	u.Path = path.Join(u.Path, prometheusReloadPath)

	service := &Service{
		logger: config.Logger,

		address: u.String(),
	}

	return service, nil
}

func (s *Service) Reload() error {
	s.logger.Log("debug", fmt.Sprintf("reloading prometheus config: %s", s.address))

	res, err := http.Post(s.address, "", nil)
	if err != nil {
		return microerror.Mask(err)
	}
	if res.StatusCode != http.StatusOK {
		return microerror.Maskf(reloadError, "a non-200 status code was returned: %s", res.StatusCode)
	}

	configurationReloadCount.Inc()

	return nil
}
