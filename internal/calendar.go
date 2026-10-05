package internal

import "github.com/egoist/mygo/ui"

// Rule draws a pixel-aligned 1 DIP horizontal line at y.
func Rule(p *ui.Painter, x0, x1, y float32, c ui.Color) {
	yy := ruleY(p, y)
	p.Line(x0, yy, x1, yy, 1, c)
}
