package mujica

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TaskGrouping chooses the sections of a TaskList.
type TaskGrouping int

const (
	TaskGroupNone TaskGrouping = iota
	TaskGroupStatus
	TaskGroupPriority
)

// TaskFilter shows all, open or done tasks.
type TaskFilter int

const (
	TaskShowAll TaskFilter = iota
	TaskShowOpen
	TaskShowDone
)

// TaskListOptions configure a TaskList.
type TaskListOptions struct {
	GroupBy TaskGrouping
	// Statuses order the status groups; required with TaskGroupStatus.
	Statuses []WorkflowStatus
	Members  []Member
	// Version must change when tasks change in place.
	Version int
	// Today anchors overdue marks; zero uses the system clock.
	Today    Date
	Disabled bool
	Label    string
}

// taskRow identifies a heading or task.
type taskRow struct {
	head string
	task int
}

// taskListCache keeps rows until the data, grouping or filter change.
type taskListCache struct {
	first    *Task
	n, ver   int
	group    TaskGrouping
	filter   TaskFilter
	rows     []taskRow
	counted  [3]int
	statuses []WorkflowStatus
	locale   Locale
}

// TaskList groups and filters tasks, virtualizing visible rows.
func TaskList(c *ui.Context, tasks []Task, filter *TaskFilter, opts TaskListOptions) ProjectResult {
	if opts.GroupBy < TaskGroupNone || opts.GroupBy > TaskGroupPriority || *filter < TaskShowAll || *filter > TaskShowDone {
		panic(fmt.Sprintf("mujica: TaskList grouping %d or filter %d out of range", opts.GroupBy, *filter))
	}
	if opts.GroupBy == TaskGroupStatus && len(opts.Statuses) == 0 {
		panic("mujica: TaskList grouping by status needs Statuses")
	}
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "task_list.label")
	}
	r := ProjectResult{}
	root := ui.Column(c).Grow(1).MinHeight(0).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	cache := ui.Local(root, "task-list", func() *taskListCache { return &taskListCache{} })
	var first *Task
	if len(tasks) > 0 {
		first = &tasks[0]
	}
	if x := *cache; x.rows == nil || x.first != first || x.n != len(tasks) || x.ver != opts.Version || x.group != opts.GroupBy || x.filter != *filter || !slices.Equal(x.statuses, opts.Statuses) || x.locale != current(c).locale {
		*cache = &taskListCache{first: first, n: len(tasks), ver: opts.Version, group: opts.GroupBy, filter: *filter, statuses: slices.Clone(opts.Statuses), locale: current(c).locale}
		(*cache).rows, (*cache).counted = taskListRows(c, tasks, opts, *filter)
	}
	rows := (*cache).rows
	state := ui.Local(root, "task-list-state", func() *ui.ListState { return &ui.ListState{} })
	st := *state
	st.Key = func(i int) any {
		if rows[i].task < 0 {
			return "head:" + rows[i].head
		}
		return tasks[rows[i].task].ID
	}
	st.Header = func(i int) bool { return rows[i].task < 0 }
	root.Children(func() {
		counts := (*cache).counted
		sel := int(*filter)
		taskListFilter(c, &sel, []string{
			fmt.Sprintf(msg(c, "task_list.all"), counts[0]),
			fmt.Sprintf(msg(c, "task_list.open"), counts[1]),
			fmt.Sprintf(msg(c, "task_list.done"), counts[2]),
		})
		if sel != int(*filter) && !root.IsDisabled() {
			*filter = TaskFilter(sel)
			slog.Debug("mujica: task_list filter", "filter", *filter)
		}
		list := ui.List(c, st, len(rows), func(i int) {
			row := rows[i]
			if row.task < 0 {
				h := ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Center).PaddingX(d.Unit() * 2).Background(k.Background).Role(ui.RoleColumnHeader)
				h.Children(func() { ui.Text(c, row.head).Font(theme.SerifFont).TextColor(k.Text) })
				h.DrawOver(func(p *ui.Painter, rr ui.Rect) { internal.Rule(p, rr.X, rr.X+rr.W, rr.Y+rr.H-1, k.Ornament) })
				return
			}
			ir := TaskItem(c, &tasks[row.task], TaskItemOptions{Members: opts.Members, Today: opts.Today, Disabled: opts.Disabled})
			if id, ok := ir.Changed(); ok {
				r.changed = id
				(*cache).rows = nil
				c.Invalidate()
			}
			if id, ok := ir.Opened(); ok {
				r.opened = id
			}
		}).Grow(1).MinHeight(d.ControlHeight() * 4)
		list.Children(func() {
			if len(rows) == 0 {
				ui.Text(c, msg(c, "task_list.empty")).TextColor(k.TextMuted).Padding(d.Unit() * 3)
			}
		})
	})
	return r
}

