package chart_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func ExampleCandlestickChart() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleOHLCChart() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.OHLCChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleHeikinAshiChart() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	raw := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.HeikinAshiChart(c, candles, &window, &raw, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleLineQuoteChart() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	points := []chart.QuotePoint{{Time: at, Price: 101, Volume: 100}, {Time: at.Add(time.Hour), Price: 103, Volume: 200}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.LineQuoteChart(c, points, chart.LineQuoteChartOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, PrevClose: 100, RegularOpen: at.Add(30 * time.Minute), RegularClose: at.Add(45 * time.Minute)})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleVolumeChart() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.VolumeChart(c, candles, &window, chart.VolumeChartOptions{Frame: chart.ChartFrame{Label: "Volume"}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleDepthChart() {
	bids := []chart.DepthLevel{{Price: 99, Size: 10}, {Price: 98, Size: 20}}
	asks := []chart.DepthLevel{{Price: 101, Size: 15}, {Price: 102, Size: 25}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.DepthChart(c, bids, asks, chart.DepthChartOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleVolumeProfile() {
	bins := []chart.ProfileBin{{Low: 99, High: 100, Volume: 30}, {Low: 100, High: 101, Volume: 80}, {Low: 101, High: 102, Volume: 40}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.VolumeProfile(c, bins, chart.VolumeProfileOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleMultiChartLayout() {
	state := chart.MultiChartState{Count: 2}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.MultiChartLayout(c, &state, chart.MultiChartLayoutOptions{Height: 200}, func(i int, cell *chart.ChartCell) {
			input.TextInput(c, &cell.Symbol, input.TextInputOptions{Label: "Symbol"})
		})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleChartSync() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	var sync chart.ChartSync
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		sync.Link(c, "price")
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Sync: &sync, SyncID: "price"})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleTimeRangeSelector() {
	value := chart.MarketRange{Preset: chart.Range1M}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.TimeRangeSelector(c, &value, chart.TimeRangeSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleIntervalSelector() {
	value := chart.Interval("1h")
	favorites := []chart.Interval{"1m", "1h", "1D"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.IntervalSelector(c, &value, &favorites, chart.IntervalSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleChartTypeSwitcher() {
	value := chart.ChartCandles
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.ChartTypeSwitcher(c, &value, chart.ChartTypeSwitcherOptions{})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleIndicatorOverlay() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	overlays := []chart.IndicatorOverlay{{Kind: chart.IndicatorMA, Params: []float64{2}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Overlays: overlays})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleIndicatorPane() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	panes := chart.IndicatorPanes{Panes: []chart.IndicatorPaneSpec{{Kind: chart.IndicatorRSI, Params: []float64{2}}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.IndicatorPane(c, candles, &window, &panes, chart.IndicatorPaneOptions{Height: 200, Label: "BTC"})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleIndicatorSelector() {
	var overlays []chart.IndicatorOverlay
	var panes chart.IndicatorPanes
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.IndicatorSelector(c, &overlays, &panes, chart.IndicatorSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleDrawingTools() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	drawings := chart.DrawingTools{Tool: chart.ToolTrendLine}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Drawings: &drawings})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleDrawingToolbar() {
	var drawings chart.DrawingTools
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.DrawingToolbar(c, &drawings, chart.DrawingToolbarOptions{Horizontal: true})
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleCompareSymbol() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	symbols := []chart.CompareSymbol{{Name: "ETH", Closes: []float64{20, 21, 23}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CompareLegend(c, &symbols)
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Compare: &symbols})
	}
	ui.Render(view, 640, 400, 1)
}

func ExamplePriceAlertLine() {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	alerts := []chart.PriceAlertLine{{ID: "take-profit", Price: 105}}
	triggered := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Alerts: &alerts, OnAlert: func(alert chart.PriceAlertLine, last chart.Candle) {
			triggered = alert.ID
		}})
		ui.Text(c, triggered)
	}
	ui.Render(view, 640, 400, 1)
}

func ExampleCandleCountdown() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandleCountdown(c, "1h", chart.CandleCountdownOptions{Location: time.UTC})
	}
	ui.Render(view, 640, 400, 1)
}
