package app

import (
	"context"
	"fmt"
	"time"

	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/events"
)

type BacktestOutputPublisher interface {
	PublishMarkerCreated(
		ctx context.Context,
		event events.Envelope[events.BacktestMarkerCreated],
	) error

	PublishRunCompleted(
		ctx context.Context,
		event events.Envelope[events.BacktestRunCompleted],
	) error
}

type MarkerPublisher struct {
	publisher BacktestOutputPublisher
}

func NewMarkerPublisher(
	publisher BacktestOutputPublisher,
) *MarkerPublisher {
	return &MarkerPublisher{
		publisher: publisher,
	}
}

func (p *MarkerPublisher) PublishMarkers(
	ctx context.Context,
	input events.Envelope[events.BacktestInputCompleted],
	markers []*backtest.Marker,
) error {
	for index, marker := range markers {
		event := events.Envelope[events.BacktestMarkerCreated]{
			EventID:    fmt.Sprintf("%s-marker-%06d", input.RunID, index),
			EventType:  events.EventTypeBacktestMarkerCreated,
			Producer:   "backtest-service",
			RunID:      input.RunID,
			OccurredAt: time.Now().UTC(),
			Payload: events.BacktestMarkerCreated{
				Marker: backtest.Marker{
					Time:     marker.Time,
					Price:    marker.Price,
					Position: marker.Position,
					Shape:    marker.Shape,
					Color:    marker.Color,
					Size:     marker.Size,
					Text:     marker.Text,
				},
			},
		}

		if err := p.publisher.PublishMarkerCreated(ctx, event); err != nil {
			return fmt.Errorf(
				"publish marker %d: %w",
				index,
				err,
			)
		}
	}

	return nil
}

func (p *MarkerPublisher) PublishRunCompleted(
	ctx context.Context,
	input events.Envelope[events.BacktestInputCompleted],
	strategyName string,
	strategyVersion string,
	markerCount int,
	startedAt time.Time,
	metrics *backtest.Metrics,
) error {
	completedAt := time.Now().UTC()

	event := events.Envelope[events.BacktestRunCompleted]{
		EventID:    input.RunID + ":completed",
		EventType:  events.EventTypeBacktestRunCompleted,
		Producer:   "backtest-service",
		RunID:      input.RunID,
		OccurredAt: completedAt,
		Payload: events.BacktestRunCompleted{
			CandleCount:     input.Payload.CandleCount,
			MarkerCount:     markerCount,
			StrategyName:    strategyName,
			StrategyVersion: strategyVersion,
			Symbol:          input.Payload.Symbol,
			Timeframe:       input.Payload.Timeframe,
			StartedAt:       startedAt,
			Metrics:         *metrics,
		},
	}

	if err := p.publisher.PublishRunCompleted(ctx, event); err != nil {
		return fmt.Errorf("publish backtest run completed: %w", err)
	}

	return nil
}
