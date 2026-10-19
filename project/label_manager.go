package project

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// LabelManagerOptions configure a LabelManager.
type LabelManagerOptions struct {
	Disabled bool
	Label    string
	// Version changes when labels or selection change in place.
	Version int
}

// labelManagerState is the label being edited and the new-label draft.
type labelManagerState struct {
	editing string
	name    string
	draft   string
	problem string
}

// LabelManager edits labels and stores chosen IDs in selected.
func LabelManager(c *ui.Context, labels *[]TaskLabel, selected *[]string, opts LabelManagerOptions) core.ChoiceResult {
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "label_manager.label")
	}
	changed := false
	col := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	st := ui.Local(col, "label-manager", func() *labelManagerState { return &labelManagerState{} })
	s := *st
	cache := ui.Local(col, "label-data", func() labelDataCache { return labelDataCache{} })
	var first *TaskLabel
	var selectedFirst *string
	if len(*labels) > 0 {
		first = &(*labels)[0]
	}
	if len(*selected) > 0 {
		selectedFirst = &(*selected)[0]
	}
	if !cache.ready || cache.first != first || cache.n != len(*labels) || cache.selectedFirst != selectedFirst || cache.selectedN != len(*selected) || cache.version != opts.Version {
		*cache = labelDataCache{first: first, n: len(*labels), selectedFirst: selectedFirst, selectedN: len(*selected), version: opts.Version, ready: true, selected: map[string]bool{}}
		known := map[string]bool{}
		for _, l := range *labels {
			datetime.EventTone(k, l.Color)
			if l.ID == "" || known[l.ID] {
				panic("mujica: LabelManager needs unique non-empty IDs")
			}
			known[l.ID] = true
		}
		for _, id := range *selected {
			if !known[id] {
				panic("mujica: LabelManager selection contains an unknown ID")
			}
			cache.selected[id] = true
		}
	}
	state := ui.Local(col, "label-rows", func() ui.ListState { return ui.ListState{} })
	state.Key = func(i int) any { return (*labels)[i].ID }

	live := func() bool { return !col.IsDisabled() }
	taken := func(name, except string) bool {
		return slices.ContainsFunc(*labels, func(l TaskLabel) bool { return l.ID != except && strings.EqualFold(l.Name, name) })
	}
	col.Children(func() {
		removed := -1
		ui.List(c, state, len(*labels), func(i int) {
			l := &(*labels)[i]
			ui.Row(c.Key(l.ID)).MinHeight(d.ControlHeight()).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				on := cache.selected[l.ID]
				box := ui.CheckboxBase(c, &on).Size(theme.MinHitSize, theme.MinHitSize).Center().Label(l.Name)
				internal.Ring(box, k.Focus, theme.ControlRadius)
				if box.Draw(func(p *ui.Painter, r ui.Rect) { labelSwatch(p, r, datetime.EventTone(k, l.Color), on, k) }).Changed() && live() {
					if on {
						*selected = append(*selected, l.ID)
					} else {
						*selected = slices.DeleteFunc(*selected, func(id string) bool { return id == l.ID })
					}
					changed = true
					slog.Debug("mujica: label_manager select", "id", l.ID, "selected", on)
				}
				if s.editing == l.ID && live() {
					labelEditRow(c, s, l, taken, &changed)
					return
				}
				ui.Text(c, l.Name).Grow(1).TextColor(k.Text).SingleLine().Ellipsis("…")
				if b := core.IconAction(c, icons.Must("palette"), fmt.Sprintf(core.Msg(c, "label_manager.edit"), l.Name)); b.Clicked() && live() {
					s.editing, s.name, s.problem = l.ID, l.Name, ""
				}
				if b := core.IconAction(c, icons.Must("trash-2"), fmt.Sprintf(core.Msg(c, "label_manager.delete"), l.Name)); b.Clicked() && live() {
					removed = i
				}
			})
		}).Height(min(float32(len(*labels)), 8) * d.ControlHeight())
		if removed >= 0 {
			id := (*labels)[removed].ID
			slog.Debug("mujica: label_manager delete", "id", id)
			*labels = slices.Delete(*labels, removed, removed+1)
			*selected = slices.DeleteFunc(*selected, func(x string) bool { return x == id })
			changed = true
		}
		if len(*labels) == 0 {
			ui.Text(c, core.Msg(c, "label_manager.empty")).TextColor(k.TextMuted)
		}
		in := input.TextInput(c, &s.draft, input.TextInputOptions{Placeholder: core.Msg(c, "label_manager.new"), Label: core.Msg(c, "label_manager.new"),
			Disabled: opts.Disabled, Error: s.problem})
		if in.Submitted() && live() {
			name := strings.TrimSpace(s.draft)
			switch {
			case name == "":
			case taken(name, ""):
				s.problem = fmt.Sprintf(core.Msg(c, "label_manager.taken"), name)
			default:
				l := TaskLabel{ID: labelNewID(*labels), Name: name, Color: datetime.EventColor(len(*labels) % 6)}
				*labels = append(*labels, l)
				s.draft, s.problem, changed = "", "", true
				slog.Debug("mujica: label_manager create", "id", l.ID)
			}
		}
	})
	if changed {
		cache.ready = false
		c.Invalidate()
	}
	return core.NewChoiceResult(col, changed)
}

