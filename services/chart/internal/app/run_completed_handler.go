package app

import (
	"context"
	"fmt"

	"github.com/inpour/event-driven-trader/services/chart/internal/render"
	"github.com/inpour/event-driven-trader/shared/events"
)

type RunCompletedHandler struct {
	candleLoader  *CandleLoader
	markerLoader  *MarkerLoader
	chartRenderer *render.ChartRenderer
}

func NewRunCompletedHandler(
	candleLoader *CandleLoader,
	markerLoader *MarkerLoader,
	chartRenderer *render.ChartRenderer,
) *RunCompletedHandler {
	return &RunCompletedHandler{
		candleLoader:  candleLoader,
		markerLoader:  markerLoader,
		chartRenderer: chartRenderer,
	}
}

func (h *RunCompletedHandler) Handle(
	ctx context.Context,
	input events.Envelope[events.BacktestRunCompleted],
) error {
	candles, err := h.candleLoader.Load(ctx, input)
	if err != nil {
		return fmt.Errorf("load candles: %w", err)
	}

	markers, err := h.markerLoader.Load(ctx, input)
	if err != nil {
		return fmt.Errorf("load markers: %w", err)
	}

	strategyInfo := input.Payload.StrategyInfo
	Metrics := input.Payload.Metrics

	err = h.chartRenderer.HTML(candles, markers, Metrics, strategyInfo)
	if err != nil {
		return fmt.Errorf("render chart: %w", err)
	}

	return nil
}
