package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/inpour/event-driven-trader/shared/events"
)

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(
	broker string,
	groupID string,
	startOffset string,
) (*Consumer, error) {
	offset, err := parseStartOffset(startOffset)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:     []string{broker},
			GroupID:     groupID,
			Topic:       TopicCandleClosed,
			StartOffset: offset,
		}),
	}, nil
}

// FetchCandleClosed waits for the next candle-close event.
//
// It does not commit the Kafka offset. Call Commit after successful handling.
func (c *Consumer) FetchCandleClosed(
	ctx context.Context,
) (events.Envelope[events.CandleClosed], kafkago.Message, error) {
	message, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return events.Envelope[events.CandleClosed]{}, kafkago.Message{}, fmt.Errorf(
			"fetch Kafka message: %w",
			err,
		)
	}

	var event events.Envelope[events.CandleClosed]

	if err := json.Unmarshal(message.Value, &event); err != nil {
		return events.Envelope[events.CandleClosed]{}, kafkago.Message{}, fmt.Errorf(
			"decode CandleClosed event at partition %d offset %d: %w",
			message.Partition,
			message.Offset,
			err,
		)
	}

	if event.EventType != events.EventTypeCandleClosed {
		return events.Envelope[events.CandleClosed]{}, kafkago.Message{}, fmt.Errorf(
			"unexpected event type %q at partition %d offset %d",
			event.EventType,
			message.Partition,
			message.Offset,
		)
	}

	return event, message, nil
}

// Commit marks a fetched message as successfully processed.
func (c *Consumer) Commit(
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

// Close closes the underlying Kafka reader.
func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("close Kafka reader: %w", err)
	}

	return nil
}

func parseStartOffset(value string) (int64, error) {
	switch value {
	case "earliest":
		return kafkago.FirstOffset, nil
	case "latest":
		return kafkago.LastOffset, nil
	default:
		return 0, fmt.Errorf("unsupported start offset %q", value)
	}
}
