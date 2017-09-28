package endpoint

type Config struct{}

func DefaultConfig() Config {
	return Config{}
}

func New(config Config) (*Endpoint, error) {
	newEndpoint := &Endpoint{}

	return newEndpoint, nil
}

type Endpoint struct{}
