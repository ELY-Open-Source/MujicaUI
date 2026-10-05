package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// WaterfallKind says how a step's bar is drawn.
type WaterfallKind int

const (
	// WaterfallDelta adds Value to the running total.
	WaterfallDelta WaterfallKind = iota
	// WaterfallSubtotal shows the running total; its Value must be zero.
	WaterfallSubtotal
	// WaterfallTotal shows the final total; Value must be zero.
	WaterfallTotal
)

// WaterfallStep is one bar of a WaterfallChart.
type WaterfallStep struct {
	Label string
	Value float64
	Kind  WaterfallKind
}

// WaterfallChartOptions configure a WaterfallChart.
type WaterfallChartOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// WaterfallChart builds a total from rises and falls.
func WaterfallChart(c *ui.Context, steps []WaterfallStep, opts WaterfallChartOptions) ChartResult {
	from, to := make([]float64, len(steps)), make([]float64, len(steps))
	run, lo, hi := 0.0, 0.0, 0.0
	for i, s := range steps {
		if s.Label == "" || math.IsNaN(s.Value) || math.IsInf(s.Value, 0) || s.Kind != WaterfallDelta && s.Value != 0 || s.Kind > WaterfallTotal {
			panic("mujica: WaterfallChart steps need a label, a finite value, and zero for totals")
		}
		if s.Kind == WaterfallDelta {
			from[i], to[i], run = run, run+s.Value, run+s.Value
		} else {
			from[i], to[i] = 0, run
		}
		lo, hi = min(lo, from[i], to[i]), max(hi, from[i], to[i])
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	size := FontSize(c, theme.CaptionSize)
	var step float64
	if len(steps) > 0 {
		pad := (hi - lo) * 0.12
		lo, hi, step = internal.Nice(lo-pad*chartBool(lo < 0), hi+pad, 4)
	}
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	signed := func(v float64) string {
		if v >= 0 {
			return "+" + val(v)
		}
		return "−" + val(-v)
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear) {
		plot := chartInset(r, 52, 26)
		return plot, internal.Band{N: len(steps), R0: plot.X, R1: plot.X + plot.W, Pad: 0.3},
			internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: len(steps)}
	spec.Readout = func(i int) string {
		if steps[i].Kind != WaterfallDelta {
			return steps[i].Label + " · " + val(to[i])
		}
		return steps[i].Label + " · " + signed(steps[i].Value) + " · " + msg(c, "waterfall_chart.running") + " " + val(to[i])
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, band, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return min(len(steps)-1, int((x-plot.X)/band.Step()))
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band, ys := geom(r)
		internal.AxisY(p, plot, ys, internal.Ticks(lo, hi, step), func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		var xs []float32
		var names []string
		for i, s := range steps {
			xs, names = append(xs, band.Center(i)), append(names, s.Label)
		}
		internal.AxisX(p, plot, xs, names, size, k)
		internal.Rule(p, plot.X, plot.X+plot.W, ys.Map(0), k.ControlBorder)
		for i, s := range steps {
			a, b := ys.Map(from[i]), ys.Map(to[i])
			bar := ui.Rect{X: band.Start(i), Y: min(a, b), W: band.Width(), H: max(abs32(a-b), 1)}
			col, shape, text := internal.Palette(k)[5], 1, val(to[i])
			switch {
			case s.Kind != WaterfallDelta:
			case s.Value >= 0:
				col, shape, text = k.Success, 3, signed(s.Value)
			default:
				col, shape, text = k.Danger, 4, signed(s.Value)
			}
			p.Fill(bar, col, 1)
			internal.Marker(p, shape, bar.X+bar.W/2, bar.Y-14+(bar.H+28)*float32(chartBool(s.Kind == WaterfallDelta && s.Value < 0)), 8, col)
			ty, ay := bar.Y-22, float32(1)
			if s.Kind == WaterfallDelta && s.Value < 0 {
				ty, ay = bar.Y+bar.H+22, 0
			}
			chartText(p, text, bar.X+bar.W/2, ty, size, k.Text, 0.5, ay)
			if i+1 < len(steps) {
				y := internal.Snap(p, ys.Map(to[i]))
				p.Line(bar.X+bar.W, y, band.Start(i+1), y, 1, k.TextMuted)
			}
			if i == active {
				p.Stroke(ui.Rect{X: bar.X - 2, Y: bar.Y - 2, W: bar.W + 4, H: bar.H + 4}, k.Text, 2, 1.5)
			}
		}
	}
	res, _ := ch.build(spec)
	return res
}
