package backtest

import (
	"errors"
)

var (
	ErrInvalidMetricsParameter = errors.New("metrics parameter is invalid")
	ErrNoPositionFound         = errors.New("there is no position")
)

type Metrics struct {
	TotalCandles, TpCnt, SlCnt int
	RR                         float64
}

func NewMetrics(totalCandles, tpCnt, slCnt int, rr float64) (*Metrics, error) {
	if totalCandles < 0 || tpCnt < 0 || slCnt < 0 {
		return nil, ErrInvalidMetricsParameter
	}

	metrics := &Metrics{
		TotalCandles: totalCandles,
		TpCnt:        tpCnt,
		SlCnt:        slCnt,
		RR:           rr,
	}

	return metrics, nil
}

func (m *Metrics) WinRate() (float64, error) {
	if m.SlCnt+m.TpCnt == 0 {
		return 0, ErrNoPositionFound
	}
	return float64(m.TpCnt) / float64(m.SlCnt+m.TpCnt), nil
}

func (m *Metrics) WeightedWinRate() (float64, error) {
	winRate, err := m.WinRate()
	if err != nil {
		return 0, err
	}
	return winRate * m.RR / (winRate*m.RR + (1 - winRate)), nil
}

func (m *Metrics) Profit() (float64, error) {
	winRate, err := m.WinRate()
	if err != nil {
		return 0, err
	}
	profit := winRate * m.RR
	loss := 1 - winRate
	return profit - loss, nil
}
