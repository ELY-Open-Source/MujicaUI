package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestMarketChartsStatesDisabledAndBoundaries(t *testing.T) {
	for _, name := range []string{"Candlestick", "OHLC", "HeikinAshi", "Quote", "Volume", "Depth", "Profile", "Overlay", "Compare", "Alerts", "Drawings"} {
		t.Run(name, func(t *testing.T) {
			cs := sampleCandles(30)
			v := MarketView{Count: 30}
			raw := false
			bids, asks := sampleBook()
			bins := VolumeProfileOf(cs, 12)
			quotes, _, _ := samplePoints()
			symbols := []CompareSymbol{{Name: "Other", Closes: closesOf(cs)}}
			alerts := []PriceAlertLine{{ID: "a", Price: 100}}
			drawings := DrawingTools{}
			status := data.DataReady
			disabled, empty := false, false
			var r ChartResult
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{Density: theme.Compact})
				ui.Column(c).Disabled(disabled).Width(600).Children(func() {
					candles := cs
					if empty {
						candles = nil
					}
					frame := ChartFrame{Label: "Market", Height: 180, Status: status, Error: "Archive unavailable"}
					opts := CandlestickOptions{Frame: frame}
					switch name {
					case "OHLC":
						r = OHLCChart(c, candles, &v, opts)
					case "HeikinAshi":
						r = HeikinAshiChart(c, candles, &v, &raw, opts)
					case "Quote":
						q := quotes
						if empty {
							q = nil
						}
						r = LineQuoteChart(c, q, LineQuoteChartOptions{Frame: frame, PrevClose: 100})
					case "Volume":
						r = VolumeChart(c, candles, &v, VolumeChartOptions{Frame: frame})
					case "Depth":
						b, a := bids, asks
						if empty {
							b = nil
							a = nil
						}
						r = DepthChart(c, b, a, DepthChartOptions{Frame: frame})
					case "Profile":
						b := bins
						if empty {
							b = nil
						}
						r = VolumeProfile(c, b, VolumeProfileOptions{Frame: frame})
					default:
						if name == "Overlay" {
							opts.Overlays = []IndicatorOverlay{{Kind: IndicatorMA}}
						}
						if name == "Compare" {
							opts.Compare = &symbols
						}
						if name == "Alerts" {
							opts.Alerts = &alerts
						}
						if name == "Drawings" {
							opts.Drawings = &drawings
						}
						r = CandlestickChart(c, candles, &v, opts)
					}
				})
			}, 620, 350)
			for range 10 {
				if tt.Focused("Market") {
					break
				}
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, ui.KeyHome)
			tt.Key(0, ui.KeyRight)
			if r.Point < 0 || len(tt.Announcements()) == 0 {
				t.Fatalf("keyboard path point %d", r.Point)
			}
			tt.Key(0, ui.KeyHome)
			first := r.Point
			tt.Key(0, ui.KeyLeft)
			if name != "Depth" && name != "Profile" && r.Point != first {
				t.Fatalf("left boundary %d", r.Point)
			}
			disabled = true
			before := v
			tt.Frame()
			tt.Key(0, ui.KeyEqual)
			tt.Scroll(200, 120, 0, 20)
			tt.ClickAt(200, 120)
			if v != before {
				t.Fatalf("disabled window changed: %+v", v)
			}
			disabled = false
			tt.Frame()
			height := r.Element.Bounds().H
			for _, st := range []data.DataStatus{data.DataLoading, data.DataFailed, data.DataReady} {
				status = st
				empty = st == data.DataReady
				tt.Frame()
				if r.Point != -1 || r.Element.Bounds().H != height {
					t.Fatalf("state %v: point %d height %v want %v", st, r.Point, r.Element.Bounds().H, height)
				}
			}
		})
	}
}
