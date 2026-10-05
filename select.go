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

// Select chooses *value from options in a popup; typing on the closed select jumps to a match.
func Select[T comparable](c *ui.Context, value *T, options []SelectOption[T], opts SelectOptions) ChoiceResult {
	k := Tokens(c)
	old := *value
	// SelectBase chooses into a copy as it is built; a choice counts once the option is known enabled.
	pick := old
	var s *ui.SelectParts[T]
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() { s = ui.SelectBase(c, &pick) })
	tr := s.Trigger
	commit := func() {
		if pick != *value && !tr.IsDisabled() && selectEnabled(options, pick) {
			*value = pick
		}
		pick = *value
	}
	commit()
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
		if !s.Open() && !tr.IsDisabled() {
			if i := selectFind(options, *value, ta.prefix); i >= 0 {
				*value = options[i].Value
			}
		}
	}
	var zero T
	label, has := selectLabel(options, *value)
	if opts.Clearable && has && (tr.Shortcut(0, ui.KeyBackspace) || tr.Shortcut(0, ui.KeyDelete)) && !tr.IsDisabled() {
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
	pick = *value
	lastHi := ui.Local(tr, selectKey{"shown"}, func() int { return -1 })
	if opts.Disabled {
		// SelectParts has no way to close it (ui/base.go:310): a disabled select builds no popup.
		return ChoiceResult{tr, false}
	}
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
	commit()
	if !s.Open() {
		*lastHi = -1
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: select changed", "label", opts.Label, "value", *value)
	}
	return ChoiceResult{tr, changed}
}

// selectEnabled reports whether v is an option that is not disabled.
func selectEnabled[T comparable](options []SelectOption[T], v T) bool {
	for _, o := range options {
		if o.Value == v {
			return !o.Disabled
		}
	}
	return false
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
			internal.Mark(p, r, selectChevron(open), k.TextMuted)
		})
	})
}

// selectChevron names the arrow of a field: up while its popup shows.
func selectChevron(open bool) string {
	if open {
		return "chevron-up"
	}
	return "chevron-down"
}

// selectClear draws a clear cross outside the Tab order, hit area 28 DIP laid out as 18; it reports a click.
func selectClear(c *ui.Context, k theme.Tokens, name string) bool {
	x := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Margin(-5).Shrink(0).Radius(theme.ControlRadius).
		Role(ui.RoleButton).Label(name).Cursor(ui.CursorPointer)
	clicked := x.Clicked()
	if x.Hovered() {
		x.Background(k.SurfaceHover)
	}
	x.Draw(func(p *ui.Painter, r ui.Rect) { internal.Mark(p, selectCenter(r, 14), "x", k.TextMuted) })
	return clicked
}

// selectCenter returns the size×size square centered in r.
func selectCenter(r ui.Rect, size float32) ui.Rect {
	return ui.Rect{X: r.X + (r.W-size)/2, Y: r.Y + (r.H-size)/2, W: size, H: size}
}

// selectCheck reserves room for the chosen mark and draws it when on.
func selectCheck(c *ui.Context, on bool) {
	k := Tokens(c)
	ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
		if on {
			internal.Mark(p, r, "check", k.AccentText)
		}
	})
}

// selectPanel styles a popup panel.
func selectPanel(c *ui.Context, panel *ui.Element) {
	k := Tokens(c)
	panel.Margin(4, 0, 0, 0).Padding(4).Radius(theme.PopupRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.14))
}

// selectRow styles an option row: chosen has the selection fill and a check, highlighted a side bar.
func selectRow(c *ui.Context, row *ui.Element, highlighted, chosen bool, label string) {
	k := Tokens(c)
	d := Density(c)
	row.MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3, 0, d.Unit()*2).Gap(d.Unit() * 2).
		AlignItems(ui.Center).Radius(theme.ControlRadius).Role(ui.RoleListItem)
	switch {
	case chosen:
		row.Background(k.Selection).TextColor(k.AccentText)
	case highlighted:
		row.Background(k.SurfaceHover)
	}
	if highlighted {
		row.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 2, H: r.H - 8}, k.AccentText, 1)
		})
	}
	row.Children(func() {
		selectCheck(c, chosen)
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
