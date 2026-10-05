package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// SpinnerOptions configure a Spinner. Size is 16 DIPs when zero.
type SpinnerOptions struct {
	Size float32
	// Label is read out and shown beside the mark; empty uses "Loading" unread on screen.
	Label     string
	ShowLabel bool
	// Color paints the mark; zero uses AccentText.
	Color ui.Color
}

// Spinner draws four diamonds turning around a center while it shows.
func Spinner(c *ui.Context, opts SpinnerOptions) *ui.Element {
	k := Tokens(c)
	size := opts.Size
	if size <= 0 {
		size = 16
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "loading")
	}
	color := opts.Color
	if color.A == 0 {
		color = k.AccentText
	}
	still := c.Preferences().ReduceMotion
	row := ui.Row(c).Gap(Density(c).Unit() * 2).AlignItems(ui.Center).Role(ui.RoleProgress).Label(label)
	row.Children(func() {
		ui.Box(c).Size(size, size).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			turn := 0.0
			if !still {
				turn = float64(p.Now().UnixMilli()%1200) / 1200
				p.AnimationFrame()
			}
			cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W*0.32
			for i := range 4 {
				a := 2 * math.Pi * (turn + float64(i)/4)
				x := cx + rad*float32(math.Cos(a))
				y := cy + rad*float32(math.Sin(a))
				alpha := 1 - float32(i)*0.22
				p.FillPath(internal.Diamond(x, y, r.W*0.3), color.Alpha(alpha))
			}
		})
		if opts.ShowLabel {
			ui.Text(c, label).TextColor(k.TextMuted).SingleLine()
		}
	})
	return row
}
