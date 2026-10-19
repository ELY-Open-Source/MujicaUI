package feedback

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/internal"
)

// OrbitLoader moves a diamond around a ring of diamonds.
func OrbitLoader(c *ui.Context, opts LoaderOptions) ui.Element {
	return loaderMark(c, opts, "orbit_loader", 1, 1600, func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W*0.36
		at := func(turn float64) (float32, float32) {
			a := 2*math.Pi*turn - math.Pi/2
			return internal.Snap(p, cx+rad*float32(math.Cos(a))), internal.Snap(p, cy+rad*float32(math.Sin(a)))
		}
		for i := range 8 {
			x, y := at(float64(i) / 8)
			p.FillPath(internal.Diamond(x, y, r.W*0.12), col.Alpha(0.3))
		}
		x, y := at(t)
		p.FillPath(internal.Diamond(x, y, internal.Snap(p, r.W*0.28)), col)
	})
}
