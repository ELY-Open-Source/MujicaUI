package feedback

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
)

// EvalCase is one evaluation case and its outcome.
type EvalCase struct {
	ID     string
	Name   string
	Score  float64
	Passed bool
}

// EvalResultTableOptions configure an EvalResultTable.
type EvalResultTableOptions struct {
	// Decimals of scores; zero means 2.
	Decimals int
	Status   data.DataStatus
	Error    string
	// Version changes when cases are edited in place.
	Version int
}

// EvalResultTableResult is a built EvalResultTable and its events.
type EvalResultTableResult struct {
	Element *ui.Element
	Table   data.DataTableView
	diff    string
}

// Diff returns the case to show differences for.
func (r EvalResultTableResult) Diff() (string, bool) { return r.diff, r.diff != "" }

type evalSummary struct {
	src                *EvalCase
	n, version, passed int
	built              bool
}

// EvalResultTable lists scored cases with a pass summary.
func EvalResultTable(c *ui.Context, s *data.DataTableState[string], cases []EvalCase, opts EvalResultTableOptions) EvalResultTableResult {
	if opts.Decimals < 0 {
		panic("mujica: EvalResultTable Decimals is negative")
	}
	dec := opts.Decimals
	if dec == 0 {
		dec = 2
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2)
	res := EvalResultTableResult{Element: col}
	summary := ui.Local(col, "summary", func() evalSummary { return evalSummary{} })
	var first *EvalCase
	if len(cases) > 0 {
		first = &cases[0]
	}
	if !summary.built || summary.src != first || summary.n != len(cases) || summary.version != opts.Version {
		passed := 0
		seen := make(map[string]bool, len(cases))
		for _, ec := range cases {
			if ec.ID == "" || seen[ec.ID] || math.IsNaN(ec.Score) || math.IsInf(ec.Score, 0) {
				panic("mujica: EvalCase needs a unique ID and finite Score")
			}
			seen[ec.ID] = true
			if ec.Passed {
				passed++
			}
		}
		*summary = evalSummary{src: first, n: len(cases), version: opts.Version, passed: passed, built: true}
	}
	passed := summary.passed
	verdict := func(ec EvalCase) string {
		if ec.Passed {
			return core.Msg(c, "eval_result_table.pass")
		}
		return core.Msg(c, "eval_result_table.fail")
	}
	cols := []data.DataColumn[EvalCase]{
		{ID: "name", Title: core.Msg(c, "eval_result_table.case"), Sortable: true, Compare: func(a, b EvalCase) int { return cmp.Compare(a.Name, b.Name) }, Text: func(e EvalCase) string { return e.Name }},
		{ID: "score", Title: core.Msg(c, "eval_result_table.score"), Width: 90, Align: ui.End, Sortable: true, Compare: func(a, b EvalCase) int { return cmp.Compare(a.Score, b.Score) },
			Text: func(e EvalCase) string { return strconv.FormatFloat(e.Score, 'f', dec, 64) }},
		{ID: "result", Title: core.Msg(c, "eval_result_table.result"), Width: 110, Sortable: true,
			Compare: func(a, b EvalCase) int { return cmp.Compare(evalRank(a), evalRank(b)) }, Text: verdict,
			Cell: func(e EvalCase) {
				icon, col := "circle-check", k.Success
				if !e.Passed {
					icon, col = "circle-x", k.Danger
				}
				ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(col).Children(func() {
					ui.Icon(c, icons.Must(icon)).Size(14, 14)
					ui.Text(c, verdict(e)).FontWeight(600)
				})
			}},
		{ID: "diff", Title: "", Width: 48, Text: func(EvalCase) string { return "" },
			Cell: func(e EvalCase) {
				b := core.IconAction(c, icons.Must("arrow-left-right"), fmt.Sprintf(core.Msg(c, "eval_result_table.diff"), e.Name)).Tooltip(core.Msg(c, "eval_result_table.view"))
				if b.Clicked() && !b.IsDisabled() {
					res.diff = e.ID
					slog.Debug("mujica: eval result diff", "case", e.ID)
				}
			}},
	}
	col.Children(func() {
		if len(cases) > 0 {
			ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "eval_result_table.summary"), passed, len(cases))).FontWeight(600)
				ui.Box(c).Height(6).Width(160).Radius(1).Background(k.Danger.Mix(k.Surface, 0.6)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(passed) / float32(len(cases)), H: r.H}, k.Success, 1)
				})
			})
		}
		res.Table = data.DataTable(c, s, cases, data.DataTableOptions[EvalCase, string]{Columns: cols, Key: func(e EvalCase) string { return e.ID },
			Empty: core.Msg(c, "eval_result_table.empty"), Status: opts.Status, Error: opts.Error, Version: opts.Version})
		res.Table.Element.Grow(1)
	})
	res.Table.Element = col
	return res
}

func evalRank(e EvalCase) int {
	if e.Passed {
		return 1
	}
	return 0
}
