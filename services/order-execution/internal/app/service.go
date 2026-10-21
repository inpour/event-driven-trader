package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/inpour/event-driven-trader/services/order-execution/internal/kafka"
	"github.com/inpour/event-driven-trader/services/order-execution/internal/logging"
	"github.com/inpour/event-driven-trader/services/order-execution/internal/metrics"
)

type Service struct {
	kafkaConsumer *kafka.Consumer
}

func NewService(kafkaConsumer *kafka.Consumer) *Service {
	return &Service{
		kafkaConsumer: kafkaConsumer,
	}
}

// Run continuously receives CandleClosed events until ctx is canceled.
func (c *Service) Run(ctx context.Context) error {
	// Mark consumer as actively listening (1), and reset to 0 when Run exits
	metrics.ConsumerListening.Set(1)
	defer metrics.ConsumerListening.Set(0)

	for {
		event, message, err := c.kafkaConsumer.FetchCandleClosed(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) ||
				errors.Is(err, io.EOF) {
				return nil
			}

			metrics.ConsumerErrorsTotal.WithLabelValues("fetch").Inc()
			return fmt.Errorf("receive candle-close event: %w", err)
		}

		// Record that an event arrived
		metrics.CandlesReceivedTotal.WithLabelValues(
			event.Payload.Symbol,
			event.Payload.Timeframe,
		).Inc()
		metrics.LastEventReceivedTimestamp.SetToCurrentTime()

		logging.LogCandleReceived(event, message)

		if err := c.kafkaConsumer.Commit(ctx, message); err != nil {
			metrics.ConsumerErrorsTotal.WithLabelValues("commit").Inc()
			return fmt.Errorf("commit candle-close event: %w", err)
		}

		metrics.CandlesCommittedTotal.Inc()
	}
}
