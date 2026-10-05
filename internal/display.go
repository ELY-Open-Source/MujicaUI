package internal

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// Rule draws a pixel-aligned 1 DIP horizontal line at y.
func Rule(p *ui.Painter, x0, x1, y float32, c ui.Color) {
	yy := ruleY(p, y)
	p.Line(x0, yy, x1, yy, 1, c)
}

// ruleY moves y to where a 1 DIP line fills whole device pixels.
func ruleY(p *ui.Painter, y float32) float32 {
	s := float64(p.Scale())
	px := math.Floor(float64(y) * s)
	if int(math.Round(s))%2 == 1 {
		px += 0.5
	}
	return float32(px / s)
}

// RuleDiamond draws a rule across the middle of r, broken by gap on each side of a size diamond.
func RuleDiamond(p *ui.Painter, r ui.Rect, gap, size float32, col ui.Color) {
	cx, cy := r.X+r.W/2, ruleY(p, r.Y+r.H/2)
	p.Line(r.X, cy, cx-gap, cy, 1, col)
	p.Line(cx+gap, cy, r.X+r.W, cy, 1, col)
	p.FillPath(Diamond(Snap(p, cx), Snap(p, cy), size), col)
}

// Frame strokes a 1 DIP line inset inside r with a size diamond on each corner.
func Frame(p *ui.Painter, r ui.Rect, inset, size float32, col ui.Color) {
	f := ui.Rect{X: r.X + inset, Y: r.Y + inset, W: r.W - 2*inset, H: r.H - 2*inset}
	p.Stroke(f, col, 1, 1)
	for _, pt := range [4][2]float32{{f.X, f.Y}, {f.X + f.W, f.Y}, {f.X, f.Y + f.H}, {f.X + f.W, f.Y + f.H}} {
		p.FillPath(Diamond(Snap(p, pt[0]), Snap(p, pt[1]), size), col)
	}
}
