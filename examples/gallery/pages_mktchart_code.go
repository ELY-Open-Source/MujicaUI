package main

const mktCode141 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode142 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.OHLCChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode143 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	raw := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.HeikinAshiChart(c, candles, &window, &raw, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode144 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	points := []mujica.QuotePoint{{Time: at, Price: 101, Volume: 100}, {Time: at.Add(time.Hour), Price: 103, Volume: 200}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.LineQuoteChart(c, points, mujica.LineQuoteChartOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, PrevClose: 100, RegularOpen: at.Add(30 * time.Minute), RegularClose: at.Add(45 * time.Minute)})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode145 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VolumeChart(c, candles, &window, mujica.VolumeChartOptions{Frame: mujica.ChartFrame{Label: "Volume"}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode146 = `	bids := []mujica.DepthLevel{{Price: 99, Size: 10}, {Price: 98, Size: 20}}
	asks := []mujica.DepthLevel{{Price: 101, Size: 15}, {Price: 102, Size: 25}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DepthChart(c, bids, asks, mujica.DepthChartOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode147 = `	bins := []mujica.ProfileBin{{Low: 99, High: 100, Volume: 30}, {Low: 100, High: 101, Volume: 80}, {Low: 101, High: 102, Volume: 40}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.VolumeProfile(c, bins, mujica.VolumeProfileOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode148 = `	state := mujica.MultiChartState{Count: 2}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MultiChartLayout(c, &state, mujica.MultiChartLayoutOptions{Height: 200}, func(i int, cell *mujica.ChartCell) {
			mujica.TextInput(c, &cell.Symbol, mujica.TextInputOptions{Label: "Symbol"})
		})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode149 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	var sync mujica.ChartSync
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		sync.Link(c, "price")
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Sync: &sync, SyncID: "price"})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode150 = `	value := mujica.MarketRange{Preset: mujica.Range1M}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TimeRangeSelector(c, &value, mujica.TimeRangeSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode151 = `	value := mujica.Interval("1h")
	favorites := []mujica.Interval{"1m", "1h", "1D"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.IntervalSelector(c, &value, &favorites, mujica.IntervalSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode152 = `	value := mujica.ChartCandles
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ChartTypeSwitcher(c, &value, mujica.ChartTypeSwitcherOptions{})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode153 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	overlays := []mujica.IndicatorOverlay{{Kind: mujica.IndicatorMA, Params: []float64{2}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Overlays: overlays})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode154 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	panes := mujica.IndicatorPanes{Panes: []mujica.IndicatorPaneSpec{{Kind: mujica.IndicatorRSI, Params: []float64{2}}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.IndicatorPane(c, candles, &window, &panes, mujica.IndicatorPaneOptions{Height: 200, Label: "BTC"})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode155 = `	var overlays []mujica.IndicatorOverlay
	var panes mujica.IndicatorPanes
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.IndicatorSelector(c, &overlays, &panes, mujica.IndicatorSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode156 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	drawings := mujica.DrawingTools{Tool: mujica.ToolTrendLine}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Drawings: &drawings})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode157 = `	var drawings mujica.DrawingTools
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DrawingToolbar(c, &drawings, mujica.DrawingToolbarOptions{Horizontal: true})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode158 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	symbols := []mujica.CompareSymbol{{Name: "ETH", Closes: []float64{20, 21, 23}}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CompareLegend(c, &symbols)
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Compare: &symbols})
	}
	ui.Render(view, 640, 400, 1)`

const mktCode159 = `	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	alerts := []mujica.PriceAlertLine{{ID: "take-profit", Price: 105}}
	triggered := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Alerts: &alerts, OnAlert: func(alert mujica.PriceAlertLine, last mujica.Candle) {
			triggered = alert.ID
		}})
		ui.Text(c, triggered)
	}
	ui.Render(view, 640, 400, 1)`

const mktCode160 = `	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandleCountdown(c, "1h", mujica.CandleCountdownOptions{Location: time.UTC})
	}
	ui.Render(view, 640, 400, 1)`
