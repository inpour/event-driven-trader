package sp2l

import (
	"github.com/inpour/event-driven-trader/services/backtest/internal/strategy"
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

func (s *Strategy) Info() *strategy.Info {
	return &strategy.Info{
		Name:    strategyName,
		Version: strategyVersion,
		Config:  s.cfg,
	}
}

func (s *Strategy) Metrics() (*backtest.Metrics, error) {
	return backtest.NewMetrics(s.totalCandles, s.tpCnt, s.slCnt, s.cfg.RR)
}
