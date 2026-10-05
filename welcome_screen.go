package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// WelcomeScreenOptions configure a WelcomeScreen. Greeting is required.
type WelcomeScreenOptions struct {
	Greeting string
	Subtitle string
	Cards    []CapabilityCard
	// Prompts is the slot for starter prompts, such as SuggestionChips.
	Prompts  func()
	Disabled bool
}

// WelcomeScreen greets a new conversation with capability cards and a starter-prompt slot; it reports a chosen example.
func WelcomeScreen(c *ui.Context, opts WelcomeScreenOptions) CapabilityCardsView {
	if opts.Greeting == "" {
		panic("mujica: WelcomeScreen needs a Greeting")
	}
	k := Tokens(c)
	d := Density(c)
	var v CapabilityCardsView
	col := ui.Column(c).FillWidth().MaxWidth(880).AlignSelf(ui.Center).AlignItems(ui.Center).Gap(d.Unit()*4).
		Padding(d.Unit()*8, d.Unit()*4).Disabled(opts.Disabled)
	col.Children(func() {
		ui.Text(c, opts.Greeting).Font(theme.SerifFont).FontSize(FontSize(c, theme.H1Size)).TextColor(k.Text).
			TextAlign(ui.Center).Role(ui.RoleText)
		ui.Box(c).Width(220).Height(9).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.RuleDiamond(p, r, 8, 6, k.Ornament)
		})
		if opts.Subtitle != "" {
			ui.Text(c, opts.Subtitle).TextColor(k.TextMuted).TextAlign(ui.Center)
		}
		if len(opts.Cards) > 0 {
			v = CapabilityCards(c, opts.Cards, CapabilityCardsOptions{Columns: min(3, len(opts.Cards)), Disabled: opts.Disabled})
			v.Element.FillWidth()
		}
		if opts.Prompts != nil {
			ui.Column(c).FillWidth().AlignItems(ui.Center).Children(opts.Prompts)
		}
	})
	v.Element = col
	return v
}
