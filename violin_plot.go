package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ViolinPlotOptions configure a ViolinPlot.
type ViolinPlotOptions struct {
	Frame ChartFrame
	// Bandwidth is the kernel width; 0 uses Silverman's rule per group.
	Bandwidth float64
	Format    func(v float64) string
}

// violinShape is one group's density outline over its value range, and its quartiles.
type violinShape struct {
	ys, dens []float64
	stats    boxStats
}

// violinPoints is how many values each outline is sampled at.
const violinPoints = 64

// ViolinPlot mirrors each group's density around its middle, with the median and quartiles marked.
func ViolinPlot(c *ui.Context, groups []BoxGroup, opts ViolinPlotOptions) ChartResult {
	if opts.Bandwidth < 0 {
		panic("mujica: ViolinPlot Bandwidth must be ≥ 0")
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	pal := internal.Palette(k)
	size := FontSize(c, theme.CaptionSize)
	ids := ""
	for _, g := range groups {
		if g.Name == "" {
			panic("mujica: ViolinPlot groups need names")
		}
		ids += chartID(g.Values) + ";"
	}
	shapes := chartMemo(ch, "violins", fmt.Sprint(ids, opts.Bandwidth), func() []violinShape {
		out := make([]violinShape, len(groups))
		for i, g := range groups {
			v := densitySample(g.Values)
			if len(v) == 0 {
				continue
			}
			h := opts.Bandwidth
			if h == 0 {
				h = densityBandwidth(v)
			}
			s := boxSummary(v)
			lo, hi := s.lo, s.hi
			for _, o := range s.outliers {
				lo, hi = min(lo, o), max(hi, o)
			}
			ys := make([]float64, violinPoints)
			for j := range ys {
				ys[j] = lo + (hi-lo)*float64(j)/float64(violinPoints-1)
			}
			out[i] = violinShape{ys: ys, dens: densityKDE(v, h, ys), stats: s}
		}
		return out
	})
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, s := range shapes {
		if s.stats.n > 0 {
			lo, hi = min(lo, s.ys[0]), max(hi, s.ys[len(s.ys)-1])
		}
	}
	var step float64
	n := len(groups)
	if math.IsInf(lo, 1) {
		n = 0
	} else {
		lo, hi, step = internal.Nice(lo, hi, 4)
	}
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear) {
		plot := chartInset(r, 52, 26)
		return plot, internal.Band{N: len(groups), R0: plot.X, R1: plot.X + plot.W, Pad: 0.15},
			internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		s := shapes[i].stats
		if s.n == 0 {
			return groups[i].Name + " · " + msg(c, "chart.missing")
		}
		return fmt.Sprintf(msg(c, "violin_plot.readout"), groups[i].Name, val(s.median), val(s.q1), val(s.q3), s.n)
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
		for i, s := range shapes {
			if s.stats.n == 0 {
				continue
			}
			// Each violin fills its band at its own widest point.
			cx, half, top := band.Center(i), band.Width()/2, 0.0
			for _, d := range s.dens {
				top = max(top, d)
			}
			var outline ui.Path
			for j, y := range s.ys {
				x := cx + half*float32(s.dens[j]/top)
				if j == 0 {
					outline.MoveTo(x, ys.Map(y))
				} else {
					outline.LineTo(x, ys.Map(y))
				}
			}
			for j := len(s.ys) - 1; j >= 0; j-- {
				outline.LineTo(cx-half*float32(s.dens[j]/top), ys.Map(s.ys[j]))
			}
			outline.Close()
			col := pal[i%len(pal)]
			p.FillPath(&outline, col.Alpha(0.22))
			w := float32(1.5)
			if i == active {
				w = 3
			}
			p.StrokePath(&outline, w, col)
			// The quartile bar and a diamond at the median read without color.
			p.Fill(ui.Rect{X: cx - 3, Y: ys.Map(s.stats.q3), W: 6, H: max(ys.Map(s.stats.q1)-ys.Map(s.stats.q3), 1)}, k.Text, 1)
			p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, ys.Map(s.stats.median)), 10), k.Surface)
			p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, ys.Map(s.stats.median)), 7), k.Text)
		}
	}
	res, _ := ch.build(spec)
	return res
}
