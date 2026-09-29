package events

import "github.com/inpour/event-driven-trader/shared/backtest"

const EventTypeBacktestMarkerCreated = "backtest.marker.created"

// BacktestMarkerCreated is published when a strategy creates a chart marker.
type BacktestMarkerCreated struct {
	Marker backtest.Marker `json:"marker"`
}
