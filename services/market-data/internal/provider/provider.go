package provider

import (
	"github.com/inpour/event-driven-trader/shared/market"
)

type Provider interface {
	Candles() ([]*market.Candle, error)
}
