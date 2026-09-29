package sp2l

import (
	"errors"

	"github.com/inpour/event-driven-trader/shared/market"
)

var (
	ErrInvalidRiskRewardRatio     = errors.New("risk-reward Ratio must be greater than zero")
	ErrInvalidMinGapDistance      = errors.New("minimum gap distance must be greater than or equal to zero")
	ErrInvalidMaxEntranceDistance = errors.New("maximum entrance distance must be greater than zero")
	ErrInvalidEMAPeriod           = errors.New("exponential moving average must be greater period than zero")
)

type Config struct {
	RR                  float64        // Risk-Reward ratio
	MinGapDistance      float64        // distance between previous and next candle of gap-candle (percent %)
	MaxEntranceDistance float64        // maximum distance between sl and entry (percent %)
	SlExtraSpace        float64        // extra space for sl, avoid being tangential (percent %)
	Session             market.Session // out of session will be ignored
	EMAPeriod           int            // exponential moving average period
}

func (c *Config) Validate() error {
	if c.RR <= 0 {
		return ErrInvalidRiskRewardRatio
	}
	if c.MinGapDistance < 0 {
		return ErrInvalidMinGapDistance
	}
	if c.MinGapDistance <= 0 {
		return ErrInvalidMaxEntranceDistance
	}
	if c.EMAPeriod <= 0 {
		return ErrInvalidEMAPeriod
	}

	return nil
}
