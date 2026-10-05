package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SourceCard displays source metadata; Clicked requests opening.
// Selection and an edge mark identify active or hovered sources.
func SourceCard(c *ui.Context, n int, src Source, st *CitationState) *ui.Element {
	if n < 1 || st == nil || src.Title == "" || src.URL == "" {
		panic("mujica: SourceCard needs n >= 1, a state and a Source with Title and URL")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	b := ui.ButtonBase(c).Padding(u*2, u*3).AlignItems(ui.Start).Justify(ui.Start).Radius(theme.CardRadius).
		Label(fmt.Sprintf(core.Msg(c, "source_card.open"), src.Title)).Tooltip(src.URL).Transition(core.ColorFade)
	internal.Ring(b, k.Focus, theme.CardRadius)
	st.track(c, b, n)
	frame := c.Now()
	b.Draw(func(_ *ui.Painter, _ ui.Rect) { st.expire(c, frame) })
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
