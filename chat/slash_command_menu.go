package chat

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SlashCommand is one command of a SlashCommandMenu.
type SlashCommand struct {
	// Name is typed after the slash; required, without spaces.
	Name        string
	Description string
}

// SlashCommandResult carries the popup and chosen command.
type SlashCommandResult struct {
	Element ui.Element
	chosen  SlashCommand
	ok      bool
}

// Chosen returns the command picked by Enter or a click.
func (r SlashCommandResult) Chosen() (SlashCommand, bool) { return r.chosen, r.ok }

// SlashCommandMenu filters commands after a slash trigger.
// Arrows navigate; Enter chooses; Tab completes; Escape dismisses.
func SlashCommandMenu(c *ui.Context, composer input.InputResult, value *string, commands []SlashCommand) SlashCommandResult {
	if !composer.Element.Valid() || !composer.Input.Valid() || value == nil {
		panic("mujica: SlashCommandMenu needs composer and value")
	}
	data := ui.Local(composer.Element, "slash-data", func() slashMenuData { return slashMenuData{} })
	query, typing := strings.CutPrefix(*value, "/")
	typing = typing && !strings.ContainsAny(query, " \t\n")
	if !data.built || query != data.query || typing != data.typing || !slices.Equal(data.commands, commands) {
		for _, cmd := range commands {
			if cmd.Name == "" || strings.ContainsAny(cmd.Name, " \t\n") {
				panic("mujica: SlashCommand Name must be one word")
			}
		}
		data.built, data.query, data.typing = true, query, typing
		data.commands = append(data.commands[:0], commands...)

		data.shown, data.rows = data.shown[:0], data.rows[:0]
		if typing {
			q := strings.ToLower(query)
			for pass := range 2 {
				for _, cmd := range commands {
					n := strings.ToLower(cmd.Name)
					if pass == 0 && strings.HasPrefix(n, q) || pass == 1 && !strings.HasPrefix(n, q) && strings.Contains(n, q) {
						data.shown = append(data.shown, cmd)
						data.rows = append(data.rows, slashMenuRow{label: "/" + cmd.Name, detail: cmd.Description})
					}
				}
			}
		}
	}
	var r SlashCommandResult
	st := ui.Local(composer.Element, slashMenuKey{"slash"}, func() slashMenuState { return slashMenuState{} })
	shown, rows := data.shown, data.rows
	i, how, panel := slashMenuShow(c, composer, st, query, typing, rows, slashMenuOptions{empty: core.Msg(c, "slash_command_menu.empty"), label: core.Msg(c, "slash_command_menu.label")})
	r.Element = panel
	if i >= 0 {
		*value = "/" + shown[i].Name + " "
		if how == slashMenuChose {
			r.chosen, r.ok = shown[i], true
		}
		slog.Debug("mujica: slash_command_menu pick", "name", shown[i].Name, "chosen", r.ok)
	}
	return r
}

type slashMenuData struct {
	commands, shown []SlashCommand
	query           string
	built, typing   bool
	rows            []slashMenuRow
}

type slashMenuKey struct{ name string }

// slashMenuRow represents an item or titled group.
type slashMenuRow struct {
	label, detail, icon, group string
}

type slashMenuOptions struct {
	empty, label string
	loading      bool
}

// slashMenuState tracks navigation, dismissal and pending choices.
type slashMenuState struct {
	hi, lastHi int
	query      string
	dismissed  *string
	pick       slashMenuHow
	held, open bool
	list       ui.ListState
}

type slashMenuHow int

const (
	slashMenuNone slashMenuHow = iota
	slashMenuChose
	slashMenuCompleted
)

