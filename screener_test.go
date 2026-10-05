package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestScreenerConditionsResultsSave(t *testing.T) {
	s := ScreenerState{Conditions: []ScreenerCondition{{Field: "pe", Op: ScreenerBelow, Value: 30}}}
	fields := []ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}}
	rows := []ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13}}}
	var got ScreenerView
	disabled := false
	tt := chatTester(800, 500, func(c *ui.Context) {
		v := Screener(c, &s, fields, rows, ScreenerOptions{Disabled: disabled})
		v.Element.Height(460)
		if v.Changed || v.Save {
			got = v
		}
	})
	if s.Table.Rows() != 1 || !tt.HasText("XOM") {
		t.Fatalf("rows %d", s.Table.Rows())
	}
	_ = tt.Click("Save conditions")
	if !got.Save {
		t.Fatal("save missing")
	}
	_ = tt.Click("Add condition")
	if len(s.Conditions) != 2 || !got.Changed {
		t.Fatalf("conditions %v", s.Conditions)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Add condition")
	if len(s.Conditions) != 2 {
		t.Fatal("disabled added")
	}
	for op := ScreenerAbove; op <= ScreenerEqual; op++ {
		want := []bool{false, true, false, true, true}[op]
		if got := ScreenerMatch(ScreenerRow{Values: map[string]float64{"x": 10}}, []ScreenerCondition{{Field: "x", Op: op, Value: 10}}); got != want {
			t.Fatalf("op %d %v", op, got)
		}
	}
	if ScreenerMatch(ScreenerRow{}, []ScreenerCondition{{Field: "x", Op: ScreenerAbove, Value: 0}}) {
		t.Fatal("missing passed")
	}
}
func TestScreenerSnapshot(t *testing.T) {
	snapshot(t, "screener", 820, 480, func(c *ui.Context) {
		sc := ui.Local(c.Root(), "snapshot.screener", func() ScreenerState {
			return ScreenerState{Conditions: []ScreenerCondition{{Field: "pe", Op: ScreenerBelow, Value: 30}}}
		})
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			Screener(c, sc, []ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}, {ID: "cap", Label: "Market cap (B)"}},
				[]ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34, "cap": 3000}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13, "cap": 450}}}, ScreenerOptions{}).Element.Height(440)

		})
	})
}
