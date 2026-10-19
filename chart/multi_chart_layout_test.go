package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestMultiChartLayoutKeepsCellsAndDisabled(t *testing.T) {
	state := MultiChartState{Count: 6, Active: 5}
	for i := range state.Cells {
		state.Cells[i] = ChartCell{Symbol: "BTC", Interval: "1h", View: MarketView{Count: 20 + i}}
	}
	disabled := false
	built := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		built = 0
		ui.Column(c).Disabled(disabled).Children(func() {
			MultiChartLayout(c, &state, MultiChartLayoutOptions{Height: 200}, func(_ int, cell *ChartCell) { built++; ui.Text(c, cell.Symbol) })
		})
	}, 700, 300)
	if built != 6 {
		t.Fatalf("built %d", built)
	}
	if err := tt.Click("6 charts"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if state.Count != 4 || state.Active != 3 || built != 4 {
		t.Fatalf("layout %+v built %d", state, built)
	}
	if state.Cells[5].View.Count != 25 {
		t.Fatal("hidden state lost")
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if state.Count != 6 || state.Cells[5].View.Count != 25 {
		t.Fatal("restored state lost")
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if state.Count != 6 {
		t.Fatal("disabled layout changed")
	}
}
func TestMultiChartLayoutSnapshot(t *testing.T) {
	state := MultiChartState{Count: 4}
	snapshot(t, "multi_chart_layout", 640, 360, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			MultiChartLayout(c, &state, MultiChartLayoutOptions{Height: 280}, func(i int, cell *ChartCell) { ui.Text(c, "BTC · 1h") })
		})
	})
}
