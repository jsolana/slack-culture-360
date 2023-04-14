package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics defined for the job
var (
	BatchOpsProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "batch_ops_processed",
		Help: "The total number of processed registers.",
	})

	BatchResult = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "batch_result",
		Help: "Result of the batch job execution: 0 if OK other error.",
	})

	BatchHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "batch_duration_seconds",
		Help: "Histogram of bach processing duration in seconds.",
	})
)
