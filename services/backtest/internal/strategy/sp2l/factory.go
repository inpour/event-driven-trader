package sp2l

import (
	"errors"
	"fmt"

	"github.com/inpour/event-driven-trader/shared/market"
)

var ErrInvalidConfig = errors.New("invalid " + strategyName + " " + strategyVersion + " config")

func New(cfg *Config) (*Strategy, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	return &Strategy{
		tpCnt:       0,
		slCnt:       0,
		positionCnt: 0,
		cfg:         *cfg,
	}, nil
}

func DefaultConfig() *Config {
	return &Config{
		RR:                  1,
		MinGapDistance:      0.03,
		MaxEntranceDistance: 0.3,
		SlExtraSpace:        0.003,
		Session: market.Session{
			StartH:  0,
			StartM:  0,
			StopH:   23,
			StopM:   59,
			Reverse: false,
		},
		EMAPeriod: 20,
	}
}
