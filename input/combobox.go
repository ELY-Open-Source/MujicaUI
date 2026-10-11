package input

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ComboboxOptions configure a Combobox.
type ComboboxOptions struct {
	Placeholder string
	// AllowCustom lets typed text that matches no option become the value.
	AllowCustom bool
	// ExternalFilter shows options as given: the caller already filtered them for Query.
	ExternalFilter bool
	// Loading shows that the caller is fetching options for Query.
	Loading  bool
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// ComboboxResult is a Combobox's field, its value change and its query.
type ComboboxResult struct {
	core.ChoiceResult
	query        string
	queryChanged bool
}

// Query returns the text in the field.
func (r ComboboxResult) Query() string { return r.query }

// QueryChanged reports whether the user typed in this frame: fetch new options then.
func (r ComboboxResult) QueryChanged() bool { return r.queryChanged }

// comboboxKey keeps this file's state apart from MyGo's own.
type comboboxKey struct{ name string }

// Combobox chooses *value by typing to filter options and picking one, or typed text with AllowCustom.
func Combobox(c *ui.Context, value *string, options []string, opts ComboboxOptions) ComboboxResult {
	k := core.Tokens(c)
	d := core.Density(c)
	old := *value
	f := ui.Row(c).MinHeight(d.ControlHeight()).MinWidth(200).AlignItems(ui.Center).Padding(0, d.Unit()*2, 0, 0).
		Gap(d.Unit() * 2).Radius(theme.ControlRadius).Background(k.Surface).Cursor(ui.CursorText).Disabled(opts.Disabled)
	query := ui.Local(f, comboboxKey{"query"}, func() string { return *value })
	wasFocused := ui.Local(f, comboboxKey{"focused"}, func() bool { return false })
	// MyGo keeps the caret on replaced text (ui/editor.go:117): a choice builds a fresh input.
	gen := ui.Local(f, comboboxKey{"input"}, func() int { return 0 })
	refocus := ui.Local(f, comboboxKey{"refocus"}, func() bool { return false })
	var p *ui.ComboboxParts
	f.Children(func() {
		// A focused ComboboxBase takes typing under a disabled wrapper, so it edits a copy.
		target := query
		if opts.Disabled {
			cp := *query
			target = &cp
		}
		ui.Row(c).Key(*gen).Grow(1).Children(func() {
			p = ui.ComboboxBase(c, target)
		})
		if *refocus {
			p.Input.Focus()
			*refocus = false
		}
		if opts.Disabled {
			p.SetOpen(false)
		}
		pad := max(0, (d.ControlHeight()-2-theme.BodySize*theme.LineHeight)/2)
		p.Input.Grow(1).MinWidth(80).Padding(pad, d.Unit()*3).Placeholder(opts.Placeholder)
		if opts.Label != "" {
			p.Input.Label(opts.Label)
		}
		if opts.Loading {
			core.Spinner(c, core.SpinnerOptions{Size: 14})
		}
		arrow := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0).Cursor(ui.CursorPointer).
			Role(ui.RoleButton).Label(core.Msg(c, "combobox.toggle"))
		if arrow.Clicked() {
			p.SetOpen(!p.Open())
			p.Input.Focus()
		}
		open := p.Open()
		arrow.Draw(func(pt *ui.Painter, r ui.Rect) {
			internal.Mark(pt, selectCenter(r, 16), selectChevron(open), k.TextMuted)
		})
	})
	in := p.Input
	if f.Clicked() {
		in.Focus()
	}
	border := k.ControlBorder
	switch {
	case in.Focused():
		border = k.AccentText
	case f.Hovered():
		border = k.Text
	}
	f.Border(theme.BorderWidth, border)
	queryChanged := in.Changed()
	shown := comboboxShown(options, *query, p.Filtering() && !opts.ExternalFilter)
	custom := ""
	if q := strings.TrimSpace(*query); opts.AllowCustom && p.Filtering() && q != "" && !comboboxHas(shown, q) {
		custom = q
	}
	lastHi := ui.Local(f, comboboxKey{"shown"}, func() int { return -1 })
	if !p.Open() {
		*lastHi = -1
	}
	p.Popup(func(panel *ui.Element) {
		SelectPanel(c, panel)
		ui.Scroll(c).MaxHeight(240).Children(func() {
			if custom != "" {
				it := p.Item(custom)
				SelectRow(c, it, it.Highlighted(), false, fmt.Sprintf(core.Msg(c, "combobox.use"), custom))
			}
			for i, o := range shown {
				it := p.Item(o)
				if it.Highlighted() && *lastHi != i {
					*lastHi = i
					it.ScrollIntoView()
				}
				SelectRow(c, it, it.Highlighted(), o == *value, o)
			}
		})
	})
	if p.Open() && len(shown) == 0 && custom == "" {
		comboboxStatus(c, f.Bounds(), opts.Loading)
	}
	focused := in.Focused()
	if in.IsDisabled() {
		// Input delivered as the field turned disabled commits nothing.
		*query = *value
	} else if v, ok := p.Chosen(); ok {
		*value, *query = v, v
		*gen, *refocus = *gen+1, focused
	} else if in.Submitted() && opts.AllowCustom {
		*value = strings.TrimSpace(*query)
	} else if in.Submitted() && p.Open() && len(shown) > 0 {
		// Results came after typing left nothing highlighted: Enter takes the first.
		*value, *query = shown[0], shown[0]
		*gen, *refocus = *gen+1, true
		p.SetOpen(false)
	} else if *wasFocused && !focused {
		if opts.AllowCustom {
			*value = strings.TrimSpace(*query)
		}
		*query = *value
	} else if !focused && *query != *value {
		// The caller set the value.
		*query = *value
	}
	*wasFocused = focused
	changed := *value != old
	if changed {
		slog.Debug("mujica: combobox changed", "label", opts.Label, "value", *value)
	}
	return ComboboxResult{core.NewChoiceResult(f, changed), *query, queryChanged}
}

// comboboxShown returns the options containing query, those starting with it first.
func comboboxShown(options []string, query string, filter bool) []string {
	q := strings.ToLower(strings.TrimSpace(query))
	if !filter || q == "" {
		return options
	}
	var starts, contains []string
	for _, o := range options {
		switch l := strings.ToLower(o); {
		case strings.HasPrefix(l, q):
			starts = append(starts, o)
		case strings.Contains(l, q):
			contains = append(contains, o)
		}
	}
	return append(starts, contains...)
}

func comboboxHas(options []string, s string) bool {
	for _, o := range options {
		if strings.EqualFold(o, s) {
			return true
		}
	}
	return false
}

// comboboxStatus shows loading or no matches under the field r; it takes no input.
func comboboxStatus(c *ui.Context, r ui.Rect, loading bool) {
	ui.Overlay(c, func() {
		panel := ui.Row(c).Absolute().Left(r.X).Top(r.Y + r.H).MinWidth(r.W).PassThrough().Role(ui.RoleStatus)
		SelectPanel(c, panel)
		panel.Children(func() {
			if loading {
				ui.Row(c).Gap(8).Padding(8, 12).AlignItems(ui.Center).Children(func() {
					core.Spinner(c, core.SpinnerOptions{})
					ui.Text(c, core.Msg(c, "combobox.loading")).TextColor(core.Tokens(c).TextMuted)
				})
				return
			}
			SelectNote(c, core.Msg(c, "combobox.no_results"))
		})
	})
}
