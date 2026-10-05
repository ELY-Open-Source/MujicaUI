package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// BarsLoader shows four bars rising and falling in turn.
func BarsLoader(c *ui.Context, opts LoaderOptions) *ui.Element {
	return loaderMark(c, opts, 1.25, 1000, func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool) {
		w := r.W / 7
		for i := range 4 {
			h := float32(0.6)
			if !still {
				h = 0.3 + 0.7*loaderWave(t-float64(i)/8)
			}
			bh := internal.Snap(p, r.H*h)
			x := r.X + float32(i)*w*2
			p.Fill(ui.Rect{X: x, Y: internal.Snap(p, r.Y+(r.H-bh)/2), W: w, H: bh}, col, 1)
		}
	})
}
