package mujica

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// QuotePoint is one tick of an intraday line: time, price and the volume traded since the last point.
type QuotePoint struct {
	Time   time.Time
	Price  float64
	Volume float64
}

// LineQuoteChartOptions configure a LineQuoteChart.
type LineQuoteChartOptions struct {
	Frame ChartFrame
	// PrevClose draws the baseline the day's change is measured from; required.
	PrevClose float64
	// RegularOpen and RegularClose bound the regular session; points outside are pre- and after-market. Zero draws one session.
	RegularOpen, RegularClose time.Time
	Decimals                  int
}

// lineQuoteSession names a point's session: -1 pre, 0 regular, 1 after.
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

// LineQuoteChart is an intraday line with the previous close as baseline, the average price, and shaded extended sessions.
func LineQuoteChart(c *ui.Context, points []QuotePoint, opts LineQuoteChartOptions) ChartResult {
	if !(opts.PrevClose > 0) {
		panic("mujica: LineQuoteChart needs a positive PrevClose")
	}
	if opts.RegularOpen.IsZero() != opts.RegularClose.IsZero() || opts.RegularClose.Before(opts.RegularOpen) {
		panic("mujica: LineQuoteChart needs both session bounds, in order")
	}
	dec := priceDecimals("LineQuoteChart", opts.Decimals)
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	avg := chartMemo(ch, "avg", chartID(points), func() []float64 {
		out := make([]float64, len(points))
		var pv, v float64
		for i, p := range points {
			if p.Time.IsZero() || i > 0 && !p.Time.After(points[i-1].Time) || !(p.Price > 0) || p.Volume < 0 {
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
	dom := chartMemo(ch, "y", chartID(points), func() [2]float64 {
		lo, hi := opts.PrevClose, opts.PrevClose
		for i, p := range points {
			lo, hi = min(lo, p.Price, avg[i]), max(hi, p.Price, avg[i])
		}
		// Symmetric around the previous close, so the baseline sits in the middle.
		d := max(hi-opts.PrevClose, opts.PrevClose-lo)
		return [2]float64{opts.PrevClose - d, opts.PrevClose + d}
	})
	size := FontSize(c, theme.CaptionSize)
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := chartInset(r, 60, 26)
		plot.W -= 44 // percent axis on the right
		var t0, t1 float64
		if len(points) > 0 {
			t0, t1 = float64(points[0].Time.Unix()), float64(points[len(points)-1].Time.Unix())
		}
		return plot, internal.Linear{D0: t0, D1: t1, R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: dom[0], D1: dom[1], R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	x := func(xs internal.Linear, i int) float32 { return xs.Map(float64(points[i].Time.Unix())) }
	spec := chartSpec{n: len(points), keys: []chartKey{
		{Name: msg(c, "line_quote_chart.price"), Color: k.AccentText, Shape: 0},
		{Name: msg(c, "line_quote_chart.average"), Color: k.Warning, Shape: 3},
	}}
	spec.pick = func(r ui.Rect, px, py float32) int {
		plot, xs, _ := geom(r)
		if !plot.Contains(px, py) {
			return -1
		}
		return internal.Nearest(len(points), func(i int) (float32, float32) { return x(xs, i), 0 }, px, 0, plot.W, true)
	}
	spec.readout = func(i int) string {
		p := points[i]
		chg := (p.Price/opts.PrevClose - 1) * 100
		parts := []string{p.Time.Format("15:04"), priceUngrouped(p.Price, dec), fmt.Sprintf("%+.2f%% %s", chg, trendWord(c, trendSign(chg))),
			msg(c, "line_quote_chart.average") + " " + priceUngrouped(avg[i], dec)}
		if s := lineQuoteSession(p.Time, opts); s != 0 {
			parts = append(parts, msg(c, map[int]string{-1: "line_quote_chart.pre", 1: "line_quote_chart.after"}[s]))
		}
		return strings.Join(parts, " · ")
	}
	spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		_, _, step := internal.Nice(dom[0], dom[1], 6)
		ticks := internal.Ticks(dom[0], dom[1], step)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		for _, t := range ticks {
			chartText(p, fmt.Sprintf("%+.2f%%", (t/opts.PrevClose-1)*100), plot.X+plot.W+6, ys.Map(t), size, trendColor(c, trendSign(t-opts.PrevClose)), 0, 0.5)
		}
		lineQuoteSessions(p, plot, xs, opts, size, c)
		var tx []float32
		var tl []string
		if len(points) > 1 {
			ts, layout := internal.TimeTicks(points[0].Time, points[len(points)-1].Time, max(2, int(plot.W/90)))
			for _, t := range ts {
				tx, tl = append(tx, xs.Map(float64(t.Unix()))), append(tl, t.Format(layout))
			}
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		by := internal.Snap(p, ys.Map(opts.PrevClose))
		for xx := plot.X; xx < plot.X+plot.W; xx += 6 {
			p.Line(xx, by, min(xx+3, plot.X+plot.W), by, 1, k.TextMuted)
		}
		price := make([]float64, len(points))
		for i, pt := range points {
			price[i] = pt.Price
		}
		at := func(i int) float32 { return x(xs, i) }
		for li, vals := range [][]float64{price, avg} {
			col := []ui.Color{k.AccentText, k.Warning}[li]
			chartRuns(vals, at, func(idx []int) {
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
			})
		}
		if active >= 0 {
			ax := internal.Snap(p, at(active))
			p.Line(ax, plot.Y, ax, plot.Y+plot.H, 1, k.TextMuted)
			internal.Marker(p, 0, ax, ys.Map(price[active]), 9, k.AccentText)
			internal.Marker(p, 3, ax, ys.Map(avg[active]), 9, k.Warning)
		}
	}
	res, _ := ch.build(spec)
	return res
}

// lineQuoteSessions shades and labels the pre- and after-market spans.
func lineQuoteSessions(p *ui.Painter, plot ui.Rect, xs internal.Linear, opts LineQuoteChartOptions, size float32, c *ui.Context) {
	if opts.RegularOpen.IsZero() {
		return
	}
	k := Tokens(c)
	x0, x1 := max(plot.X, xs.Map(float64(opts.RegularOpen.Unix()))), min(plot.X+plot.W, xs.Map(float64(opts.RegularClose.Unix())))
	for _, span := range []struct {
		a, b float32
		key  string
	}{{plot.X, x0, "line_quote_chart.pre"}, {x1, plot.X + plot.W, "line_quote_chart.after"}} {
		if span.b-span.a < 1 {
			continue
		}
		p.Fill(ui.Rect{X: span.a, Y: plot.Y, W: span.b - span.a, H: plot.H}, k.SurfaceHover.Alpha(0.6), 0)
		chartText(p, msg(c, span.key), span.a+4, plot.Y+2, size, k.TextMuted, 0, 0)
	}
}
