package events

import "github.com/inpour/event-driven-trader/shared/backtest"

const EventTypeBacktestMarkerCreated = "backtest.marker.created"

// BacktestMarkerCreated is published when a strategy creates a chart marker.
type BacktestMarkerCreated struct {
	MarkerID        string          `json:"marker_id"`
	StrategyName    string          `json:"strategy_name"`
	StrategyVersion string          `json:"strategy_version"`
	Marker          backtest.Marker `json:"marker"`
}
