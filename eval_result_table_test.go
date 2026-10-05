package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func evalSample() []EvalCase {
	return []EvalCase{
		{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true},
		{ID: "c2", Name: "Refuses forgery", Score: 0.41},
		{ID: "c3", Name: "Summarizes letter", Score: 0.87, Passed: true},
	}
}

func TestEvalResultTableSummaryAndDiff(t *testing.T) {
	var s DataTableState[string]
	diff := ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
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
				Use(c, Settings{})
				var s DataTableState[string]
				EvalResultTable(c, &s, cases, EvalResultTableOptions{})
			}, 300, 200, 1)
		}()
	}
}

func TestEvalResultTableSnapshot(t *testing.T) {
	cases := evalSample()
	snapshot(t, "eval_result_table", 640, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *DataTableState[string] { return &DataTableState[string]{} })
		ui.Column(c).Padding(16).Height(300).Children(func() {
			EvalResultTable(c, *s, cases, EvalResultTableOptions{}).Table.Element.Grow(1)
		})
	})
}
