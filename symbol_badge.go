package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MarketSymbol is a tradable instrument as the caller lists it.
type MarketSymbol struct {
	Code, Name string
	// Market is the venue, such as NASDAQ or SSE.
	Market string
	Type   string
}

// SymbolBadgeOptions configure a SymbolBadge.
type SymbolBadgeOptions struct {
	// Clickable makes the badge a button reporting clicks.
	Clickable bool
	Disabled  bool
}

// SymbolBadge shows a code with its market and type tags.
func SymbolBadge(c *ui.Context, s MarketSymbol, opts SymbolBadgeOptions) ChoiceResult {
	if s.Code == "" {
		panic("mujica: SymbolBadge needs a Code")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var e *ui.Element
	wrap := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		if opts.Clickable {
			e = ui.ButtonBase(c)
			internal.Ring(e, k.Focus, theme.ControlRadius)
			if e.Hovered() {
				e.Background(k.SurfaceHover)
			}
			e.Transition(colorFade)
		} else {
			e = ui.Row(c).Role(ui.RoleText)
		}
	})
	e.MinHeight(theme.MinHitSize).Gap(u).AlignItems(ui.Center).Padding(0, u).Radius(theme.ControlRadius).Shrink(0).
		Label(symbolBadgeName(s))
	e.Children(func() {
		ui.Text(c, s.Code).Font(theme.MonoFont).FontWeight(700).TextColor(k.Text).SingleLine()
		for _, tag := range []string{s.Market, s.Type} {
			if tag != "" {
				ui.Row(c).Padding(1, u).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.Border).Background(k.Background).Children(func() {
					ui.Text(c, tag).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				})
			}
		}
	})
	clicked := opts.Clickable && e.Clicked() && !e.IsDisabled()
	if clicked {
		slog.Debug("mujica: symbol_badge click", "code", s.Code)
	}
	return ChoiceResult{Element: e, changed: clicked}
}

// symbolBadgeName reads code, name, market and type.
func symbolBadgeName(s MarketSymbol) string {
	var parts []string
	for _, p := range []string{s.Code, s.Name, s.Market, s.Type} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ", ")
}
