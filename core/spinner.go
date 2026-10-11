package core

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
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

// Spinner turns four diamonds only while painted; hidden or under reduced motion it stops.
func Spinner(c *ui.Context, opts SpinnerOptions) *ui.Element {
	if !(opts.Size >= 0) || math.IsInf(float64(opts.Size), 1) {
		panic(fmt.Sprintf("mujica: Spinner Size %v invalid", opts.Size))
	}
	k := Tokens(c)
	size := opts.Size
	if size == 0 {
		size = 16
	}
	label := opts.Label
	if label == "" {
		label = Msg(c, "loading")
	}
	color := opts.Color
	if color.A == 0 {
		color = k.AccentText
	}
	still := c.Preferences().ReduceMotion
	row := ui.Row(c).Gap(Density(c).Unit() * 2).AlignItems(ui.Center)
	row.Children(func() {
		// Keyed by motion so a flip repaints; MyGo's spinner carries the progress role.
		ui.Box(c).Key(still).Shrink(0).Children(func() { spinnerMark(c, label, size, color, still) })
		if opts.ShowLabel {
			ui.Text(c, label).TextColor(k.TextMuted).SingleLine()
		}
	})
	return row
}

func spinnerMark(c *ui.Context, label string, size float32, color ui.Color, still bool) {
	ui.Spinner(c).Label(label).Size(size, size).Draw(func(p *ui.Painter, r ui.Rect) {
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
			p.FillPath(internal.Diamond(internal.Snap(p, x), internal.Snap(p, y), internal.Snap(p, r.W*0.3)), color.Alpha(alpha))
		}
	})
}
