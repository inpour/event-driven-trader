package sp2l

import (
	"github.com/inpour/event-driven-trader/shared/backtest"
)

type Strategy struct {
	totalCandles, tpCnt, slCnt, positionCnt int
	cfg                                     Config
}

const (
	strategyName    = "SP2L"
	strategyVersion = "v3"
)

func (s *Strategy) Info() *backtest.StrategyInfo {
	return &backtest.StrategyInfo{
		Name:    strategyName,
		Version: strategyVersion,
		Config:  s.cfg,
	}
}

func (s *Strategy) Metrics() (*backtest.Metrics, error) {
	return backtest.NewMetrics(s.totalCandles, s.tpCnt, s.slCnt, s.cfg.RR)
}