// taskListRows computes groups and completion counts.
func taskListRows(c *ui.Context, tasks []Task, opts TaskListOptions, filter TaskFilter) ([]taskRow, [3]int) {
	var counts [3]int
	keep := make([]int, 0, len(tasks))
	for i, t := range tasks {
		counts[0]++
		if t.Done {
			counts[2]++
		} else {
			counts[1]++
		}
		if filter == TaskShowAll || (filter == TaskShowDone) == t.Done {
			keep = append(keep, i)
		}
	}
	rows := make([]taskRow, 0, len(keep)+8)
	switch opts.GroupBy {
	case TaskGroupNone:
		for _, i := range keep {
			rows = append(rows, taskRow{task: i})
		}
	case TaskGroupStatus:
		for _, s := range opts.Statuses {
			rows = taskListSection(rows, s.Name, keep, func(i int) bool { return tasks[i].Status == s.ID })
		}
		for _, i := range keep {
			projectStatus(opts.Statuses, tasks[i].Status)
		}
	case TaskGroupPriority:
		for p := PriorityUrgent; p >= PriorityNone; p-- {
			rows = taskListSection(rows, PriorityName(c, p), keep, func(i int) bool { return tasks[i].Priority == p })
		}
	}
	return slices.Clip(rows), counts
}

// taskListSection appends a heading and its tasks when any match.
func taskListSection(rows []taskRow, head string, keep []int, in func(int) bool) []taskRow {
	start := len(rows)
	for _, i := range keep {
		if in(i) {
			if len(rows) == start {
				rows = append(rows, taskRow{head: head, task: -1})
			}
			rows = append(rows, taskRow{task: i})
		}
	}
	return rows
}

// taskListFilter marks the selected completion filter with an inset.
func taskListFilter(c *ui.Context, value *int, labels []string) {
	k, d := Tokens(c), Density(c)
	parts := ui.SegmentedBase(c, value, len(labels))
	parts.Track.MinHeight(max(d.ControlHeight(), theme.MinHitSize+4)).AlignItems(ui.Stretch).Padding(2).Gap(2).Radius(theme.ControlRadius+2).
		Background(k.Background).Border(theme.BorderWidth, k.ControlBorder).Label(msg(c, "task_list.filter"))
	parts.Track.Children(func() {
		for i, label := range labels {
			selected := i == *value
			seg := parts.Segment(i).MinWidth(theme.MinHitSize).MinHeight(theme.MinHitSize).PaddingX(d.Unit() * 3).Radius(theme.ControlRadius).Label(label)
			fg := k.TextMuted
			if selected {
				fg = k.AccentText
				seg.Background(k.Selection).FontWeight(600)
				seg.Draw(func(p *ui.Painter, r ui.Rect) {
					p.Stroke(ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}, k.AccentText, theme.ControlRadius, 1)
				})
			} else if seg.Hovered() {
				fg = k.Text
				seg.Background(k.SurfaceHover)
			}
			internal.Ease(seg, Motion(c, theme.StateDuration), fg)
			internal.Ring(seg, k.Focus, theme.ControlRadius)
			seg.Children(func() { ui.Text(c, label).SingleLine() })
		}
	})
}
