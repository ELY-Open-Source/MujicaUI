package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// FunnelChartOptions configure a FunnelChart.
type FunnelChartOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// FunnelChart narrows by stage, labeled with conversion rates.
func FunnelChart(c *ui.Context, stages []ChartSlice, opts FunnelChartOptions) ChartResult {
	top := 0.0
	for _, s := range stages {
		if s.Label == "" || s.Value < 0 || math.IsNaN(s.Value) || math.IsInf(s.Value, 0) {
			panic("mujica: FunnelChart stages need a label and a finite value ≥ 0")
		}
		top = max(top, s.Value)
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	col := internal.Palette(k)[0]
	size := core.FontSize(c, theme.CaptionSize)
	n := len(stages)
	if top == 0 {
		n = 0
	}
	pct := func(a, b float64) string {
		if b == 0 {
			return "—"
		}
		return fmt.Sprintf("%.1f%%", 100*a/b)
	}
	rate := func(i int) string {
		if i == 0 {
			return pct(stages[0].Value, stages[0].Value)
		}
		return pct(stages[i].Value, stages[i-1].Value)
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Band) {
		plot := ui.Rect{X: r.X + 8, Y: r.Y + 4, W: r.W * 0.58, H: r.H - 8}
		return plot, internal.Band{N: len(stages), R0: plot.Y, R1: plot.Y + plot.H, Pad: 0.08}
	}
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		s := stages[i]
		return fmt.Sprintf(core.Msg(c, "funnel_chart.readout"), s.Label, ChartValue(c, s.Value, opts.Format), rate(i), pct(s.Value, stages[0].Value))
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, band := geom(r)
		if y < plot.Y || y >= plot.Y+plot.H || x < r.X || x > r.X+r.W {
			return -1
		}
		return min(len(stages)-1, int((y-plot.Y)/band.Step()))
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band := geom(r)
		cx := plot.X + plot.W/2
		width := func(i int) float32 { return max(plot.W*float32(stages[min(i, len(stages)-1)].Value/top), 2) }
		for i, s := range stages {
			y0, y1 := band.Start(i), band.Start(i)+band.Width()
			w0, w1 := width(i), width(i+1)
			var shape ui.Path
			shape.MoveTo(cx-w0/2, y0).LineTo(cx+w0/2, y0).LineTo(cx+w1/2, y1).LineTo(cx-w1/2, y1).Close()
			p.FillPath(&shape, col)
			if i == active {
				p.StrokePath(&shape, 2, k.Text)
			}
			ly := (y0 + y1) / 2
			lx := plot.X + plot.W + 16
			ChartText(p, s.Label+" · "+ChartValue(c, s.Value, opts.Format), lx, ly, size, k.Text, 0, 1)
			ChartText(p, rate(i), lx, ly+2, size, k.TextMuted, 0, 0)
		}
	}
	res, _ := ch.Build(spec)
	return res
}
