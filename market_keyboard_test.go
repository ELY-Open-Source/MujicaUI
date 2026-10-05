package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestDrawingToolsKeyboardCreateAndSelect(t *testing.T) {
	cs := sampleCandles(30)
	v := MarketView{Count: 30}
	d := DrawingTools{Tool: ToolTrendLine}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Drawing"}, Drawings: &d})
	}, 640, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if len(d.Items) != 1 || d.Items[0].A.Time != cs[0].Time || d.Items[0].B.Time != cs[1].Time {
		t.Fatalf("keyboard drawing %+v", d)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Key(ui.Alt, ui.KeyEnter)
	if d.Selected() != 0 {
		t.Fatalf("selection %d", d.Selected())
	}
	tt.Key(0, ui.KeyDelete)
	if len(d.Items) != 0 {
		t.Fatal("keyboard deletion failed")
	}
}

func TestPriceAlertLineKeyboard(t *testing.T) {
	cs := sampleCandles(30)
	v := MarketView{Count: 30}
	alerts := []PriceAlertLine{{ID: "a", Price: 100}, {ID: "b", Price: 105}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Alerts"}, Alerts: &alerts})
	}, 640, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyUp)
	if alerts[0].Price != 100.01 || alerts[1].Price != 105 {
		t.Fatalf("first alert %+v", alerts)
	}
	tt.Key(0, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyDown)
	if alerts[1].Price != 104.99 {
		t.Fatalf("second alert %+v", alerts)
	}
}
