package project

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TaskItemOptions configure a TaskItem.
type TaskItemOptions struct {
	Members []ProjectMember
	// Today anchors the overdue mark; zero uses the system clock.
	Today    core.Date
	Disabled bool
}

// taskItemEdit is a title being edited inline.
type taskItemEdit struct {
	on    bool
	draft string
}

// TaskItem displays task fields with completion and inline editing.
func TaskItem(c *ui.Context, task *Task, opts TaskItemOptions) ProjectResult {
	projectCheckPriority(task.Priority)
	for _, id := range task.Assignees {
		projectMember(opts.Members, id)
	}
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	r := ProjectResult{}
	row := ui.Row(c).MinHeight(d.ControlHeight()).Gap(d.Unit() * 3).AlignItems(ui.Center).PaddingX(d.Unit() * 2).
		Role(ui.RoleListItem).Label(task.Title).Disabled(opts.Disabled)
	r.Element = row
	edit := ui.Local(row, "task-item-edit", func() taskItemEdit { return taskItemEdit{} })
	if row.Hovered() {
		row.Background(k.SurfaceHover)
	}
	row.Children(func() {
		done := task.Done
		if input.Checkbox(c, &done, "", input.CheckboxOptions{Disabled: opts.Disabled}).Label(task.Title).Changed() && !row.IsDisabled() {
			task.Done, r.changed = done, task.ID
			slog.Debug("mujica: task_item done", "id", task.ID, "done", done)
		}
		if edit.on && !row.IsDisabled() {
			in := input.TextInput(c, &edit.draft, input.TextInputOptions{Label: core.Msg(c, "task_item.rename")})
			in.Element.Grow(1)
			core.AutoFocus(c, in.Input)
			switch {
			case in.Submitted() && strings.TrimSpace(edit.draft) != "":
				if t := strings.TrimSpace(edit.draft); t != task.Title {
					task.Title, r.changed = t, task.ID
					slog.Debug("mujica: task_item rename", "id", task.ID)
				}
				edit.on = false
			case in.Input.Shortcut(0, ui.KeyEscape):
				edit.on = false
			}
			return
		}
		edit.on = false
		title := ui.ButtonBase(c).Grow(1).Shrink(1).Justify(ui.Start).MinHeight(theme.MinHitSize).Radius(theme.ControlRadius).Label(task.Title)
		internal.Ring(title, k.Focus, theme.ControlRadius)
		title.Children(func() {
			t := ui.Text(c, task.Title).SingleLine().Ellipsis("…").TextColor(k.Text)
			if task.Done {
				t.Strikethrough().TextColor(k.TextMuted)
			}
		})
		if !title.IsDisabled() {
			if title.DoubleClicked() || title.Shortcut(0, ui.KeyF2) {
				*edit = taskItemEdit{on: true, draft: task.Title}
			} else if title.Clicked() {
				r.opened = task.ID
			}
		}
		PriorityIndicator(c, task.Priority, PriorityIndicatorOptions{IconOnly: true})
		if text, late := projectDueText(c, task.Due, today, task.Done); text != "" {
			fg := k.TextMuted
			if late {
				fg = k.Danger
			}
			ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(fg).SingleLine().Shrink(0)
		}
		if len(task.Assignees) > 0 {
			items := make([]display.AvatarItem, len(task.Assignees))
			for i, id := range task.Assignees {
				items[i] = display.AvatarItem{Name: projectMember(opts.Members, id).Name}
			}
			display.AvatarGroup(c, items, display.AvatarGroupOptions{Max: 2, Size: 24})
		}
	})
	return r
}
