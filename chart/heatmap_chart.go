package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// HeatmapChartOptions configure a HeatmapChart.
type HeatmapChartOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// heatSteps is the number of color steps.
const heatSteps = 5

// HeatmapChart colors a matrix by value; NaN cells are dashed.
func HeatmapChart(c *ui.Context, rows, cols []string, values [][]float64, opts HeatmapChartOptions) ChartResult {
	if len(values) != len(rows) {
		panic("mujica: HeatmapChart needs one value row per row label")
	}
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, r := range values {
		if len(r) != len(cols) {
			panic("mujica: HeatmapChart needs one value per column label")
		}
		for _, v := range r {
			if !math.IsNaN(v) {
				lo, hi = min(lo, v), max(hi, v)
			}
		}
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	size := core.FontSize(c, theme.CaptionSize)
	n := len(rows) * len(cols)
	if math.IsInf(lo, 1) {
		n = 0
	}
	// level is v's color step, 0 to heatSteps-1.
	level := func(v float64) int {
		if hi == lo {
			return heatSteps - 1
		}
		return min(heatSteps-1, int((v-lo)/(hi-lo)*heatSteps))
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Band, internal.Band) {
		plot := ui.Rect{X: r.X + 72, Y: r.Y + 4, W: max(r.W-72-64, 1), H: max(r.H-28, 1)}
		return plot, internal.Band{N: len(cols), R0: plot.X, R1: plot.X + plot.W, Pad: 0.06}, internal.Band{N: len(rows), R0: plot.Y, R1: plot.Y + plot.H, Pad: 0.06}
	}
	cell := func(xb, yb internal.Band, i int) ui.Rect {
		return ui.Rect{X: xb.Start(i % len(cols)), Y: yb.Start(i / len(cols)), W: xb.Width(), H: yb.Width()}
	}
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		return rows[i/len(cols)] + " · " + cols[i%len(cols)] + " · " + ChartValue(c, values[i/len(cols)][i%len(cols)], opts.Format)
	}
	spec.Move = func(i int, key ui.Key) int {
		r, col := i/len(cols), i%len(cols)
		switch key {
		case ui.KeyLeft:
			col = max(col-1, 0)
		case ui.KeyRight:
			col = min(col+1, len(cols)-1)
		case ui.KeyUp:
			r = max(r-1, 0)
		default:
			r = min(r+1, len(rows)-1)
		}
		return r*len(cols) + col
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, xb, yb := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		col := min(len(cols)-1, int((x-plot.X)/xb.Step()))
		row := min(len(rows)-1, int((y-plot.Y)/yb.Step()))
		return row*len(cols) + col
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xb, yb := geom(r)
		for i := range len(rows) * len(cols) {
			b := cell(xb, yb, i)
			v := values[i/len(cols)][i%len(cols)]
			if math.IsNaN(v) {
				p.StrokeDashed(b, k.ControlBorder, 2, 1)
			} else {
				bg := internal.Ramp(k, (float64(level(v))+0.5)/heatSteps)
				p.Fill(b, bg, 2)
				if b.W >= 30 && b.H >= 18 {
					ChartText(p, ChartValue(c, v, opts.Format), b.X+b.W/2, b.Y+b.H/2, size, ChartInk(k, bg), 0.5, 0.5)
				}
			}
			if i == active {
				p.Stroke(ui.Rect{X: b.X - 2, Y: b.Y - 2, W: b.W + 4, H: b.H + 4}, k.Text, 3, 2)
			}
		}
		for i, name := range rows {
			ChartText(p, name, plot.X-8, yb.Center(i), size, k.TextMuted, 1, 0.5)
		}
		for i, name := range cols {
			ChartText(p, name, xb.Center(i), plot.Y+plot.H+4, size, k.TextMuted, 0.5, 0)
		}
		heatLegend(p, ui.Rect{X: plot.X + plot.W + 16, Y: plot.Y, W: 14, H: plot.H}, lo, hi, size, k, opts.Format, c)
	}
	res, _ := ch.Build(spec)
	return res
}

// heatLegend stacks the color steps with the range at ends.
func heatLegend(p *ui.Painter, r ui.Rect, lo, hi float64, size float32, k theme.Tokens, format func(float64) string, c *ui.Context) {
	h := min(r.H, 120) / heatSteps
	for i := range heatSteps {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + float32(heatSteps-1-i)*h, W: r.W, H: h - 1}, internal.Ramp(k, (float64(i)+0.5)/heatSteps), 1)
	}
	ChartText(p, ChartValue(c, hi, format), r.X+r.W+4, r.Y, size, k.TextMuted, 0, 0)
	ChartText(p, ChartValue(c, lo, format), r.X+r.W+4, r.Y+h*heatSteps, size, k.TextMuted, 0, 1)
}
