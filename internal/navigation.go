package internal

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Ease moves e's background, border and text color fg to their targets over d.
func Ease(e *ui.Element, d time.Duration, fg ui.Color) *ui.Element {
	if d > 0 {
		e.Transition(ui.ElementTransition{Colors: true, Duration: d})
	}
	ch := func(key string, v uint8) uint8 {
		return uint8(math.Round(float64(e.Animate(key, float32(v), d))))
	}
	return e.TextColor(ui.Color{R: ch("mujica.fg.r", fg.R), G: ch("mujica.fg.g", fg.G), B: ch("mujica.fg.b", fg.B), A: ch("mujica.fg.a", fg.A)})
}

// PopIn fades popup e in over d from the frame it appeared and returns the progress, 0 to 1.
func PopIn(c *ui.Context, e *ui.Element, d time.Duration) float32 {
	born := ui.Local(e, "mujica.popin", func() time.Time { return c.Now() })
	p := float32(1)
	if t := float32(c.Now().Sub(*born)) / float32(max(d, 1)); d > 0 && t < 1 {
		p = ui.EaseOut(t)
		c.AnimationFrame()
	}
	e.Opacity(p)
	return p
}
