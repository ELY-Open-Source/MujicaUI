package feedback

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ShimmerOptions configure a Shimmer placeholder.
type ShimmerOptions struct {
	// Width fills the parent when zero.
	Width float32
	// Height is ControlHeight when zero.
	Height float32
	// Label is read out while loading; empty means "Loading".
	Label string
}

// Shimmer sweeps a placeholder while loading, then builds content.
func Shimmer(c *ui.Context, loading bool, opts ShimmerOptions, content func()) ui.Element {
	if !(opts.Width >= 0) || !(opts.Height >= 0) || math.IsInf(float64(opts.Width+opts.Height), 1) {
		panic("mujica: Shimmer size must be finite and not negative")
	}
	if !loading {
		box := ui.Column(c).AlignItems(ui.Stretch)
		if content != nil {
			box.Children(content)
		}
		return box
	}
	k := core.Tokens(c)
	h := opts.Height
	if h == 0 {
		h = core.Density(c).ControlHeight()
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "shimmer.loading")
	}
	fill := k.SurfaceHover
	if k.Dark {
		fill = k.Border
	}
	still := c.Preferences().ReduceMotion
	e := ui.Column(c).Height(h).AlignItems(ui.Stretch).Role(ui.RoleProgress).Label(label)
	if opts.Width > 0 {
		e.Width(opts.Width)
	} else {
		e.FillWidth()
	}
	// Keyed by motion so a flip repaints.
	e.Children(func() {
		ui.Box(c.Key(still)).Grow(1).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, fill, theme.CardRadius)
			if still {
				return
			}
			p.AnimationFrame()
			t := float32(p.Now().UnixMilli()%1600) / 1600
			band := max(r.W*0.25, 24)
			x := internal.Snap(p, r.X-band+(r.W+band)*t)
			p.Clip(r, theme.CardRadius, func() {
				p.Fill(ui.Rect{X: x, Y: r.Y, W: band, H: r.H}, k.Surface.Alpha(0.55), 0)
			})
		})
	})
	return e
}
