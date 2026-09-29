package strategy

import (
	"context"

	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/market"
)

type Strategy interface {
	Info() *Info
	Probe(ctx context.Context, candles []*market.Candle) ([]*backtest.Marker, error)
	Metrics() (*backtest.Metrics, error)
}

type Info struct {
	Name    string
	Version string
	Config  any
}
