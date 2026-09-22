package events

import "github.com/inpour/event-driven-trader/shared/market"

const EventTypeCandleClosed = "market.candle.closed"

// CandleClosed is published when one completed OHLCV candle is available.
type CandleClosed struct {
	Symbol    string        `json:"symbol"`
	Timeframe string        `json:"timeframe"`
	Candle    market.Candle `json:"candle"`
}
