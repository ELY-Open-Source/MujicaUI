package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// Process is one row of a ProcessList.
type Process struct {
	PID  int
	Name string
	// CPU is percent of one core; Memory is resident bytes.
	CPU    float64
	Memory uint64
}

// ProcessListState is a ProcessList's sort, selection and pending end.
type ProcessListState struct {
	// Sort is by name, pid, cpu or memory.
	Sort     ui.SortOrder
	rows     ui.ListState
	selected int
	order    []int
	sorted   ui.SortOrder
	data     *Process
	n        int
	confirm  bool
	ending   Process
	started  bool
}

// ProcessListOptions configure a ProcessList.
type ProcessListOptions struct {
	Disabled bool
}

// ProcessListResult carries a confirmed end request.
type ProcessListResult struct {
	Element *ui.Element
	// End is the confirmed process to end, or nil.
	End *Process
}

// ProcessList is a sortable process table with confirmed ending.
func ProcessList(c *ui.Context, procs []Process, s *ProcessListState, opts ProcessListOptions) ProcessListResult {
	if s == nil {
		panic("mujica: ProcessList needs a state")
	}
	if !s.started {
		s.started, s.selected, s.Sort = true, -1, ui.SortOrder{Column: "cpu", Descending: true}
	}
	processListSort(procs, s)
	k := Tokens(c)
	var res ProcessListResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	if disabled {
		s.confirm = false
	}
	ask := func(p Process) {
		if !disabled {
			s.ending, s.confirm = p, true
		}
	}
	res.Element.Children(func() {
		if len(procs) == 0 {
			codeCaption(c, msg(c, "process_list.empty")).Padding(Density(c).Unit() * 4)
			return
		}
		cols := []ui.TableColumn{
			{Title: msg(c, "process_list.name"), ID: "name", Sortable: true, MinWidth: 120},
			{Title: msg(c, "process_list.pid"), ID: "pid", Width: 80, Align: ui.End, Sortable: true},
			{Title: msg(c, "process_list.cpu"), ID: "cpu", Width: 90, Align: ui.End, Sortable: true},
			{Title: msg(c, "process_list.memory"), ID: "memory", Width: 100, Align: ui.End, Sortable: true},
			{Title: msg(c, "process_list.actions"), ID: "end", Width: 56, Fixed: true, Align: ui.Center},
		}
		s.rows.Sort = &s.Sort
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any { return procs[s.order[i]].PID }
		s.rows.Label = func(i int) string { return procs[s.order[i]].Name }
		restore := dataSelectTheme(c, 0)
		table := ui.Table(c, &s.rows, cols, len(procs), func(row, col int) {
			p := procs[s.order[row]]
			switch col {
			case 0:
				ui.Row(c).Gap(Density(c).Unit() * 2).AlignItems(ui.Center).Children(func() {
					codeSelectMark(c, s.selected == row)
					ui.Text(c, p.Name).SingleLine()
				})
			case 1:
				codeMono(c, fmt.Sprint(p.PID))
			case 2:
				codeMono(c, fmt.Sprintf("%.1f%%", p.CPU))
			case 3:
				codeMono(c, processListBytes(p.Memory))
			case 4:
				if iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "process_list.end"), p.Name)).Clicked() {
					ask(p)
				}
			}
		}).Grow(1).MinHeight(Density(c).ControlHeight() * 4)
		restore()
		internal.Ring(table, k.Focus, theme.ControlRadius)
		if table.Changed() {
			processListSort(procs, s)
		}
		if table.Shortcut(0, ui.KeyDelete) && s.selected >= 0 && s.selected < len(procs) {
			ask(procs[s.order[s.selected]])
		}
	})
	if s.confirm {
		r := AlertDialog(c, &s.confirm, AlertDialogOptions{
			Title:        fmt.Sprintf(msg(c, "process_list.confirm"), s.ending.Name, s.ending.PID),
			Description:  msg(c, "process_list.warning"),
			ConfirmLabel: msg(c, "process_list.endButton"),
		})
		if r.Confirmed && !disabled {
			p := s.ending
			res.End = &p
			slog.Debug("mujica: process_list end", "pid", p.PID)
		}
	}
	return res
}

// processListSort reorders when data or sort change.
func processListSort(procs []Process, s *ProcessListState) {
	var first *Process
	if len(procs) > 0 {
		first = &procs[0]
	}
	if s.data == first && s.n == len(procs) && s.sorted == s.Sort && len(s.order) == len(procs) {
		return
	}
	s.data, s.n, s.sorted = first, len(procs), s.Sort
	s.order = s.order[:0]
	for i := range procs {
		s.order = append(s.order, i)
	}
	by := func(a, b int) int {
		pa, pb := procs[a], procs[b]
		switch s.Sort.Column {
		case "name":
			return strings.Compare(strings.ToLower(pa.Name), strings.ToLower(pb.Name))
		case "pid":
			return cmp.Compare(pa.PID, pb.PID)
		case "cpu":
			return cmp.Compare(pa.CPU, pb.CPU)
		case "memory":
			return cmp.Compare(pa.Memory, pb.Memory)
		case "":
			return 0
		}
		panic("mujica: ProcessList sorts by name, pid, cpu or memory, not " + s.Sort.Column)
	}
	slices.SortStableFunc(s.order, func(a, b int) int {
		if s.Sort.Descending {
			return by(b, a)
		}
		return by(a, b)
	})
	slog.Debug("mujica: process_list sort", "by", s.Sort.Column, "descending", s.Sort.Descending)
}

func processListBytes(b uint64) string {
	const unit = 1024
	switch {
	case b >= unit*unit*unit:
		return fmt.Sprintf("%.1f GB", float64(b)/(unit*unit*unit))
	case b >= unit*unit:
		return fmt.Sprintf("%.1f MB", float64(b)/(unit*unit))
	case b >= unit:
		return fmt.Sprintf("%.0f KB", float64(b)/unit)
	}
	return fmt.Sprintf("%d B", b)
}
