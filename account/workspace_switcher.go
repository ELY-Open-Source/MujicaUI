package account

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Workspace is a place a user can switch to.
type Workspace struct {
	ID, Name string
	// Detail supplies the workspace summary.
	Detail string
}

// WorkspaceSwitcherOptions configure a WorkspaceSwitcher.
type WorkspaceSwitcherOptions struct {
	Disabled bool
}

type workspaceSwitcherState struct {
	open     bool
	query    string
	cursor   int
	source   []Workspace
	hits     []int
	filtered string
	ready    bool
}

// WorkspaceSwitcher filters workspaces and changes the current selection.
func WorkspaceSwitcher(c *ui.Context, current *string, spaces []Workspace, opts WorkspaceSwitcherOptions) core.ChoiceResult {
	var cur *Workspace
	for i := range spaces {
		if spaces[i].ID == "" || spaces[i].Name == "" {
			panic("mujica: Workspace needs an ID and a Name")
		}
		if spaces[i].ID == *current {
			cur = &spaces[i]
		}
	}
	if cur == nil {
		panic(fmt.Sprintf("mujica: WorkspaceSwitcher current %q is not a workspace", *current))
	}
	var r core.ChoiceResult
	b := accountTrigger(c, cur.Name, cur.Detail, nil, opts.Disabled)
	r.Element = b
	st := ui.Local(b, "ws", func() workspaceSwitcherState { return workspaceSwitcherState{} })
	disabled := b.IsDisabled()
	if disabled {
		st.open = false
	} else if b.Clicked() || b.Shortcut(0, ui.KeyDown) {
		st.open, st.query, st.cursor = true, "", 0
	}
	choose := func(i int) {
		if !disabled && spaces[i].ID != *current {
			*current, r = spaces[i].ID, core.NewChoiceResult(r.Element, true)
			slog.Debug("mujica: workspace_switcher switched")
		}
		st.open = false
	}
	overlay.Popover(c, b, &st.open, overlay.PopoverOptions{Width: 300, Label: core.Msg(c, "workspace_switcher.label")}, func() {
		res := input.SearchInput(c, &st.query, input.SearchInputOptions{Placeholder: core.Msg(c, "workspace_switcher.search"), Label: core.Msg(c, "workspace_switcher.search")})
		in := res.Input
		core.AutoFocus(c, in)
		q := strings.ToLower(strings.TrimSpace(st.query))
		if !st.ready || st.filtered != q || !slices.Equal(st.source, spaces) {
			st.source, st.filtered, st.ready = slices.Clone(spaces), q, true
			st.hits, st.cursor = st.hits[:0], 0
			for i, w := range spaces {
				if q == "" || strings.Contains(strings.ToLower(w.Name+" "+w.Detail), q) {
					st.hits = append(st.hits, i)
				}
			}
		}
		hits := st.hits
		st.cursor = max(0, min(st.cursor, len(hits)-1))
		before := st.cursor
		switch {
		case disabled || in.IsDisabled() || len(hits) == 0:
		case in.Shortcut(0, ui.KeyDown):
			st.cursor = (st.cursor + 1) % len(hits)
		case in.Shortcut(0, ui.KeyUp):
			st.cursor = (st.cursor + len(hits) - 1) % len(hits)
		}
		if res.Submitted() && len(hits) > 0 {
			choose(hits[st.cursor])
		}
		if len(hits) == 0 {
			ui.Text(c, core.Msg(c, "workspace_switcher.none")).TextColor(core.Tokens(c).TextMuted).Padding(core.Density(c).Unit() * 2).Role(ui.RoleStatus)
		}
		list := accountRows(c, len(hits), func(n int) any { return spaces[hits[n]].ID }, func(n int) {
			i := hits[n]
			row := workspaceRow(c, spaces[i], spaces[i].ID == *current, n == st.cursor)
			if row.Clicked() && !row.IsDisabled() {
				choose(i)
			}
		})
		if len(hits) > 8 && (st.cursor != before || res.Changed()) {
			rows := ui.Local(list, "rows", func() ui.ListState { return ui.ListState{} })
			rows.ScrollIntoView(st.cursor)
		}
	})
	return r
}

// workspaceRow displays a clickable workspace with a selection mark.
func workspaceRow(c *ui.Context, w Workspace, current, cursor bool) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	b := ui.ButtonBase(c.Key(w.ID)).MinHeight(d.ControlHeight()).Padding(d.Unit(), d.Unit()*2).Gap(d.Unit() * 2).
		Justify(ui.Start).Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(w.Name).Selected(current)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case cursor || current:
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).Children(func() {
			if current {
				ui.Icon(c, icons.Must("check")).Size(16, 16).TextColor(k.AccentText)
			}
		})
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			ui.Text(c, w.Name).SingleLine().Ellipsis("…")
			if w.Detail != "" {
				ui.Text(c, w.Detail).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
			}
		})
	})
	return b
}
