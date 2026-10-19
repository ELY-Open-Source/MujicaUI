package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestChartSyncRangeCursorAndUnlink(t *testing.T) {
	cs := sampleCandles(300)
	sync := ChartSync{}
	v := MarketView{Count: 40, End: 200}
	sync.store("a", cs, &v, 150)
	var other MarketView
	sync.load("b", cs, &other)
	if other != v || sync.cursorIndex("b", cs) != 150 {
		t.Fatalf("sync %+v %+v", other, sync)
	}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Disabled(disabled).Children(func() { sync.Link(c, "b") })
	}, 300, 90)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if sync.linked("b") || sync.cursorIndex("b", cs) != -1 {
		t.Fatal("unlink failed")
	}
	other = MarketView{Count: 12}
	sync.load("b", cs, &other)
	if other.Count != 12 {
		t.Fatal("unlinked view overwritten")
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if sync.linked("b") {
		t.Fatal("disabled link changed")
	}
	sync.store("a", cs, &v, -1)
	if !sync.Cursor.IsZero() {
		t.Fatal("cursor remained after leaving owner")
	}
}
func TestChartSyncSnapshot(t *testing.T) {
	cs := sampleCandles(50)
	sync := ChartSync{Cursor: cs[25].Time, owner: "a"}
	v := MarketView{Count: 50}
	snapshot(t, "chart_sync", 640, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			sync.Link(c, "b")
			CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Linked", Height: 230}, Sync: &sync, SyncID: "b"})
		})
	})
}