// labelEditRow renames a label and picks its color.
func labelEditRow(c *ui.Context, s *labelManagerState, l *TaskLabel, taken func(string, string) bool, changed *bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	ui.Column(c).Grow(1).Gap(d.Unit()).Children(func() {
		in := input.TextInput(c, &s.name, input.TextInputOptions{Label: core.Msg(c, "label_manager.name"), Error: s.problem})
		core.AutoFocus(c, in.Input)
		ui.Row(c).Gap(d.Unit()).Role(ui.RoleRadioGroup).Label(core.Msg(c, "label_manager.color")).Children(func() {
			for i := range 6 {
				col := datetime.EventColor(i)
				sel := l.Color
				b := ui.RadioBase(c, &sel, col).Size(theme.MinHitSize, theme.MinHitSize).Label(core.Msg(c, labelColorKeys[i]))
				internal.Ring(b, k.Focus, theme.ControlRadius)
				on := l.Color == col
				b.Draw(func(p *ui.Painter, r ui.Rect) { labelSwatch(p, r, datetime.EventTone(k, col), on, k) })
				if b.Changed() && !b.IsDisabled() {
					l.Color, *changed = col, true
					slog.Debug("mujica: label_manager color", "id", l.ID, "color", col)
				}
			}
		})
		done := in.Submitted() || in.Input.Shortcut(0, ui.KeyEscape)
		if in.Submitted() {
			name := strings.TrimSpace(s.name)
			switch {
			case name == "":
				done = false
			case taken(name, l.ID):
				s.problem, done = fmt.Sprintf(core.Msg(c, "label_manager.taken"), name), false
			case name != l.Name:
				l.Name, *changed = name, true
				slog.Debug("mujica: label_manager rename", "id", l.ID)
			}
		}
		if done {
			s.editing, s.problem = "", ""
		}
	})
}

var labelColorKeys = [...]string{"label_manager.accent", "label_manager.info", "label_manager.success", "label_manager.warning",
	"label_manager.danger", "label_manager.neutral"}

// labelSwatch draws a color dot with a selection check.
func labelSwatch(p *ui.Painter, r ui.Rect, col ui.Color, on bool, k theme.Tokens) {
	box := ui.Rect{X: r.X + r.W/2 - 8, Y: r.Y + r.H/2 - 8, W: 16, H: 16}
	p.Fill(box, col, 8)
	if on {
		internal.Mark(p, ui.Rect{X: box.X + 2, Y: box.Y + 2, W: 12, H: 12}, "check", k.Surface)
	}
}

// labelNewID returns an ID no label uses.
func labelNewID(labels []TaskLabel) string {
	for n := len(labels) + 1; ; n++ {
		id := fmt.Sprintf("label-%d", n)
		if !slices.ContainsFunc(labels, func(l TaskLabel) bool { return l.ID == id }) {
			return id
		}
	}
}

type labelDataCache struct {
	first                 *TaskLabel
	selectedFirst         *string
	n, selectedN, version int
	ready                 bool
	selected              map[string]bool
}
