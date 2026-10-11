package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// HistogramOptions configure a Histogram.
type HistogramOptions struct {
	Frame ChartFrame
	// Bins is the bin count; 0 uses Sturges' rule.
	Bins     int
	BinWidth float64
	// Cumulative counts every value up to each bin's upper edge.
	Cumulative bool
	Format     func(v float64) string
}

// histogramMaxBins bounds the bins a width may make.
const histogramMaxBins = 10_000

type histBins struct {
	start, width float64
	counts       []int
}

// histogramBin counts finite values into equal bins.
func histogramBin(values []float64, bins int, width float64, cumulative bool) histBins {
	lo, hi := math.Inf(1), math.Inf(-1)
	n := 0
	for _, v := range values {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			lo, hi, n = min(lo, v), max(hi, v), n+1
		}
	}
	if n == 0 {
		return histBins{}
	}
	var b histBins
	switch {
	case width > 0:
		b.start, b.width = math.Floor(lo/width)*width, width
		bins = max(1, int(math.Ceil((hi-b.start)/width-1e-9)))
		if hi == b.start+float64(bins)*width {
			bins++
		}
	default:
		if bins == 0 {
			bins = int(math.Ceil(math.Log2(float64(n)))) + 1
		}
		b.start, b.width = lo, (hi-lo)/float64(bins)
		if b.width == 0 {
			b.width = 1
		}
	}
	if bins > histogramMaxBins {
		panic("mujica: Histogram would make more than 10000 bins")
	}
	b.counts = make([]int, bins)
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}
		b.counts[min(bins-1, int((v-b.start)/b.width))]++
	}
	if cumulative {
		for i := 1; i < bins; i++ {
			b.counts[i] += b.counts[i-1]
		}
	}
	return b
}

// Histogram counts values into equal bins drawn as bars.
func Histogram(c *ui.Context, values []float64, opts HistogramOptions) ChartResult {
	if opts.Bins < 0 || opts.BinWidth < 0 || math.IsNaN(opts.BinWidth) {
		panic("mujica: Histogram Bins and BinWidth must be ≥ 0")
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	col := internal.Palette(k)[0]
	size := core.FontSize(c, theme.CaptionSize)
	b := ChartMemo(ch, "bins", [4]any{ChartID(values), opts.Bins, opts.BinWidth, opts.Cumulative}, func() histBins {
		return histogramBin(values, opts.Bins, opts.BinWidth, opts.Cumulative)
	})
	top := 0
	for _, n := range b.counts {
		top = max(top, n)
	}
	var ys internal.Linear
	var hi, step float64
	if top > 0 {
		_, hi, step = internal.Nice(0, float64(top), 4)
		step = max(step, 1)
	}
	// edge writes bin edge i at a given precision.
	edge := func(i int, fine float64) string {
		v := b.start + float64(i)*b.width
		if opts.Format != nil {
			return opts.Format(v)
		}
		s := internal.FormatTick(v, b.width/fine)
		if strings.Contains(s, ".") {
			s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
		}
		return s
	}
	geom := func(r ui.Rect) (ui.Rect, internal.Band) {
		plot := ChartInset(r, 52, 26)
		ys = internal.Linear{D0: 0, D1: hi, R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
		return plot, internal.Band{N: len(b.counts), R0: plot.X, R1: plot.X + plot.W, Pad: 0.06}
	}
	spec := ChartSpec{Marks: len(b.counts)}
	spec.Readout = func(i int) string {
		if opts.Cumulative {
			return fmt.Sprintf("≤ %s · %d", edge(i+1, 100), b.counts[i])
		}
		closer := ")"
		if i == len(b.counts)-1 && opts.BinWidth == 0 {
			closer = "]"
		}
		return fmt.Sprintf("[%s, %s%s · %d", edge(i, 100), edge(i+1, 100), closer, b.counts[i])
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		plot, band := geom(r)
		if !plot.Contains(x, y) {
			return -1
		}
		return min(len(b.counts)-1, int((x-plot.X)/band.Step()))
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, band := geom(r)
		internal.AxisY(p, plot, ys, internal.Ticks(0, hi, step), func(v float64) string { return internal.FormatTick(v, step) }, size, k)
		var xs []float32
		var labels []string
		for i := 0; i <= len(b.counts); i++ {
			xs, labels = append(xs, plot.X+band.Step()*float32(i)), append(labels, edge(i, 1))
		}
		internal.AxisX(p, plot, xs, labels, size, k)
		for i, n := range b.counts {
			y := ys.Map(float64(n))
			bar := ui.Rect{X: band.Start(i), Y: y, W: band.Width(), H: plot.Y + plot.H - y}
			p.Fill(bar, col, 1)
			if i == active {
				p.Stroke(ui.Rect{X: bar.X - 2, Y: bar.Y - 2, W: bar.W + 4, H: bar.H + 4}, k.Text, 2, 1.5)
			}
		}
	}
	res, _ := ch.Build(spec)
	return res
}
