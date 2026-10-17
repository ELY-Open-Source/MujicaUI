package code

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

// ProblemSeverity ranks a Problem.
type ProblemSeverity int

const (
	ProblemError ProblemSeverity = iota
	ProblemWarning
	ProblemInfo
)

// Problem is one diagnostic.
type Problem struct {
	Severity        ProblemSeverity
	File            string
	Line, Col       int
	Message, Source string
}

// ProblemsPanelState is a ProblemsPanel's filters, folded files and place.
type ProblemsPanelState struct {
	// Hidden severities are filtered out.
	Hidden [3]bool
	// Collapsed holds folded file paths.
	Collapsed ui.Selection[string]
	rows      ui.ListState
	selected  int
	flat      []problemRow
	data      *Problem
	n, folded int
	hidden    [3]bool
	started   bool
}

// problemRow is a file header or a problem.
type problemRow struct {
	file    string
	problem int
}

// ProblemsPanelOptions configure a ProblemsPanel.
type ProblemsPanelOptions struct {
	Disabled bool
}

// ProblemsPanelResult carries the problem opened.
type ProblemsPanelResult struct {
	Element ui.Element
	Opened  *Problem
}

// ProblemsPanel lists diagnostics grouped by file.
func ProblemsPanel(c *ui.Context, problems []Problem, s *ProblemsPanelState, opts ProblemsPanelOptions) ProblemsPanelResult {
	if s == nil {
		panic("mujica: ProblemsPanel needs a state")
	}
	if !s.started {
		s.started, s.selected = true, -1
	}
	var counts [3]int
	for _, p := range problems {
		if p.Severity < ProblemError || p.Severity > ProblemInfo {
			panic("mujica: Problem severity out of range")
		}
		counts[p.Severity]++
	}
	var res ProblemsPanelResult
	res.Element = CodePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		CodeHeader(c, func() {
			for sv := ProblemError; sv <= ProblemInfo; sv++ {
				on := !s.Hidden[sv]
				name := fmt.Sprintf(core.Msg(c, problemsKeys[sv]), counts[sv])
				if input.Toggle(c, &on, name, input.ToggleOptions{Icon: icons.Must(problemsIcons[sv]), Disabled: opts.Disabled}).Changed() {
					s.Hidden[sv] = !on
					slog.Debug("mujica: problems_panel filter", "severity", sv, "shown", on)
				}
			}
		})
		problemsGroup(problems, s)
		if len(s.flat) == 0 {
			CodeCaption(c, core.Msg(c, "problems_panel.empty")).Padding(core.Density(c).Unit() * 4)
			return
		}
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any { return s.flat[i] }
		s.rows.Label = func(i int) string {
			if r := s.flat[i]; r.problem >= 0 {
				return problems[r.problem].Message
			}
			return s.flat[i].file
		}
		list := CodeRows(c, &s.rows, len(s.flat), func(i int) { problemsRow(c, problems, s, s.flat[i], disabled) })
		if disabled || s.selected < 0 {
			return
		}
		r := s.flat[s.selected]
		switch {
		case list.Submitted() && r.problem >= 0:
			p := problems[r.problem]
			res.Opened = &p
			slog.Debug("mujica: problems_panel open", "line", p.Line)
		case list.Submitted(),
			r.problem < 0 && !s.Collapsed.Has(r.file) && list.Shortcut(0, ui.KeyLeft),
			r.problem < 0 && s.Collapsed.Has(r.file) && list.Shortcut(0, ui.KeyRight):
			problemsFold(s, r.file)
		}
	})
	return res
}

var (
	problemsKeys  = [3]string{"problems_panel.errors", "problems_panel.warnings", "problems_panel.infos"}
	problemsIcons = [3]string{"circle-x", "triangle-alert", "info"}
)

func problemsFold(s *ProblemsPanelState, file string) {
	if s.Collapsed.Has(file) {
		s.Collapsed.Remove(file)
	} else {
		s.Collapsed.Add(file)
	}
	s.folded = -1
	slog.Debug("mujica: problems_panel fold", "folded", s.Collapsed.Has(file))
}

// problemsGroup regroups when data or filters change.
func problemsGroup(problems []Problem, s *ProblemsPanelState) {
	var first *Problem
	if len(problems) > 0 {
		first = &problems[0]
	}
	if s.data == first && s.n == len(problems) && s.folded == s.Collapsed.Len() && s.hidden == s.Hidden {
		return
	}
	s.data, s.n, s.folded, s.hidden = first, len(problems), s.Collapsed.Len(), s.Hidden
	var idx []int
	for i, p := range problems {
		if !s.Hidden[p.Severity] {
			idx = append(idx, i)
		}
	}
	slices.SortStableFunc(idx, func(a, b int) int {
		pa, pb := problems[a], problems[b]
		return cmp.Or(cmp.Compare(pa.File, pb.File), cmp.Compare(pa.Line, pb.Line), cmp.Compare(pa.Col, pb.Col))
	})
	s.flat = s.flat[:0]
	for j, i := range idx {
		f := problems[i].File
		if j == 0 || problems[idx[j-1]].File != f {
			s.flat = append(s.flat, problemRow{f, -1})
		}
		if !s.Collapsed.Has(f) {
			s.flat = append(s.flat, problemRow{f, i})
		}
	}
	s.selected = min(s.selected, len(s.flat)-1)
}

func problemsRow(c *ui.Context, problems []Problem, s *ProblemsPanelState, r problemRow, disabled bool) {
	k := core.Tokens(c)
	if r.problem < 0 {
		chev, key := "chevron-down", "search_panel.fold"
		if s.Collapsed.Has(r.file) {
			chev, key = "chevron-right", "search_panel.unfold"
		}
		if core.IconAction(c, icons.Must(chev), fmt.Sprintf(core.Msg(c, key), r.file)).Clicked() && !disabled {
			problemsFold(s, r.file)
		}
		ui.Text(c, r.file).FontWeight(600).SingleLine().Grow(1).MinWidth(0)
		return
	}
	p := problems[r.problem]
	col := [3]ui.Color{k.Danger, k.Warning, k.Info}[p.Severity]
	ui.Icon(c, icons.Must(problemsIcons[p.Severity])).Size(16, 16).Shrink(0).TextColor(col).Margin(0, 0, 0, 28).
		Label(core.Msg(c, [3]string{"problems_panel.error", "problems_panel.warning", "problems_panel.info"}[p.Severity]))
	ui.Text(c, p.Message).SingleLine().Grow(1).MinWidth(0)
	if p.Source != "" {
		CodeCaption(c, p.Source).Shrink(0)
	}
	CodeCaption(c, fmt.Sprintf(core.Msg(c, "problems_panel.at"), p.Line, p.Col)).FontFeatures("tnum").Shrink(0)
}
