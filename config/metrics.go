package config

// Metrics is the configuration of metrics in the service
type Metrics struct {
	// Prometheus configuration
	Prometheus PrometheusConfig
}
type PrometheusConfig struct {
	// Enabled indicates whether Prometheus instrumentation will be active.
	// Environment variable: APP_METRICS_PROMETHEUS_ENABLED
	Enabled bool `default:"true"`

	// HTTPPath indicates the URI path where metric will be served by the HTTP server
	// Environment variable: APP_METRICS_PROMETHEUS_HTTPPATH
	HTTPPath string `default:"/metrics" validate:"uri"`

	// Port indicates the port where metric will be served by the HTTP server
	// Environment variable: APP_METRICS_PROMETHEUS_PORT
	Port string `default:"8080"`
}
