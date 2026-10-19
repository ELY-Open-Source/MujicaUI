package input

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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

// MultiSelect chooses several *values from options in a popup with search and "Select all".
func MultiSelect[T comparable](c *ui.Context, values *[]T, options []SelectOption[T], opts MultiSelectOptions) core.ChoiceResult {
	if opts.Max < 0 || opts.MaxTags < 0 {
		panic("mujica: MultiSelect Max and MaxTags must not be negative")
	}
	k := core.Tokens(c)
	before := slices.Clone(*values)
	tr := ui.ButtonBase(c).Disabled(opts.Disabled).Role(ui.RolePopUpButton)
	if opts.Label != "" {
		tr.Label(opts.Label)
	}
	open := ui.Local(tr, multiSelectKey{"open"}, func() bool { return false })
	if opts.Disabled {
		*open = false
	}
	query := ui.Local(tr, multiSelectKey{"query"}, func() string { return "" })
	if tr.Clicked() || (!*open && tr.Shortcut(0, ui.KeyDown) && !tr.IsDisabled()) {
		*open = !*open
		*query = ""
	}
	if n := len(*values); n > 0 && tr.Shortcut(0, ui.KeyBackspace) && !tr.IsDisabled() {
		*values = slices.Clone((*values)[:n-1])
	}
	SelectField(c, tr, *open, func() {
		MultiSelectTags(c, values, options, opts)
	})
	ui.PopoverBase(c, tr, open, func(panel ui.Element) {
		SelectPanel(c, panel)
		panel.MinWidth(tr.Bounds().W).MaxWidth(max(tr.Bounds().W, 320))
		multiSelectPopup(c, k, values, options, query, opts)
	})
	changed := !slices.Equal(before, *values)
	if changed {
		slog.Debug("mujica: multi select changed", "label", opts.Label, "count", len(*values))
	}
	return core.NewChoiceResult(tr, changed)
}

// MultiSelectTags shows the chosen values as removable tags, collapsing the rest into "+N".
func MultiSelectTags[T comparable](c *ui.Context, values *[]T, options []SelectOption[T], opts MultiSelectOptions) {
	k := core.Tokens(c)
	if len(*values) == 0 {
		ui.Text(c, selectOr(opts.Placeholder, core.Msg(c, "select.placeholder"))).SingleLine().TextColor(k.TextMuted).Grow(1).Shrink(1)
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
			label, _ := SelectLabel(options, v)
			ui.Row(c.Key(v)).Shrink(0).Children(func() {
				if multiSelectTag(c, label, opts.Disabled) {
					*values = slices.Delete(slices.Clone(*values), i, i+1)
				}
			})
		}
		if n := len(*values) - limit; n > 0 {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "multi_select.more"), n)).TextColor(k.TextMuted).SingleLine().Shrink(0)
		}
	})
}

// multiSelectTag draws a chosen value with a remove cross outside the Tab order; it reports a removal.
func multiSelectTag(c *ui.Context, label string, disabled bool) bool {
	k := core.Tokens(c)
	removed := false
	ui.Row(c).Gap(2).Padding(1, 2, 1, 6).Radius(theme.ControlRadius).Background(k.Selection).
		Border(theme.BorderWidth, k.Border).Shrink(0).Children(func() {
		ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text).SingleLine()
		if !disabled {
			removed = selectClear(c, k, fmt.Sprintf(core.Msg(c, "multi_select.remove"), label))
		}
	})
	return removed
}

func multiSelectPopup[T comparable](c *ui.Context, k theme.Tokens, values *[]T, options []SelectOption[T], query *string, opts MultiSelectOptions) {
	d := core.Density(c)
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
		search := SelectSearch(c, query, core.Msg(c, "multi_select.search"))
		into = search.Shortcut(0, ui.KeyDown)
		if len(shown) == 0 {
			SelectNote(c, core.Msg(c, "multi_select.no_results"))
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
				b := MultiSelectCheck(c, k, &tick, state, core.Msg(c, "multi_select.all"), opts.Max > 0 && !some && full)
				if into {
					b.Focus()
				}
				if b.Changed() {
					*values = MultiSelectAll(*values, shown, tick, opts.Max)
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
					row := ui.Row(c.Key(o.Value))
					row.Children(func() {
						cb := MultiSelectCheck(c, k, &on, st, o.Label, o.Disabled || (full && !on))
						if cb.Changed() {
							*values = MultiSelectSet(*values, o.Value, on)
						}
					})
				}
			})
		})
		if opts.Max > 0 {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "multi_select.limit"), len(*values), opts.Max)).
				FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(2, d.Unit()*2)
		}
	})
}

// MultiSelectCheck is a CheckboxBase row drawing its own box.
func MultiSelectCheck(c *ui.Context, k theme.Tokens, on *bool, st internal.CheckState, label string, disabled bool) ui.Element {
	d := core.Density(c)
	var cb ui.Element
	// A disabled wrapper blocks the base, which toggles as it is built.
	ui.Row(c).Grow(1).Disabled(disabled).Children(func() {
		cb = ui.CheckboxBase(c, on).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).
			Gap(d.Unit() * 2).AlignItems(ui.Center).Radius(theme.ControlRadius).Grow(1)
	})
	if st == internal.Mixed {
		// MyGo sets the mixed state only inside its own widgets (ui/element.go:282).
		cb.Description(core.Msg(c, "multi_select.mixed"))
	}
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

// MultiSelectSet adds or removes v.
func MultiSelectSet[T comparable](values []T, v T, on bool) []T {
	if on {
		return append(slices.Clone(values), v)
	}
	return slices.DeleteFunc(slices.Clone(values), func(x T) bool { return x == v })
}

// MultiSelectAll checks every enabled shown option up to max, or unchecks them all.
func MultiSelectAll[T comparable](values []T, shown []SelectOption[T], on bool, limit int) []T {
	for _, o := range shown {
		if o.Disabled || slices.Contains(values, o.Value) == on {
			continue
		}
		if on && limit > 0 && len(values) >= limit {
			break
		}
		values = MultiSelectSet(values, o.Value, on)
	}
	return values
}

// SelectSearch is the search field at the top of a popup; it takes the focus as the popup opens.
func SelectSearch(c *ui.Context, query *string, placeholder string) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	f := ui.Row(c).Gap(d.Unit()*2).Padding(0, d.Unit()*2).MinHeight(d.ControlHeight()-4).AlignItems(ui.Center).
		Radius(theme.ControlRadius).Background(k.Background).Border(theme.BorderWidth, k.ControlBorder)
	var in ui.Element
	f.Children(func() {
		ui.Box(c).Size(14, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.Mark(p, r, "search", k.TextMuted)
		})
		in = core.AutoFocus(c, ui.TextInputBase(c, query).Grow(1).Placeholder(placeholder).Label(placeholder))
	})
	if in.Focused() {
		f.Border(theme.BorderWidth, k.AccentText)
	}
	return in
}
