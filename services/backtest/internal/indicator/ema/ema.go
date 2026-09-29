package ema

import (
	"errors"
)

var ErrInvalidPeriod = errors.New("EMA period must be greater than zero")

// EMA calculates an exponential moving average incrementally.
type EMA struct {
	period        int
	alpha         float64
	samples       int
	currentResult *Result
}

func New(period int) (*EMA, error) {
	if period <= 0 {
		return nil, ErrInvalidPeriod
	}

	return &EMA{
		period:  period,
		alpha:   2.0 / float64(period+1),
		samples: 0,
		currentResult: &Result{
			Value:   0,
			IsReady: false,
		},
	}, nil
}

// Add updates the EMA with one new value.
//
// The first input becomes the initial EMA value. Later values use:
// EMA = alpha × current + (1 - alpha) × previousEMA.
//
// The returned Result contains the current EMA value. The returned boolean is
// true when the EMA is isReady for use; it is false while the indicator is
// warming up.
func (e *EMA) Add(value float64) Result {
	e.samples++
	if !e.currentResult.IsReady {
		e.currentResult.Value += value
		if e.samples == e.period {
			e.currentResult.Value /= float64(e.period)
			e.currentResult.IsReady = true
		}
		return *e.currentResult
	}

	e.currentResult.Value = e.alpha*value + (1.0-e.alpha)*e.currentResult.Value
	return *e.currentResult
}

// Value returns the current EMA currentResult and whether the EMA is isReady for use.
func (e *EMA) Value() Result {
	return *e.currentResult
}

func (e *EMA) IsReady() bool {
	return e.currentResult.IsReady
}

// Reset removes all accumulated EMA state.
func (e *EMA) Reset() {
	e.currentResult.Value = 0
	e.samples = 0
	e.currentResult.IsReady = false
}
