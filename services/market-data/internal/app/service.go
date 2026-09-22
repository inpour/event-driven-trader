package app

import (
	"context"
	"fmt"
	"time"

	"github.com/inpour/event-driven-trader/services/market-data/internal/kafka"
	"github.com/inpour/event-driven-trader/services/market-data/internal/provider"
	"github.com/inpour/event-driven-trader/shared/events"
)

type Service struct {
	provider  provider.Provider
	producer  kafka.EventProducer
	runID     string
	symbol    string
	timeframe string
}

func NewService(
	provider provider.Provider,
	producer kafka.EventProducer,
	runID string,
	symbol string,
	timeframe string,
) *Service {
	return &Service{
		provider:  provider,
		producer:  producer,
		runID:     runID,
		symbol:    symbol,
		timeframe: timeframe,
	}
}

func (s *Service) Run(ctx context.Context) error {
	candles, err := s.provider.Candles()
	if err != nil {
		return fmt.Errorf("read candles: %w", err)
	}

	for index, candle := range candles {
		event := events.Envelope[events.CandleClosed]{
			EventID:    fmt.Sprintf("%s-candle-%06d", s.runID, index+1),
			EventType:  events.EventTypeCandleClosed,
			OccurredAt: time.Now().UTC(),
			Producer:   "market-data-service",
			RunID:      s.runID,
			Payload: events.CandleClosed{
				Symbol:    s.symbol,
				Timeframe: s.timeframe,
				Candle:    *candle,
			},
		}

		if err := s.producer.PublishCandleClosed(ctx, event); err != nil {
			return fmt.Errorf(
				"publish candle %d: %w",
				index+1,
				err,
			)
		}
	}

	completedAt := time.Now().UTC()

	completedEvent := events.Envelope[events.BacktestInputCompleted]{
		EventID:    fmt.Sprintf("%s-input-completed", s.runID),
		EventType:  events.EventTypeBacktestInputCompleted,
		OccurredAt: completedAt,
		Producer:   "market-data-service",
		RunID:      s.runID,
		Payload: events.BacktestInputCompleted{
			CandleCount: len(candles),
			Symbol:      s.symbol,
			Timeframe:   s.timeframe,
			CompletedAt: completedAt,
		},
	}

	if err := s.producer.PublishBacktestInputCompleted(ctx, completedEvent); err != nil {
		return fmt.Errorf("publish backtest-input-completed event: %w", err)
	}

	return nil
}
