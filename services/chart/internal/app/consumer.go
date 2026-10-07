package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/inpour/event-driven-trader/services/chart/internal/kafka"
)

type RunCompletedConsumer struct {
	consumer *kafka.RunCompletedConsumer
	handler  *RunCompletedHandler
}

func NewRunCompleteConsumer(
	consumer *kafka.RunCompletedConsumer,
	handler *RunCompletedHandler,
) *RunCompletedConsumer {
	return &RunCompletedConsumer{
		consumer: consumer,
		handler:  handler,
	}
}

func (c *RunCompletedConsumer) Consume(
	ctx context.Context,
) error {
	for {
		RunCompleted, message, err := c.consumer.Fetch(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			fmt.Fprintln(os.Stderr, fmt.Errorf("fetch backtest run-completed event: %w", err))
		}

		fmt.Fprintln(os.Stdout, "run-completed event received")

		if err := c.handler.Handle(ctx, RunCompleted); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			fmt.Fprintln(os.Stderr, fmt.Errorf("chart render failed: %w", err))
		}

		fmt.Fprintln(os.Stdout, "chart created")

		if err := c.consumer.Commit(ctx, message); err != nil {
			return fmt.Errorf(
				"commit successful input event for run %q at offset %d: %w",
				RunCompleted.RunID,
				message.Offset,
				err,
			)
		}
	}
}
