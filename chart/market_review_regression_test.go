package chart

import (
	"image"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMarketSyncRange(t *testing.T) {
	cs := sampleCandles(300)
	short := cs[:200]
	s := ChartSync{Count: 40, End: cs[249].Time}
	var a, b MarketView
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			CandlestickChart(c, cs, &a, CandlestickOptions{Frame: ChartFrame{Label: "A", Height: 100}, Sync: &s, SyncID: "a"})
			CandlestickChart(c, short, &b, CandlestickOptions{Frame: ChartFrame{Label: "B", Height: 100}, Sync: &s, SyncID: "b"})
		})
	}, 640, 400)
	t.Logf("A=%+v B=%+v sync.End=%v", a, b, s.End)
	if a.End != 250 || s.End != cs[249].Time {
		t.Error("B erased A's explicit historical range")
	}
}
func TestMarketPaneKeyboardGeometry(t *testing.T) {
	cs := sampleCandles(300)
	v := MarketView{Count: 40}
	panes := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorOBV}}}
	var r IndicatorPaneResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r = IndicatorPane(c, cs, &v, &panes, IndicatorPaneOptions{Label: "BTC", Height: 220})
	}, 640, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	from, to := marketWindow(v, 300, 192)
	t.Logf("Point=%d visible=%d..%d view=%+v", r.Point, from, to, v)
	if r.Point < from || r.Point >= to {
		t.Error("focused readout has no visible crosshair")
	}
}
func TestMarketAlertsWithDrawingTools(t *testing.T) {
	cs := sampleCandles(30)
	v := MarketView{}
	a := []PriceAlertLine{{ID: "a", Price: 100}}
	d := DrawingTools{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "A"}, Alerts: &a, Drawings: &d})
	}, 640, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyUp)
	t.Logf("alert=%v", a[0].Price)
	if a[0].Price != 100.01 {
		t.Error("alert keyboard adjustment inaccessible with idle drawing tools")
	}
}
func TestMarketDrawingSwitchCancelsPending(t *testing.T) {
	cs := sampleCandles(30)
	v := MarketView{}
	d := DrawingTools{Tool: ToolTrendLine}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Children(func() {
			DrawingToolbar(c, &d, DrawingToolbarOptions{})
			CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "A"}, Drawings: &d}).Element.Grow(1)
		})
	}, 640, 450)
	tt.ClickAt(250, 100)
	if d.pending == nil {
		t.Fatal("setup missing pending point")
	}
	tt.Key(ui.Alt, ui.Key1)
	t.Logf("Tool=%d pending=%+v", d.Tool, d.pending)
	before := v
	tt.Press(300, 150)
	tt.Move(400, 150)
	tt.Release(400, 150)
	t.Logf("before=%+v after=%+v pending=%+v", before, v, d.pending)
	if d.Tool == ToolSelect && d.pending != nil {
		t.Error("Select retains previous tool preview and captures pan input")
	}
}
func TestMarketCountdownDST(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 3, 8, 23, 30, 0, 0, loc)
	close := Interval("1h").Close(now, loc)
	t.Logf("now=%s closes=%s", now, close)
	if close.Hour() != 0 {
		t.Error("hourly candle closes off the wall-clock hour")
	}
}
func TestMarketPaneResizeGeometry(t *testing.T) {
	cs := sampleCandles(300)
	v := MarketView{Count: 40}
	p := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorMACD}, {Kind: IndicatorRSI}}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IndicatorPane(c, cs, &v, &p, IndicatorPaneOptions{Label: "BTC", Height: 320})
	}, 640, 400)
	for i := 0; i < 12; i++ {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Resize MACD(12,26,9) and RSI(14)") {
			break
		}
	}
	t.Logf("texts=%q", tt.Texts())
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	top, _ := tt.Find("MACD(12,26,9)")
	plot, _ := tt.Find("BTC MACD(12,26,9)")
	t.Logf("sizes=%v top=%+v plot=%+v", p.Sizes, top, plot)
	if plot.H == 0 {
		t.Error("resizing to advertised minimum completely hides the chart")
	}
}
func TestMarketMissingComparisonCache(t *testing.T) {
	cs := sampleCandles(100000)
	v := MarketView{}
	s := []CompareSymbol{{Name: "missing", Closes: make([]float64, len(cs))}}
	for i := range s[0].Closes {
		s[0].Closes[i] = math.NaN()
	}
	var r ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r = CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "A"}, Compare: &s})
	}, 640, 320)
	memo := ui.Local(r.Element, "mujica.memo.compare.missing", func() ChartMemoEntry { return ChartMemoEntry{} })
	before := ChartID(memo.Val.([]float64))
	tt.Frame()
	memo = ui.Local(r.Element, "mujica.memo.compare.missing", func() ChartMemoEntry { return ChartMemoEntry{} })
	after := ChartID(memo.Val.([]float64))
	t.Logf("before=%s after=%s", before, after)
	if before != after {
		t.Error("unchanged missing comparison allocated and recomputed 100000 values")
	}
}

func TestMarketDailyAxisLocale(t *testing.T) {
	candles := sampleCandles(30)
	for i := range candles {
		candles[i].Time = candles[0].Time.AddDate(0, 0, i)
	}
	render := func(locale core.Locale) *image.RGBA {
		view := MarketView{}
		return ui.Render(func(c *ui.Context) {
			core.Use(c, core.Settings{Locale: locale})
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC", Height: 240}})
		}, 640, 320, 1)
	}
	en, zh := render(core.EnUS), render(core.ZhCN)
	for y := 214; y < 240; y++ {
		for x := 60; x < 620; x++ {
			if en.RGBAAt(x, y) != zh.RGBAAt(x, y) {
				return
			}
		}
	}
	t.Fatal("Chinese daily time axis retained English labels")
}

func TestMarketSyncDifferentIntervals(t *testing.T) {
	cs := sampleCandles(300)
	var slower []Candle
	for i := 0; i < len(cs); i += 2 {
		slower = append(slower, cs[i])
	}
	s := ChartSync{Count: 40, End: cs[199].Time}
	var a, b MarketView
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			CandlestickChart(c, cs, &a, CandlestickOptions{Frame: ChartFrame{Label: "A", Height: 100}, Sync: &s, SyncID: "a"})
			CandlestickChart(c, slower, &b, CandlestickOptions{Frame: ChartFrame{Label: "B", Height: 100}, Sync: &s, SyncID: "b"})
		})
	}, 640, 400)
	tt.Frame()
	from, to := marketWindow(b, len(slower), 200)
	if b.Count != 20 || slower[from].Time != cs[160].Time || slower[to-1].Time != cs[198].Time || s.End != cs[199].Time {
		t.Fatalf("intervals drifted: a=%+v b=%+v sync=%+v", a, b, s)
	}
}
