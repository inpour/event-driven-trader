package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/inpour/event-driven-trader/shared/events"
)

type InputCompletedConsumer struct {
	reader *kafkago.Reader
}

func NewInputCompletedConsumer(
	broker,
	groupID string,
	topic string,
) *InputCompletedConsumer {
	return &InputCompletedConsumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        []string{broker},
			GroupID:        groupID,
			Topic:          topic,
			StartOffset:    kafkago.FirstOffset,
			CommitInterval: 0,
			MinBytes:       1,
			MaxBytes:       10e6,
		}),
	}
}

func (c *InputCompletedConsumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("close Kafka reader: %w", err)
	}

	return nil
}

func (c *InputCompletedConsumer) Fetch(
	ctx context.Context,
) (
	events.Envelope[events.BacktestInputCompleted],
	kafkago.Message,
	error,
) {
	message, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return events.Envelope[events.BacktestInputCompleted]{},
			kafkago.Message{},
			fmt.Errorf("fetch input-completed event: %w", err)
	}

	var event events.Envelope[events.BacktestInputCompleted]
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return events.Envelope[events.BacktestInputCompleted]{},
			kafkago.Message{},
			fmt.Errorf(
				"decode input-completed event at offset %d: %w",
				message.Offset,
				err,
			)
	}

	return event, message, nil
}

func (c *InputCompletedConsumer) Commit(
	ctx context.Context,
	message kafkago.Message,
) error {
	if err := c.reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf(
			"commit Kafka message at partition %d offset %d: %w",
			message.Partition,
			message.Offset,
			err,
		)
	}

	return nil
}
