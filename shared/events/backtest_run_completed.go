package events

import (
	"time"

	"github.com/inpour/event-driven-trader/shared/backtest"
)

const EventTypeBacktestRunCompleted = "backtest.run.completed"

// BacktestRunCompleted is published after backtest-service has finished.
type BacktestRunCompleted struct {
	CandleCount int       `json:"candle_count"`
	MarkerCount int       `json:"marker_count"`
	Symbol      string    `json:"symbol"`
	Timeframe   string    `json:"timeframe"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`

	StrategyInfo backtest.StrategyInfo `json:"strategy_info"`
	Metrics      backtest.Metrics      `json:"metrics"`
}
