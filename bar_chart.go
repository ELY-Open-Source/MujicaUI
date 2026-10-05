package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// BarChartOptions configure a BarChart.
type BarChartOptions struct {
	Frame ChartFrame
	// Horizontal lays categories down the side and bars across.
	Horizontal bool
	// Stacked piles series into one bar per category.
	Stacked bool
	// Labels writes each bar's value at its end.
	Labels bool
	Format func(v float64) string
}

type barMark struct{ cat, series int }

// BarChart draws grouped or stacked bars; NaN draws none.
func BarChart(c *ui.Context, categories []string, series []ChartSeries, opts BarChartOptions) ChartResult {
	if n := chartCheckSeries(series); n != len(categories) && len(series) > 0 {
		panic("mujica: BarChart series length differs from categories")
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	pal := internal.Palette(k)
	size := FontSize(c, theme.CaptionSize)
	keys := make([]ChartLegendItem, len(series))
	var vis []int
	for i, s := range series {
		keys[i] = ChartLegendItem{Name: s.Name, Color: pal[i%len(pal)], Shape: i}
		if !ch.hidden(s.Name) {
			vis = append(vis, i)
		}
	}
	// Bars stacked or not; lo and hi always include zero.
	type layout struct {
		marks  []barMark
		from   []float64
		lo, hi float64
	}
	lay := chartMemo(ch, "bars", [2]any{chartSeriesID(ch, series), opts.Stacked}, func() layout {
		var l layout
		for ci := range categories {
			pos, neg := 0.0, 0.0
			for _, si := range vis {
				v := series[si].Values[ci]
				if math.IsNaN(v) {
					continue
				}
				from := 0.0
				if opts.Stacked && v >= 0 {
					from, pos = pos, pos+v
				} else if opts.Stacked {
					from, neg = neg, neg+v
				}
				l.marks, l.from = append(l.marks, barMark{ci, si}), append(l.from, from)
				l.lo, l.hi = min(l.lo, from, from+v), max(l.hi, from, from+v)
			}
		}
		return l
	})
	var lo, hi, step float64
	var ticks []float64
	if len(lay.marks) > 0 {
		lo, hi = lay.lo, lay.hi
		if opts.Labels {
			// Headroom keeps value labels inside the plot.
			pad := (hi - lo) * 0.15
			lo, hi = lo-pad*chartBool(lo < 0), hi+pad*chartBool(hi > 0)
		}
		lo, hi, step = internal.Nice(lo, hi, 4)
		ticks = internal.Ticks(lo, hi, step)
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear) {
		if opts.Horizontal {
			plot := chartInset(r, 96, 26)
			return plot, internal.Band{N: len(categories), R0: plot.Y, R1: plot.Y + plot.H, Pad: 0.3},
				internal.Linear{D0: lo, D1: hi, R0: float64(plot.X), R1: float64(plot.X + plot.W)}
		}
		plot := chartInset(r, 52, 26)
		return plot, internal.Band{N: len(categories), R0: plot.X, R1: plot.X + plot.W, Pad: 0.3},
			internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	// rect is bar i's box in plot space.
	rect := func(band internal.Band, vs internal.Linear, i int) ui.Rect {
		m := lay.marks[i]
		at, w := band.Start(m.cat), band.Width()
		if !opts.Stacked && len(vis) > 0 {
			slot := 0
			for j, si := range vis {
				if si == m.series {
					slot = j
				}
			}
			w /= float32(len(vis))
			at += w * float32(slot)
		}
		a, b := vs.Map(lay.from[i]), vs.Map(lay.from[i]+series[m.series].Values[m.cat])
		a, b = min(a, b), max(a, b)
		if opts.Horizontal {
			return ui.Rect{X: a, Y: at, W: max(b-a, 1), H: w}
		}
		return ui.Rect{X: at, Y: a, W: w, H: max(b-a, 1)}
	}
	value := func(i int) float64 { m := lay.marks[i]; return series[m.series].Values[m.cat] }
	spec := chartSpec{n: len(lay.marks), keys: keys, toggle: len(series) > 1}
	spec.readout = func(i int) string {
		m := lay.marks[i]
		return categories[m.cat] + " · " + series[m.series].Name + " " + chartValue(c, value(i), opts.Format)
	}
	spec.pick = func(r ui.Rect, x, y float32) int {
		_, band, vs := geom(r)
		for i := range lay.marks {
			b := rect(band, vs, i)
			// Thin bars still take a 6 DIP target.
			if (ui.Rect{X: b.X - 3, Y: b.Y - 3, W: b.W + 6, H: b.H + 6}).Contains(x, y) {
				return i
			}
		}
		return -1
	}
	spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band, vs := geom(r)
		label := func(v float64) string { return internal.FormatTick(v, step) }
		var pos []float32
		for i := range categories {
			pos = append(pos, band.Center(i))
		}
		zero := vs.Map(0)
		if opts.Horizontal {
			for _, t := range ticks {
				x := internal.Snap(p, vs.Map(t))
				p.Line(x, plot.Y, x, plot.Y+plot.H, 1, k.Border)
				chartText(p, label(t), x, plot.Y+plot.H+6, size, k.TextMuted, 0.5, 0)
			}
			for i, y := range pos {
				chartText(p, categories[i], plot.X-8, y, size, k.TextMuted, 1, 0.5)
			}
			p.Line(internal.Snap(p, zero), plot.Y, internal.Snap(p, zero), plot.Y+plot.H, 1, k.ControlBorder)
		} else {
			internal.AxisY(p, plot, vs, ticks, label, size, k)
			internal.AxisX(p, plot, pos, categories, size, k)
			internal.Rule(p, plot.X, plot.X+plot.W, zero, k.ControlBorder)
		}
		for i, m := range lay.marks {
			b := rect(band, vs, i)
			col := keys[m.series].Color
			p.Fill(b, col, 1)
			if i == active {
				p.Stroke(ui.Rect{X: b.X - 2, Y: b.Y - 2, W: b.W + 4, H: b.H + 4}, k.Text, 2, 1.5)
			}
			// The series shape marks the bar besides its color.
			if min(b.W, b.H) >= 10 && max(b.W, b.H) >= 14 {
				internal.Marker(p, m.series, b.X+b.W/2, b.Y+b.H/2, min(b.W, b.H, 14)*0.6, k.Surface)
			}
			if opts.Labels && !opts.Stacked {
				barLabel(p, chartValue(c, value(i), opts.Format), b, value(i) < 0, opts.Horizontal, size, k.Text)
			}
		}
	}
	res, _ := ch.build(spec)
	return res
}

// barLabel writes a value just past the bar's far end.
func barLabel(p *ui.Painter, s string, b ui.Rect, neg, horizontal bool, size float32, col ui.Color) {
	switch {
	case horizontal && neg:
		chartText(p, s, b.X-4, b.Y+b.H/2, size, col, 1, 0.5)
	case horizontal:
		chartText(p, s, b.X+b.W+4, b.Y+b.H/2, size, col, 0, 0.5)
	case neg:
		chartText(p, s, b.X+b.W/2, b.Y+b.H+2, size, col, 0.5, 0)
	default:
		chartText(p, s, b.X+b.W/2, b.Y-2, size, col, 0.5, 1)
	}
}
