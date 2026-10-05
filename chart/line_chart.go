package chart

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// LineChartOptions configure a LineChart.
type LineChartOptions struct {
	Frame ChartFrame
	// X places point i at ascending X[i]; nil means i.
	X []float64
	// Times makes a time axis; set X or Times.
	Times []time.Time
	// Format writes readout values; nil uses six digits.
	Format func(v float64) string
}

// LineChart draws series as lines; NaN breaks them.
func LineChart(c *ui.Context, series []ChartSeries, opts LineChartOptions) ChartResult {
	n := chartCheckSeries(series)
	ch := NewChart(c, opts.Frame)
	ax := ChartMemo(ch, "x", [3]string{ChartID(opts.X), ChartID(opts.Times), ChartID(series)}, func() ChartX { return NewChartX(n, opts.X, opts.Times) })
	lo, hi, ok := ChartYRange(ch, series, false, nil)
	if !ok {
		n = 0
	}
	tops := make([][]float64, len(series))
	for i, s := range series {
		tops[i] = s.Values
	}
	return chartLines(ch, chartLineData{series: series, tops: tops, ax: ax, lo: lo, hi: hi, n: n, format: opts.Format})
}

// ChartYRange is the memoized y domain of visible series.
func ChartYRange(ch *chart, series []ChartSeries, zero bool, tops [][]float64) (lo, hi float64, ok bool) {
	d := ChartMemo(ch, "y", chartSeriesID(ch, series)+ChartID(tops), func() [3]float64 {
		src := series
		if tops != nil {
			src = make([]ChartSeries, len(series))
			for i, s := range series {
				src[i] = ChartSeries{Name: s.Name, Values: tops[i]}
			}
		}
		lo, hi, ok := chartYDomain(src, ch.Hidden, zero)
		if !ok {
			return [3]float64{}
		}
		return [3]float64{lo, hi, 1}
	})
	return d[0], d[1], d[2] == 1
}

// chartLineData is what line and area charts draw.
type chartLineData struct {
	series      []ChartSeries
	tops, bases [][]float64
	ax          ChartX
	lo, hi      float64
	n           int
	format      func(float64) string
}

func chartLines(ch *chart, d chartLineData) ChartResult {
	c, k := ch.c, core.Tokens(ch.c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	keys := make([]ChartLegendItem, len(d.series))
	for i, s := range d.series {
		keys[i] = ChartLegendItem{Name: s.Name, Color: pal[i%len(pal)], Shape: i}
	}
	var ticks []float64
	var step float64
	lo, hi := d.lo, d.hi
	if d.n > 0 {
		lo, hi, step = internal.Nice(lo, hi, 4)
		ticks = internal.Ticks(lo, hi, step)
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := ChartInset(r, 52, 26)
		return plot, d.ax.Scale(plot), internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: d.n, Legend: keys, Hidden: ch.Toggles(true)}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, xs, _ := geom(r)
		if !(ui.Rect{X: plot.X - 6, Y: plot.Y, W: plot.W + 12, H: plot.H}).Contains(x, y) {
			return -1
		}
		return d.ax.Nearest(xs.Invert(x))
	}
	spec.Readout = func(i int) string {
		parts := []string{d.ax.label(i)}
		for _, s := range d.series {
			if !ch.Hidden(s.Name) {
				parts = append(parts, s.Name+" "+ChartValue(c, s.Values[i], d.format))
			}
		}
		return strings.Join(parts, " · ")
	}
	spec.Tooltip = func(i int) []ChartTooltipLine {
		lines := []ChartTooltipLine{{Text: d.ax.label(i)}}
		for si, s := range d.series {
			if !ch.Hidden(s.Name) {
				lines = append(lines, ChartTooltipLine{Text: s.Name + " " + ChartValue(c, s.Values[i], d.format), Color: keys[si].Color, Shape: si})
			}
		}
		return lines
	}
	spec.Anchor = func(r ui.Rect, i int) (float32, float32) {
		plot, xs, _ := geom(r)
		return xs.Map(d.ax.Vals[i]), plot.Y
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		tx, tl := d.ax.Ticks(plot)
		internal.AxisX(p, plot, tx, tl, size, k)
		x := func(i int) float32 { return xs.Map(d.ax.Vals[i]) }
		pt := func(i int, v float64) (float32, float32) { return internal.Snap(p, x(i)), internal.Snap(p, ys.Map(v)) }
		for si, s := range d.series {
			if ch.Hidden(s.Name) {
				continue
			}
			col, top := keys[si].Color, d.tops[si]
			chartRuns(top, x, func(idx []int) {
				var line ui.Path
				for j, i := range idx {
					if j == 0 {
						line.MoveTo(pt(i, top[i]))
					} else {
						line.LineTo(pt(i, top[i]))
					}
				}
				if d.bases != nil {
					fill := line
					for j := len(idx) - 1; j >= 0; j-- {
						fill.LineTo(pt(idx[j], max(lo, min(hi, d.bases[si][idx[j]]))))
					}
					fill.Close()
					p.FillPath(&fill, col.Alpha(0.22))
				}
				p.StrokePath(&line, 2, col)
			})
			// Sparse markers name each series by shape, not only color.
			every := max(1, int(float32(d.n)/max(plot.W/48, 1)))
			for i := 0; i < d.n; i += every {
				if !math.IsNaN(top[i]) {
					internal.Marker(p, si, x(i), ys.Map(top[i]), 7, col)
				}
			}
		}
		if active < 0 {
			return
		}
		ax0 := internal.Snap(p, x(active))
		p.Line(ax0, plot.Y, ax0, plot.Y+plot.H, 1, k.TextMuted)
		for si, s := range d.series {
			if v := d.tops[si][active]; !ch.Hidden(s.Name) && !math.IsNaN(v) {
				internal.Marker(p, si, ax0, ys.Map(v), 13, k.Surface)
				internal.Marker(p, si, ax0, ys.Map(v), 9, keys[si].Color)
			}
		}
	}
	res, _ := ch.Build(spec)
	return res
}
