package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ViolinPlotOptions configure a ViolinPlot.
type ViolinPlotOptions struct {
	Frame ChartFrame
	// Bandwidth defaults to the per-group estimate.
	Bandwidth float64
	Format    func(v float64) string
}

// violinShape stores density samples and quartiles.
type violinShape struct {
	ys, dens []float64
	stats    boxStats
}

// violinPoints is how many values each outline is sampled at.
const violinPoints = 64

// ViolinPlot compares density outlines, medians, and quartiles.
func ViolinPlot(c *ui.Context, groups []BoxGroup, opts ViolinPlotOptions) ChartResult {
	if opts.Bandwidth < 0 || math.IsNaN(opts.Bandwidth) || math.IsInf(opts.Bandwidth, 0) {
		panic("mujica: ViolinPlot Bandwidth must be ≥ 0")
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	ids := ""
	seen := map[string]bool{}
	for _, g := range groups {
		if g.Name == "" || seen[g.Name] {
			panic("mujica: ViolinPlot groups need names")
		}
		seen[g.Name] = true
		ids += ChartID(g.Values) + ";"
	}
	shapes := ChartMemo(ch, "violins", fmt.Sprint(ids, opts.Bandwidth), func() []violinShape {
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
			if lo == hi {
				lo, hi = lo-3*h, hi+3*h
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
	val := func(v float64) string { return ChartValue(c, v, opts.Format) }
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Linear) {
		plot := ChartInset(r, 52, 26)
		return plot, internal.Band{N: len(groups), R0: plot.X, R1: plot.X + plot.W, Pad: 0.15},
			internal.Linear{D0: lo, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	visible := ui.Local(ch.Res.Element, "mujica.violin.visible", func() []int { return nil })
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		s := shapes[i].stats
		if s.n == 0 {
			return groups[i].Name + " · " + core.Msg(c, "chart.missing")
		}
		return fmt.Sprintf(core.Msg(c, "violin_plot.readout"), groups[i].Name, val(s.median), val(s.q1), val(s.q3), s.n)
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
		*visible = chartVisibleIndices(len(groups), max(1, int(plot.W)), active)
		for _, i := range *visible {
			xs, names = append(xs, band.Center(i)), append(names, groups[i].Name)
		}
		internal.AxisX(p, plot, xs, names, size, k)
		for _, i := range *visible {
			s := shapes[i]
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
				x := internal.Snap(p, cx+half*float32(s.dens[j]/top))
				py := internal.Snap(p, ys.Map(y))
				if j == 0 {
					outline.MoveTo(x, py)
				} else {
					outline.LineTo(x, py)
				}
			}
			for j := len(s.ys) - 1; j >= 0; j-- {
				outline.LineTo(internal.Snap(p, cx-half*float32(s.dens[j]/top)), internal.Snap(p, ys.Map(s.ys[j])))
			}
			outline.Close()
			col := pal[i%len(pal)]
			p.FillPath(&outline, col.Alpha(0.22))
			w := float32(1.5)
			if i == active {
				w = 3
			}
			p.StrokePath(&outline, w, col)
			// Quartile bars and median diamonds distinguish values.
			p.Fill(ui.Rect{X: cx - 3, Y: ys.Map(s.stats.q3), W: 6, H: max(ys.Map(s.stats.q1)-ys.Map(s.stats.q3), 1)}, k.Text, 1)
			p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, ys.Map(s.stats.median)), 10), k.Surface)
			p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, ys.Map(s.stats.median)), 7), k.Text)
		}
	}
	res, _ := ch.Build(spec)
	return res
}
