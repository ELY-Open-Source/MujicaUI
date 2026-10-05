package feedback

import (
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
)

func evalSample() []EvalCase {
	return []EvalCase{
		{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true},
		{ID: "c2", Name: "Refuses forgery", Score: 0.41},
		{ID: "c3", Name: "Summarizes letter", Score: 0.87, Passed: true},
	}
}

func TestEvalResultTableSummaryAndDiff(t *testing.T) {
	var s data.DataTableState[string]
	diff := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Height(400).Children(func() {
			r := EvalResultTable(c, &s, evalSample(), EvalResultTableOptions{})
			r.Table.Element.Grow(1)
			if id, ok := r.Diff(); ok {
				diff = id
			}
		})
	}, 700, 420)
	for _, want := range []string{"2 of 3 passed", "0.41", "Fail", "Pass"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q: %q", want, tt.Texts())
		}
	}
	if err := tt.Click("View differences for Refuses forgery"); err != nil {
		t.Fatal(err)
	}
	if diff != "c2" {
		t.Fatalf("diff %q", diff)
	}
	if err := tt.Click("Score"); err != nil {
		t.Fatal(err)
	}
	if s.Sort.Column != "score" {
		t.Fatalf("sort %+v", s.Sort)
	}
}

func TestEvalResultTableInvalid(t *testing.T) {
	for _, cases := range [][]EvalCase{{{Name: "x"}}, {{ID: "a", Score: nanValue()}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", cases)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				var s data.DataTableState[string]
				EvalResultTable(c, &s, cases, EvalResultTableOptions{})
			}, 300, 200, 1)
		}()
	}
}

func TestEvalResultTableSnapshot(t *testing.T) {
	cases := evalSample()
	snapshot(t, "eval_result_table", 640, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *data.DataTableState[string] { return &data.DataTableState[string]{} })
		ui.Column(c).Padding(16).Height(300).Children(func() {
			EvalResultTable(c, *s, cases, EvalResultTableOptions{}).Table.Element.Grow(1)
		})
	})
}

func TestEvalResultTableDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate case IDs accepted")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var state data.DataTableState[string]
		EvalResultTable(c, &state, []EvalCase{{ID: "a"}, {ID: "a"}}, EvalResultTableOptions{})
	}, 500, 300, 1)
}

func TestEvalResultTableDisabledKeyboardDiff(t *testing.T) {
	var state data.DataTableState[string]
	disabled := false
	diffs := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := EvalResultTable(c, &state, evalSample(), EvalResultTableOptions{})
			r.Table.Element.Height(300)
			if _, ok := r.Diff(); ok {
				diffs++
			}
		})
	}, 700, 400)
	tabTo(t, tt, "View differences for Greets by title")
	for i := 0; i < 8 && diffs == 0; i++ {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
	}
	if diffs != 1 {
		t.Fatalf("keyboard diffs %d", diffs)
	}
	disabled = true
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("View differences for Greets by title")
	if diffs != 1 {
		t.Fatalf("disabled diffs %d", diffs)
	}
}

func TestEvalResultTableLargeDataAndVersion(t *testing.T) {
	cases := make([]EvalCase, 100000)
	for i := range cases {
		cases[i] = EvalCase{ID: strconv.Itoa(i), Name: "Case " + strconv.Itoa(i), Score: 1, Passed: true}
	}
	var s data.DataTableState[string]
	version := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		EvalResultTable(c, &s, cases, EvalResultTableOptions{Version: version}).Element.Height(300)
	}, 700, 400)
	a, b := s.VisibleRows()
	if b-a > 30 || !tt.HasText("100000 of 100000 passed") {
		t.Fatalf("rows %d summary missing", b-a)
	}
	cases[0].Passed = false
	version++
	tt.Frame()
	if !tt.HasText("99999 of 100000 passed") {
		t.Fatal("stale summary")
	}
}
