package navigation

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Command is an action of a CommandPalette.
type Command struct {
	ID    string
	Label string
	// Group heads the command in the list; commands keep their order within groups.
	Group string
	// Keywords also match the search.
	Keywords []string
	Icon     *ui.SVG
}

// CommandPaletteOptions configure a CommandPalette.
type CommandPaletteOptions struct {
	// OpenMods and OpenKey open the palette; zero key means Cmd+K (Ctrl+K off macOS).
	OpenMods ui.Modifiers
	OpenKey  ui.Key
}

// CommandPaletteResult is the built CommandPalette; Element is nil while closed.
type CommandPaletteResult struct {
	Element ui.Element
	chosen  string
}

// Chosen returns the ID of the command run in this frame; running closes the palette.
func (r CommandPaletteResult) Chosen() (string, bool) { return r.chosen, r.chosen != "" }

type commandPaletteKey struct{ open *bool }

type commandPaletteState struct {
	query   string
	hl      int
	ptr     [2]float32
	wasOpen bool
	reveal  bool
}

// CommandPalette is a searchable command dialog opened by a shortcut or by setting *open.
func CommandPalette(c *ui.Context, open *bool, commands []Command, opts CommandPaletteOptions) CommandPaletteResult {
	seen := map[string]bool{}
	for _, cmd := range commands {
		if cmd.ID == "" || cmd.Label == "" {
			panic("mujica: Command needs an ID and a Label")
		}
		if seen[cmd.ID] {
			panic("mujica: duplicate Command ID " + cmd.ID)
		}
		seen[cmd.ID] = true
	}
	mods, key := opts.OpenMods, opts.OpenKey
	if key == 0 {
		mods, key = ui.Cmd, ui.KeyK
	}
	st := ui.Local(c.Root(), commandPaletteKey{open}, func() commandPaletteState { return commandPaletteState{} })
	if c.Shortcut(mods, key) && !*open {
		*open = true
	}
	if *open && !st.wasOpen {
		st.query, st.hl = "", 0
		slog.Debug("mujica: command palette open")
	}
	st.wasOpen = *open
	res := CommandPaletteResult{}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	res.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		ww, wh := c.Size()
		p := core.OverlayEnter(c, panel)
		back.Background(ui.RGBA(0, 0, 0, 0.4*p)).Justify(ui.Start)
		core.OverlayStyle(c, panel, p)
		panel.Margin(wh*0.12, 0, 0, 0).Width(560).MaxWidth(ww - 32).MaxHeight(wh * 0.76).
			AlignItems(ui.Stretch).Label(core.Msg(c, "command_palette.label"))
		found := commandPaletteFilter(commands, st.query)
		var in ui.Element
		ui.Row(c).AlignItems(ui.Center).Gap(u*3).Padding(u*3, u*4).Children(func() {
			ui.Icon(c, icons.Must("search")).Size(16, 16).TextColor(k.TextMuted)
			in = core.AutoFocus(c, ui.TextInputBase(c, &st.query).FillWidth()).Label(core.Msg(c, "command_palette.search"))
		})
		ui.Box(c).Height(1).Background(k.Ornament).Role(ui.RoleNone)
		if in.Changed() {
			st.hl, st.reveal = 0, true
		}
		n := len(found)
		switch {
		case n == 0:
		case in.Shortcut(0, ui.KeyDown):
			st.hl, st.reveal = (st.hl+1)%n, true
		case in.Shortcut(0, ui.KeyUp):
			st.hl, st.reveal = (st.hl-1+n)%n, true
		}
		st.hl = max(0, min(st.hl, n-1))
		// Focus stays in the search box, so the highlighted command is read out.
		if n > 0 {
			in.Description(found[st.hl].Label)
			if st.reveal {
				c.Announce(found[st.hl].Label)
			}
		} else {
			in.Description(core.Msg(c, "command_palette.empty"))
		}
		run := -1
		if in.Submitted() && n > 0 {
			run = st.hl
		}
		list := ui.Scroll(c).Padding(u * 2).Gap(2).MaxHeight(360).Role(ui.RoleList)
		list.Children(func() {
			if n == 0 {
				ui.Text(c, core.Msg(c, "command_palette.empty")).TextColor(k.TextMuted).Padding(u*6, u*2).TextAlign(ui.Center).FillWidth()
				return
			}
			gutter := slices.ContainsFunc(found, func(cmd Command) bool { return cmd.Icon != nil })
			group := "\x00"
			for i, cmd := range found {
				if cmd.Group != group {
					group = cmd.Group
					if group != "" {
						ui.Text(c, group).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).
							TextColor(k.Ornament.Mix(k.Text, 0.35)).Padding(u*2, u*2, u, u*2).SingleLine()
					}
				}
				if commandPaletteRow(c, cmd, i, st, gutter) {
					run = i
				}
			}
		})
		ui.Row(c).Gap(u*4).Padding(u*2, u*4).Background(k.Background).Children(func() {
			for _, hint := range []string{"command_palette.hint_move", "command_palette.hint_run", "command_palette.hint_close"} {
				ui.Text(c, core.Msg(c, hint)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			}
		})
		if run >= 0 {
			res.chosen = found[run].ID
			*open = false
			slog.Debug("mujica: command palette run", "id", res.chosen)
		}
	})
	if !*open {
		st.wasOpen = false
	}
	return res
}

// commandPaletteRow builds row i and reports a click on it.
func commandPaletteRow(c *ui.Context, cmd Command, i int, st *commandPaletteState, gutter bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	row := ui.Row(c).AlignItems(ui.Center).Gap(u*3).MinHeight(d.ControlHeight()).Padding(0, u*3).
		Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(cmd.Label).TextColor(k.Text).Cursor(ui.CursorPointer)
	on := i == st.hl
	if on {
		row.Background(k.Selection).FontWeight(600)
		row.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 3, H: r.H - 8}, k.AccentText, 1.5)
		})
		if st.reveal {
			row.ScrollIntoView()
			st.reveal = false
		}
	}
	internal.Ease(row, core.Motion(c, theme.StateDuration), k.Text)
	px, py, over := row.PointerPosition()
	b := row.Bounds()
	if at := [2]float32{b.X + px, b.Y + py}; over && at != st.ptr {
		st.ptr = at
		if !on {
			st.hl = i
			c.Invalidate()
		}
	}
	row.Children(func() {
		if cmd.Icon != nil {
			ui.Icon(c, cmd.Icon).Size(16, 16).TextColor(k.TextMuted)
		} else if gutter {
			ui.Box(c).Size(16, 16).Role(ui.RoleNone)
		}
		ui.Text(c, cmd.Label).SingleLine()
	})
	return row.Clicked()
}

// commandPaletteFilter keeps commands matching every word of query, grouped by first appearance.
func commandPaletteFilter(commands []Command, query string) []Command {
	words := strings.Fields(strings.ToLower(query))
	var groups []string
	byGroup := map[string][]Command{}
	for _, cmd := range commands {
		hay := strings.ToLower(cmd.Label + " " + cmd.Group + " " + strings.Join(cmd.Keywords, " "))
		ok := true
		for _, w := range words {
			if !strings.Contains(hay, w) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if _, seen := byGroup[cmd.Group]; !seen {
			groups = append(groups, cmd.Group)
		}
		byGroup[cmd.Group] = append(byGroup[cmd.Group], cmd)
	}
	var out []Command
	for _, g := range groups {
		out = append(out, byGroup[g]...)
	}
	return out
}
