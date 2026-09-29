package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/inpour/event-driven-trader/shared/events"
)

type BatchReader struct {
	brokers     []string
	candleTopic string
}

func NewBatchReader(
	broker string,
	candleTopic string,
) *BatchReader {
	return &BatchReader{
		brokers:     []string{broker},
		candleTopic: candleTopic,
	}
}

func (r *BatchReader) ReadCandles(
	ctx context.Context,
	runID string,
	expectedCount int,
) ([]events.Envelope[events.CandleClosed], error) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     r.brokers,
		Topic:       r.candleTopic,
		Partition:   0,
		StartOffset: kafkago.FirstOffset, // to read all candles from beginning
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()

	candleEvents := make(
		[]events.Envelope[events.CandleClosed],
		0,
		expectedCount,
	)

	for len(candleEvents) < expectedCount {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("read candle message: %w", err)
		}

		var event events.Envelope[events.CandleClosed]
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return nil, fmt.Errorf(
				"decode candle event at offset %d: %w",
				message.Offset,
				err,
			)
		}

		if event.EventType != events.EventTypeCandleClosed {
			continue
		}

		if event.RunID != runID {
			continue
		}

		candleEvents = append(candleEvents, event)
	}

	return candleEvents, nil
}
