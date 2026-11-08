package api

import "github.com/kelseyhightower/envconfig"

// Config holds the application settings, loaded from environment variables.
type Config struct {
	AppName string `default:"bookmark_service" envconfig:"APP_NAME"`
	AppPort int    `default:"8080" envconfig:"APP_PORT"`
	// InstanceID is generated as a UUID at start-up when left empty.
	InstanceID string `envconfig:"INSTANCE_ID"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