// slashMenuShow returns the selected row and activation method.
func slashMenuShow(c *ui.Context, composer input.InputResult, st *slashMenuState, query string, open bool, rows []slashMenuRow, o slashMenuOptions) (int, slashMenuHow, ui.Element) {
	if !open {
		st.dismissed = nil
	}
	if query != st.query {
		st.hi, st.query = 0, query
		st.list.ScrollTo(0, ui.Start)
		if st.dismissed != nil && *st.dismissed != query {
			st.dismissed = nil
		}
	}
	// Panel presses preserve the menu while transferring focus.
	open = open && st.dismissed == nil && !composer.Element.IsDisabled() && (composer.Input.Focused() || st.held)
	pick, how := -1, slashMenuNone
	st.open = open
	if !open {
		st.pick = slashMenuNone
		return pick, how, ui.Element{}
	}
	st.hi = min(max(st.hi, 0), len(rows)-1)
	if st.pick != slashMenuNone && st.hi >= 0 {
		pick, how = st.hi, st.pick
	}
	st.pick = slashMenuNone
	n := len(rows)
	promptComposerHook(composer, func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputKeyDown || ev.Mods != 0 {
			return false
		}
		switch ev.Key {
		case ui.KeyDown, ui.KeyUp:
			if n > 0 {
				d := map[ui.Key]int{ui.KeyDown: 1, ui.KeyUp: -1}[ev.Key]
				st.hi = (max(st.hi, 0) + d + n) % n
			}
		case ui.KeyEnter:
			st.pick = slashMenuChose
		case ui.KeyTab:
			st.pick = slashMenuCompleted
		case ui.KeyEscape:
			q := query
			st.dismissed = &q
		default:
			return false
		}
		c.Invalidate()
		return true
	})
	i, panel := slashMenuPanel(c, composer.Element.Bounds(), st, rows, o)
	if i >= 0 {
		pick, how = i, slashMenuChose
	}
	if pick >= 0 {
		composer.Input.Focus()
	}
	return pick, how, panel
}

// slashMenuPanel builds visible rows and returns clicked choices.
func slashMenuPanel(c *ui.Context, a ui.Rect, st *slashMenuState, rows []slashMenuRow, o slashMenuOptions) (int, ui.Element) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	rowH := core.Density(c).ControlHeight()
	headH := rowH * 0.75
	h := 2 * u
	group := ""
	for _, r := range rows {
		if r.group != group {
			h, group = h+headH, r.group
		}
		h += rowH
	}
	if len(rows) == 0 {
		h += rowH
	}
	ww, wh := c.Size()
	h = min(h+2*theme.BorderWidth, 280, max(0, wh-8))
	w := min(max(a.W*0.6, 320), max(a.W, 200), max(0, ww-8))
	x, y := core.PlaceOverlay(c, a, w, h, u)
	clicked, pressed := -1, false
	var result ui.Element
	ui.Overlay(c, func() {
		panel := ui.Column(c).Absolute().Left(x).Top(y).Width(w).Height(h).Role(ui.RoleList).Label(o.label)
		core.OverlayStyle(c, panel, overlay.OverlayAnchored(c, panel, a))
		result = panel
		panel.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputPointerDown {
				st.held = true
			}
			return false
		})
		panel.Children(func() {
			if st.lastHi != st.hi {
				st.lastHi = st.hi
				st.list.ScrollIntoView(max(st.hi, 0))
			}
			ui.List(c, &st.list, max(1, len(rows)), func(i int) {
				if len(rows) == 0 {
					ui.Row(c).MinHeight(rowH).Padding(0, u*3).Gap(u * 2).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
						if o.loading {
							core.Spinner(c, core.SpinnerOptions{Size: 14})
						}
						ui.Text(c, o.empty).SingleLine()
					})
					return
				}
				r := rows[i]
				ui.Column(c).Children(func() {
					if r.group != "" && (i == 0 || r.group != rows[i-1].group) {
						ui.Text(c, r.group).Height(headH).Padding(u*2, u*3, 0, u*3).FontSize(core.FontSize(c, theme.CaptionSize)).
							FontWeight(600).TextColor(k.TextMuted).Role(ui.RoleText)
					}
					row := ui.Row(c.Key(i)).MinHeight(rowH).Padding(0, u*3, 0, u*2).Gap(u * 2).AlignItems(ui.Center).
						Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(r.label).Cursor(ui.CursorPointer).Selected(i == st.hi).TextColor(k.Text)
					if row.Clicked() {
						clicked = i
					}
					pressed = pressed || row.Pressed()
					hi := i == st.hi
					if hi {
						row.Background(k.Selection).Draw(func(p *ui.Painter, b ui.Rect) {
							p.Fill(ui.Rect{X: b.X, Y: b.Y + 4, W: 2, H: b.H - 8}, k.AccentText, 1)
						})
					} else if row.Hovered() {
						row.Background(k.SurfaceHover)
					}
					row.Children(func() {
						if r.icon != "" {
							ui.Icon(c, icons.Must(r.icon)).Size(16, 16).Shrink(0).TextColor(k.TextMuted)
						}
						ui.Text(c, r.label).SingleLine().Shrink(0).FontWeight(600)
						ui.Text(c, r.detail).SingleLine().MinWidth(0).Shrink(1).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
					})
				})
			}).Fill().Padding(u)
		})
	})
	st.held = st.held && pressed
	return clicked, result
}
