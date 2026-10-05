package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func costSample() []CostItem {
	return []CostItem{
		{Model: "opus", Tool: "search", Step: "plan", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
		{Model: "haiku", Tool: "search", Step: "gather", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
		{Model: "opus", Tool: "write", Step: "write", InputTokens: 2000, OutputTokens: 3000, Cost: 1.15},
	}
}

func TestCostSum(t *testing.T) {
	rows, total := costSum(costSample(), CostByModel, "?")
	if len(rows) != 2 || rows[0].name != "opus" || rows[0].calls != 2 || rows[0].input != 3000 || total.calls != 3 {
		t.Fatalf("rows %+v total %+v", rows, total)
	}
	if d := total.cost - 1.5; d > 1e-9 || d < -1e-9 {
		t.Fatalf("total cost %v", total.cost)
	}
	rows, _ = costSum([]CostItem{{Cost: 1}}, CostByTool, "(none)")
	if rows[0].name != "(none)" {
		t.Fatalf("unnamed group %q", rows[0].name)
	}
}

func TestCostBreakdownGroupsAndTotal(t *testing.T) {
	var s CostBreakdownState
	by := CostByModel
	items := costSample()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(400).Children(func() {
			CostBreakdown(c, &s, items, CostBreakdownOptions{GroupBy: by}).Element.Grow(1)
		})
	}, 800, 420)
	for _, want := range []string{"opus", "haiku", "$1.45", "$0.05", "97%", "Total", "$1.50", "7000 in · 3700 out"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if _, ok := tt.Find("Total $1.50 over 3 calls"); !ok {
		t.Error("total not named")
	}
	rows := &s.rows[0]
	tt.Frame()
	if &s.rows[0] != rows {
		t.Fatal("sums recomputed without a change")
	}
	by = CostByStep
	tt.Frame()
	if !tt.HasText("gather") || tt.HasText("haiku") {
		t.Fatalf("group by step: %q", tt.Texts())
	}
	if err := tt.Click("Cost"); err != nil {
		t.Fatal(err)
	}
	if s.Table.Sort.Column != "cost" {
		t.Fatalf("sort %+v", s.Table.Sort)
	}
}

func TestCostBreakdownInvalid(t *testing.T) {
	for _, tc := range []struct {
		items []CostItem
		opts  CostBreakdownOptions
	}{{[]CostItem{{Cost: -1}}, CostBreakdownOptions{}}, {nil, CostBreakdownOptions{GroupBy: 7}}, {nil, CostBreakdownOptions{Decimals: -1}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", tc)
				}
			}()
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				var s CostBreakdownState
				CostBreakdown(c, &s, tc.items, tc.opts)
			}, 300, 200, 1)
		}()
	}
}

func TestCostBreakdownSnapshot(t *testing.T) {
	items := costSample()
	snapshot(t, "cost_breakdown", 760, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *CostBreakdownState { return &CostBreakdownState{} })
		ui.Column(c).Padding(16).Height(300).Children(func() {
			CostBreakdown(c, *s, items, CostBreakdownOptions{}).Element.Grow(1)
		})
	})
}
