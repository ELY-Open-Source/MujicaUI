package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type mktDemoState struct {
	candles                                []chart.Candle
	points                                 []chart.QuotePoint
	bids, asks                             []chart.DepthLevel
	bins                                   []chart.ProfileBin
	window                                 chart.MarketView
	second                                 chart.MarketView
	sync                                   chart.ChartSync
	layout                                 chart.MultiChartState
	interval                               chart.Interval
	favorites                              []chart.Interval
	rangeValue                             chart.MarketRange
	chartType                              chart.ChartType
	overlays                               []chart.IndicatorOverlay
	panes                                  chart.IndicatorPanes
	drawings                               chart.DrawingTools
	compare                                []chart.CompareSymbol
	alerts                                 []chart.PriceAlertLine
	status                                 string
	disabled, raw, redUp, hidden, extended bool
	previous, fraction                     float64
}

func mktInitial() mktDemoState {
	s := mktDemoState{status: "Ready", interval: "1h", favorites: []chart.Interval{"1m", "1h", "1D"}, previous: 100, fraction: 0.7, extended: true}
	start := time.Date(2026, 10, 5, 8, 0, 0, 0, time.UTC)
	price := 100.0
	compare := make([]float64, 180)
	for i := range 180 {
		closeAt := price + math.Sin(float64(i)*0.7)*1.5 + math.Cos(float64(i)*0.13)*0.8
		s.candles = append(s.candles, chart.Candle{Time: start.Add(time.Duration(i) * time.Hour), Open: price, High: max(price, closeAt) + 0.8, Low: min(price, closeAt) - 0.7, Close: closeAt, Volume: 1000 + float64(i%17)*130})
		s.points = append(s.points, chart.QuotePoint{Time: start.Add(time.Duration(i) * 5 * time.Minute), Price: 100 + 2*math.Sin(float64(i)/10), Volume: 100 + float64(i%7)*20})
		compare[i] = 20 + float64(i)*0.08
		price = closeAt
	}
	for i := range 40 {
		s.bids = append(s.bids, chart.DepthLevel{Price: 99.9 - float64(i)*0.1, Size: float64(10 + i*3)})
		s.asks = append(s.asks, chart.DepthLevel{Price: 100.1 + float64(i)*0.1, Size: float64(8 + i*4)})
	}
	s.bins = chart.VolumeProfileOf(s.candles, 24)
	s.compare = []chart.CompareSymbol{{Name: "ETH", Closes: compare}}
	s.alerts = []chart.PriceAlertLine{{ID: "target", Price: 104, Label: "Target"}}
	s.overlays = []chart.IndicatorOverlay{{Kind: chart.IndicatorMA, Params: []float64{10}}}
	s.panes = chart.IndicatorPanes{Panes: []chart.IndicatorPaneSpec{{Kind: chart.IndicatorMACD}, {Kind: chart.IndicatorRSI}}}
	s.layout.Count = 2
	for i := range s.layout.Cells {
		s.layout.Cells[i] = chart.ChartCell{Symbol: fmt.Sprintf("Symbol %d", i+1), Interval: "1h"}
	}
	return s
}

func mktControls(c *ui.Context, s *mktDemoState, log func(string)) chart.ChartFrame {
	ui.Row(c).Gap(12).Wrap().Children(func() {
		if input.Select(c, &s.status, []input.SelectOption[string]{{Value: "Ready", Label: "Ready"}, {Value: "Loading", Label: "Loading"}, {Value: "Failed", Label: "Failed"}, {Value: "Empty", Label: "Empty"}}, input.SelectOptions{Label: "Data state"}).Changed() {
			log("Data state: " + s.status)
		}
		if input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}).Changed() {
			log(fmt.Sprintf("Disabled: %v", s.disabled))
		}
		if input.Checkbox(c, &s.redUp, "Red up", input.CheckboxOptions{}).Changed() {
			log(fmt.Sprintf("Red up: %v", s.redUp))
		}
	})
	trend := core.GreenUp
	if s.redUp {
		trend = core.RedUp
	}
	core.UseTrend(c, trend)
	f := chart.ChartFrame{Label: "BTC", Height: 220}
	switch s.status {
	case "Loading":
		f.Status = data.DataLoading
	case "Failed":
		f.Status = data.DataFailed
		f.Error = "Archive unavailable"
	}
	return f
}

func mktLog(log func(string), name string, r chart.ChartResult) {
	if r.Changed() {
		log(fmt.Sprintf("%s: point %d", name, r.Point))
	}
	if r.Retried() {
		log(name + ": retry")
	}
}
