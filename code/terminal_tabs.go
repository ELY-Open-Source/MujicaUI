package code

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TerminalTab is one terminal of a TerminalTabs.
type TerminalTab struct {
	ID, Title string
	// Busy asks for confirmation before closing.
	Busy bool
}

// TerminalTabsOptions configure TerminalTabs.
type TerminalTabsOptions struct {
	Disabled bool
}

// TerminalTabsResult is TerminalTabs and what happened in this frame.
type TerminalTabsResult struct {
	Element *ui.Element
	// Renamed and Closed hold a tab ID.
	Changed bool
	Renamed string
	Closed  string
}

type terminalTabsState struct {
	editing, draft, closing string
	confirm                 bool
}

// TerminalTabs shows renamable, closable terminal tabs.
func TerminalTabs(c *ui.Context, tabs *[]TerminalTab, active *string, opts TerminalTabsOptions, content func(id string)) TerminalTabsResult {
	if tabs == nil || active == nil || content == nil {
		panic("mujica: TerminalTabs needs tabs, active and content")
	}
	at := slices.IndexFunc(*tabs, func(t TerminalTab) bool { return t.ID == *active })
	if len(*tabs) > 0 && at < 0 {
		panic("mujica: TerminalTabs active is not a tab ID")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res TerminalTabsResult
	res.Element = ui.Column(c).MinHeight(0).Disabled(opts.Disabled)
	st := ui.Local(res.Element, "terminal_tabs", func() terminalTabsState { return terminalTabsState{} })
	disabled := res.Element.IsDisabled()
	if disabled {
		st.editing, st.confirm = "", false
	}
	remove := func(id string) {
		i := slices.IndexFunc(*tabs, func(t TerminalTab) bool { return t.ID == id })
		if i < 0 {
			return
		}
		*tabs = slices.Delete(*tabs, i, i+1)
		if *active == id && len(*tabs) > 0 {
			*active = (*tabs)[min(i, len(*tabs)-1)].ID
			res.Changed = true
		}
		res.Closed = id
		slog.Debug("mujica: terminal_tabs close", "id", id)
	}
	askClose := func(t TerminalTab) {
		if t.Busy {
			st.closing, st.confirm = t.ID, true
			return
		}
		remove(t.ID)
	}
	res.Element.Children(func() {
		strip := ui.ScrollHorizontal(c).Shrink(0).MinWidth(0).Padding(theme.FocusGap+theme.FocusWidth, theme.FocusGap+theme.FocusWidth, 0, theme.FocusGap+theme.FocusWidth)
		strip.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, k.Border, 0) })
		strip.Children(func() {
			idx := max(at, 0)
			base := ui.TabsBase(c, &idx, len(*tabs))
			base.List.Gap(u).Label(core.Msg(c, "terminal_tabs.label"))
			base.List.Children(func() {
				for i, t := range *tabs {
					rename := func(title string) {
						(*tabs)[i].Title = title
						res.Renamed = t.ID
						slog.Debug("mujica: terminal_tabs rename", "id", t.ID)
					}
					terminalTabsTab(c, base.Tab(i), t, i == at, st, rename, askClose, disabled)
				}
			})
			if base.List.Changed() && !disabled {
				*active = (*tabs)[idx].ID
				res.Changed = true
				slog.Debug("mujica: terminal_tabs select", "id", *active)
			}
		})
		if len(*tabs) == 0 {
			CodeCaption(c, core.Msg(c, "terminal_tabs.empty")).Padding(u * 4)
			return
		}
		ui.Column(c).Grow(1).MinHeight(0).Padding(u*2, 0, 0, 0).Children(func() { content(*active) })
	})
	if st.confirm {
		title := ""
		for _, t := range *tabs {
			if t.ID == st.closing {
				title = t.Title
			}
		}
		r := overlay.AlertDialog(c, &st.confirm, overlay.AlertDialogOptions{Title: fmt.Sprintf(core.Msg(c, "terminal_tabs.confirm"), title),
			Description: core.Msg(c, "terminal_tabs.busy"), ConfirmLabel: core.Msg(c, "terminal_tabs.close")})
		if r.Confirmed && !disabled {
			remove(st.closing)
		}
	}
	return res
}

// terminalTabsTab builds one tab with title and close.
func terminalTabsTab(c *ui.Context, tab *ui.Element, t TerminalTab, on bool, st *terminalTabsState, rename func(string), askClose func(TerminalTab), disabled bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	tab.MinHeight(d.ControlHeight()).Padding(0, u, 0, u*3).Gap(u).AlignItems(ui.Center).Radius(theme.ControlRadius, theme.ControlRadius, 0, 0).Label(t.Title)
	if t.Busy {
		tab.Description(core.Msg(c, "terminal_tabs.running"))
	}
	fg := k.TextMuted
	switch {
	case on:
		tab.Background(k.Selection).FontWeight(600)
		fg = k.Text
	case tab.Hovered():
		tab.Background(k.SurfaceHover)
		fg = k.Text
	}
	internal.Ease(tab, core.Motion(c, theme.StateDuration), fg)
	tab.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, k.AccentText, 0)
		}
		if tab.FocusVisible() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	if !disabled && st.editing == "" && (tab.DoubleClicked() || tab.Shortcut(0, ui.KeyF2)) {
		st.editing, st.draft = t.ID, t.Title
	}
	if !disabled && tab.Shortcut(0, ui.KeyDelete) {
		askClose(t)
	}
	tab.Children(func() {
		if st.editing == t.ID {
			in := input.TextInput(c, &st.draft, input.TextInputOptions{Label: core.Msg(c, "terminal_tabs.rename")})
			in.Element.Width(u * 36)
			in.Input.AutoFocus()
			switch {
			case in.Input.Shortcut(0, ui.KeyEscape):
				st.editing = ""
			case in.Submitted() && strings.TrimSpace(st.draft) != "":
				rename(strings.TrimSpace(st.draft))
				st.editing = ""
				tab.Focus()
			}
			return
		}
		ui.Text(c, t.Title).SingleLine()
		if t.Busy {
			ui.Box(c).Size(6, 6).Shrink(0).Role(ui.RoleNone).Background(k.AccentText).Radius(3)
		}
		if core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "terminal_tabs.closeTab"), t.Title)).Clicked() && !disabled {
			askClose(t)
		}
	})
}
