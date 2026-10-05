package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MultiSelectOptions configure a MultiSelect.
type MultiSelectOptions struct {
	// Placeholder shows while nothing is chosen; empty uses "Select…".
	Placeholder string
	// Max limits how many values may be chosen; zero is no limit.
	Max int
	// MaxTags is how many chosen values show as tags before "+N"; zero is 3.
	MaxTags  int
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// multiSelectKey keeps this file's state apart from MyGo's own.
type multiSelectKey struct{ name string }

// MultiSelect chooses several *values from options in a popup with a
// search field, "Select all" for the matches, and an optional limit.
// Down from the search moves into the list; Space checks; Escape closes.
func MultiSelect[T comparable](c *ui.Context, values *[]T, options []SelectOption[T], opts MultiSelectOptions) ChoiceResult {
	if opts.Max < 0 || opts.MaxTags < 0 {
		panic("mujica: MultiSelect Max and MaxTags must not be negative")
	}
	k := Tokens(c)
	before := slices.Clone(*values)
	tr := ui.ButtonBase(c).Disabled(opts.Disabled).Role(ui.RolePopUpButton)
	if opts.Label != "" {
		tr.Label(opts.Label)
	}
	open := ui.Local(tr, multiSelectKey{"open"}, func() bool { return false })
	query := ui.Local(tr, multiSelectKey{"query"}, func() string { return "" })
	if tr.Clicked() || (!*open && tr.Shortcut(0, ui.KeyDown)) {
		*open = !*open
		*query = ""
	}
	if n := len(*values); n > 0 && tr.Shortcut(0, ui.KeyBackspace) {
		*values = slices.Clone((*values)[:n-1])
	}
	selectField(c, tr, *open, func() {
		multiSelectTags(c, values, options, opts)
	})
	ui.PopoverBase(c, tr, open, func(panel *ui.Element) {
		selectPanel(c, panel)
		panel.MinWidth(tr.Bounds().W).MaxWidth(max(tr.Bounds().W, 320))
		multiSelectPopup(c, k, values, options, query, opts)
	})
	changed := !slices.Equal(before, *values)
	if changed {
		slog.Debug("mujica: multi select changed", "label", opts.Label, "count", len(*values))
	}
	return ChoiceResult{tr, changed}
}

// multiSelectTags shows the chosen values as removable tags, collapsing the rest into "+N".
func multiSelectTags[T comparable](c *ui.Context, values *[]T, options []SelectOption[T], opts MultiSelectOptions) {
	k := Tokens(c)
	if len(*values) == 0 {
		ui.Text(c, selectOr(opts.Placeholder, msg(c, "select.placeholder"))).SingleLine().TextColor(k.TextMuted).Grow(1).Shrink(1)
		return
	}
	limit := opts.MaxTags
	if limit == 0 {
		limit = 3
	}
	ui.Row(c).Gap(4).Grow(1).Shrink(1).ClipX().Children(func() {
		for i, v := range *values {
			if i == limit {
				break
			}
			label, _ := selectLabel(options, v)
			ui.Row(c).Key(v).Shrink(0).Children(func() {
				if multiSelectTag(c, label, opts.Disabled) {
					*values = slices.Delete(slices.Clone(*values), i, i+1)
				}
			})
		}
		if n := len(*values) - limit; n > 0 {
			ui.Text(c, fmt.Sprintf(msg(c, "multi_select.more"), n)).TextColor(k.TextMuted).SingleLine().Shrink(0)
		}
	})
}

// multiSelectTag draws a chosen value with a remove cross outside the Tab order; it reports a removal.
func multiSelectTag(c *ui.Context, label string, disabled bool) bool {
	k := Tokens(c)
	removed := false
	ui.Row(c).Gap(2).Padding(1, 2, 1, 6).Radius(theme.ControlRadius).Background(k.Selection).
		Border(theme.BorderWidth, k.Border).Shrink(0).Children(func() {
		ui.Text(c, label).FontSize(theme.CaptionSize).TextColor(k.Text).SingleLine()
		if !disabled {
			removed = selectClear(c, k, fmt.Sprintf(msg(c, "multi_select.remove"), label))
		}
	})
	return removed
}

func multiSelectPopup[T comparable](c *ui.Context, k theme.Tokens, values *[]T, options []SelectOption[T], query *string, opts MultiSelectOptions) {
	d := Density(c)
	var shown []SelectOption[T]
	q := strings.ToLower(strings.TrimSpace(*query))
	for _, o := range options {
		if strings.Contains(strings.ToLower(o.Label), q) {
			shown = append(shown, o)
		}
	}
	full := opts.Max > 0 && len(*values) >= opts.Max
	into := false
	ui.Column(c).Gap(4).Children(func() {
		search := selectSearch(c, query, msg(c, "multi_select.search"))
		into = search.Shortcut(0, ui.KeyDown)
		if len(shown) == 0 {
			selectNote(c, msg(c, "multi_select.no_results"))
			return
		}
		list := ui.Column(c).FocusGroup(ui.Vertical)
		list.Children(func() {
			all, some := true, false
			for _, o := range shown {
				on := slices.Contains(*values, o.Value)
				all = all && (on || o.Disabled)
				some = some || on
			}
			state := internal.Unchecked
			if all && some {
				state = internal.Checked
			} else if some {
				state = internal.Mixed
			}
			tick := state == internal.Checked
			ui.Column(c).PaddingX(4).Children(func() {
				b := multiSelectCheck(c, k, &tick, state, msg(c, "multi_select.all"), opts.Max > 0 && !some && full)
				if into {
					b.Focus()
				}
				if b.Changed() {
					*values = multiSelectAll(*values, shown, tick, opts.Max)
				}
			})
			ui.Box(c).Height(1).Background(k.Border).MarginY(2)
			ui.Scroll(c).MaxHeight(240).Padding(4).Children(func() {
				for _, o := range shown {
					on := slices.Contains(*values, o.Value)
					st := internal.Unchecked
					if on {
						st = internal.Checked
					}
					row := ui.Row(c).Key(o.Value)
					row.Children(func() {
						cb := multiSelectCheck(c, k, &on, st, o.Label, o.Disabled || (full && !on))
						if cb.Changed() {
							*values = multiSelectSet(*values, o.Value, on)
						}
					})
				}
			})
		})
		if opts.Max > 0 {
			ui.Text(c, fmt.Sprintf(msg(c, "multi_select.limit"), len(*values), opts.Max)).
				FontSize(theme.CaptionSize).TextColor(k.TextMuted).Padding(2, d.Unit()*2)
		}
	})
}

// multiSelectCheck is a CheckboxBase row drawing its own box.
func multiSelectCheck(c *ui.Context, k theme.Tokens, on *bool, st internal.CheckState, label string, disabled bool) *ui.Element {
	d := Density(c)
	scratch := *on
	if disabled {
		on = &scratch
	}
	cb := ui.CheckboxBase(c, on).Disabled(disabled).MinHeight(d.ControlHeight()-4).Padding(0, d.Unit()*2).
		Gap(d.Unit() * 2).AlignItems(ui.Center).Radius(theme.ControlRadius).Grow(1)
	internal.Ring(cb, k.Focus, theme.ControlRadius)
	hover := cb.Hovered()
	if hover {
		cb.Background(k.SurfaceHover)
	}
	cb.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.CheckBox(p, r, st, k, false)
		})
		ui.Text(c, label).SingleLine()
	})
	return cb
}

