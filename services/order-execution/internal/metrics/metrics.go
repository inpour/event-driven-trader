package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ConsumerListening is 1 when the service loop is actively waiting for Kafka events, 0 when stopped.
	ConsumerListening = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "trader_order_execution_consumer_listening",
		Help: "1 if the order-execution Kafka consumer loop is actively listening for events, 0 otherwise.",
	})

	// CandlesReceivedTotal counts every CandleClosed event fetched from Kafka.
	CandlesReceivedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "trader_order_execution_candles_received_total",
		Help: "Total number of candle-closed events received by order-execution-service.",
	}, []string{"symbol", "timeframe"})

	// CandlesCommittedTotal counts every successfully committed offset.
	CandlesCommittedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "trader_order_execution_candles_committed_total",
		Help: "Total number of candle-closed Kafka offsets committed.",
	})

	// ConsumerErrorsTotal tracks fetch, decode, or commit errors.
	ConsumerErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "trader_order_execution_errors_total",
		Help: "Total number of errors in order-execution consumer loop by stage.",
	}, []string{"stage"}) // stage="fetch" or stage="commit"

	// LastEventReceivedTimestamp records the Unix timestamp of the latest event received.
	LastEventReceivedTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "trader_order_execution_last_event_timestamp_seconds",
		Help: "Unix timestamp of the last candle-closed event received.",
	})
)
