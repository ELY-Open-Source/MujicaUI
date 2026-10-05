package mujica

import "github.com/egoist/mygo/ui"

// AspectRatioOptions configure an AspectRatio.
type AspectRatioOptions struct {
	// Ratio is width over height; zero is 1.
	Ratio float32
}

// AspectRatio fills its parent's box and centers a frame of Ratio inside, as large as fits.
// The child is centered in the frame without stretching and clipped to it.
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
	return outer.Children(func() {
		frame := ui.Box(c).AspectRatio(r).Shrink(0).Center().Clip()
		was := frame.Bounds()
		w := room.W
		// A room whose height differs from the frame's, or a frame narrower than it, has a set height.
		if room.H-was.H > 0.5 || was.H-room.H > 0.5 || was.W < room.W-0.5 {
			w = min(w, room.H*r)
		}
		if w <= 0 {
			// Not laid out yet: lay out once, then size from the room.
			frame.WidthPercent(100)
			if first := ui.Local(outer, "aspect_ratio.first", func() bool { return true }); *first {
				*first = false
				c.Invalidate()
			}
		} else {
			frame.Width(w)
			if w-was.W > 0.5 || was.W-w > 0.5 {
				c.Invalidate()
			}
		}
		frame.Children(child)
	})
}
