package app

import (
	"context"
	"fmt"

	"github.com/inpour/event-driven-trader/shared/events"
	"github.com/inpour/event-driven-trader/shared/market"
)

type CandleBatchReader interface {
	ReadCandles(
		ctx context.Context,
		runID string,
		expectedCount int,
	) ([]events.Envelope[events.CandleClosed], error)
}

type CandleLoader struct {
	reader CandleBatchReader
}

func NewCandleLoader(reader CandleBatchReader) *CandleLoader {
	return &CandleLoader{
		reader: reader,
	}
}

func (l *CandleLoader) Load(
	ctx context.Context,
	input events.Envelope[events.BacktestInputCompleted],
) ([]*market.Candle, error) {
	candleEvents, err := l.reader.ReadCandles(
		ctx,
		input.RunID,
		input.Payload.CandleCount,
	)
	if err != nil {
		return nil, fmt.Errorf("read candle batch: %w", err)
	}

	if len(candleEvents) != input.Payload.CandleCount {
		return nil, fmt.Errorf(
			"unexpected candle count: expected=%d actual=%d",
			input.Payload.CandleCount,
			len(candleEvents),
		)
	}

	candles := make([]*market.Candle, 0, len(candleEvents))

	for _, event := range candleEvents {
		if event.RunID != input.RunID {
			return nil, fmt.Errorf(
				"candle event belongs to another run: expected=%q actual=%q",
				input.RunID,
				event.RunID,
			)
		}

		candles = append(candles, &market.Candle{
			Time:   event.Payload.Candle.Time,
			Open:   event.Payload.Candle.Open,
			High:   event.Payload.Candle.High,
			Low:    event.Payload.Candle.Low,
			Close:  event.Payload.Candle.Close,
			Volume: event.Payload.Candle.Volume,
		})
	}

	return candles, nil
}
