package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestPriceAlertLineDragAndTrigger(t *testing.T) {
	candles := sampleCandles(60)
	view := MarketView{Count: 60}
	alerts := []PriceAlertLine{{ID: "a", Price: 100, Label: "Buy"}}
	var fired []string
	var res ChartResult
	_ = res
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}, Alerts: &alerts,
			OnAlert: func(a PriceAlertLine, last Candle) { fired = append(fired, a.ID) }})
	}, 800, 400)
	if !tt.HasText("BTC") {
		t.Fatal("not built")
	}
	// Drag the line from its own y, as painted.
	st := mktState0(res)
	plot := st.plot
	y := st.ys.Map(100)
	tt.Press(plot.X+100, y)
	tt.Move(plot.X+100, y-40)
	tt.Release(plot.X+100, y-40)
	if alerts[0].Price <= 100 || view.End != 0 {
		t.Fatalf("drag did not move the alert up, or panned: %+v %+v", alerts, view)
	}
	// A new close crossing the alert fires it once.
	alerts[0].Price = candles[59].Close + 0.5
	next := candles[59]
	next.Time = next.Time.Add(3600e9)
	next.Open, next.Close, next.High, next.Low = next.Close, next.Close+1, next.Close+1.2, next.Close-0.1
	candles = append(candles, next)
	tt.Frame()
	tt.Frame()
	if len(fired) != 1 || fired[0] != "a" {
		t.Fatalf("fired %v", fired)
	}
}

func TestPriceAlertLineSnapshot(t *testing.T) {
	candles := sampleCandles(60)
	view := MarketView{Count: 60}
	alerts := []PriceAlertLine{{ID: "a", Price: 104, Label: "Take profit"}, {ID: "b", Price: 96}}
	snapshot(t, "price_alert_line", 640, 320, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC", Height: 260}, Alerts: &alerts})
		})
	})
}
