package mujica

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// PerformanceChartOptions configure a PerformanceChart.
type PerformanceChartOptions struct {
	Frame   ChartFrame
	Version int
	// Times dates each value; required, ascending.
	Times []time.Time
	// Benchmark is compared against the portfolio; nil shows none.
	Benchmark           []float64
	Name, BenchmarkName string
}

type Drawdown struct {
	Peak, Trough int
	Percent      float64
}

func MaxDrawdown(values []float64) Drawdown {
	var dd Drawdown
	peak := 0
	for i, v := range values {
		if v > values[peak] {
			peak = i
		}
		if f := (v/values[peak] - 1) * 100; f < dd.Percent {
			dd = Drawdown{Peak: peak, Trough: i, Percent: f}
		}
	}
	return dd
}

func PerformanceChart(c *ui.Context, values []float64, opts PerformanceChartOptions) ChartResult {
	if opts.Benchmark != nil && len(opts.Benchmark) != len(values) {
		panic("mujica: PerformanceChart Benchmark length differs")
	}
	if len(opts.Times) != len(values) {
		panic("mujica: PerformanceChart Times length differs")
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	pal := internal.Palette(k)
	size := FontSize(c, theme.CaptionSize)
	name := opts.Name
	if name == "" {
		name = msg(c, "performance_chart.portfolio")
	}
	bname := opts.BenchmarkName
	if bname == "" {
		bname = msg(c, "performance_chart.benchmark")
	}
	if opts.Benchmark != nil && name == bname {
		panic("mujica: PerformanceChart series names must differ")
	}
	series := chartMemo(ch, "returns", [4]string{chartID(values), chartID(opts.Benchmark), name + fmt.Sprint(opts.Version), bname}, func() []ChartSeries {
		for _, vs := range [][]float64{values, opts.Benchmark} {
			for _, v := range vs {
				if !(v > 0) || math.IsInf(v, 0) {
					panic("mujica: PerformanceChart values must be positive")
				}
			}
		}

		out := []ChartSeries{{Name: name, Values: performanceReturns(values)}}
		if opts.Benchmark != nil {
			out = append(out, ChartSeries{Name: bname, Values: performanceReturns(opts.Benchmark)})
		}
		return out
	})
	n := len(values)
	ax := chartMemo(ch, "x", chartID(opts.Times)+fmt.Sprint(opts.Version), func() chartX { return newChartX(n, nil, opts.Times) })
	dd := chartMemo(ch, "dd", chartID(values)+fmt.Sprint(opts.Version), func() Drawdown { return MaxDrawdown(values) })
	lo, hi, ok := chartYRange(ch, series, true, nil)
	if !ok {
		n = 0
	}
	keys := make([]ChartLegendItem, len(series))
	for i, s := range series {
		keys[i] = ChartLegendItem{Name: s.Name, Color: pal[i], Shape: i}
	}
	var ticks []float64
	var step float64
	if n > 0 {
		lo, hi, step = internal.Nice(lo, hi, 4)
		ticks = internal.Ticks(lo, hi, step)
	}
	pct := func(v float64) string { return statisticSigned(v, 2) + "%" }
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := chartInset(r, 56, 26)
		return plot, ax.scale(plot), internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: n, Legend: keys, Hidden: ch.toggles(true)}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, xs, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return ax.nearest(xs.Invert(x))
	}
	spec.Readout = func(i int) string {
		parts := []string{opts.Times[i].Format("2006-01-02")}
		for _, s := range series {
			if !ch.hidden(s.Name) {
				parts = append(parts, s.Name+" "+pct(s.Values[i]))
			}
		}
		return strings.Join(parts, " · ")
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		x := func(i int) float32 { return internal.Snap(p, xs.Map(ax.vals[i])) }
		if dd.Percent < 0 && !ch.hidden(name) {
			band := ui.Rect{X: x(dd.Peak), Y: plot.Y, W: x(dd.Trough) - x(dd.Peak), H: plot.H}
			p.Fill(band, trendColor(c, -1).Alpha(0.1), 0)
			p.Line(band.X, plot.Y, band.X, plot.Y+plot.H, 1, trendColor(c, -1).Alpha(0.5))
			p.Line(band.X+band.W, plot.Y, band.X+band.W, plot.Y+plot.H, 1, trendColor(c, -1).Alpha(0.5))
			chartText(p, fmt.Sprintf(msg(c, "performance_chart.drawdown"), pct(dd.Percent)), band.X+4, plot.Y+2, size, chartInk(k, k.Surface), 0, 0)
		}
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) + "%" }, size, k)
		tx, tl := ax.ticks(plot)
		internal.AxisX(p, plot, tx, tl, size, k)
		zero := internal.Snap(p, ys.Map(0))
		p.Line(plot.X, zero, plot.X+plot.W, zero, 1, k.ControlBorder)
		for si, s := range series {
			if ch.hidden(s.Name) {
				continue
			}
			var line ui.Path
			idx := chartMemo(ch, "performance.geometry."+s.Name, [3]any{chartID(s.Values), chartID(ax.vals), plot.W}, func() []int {
				return internal.Columns(n, func(i int) float32 { return xs.Map(ax.vals[i]) - plot.X }, func(i int) float64 { return s.Values[i] })
			})
			for j, i := range idx {
				if j == 0 {
					line.MoveTo(x(i), internal.Snap(p, ys.Map(s.Values[i])))
				} else {
					line.LineTo(x(i), internal.Snap(p, ys.Map(s.Values[i])))
				}
			}
			p.StrokePath(&line, 2, keys[si].Color)
			every := max(1, int(float32(n)/max(plot.W/48, 1)))
			for i := 0; i < n; i += every {
				internal.Marker(p, si, x(i), ys.Map(s.Values[i]), 7, keys[si].Color)
			}
		}
		if active >= 0 {
			p.Line(x(active), plot.Y, x(active), plot.Y+plot.H, 1, k.TextMuted)
			for si, s := range series {
				if !ch.hidden(s.Name) {
					internal.Marker(p, si, x(active), ys.Map(s.Values[active]), 13, k.Surface)
					internal.Marker(p, si, x(active), ys.Map(s.Values[active]), 9, keys[si].Color)
				}
			}
		}
	}
	res, _ := ch.build(spec)
	if dd.Percent < 0 && n > 0 {
		res.Element.Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, "performance_chart.drawdown_range"), pct(dd.Percent),
				opts.Times[dd.Peak].Format("2006-01-02"), opts.Times[dd.Trough].Format("2006-01-02"))).FontSize(size).TextColor(k.TextMuted)
		})
	}
	return res
}

// performanceReturns turns values into cumulative percent return from the first.
func performanceReturns(values []float64) []float64 {
	out := make([]float64, len(values))
	for i, v := range values {
		out[i] = (v/values[0] - 1) * 100
	}
	return out
}
