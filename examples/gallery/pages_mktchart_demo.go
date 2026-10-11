package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func mktDemo(c *ui.Context, log func(string), id string) {
	s := ui.Local(c.Root(), id, mktInitial)
	ui.Column(c).Gap(12).FillWidth().Children(func() {
		frame := mktControls(c, s, log)
		candles, points, bids, asks, bins := s.candles, s.points, s.bids, s.asks, s.bins
		if s.status == "Empty" {
			candles = nil
			points = nil
			bids = nil
			asks = nil
			bins = nil
		}
		ui.Column(c).Disabled(s.disabled).Gap(12).FillWidth().Children(func() {
			opts := chart.CandlestickOptions{Frame: frame, Type: s.chartType}
			draw := func() { mktLog(log, id, chart.CandlestickChart(c, candles, &s.window, opts)) }
			switch id {
			case "141":
				if chart.ChartTypeSwitcher(c, &s.chartType, chart.ChartTypeSwitcherOptions{}).Changed() {
					log(fmt.Sprintf("Chart type: %d", s.chartType))
				}
				opts.Type = s.chartType
				draw()
			case "142":
				mktLog(log, id, chart.OHLCChart(c, candles, &s.window, opts))
			case "143":
				before := s.raw
				mktLog(log, id, chart.HeikinAshiChart(c, candles, &s.window, &s.raw, opts))
				if before != s.raw {
					log(fmt.Sprintf("Raw: %v", s.raw))
				}
			case "144":
				if input.NumberInput(c, &s.previous, input.NumberInputOptions{Min: 80, Max: 120, Step: 0.5, Label: "Previous close"}).Changed() {
					log(fmt.Sprintf("Previous close: %.2f", s.previous))
				}
				if input.Checkbox(c, &s.extended, "Extended sessions", input.CheckboxOptions{}).Changed() {
					log(fmt.Sprintf("Extended sessions: %v", s.extended))
				}
				o := chart.LineQuoteChartOptions{Frame: frame, PrevClose: s.previous}
				if s.extended {
					o.RegularOpen = s.points[0].Time.Add(90 * time.Minute)
					o.RegularClose = s.points[0].Time.Add(8 * time.Hour)
				}
				mktLog(log, id, chart.LineQuoteChart(c, points, o))
			case "145":
				opts.Frame.Height = 150
				draw()
				frame.Label = "Volume"
				frame.Height = 100
				mktLog(log, id, chart.VolumeChart(c, candles, &s.window, chart.VolumeChartOptions{Frame: frame}))
			case "146":
				mktLog(log, id, chart.DepthChart(c, bids, asks, chart.DepthChartOptions{Frame: frame}))
			case "147":
				if input.NumberInput(c, &s.fraction, input.NumberInputOptions{Min: 0.1, Max: 1, Step: 0.1, Label: "Value area fraction"}).Changed() {
					log(fmt.Sprintf("Value area: %.1f", s.fraction))
				}
				mktLog(log, id, chart.VolumeProfile(c, bins, chart.VolumeProfileOptions{Frame: frame, ValueArea: s.fraction}))
			case "148":
				if chart.MultiChartLayout(c, &s.layout, chart.MultiChartLayoutOptions{Height: 420}, func(i int, cell *chart.ChartCell) {
					if input.TextInput(c, &cell.Symbol, input.TextInputOptions{Label: fmt.Sprintf("Symbol %d", i+1)}).Changed() {
						log(fmt.Sprintf("Cell %d: %s", i+1, cell.Symbol))
					}
					items := make([]input.SelectOption[chart.Interval], len(chart.Intervals))
					for j, iv := range chart.Intervals {
						items[j] = input.SelectOption[chart.Interval]{Value: iv, Label: string(iv)}
					}
					if input.Select(c, &cell.Interval, items, input.SelectOptions{Label: fmt.Sprintf("Interval %d", i+1)}).Changed() {
						log(fmt.Sprintf("Cell %d: %s", i+1, cell.Interval))
					}
					f := frame
					f.Label = cell.Symbol
					f.Height = 70
					mktLog(log, fmt.Sprint(i), chart.CandlestickChart(c, candles, &cell.View, chart.CandlestickOptions{Frame: f}))
				}).Changed() {
					log(fmt.Sprintf("Layout: %d, active %d", s.layout.Count, s.layout.Active))
				}
			case "149":
				for _, chartID := range []string{"price", "other"} {
					ui.Column(c).Key(chartID).Children(func() {
						if s.sync.Link(c, chartID).Changed() {
							log("Sync: " + chartID)
						}
						v := &s.window
						if chartID == "other" {
							v = &s.second
						}
						o := opts
						o.Frame.Height = 100
						o.Frame.Label = chartID
						o.Sync = &s.sync
						o.SyncID = chartID
						mktLog(log, chartID, chart.CandlestickChart(c, candles, v, o))
					})
				}
			case "150":
				if chart.TimeRangeSelector(c, &s.rangeValue, chart.TimeRangeSelectorOptions{}).Changed() {
					log(fmt.Sprintf("Range: %d", s.rangeValue.Preset))
				}
				from, to := s.rangeValue.Bounds(c.Now(), time.UTC)
				ui.Text(c, from.Format("2006-01-02")+" — "+to.Format("2006-01-02"))
			case "151":
				before := fmt.Sprint(s.favorites)
				if chart.IntervalSelector(c, &s.interval, &s.favorites, chart.IntervalSelectorOptions{}).Changed() {
					log("Interval: " + string(s.interval))
				}
				if before != fmt.Sprint(s.favorites) {
					log("Favorites: " + fmt.Sprint(s.favorites))
				}
				chart.CandleCountdown(c, s.interval, chart.CandleCountdownOptions{})
			case "152":
				if chart.ChartTypeSwitcher(c, &s.chartType, chart.ChartTypeSwitcherOptions{}).Changed() {
					log(fmt.Sprintf("Chart type: %d", s.chartType))
				}
				opts.Type = s.chartType
				draw()
			case "153":
				if chart.IndicatorSelector(c, &s.overlays, &s.panes, chart.IndicatorSelectorOptions{}).Changed() {
					log("Overlays changed")
				}
				opts.Overlays = s.overlays
				draw()
			case "154":
				if chart.IndicatorSelector(c, &s.overlays, &s.panes, chart.IndicatorSelectorOptions{}).Changed() {
					log("Indicators changed")
				}
				r := chart.IndicatorPane(c, candles, &s.window, &s.panes, chart.IndicatorPaneOptions{Label: "BTC", Height: 320, Status: frame.Status, Error: frame.Error})
				if r.PointChanged {
					log(fmt.Sprintf("Pane %d: point %d", r.Pane, r.Point))
				}
				if r.Removed >= 0 {
					log(fmt.Sprintf("Removed pane %d", r.Removed))
				}
				if r.Resized {
					log("Panes resized")
				}
			case "155":
				if chart.IndicatorSelector(c, &s.overlays, &s.panes, chart.IndicatorSelectorOptions{}).Changed() {
					log("Indicators changed")
				}
				opts.Overlays = s.overlays
				draw()
			case "156", "157":
				if input.TextInput(c, &s.drawings.Text, input.TextInputOptions{Label: "Annotation text"}).Changed() {
					log("Annotation text changed")
				}
				if chart.DrawingToolbar(c, &s.drawings, chart.DrawingToolbarOptions{Horizontal: true}).Changed() {
					log(fmt.Sprintf("Tool: %d, locked: %v", s.drawings.Tool, s.drawings.Locked))
				}
				before := fmt.Sprint(s.drawings.Items)
				opts.Drawings = &s.drawings
				draw()
				if before != fmt.Sprint(s.drawings.Items) {
					log(fmt.Sprintf("Drawings: %d", len(s.drawings.Items)))
				}
			case "158":
				if chart.CompareLegend(c, &s.compare).Changed() {
					log("Compared symbol removed")
				}
				if input.Button(c, "Restore ETH", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					s.compare = mktInitial().compare
					log("Restored ETH")
				}
				if len(candles) > 0 {
					opts.Compare = &s.compare
				}
				draw()
			case "159":
				before := s.alerts[0].Price
				opts.Alerts = &s.alerts
				opts.OnAlert = func(a chart.PriceAlertLine, last chart.Candle) {
					log(fmt.Sprintf("Alert %s crossed %.2f", a.ID, last.Close))
				}
				if input.Button(c, "Append crossing candle", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					last := s.candles[len(s.candles)-1]
					next := last
					next.Time = last.Time.Add(time.Hour)
					next.Open = last.Close
					next.Close = s.alerts[0].Price + 2
					if last.Close > s.alerts[0].Price {
						next.Close = s.alerts[0].Price - 2
					}
					next.High = max(next.Open, next.Close) + 1
					next.Low = min(next.Open, next.Close) - 1
					s.candles = append(s.candles[:len(s.candles):len(s.candles)], next)
					candles = s.candles
					log("Appended candle")
				}
				draw()
				if before != s.alerts[0].Price {
					log(fmt.Sprintf("Alert price: %.2f", s.alerts[0].Price))
				}
			case "160":
				if chart.IntervalSelector(c, &s.interval, &s.favorites, chart.IntervalSelectorOptions{}).Changed() {
					log("Interval: " + string(s.interval))
				}
				if input.Checkbox(c, &s.hidden, "Hide countdown", input.CheckboxOptions{}).Changed() {
					log(fmt.Sprintf("Hidden: %v", s.hidden))
				}
				if !s.hidden {
					chart.CandleCountdown(c, s.interval, chart.CandleCountdownOptions{})
				}
			}
		})
	})
}
