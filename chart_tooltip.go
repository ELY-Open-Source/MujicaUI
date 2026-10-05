package mujica

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartTooltipLine is one row of a tooltip; a zero Color shows no marker.
type ChartTooltipLine struct {
	Text  string
	Color ui.Color
	Shape int
}

// ChartTooltip floats painted lines beside window point (x, y), avoiding edges.
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
	left, top := x+gap, y+gap
	if left+w > ww-4 {
		left = x - gap - w
	}
	if top+h > wh-4 {
		top = y - gap - h
	}
	left, top = max(4, left), max(4, top)
	var tip *ui.Element
	ui.Overlay(c, func() {
		tip = ui.Box(c).Absolute().Left(left).Top(top).Size(w, h).PassThrough().Role(ui.RoleNone)
		overlayStyle(c, tip, 1)
		tip.Draw(func(p *ui.Painter, r ui.Rect) {
			at := r.Y + pad
			for i, l := range lines {
				tx := r.X + pad
				if l.Color.A > 0 {
					internal.Marker(p, l.Shape, tx+5, at+lh[i]/2, 9, l.Color)
					tx += 14
				}
				p.RichText(tx, at, 0, ui.Span{Text: l.Text, Size: size, Weight: chartTipWeight(i), Color: k.Text})
				at += lh[i]
			}
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

// chartTipLines splits a readout at " · " into a heading and rows.
func chartTipLines(readout string) []ChartTooltipLine {
	var out []ChartTooltipLine
	for _, s := range strings.Split(readout, " · ") {
		out = append(out, ChartTooltipLine{Text: s})
	}
	return out
}
