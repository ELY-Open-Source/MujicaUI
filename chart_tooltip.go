package mujica

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartTooltipLine pairs readout text with an optional marker.
type ChartTooltipLine struct {
	Text  string
	Color ui.Color
	Shape int
}

// ChartTooltip shows readouts beside a window point, avoiding edges.
func ChartTooltip(c *ui.Context, x, y float32, lines []ChartTooltipLine) *ui.Element {
	if len(lines) == 0 {
		panic("mujica: ChartTooltip needs a line")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	size := FontSize(c, theme.CaptionSize)
	pad, gap := u*2, float32(12)
	w, h := float32(0), pad*2
	lh := make([]float32, len(lines))
	for i, l := range lines {
		tw, th := c.MeasureText(0, ui.Span{Text: l.Text, Size: size, Weight: chartTipWeight(i)})
		if l.Color.A > 0 {
			tw += 14
		}
		w, lh[i] = max(w, tw), th
		h += th
	}
	w += pad * 2
	ww, wh := c.Size()
	w, h = min(w, max(1, ww-8)), min(h, max(1, wh-8))
	left, top := overlayPlace(c, ui.Rect{X: x, Y: y}, w, h, gap)
	if x+gap+w <= ww-4 {
		left = max(4, x+gap)
	} else {
		left = max(4, min(x-gap-w, ww-4-w))
	}
	var tip *ui.Element
	ui.Overlay(c, func() {
		tip = ui.Column(c).Absolute().Left(left).Top(top).Size(w, h).Role(ui.RoleNone)
		overlayStyle(c, tip, overlayEnter(c, tip))
		tip.Children(func() {
			inner := max(1, w-2*pad)
			content := float32(0)
			for i, l := range lines {
				width := inner
				if l.Color.A > 0 {
					width = max(1, width-14)
				}
				_, lh[i] = c.MeasureText(width, ui.Span{Text: l.Text, Size: size, Weight: chartTipWeight(i)})
				content += lh[i]
			}
			ui.Scroll(c).Fill().Padding(pad).Children(func() {
				ui.Box(c).Width(inner).Height(content).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
					at := r.Y
					for i, l := range lines {
						tx, width := r.X, inner
						if l.Color.A > 0 {
							internal.Marker(p, l.Shape, tx+5, at+lh[i]/2, 9, l.Color)
							tx, width = tx+14, max(1, width-14)
						}
						p.RichText(tx, at, width, ui.Span{Text: l.Text, Size: size, Weight: chartTipWeight(i), Color: k.Text})
						at += lh[i]
					}
				})
			})
		})
	})
	return tip
}

// chartTipWeight sets the first line, the heading, in semibold.
func chartTipWeight(i int) int {
	if i == 0 {
		return 600
	}
	return 400
}

// chartTipLines splits a readout into tooltip lines.
func chartTipLines(readout string) []ChartTooltipLine {
	var out []ChartTooltipLine
	for _, s := range strings.Split(readout, " · ") {
		out = append(out, ChartTooltipLine{Text: s})
	}
	return out
}
