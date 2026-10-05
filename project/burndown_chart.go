package project

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// BurndownChartOptions configure a BurndownChart.
type BurndownChartOptions struct {
	Frame chart.ChartFrame
	// Start and End are the sprint's first and last days.
	Start, End core.Date
	// Total starts the ideal line, which reaches zero.
	Total float64
	// Format writes amounts; nil uses up to two decimals.
	Format func(v float64) string
	// Version changes when remaining work changes in place.
	Version int
}

// BurndownChart compares remaining work with the ideal line.
func BurndownChart(c *ui.Context, remaining []float64, opts BurndownChartOptions) chart.ChartResult {
	core.DatetimeCheckDate("BurndownChart Start", opts.Start)
	core.DatetimeCheckDate("BurndownChart End", opts.End)
	if opts.Start.IsZero() || opts.End.Compare(opts.Start) <= 0 || !(opts.Total > 0) || math.IsInf(opts.Total, 0) {
		panic(fmt.Sprintf("mujica: BurndownChart needs Start < End and a positive Total, got %v..%v %v", opts.Start, opts.End, opts.Total))
	}
	days := int((opts.End.Civil().Unix()-opts.Start.Civil().Unix())/86400) + 1
	if len(remaining) > days {
		panic(fmt.Sprintf("mujica: BurndownChart has %d values for %d days", len(remaining), days))
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string {
			return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', 2, 64), "0"), ".")
		}
	}
	ch := chart.NewChart(c, opts.Frame)
	summary := chart.ChartMemo(ch, "burndown-domain", struct {
		data    string
		total   float64
		version int
	}{chart.ChartID(remaining), opts.Total, opts.Version}, func() float64 {
		hi := opts.Total
		for _, v := range remaining {
			if v < 0 || math.IsInf(v, 0) {
				panic("mujica: BurndownChart amounts must be non-negative and finite or NaN")
			}
			if !math.IsNaN(v) {
				hi = max(hi, v)
			}
		}
		return hi
	})
	cache := ui.Local(ch.Res.Element, "burndown-samples", func() burndownCache { return burndownCache{} })
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	ideal := func(i int) float64 { return opts.Total * (1 - float64(i)/float64(days-1)) }
	actual := func(i int) float64 {
		if i < len(remaining) {
			return remaining[i]
		}
		return math.NaN()
	}
	lo, top, step := internal.Nice(0, summary, 4)
	ticks := internal.Ticks(lo, top, step)
	keys := []chart.ChartLegendItem{{Name: core.Msg(c, "burndown_chart.actual"), Color: pal[0], Shape: 0}, {Name: core.Msg(c, "burndown_chart.ideal"), Color: pal[5], Shape: 1}}
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := chart.ChartInset(r, 48, 26)
		return plot, internal.Linear{D0: 0, D1: float64(days - 1), R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: lo, D1: top, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	n := days
	if len(remaining) == 0 {
		n = 0
	}
	spec := chart.ChartSpec{Marks: n, Legend: keys, Hidden: ch.Toggles(true)}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, xs, _ := geom(r)
		if x < plot.X-8 || x > plot.X+plot.W+8 {
			return -1
		}
		return max(0, min(days-1, int(math.Round(xs.Invert(x)))))
	}
	spec.Readout = func(i int) string {
		day := opts.Start.AddDays(i)
		parts := core.DatetimeMonthName(c, day.Month) + " " + strconv.Itoa(day.Day)
		if v := actual(i); !math.IsNaN(v) && !ch.Hidden(keys[0].Name) {
			parts += " · " + keys[0].Name + " " + format(v)
		}
		if !ch.Hidden(keys[1].Name) {
			parts += " · " + keys[1].Name + " " + format(math.Round(ideal(i)*100)/100)
		}
		return parts
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, xs, ys := geom(r)
		internal.AxisY(p, plot, ys, ticks, func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		var tx []float32
		var tl []string
		stride := max(1, int(math.Ceil(float64(days)/max(float64(plot.W)/70, 1))))
		for i := 0; i < days; i += stride {
			tx = append(tx, xs.Map(float64(i)))
			tl = append(tl, strconv.Itoa(opts.Start.AddDays(i).Day))
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		pt := func(i int, v float64) (float32, float32) {
			return internal.Snap(p, xs.Map(float64(i))), internal.Snap(p, ys.Map(v))
		}
		if !ch.Hidden(keys[1].Name) {
			x0, y0 := pt(0, ideal(0))
			x1, y1 := pt(days-1, 0)
			burndownDashed(p, x0, y0, x1, y1, keys[1].Color)
		}
		if !ch.Hidden(keys[0].Name) {
			var line ui.Path
			open := false
			width := max(1, int(plot.W*p.Scale()))
			data := chart.ChartID(remaining)
			if cache.data != data || cache.version != opts.Version || cache.days != days || cache.width != width {
				cache.data, cache.version, cache.days, cache.width = data, opts.Version, days, width
				cache.points = burndownSamples(remaining, days, width)
			}
			lastRun := -1
			for _, point := range cache.points {
				x, y := pt(point.index, remaining[point.index])
				if open && lastRun == point.run {
					line.LineTo(x, y)
				} else {
					line.MoveTo(x, y)
					open = true
				}
				lastRun = point.run
				internal.Marker(p, 0, x, y, 7, keys[0].Color)
			}
			p.StrokePath(&line, 2, keys[0].Color)
		}
		if active >= 0 {
			ax, _ := pt(active, 0)
			p.Line(ax, plot.Y, ax, plot.Y+plot.H, 1, k.TextMuted)
		}
	}
	res, _ := ch.Build(spec)
	return res
}

// burndownDashed draws the ideal line with a distinct pattern.
func burndownDashed(p *ui.Painter, x0, y0, x1, y1 float32, col ui.Color) {
	l := float32(math.Hypot(float64(x1-x0), float64(y1-y0)))
	var path ui.Path
	for t := float32(0); t < l; t += 10 {
		e := min(t+6, l)
		path.MoveTo(internal.Snap(p, x0+(x1-x0)*t/l), internal.Snap(p, y0+(y1-y0)*t/l))
		path.LineTo(internal.Snap(p, x0+(x1-x0)*e/l), internal.Snap(p, y0+(y1-y0)*e/l))
	}
	p.StrokePath(&path, 1.5, col)
}

type burndownPoint struct{ index, run int }
type burndownCache struct {
	data                 string
	version, days, width int
	points               []burndownPoint
}

// burndownSamples keeps four finite points per device pixel.
func burndownSamples(values []float64, days, width int) []burndownPoint {
	var out []burndownPoint
	run := 0
	for i := 0; i < len(values); {
		bucket := i * width / max(days-1, 1)
		var first, low, high, last burndownPoint
		count := 0
		for i < len(values) && i*width/max(days-1, 1) == bucket {
			v := values[i]
			if math.IsNaN(v) {
				run++
				i++
				continue
			}
			point := burndownPoint{index: i, run: run}
			if count == 0 {
				first, low, high = point, point, point
			}
			if v < values[low.index] {
				low = point
			}
			if v > values[high.index] {
				high = point
			}
			last = point
			count++
			i++
		}
		if count == 0 {
			continue
		}
		points := []burndownPoint{first, low, high, last}
		slices.SortFunc(points, func(a, b burndownPoint) int { return a.index - b.index })
		for j, point := range points {
			if j == 0 || points[j-1].index != point.index {
				out = append(out, point)
			}
		}
	}
	return out
}
