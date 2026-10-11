package project

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// KanbanColumn is one column, holding the tasks in Status.
type KanbanColumn struct {
	Status string
	// WIP limits card count; zero allows any count.
	WIP int
}

// KanbanBoardOptions configure a KanbanBoard.
type KanbanBoardOptions struct {
	Columns  []KanbanColumn
	Statuses []WorkflowStatus
	Members  []ProjectMember
	Labels   []TaskLabel
	// Version must change when tasks change in place.
	Version  int
	Disabled bool
	Label    string
}

// KanbanMove identifies a moved card and its destination.
type KanbanMove struct {
	ID, From, To string
	Index        int
}

// KanbanResult is the board and what the user did.
type KanbanResult struct {
	Element *ui.Element
	move    *KanbanMove
	opened  string
}

// Moved returns this frame's applied card move.
func (r KanbanResult) Moved() (KanbanMove, bool) {
	if r.move == nil {
		return KanbanMove{}, false
	}
	return *r.move, true
}

// Opened returns the card the user activated.
func (r KanbanResult) Opened() (string, bool) { return r.opened, r.opened != "" }

// kanbanDrag is the ID of a card being dragged.
type kanbanDrag string

// kanbanState caches column indexes and pending focus.
type kanbanState struct {
	first   *Task
	n, ver  int
	rev, at int
	cols    [][]int
	columns []KanbanColumn
	focus   string
	problem string
	// problemAt is the column the problem belongs to.
	problemAt string
	lists     []*ui.ListState
}

// KanbanBoard supports status columns, WIP limits and card moves.
func KanbanBoard(c *ui.Context, tasks *[]Task, opts KanbanBoardOptions) KanbanResult {
	if len(opts.Columns) == 0 {
		panic("mujica: KanbanBoard needs Columns")
	}
	for i, col := range opts.Columns {
		projectStatus(opts.Statuses, col.Status)
		if col.WIP < 0 || slices.ContainsFunc(opts.Columns[:i], func(o KanbanColumn) bool { return o.Status == col.Status }) {
			panic(fmt.Sprintf("mujica: KanbanBoard column %q repeats or has a negative WIP", col.Status))
		}
	}
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "kanban_board.label")
	}
	r := KanbanResult{}
	board := ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Stretch).Grow(1).MinHeight(0).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = board
	st := *ui.Local(board, "kanban", func() *kanbanState { return &kanbanState{} })
	cols := kanbanColumns(st, *tasks, opts)
	apply := func(id, to, before, how string) {
		from := slices.IndexFunc(*tasks, func(t Task) bool { return t.ID == id })
		ci := slices.IndexFunc(opts.Columns, func(c KanbanColumn) bool { return c.Status == to })
		if from < 0 || ci < 0 {
			return
		}
		t := (*tasks)[from]
		if wip := opts.Columns[ci].WIP; wip > 0 && t.Status != to && len(cols[ci]) >= wip {
			st.problem, st.problemAt = fmt.Sprintf(core.Msg(c, "kanban_board.full"), projectStatus(opts.Statuses, to).Name, wip), to
			slog.Debug("mujica: kanban_board wip refused", "id", id, "to", to)
			return
		}
		var dest []string
		for _, i := range cols[ci] {
			if i != from {
				dest = append(dest, (*tasks)[i].ID)
			}
		}
		rest := slices.Delete(*tasks, from, from+1)
		pos := func(id string) int { return slices.IndexFunc(rest, func(x Task) bool { return x.ID == id }) }
		index, at := len(dest), len(rest)
		switch {
		case before != "":
			index, at = slices.Index(dest, before), pos(before)
		case len(dest) > 0:
			at = pos(dest[len(dest)-1]) + 1
		}
		from0 := t.Status
		t.Status = to
		*tasks = slices.Insert(rest, at, t)
		st.rev++
		st.focus, st.problem = id, ""
		st.lists[ci].ScrollIntoView(index)
		r.move = &KanbanMove{ID: id, From: from0, To: to, Index: index}
		slog.Debug("mujica: kanban_board move", "id", id, "from", from0, "to", to, "index", index, "by", how)
	}
	// Apply moves after building all rows.
	var pending func()
	move := func(id, to, before, how string) {
		if !board.IsDisabled() && id != before && pending == nil {
			pending = func() { apply(id, to, before, how) }
		}
	}
	board.Children(func() {
		for ci, col := range opts.Columns {
			kanbanColumn(c, st, ci, col, func(i int) []int { return cols[i] }, *tasks, opts, &r, move)
		}
	})
	if pending != nil {
		pending()
		c.Invalidate()
	}
	return r
}

