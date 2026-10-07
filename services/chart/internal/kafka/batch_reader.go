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
	markerTopic string
}

func NewBatchReader(
	broker string,
	candleTopic string,
	markerTopic string,
) *BatchReader {
	return &BatchReader{
		brokers:     []string{broker},
		candleTopic: candleTopic,
		markerTopic: markerTopic,
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

func (r *BatchReader) ReadMarkers(
	ctx context.Context,
	runID string,
	expectedCount int,
) ([]events.Envelope[events.BacktestMarkerCreated], error) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     r.brokers,
		Topic:       r.markerTopic,
		Partition:   0,
		StartOffset: kafkago.FirstOffset, // to read all markers from beginning
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()

	markerEvents := make(
		[]events.Envelope[events.BacktestMarkerCreated],
		0,
		expectedCount,
	)

	for len(markerEvents) < expectedCount {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("read marker message: %w", err)
		}

		var event events.Envelope[events.BacktestMarkerCreated]
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return nil, fmt.Errorf(
				"decode marker event at offset %d: %w",
				message.Offset,
				err,
			)
		}

		if event.EventType != events.EventTypeBacktestMarkerCreated {
			continue
		}

		if event.RunID != runID {
			continue
		}

		markerEvents = append(markerEvents, event)
	}

	return markerEvents, nil
}
