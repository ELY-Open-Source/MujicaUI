package mujica

import (
	"github.com/egoist/mygo/ui"
)

// ShimmerTextOptions configure a ShimmerText.
type ShimmerTextOptions struct {
	// Active sweeps the light; false shows plain muted text.
	Active bool
}

// ShimmerText shows a status such as "Thinking" with light sweeping across its letters.
func ShimmerText(c *ui.Context, text string, opts ShimmerTextOptions) *ui.Element {
	if text == "" {
		panic("mujica: ShimmerText needs text")
	}
	k := Tokens(c)
	still := c.Preferences().ReduceMotion
	if !opts.Active || still {
		e := ui.Text(c, text).TextColor(k.TextMuted)
		if opts.Active {
			e.Role(ui.RoleStatus)
		}
		return e
	}
	runes := []rune(text)
	size := c.Theme().FontSize
	// The text keeps its layout and name; the sweep is painted over it.
	e := ui.Text(c, text).TextColor(ui.Transparent).Role(ui.RoleStatus).SingleLine()
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		t := float32(p.Now().UnixMilli()%1800) / 1800
		at := t*float32(len(runes)+6) - 3
		spans := make([]ui.Span, len(runes))
		for i, ch := range runes {
			d := max(0, 1-abs32(float32(i)-at)/3)
			spans[i] = ui.Span{Text: string(ch), Size: size, Color: k.TextMuted.Mix(k.Text, d)}
		}
		p.RichText(r.X, r.Y, 0, spans...)
	})
	return e
}