// multiSelectSet adds or removes v.
func multiSelectSet[T comparable](values []T, v T, on bool) []T {
	if on {
		return append(slices.Clone(values), v)
	}
	return slices.DeleteFunc(slices.Clone(values), func(x T) bool { return x == v })
}

// multiSelectAll checks every enabled shown option up to max, or unchecks them all.
func multiSelectAll[T comparable](values []T, shown []SelectOption[T], on bool, limit int) []T {
	for _, o := range shown {
		if o.Disabled || slices.Contains(values, o.Value) == on {
			continue
		}
		if on && limit > 0 && len(values) >= limit {
			break
		}
		values = multiSelectSet(values, o.Value, on)
	}
	return values
}

// selectSearch is the search field at the top of a popup; it takes the focus as the popup opens.
func selectSearch(c *ui.Context, query *string, placeholder string) *ui.Element {
	k := Tokens(c)
	d := Density(c)
	f := ui.Row(c).Gap(d.Unit()*2).Padding(0, d.Unit()*2).MinHeight(d.ControlHeight()-4).AlignItems(ui.Center).
		Radius(theme.ControlRadius).Background(k.Background).Border(theme.BorderWidth, k.ControlBorder)
	var in *ui.Element
	f.Children(func() {
		ui.Box(c).Size(14, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.Mark(p, r, "search", k.TextMuted)
		})
		in = ui.TextInputBase(c, query).Grow(1).Placeholder(placeholder).Label(placeholder).AutoFocus()
	})
	if in.Focused() {
		f.Border(theme.BorderWidth, k.AccentText)
	}
	return in
}
