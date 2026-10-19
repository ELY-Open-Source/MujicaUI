package chart

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// QuotePoint holds time, price, and incremental volume.
type QuotePoint struct {
	Time   time.Time
	Price  float64
	Volume float64
}

// LineQuoteChartOptions configure a LineQuoteChart.
type LineQuoteChartOptions struct {
	Frame ChartFrame
	// PrevClose is the required positive baseline price.
	PrevClose float64
	// RegularOpen and RegularClose optionally delimit extended trading sessions.
	RegularOpen, RegularClose time.Time
	Decimals                  int
}

// lineQuoteSession returns pre-market, regular, or after-market.
func lineQuoteSession(t time.Time, opts LineQuoteChartOptions) int {
	switch {
	case opts.RegularOpen.IsZero():
		return 0
	case t.Before(opts.RegularOpen):
		return -1
	case !t.Before(opts.RegularClose):
		return 1
	}
	return 0
}

// LineQuoteChart draws price, weighted average, and extended sessions.
func LineQuoteChart(c *ui.Context, points []QuotePoint, opts LineQuoteChartOptions) ChartResult {
	if math.IsInf(opts.PrevClose, 0) || !(opts.PrevClose > 0) {
		panic("mujica: LineQuoteChart needs a positive PrevClose")
	}
	if opts.RegularOpen.IsZero() != opts.RegularClose.IsZero() || opts.RegularClose.Before(opts.RegularOpen) {
		panic("mujica: LineQuoteChart needs both session bounds, in order")
	}
	dec := core.PriceDecimals("LineQuoteChart", opts.Decimals)
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	avg := ChartMemo(ch, "avg", ChartID(points), func() []float64 {
		out := make([]float64, len(points))
		var pv, v float64
		for i, p := range points {
			if p.Time.IsZero() || i > 0 && !p.Time.After(points[i-1].Time) || !(p.Price > 0) || p.Volume < 0 || math.IsNaN(p.Volume) || math.IsInf(p.Volume, 0) || math.IsInf(p.Price, 0) {
				panic(fmt.Sprintf("mujica: LineQuoteChart point %d is out of order or invalid", i))
			}
			pv += p.Price * p.Volume
			v += p.Volume
			out[i] = p.Price
			if v > 0 {
				out[i] = pv / v
			}
		}
		return out
	})
	dom := ChartMemo(ch, "y", [2]any{ChartID(points), opts.PrevClose}, func() [2]float64 {
		lo, hi := opts.PrevClose, opts.PrevClose
		for i, p := range points {
			lo, hi = min(lo, p.Price, avg[i]), max(hi, p.Price, avg[i])
		}
		// Center the scale on the previous close.
		d := max(hi-opts.PrevClose, opts.PrevClose-lo)
		return [2]float64{opts.PrevClose - d, opts.PrevClose + d}
	})
	size := core.FontSize(c, theme.CaptionSize)
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := ChartInset(r, 60, 26)
		plot.W -= 44 // percent axis on the right
		var t0, t1 float64
		if len(points) > 0 {
			t0, t1 = float64(points[0].Time.Unix()), float64(points[len(points)-1].Time.Unix())
		}
		return plot, internal.Linear{D0: t0, D1: t1, R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: dom[0], D1: dom[1], R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	x := func(xs internal.Linear, i int) float32 { return xs.Map(float64(points[i].Time.Unix())) }
	spec := ChartSpec{Marks: len(points), Legend: []ChartLegendItem{
		{Name: core.Msg(c, "line_quote_chart.price"), Color: k.AccentText, Shape: 0},
		{Name: core.Msg(c, "line_quote_chart.average"), Color: k.Warning, Shape: 3},
	}}
	spec.Pick = func(r ui.Rect, px, py float32) int {
		plot, xs, _ := geom(r)
		if !plot.Contains(px, py) {
			return -1
		}
		return internal.Nearest(len(points), func(i int) (float32, float32) { return x(xs, i), 0 }, px, 0, plot.W, true)
	}
	spec.Readout = func(i int) string {
		p := points[i]
		chg := (p.Price/opts.PrevClose - 1) * 100
		parts := []string{p.Time.Format("15:04"), core.PriceUngrouped(p.Price, dec), fmt.Sprintf("%+.2f%% %s", chg, core.TrendWord(c, core.TrendSign(chg))),
			core.Msg(c, "line_quote_chart.average") + " " + core.PriceUngrouped(avg[i], dec)}
		if s := lineQuoteSession(p.Time, opts); s != 0 {
			parts = append(parts, core.Msg(c, map[int]string{-1: "line_quote_chart.pre", 1: "line_quote_chart.after"}[s]))
		}
		return strings.Join(parts, " · ")
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		_, _, step := internal.Nice(dom[0], dom[1], 6)
		ticks := internal.Ticks(dom[0], dom[1], step)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		for _, t := range ticks {
			ChartText(p, fmt.Sprintf("%+.2f%%", (t/opts.PrevClose-1)*100), plot.X+plot.W+6, ys.Map(t), size, core.TrendColor(c, core.TrendSign(t-opts.PrevClose)), 0, 0.5)
		}
		lineQuoteSessions(p, plot, xs, opts, size, c)
		var tx []float32
		var tl []string
		if len(points) > 1 {
			ts, layout := internal.TimeTicks(points[0].Time, points[len(points)-1].Time, max(2, int(plot.W/90)))
			for _, t := range ts {
				if layout == "Jan 2" {
					layout = core.Msg(c, "candlestick_chart.daily_time")
				}
				tx, tl = append(tx, xs.Map(float64(t.Unix()))), append(tl, t.Format(layout))
			}
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		by := internal.Snap(p, ys.Map(opts.PrevClose))
		for xx := plot.X; xx < plot.X+plot.W; xx += 6 {
			p.Line(xx, by, min(xx+3, plot.X+plot.W), by, 1, k.TextMuted)
		}
		price := ChartMemo(ch, "prices", ChartID(points), func() []float64 {
			out := make([]float64, len(points))
			for i, pt := range points {
				out[i] = pt.Price
			}
			return out
		})
		at := func(i int) float32 { return x(xs, i) }
		for li, vals := range [][]float64{price, avg} {
			col := []ui.Color{k.AccentText, k.Warning}[li]
			runs := ChartMemo(ch, fmt.Sprint("quote.runs.", li), [3]any{ChartID(points), plot.X, plot.W}, func() [][]int {
				var out [][]int
				chartRuns(vals, at, func(idx []int) { out = append(out, idx) })
				return out
			})
			for _, idx := range runs {
				var line ui.Path
				for j, i := range idx {
					px, py := internal.Snap(p, at(i)), internal.Snap(p, ys.Map(vals[i]))
					if j == 0 {
						line.MoveTo(px, py)
					} else {
						line.LineTo(px, py)
					}
				}
				p.StrokePath(&line, []float32{1.75, 1.25}[li], col)
				for j := 0; j < len(idx); j += max(1, len(idx)/12) {
					i := idx[j]
					internal.Marker(p, []int{0, 3}[li], at(i), ys.Map(vals[i]), 5, col)
				}
			}
		}
		if active >= 0 {
			ax := internal.Snap(p, at(active))
			p.Line(ax, plot.Y, ax, plot.Y+plot.H, 1, k.TextMuted)
			internal.Marker(p, 0, ax, ys.Map(price[active]), 9, k.AccentText)
			internal.Marker(p, 3, ax, ys.Map(avg[active]), 9, k.Warning)
		}
	}
	res, _ := ch.Build(spec)
	return res
}

// lineQuoteSessions shades and labels the pre- and after-market spans.
func lineQuoteSessions(p *ui.Painter, plot ui.Rect, xs internal.Linear, opts LineQuoteChartOptions, size float32, c *ui.Context) {
	if opts.RegularOpen.IsZero() {
		return
	}
	k := core.Tokens(c)
	x0, x1 := max(plot.X, xs.Map(float64(opts.RegularOpen.Unix()))), min(plot.X+plot.W, xs.Map(float64(opts.RegularClose.Unix())))
	for _, span := range []struct {
		a, b float32
		key  string
	}{{plot.X, x0, "line_quote_chart.pre"}, {x1, plot.X + plot.W, "line_quote_chart.after"}} {
		if span.b-span.a < 1 {
			continue
		}
		p.Fill(ui.Rect{X: span.a, Y: plot.Y, W: span.b - span.a, H: plot.H}, k.SurfaceHover.Alpha(0.6), 0)
		ChartText(p, core.Msg(c, span.key), span.a+4, plot.Y+2, size, k.TextMuted, 0, 0)
	}
}
