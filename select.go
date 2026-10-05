package mujica

import (
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SelectOption is one choice of a Select or MultiSelect.
type SelectOption[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
}

// SelectOptions configure a Select.
type SelectOptions struct {
	// Placeholder shows while no option holds the value; empty uses "Select…".
	Placeholder string
	// Clearable adds a clear button and Backspace/Delete, which set the zero value.
	Clearable bool
	Disabled  bool
	// Empty shows in the popup when there are no options; empty uses "No options".
	Empty string
	// Label names the select for assistive technology.
	Label string
}

// selectKey keeps this file's state apart from SelectBase's own.
type selectKey struct{ name string }

// selectTypeahead is the prefix typed on a closed select, and when.
type selectTypeahead struct {
	prefix  string
	pending bool
	at      time.Time
}

// Select chooses *value from options in a popup. Enter, Space and the
// arrows open it; the arrows, Home and End move; Enter chooses; Escape
// closes. Typing on the closed select chooses the option starting with it.
func Select[T comparable](c *ui.Context, value *T, options []SelectOption[T], opts SelectOptions) ChoiceResult {
	k := Tokens(c)
	old := *value
	if opts.Disabled {
		// SelectBase takes clicks as it is built: a disabled one works on a copy.
		scratch := old
		value = &scratch
	}
	s := ui.SelectBase(c, value)
	tr := s.Trigger.Disabled(opts.Disabled)
	if opts.Label != "" {
		tr.Label(opts.Label)
	}
	ta := ui.Local(tr, selectKey{"typeahead"}, func() selectTypeahead { return selectTypeahead{} })
	tr.TextCaret(ui.Rect{W: 1, H: Density(c).ControlHeight()}).HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputText || strings.TrimSpace(ev.Text) == "" {
			return false
		}
		if time.Since(ta.at) > time.Second {
			ta.prefix = ""
		}
		ta.prefix += strings.ToLower(ev.Text)
		ta.pending, ta.at = true, time.Now()
		return true
	})
	if ta.pending {
		ta.pending = false
		if !s.Open() {
			if i := selectFind(options, *value, ta.prefix); i >= 0 {
				*value = options[i].Value
			}
		}
	}
	var zero T
	label, has := selectLabel(options, *value)
	if opts.Clearable && has && (tr.Shortcut(0, ui.KeyBackspace) || tr.Shortcut(0, ui.KeyDelete)) {
		*value, has = zero, false
	}
	selectField(c, tr, s.Open(), func() {
		if has {
			ui.Text(c, label).SingleLine().Grow(1).Shrink(1)
		} else {
			ui.Text(c, selectOr(opts.Placeholder, msg(c, "select.placeholder"))).SingleLine().TextColor(k.TextMuted).Grow(1).Shrink(1)
		}
		if opts.Clearable && has && !opts.Disabled && selectClear(c, k, msg(c, "select.clear")) {
			*value = zero
		}
	})
	lastHi := ui.Local(tr, selectKey{"shown"}, func() int { return -1 })
	s.Popup(func(panel *ui.Element) {
		selectPanel(c, panel)
		ui.Scroll(c).MaxHeight(280).Children(func() {
			if len(options) == 0 {
				selectNote(c, selectOr(opts.Empty, msg(c, "select.empty")))
			}
			for i, o := range options {
				if o.Disabled {
					selectRow(c, ui.Row(c).Disabled(true), false, false, o.Label)
					continue
				}
				it := s.Item(o.Value)
				if it.Highlighted() && *lastHi != i {
					*lastHi = i
					it.ScrollIntoView()
				}
				selectRow(c, it, it.Highlighted(), o.Value == *value, o.Label)
			}
		})
	})
	if !s.Open() {
		*lastHi = -1
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: select changed", "label", opts.Label, "value", *value)
	}
	return ChoiceResult{tr, changed}
}

