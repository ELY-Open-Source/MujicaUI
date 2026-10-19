package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestRiskMeterLevelsAndText(t *testing.T) {
	level := 4
	levels := []RiskLevel{{Label: "Conservative", Description: "Capital first."}, {Label: "Balanced"}, {Label: "Aggressive", Description: "Losses can exceed half."}}
	custom := false
	tt := tradingTester(core.GreenUp, 400, 160, func(c *ui.Context) {
		if custom {
			RiskMeter(c, level, RiskMeterOptions{Levels: levels})
		} else {
			RiskMeter(c, level, RiskMeterOptions{})
		}
	})
	if !tt.HasText("Moderately high") || !tt.HasText("4 / 5") {
		t.Fatalf("default: %q", tt.Texts())
	}
	custom, level = true, 3
	tt.Frame()
	if !tt.HasText("Aggressive") || !tt.HasText("Losses can exceed half.") || !tt.HasText("3 / 3") {
		t.Fatalf("custom: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("level 0 accepted")
		}
	}()
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) { RiskMeter(c, 0, RiskMeterOptions{}) })
}

func TestRiskMeterSnapshot(t *testing.T) {
	snapshot(t, "risk_meter", 320, 110, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			RiskMeter(c, 4, RiskMeterOptions{Levels: []RiskLevel{{Label: "Low"}, {Label: "Medium"}, {Label: "Elevated"}, {Label: "High", Description: "Suited to experienced traders."}, {Label: "Extreme"}}})
		})
	})
}
