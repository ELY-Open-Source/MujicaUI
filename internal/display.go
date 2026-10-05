package internal

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Fade animates e's background and border colors over d; zero skips it.
func Fade(e *ui.Element, d time.Duration) *ui.Element {
	if d > 0 {
		e.Transition(ui.ElementTransition{Colors: true, Duration: d})
	}
	return e
}

// Rule draws a 1 DIP horizontal line from x0 to x1 at y, on device pixels.
func Rule(p *ui.Painter, x0, x1, y float32, c ui.Color) {
	s := float64(p.Scale())
	px := math.Floor(float64(y) * s)
	if int(math.Round(s))%2 == 1 {
		px += 0.5
	}
	yy := float32(px / s)
	p.Line(x0, yy, x1, yy, 1, c)
}
