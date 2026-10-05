package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// SlashCommand is one command of a SlashCommandMenu.
type SlashCommand struct {
	// Name is typed after the slash; required, without spaces.
	Name        string
	Description string
}

// SlashCommandResult is the menu's popup panel, nil while closed, and the command chosen this frame.
type SlashCommandResult struct {
	Element *ui.Element
	chosen  SlashCommand
	ok      bool
}

// Chosen returns the command picked by Enter or a click.
func (r SlashCommandResult) Chosen() (SlashCommand, bool) { return r.chosen, r.ok }

// SlashCommandMenu pops up commands while the composer's text is "/" plus a filter.
// Up and Down move, Enter or a click chooses, Tab completes, Escape closes; typing stays in the composer.
func SlashCommandMenu(c *ui.Context, composer InputResult, value *string, commands []SlashCommand) SlashCommandResult {
	for _, cmd := range commands {
		if cmd.Name == "" || strings.ContainsAny(cmd.Name, " \t\n") {
			panic("mujica: SlashCommand Name must be one word")
		}
	}
	var r SlashCommandResult
	st := ui.Local(composer.Element, slashMenuKey{"slash"}, func() slashMenuState { return slashMenuState{} })
	query, typing := strings.CutPrefix(*value, "/")
	typing = typing && !strings.ContainsAny(query, " \t\n")
	var shown []SlashCommand
	var rows []slashMenuRow
	if typing {
		q := strings.ToLower(query)
		for pass := range 2 {
			for _, cmd := range commands {
				n := strings.ToLower(cmd.Name)
				if pass == 0 && strings.HasPrefix(n, q) || pass == 1 && !strings.HasPrefix(n, q) && strings.Contains(n, q) {
					shown = append(shown, cmd)
					rows = append(rows, slashMenuRow{label: "/" + cmd.Name, detail: cmd.Description})
				}
			}
		}
	}
	i, how := slashMenuShow(c, composer, st, query, typing, rows, slashMenuOptions{empty: msg(c, "slash_command_menu.empty"), label: msg(c, "slash_command_menu.label")})
	r.Element = st.panel
	if i >= 0 {
		*value = "/" + shown[i].Name + " "
		if how == slashMenuChose {
			r.chosen, r.ok = shown[i], true
		}
		slog.Debug("mujica: slash_command_menu pick", "name", shown[i].Name, "chosen", r.ok)
	}
	return r
}

type slashMenuKey struct{ name string }

// slashMenuRow is one line of a composer menu; a group starts a titled section.
type slashMenuRow struct {
	label, detail, icon, group string
}

type slashMenuOptions struct {
	empty, label string
	loading      bool
}

// slashMenuState is a composer menu's highlight, Escape and pending key choice.
type slashMenuState struct {
	hi, lastHi int
	query      string
	dismissed  *string
	pick       slashMenuHow
	panel      *ui.Element
	held, open bool
}

type slashMenuHow int

const (
	slashMenuNone slashMenuHow = iota
	slashMenuChose
	slashMenuCompleted
)

// slashMenuShow runs a composer menu while open and returns the row picked this frame and how.
func slashMenuShow(c *ui.Context, composer InputResult, st *slashMenuState, query string, open bool, rows []slashMenuRow, o slashMenuOptions) (int, slashMenuHow) {
	st.panel = nil
	if query != st.query {
		st.hi, st.query = 0, query
		if st.dismissed != nil && *st.dismissed != query {
			st.dismissed = nil
		}
	}
	// A press on the panel takes the focus from the composer; the press keeps it open.
	open = open && st.dismissed == nil && !composer.Element.IsDisabled() && (composer.Input.Focused() || st.held)
	pick, how := -1, slashMenuNone
	st.open = open
	if !open {
		st.pick = slashMenuNone
		return pick, how
	}
	st.hi = min(st.hi, len(rows)-1)
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
	if i := slashMenuPanel(c, composer.Element.Bounds(), st, rows, o); i >= 0 {
		pick, how = i, slashMenuChose
	}
	if pick >= 0 {
		composer.Input.Focus()
	}
	return pick, how
}

// slashMenuPanel draws the rows over or under the composer a and returns a clicked row.
func slashMenuPanel(c *ui.Context, a ui.Rect, st *slashMenuState, rows []slashMenuRow, o slashMenuOptions) int {
	k := Tokens(c)
	u := Density(c).Unit()
	rowH := Density(c).ControlHeight()
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
	h = min(h, 280) + 2*theme.BorderWidth
	w := min(max(a.W*0.6, 320), max(a.W, 200))
	x, y := overlayPlace(c, a, w, h, u)
	clicked, pressed := -1, false
	ui.Overlay(c, func() {
		panel := ui.Column(c).Absolute().Left(x).Top(y).Width(w).Height(h).Role(ui.RoleList).Label(o.label)
		overlayStyle(c, panel, overlayAnchored(c, panel, a))
		st.panel = panel
		panel.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputPointerDown {
				st.held = true
			}
			return false
		})
		panel.Children(func() {
			ui.Scroll(c).Fill().Padding(u).Children(func() {
				if len(rows) == 0 {
					ui.Row(c).MinHeight(rowH).Padding(0, u*3).Gap(u * 2).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
						if o.loading {
							Spinner(c, SpinnerOptions{Size: 14})
						}
						ui.Text(c, o.empty).SingleLine()
					})
				}
				group := ""
				for i, r := range rows {
					if r.group != group {
						group = r.group
						ui.Text(c, group).Height(headH).Padding(u*2, u*3, 0, u*3).FontSize(FontSize(c, theme.CaptionSize)).
							FontWeight(600).TextColor(k.TextMuted).Role(ui.RoleText)
					}
					row := ui.Row(c).Key(i).Height(rowH).Padding(0, u*3, 0, u*2).Gap(u * 2).AlignItems(ui.Center).
						Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(r.label).Cursor(ui.CursorPointer)
					if row.Clicked() {
						clicked = i
					}
					pressed = pressed || row.Pressed()
					hi := i == st.hi
					if hi {
						row.Background(k.Selection).Draw(func(p *ui.Painter, b ui.Rect) {
							p.Fill(ui.Rect{X: b.X, Y: b.Y + 4, W: 2, H: b.H - 8}, k.AccentText, 1)
						})
						if st.lastHi != i {
							st.lastHi = i
							row.ScrollIntoView()
						}
					} else if row.Hovered() {
						row.Background(k.SurfaceHover)
					}
					row.Children(func() {
						if r.icon != "" {
							ui.Icon(c, icons.Must(r.icon)).Size(16, 16).Shrink(0).TextColor(k.TextMuted)
						}
						ui.Text(c, r.label).SingleLine().Shrink(0).FontWeight(600)
						ui.Text(c, r.detail).SingleLine().MinWidth(0).Shrink(1).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize))
					})
				}
			})
		})
	})
	st.held = st.held && pressed
	return clicked
}