// selectFind returns the enabled option matching typed prefix.
func selectFind[T comparable](options []SelectOption[T], current T, prefix string) int {
	start := 0
	// Repeating one letter cycles through the options starting with it.
	if r := []rune(prefix); strings.Count(prefix, string(r[0])) == len(r) {
		prefix = string(r[0])
		for i, o := range options {
			if o.Value == current {
				start = i + 1
			}
		}
	}
	for n := range options {
		i := (start + n) % len(options)
		if !options[i].Disabled && strings.HasPrefix(strings.ToLower(options[i].Label), prefix) {
			return i
		}
	}
	return -1
}

func selectLabel[T comparable](options []SelectOption[T], v T) (string, bool) {
	for _, o := range options {
		if o.Value == v {
			return o.Label, true
		}
	}
	return "", false
}

func selectOr(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

// selectField styles tr as a field that opens a popup, with content and a chevron.
func selectField(c *ui.Context, tr *ui.Element, open bool, content func()) {
	k := Tokens(c)
	d := Density(c)
	tr.MinHeight(d.ControlHeight()).MinWidth(180).Padding(0, d.Unit()*2, 0, d.Unit()*3).Gap(d.Unit() * 2).
		Justify(ui.Start).Radius(theme.ControlRadius).Background(k.Surface)
	internal.Ring(tr, k.Focus, theme.ControlRadius)
	border := k.ControlBorder
	switch {
	case open:
		border = k.AccentText
	case tr.Hovered():
		border = k.Text
	}
	tr.Border(theme.BorderWidth, border)
	tr.Children(func() {
		content()
		ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			if open {
				// Mirrored: the chevron points up while the popup shows.
				r.Y, r.H = r.Y+r.H, -r.H
			}
			internal.Chevron(p, r, false, k.TextMuted)
		})
	})
}

// selectClear draws a small clear button that takes no Tab stop; it reports a click.
func selectClear(c *ui.Context, k theme.Tokens, name string) bool {
	x := ui.Box(c).Size(18, 18).Shrink(0).Radius(theme.ControlRadius).Role(ui.RoleButton).Label(name)
	clicked := x.Clicked()
	if x.Hovered() {
		x.Background(k.SurfaceHover)
	}
	x.Draw(func(p *ui.Painter, r ui.Rect) { internal.Cross(p, r, k.TextMuted) })
	return clicked
}

// selectPanel styles a popup panel.
func selectPanel(c *ui.Context, panel *ui.Element) {
	k := Tokens(c)
	panel.Margin(4, 0, 0, 0).Padding(4).Radius(theme.PopupRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.14))
}

// selectRow styles an option row: highlight background, check mark when chosen.
func selectRow(c *ui.Context, row *ui.Element, highlighted, chosen bool, label string) {
	k := Tokens(c)
	d := Density(c)
	row.MinHeight(d.ControlHeight()-4).Padding(0, d.Unit()*3, 0, d.Unit()*2).Gap(d.Unit() * 2).
		AlignItems(ui.Center).Radius(theme.ControlRadius)
	if highlighted {
		row.Background(k.SurfaceHover)
	}
	if chosen {
		row.TextColor(k.AccentText)
	}
	row.Children(func() {
		ui.Box(c).Size(14, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			if chosen {
				var path ui.Path
				path.MoveTo(r.X+r.W*0.15, r.Y+r.H*0.52).LineTo(r.X+r.W*0.4, r.Y+r.H*0.76).LineTo(r.X+r.W*0.85, r.Y+r.H*0.26)
				p.StrokePath(&path, 1.75, k.AccentText)
			}
		})
		ui.Text(c, label).SingleLine()
	})
}

// selectNote shows a muted line in a popup, as an empty state.
func selectNote(c *ui.Context, s string) {
	d := Density(c)
	ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).AlignItems(ui.Center).Children(func() {
		ui.Text(c, s).TextColor(Tokens(c).TextMuted)
	})
}
