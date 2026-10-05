package mujica

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// BulletChartOptions configure a BulletChart; Frame.Height 0 means 72 DIPs.
type BulletChartOptions struct {
	Frame ChartFrame
	// Target, when set, draws a cross line at the goal.
	Target *float64
	// Bands are the ascending upper ends of the qualitative ranges.
	Bands []float64
	// Max ends the scale; 0 means automatic.
	Max      float64
	Vertical bool
	Format   func(v float64) string
}

// BulletChart shows a value against a target and ranges.
func BulletChart(c *ui.Context, value float64, opts BulletChartOptions) ChartResult {
	if math.IsNaN(value) || value < 0 || opts.Max < 0 || !slices.IsSorted(opts.Bands) || len(opts.Bands) > 0 && opts.Bands[0] <= 0 {
		panic("mujica: BulletChart needs a value ≥ 0, Max ≥ 0 and positive ascending Bands")
	}
	if opts.Frame.Height == 0 {
		opts.Frame.Height = 72
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	size := FontSize(c, theme.CaptionSize)
	top := opts.Max
	if top == 0 {
		top = value
		if opts.Target != nil {
			top = max(top, *opts.Target)
		}
		if len(opts.Bands) > 0 {
			top = max(top, opts.Bands[len(opts.Bands)-1])
		}
		_, top, _ = internal.Nice(0, max(top, 1e-9), 4)
	}
	_, _, step := internal.Nice(0, top, 4)
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	marks := []string{msg(c, "bullet_chart.actual") + " " + val(value)}
	if opts.Target != nil {
		marks = append(marks, msg(c, "bullet_chart.target")+" "+val(*opts.Target))
	}
	for i, b := range opts.Bands {
		lo := 0.0
		if i > 0 {
			lo = opts.Bands[i-1]
		}
		marks = append(marks, fmt.Sprintf(msg(c, "bullet_chart.band"), i+1, val(lo), val(b)))
	}
	// geom returns the plot and its value scale.
	geom := func(r ui.Rect) (ui.Rect, internal.Linear) {
		if opts.Vertical {
			plot := ui.Rect{X: r.X + r.W/2 - 14, Y: r.Y + 8, W: 28, H: r.H - 16}
			return plot, internal.Linear{D0: 0, D1: top, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
		}
		plot := ui.Rect{X: r.X + 8, Y: r.Y + 6, W: r.W - 24, H: max(r.H-30, 12)}
		return plot, internal.Linear{D0: 0, D1: top, R0: float64(plot.X), R1: float64(plot.X + plot.W)}
	}
	span := func(plot ui.Rect, s internal.Linear, a, b, inset float32) ui.Rect {
		if opts.Vertical {
			return ui.Rect{X: plot.X + inset, Y: s.Map(float64(b)), W: plot.W - 2*inset, H: s.Map(float64(a)) - s.Map(float64(b))}
		}
		return ui.Rect{X: s.Map(float64(a)), Y: plot.Y + inset, W: s.Map(float64(b)) - s.Map(float64(a)), H: plot.H - 2*inset}
	}
	spec := chartSpec{n: len(marks)}
	spec.readout = func(i int) string { return marks[i] }
	spec.pick = func(r ui.Rect, x, y float32) int {
		plot, s := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		at := x
		if opts.Vertical {
			at = y
		}
		v := s.Invert(at)
		if opts.Target != nil && math.Abs(float64(s.Map(*opts.Target)-at)) < 4 {
			return 1
		}
		if span(plot, s, 0, float32(min(value, top)), bulletInset(plot, opts.Vertical)).Contains(x, y) {
			return 0
		}
		for i, b := range opts.Bands {
			if v <= b {
				return len(marks) - len(opts.Bands) + i
			}
		}
		return -1
	}
	spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, s := geom(r)
		inset := bulletInset(plot, opts.Vertical)
		prev := 0.0
		for i, b := range opts.Bands {
			shade := k.SurfaceHover.Mix(k.Border, float32(len(opts.Bands)-1-i)/float32(max(len(opts.Bands)-1, 1)))
			p.Fill(span(plot, s, float32(prev), float32(min(b, top)), 0), shade, 0)
			if active == len(marks)-len(opts.Bands)+i {
				p.Stroke(span(plot, s, float32(prev), float32(min(b, top)), 0), k.Text, 0, 1.5)
			}
			prev = b
		}
		bar := span(plot, s, 0, float32(min(value, top)), inset)
		p.Fill(bar, internal.Palette(k)[0], 0)
		if active == 0 {
			p.Stroke(ui.Rect{X: bar.X - 2, Y: bar.Y - 2, W: bar.W + 4, H: bar.H + 4}, k.Text, 1, 1.5)
		}
		if opts.Target != nil {
			t := internal.Snap(p, s.Map(min(*opts.Target, top)))
			w := float32(2.5)
			if active == 1 {
				w = 4
			}
			if opts.Vertical {
				p.Line(plot.X+2, t, plot.X+plot.W-2, t, w, k.Text)
			} else {
				p.Line(t, plot.Y+2, t, plot.Y+plot.H-2, w, k.Text)
			}
		}
		for _, tk := range internal.Ticks(0, top, step) {
			if opts.Vertical {
				chartText(p, internal.FormatTick(tk, step), plot.X+plot.W+6, s.Map(tk), size, k.TextMuted, 0, 0.5)
			} else {
				chartText(p, internal.FormatTick(tk, step), s.Map(tk), plot.Y+plot.H+4, size, k.TextMuted, 0.5, 0)
			}
		}
	}
	res, _ := ch.build(spec)
	return res
}

// bulletInset is a third of the band's thickness.
func bulletInset(plot ui.Rect, vertical bool) float32 {
	if vertical {
		return plot.W / 3
	}
	return plot.H / 3
}
