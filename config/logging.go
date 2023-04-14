package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// LoggingConfig is the configuration of logging in the service
type LoggingConfig struct {
	// Level indicates the level at which the logger will filter the events. Any event logged at a level
	// below this will be skipped.
	// Environment variable: APP_LOGGING_LEVEL
	Level string `default:"info"`
}

// Config the log format used in the service
func (*Config) NewLoggingFormatter() log.Formatter {
	return &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
}

// Map the logging level defined by configuration to its correspondent `log.Level`
func (c *Config) LoggingLevel() (log.Level, error) {
	return log.ParseLevel(c.Logging.Level)
}
