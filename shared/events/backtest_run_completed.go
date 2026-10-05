package events

import (
	"time"

	"github.com/inpour/event-driven-trader/shared/backtest"
)

const EventTypeBacktestRunCompleted = "backtest.run.completed"

// BacktestRunCompleted is published after backtest-service has:
// 1. loaded all candle events,
// 2. evaluated the strategy,
// 3. published all marker events.
type BacktestRunCompleted struct {
	CandleCount     int              `json:"candle_count"`
	MarkerCount     int              `json:"marker_count"`
	StrategyName    string           `json:"strategy_name"`
	StrategyVersion string           `json:"strategy_version"`
	Symbol          string           `json:"symbol"`
	Timeframe       string           `json:"timeframe"`
	StartedAt       time.Time        `json:"started_at"`
	CompletedAt     time.Time        `json:"completed_at"`
	Metrics         backtest.Metrics `json:"metrics"`
}
