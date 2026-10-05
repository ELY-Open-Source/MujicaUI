package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

type mktDemoState struct {
	candles                                []mujica.Candle
	points                                 []mujica.QuotePoint
	bids, asks                             []mujica.DepthLevel
	bins                                   []mujica.ProfileBin
	window                                 mujica.MarketView
	second                                 mujica.MarketView
	sync                                   mujica.ChartSync
	layout                                 mujica.MultiChartState
	interval                               mujica.Interval
	favorites                              []mujica.Interval
	rangeValue                             mujica.MarketRange
	chartType                              mujica.ChartType
	overlays                               []mujica.IndicatorOverlay
	panes                                  mujica.IndicatorPanes
	drawings                               mujica.DrawingTools
	compare                                []mujica.CompareSymbol
	alerts                                 []mujica.PriceAlertLine
	status                                 string
	disabled, raw, redUp, hidden, extended bool
	previous, fraction                     float64
}

func mktInitial() mktDemoState {
	s := mktDemoState{status: "Ready", interval: "1h", favorites: []mujica.Interval{"1m", "1h", "1D"}, previous: 100, fraction: 0.7, extended: true}
	start := time.Date(2026, 10, 5, 8, 0, 0, 0, time.UTC)
	price := 100.0
	compare := make([]float64, 180)
	for i := range 180 {
		closeAt := price + math.Sin(float64(i)*0.7)*1.5 + math.Cos(float64(i)*0.13)*0.8
		s.candles = append(s.candles, mujica.Candle{Time: start.Add(time.Duration(i) * time.Hour), Open: price, High: max(price, closeAt) + 0.8, Low: min(price, closeAt) - 0.7, Close: closeAt, Volume: 1000 + float64(i%17)*130})
		s.points = append(s.points, mujica.QuotePoint{Time: start.Add(time.Duration(i) * 5 * time.Minute), Price: 100 + 2*math.Sin(float64(i)/10), Volume: 100 + float64(i%7)*20})
		compare[i] = 20 + float64(i)*0.08
		price = closeAt
	}
	for i := range 40 {
		s.bids = append(s.bids, mujica.DepthLevel{Price: 99.9 - float64(i)*0.1, Size: float64(10 + i*3)})
		s.asks = append(s.asks, mujica.DepthLevel{Price: 100.1 + float64(i)*0.1, Size: float64(8 + i*4)})
	}
	s.bins = mujica.VolumeProfileOf(s.candles, 24)
	s.compare = []mujica.CompareSymbol{{Name: "ETH", Closes: compare}}
	s.alerts = []mujica.PriceAlertLine{{ID: "target", Price: 104, Label: "Target"}}
	s.overlays = []mujica.IndicatorOverlay{{Kind: mujica.IndicatorMA, Params: []float64{10}}}
	s.panes = mujica.IndicatorPanes{Panes: []mujica.IndicatorPaneSpec{{Kind: mujica.IndicatorMACD}, {Kind: mujica.IndicatorRSI}}}
	s.layout.Count = 2
	for i := range s.layout.Cells {
		s.layout.Cells[i] = mujica.ChartCell{Symbol: fmt.Sprintf("Symbol %d", i+1), Interval: "1h"}
	}
	return s
}

func mktControls(c *ui.Context, s *mktDemoState, log func(string)) mujica.ChartFrame {
	ui.Row(c).Gap(12).Wrap().Children(func() {
		if mujica.Select(c, &s.status, []mujica.SelectOption[string]{{Value: "Ready", Label: "Ready"}, {Value: "Loading", Label: "Loading"}, {Value: "Failed", Label: "Failed"}, {Value: "Empty", Label: "Empty"}}, mujica.SelectOptions{Label: "Data state"}).Changed() {
			log("Data state: " + s.status)
		}
		if mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}).Changed() {
			log(fmt.Sprintf("Disabled: %v", s.disabled))
		}
		if mujica.Checkbox(c, &s.redUp, "Red up", mujica.CheckboxOptions{}).Changed() {
			log(fmt.Sprintf("Red up: %v", s.redUp))
		}
	})
	trend := mujica.GreenUp
	if s.redUp {
		trend = mujica.RedUp
	}
	mujica.UseTrend(c, trend)
	f := mujica.ChartFrame{Label: "BTC", Height: 220}
	switch s.status {
	case "Loading":
		f.Status = mujica.DataLoading
	case "Failed":
		f.Status = mujica.DataFailed
		f.Error = "Archive unavailable"
	}
	return f
}

func mktLog(log func(string), name string, r mujica.ChartResult) {
	if r.Changed() {
		log(fmt.Sprintf("%s: point %d", name, r.Point))
	}
	if r.Retried() {
		log(name + ": retry")
	}
}
