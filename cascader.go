package mujica

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// CascaderNode is one choice at a level of a Cascader.
type CascaderNode struct {
	Value, Label string
	// Leaf marks a node without children.
	Leaf     bool
	Disabled bool
}

// CascaderStatus is how far the caller got loading a level. The zero value is loaded.
type CascaderStatus int

const (
	CascaderLoaded CascaderStatus = iota
	// CascaderIdle asks the Cascader to call Load for the level.
	CascaderIdle
	CascaderLoading
	CascaderFailed
)

// CascaderLevel is what the caller knows of the children under a path.
type CascaderLevel struct {
	Nodes  []CascaderNode
	Status CascaderStatus
	// Err explains a failed load; it shows beside the retry button.
	Err error
}

// CascaderOptions configure a Cascader. Level is required.
type CascaderOptions struct {
	// Level returns the children under path; the empty path is the roots.
	Level func(path []string) CascaderLevel
	// Load fetches the children under path while Idle and on retry; it marks the level Loading.
	Load func(path []string)
	// AnyLevel lets branches be chosen too; by default only leaves are.
	AnyLevel bool
	// Placeholder shows while nothing is chosen; empty uses "Select…".
	Placeholder string
	Disabled    bool
	// Label names the field for assistive technology.
	Label string
}

// cascaderState is the popup's open path and the row the keys move to.
type cascaderState struct {
	open      bool
	active    []string
	requested map[string]bool
	focus     []string
	focusSet  bool
}

type cascaderKey struct{}

// Cascader chooses a path through levels the caller supplies, loading children on demand.
func Cascader(c *ui.Context, path *[]string, opts CascaderOptions) ChoiceResult {
	if opts.Level == nil {
		panic("mujica: Cascader needs Level")
	}
	before := slices.Clone(*path)
	tr := ui.ButtonBase(c).Disabled(opts.Disabled).Role(ui.RolePopUpButton)
	if opts.Label != "" {
		tr.Label(opts.Label)
	}
	st := *ui.Local(tr, cascaderKey{}, func() *cascaderState { return &cascaderState{requested: map[string]bool{}} })
	if opts.Disabled {
		st.open = false
	}
	if tr.Clicked() || (!st.open && tr.Shortcut(0, ui.KeyDown) && !tr.IsDisabled()) {
		st.open = !st.open
		if st.open {
			st.active = slices.Clone(*path)
			st.focus, st.focusSet = slices.Clone(*path), true
			if len(*path) == 0 {
				st.focus = []string{""}
			}
		}
	}
	selectField(c, tr, st.open, func() {
		if len(*path) == 0 {
			ui.Text(c, selectOr(opts.Placeholder, msg(c, "select.placeholder"))).SingleLine().TextColor(Tokens(c).TextMuted).Grow(1)
			return
		}
		ui.Text(c, cascaderLabels(opts, *path)).SingleLine().Grow(1).Shrink(1)
	})
	ui.PopoverBase(c, tr, &st.open, func(panel *ui.Element) {
		selectPanel(c, panel)
		panel.Padding(0)
		ui.Row(c).AlignItems(ui.Stretch).Dividers(theme.BorderWidth, Tokens(c).Border).Children(func() {
			for depth := 0; depth <= len(st.active); depth++ {
				if !cascaderColumn(c, st, path, opts, depth) {
					break
				}
			}
		})
	})
	changed := !slices.Equal(before, *path)
	if changed {
		slog.Debug("mujica: cascader changed", "label", opts.Label, "path", *path)
	}
	return ChoiceResult{tr, changed}
}

// cascaderLabels writes the full path of labels, as "Europe / France / Paris".
func cascaderLabels(opts CascaderOptions, path []string) string {
	labels := make([]string, len(path))
	for i, v := range path {
		labels[i] = v
		for _, n := range opts.Level(path[:i]).Nodes {
			if n.Value == v {
				labels[i] = n.Label
			}
		}
	}
	return strings.Join(labels, " / ")
}