// kanbanColumns returns each column's task indexes, rebuilt when tasks change.
func kanbanColumns(st *kanbanState, tasks []Task, opts KanbanBoardOptions) [][]int {
	var first *Task
	if len(tasks) > 0 {
		first = &tasks[0]
	}
	if st.cols != nil && st.first == first && st.n == len(tasks) && st.ver == opts.Version && st.at == st.rev && slices.Equal(st.columns, opts.Columns) {
		return st.cols
	}
	st.first, st.n, st.ver, st.at = first, len(tasks), opts.Version, st.rev
	st.columns = slices.Clone(opts.Columns)
	st.cols = make([][]int, len(opts.Columns))
	for i, t := range tasks {
		ci := slices.IndexFunc(opts.Columns, func(c KanbanColumn) bool { return c.Status == t.Status })
		if ci < 0 {
			projectStatus(opts.Statuses, t.Status)
			continue
		}
		st.cols[ci] = append(st.cols[ci], i)
	}
	for len(st.lists) < len(opts.Columns) {
		st.lists = append(st.lists, &ui.ListState{})
	}
	return st.cols
}

// kanbanColumn draws a header, WIP count and virtual cards.
func kanbanColumn(c *ui.Context, st *kanbanState, ci int, col KanbanColumn, colsOf func(int) []int, tasks []Task, opts KanbanBoardOptions,
	r *KanbanResult, move func(id, to, before, how string)) {
	k := core.Tokens(c)
	d := core.Density(c)
	idx := colsOf(ci)
	status := projectStatus(opts.Statuses, col.Status)
	over := col.WIP > 0 && len(idx) > col.WIP
	full := col.WIP > 0 && len(idx) >= col.WIP
	box := ui.Column(c).Key(col.Status).Width(d.ControlHeight()*7.5).Shrink(0).Gap(d.Unit()*2).Padding(d.Unit()*2).
		Radius(theme.CardRadius).Background(k.Background).Border(theme.BorderWidth, k.Border).Role(ui.RoleGroup).Label(status.Name)
	if id, ok := ui.Drop[kanbanDrag](box); ok && !box.IsDisabled() {
		move(string(id), col.Status, "", "drag")
	}
	if _, ok := ui.DragOver[kanbanDrag](box); ok {
		box.Border(theme.BorderWidth, k.AccentText)
	}
	box.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Children(func() {
			statusGlyphBox(c, status.Group)
			ui.Text(c, status.Name).Font(theme.SerifFont).TextColor(k.Text).Grow(1).SingleLine()
			count := fmt.Sprint(len(idx))
			if col.WIP > 0 {
				count = fmt.Sprintf("%d / %d", len(idx), col.WIP)
			}
			fg := k.TextMuted
			if over {
				fg = k.Danger
				ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14).TextColor(fg)
			}
			wip := ui.Text(c, count).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(fg).FontFeatures("tnum")
			if col.WIP > 0 {
				wip.Label(fmt.Sprintf(core.Msg(c, "kanban_board.wip"), len(idx), col.WIP))
			}
		})
		if over {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "kanban_board.over"), len(idx)-col.WIP)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Danger).Role(ui.RoleStatus)
		} else if full && st.problem != "" && st.problemAt == col.Status {
			ui.Text(c, st.problem).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Danger).Role(ui.RoleStatus)
		}
		ls := st.lists[ci]
		ls.Key = func(i int) any { return tasks[idx[i]].ID }
		at := func(cidx []int, i int) string {
			if i < 0 || i >= len(cidx) {
				return ""
			}
			return tasks[cidx[i]].ID
		}
		ui.List(c, ls, len(idx), func(i int) {
			t := tasks[idx[i]]
			cell := ui.Box(c).PaddingY(d.Unit())
			if id, ok := ui.Drop[kanbanDrag](cell); ok && !cell.IsDisabled() {
				move(string(id), col.Status, t.ID, "drag")
			}
			if _, ok := ui.DragOver[kanbanDrag](cell); ok {
				cell.DrawOver(func(p *ui.Painter, rr ui.Rect) { p.Fill(ui.Rect{X: rr.X, Y: rr.Y, W: rr.W, H: 2}, k.AccentText, 1) })
			}
			cell.Children(func() {
				card := IssueCard(c, t, IssueCardOptions{Statuses: opts.Statuses, Members: opts.Members, Labels: opts.Labels})
				if card.IsDisabled() {
					return
				}
				card.Drag(kanbanDrag(t.ID))
				if st.focus == t.ID {
					card.Focus()
					st.focus = ""
				}
				if card.Clicked() {
					r.opened = t.ID
				}
				switch {
				case card.Shortcut(ui.Alt, ui.KeyUp) && i > 0:
					move(t.ID, col.Status, at(idx, i-1), "keys")
				case card.Shortcut(ui.Alt, ui.KeyDown) && i < len(idx)-1:
					move(t.ID, col.Status, at(idx, i+2), "keys")
				case card.Shortcut(ui.Alt, ui.KeyLeft) && ci > 0:
					move(t.ID, opts.Columns[ci-1].Status, at(colsOf(ci-1), i), "keys")
				case card.Shortcut(ui.Alt, ui.KeyRight) && ci < len(opts.Columns)-1:
					move(t.ID, opts.Columns[ci+1].Status, at(colsOf(ci+1), i), "keys")
				}
			})
		}).Grow(1).MinHeight(d.ControlHeight() * 3).Children(func() {
			if len(idx) == 0 {
				ui.Text(c, core.Msg(c, "kanban_board.empty")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(d.Unit() * 2)
			}
		})
	})
}
