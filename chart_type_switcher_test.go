package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestChartTypeSwitcherKeyboardAndDisabled(t *testing.T) {
	value := ChartCandles
	disabled := false
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if ChartTypeSwitcher(c, &value, ChartTypeSwitcherOptions{Disabled: disabled}).Changed() {
			changes++
		}
	}, 800, 100)
	tt.Key(0, ui.KeyTab)
	for _, want := range []ChartType{ChartOHLC, ChartLine, ChartArea, ChartHeikinAshi} {
		tt.Key(0, ui.KeyRight)
		tt.Key(0, ui.KeySpace)
		if value != want {
			t.Fatalf("type %v, want %v", value, want)
		}
	}
	if changes != 4 {
		t.Fatalf("changes %d", changes)
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	if value != ChartHeikinAshi || changes != 4 {
		t.Fatal("disabled type changed")
	}
}
func TestChartTypeSwitcherSnapshot(t *testing.T) {
	value := ChartCandles
	snapshot(t, "chart_type_switcher", 680, 90, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() { ChartTypeSwitcher(c, &value, ChartTypeSwitcherOptions{}) })
	})
}
