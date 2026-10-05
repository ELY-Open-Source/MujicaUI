package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// CapabilityCard explains one thing the assistant can do, with example prompts.
type CapabilityCard struct {
	Title, Description string
	Icon               *ui.SVG
	Examples           []string
}

// CapabilityCardsOptions configure CapabilityCards.
type CapabilityCardsOptions struct {
	// Columns is the most cards per row; zero is 3.
	Columns  int
	Disabled bool
}

// CapabilityCardsView is built CapabilityCards and the example clicked this frame.
type CapabilityCardsView struct {
	Element *ui.Element
	// Example is the example prompt chosen, empty for none.
	Example string
}

// CapabilityCards shows what the assistant can do; clicking an example hands it to the caller.
func CapabilityCards(c *ui.Context, cards []CapabilityCard, opts CapabilityCardsOptions) CapabilityCardsView {
	if opts.Columns < 0 {
		panic("mujica: CapabilityCards Columns is negative")
	}
	cols := opts.Columns
	if cols == 0 {
		cols = 3
	}
	k := Tokens(c)
	d := Density(c)
	g := ui.Grid(c).Columns(cols).Gap(d.Unit() * 3).Disabled(opts.Disabled)
	v := CapabilityCardsView{Element: g}
	g.Children(func() {
		for i, card := range cards {
			if card.Title == "" {
				panic("mujica: CapabilityCard needs a Title")
			}
			box := ui.Column(c).Key(i).Gap(d.Unit()*2).Padding(d.Unit()*4).Background(k.Surface).
				Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).Role(ui.RoleGroup).Label(card.Title)
			box.Children(func() {
				ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
					if card.Icon != nil {
						ui.Icon(c, card.Icon).Size(20, 20).TextColor(k.AccentText)
					}
					ui.Text(c, card.Title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text)
				})
				if card.Description != "" {
					ui.Text(c, card.Description).TextColor(k.TextMuted)
				}
				ui.Column(c).Gap(d.Unit()).FocusGroup(ui.Vertical).Children(func() {
					for _, ex := range card.Examples {
						if capabilityExample(c, k, ex) {
							v.Example = ex
						}
					}
				})
			})
		}
	})
	if v.Example != "" {
		slog.Debug("mujica: capability_cards example", "len", len(v.Example))
	}
	return v
}

// capabilityExample is a quoted example prompt button; it reports a click.
func capabilityExample(c *ui.Context, k theme.Tokens, text string) bool {
	d := Density(c)
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(d.Unit(), d.Unit()*3).Justify(ui.Start).
		Radius(theme.ControlRadius).Background(k.Background).Border(theme.BorderWidth, k.Border).Label(text)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover).Border(theme.BorderWidth, k.ControlBorder)
	}
	b.Transition(colorFade)
	b.Children(func() { ui.Text(c, "“"+text+"”").TextColor(k.Text) })
	return b.Clicked() && !b.IsDisabled()
}
