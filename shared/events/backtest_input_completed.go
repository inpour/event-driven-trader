package events

import "time"

const EventTypeBacktestInputCompleted = "backtest.input.completed"

// BacktestInputCompleted is published after market-data-service has sent
// every CandleClosed event for the current run.
type BacktestInputCompleted struct {
	CandleCount int       `json:"candle_count"`
	Symbol      string    `json:"symbol"`
	Timeframe   string    `json:"timeframe"`
	CompletedAt time.Time `json:"completed_at"`
}
