package redis

import "github.com/kelseyhightower/envconfig"

// Config holds the Redis connection settings, loaded from environment variables.
type Config struct {
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

func newConfig(envPrefix string) (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
