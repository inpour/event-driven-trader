package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/inpour/event-driven-trader/shared/backtest"
	"github.com/inpour/event-driven-trader/shared/market"
)

type templateData struct {
	Data    template.JS
	Markers template.JS
	Title   string
}

//go:embed template.html
var chartTemplate string

type ChartRenderer struct {
	path string
}

func NewChartRenderer(path string) *ChartRenderer {
	return &ChartRenderer{path: path}
}

func (c *ChartRenderer) HTML(
	candles []*market.Candle,
	markers []*backtest.Marker,
	metrics backtest.Metrics,
	info backtest.StrategyInfo,
) error {

	jsonCandles, err := json.Marshal(candles)
	if err != nil {
		return err
	}
	jsonMarkers, err := json.Marshal(markers)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("chart").Parse(chartTemplate))

	chart, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer chart.Close()

	winRate, _ := metrics.WinRate()
	title := fmt.Sprintf("%s.%s, WR:%.1f%%, TP:%d, SL:%d; Config: %+v",
		info.Name, info.Version, winRate*100,
		metrics.TpCnt, metrics.SlCnt, info.Config)

	pageData := templateData{
		Data:    template.JS(jsonCandles),
		Markers: template.JS(jsonMarkers),
		Title:   title,
	}
	if err := tmpl.Execute(chart, pageData); err != nil {
		return err
	}

	return nil
}
