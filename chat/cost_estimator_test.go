package chat

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCostEstimatorTable(t *testing.T) {
	tt := chatTester(420, 260, func(c *ui.Context) {
		CostEstimator(c, CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
	})
	for _, s := range []string{"12000", "3400", "15400", "$0.0360", "$0.0510", "$0.0870", "input $3, output $15"} {
		if !tt.HasText(s) {
			t.Errorf("%q missing from %q", s, tt.Texts())
		}
	}
	tt = chatTester(420, 260, func(c *ui.Context) {
		CostEstimator(c, CostEstimatorOptions{InputTokens: 1e6, InputPrice: 2.5, Currency: "¥", Digits: 2})
	})
	if !tt.HasText("¥2.50") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestEstimateCost(t *testing.T) {
	got := EstimateCost(1_000_000, 500_000, 3, 15)
	if got.Input != 3 || got.Output != 7.5 || got.Total != 10.5 {
		t.Fatalf("estimate %+v", got)
	}
	if z := EstimateCost(0, 0, 3, 15); z.Total != 0 {
		t.Fatalf("zero tokens cost %v", z.Total)
	}
}

func TestCostEstimatorRejectsBadInput(t *testing.T) {
	for _, opts := range []CostEstimatorOptions{{InputTokens: -1}, {InputPrice: -1}, {OutputPrice: math.NaN()}, {Digits: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", opts)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				CostEstimator(c, opts)
			}, 100, 100, 1)
		}()
	}
}

func TestCostEstimatorSnapshot(t *testing.T) {
	snapshot(t, "cost_estimator", 420, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CostEstimator(c, CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
		})
	})
}
