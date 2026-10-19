package project

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// TaskDetailPanelOptions configure a TaskDetailPanel.
type TaskDetailPanelOptions struct {
	Statuses []WorkflowStatus
	Members  []ProjectMember
	Labels   []TaskLabel
	// Activity builds the caller's activity log below the fields.
	Activity func()
	// Today anchors the due date; zero uses the system clock.
	Today    core.Date
	Disabled bool
	// Version changes when task or related data change in place.
	Version int
}

// TaskDetailResult exposes the panel, edits and close events.
type TaskDetailResult struct {
	Element ui.Element
	changed bool
	closed  bool
}

// Changed reports a field edit this frame.
func (r TaskDetailResult) Changed() bool { return r.changed }

// Closed reports a click on Close or Escape.
func (r TaskDetailResult) Closed() bool { return r.closed }

// TaskDetailPanel edits task fields and hosts an activity slot.
func TaskDetailPanel(c *ui.Context, task *Task, opts TaskDetailPanelOptions) TaskDetailResult {
	projectStatus(opts.Statuses, task.Status)
	k := core.Tokens(c)
	d := core.Density(c)
	r := TaskDetailResult{}
	panel := ui.Column(c).Gap(d.Unit()*4).Padding(d.Unit()*4).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Role(ui.RoleGroup).Label(task.ID + " " + task.Title).Disabled(opts.Disabled)
	r.Element = panel
	opts.Disabled = panel.IsDisabled()
	if opts.Disabled {
		copied := *task
		copied.Assignees, copied.Labels, copied.Subtasks = slices.Clone(task.Assignees), slices.Clone(task.Labels), slices.Clone(task.Subtasks)
		task = &copied
	}
	panel.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			IssueIdBadge(c, task.ID, IssueIdBadgeOptions{Copyable: true, Disabled: opts.Disabled})
			ui.Spacer(c)
			if b := core.IconAction(c, icons.Must("x"), core.Msg(c, "task_detail_panel.close")); b.Clicked() && !b.IsDisabled() {
				r.closed = true
			}
		})
		title := input.TextInput(c, &task.Title, input.TextInputOptions{Label: core.Msg(c, "task_detail_panel.title"), Disabled: opts.Disabled})
		title.Input.Font(theme.SerifFont)
		r.changed = title.Changed()
		taskField(c, core.Msg(c, "task_detail_panel.status"), func() {
			r.changed = StatusSelect(c, &task.Status, opts.Statuses, StatusSelectOptions{Disabled: opts.Disabled}).Changed() || r.changed
		})
		taskField(c, core.Msg(c, "task_detail_panel.priority"), func() {
			prios := make([]input.SelectOption[Priority], 0, 5)
			for p := PriorityUrgent; p >= PriorityNone; p-- {
				prios = append(prios, input.SelectOption[Priority]{Value: p, Label: PriorityName(c, p)})
			}
			r.changed = input.Select(c, &task.Priority, prios, input.SelectOptions{Label: core.Msg(c, "task_detail_panel.priority"), Disabled: opts.Disabled}).Changed() || r.changed
		})
		taskField(c, core.Msg(c, "assignee_picker.label"), func() {
			r.changed = AssigneePicker(c, &task.Assignees, opts.Members, AssigneePickerOptions{Disabled: opts.Disabled, Version: opts.Version}).Changed() || r.changed
		})
		taskField(c, core.Msg(c, "due_date_picker.label"), func() {
			r.changed = DueDatePicker(c, &task.Due, DueDatePickerOptions{Done: task.Done, Today: opts.Today, Disabled: opts.Disabled}).Changed() || r.changed
		})
		taskField(c, core.Msg(c, "label_manager.label"), func() {
			options := make([]input.SelectOption[string], len(opts.Labels))
			for i, l := range opts.Labels {
				options[i] = input.SelectOption[string]{Value: l.ID, Label: l.Name}
			}
			r.changed = input.MultiSelect(c, &task.Labels, options, input.MultiSelectOptions{Label: core.Msg(c, "label_manager.label"), Disabled: opts.Disabled}).Changed() || r.changed
		})
		r.changed = SubtaskList(c, &task.Subtasks, SubtaskListOptions{Disabled: opts.Disabled, Version: opts.Version}).Changed() || r.changed
		if opts.Activity != nil {
			ui.Box(c).Height(16).FillWidth().Draw(func(p *ui.Painter, rr ui.Rect) { internal.RuleDiamond(p, rr, 8, 6, k.Ornament) })
			ui.Text(c, core.Msg(c, "task_detail_panel.activity")).Font(theme.SerifFont).TextColor(k.Text)
			opts.Activity()
		}
	})
	if !panel.IsDisabled() && panel.FocusWithin() && c.Shortcut(0, ui.KeyEscape) {
		r.closed = true
	}
	if r.changed {
		slog.Debug("mujica: task_detail_panel change", "id", task.ID)
	}
	if r.closed {
		slog.Debug("mujica: task_detail_panel close", "id", task.ID)
	}
	return r
}

// taskField puts a caption above a control.
func taskField(c *ui.Context, label string, control func()) {
	ui.Column(c).Gap(core.Density(c).Unit()).Children(func() {
		ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(core.Tokens(c).TextMuted)
		control()
	})
}
