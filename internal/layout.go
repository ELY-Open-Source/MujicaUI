package internal

import "github.com/egoist/mygo/ui"

// Remeasure rebuilds the view when layout gives e another size than the w×h it was built for.
func Remeasure(c *ui.Context, e ui.Element, w, h float32) ui.Element {
	return e.DrawOver(func(_ *ui.Painter, r ui.Rect) {
		if abs(r.W-w) > 0.5 || abs(r.H-h) > 0.5 {
			c.Invalidate()
		}
	})
}

func abs(v float32) float32 { return max(v, -v) }
