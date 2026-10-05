package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// PulseLoader shows one dot breathing between 60% and full size; never a halo.
func PulseLoader(c *ui.Context, opts LoaderOptions) *ui.Element {
	return loaderMark(c, opts, 1, 1600, func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool) {
		s, a := float32(1), float32(1)
		if !still {
			w := loaderWave(t)
			s, a = 0.6+0.4*w, 0.55+0.45*w
		}
		var dot ui.Path
		dot.Circle(r.X+r.W/2, r.Y+r.H/2, internal.Snap(p, r.H*0.4*s))
		p.FillPath(&dot, col.Alpha(a))
	})
}
