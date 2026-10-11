package feedback

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
)

// WaveLoader draws a sine line flowing to the right.
func WaveLoader(c *ui.Context, opts LoaderOptions) *ui.Element {
	return loaderMark(c, opts, "wave_loader", 2.5, 1400, func(p *ui.Painter, r ui.Rect, t float64, col ui.Color, still bool) {
		var path ui.Path
		amp := r.H * 0.3
		const steps = 24
		for i := 0; i <= steps; i++ {
			f := float64(i) / steps
			x := r.X + r.W*float32(f)
			y := r.Y + r.H/2 - amp*float32(math.Sin(2*math.Pi*(1.5*f-t)))
			x, y = internal.Snap(p, x), internal.Snap(p, y)
			if i == 0 {
				path.MoveTo(x, y)
			} else {
				path.LineTo(x, y)
			}
		}
		p.StrokePath(&path, max(1.5, r.H/10), col)
	})
}
