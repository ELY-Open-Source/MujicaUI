package project

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AssigneePickerOptions configure an AssigneePicker.
type AssigneePickerOptions struct {
	// Max limits selections; zero allows any number.
	Max      int
	Disabled bool
	Label    string
	// Version changes when members or selection change in place.
	Version int
}

type assigneeState struct {
	first            *ProjectMember
	n, version       int
	ready, open      bool
	query, filtered  string
	options, shown   []input.SelectOption[string]
	members          map[string]ProjectMember
	selected         map[string]bool
	selectedFirst    *string
	selectedN, count int
	list             ui.ListState
	focus            string
}

// AssigneePicker searches and selects members, showing their avatars.
func AssigneePicker(c *ui.Context, ids *[]string, members []ProjectMember, opts AssigneePickerOptions) core.ChoiceResult {
	if opts.Max < 0 {
		panic("mujica: AssigneePicker Max must be non-negative")
	}
	d, k := core.Density(c), core.Tokens(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "assignee_picker.label")
	}
	row := ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Disabled(opts.Disabled)
	st := ui.Local(row, "assignees", func() assigneeState { return assigneeState{} })
	var first *ProjectMember
	var selectedFirst *string
	if len(members) > 0 {
		first = &members[0]
	}
	if len(*ids) > 0 {
		selectedFirst = &(*ids)[0]
	}
	dataChanged := !st.ready || st.first != first || st.n != len(members) || st.version != opts.Version
	if dataChanged {
		st.first, st.n, st.version, st.ready = first, len(members), opts.Version, true
		st.members = make(map[string]ProjectMember, len(members))
		st.options = make([]input.SelectOption[string], len(members))
		for i, m := range members {
			if _, exists := st.members[m.ID]; exists || m.ID == "" {
				panic("mujica: AssigneePicker needs unique non-empty IDs")
			}
			st.members[m.ID] = m
			st.options[i] = input.SelectOption[string]{Value: m.ID, Label: m.Name}
		}
		st.shown = nil
	}
	before := *ids
	var trigger *ui.Element
	row.Children(func() {
		if len(*ids) > 0 {
			items := make([]display.AvatarItem, min(3, len(*ids)))
			for i := range items {
				items[i] = display.AvatarItem{Name: st.members[(*ids)[i]].Name}
			}
			display.AvatarGroup(c, items, display.AvatarGroupOptions{Max: 3, Size: 28})
		}
		trigger = ui.ButtonBase(c).Disabled(row.IsDisabled()).Role(ui.RolePopUpButton).Label(label)
		if trigger.IsDisabled() {
			st.open = false
		}
		if !trigger.IsDisabled() && (trigger.Clicked() || !st.open && trigger.Shortcut(0, ui.KeyDown)) {
			st.open, st.query = !st.open, ""
		}
		if len(*ids) > 0 && !trigger.IsDisabled() && trigger.Shortcut(0, ui.KeyBackspace) {
			*ids = slices.Clone((*ids)[:len(*ids)-1])
		}
		input.SelectField(c, trigger, st.open, func() {
			options := make([]input.SelectOption[string], min(2, len(*ids)))
			for i := range options {
				options[i] = input.SelectOption[string]{Value: (*ids)[i], Label: st.members[(*ids)[i]].Name}
			}
			input.MultiSelectTags(c, ids, options, input.MultiSelectOptions{MaxTags: 2, Disabled: trigger.IsDisabled(), Placeholder: core.Msg(c, "assignee_picker.placeholder")})
		})
		ui.PopoverBase(c, trigger, &st.open, func(panel *ui.Element) {
			input.SelectPanel(c, panel)
			ww, wh := c.Size()
			panel.MinWidth(min(ww-8, trigger.Bounds().W)).MaxWidth(max(1, ww-8)).MaxHeight(max(1, wh-8))
			ui.Scroll(c).MaxHeight(max(1, wh-24)).Children(func() {
				search := input.SelectSearch(c, &st.query, core.Msg(c, "multi_select.search"))
				q := strings.ToLower(strings.TrimSpace(st.query))
				if st.shown == nil || st.filtered != q {
					st.shown = make([]input.SelectOption[string], 0)
					for _, o := range st.options {
						if strings.Contains(strings.ToLower(o.Label), q) {
							st.shown = append(st.shown, o)
						}
					}
					st.filtered = q
					st.list.ScrollTo(0, ui.Start)
					dataChanged = true
				}
				if dataChanged || st.selectedFirst != selectedFirst || st.selectedN != len(*ids) || st.selected == nil {
					st.selected = make(map[string]bool, len(*ids))
					for _, id := range *ids {
						if _, ok := st.members[id]; !ok {
							panic("mujica: AssigneePicker selection contains an unknown ID")
						}
						if st.selected[id] {
							panic("mujica: AssigneePicker selection contains duplicate IDs")
						}
						st.selected[id] = true
					}
					st.selectedFirst, st.selectedN, st.count = selectedFirst, len(*ids), 0
					for _, o := range st.shown {
						if st.selected[o.Value] {
							st.count++
						}
					}
				}
				if len(st.shown) == 0 {
					input.SelectNote(c, core.Msg(c, "multi_select.no_results"))
					return
				}
				on, full := st.count == len(st.shown), opts.Max > 0 && len(*ids) >= opts.Max
				mark := internal.Unchecked
				if on {
					mark = internal.Checked
				} else if st.count > 0 {
					mark = internal.Mixed
				}
				var all *ui.Element
				ui.Row(c).MinHeight(d.ControlHeight()).Children(func() {
					all = input.MultiSelectCheck(c, k, &on, mark, core.Msg(c, "multi_select.all"), full && st.count == 0)
				})
				if search.Shortcut(0, ui.KeyDown) {
					all.Focus()
				}
				if all.Changed() && !all.IsDisabled() {
					*ids = input.MultiSelectAll(*ids, st.shown, on, opts.Max)
				}
				if all.Shortcut(0, ui.KeyDown) && !all.IsDisabled() {
					for i, o := range st.shown {
						if full && !st.selected[o.Value] {
							continue
						}
						st.focus = o.Value
						st.list.ScrollIntoView(i)
						c.Invalidate()
						break
					}
				}
				st.list.Key = func(i int) any { return st.shown[i].Value }
				ui.List(c, &st.list, len(st.shown), func(i int) {
					o := st.shown[i]
					selected := st.selected[o.Value]
					state := internal.Unchecked
					if selected {
						state = internal.Checked
					}
					var cb *ui.Element
					ui.Row(c).MinHeight(d.ControlHeight()).Children(func() { cb = input.MultiSelectCheck(c, k, &selected, state, o.Label, full && !selected) })
					if st.focus == o.Value && !cb.IsDisabled() {
						cb.Focus()
						st.focus = ""
					}
					if cb.Changed() && !cb.IsDisabled() {
						*ids = input.MultiSelectSet(*ids, o.Value, selected)
					}
					if cb.IsDisabled() {
						return
					}
					to := i
					switch {
					case cb.Shortcut(0, ui.KeyDown):
						to = assigneeNext(st, i, 1, full)
					case cb.Shortcut(0, ui.KeyUp):
						to = assigneeNext(st, i, -1, full)
					case cb.Shortcut(0, ui.KeyHome):
						to = assigneeNext(st, -1, 1, full)
					case cb.Shortcut(0, ui.KeyEnd):
						to = assigneeNext(st, len(st.shown), -1, full)
					}
					if to != i {
						st.focus = st.shown[to].Value
						st.list.ScrollIntoView(to)
						c.Invalidate()
					}
				}).Height(min(240, min(float32(len(st.shown))*d.ControlHeight(), max(d.ControlHeight(), wh-140)))).Padding(4)
				if opts.Max > 0 {
					ui.Text(c, fmt.Sprintf(core.Msg(c, "multi_select.limit"), len(*ids), opts.Max)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				}
			})
		})
	})
	changed := !slices.Equal(before, *ids)
	if dataChanged && !st.open {
		for _, id := range *ids {
			if _, ok := st.members[id]; !ok {
				panic("mujica: AssigneePicker selection contains an unknown ID")
			}
		}
	}
	if changed {
		st.selected = nil
		slog.Debug("mujica: assignee_picker change", "count", len(*ids))
		c.Invalidate()
	}
	return core.NewChoiceResult(row, changed)
}

// assigneeNext skips choices disabled by the selection limit.
func assigneeNext(st *assigneeState, from, direction int, full bool) int {
	for i := from + direction; i >= 0 && i < len(st.shown); i += direction {
		if !full || st.selected[st.shown[i].Value] {
			return i
		}
	}
	return from
}
