package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// SubtaskListOptions configure a SubtaskList.
type SubtaskListOptions struct {
	Disabled bool
	Label    string
	// Version changes when subtasks change in place.
	Version int
}

// subtaskDrag is the ID of a subtask being dragged.
type subtaskDrag string

// SubtaskList edits and reorders subtasks with completion progress.
func SubtaskList(c *ui.Context, subs *[]Subtask, opts SubtaskListOptions) ChoiceResult {
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "subtask_list.label")
	}
	changed := false
	col := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	draft := ui.Local(col, "subtask-draft", func() string { return "" })
	cache := ui.Local(col, "subtask-summary", func() subtaskCache { return subtaskCache{} })
	var first *Subtask
	if len(*subs) > 0 {
		first = &(*subs)[0]
	}
	if !cache.ready || cache.first != first || cache.n != len(*subs) || cache.version != opts.Version {
		*cache = subtaskCache{first: first, n: len(*subs), version: opts.Version, ready: true}
		seen := map[string]bool{}
		for _, s := range *subs {
			if s.ID == "" || seen[s.ID] {
				panic(fmt.Sprintf("mujica: SubtaskList subtask %q is empty or repeated", s.ID))
			}
			seen[s.ID] = true
			if s.Done {
				cache.done++
			}
		}

	}
	state := ui.Local(col, "subtask-rows", func() ui.ListState { return ui.ListState{} })
	state.Key = func(i int) any { return (*subs)[i].ID }
	var pending func()

	move := func(from, to int, how string) {
		if col.IsDisabled() || to < 0 || to >= len(*subs) || from == to {
			return
		}

		if pending == nil {
			pending = func() {
				step := (*subs)[from]
				*subs = slices.Insert(slices.Delete(*subs, from, from+1), to, step)
				changed = true
				state.ScrollIntoView(to)
				slog.Debug("mujica: subtask_list move", "id", step.ID, "to", to, "by", how)
			}
		}
	}
	col.Children(func() {
		done := cache.done
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			frac := 0.0
			if len(*subs) > 0 {
				frac = float64(done) / float64(len(*subs))
			}
			ui.Box(c).Grow(1).Children(func() { projectBar(c, frac, label) })
			ui.Text(c, fmt.Sprintf(msg(c, "subtask_list.progress"), done, len(*subs))).FontSize(FontSize(c, theme.CaptionSize)).
				TextColor(k.TextMuted).FontFeatures("tnum")
		})
		removed := -1
		ui.List(c, state, len(*subs), func(i int) {
			s := &(*subs)[i]
			row := ui.Row(c).Key(s.ID).MinHeight(d.ControlHeight()).Gap(d.Unit() * 2).AlignItems(ui.Center).Role(ui.RoleListItem)
			if id, ok := ui.Drop[subtaskDrag](row); ok {
				if from := slices.IndexFunc(*subs, func(o Subtask) bool { return o.ID == string(id) }); from >= 0 {
					move(from, i, "drag")
				}
			}
			if _, over := ui.DragOver[subtaskDrag](row); over {
				row.Background(k.Selection)
			}
			row.Children(func() {
				grip := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Center().Cursor(ui.CursorGrab).Label(fmt.Sprintf(msg(c, "subtask_list.drag"), s.Title))
				if !grip.IsDisabled() {
					grip.Drag(subtaskDrag(s.ID))
				}
				grip.Children(func() { ui.Icon(c, icons.Must("grip-vertical")).Size(14, 14).TextColor(k.TextMuted) })
				v := s.Done
				cb := Checkbox(c, &v, s.Title, CheckboxOptions{Disabled: opts.Disabled})
				if cb.Changed() && !cb.IsDisabled() {
					s.Done, changed = v, true
				}
				if !cb.IsDisabled() {
					if cb.Shortcut(ui.Alt, ui.KeyUp) {
						move(i, i-1, "keys")
					}
					if cb.Shortcut(ui.Alt, ui.KeyDown) {
						move(i, i+1, "keys")
					}
				}
				ui.Spacer(c)
				if b := iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "subtask_list.remove"), s.Title)); b.Clicked() && !b.IsDisabled() {
					removed = i
				}
			})
		}).Height(min(float32(len(*subs)), 8) * d.ControlHeight())
		if removed >= 0 && !col.IsDisabled() {
			slog.Debug("mujica: subtask_list remove", "id", (*subs)[removed].ID)
			*subs, changed = slices.Delete(*subs, removed, removed+1), true
		}
		in := TextInput(c, draft, TextInputOptions{Placeholder: msg(c, "subtask_list.add"), Label: msg(c, "subtask_list.add"), Disabled: opts.Disabled})
		if in.Submitted() && strings.TrimSpace(*draft) != "" && !col.IsDisabled() {
			*subs = append(*subs, Subtask{ID: subtaskNewID(*subs), Title: strings.TrimSpace(*draft)})
			*draft, changed = "", true
			slog.Debug("mujica: subtask_list add")
		}
	})
	if pending != nil {
		pending()
	}
	if changed {
		cache.ready = false
		c.Invalidate()
	}
	return ChoiceResult{Element: col, changed: changed}
}

// subtaskNewID returns an ID no subtask uses.
func subtaskNewID(subs []Subtask) string {
	for n := len(subs) + 1; ; n++ {
		id := fmt.Sprintf("sub-%d", n)
		if !slices.ContainsFunc(subs, func(s Subtask) bool { return s.ID == id }) {
			return id
		}
	}
}

type subtaskCache struct {
	first            *Subtask
	n, version, done int
	ready            bool
}
