package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/inpour/event-driven-trader/services/backtest/internal/metrics"
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
	symbol := input.Payload.Symbol
	timeframe := input.Payload.Timeframe

	// Measure total end-to-end backtest duration when Handle exits
	defer func() {
		metrics.BacktestDurationSeconds.
			WithLabelValues("total", symbol, timeframe).
			Observe(time.Since(startedAt).Seconds())
	}()

	h.logger.Info(
		"starting backtest",
		"run_id", input.RunID,
		"symbol", symbol,
		"timeframe", timeframe,
		"expected_candle_count", input.Payload.CandleCount,
	)

	// 1. Measure Candle Loading from Kafka
	loadStart := time.Now()
	candles, err := h.candleLoader.Load(ctx, input)
	metrics.BacktestDurationSeconds.
		WithLabelValues("load_candles", symbol, timeframe).
		Observe(time.Since(loadStart).Seconds())
	if err != nil {
		return fmt.Errorf("load candles: %w", err)
	}

	// 2. Measure Strategy Execution (SP2L Probe)
	runStart := time.Now()
	markers, strategyMetrics, err := h.runner.Run(ctx, candles)
	metrics.BacktestDurationSeconds.
		WithLabelValues("strategy_probe", symbol, timeframe).
		Observe(time.Since(runStart).Seconds())
	if err != nil {
		return fmt.Errorf("run backtest: %w", err)
	}

	// 3. Measure Publishing Markers + RunCompleted to Kafka
	pubStart := time.Now()
	if err := h.markerPublisher.PublishMarkers(ctx, input, markers); err != nil {
		return fmt.Errorf("publish markers: %w", err)
	}

	if err := h.markerPublisher.PublishRunCompleted(
		ctx,
		input,
		h.runner.strategy.Info(),
		len(markers),
		startedAt,
		strategyMetrics,
	); err != nil {
		return err
	}
	metrics.BacktestDurationSeconds.
		WithLabelValues("publish_results", symbol, timeframe).
		Observe(time.Since(pubStart).Seconds())

	h.logger.Info(
		"backtest completed",
		"run_id", input.RunID,
		"candle_count", len(candles),
		"marker_count", len(markers),
	)

	return nil
}
