package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSpreadIndicatorValueAndThreshold(t *testing.T) {
	bid, ask := 99.95, 100.05
	tt := tradingTester(core.GreenUp, 400, 100, func(c *ui.Context) {
		SpreadIndicator(c, bid, ask, SpreadIndicatorOptions{WarnPercent: 0.5})
	})
	if !tt.HasText("0.10 (0.100%)") || tt.HasText("Wide spread") {
		t.Fatalf("narrow: %q", tt.Texts())
	}
	ask = 101
	tt.Frame()
	if !tt.HasText("Wide spread") {
		t.Fatalf("wide: %q", tt.Texts())
	}
	ask = 99
	tt.Frame()
	if !tt.HasText("Crossed market") {
		t.Fatalf("crossed: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative threshold accepted")
		}
	}()
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) { SpreadIndicator(c, 1, 2, SpreadIndicatorOptions{WarnPercent: -1}) })
}

func TestSpreadIndicatorSnapshot(t *testing.T) {
	snapshot(t, "spread_indicator", 380, 100, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(10).Children(func() {
			SpreadIndicator(c, 99.95, 100.05, SpreadIndicatorOptions{WarnPercent: 0.5})
			SpreadIndicator(c, 99, 101, SpreadIndicatorOptions{WarnPercent: 0.5})
		})
	})
}
