package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/inpour/event-driven-trader/services/order-execution/internal/kafka"
	"github.com/inpour/event-driven-trader/services/order-execution/internal/logging"
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
	for {
		event, message, err := c.kafkaConsumer.FetchCandleClosed(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) ||
				errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("receive candle-close event: %w", err)
		}

		logging.LogCandleReceived(event, message)

		if err := c.kafkaConsumer.Commit(ctx, message); err != nil {
			return fmt.Errorf("commit candle-close event: %w", err)
		}
	}
}
