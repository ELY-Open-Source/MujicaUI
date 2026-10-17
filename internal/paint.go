// Package internal holds drawing shared by MujicaUI components.
package internal

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// FocusRing strokes the keyboard ring outside r, leaving the focus gap.
func FocusRing(p *ui.Painter, r ui.Rect, radius float32, c ui.Color) {
	o := theme.FocusGap + theme.FocusWidth
	p.Stroke(ui.Rect{X: r.X - o, Y: r.Y - o, W: r.W + 2*o, H: r.H + 2*o}, c, radius+o, theme.FocusWidth)
}

// Ring replaces MyGo's ring on e with the MujicaUI ring.
func Ring(e ui.Element, c ui.Color, radius float32) ui.Element {
	return e.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if e.FocusVisible() {
			FocusRing(p, r, radius, c)
		}
	})
}

// Diamond returns a rhombus of size s centered at (cx, cy).
func Diamond(cx, cy, s float32) *ui.Path {
	var p ui.Path
	h := s / 2
	p.MoveTo(cx, cy-h)
	p.LineTo(cx+h, cy)
	p.LineTo(cx, cy+h)
	p.LineTo(cx-h, cy)
	p.Close()
	return &p
}

// Snap rounds v to MyGo's quarter-device-pixel path cache grid.
func Snap(p *ui.Painter, v float32) float32 {
	g := 4 * p.Scale()
	return float32(math.Round(float64(v*g))) / g
}
