package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValues(t *testing.T) {
	cfg := LoadFromEnv()
	// Check default values
	assert.NotEmpty(t, cfg.Name)
	assert.Equal(t, "Culture 360", cfg.Name)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, true, cfg.Metrics.Prometheus.Enabled)
	assert.Equal(t, "/metrics", cfg.Metrics.Prometheus.HTTPPath)
	assert.Equal(t, "8080", cfg.Metrics.Prometheus.Port)
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("APP_NAME", "Culture 360")
	defer os.Unsetenv("APP_NAME")
	os.Setenv("APP_MYSQL_USERNAME", "username")
	defer os.Unsetenv("APP_MYSQL_USERNAME")
	os.Setenv("APP_MYSQL_PASSWORD", "password")
	defer os.Unsetenv("APP_MYSQL_PASSWORD")
	os.Setenv("APP_MYSQL_TESTMODE", "true")
	defer os.Unsetenv("APP_MYSQL_TESTMODE")
	os.Setenv("APP_MYSQL_CONNECTION", "connection")
	defer os.Unsetenv("APP_MYSQL_CONNECTION")
	os.Setenv("APP_MYSQL_CONNRETRYSLEEP", "20s")
	defer os.Unsetenv("APP_MYSQL_CONNRETRYSLEEP")
	os.Setenv("APP_MYSQL_CONNMAXLIFETIME", "20s")
	defer os.Unsetenv("APP_MYSQL_CONNMAXLIFETIME")
	os.Setenv("APP_MYSQL_MAXIDLECONNS", "20")
	defer os.Unsetenv("APP_MYSQL_MAXIDLECONNS")
	os.Setenv("APP_MYSQL_MAXOPENCONNS", "20")
	defer os.Unsetenv("APP_MYSQL_MAXOPENCONNS")
	os.Setenv("APP_LOGGING_LEVEL", "debug")
	defer os.Unsetenv("APP_LOGGING_LEVEL")
	os.Setenv("APP_METRICS_PROMETHEUS_ENABLED", "false")
	defer os.Unsetenv("APP_METRICS_PROMETHEUS_ENABLED")
	os.Setenv("APP_METRICS_PROMETHEUS_HTTPPATH", "/actuator/metrics")
	defer os.Unsetenv("APP_METRICS_PROMETHEUS_HTTPPATH")
	os.Setenv("APP_METRICS_PROMETHEUS_PORT", "4040")
	defer os.Unsetenv("APP_METRICS_PROMETHEUS_PORT")
	os.Setenv("APP_NOTIFICATION_AUTHTOKEN", "xoxb")
	defer os.Unsetenv("APP_NOTIFICATION_AUTHTOKEN")
	os.Setenv("APP_NOTIFICATION_APPTOKEN", "xapp")
	defer os.Unsetenv("APP_NOTIFICATION_APPTOKEN")

	cfg := LoadFromEnv()

	// Check custom values
	assert.NotEmpty(t, cfg.Name)
	assert.Equal(t, "Culture 360", cfg.Name)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, false, cfg.Metrics.Prometheus.Enabled)
	assert.Equal(t, "/actuator/metrics", cfg.Metrics.Prometheus.HTTPPath)
	assert.Equal(t, "4040", cfg.Metrics.Prometheus.Port)
	assert.Equal(t, "xoxb", cfg.Notification.AuthToken)
	assert.Equal(t, "xapp", cfg.Notification.AppToken)
	assert.Equal(t, true, cfg.MySQL.TestMode)
	assert.Equal(t, "connection", cfg.MySQL.Connection)
	assert.Equal(t, 20*time.Second, cfg.MySQL.ConnRetrySleep)
	assert.Equal(t, 20*time.Second, cfg.MySQL.ConnMaxLifetime)
	assert.Equal(t, uint(20), cfg.MySQL.MaxIdleConns)
	assert.Equal(t, uint(20), cfg.MySQL.MaxOpenConns)
}
