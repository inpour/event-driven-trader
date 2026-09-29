package app

import (
	"context"
	"fmt"

	"github.com/inpour/event-driven-trader/services/backtest/internal/strategy"
	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/market"
)

type Runner struct {
	strategy strategy.Strategy
}

func NewRunner(strategy strategy.Strategy) *Runner {
	return &Runner{
		strategy: strategy,
	}
}

func (r *Runner) Run(ctx context.Context, candles []*market.Candle) ([]*backtest.Marker, error) {
	markers, err := r.strategy.Probe(ctx, candles)
	if err != nil {
		info := r.strategy.Info()
		return nil, fmt.Errorf(
			"run strategy %s@%s: %w",
			info.Name,
			info.Version,
			err,
		)
	}

	return markers, nil
}
