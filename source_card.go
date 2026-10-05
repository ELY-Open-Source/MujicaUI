package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SourceCard shows source n's title, site, time and snippet; Clicked opens it.
// It lights up with Selection and an accent edge while st holds n Active or hovered.
func SourceCard(c *ui.Context, n int, src Source, st *CitationState) *ui.Element {
	if n < 1 || st == nil || src.Title == "" || src.URL == "" {
		panic("mujica: SourceCard needs n >= 1, a state and a Source with Title and URL")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	b := ui.ButtonBase(c).Padding(u*2, u*3).AlignItems(ui.Start).Justify(ui.Start).Radius(theme.CardRadius).
		Label(fmt.Sprintf(msg(c, "source_card.open"), src.Title)).Tooltip(src.URL).Transition(colorFade)
	internal.Ring(b, k.Focus, theme.CardRadius)
	st.track(c, b, n)
	switch {
	case st.lit(n):
		b.Background(k.Selection).Border(theme.BorderWidth, k.Border).BorderWidth(theme.BorderWidth, theme.BorderWidth, theme.BorderWidth, 3).
			BorderColor(k.AccentText)
	default:
		b.Background(k.Surface).Border(theme.BorderWidth, k.Border)
	}
	b.Children(func() { sourceSummary(c, n, src) })
	if b.Clicked() && !b.IsDisabled() {
		slog.Debug("mujica: source_card open", "n", n)
	}
	return b
}
