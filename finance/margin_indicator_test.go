package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestMarginIndicatorStates(t *testing.T) {
	ratio := 0.42
	tt := tradingTester(core.GreenUp, 400, 160, func(c *ui.Context) {
		ui.Box(c).Width(300).Children(func() { MarginIndicator(c, ratio, MarginIndicatorOptions{}) })
	})
	for _, tc := range []struct {
		ratio float64
		want  string
	}{{0.42, "Safe"}, {0.8, "Margin call"}, {1.2, "Liquidation"}} {
		ratio = tc.ratio
		tt.Frame()
		if !tt.HasText(tc.want) {
			t.Fatalf("%v: %q", tc.ratio, tt.Texts())
		}
	}
	if !tt.HasText("Warning line 80.0%") || !tt.HasText("Liquidation line 100.0%") || !tt.HasText("120.0%") {
		t.Fatalf("lines: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("warning above liquidation accepted")
		}
	}()
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) { MarginIndicator(c, 1, MarginIndicatorOptions{Warning: 2, Liquidation: 1}) })
}

func TestMarginIndicatorSnapshot(t *testing.T) {
	snapshot(t, "margin_indicator", 340, 120, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { MarginIndicator(c, 0.86, MarginIndicatorOptions{}) })
	})
}
