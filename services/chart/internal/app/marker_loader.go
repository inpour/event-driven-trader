package app

import (
	"context"
	"fmt"

	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/events"
)

type MarkerBatchReader interface {
	ReadMarkers(
		ctx context.Context,
		runID string,
		expectedCount int,
	) ([]events.Envelope[events.BacktestMarkerCreated], error)
}

type MarkerLoader struct {
	reader MarkerBatchReader
}

func NewMarkerLoader(reader MarkerBatchReader) *MarkerLoader {
	return &MarkerLoader{
		reader: reader,
	}
}

func (l *MarkerLoader) Load(
	ctx context.Context,
	input events.Envelope[events.BacktestRunCompleted],
) ([]*backtest.Marker, error) {
	markerEvents, err := l.reader.ReadMarkers(
		ctx,
		input.RunID,
		input.Payload.MarkerCount,
	)
	if err != nil {
		return nil, fmt.Errorf("read marker batch: %w", err)
	}

	if len(markerEvents) != input.Payload.MarkerCount {
		return nil, fmt.Errorf(
			"unexpected marker count: expected=%d actual=%d",
			input.Payload.MarkerCount,
			len(markerEvents),
		)
	}

	markers := make([]*backtest.Marker, 0, len(markerEvents))

	for _, event := range markerEvents {
		if event.RunID != input.RunID {
			return nil, fmt.Errorf(
				"marker event belongs to another run: expected=%q actual=%q",
				input.RunID,
				event.RunID,
			)
		}

		markers = append(markers, &backtest.Marker{
			Time:     event.Payload.Marker.Time,
			Price:    event.Payload.Marker.Price,
			Position: event.Payload.Marker.Position,
			Shape:    event.Payload.Marker.Shape,
			Color:    event.Payload.Marker.Color,
			Size:     event.Payload.Marker.Size,
			Text:     event.Payload.Marker.Text,
		})
	}

	return markers, nil
}
