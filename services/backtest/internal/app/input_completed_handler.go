package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/inpour/event-driven-trader/shared/events"
)

type InputCompletedHandler struct {
	candleLoader    *CandleLoader
	runner          *Runner
	markerPublisher *MarkerPublisher
	logger          *slog.Logger
}

func NewInputCompletedHandler(
	candleLoader *CandleLoader,
	runner *Runner,
	markerPublisher *MarkerPublisher,
	logger *slog.Logger,
) *InputCompletedHandler {
	return &InputCompletedHandler{
		candleLoader:    candleLoader,
		runner:          runner,
		markerPublisher: markerPublisher,
		logger:          logger,
	}
}

func (h *InputCompletedHandler) Handle(
	ctx context.Context,
	input events.Envelope[events.BacktestInputCompleted],
) error {
	startedAt := time.Now().UTC()

	h.logger.Info(
		"starting backtest",
		"run_id", input.RunID,
		"symbol", input.Payload.Symbol,
		"timeframe", input.Payload.Timeframe,
		"expected_candle_count", input.Payload.CandleCount,
	)

	candles, err := h.candleLoader.Load(ctx, input)
	if err != nil {
		return fmt.Errorf("load candles: %w", err)
	}

	markers, err := h.runner.Run(ctx, candles)
	if err != nil {
		return fmt.Errorf("run backtest: %w", err)
	}

	if err := h.markerPublisher.PublishMarkers(ctx, input, markers); err != nil {
		return fmt.Errorf("publish markers: %w", err)
	}

	if err := h.markerPublisher.PublishRunCompleted(
		ctx,
		input,
		h.runner.strategy.Info().Name,
		h.runner.strategy.Info().Version,
		len(markers),
		startedAt,
	); err != nil {
		return err
	}

	h.logger.Info(
		"backtest completed",
		"run_id", input.RunID,
		"candle_count", len(candles),
		"marker_count", len(markers),
	)

	return nil
}
