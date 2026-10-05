package mujica

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// BoxGroup is one named sample of a BoxPlot.
type BoxGroup struct {
	Name   string
	Values []float64
}

// BoxPlotOptions configure a BoxPlot.
type BoxPlotOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// boxStats are quartiles, whisker ends and outliers.
type boxStats struct {
	q1, median, q3, lo, hi float64
	outliers               []float64
	n                      int
}

// boxSummary summarizes finite values; n is 0 without any.
func boxSummary(values []float64) boxStats {
	var v []float64
	for _, x := range values {
		if !math.IsNaN(x) && !math.IsInf(x, 0) {
			v = append(v, x)
		}
	}
	if len(v) == 0 {
		return boxStats{}
	}
	slices.Sort(v)
	q := func(p float64) float64 {
		at := p * float64(len(v)-1)
		i := int(at)
		if i+1 >= len(v) {
			return v[i]
		}
		return v[i] + (at-float64(i))*(v[i+1]-v[i])
	}
	s := boxStats{q1: q(0.25), median: q(0.5), q3: q(0.75), n: len(v)}
	fence := 1.5 * (s.q3 - s.q1)
	s.lo, s.hi = s.q1, s.q3
	for _, x := range v {
		if x < s.q1-fence || x > s.q3+fence {
			s.outliers = append(s.outliers, x)
			continue
		}
		s.lo, s.hi = min(s.lo, x), max(s.hi, x)
	}
	return s
}

// BoxPlot compares groups by quartile box, whiskers and outlier points.
func BoxPlot(c *ui.Context, groups []BoxGroup, opts BoxPlotOptions) ChartResult {
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	pal := internal.Palette(k)
	size := FontSize(c, theme.CaptionSize)
	ids := ""
	for _, g := range groups {
		if g.Name == "" {
			panic("mujica: BoxPlot groups need names")
		}
		ids += chartID(g.Values) + ";"
	}
	stats := chartMemo(ch, "stats", ids, func() []boxStats {
		out := make([]boxStats, len(groups))
		for i, g := range groups {
			out[i] = boxSummary(g.Values)
		}
		return out
	})
	lo, hi, n := math.Inf(1), math.Inf(-1), 0
	for _, s := range stats {
		if s.n > 0 {
			n++
			lo, hi = min(lo, s.lo, slices.Min(append([]float64{s.lo}, s.outliers...))), max(hi, s.hi, slices.Max(append([]float64{s.hi}, s.outliers...)))
		}
	}
	var step float64
	if n > 0 {
		lo, hi, step = internal.Nice(lo, hi, 4)
	} else {
		groups = nil
	}
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear) {
		plot := chartInset(r, 52, 26)
		return plot, internal.Band{N: len(groups), R0: plot.X, R1: plot.X + plot.W, Pad: 0.45},
			internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: len(groups)}
	spec.Readout = func(i int) string {
		s := stats[i]
		if s.n == 0 {
			return groups[i].Name + " · " + msg(c, "chart.missing")
		}
		return fmt.Sprintf(msg(c, "box_plot.readout"), groups[i].Name, val(s.median), val(s.q1), val(s.q3), val(s.lo), val(s.hi), len(s.outliers))
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, band, _ := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return min(len(groups)-1, int((x-plot.X)/band.Step()))
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band, ys := geom(r)
		internal.AxisY(p, plot, ys, internal.Ticks(lo, hi, step), func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		var xs []float32
		var names []string
		for i, g := range groups {
			xs, names = append(xs, band.Center(i)), append(names, g.Name)
		}
		internal.AxisX(p, plot, xs, names, size, k)
		for i, s := range stats {
			if s.n == 0 {
				continue
			}
			col := pal[i%len(pal)]
			x, w := band.Start(i), band.Width()
			cx := internal.Snap(p, x+w/2)
			whisker := func(a, b float64) {
				var path ui.Path
				path.MoveTo(cx, ys.Map(a)).LineTo(cx, ys.Map(b))
				p.StrokePath(&path, 1.5, k.TextMuted)
				p.Line(cx-w/4, ys.Map(b), cx+w/4, ys.Map(b), 1.5, k.TextMuted)
			}
			whisker(s.q3, s.hi)
			whisker(s.q1, s.lo)
			box := ui.Rect{X: x, Y: ys.Map(s.q3), W: w, H: max(ys.Map(s.q1)-ys.Map(s.q3), 1)}
			p.Fill(box, col.Alpha(0.22), 2)
			p.Stroke(box, col, 2, 1.5)
			p.Line(x, ys.Map(s.median), x+w, ys.Map(s.median), 2.5, k.Text)
			for _, o := range s.outliers {
				var dot ui.Path
				dot.Circle(cx, internal.Snap(p, ys.Map(o)), 3.5)
				p.StrokePath(&dot, 1.5, col)
			}
			if i == active {
				p.Stroke(ui.Rect{X: x - 6, Y: plot.Y, W: w + 12, H: plot.H}, k.Text, 3, 1.5)
			}
		}
	}
	res, _ := ch.build(spec)
	return res
}
