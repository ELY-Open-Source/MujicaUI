package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/internal"
)

// AspectRatioOptions configure an AspectRatio.
type AspectRatioOptions struct {
	// Ratio is width over height; zero is 1.
	Ratio float32
}

// AspectRatio centers the largest Ratio frame that fits; the child is centered, not stretched.
func AspectRatio(c *ui.Context, opts AspectRatioOptions, child func()) *ui.Element {
	if opts.Ratio < 0 {
		panic("mujica: AspectRatio Ratio must be positive")
	}
	if child == nil {
		panic("mujica: AspectRatio needs a child")
	}
	r := opts.Ratio
	if r == 0 {
		r = 1
	}
	outer := ui.Box(c).Fill().Center()
	room := outer.Bounds()
	internal.Remeasure(c, outer, room.W, room.H)
	return outer.Children(func() {
		frame := ui.Box(c).AspectRatio(r).Shrink(0).Center().Clip()
		was := frame.Bounds()
		w := room.W
		// A room whose height differs from the frame's, or a frame narrower than it, has a set height.
		if room.H-was.H > 0.5 || was.H-room.H > 0.5 || was.W < room.W-0.5 {
			w = min(w, room.H*r)
		}
		if w <= 0 {
			// Not laid out yet: Remeasure rebuilds once the room is known.
			frame.WidthPercent(100)
		} else {
			frame.Width(w)
			if w-was.W > 0.5 || was.W-w > 0.5 {
				c.Invalidate()
			}
		}
		frame.Children(child)
	})
}
