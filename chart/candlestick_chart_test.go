package chart

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// sampleCandles is a deterministic random walk, hourly from 2026-01-01.
func sampleCandles(n int) []Candle {
	out := make([]Candle, n)
	price := 100.0
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range out {
		o := price
		c := o + math.Sin(float64(i)*0.7)*1.5 + math.Cos(float64(i)*0.13)*0.8
		hi, lo := max(o, c)+0.6+math.Abs(math.Sin(float64(i)))*0.5, min(o, c)-0.5-math.Abs(math.Cos(float64(i)))*0.4
		out[i] = Candle{Time: base.Add(time.Duration(i) * time.Hour), Open: o, High: hi, Low: lo, Close: c, Volume: 1000 + float64(i%17)*130}
		price = c
	}
	return out
}

func mktState0(r ChartResult) *mktState {
	return *ui.Local(r.Element, "mujica.mkt", func() *mktState { return nil })
}

func TestCandlestickChartDrawsOnlyVisible(t *testing.T) {
	candles := sampleCandles(100000)
	var view MarketView
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}})
	}, 900, 400)
	tt.Frame()
	st := mktState0(res)
	if st.drawn != 120 {
		t.Fatalf("drew %d candles of 100000", st.drawn)
	}
	tt.Key(0, ui.KeyTab)
	for range 40 {
		tt.Key(0, ui.KeyMinus)
	}
	if st.drawn > 900/3 || view.Count > 900/3 {
		t.Fatalf("zoomed out: drew %d, view %+v", st.drawn, view)
	}
}

func TestCandlestickChartKeyboardPanZoomReadout(t *testing.T) {
	candles := sampleCandles(500)
	view := MarketView{Count: 50}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}})
	}, 900, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if res.Point != 499 || !tt.HasText("2026-01-21 19:00") {
		t.Fatalf("End: point %d texts %q", res.Point, tt.Texts())
	}
	tt.Key(ui.Shift, ui.KeyLeft)
	if view.End != 495 || view.Count != 50 {
		t.Fatalf("Shift+Left: %+v", view)
	}
	tt.Key(0, ui.KeyEqual)
	if view.Count != 40 {
		t.Fatalf("zoom in: %+v", view)
	}
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	if f, _ := marketWindow(view, 500, 300); f != 0 {
		t.Fatalf("Home did not bring the first candle into view: %+v", view)
	}
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.Contains(s, "O 100.00") && strings.Contains(s, "Vol 1,000")
	}
	if !found {
		t.Fatalf("readout %q", tt.Texts())
	}
}

func TestCandlestickChartDragPansAndRejects(t *testing.T) {
	candles := sampleCandles(300)
	view := MarketView{Count: 100}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}}).Element.Width(860)
	}, 900, 400)
	tt.Frame()
	tt.Press(400, 150)
	tt.Move(560, 150) // about 20 candles of 8 DIPs
	tt.Release(560, 150)
	if view.End < 275 || view.End > 285 {
		t.Fatalf("drag right shows older candles: %+v", view)
	}
	for name, f := range map[string]func(c *ui.Context){
		"type": func(c *ui.Context) {
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "x"}, Type: 9})
		},
		"order": func(c *ui.Context) {
			bad := []Candle{candles[1], candles[0]}
			CandlestickChart(c, bad, &view, CandlestickOptions{Frame: ChartFrame{Label: "x"}})
		},
		"prices": func(c *ui.Context) {
			bad := []Candle{{Time: candles[0].Time, Open: 5, High: 4, Low: 3, Close: 4}}
			CandlestickChart(c, bad, &view, CandlestickOptions{Frame: ChartFrame{Label: "x"}})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); f(c) }, 300, 200)
		}()
	}
}

func TestCandlestickChartTrendColors(t *testing.T) {
	candles := sampleCandles(30)
	view := MarketView{Count: 30}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		core.UseTrend(c, core.RedUp)
		res = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "沪深300"}})
	}, 600, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if res.Point != 0 || !tt.HasText("开 100.00") {
		t.Fatalf("zh readout %q", tt.Texts())
	}
}

func TestCandlestickChartSnapshot(t *testing.T) {
	candles := sampleCandles(80)
	view := MarketView{Count: 60}
	snapshot(t, "candlestick_chart", 640, 320, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC", Height: 260}})
		})
	})
}
