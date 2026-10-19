package feedback

import (
	"math"
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
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
		core.Use(c, core.Settings{})
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
				core.Use(c, core.Settings{})
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

func TestCostBreakdownVersionAndToolGroups(t *testing.T) {
	var s CostBreakdownState
	items := costSample()
	version := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CostBreakdown(c, &s, items, CostBreakdownOptions{GroupBy: CostByTool, Version: version}).Element.Height(300)
	}, 800, 400)
	if len(s.rows) != 2 || s.rows[0].name != "search" || s.rows[0].calls != 2 {
		t.Fatalf("tool groups %+v", s.rows)
	}
	items[0].Cost = 1.30
	version++
	tt.Frame()
	if !tt.HasText("$2.50") {
		t.Fatalf("stale total: %q", tt.Texts())
	}
}

func TestCostBreakdownLargeData(t *testing.T) {
	items := make([]CostItem, 100000)
	for i := range items {
		items[i] = CostItem{Model: strconv.Itoa(i), Cost: 1}
	}
	var s CostBreakdownState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CostBreakdown(c, &s, items, CostBreakdownOptions{}).Element.Height(300)
	}, 800, 400)
	if s.total.cost != 100000 || s.Table.Rows() != 100000 || len(tt.Texts()) > 300 {
		t.Fatalf("rows %d texts %d total %v", s.Table.Rows(), len(tt.Texts()), s.total.cost)
	}
	a, b := s.Table.VisibleRows()
	if b-a > 30 {
		t.Fatalf("built %d rows", b-a)
	}
}

func TestCostSumOverflowPanics(t *testing.T) {
	for _, items := range [][]CostItem{
		{{InputTokens: math.MaxInt}, {InputTokens: 1}},
		{{OutputTokens: math.MaxInt}, {OutputTokens: 1}},
		{{Cost: math.MaxFloat64}, {Cost: math.MaxFloat64}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("unrepresentable cost total accepted")
				}
			}()
			costSum(items, CostByModel, "unnamed")
		}()
	}
}
