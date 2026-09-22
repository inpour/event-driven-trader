package csv

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/inpour/event-driven-trader/shared/market"
)

type Provider struct {
	path string
}

func NewProvider(path string) *Provider {
	return &Provider{
		path: path,
	}
}

func (p *Provider) Candles() ([]*market.Candle, error) {
	file, err := os.Open(p.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var candles []*market.Candle

	if len(rows) > 10000 {
		rows = rows[len(rows)-10000:]
	}

	for _, row := range rows {
		open, _ := strconv.ParseFloat(row[2], 64)
		high, _ := strconv.ParseFloat(row[3], 64)
		low, _ := strconv.ParseFloat(row[4], 64)
		close_, _ := strconv.ParseFloat(row[5], 64)
		volume, _ := strconv.ParseFloat(row[6], 64)

		unixTime, err := timeToUnixTime(row[0] + " " + row[1] + ":00")
		if err != nil {
			return nil, err
		}
		candles = append(candles, &market.Candle{
			Time:   unixTime,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close_,
			Volume: volume,
		})
	}

	return candles, nil
}

func timeToUnixTime(date string) (int64, error) {
	t, err := time.Parse("2006.01.02 15:04:05", date)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
