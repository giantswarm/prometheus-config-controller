package middleware

type Config struct{}

func DefaultConfig() Config {
	return Config{}
}

func New(config Config) (*Middleware, error) {
	return &Middleware{}, nil
}

type Middleware struct{}
