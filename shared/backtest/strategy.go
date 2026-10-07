package backtest

import (
	"context"

	"github.com/inpour/event-driven-trader/shared/market"
)

type Strategy interface {
	Info() *StrategyInfo
	Probe(ctx context.Context, candles []*market.Candle) ([]*Marker, error)
	Metrics() (*Metrics, error)
}

type StrategyInfo struct {
	Name    string
	Version string
	Config  any
}
