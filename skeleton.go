package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// SkeletonShape is the placeholder's form. The zero value is SkeletonText.
type SkeletonShape int

const (
	SkeletonText SkeletonShape = iota
	SkeletonCircle
	SkeletonRect
)

// SkeletonOptions configure a Skeleton. Sizes are fixed so content does not shift.
type SkeletonOptions struct {
	Shape SkeletonShape
	// Width fills the parent when zero; a circle uses it as its diameter (40 DIP when zero).
	Width float32
	// Height is the rect height (ControlHeight when zero); text derives it from Lines.
	Height float32
	// Lines is the text line count, 1 when zero; the last of several is shorter.
	Lines int
}

// Skeleton stands in for content that is loading. It pulses while painted
// and shows a static fill under reduced motion.
func Skeleton(c *ui.Context, opts SkeletonOptions) *ui.Element {
	if opts.Width < 0 || opts.Height < 0 || opts.Lines < 0 {
		panic("mujica: negative Skeleton size")
	}
	k := Tokens(c)
	still := c.Preferences().ReduceMotion
	fill := k.SurfaceHover
	if k.Dark {
		fill = k.Border
	}
	paint := func(p *ui.Painter, r ui.Rect, radius float32) {
		a := float32(1)
		if !still {
			p.AnimationFrame()
			t := float64(p.Now().UnixMilli()%1600) / 1600
			a = 0.55 + 0.45*float32(0.5+0.5*math.Cos(2*math.Pi*t))
		}
		p.Fill(r, fill.Alpha(a), radius)
	}
	e := ui.Box(c).Role(ui.RoleNone)
	if opts.Width > 0 {
		e.Width(opts.Width)
	} else if opts.Shape != SkeletonCircle {
		e.FillWidth()
	}
	switch opts.Shape {
	case SkeletonCircle:
		d := opts.Width
		if d == 0 {
			d = 40
		}
		e.Size(d, d).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) { paint(p, r, r.W/2) })
	case SkeletonRect:
		h := opts.Height
		if h == 0 {
			h = Density(c).ControlHeight()
		}
		e.Height(h).Draw(func(p *ui.Painter, r ui.Rect) { paint(p, r, theme.CardRadius) })
	default:
		lines := max(opts.Lines, 1)
		line := theme.BodySize * theme.LineHeight
		e.Height(line * float32(lines)).Draw(func(p *ui.Painter, r ui.Rect) {
			bar := theme.BodySize * 0.8
			for i := range lines {
				w := r.W
				if lines > 1 && i == lines-1 {
					w *= 0.6
				}
				y := r.Y + float32(i)*line + (line-bar)/2
				paint(p, ui.Rect{X: r.X, Y: y, W: w, H: bar}, theme.ControlRadius)
			}
		})
	}
	return e
}
