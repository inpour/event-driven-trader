package sp2l

import (
	"context"
	"fmt"
	"sync"

	"github.com/inpour/event-driven-trader/services/backtest/internal/indicator/ema"
	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/market"
	"github.com/inpour/event-driven-trader/shared/trade"
)

func (s *Strategy) Probe(ctx context.Context, candles []*market.Candle) ([]*backtest.Marker, error) {
	if len(candles) < 25 {
		return nil, fmt.Errorf("expected at least %d candles, got %d", 25, len(candles))
	}

	s.totalCandles = len(candles)

	emaResults, err := s.prepareEMA(candles)
	if err != nil {
		return nil, err
	}

	ch := make(chan *backtest.Marker)

	go s.probe(candles, emaResults, ch)

	markers := make([]*backtest.Marker, 0)
	for marker_ := range ch {
		markers = append(markers, marker_)
	}

	return markers, nil
}

func (s *Strategy) prepareEMA(candles []*market.Candle) ([]*ema.Result, error) {
	results := make([]*ema.Result, 0, len(candles))
	ema_, err := ema.New(s.cfg.EMAPeriod)
	if err != nil {
		return nil, err
	}

	for _, candle := range candles {
		result := ema_.Add(candle.Close)
		results = append(results, &result)
	}

	return results, nil
}

func (s *Strategy) probe(candles []*market.Candle, emaResults []*ema.Result, ch chan<- *backtest.Marker) {
	wg := new(sync.WaitGroup)
	for i := 4; i < len(candles)-1; i++ {
		// assume that candles[i-1] is gap-candle
		if !emaResults[i-2].IsReady {
			continue // EMA indicator is not ready
		}
		if !s.cfg.Session.In(candles[i-1].Time) {
			continue
		}
		gap := s.findGap(candles[i-2 : i+1]) // candles[i-1] is gap or not
		if gap == noGap {
			continue
		}
		ok := s.checkEMA(candles[i-2:i+1], emaResults[i-2:i+1])
		if !ok {
			continue
		}
		s.positionCnt++
		side, sl := s.gapData(gap, candles[i-4:i])
		s.gapSignal(s.positionCnt, candles[i-1].Time, candles[i-1].Open, ch)
		wg.Add(1)
		go s.enter(wg, s.positionCnt, side, candles[i:], sl, ch)
		i = exitGapArea(candles[i:], i, side)
	}
	wg.Wait()
	close(ch)
}

type gapType int

const (
	noGap gapType = iota
	risingGap
	fallingGap
)

// window[1] is gap-candle or not?
func (s *Strategy) findGap(window []*market.Candle) gapType {
	checkDistance := func(gap gapType, low, high float64) gapType {
		distance := 1 + s.cfg.MinGapDistance/100
		if high > low*distance {
			return gap
		}
		return noGap
	}
	prev := window[0]
	next := window[2]

	switch {
	case prev.High < next.Low:
		return checkDistance(risingGap, prev.High, next.Low)
	case next.High < prev.Low:
		return checkDistance(fallingGap, next.High, prev.Low)
	default:
		return noGap
	}
}

func (s *Strategy) checkEMA(window []*market.Candle, emaWindow []*ema.Result) bool {
	prev := window[0]
	current := window[1]
	next := window[2]

	emaPrev := emaWindow[0]
	emaCurrent := emaWindow[1]
	emaNext := emaWindow[2]

	//fmt.Println()
	//fmt.Println("prev:", prev.Low, prev.High, emaPrev.Value)
	//fmt.Println("curr:", current.Low, current.High, emaCurrent.Value)
	//fmt.Println("next:", next.Low, next.High, emaNext.Value)
	//t := time.Unix(current.Time, 0).Add(-3*time.Hour - 30*time.Minute)
	//fmt.Println("**** time:", t.Month(), t.Day(), t.Hour(), t.Minute())

	if (prev.Low < emaPrev.Value && prev.High > emaPrev.Value) ||
		(current.Low < emaCurrent.Value && current.High > emaCurrent.Value) ||
		(next.Low < emaNext.Value && next.High > emaNext.Value) {
		return true
	}
	return false
}

// window[3] is gap-candle; calculate sl from window[0] & window[1] & window[2]
func (s *Strategy) gapData(gap gapType, window []*market.Candle) (side trade.Side, sl float64) {
	if gap == risingGap {
		side = trade.Buy
		sl = min(window[0].Low, window[1].Low, window[2].Low, window[3].Low) * (1 - s.cfg.SlExtraSpace/100)
	} else {
		side = trade.Sell
		sl = max(window[0].High, window[1].High, window[2].High, window[3].High) * (1 + s.cfg.SlExtraSpace/100)
	}
	return side, sl
}

func exitGapArea(candles []*market.Candle, i int, side trade.Side) int {
	i++
	for j := 0; j < len(candles)-1; j++ {
		if (side == trade.Buy && candles[j].Low < candles[j+1].Low) || (side == trade.Sell && candles[j].High > candles[j+1].High) {
			i++
		} else {
			break
		}
	}
	return i
}

func (s *Strategy) gapSignal(id int, time int64, price float64, ch chan<- *backtest.Marker) {
	ch <- marker(gapMarker, id, time, price)
}

func (s *Strategy) enter(wg *sync.WaitGroup, id int, side trade.Side, candles []*market.Candle, sl float64, ch chan<- *backtest.Marker) {
	defer wg.Done()

	maxEntranceDistance := func(sl, entryPrice float64) bool {
		distance := 1 + s.cfg.MaxEntranceDistance/100
		low := sl
		high := entryPrice
		if high < low {
			low, high = high, low
		}

		if high > low*distance {
			return true
		}
		return false
	}

	var entryPrice float64
	if entryPrice = candles[0].Low; side == trade.Sell {
		entryPrice = candles[0].High
	}

	for i, candle := range candles[1:] {
		if !s.cfg.Session.In(candle.Time) {
			return
		}
		if maxEntranceDistance(sl, entryPrice) {
			ch <- marker(MaxEntranceMarker, id, candle.Time, entryPrice)
			return
		}
		if (side == trade.Buy && entryPrice > candle.Low) || (side == trade.Sell && entryPrice < candle.High) {
			markerType_ := buyMarker
			if side == trade.Sell {
				markerType_ = sellMarker
			}
			ch <- marker(markerType_, id, candle.Time, entryPrice)
			p := &trade.Position{
				Id:    id,
				Enter: entryPrice,
				Sl:    sl,
				Tp:    entryPrice + s.cfg.RR*(entryPrice-sl),
				Exit:  -1,
				Side:  side,
			}
			wg.Add(1)
			go s.exit(wg, id, p, candles[i+1:], ch)
			return
		}
		if side == trade.Buy {
			entryPrice = candle.Low
		} else {
			entryPrice = candle.High
		}
	}
}

func (s *Strategy) exit(wg *sync.WaitGroup, id int, p *trade.Position, candles []*market.Candle, ch chan<- *backtest.Marker) {
	defer wg.Done()
	for _, candle := range candles {
		if (p.Side == trade.Buy && p.Tp < candle.High) || (p.Side == trade.Sell && p.Tp > candle.Low) { // Tp
			ch <- marker(tpMarker, id, candle.Time, p.Tp)
			s.tpCnt++
			break
		} else if (p.Side == trade.Buy && p.Sl > candle.Low) || (p.Side == trade.Sell && p.Sl < candle.High) { // Sl
			ch <- marker(slMarker, id, candle.Time, p.Sl)
			s.slCnt++
			break
		}
	}
}
