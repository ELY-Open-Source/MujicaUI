package feedback

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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

// Skeleton stands in for loading content; it pulses while painted, static under reduced motion.
func Skeleton(c *ui.Context, opts SkeletonOptions) ui.Element {
	if opts.Shape != SkeletonText && opts.Shape != SkeletonCircle && opts.Shape != SkeletonRect {
		panic(fmt.Sprintf("mujica: unknown SkeletonShape %d", int(opts.Shape)))
	}
	if !(opts.Width >= 0) || !(opts.Height >= 0) || opts.Lines < 0 || math.IsInf(float64(opts.Width+opts.Height), 1) {
		panic(fmt.Sprintf("mujica: Skeleton size %vx%v, %d lines invalid", opts.Width, opts.Height, opts.Lines))
	}
	k := core.Tokens(c)
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
			h = core.Density(c).ControlHeight()
		}
		e.Height(h).Draw(func(p *ui.Painter, r ui.Rect) { paint(p, r, theme.CardRadius) })
	default:
		lines := max(opts.Lines, 1)
		size := core.FontSize(c, theme.BodySize)
		line := size * theme.LineHeight
		e.Height(line * float32(lines)).Draw(func(p *ui.Painter, r ui.Rect) {
			bar := size * 0.8
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
