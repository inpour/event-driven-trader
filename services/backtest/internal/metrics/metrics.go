package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// BacktestDurationSeconds measures the latency of each backtest stage and the overall run.
	BacktestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "trader_backtest_duration_seconds",
			Help: "Duration of backtest pipeline stages in seconds.",
			// Buckets from 1ms (0.001s) to 10s to cover both fast CPU strategy runs and slower Kafka I/O
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 3, 4, 5, 7.5, 10},
		},
		[]string{"stage", "symbol", "timeframe"},
	)

	// MarkerPublishDurationSeconds measures individual Kafka write latency per marker.
	MarkerPublishDurationSeconds = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "trader_backtest_marker_publish_duration_seconds",
			Help:    "Duration of publishing a single marker event to Kafka in seconds.",
			Buckets: []float64{0.00005, 0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1},
		},
	)
)