// cascaderColumn builds the column at depth and reports whether a deeper one may follow.
func cascaderColumn(c *ui.Context, st *cascaderState, path *[]string, opts CascaderOptions, depth int) bool {
	k := Tokens(c)
	prefix := slices.Clone(st.active[:depth])
	lv := opts.Level(prefix)
	id := strings.Join(prefix, "\x00")
	if lv.Status == CascaderIdle && !st.requested[id] {
		if opts.Load == nil {
			panic("mujica: Cascader level is Idle but Load is nil")
		}
		st.requested[id] = true
		slog.Debug("mujica: cascader load", "path", prefix)
		opts.Load(prefix)
	}
	if lv.Status != CascaderIdle {
		delete(st.requested, id)
	}
	deeper := false
	col := ui.Column(c).Key(id).Width(180).FocusGroup(ui.Vertical)
	col.Children(func() {
		switch {
		case lv.Status == CascaderIdle || lv.Status == CascaderLoading:
			ui.Row(c).Gap(8).Padding(8).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
				Spinner(c, SpinnerOptions{})
				ui.Text(c, msg(c, "cascader.loading")).TextColor(k.TextMuted)
			})
		case lv.Status == CascaderFailed:
			ui.Column(c).Gap(8).Padding(8).AlignItems(ui.Start).Children(func() {
				text := msg(c, "cascader.failed")
				if lv.Err != nil {
					text += ": " + lv.Err.Error()
				}
				ui.Text(c, text).TextColor(k.Danger)
				if Button(c, msg(c, "cascader.retry"), ButtonOptions{Variant: Secondary, Icon: icons.Must("refresh-cw")}).Clicked() {
					slog.Debug("mujica: cascader retry", "path", prefix)
					opts.Load(prefix)
				}
			})
		case len(lv.Nodes) == 0:
			selectNote(c, msg(c, "cascader.empty"))
		default:
			ui.Scroll(c).MaxHeight(260).Padding(4).Children(func() {
				entry := slices.IndexFunc(lv.Nodes, func(n CascaderNode) bool { return !n.Disabled })
				for i, n := range lv.Nodes {
					on := depth < len(st.active) && st.active[depth] == n.Value
					deeper = deeper || on && !n.Leaf
					cascaderRow(c, st, path, opts, prefix, n, on, i == entry)
				}
			})
		}
	})
	return deeper
}

func cascaderRow(c *ui.Context, st *cascaderState, path *[]string, opts CascaderOptions, prefix []string, n CascaderNode, on, first bool) {
	k := Tokens(c)
	d := Density(c)
	full := append(slices.Clone(prefix), n.Value)
	depth := len(prefix)
	b := ui.ButtonBase(c).Key(n.Value).Disabled(n.Disabled).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).
		Gap(d.Unit() * 2).Justify(ui.Start).Radius(theme.ControlRadius).FillWidth().Role(ui.RoleListItem)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if st.focusSet && (slices.Equal(st.focus, full) || first && slices.Equal(st.focus, append(slices.Clone(prefix), ""))) {
		b.Focus()
		st.focusSet = false
	}
	open := func(focusChild bool) {
		st.active = full
		if focusChild {
			st.focus, st.focusSet = append(slices.Clone(full), ""), true
		}
		if opts.AnyLevel {
			*path = slices.Clone(full)
		}
	}
	clicked := b.Clicked()
	switch {
	case clicked && n.Leaf:
		st.active, *path, st.open = full, full, false
	case clicked:
		open(false)
	case b.Shortcut(0, ui.KeyRight) && !n.Leaf && !b.IsDisabled():
		open(true)
	case b.Shortcut(0, ui.KeyLeft) && depth > 0:
		st.focus, st.focusSet = slices.Clone(prefix), true
	}
	chosen := len(*path) > depth && slices.Equal((*path)[:depth+1], full)
	b.Selected(chosen).Background(ui.Transparent).TextColor(k.Text)
	switch {
	case on:
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	if chosen {
		b.TextColor(k.AccentText)
	}
	b.Children(func() {
		selectCheck(c, chosen)
		ui.Text(c, n.Label).SingleLine().Grow(1).Shrink(1)
		if !n.Leaf {
			ui.Box(c).Size(14, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				internal.Mark(p, r, "chevron-right", k.TextMuted)
			})
		}
	})
}
