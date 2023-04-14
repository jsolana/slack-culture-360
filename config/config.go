package config

import (
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	envAppConfigPrefix = "APP"
)

// Config represents the whole app configuration available in the service
type Config struct {
	// Name is the application name.
	// Environment variable: APP_NAME
	Name string `default:"Culture 360"`
	// MySQL DB Connection configuration
	MySQL MySQL
	// Logging configuration
	Logging LoggingConfig
	// Metrics configuration
	Metrics Metrics
	// Notification configuration
	Notification SlackConfiguration
}

// Default returns the default configuration for the service.
func Default() *Config {
	// Hack: `envconfig` has no meanings to initialize the config with default values without
	// reading from the environment variables. What we do here is to use a surely unexciting prefix
	// so no value is found in the environment, but the defaults are applied.
	return loadFromEnv("_____a_prefix_that_should_not_exists")
}

// LoadFromEnv loads the configuration from the environment variables
// It panics if something goes wrong while loading the config.
func LoadFromEnv() *Config {
	return loadFromEnv(envAppConfigPrefix)
}

func loadFromEnv(prefix string) *Config {
	cfg := new(Config)
	if err := envconfig.Process(prefix, cfg); err != nil {
		log.WithError(err).Panic("Unexpected error while loading config from environment variables")
	}

	// Fix app name if unset
	if cfg.Name == "" {
		cfg.Name = filepath.Base(os.Args[0])
	}

	return cfg
}
