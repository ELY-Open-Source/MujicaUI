package finance

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func greeksRows() []GreeksRow {
	return []GreeksRow{
		{Label: "AAPL 230 C", Greeks: OptionGreeks{0.52, 0.031, -0.12, 0.21, 0.05}},
		{Label: "AAPL 220 P", Greeks: OptionGreeks{-0.38, 0.028, -0.10, 0.19, math.NaN()}},
	}
}

func TestGreeksTableValuesAndPrecision(t *testing.T) {
	tt := chatTester(640, 320, func(c *ui.Context) { GreeksTable(c, greeksRows(), GreeksTableOptions{}) })
	for _, s := range []string{"Delta", "Rho", "+0.5200", "−0.1200", "+0.0310", "—"} {
		if !tt.HasText(s) {
			t.Errorf("%q missing from %q", s, tt.Texts())
		}
	}
	tt = chatTester(640, 320, func(c *ui.Context) { GreeksTable(c, greeksRows(), GreeksTableOptions{Decimals: 2}) })
	if !tt.HasText("+0.52") || tt.HasText("+0.5200") {
		t.Fatalf("decimals: %q", tt.Texts())
	}
}

func TestGreeksTableNotes(t *testing.T) {
	tt := chatTester(640, 320, func(c *ui.Context) { GreeksTable(c, greeksRows(), GreeksTableOptions{}) })
	if tt.HasText("Price change per day as time passes.") {
		t.Fatal("notes shown without Notes")
	}
	tt = chatTester(640, 420, func(c *ui.Context) { GreeksTable(c, greeksRows(), GreeksTableOptions{Notes: true}) })
	if !tt.HasText("Price change per day as time passes.") {
		t.Fatal("notes missing")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		GreeksTable(c, nil, GreeksTableOptions{Decimals: -1})
	}, 100, 100, 1)
}

func TestGreeksTableSnapshot(t *testing.T) {
	snapshot(t, "greeks_table", 640, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { GreeksTable(c, greeksRows(), GreeksTableOptions{Notes: true}) })
	})
}
