package core

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/internal"
)

// ProgressArc is the arc from turn a to turn b, clockwise from twelve o'clock.
func ProgressArc(pt *ui.Painter, cx, cy, r float32, a, b float32) *ui.Path {
	var p ui.Path
	n := max(2, int(math.Ceil(float64((b-a)*64))))
	for i := 0; i <= n; i++ {
		t := float64(a + (b-a)*float32(i)/float32(n))
		x := cx + r*float32(math.Sin(2*math.Pi*t))
		y := cy - r*float32(math.Cos(2*math.Pi*t))
		if i == 0 {
			p.MoveTo(internal.Snap(pt, x), internal.Snap(pt, y))
		} else {
			p.LineTo(internal.Snap(pt, x), internal.Snap(pt, y))
		}
	}
	return &p
}

func Abs32(v float32) float32 { return float32(math.Abs(float64(v))) }
