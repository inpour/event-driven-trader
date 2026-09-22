package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/inpour/event-driven-trader/shared/events"
	kafkago "github.com/segmentio/kafka-go"
)

type EventProducer interface {
	PublishCandleClosed(
		ctx context.Context,
		event events.Envelope[events.CandleClosed],
	) error

	PublishBacktestInputCompleted(
		ctx context.Context,
		event events.Envelope[events.BacktestInputCompleted],
	) error
}

// Producer publishes market-data events to Kafka.
type Producer struct {
	candleWriter         *kafkago.Writer
	inputCompletedWriter *kafkago.Writer
}

// NewProducer creates a Kafka producer for market-data events.
func NewProducer(broker string) *Producer {
	return &Producer{
		candleWriter:         newWriter(broker, TopicCandleClosed),
		inputCompletedWriter: newWriter(broker, TopicInputCompleted),
	}
}

func newWriter(broker string, topic string) *kafkago.Writer {
	return &kafkago.Writer{
		Addr:         kafkago.TCP(broker),
		Topic:        topic,
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
	}
}

// PublishCandleClosed publishes one completed candle event.
func (p *Producer) PublishCandleClosed(
	ctx context.Context,
	event events.Envelope[events.CandleClosed],
) error {
	return publish(ctx, p.candleWriter, event)
}

// PublishBacktestInputCompleted publishes an event indicating that all candles
// for the current backtest input have been published.
func (p *Producer) PublishBacktestInputCompleted(
	ctx context.Context,
	event events.Envelope[events.BacktestInputCompleted],
) error {
	return publish(ctx, p.inputCompletedWriter, event)
}

func publish[T any](
	ctx context.Context,
	writer *kafkago.Writer,
	event events.Envelope[T],
) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event %q: %w", event.EventType, err)
	}

	message := kafkago.Message{
		Key:   []byte(event.RunID),
		Value: value,
	}

	if err := writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("publish event %q: %w", event.EventType, err)
	}

	return nil
}

// Close closes all Kafka writers.
func (p *Producer) Close() error {
	var errs []error

	if err := p.candleWriter.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close candle writer: %w", err))
	}

	if err := p.inputCompletedWriter.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close input-completed writer: %w", err))
	}

	return errors.Join(errs...)
}
