package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/inpour/event-driven-trader/shared/events"
)

type Producer struct {
	markerWriter       *kafkago.Writer
	runCompletedWriter *kafkago.Writer
}

func NewProducer(
	broker string,
	markerCreatedTopic string,
	runCompletedTopic string,
) (*Producer, error) {

	if markerCreatedTopic == "" {
		return nil, errors.New("marker-created topic is required")
	}

	if runCompletedTopic == "" {
		return nil, errors.New("run-completed topic is required")
	}

	return &Producer{
		markerWriter: newWriter(
			broker,
			markerCreatedTopic,
		),
		runCompletedWriter: newWriter(
			broker,
			runCompletedTopic,
		),
	}, nil
}

func (p *Producer) PublishMarkerCreated(
	ctx context.Context,
	event events.Envelope[events.BacktestMarkerCreated],
) error {
	if event.EventType != events.EventTypeBacktestMarkerCreated {
		return fmt.Errorf(
			"unexpected marker event type: %q",
			event.EventType,
		)
	}

	return publish(ctx, p.markerWriter, event)
}

func (p *Producer) PublishRunCompleted(
	ctx context.Context,
	event events.Envelope[events.BacktestRunCompleted],
) error {
	if event.EventType != events.EventTypeBacktestRunCompleted {
		return fmt.Errorf(
			"unexpected run-completed event type: %q",
			event.EventType,
		)
	}

	return publish(ctx, p.runCompletedWriter, event)
}

func (p *Producer) Close() error {
	var result error

	if p.markerWriter != nil {
		if err := p.markerWriter.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close marker writer: %w", err))
		}
	}

	if p.runCompletedWriter != nil {
		if err := p.runCompletedWriter.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close run-completed writer: %w", err))
		}
	}

	return result
}

func newWriter(
	broker string,
	topic string,
) *kafkago.Writer {
	return &kafkago.Writer{
		Addr:         kafkago.TCP(broker),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		BatchSize:    1,
		BatchTimeout: 0,
	}
}

func publish[T any](
	ctx context.Context,
	writer *kafkago.Writer,
	event events.Envelope[T],
) error {
	if writer == nil {
		return errors.New("kafka writer is nil")
	}

	if event.RunID == "" {
		return errors.New("event run_id is required")
	}

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf(
			"marshal event type %q: %w",
			event.EventType,
			err,
		)
	}

	message := kafkago.Message{
		Key:   []byte(event.RunID),
		Value: value,
	}

	if err := writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf(
			"publish event type %q for run %q: %w",
			event.EventType,
			event.RunID,
			err,
		)
	}

	return nil
}
