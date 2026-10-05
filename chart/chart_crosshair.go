package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChartCrosshairOptions configure a ChartCrosshair.
type ChartCrosshairOptions struct {
	// X draws vertical lines; Y draws horizontal lines.
	X, Y *ChartAxisPart
	// Reach defaults to 24 DIPs for point snapping.
	Reach float32
}

// ChartCrosshairPart draws and snaps the crosshair.
type ChartCrosshairPart struct {
	opts ChartCrosshairOptions
	k    theme.Tokens
	size float32
	pad  float32
}

// ChartCrosshair draws crossed lines and axis readouts.
func ChartCrosshair(c *ui.Context, opts ChartCrosshairOptions) ChartCrosshairPart {
	if opts.X == nil && opts.Y == nil || opts.Reach < 0 || math.IsNaN(float64(opts.Reach)) || math.IsInf(float64(opts.Reach), 0) {
		panic("mujica: ChartCrosshair needs an axis and Reach ≥ 0")
	}
	if opts.Reach == 0 {
		opts.Reach = 24
	}
	return ChartCrosshairPart{opts: opts, k: core.Tokens(c), size: core.FontSize(c, theme.CaptionSize), pad: core.Density(c).Unit()}
}

// Snap returns the nearest point within Reach, or -1.
func (h ChartCrosshairPart) Snap(n int, at func(i int) (x, y float32), x, y float32) int {
	return internal.Nearest(n, at, x, y, h.opts.Reach, false)
}

// Draw paints crossed lines and value labels.
func (h ChartCrosshairPart) Draw(p *ui.Painter, plot ui.Rect, x, y float32) {
	k := h.k
	if a := h.opts.X; a != nil {
		xx := internal.Snap(p, x)
		p.Line(xx, plot.Y, xx, plot.Y+plot.H, 1, k.TextMuted)
		at := plot.Y + plot.H
		if a.opts.Side == SideTop {
			at = plot.Y
		}
		h.pill(p, a.readout(a.Invert(plot, x)), x, at, true, a.opts.Side == SideTop)
	}
	if a := h.opts.Y; a != nil {
		yy := internal.Snap(p, y)
		p.Line(plot.X, yy, plot.X+plot.W, yy, 1, k.TextMuted)
		at := plot.X
		if a.opts.Side == SideRight {
			at = plot.X + plot.W
		}
		h.pill(p, a.readout(a.Invert(plot, y)), at, y, false, a.opts.Side == SideRight)
	}
}

// pill paints a value tag outside the plot.
func (h ChartCrosshairPart) pill(p *ui.Painter, s string, x, y float32, across, far bool) {
	w, ht := p.MeasureText(0, ui.Span{Text: s, Size: h.size})
	box := ui.Rect{W: w + 2*h.pad, H: ht + h.pad}
	switch {
	case across && far:
		box.X, box.Y = x-box.W/2, y-box.H-2
	case across:
		box.X, box.Y = x-box.W/2, y+2
	case far:
		box.X, box.Y = x+2, y-box.H/2
	default:
		box.X, box.Y = x-box.W-2, y-box.H/2
	}
	p.Fill(box, h.k.Text, theme.ControlRadius)
	p.Text(box.X+h.pad, box.Y+h.pad/2, s, h.size, h.k.Surface)
}
