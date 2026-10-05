package mujica

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// DensityPlotOptions configure a DensityPlot.
type DensityPlotOptions struct {
	Frame ChartFrame
	// Bandwidth defaults to the per-group Gaussian estimate.
	Bandwidth float64
	// Points defaults to 128 samples per curve.
	Points int
	Format func(v float64) string
}

// densitySample is finite values, sorted out of NaN and infinities.
func densitySample(values []float64) []float64 {
	var out []float64
	for _, v := range values {
		if math.IsInf(v, 0) {
			panic("mujica: density values must be finite or NaN")
		}
		if !math.IsNaN(v) {
			out = append(out, v)
		}
	}
	return out
}

// densityBandwidth estimates spread, using 1 for flat samples.
func densityBandwidth(v []float64) float64 {
	n := float64(len(v))
	mean := 0.0
	for _, x := range v {
		mean += x
	}
	mean /= n
	sd := 0.0
	for _, x := range v {
		sd += (x - mean) * (x - mean)
	}
	sd = math.Sqrt(sd / max(n-1, 1))
	if sd == 0 {
		return 1
	}
	return 1.06 * sd * math.Pow(n, -0.2)
}

// densityKDE estimates Gaussian density at xs.
func densityKDE(v []float64, h float64, xs []float64) []float64 {
	out := make([]float64, len(xs))
	norm := 1 / (float64(len(v)) * h * math.Sqrt(2*math.Pi))
	for i, x := range xs {
		s := 0.0
		for _, p := range v {
			u := (x - p) / h
			s += math.Exp(-0.5 * u * u)
		}
		out[i] = s * norm
	}
	return out
}

// densityCurves samples groups on a shared padded grid.
func densityCurves(groups []BoxGroup, bandwidth float64, points int) (xs []float64, series []ChartSeries) {
	lo, hi := math.Inf(1), math.Inf(-1)
	samples := make([][]float64, len(groups))
	widths := make([]float64, len(groups))
	for i, g := range groups {
		samples[i] = densitySample(g.Values)
		if len(samples[i]) == 0 {
			continue
		}
		widths[i] = bandwidth
		if widths[i] == 0 {
			widths[i] = densityBandwidth(samples[i])
		}
		for _, v := range samples[i] {
			lo, hi = min(lo, v-3*widths[i]), max(hi, v+3*widths[i])
		}
	}
	if math.IsInf(lo, 1) {
		return nil, nil
	}
	xs = make([]float64, points)
	for i := range xs {
		xs[i] = lo + (hi-lo)*float64(i)/float64(points-1)
	}
	for i, g := range groups {
		ys := make([]float64, points)
		if len(samples[i]) > 0 {
			ys = densityKDE(samples[i], widths[i], xs)
		}
		series = append(series, ChartSeries{Name: g.Name, Values: ys})
	}
	return xs, series
}

// DensityPlot overlays groups on a shared density axis.
func DensityPlot(c *ui.Context, groups []BoxGroup, opts DensityPlotOptions) ChartResult {
	if opts.Bandwidth < 0 || math.IsNaN(opts.Bandwidth) || math.IsInf(opts.Bandwidth, 0) || opts.Points < 0 || opts.Points == 1 {
		panic("mujica: DensityPlot needs Bandwidth ≥ 0 and Points of 0 or 2 or more")
	}
	points := opts.Points
	if points == 0 {
		points = 128
	}
	ch := newChart(c, opts.Frame)
	ids := ""
	seen := map[string]bool{}
	for _, g := range groups {
		if g.Name == "" || seen[g.Name] {
			panic("mujica: DensityPlot groups need unique names")
		}
		seen[g.Name] = true
		ids += g.Name + chartID(g.Values) + ";"
	}
	type curves struct {
		xs     []float64
		series []ChartSeries
	}
	cv := chartMemo(ch, "kde", fmt.Sprint(ids, opts.Bandwidth, points), func() curves {
		xs, s := densityCurves(groups, opts.Bandwidth, points)
		return curves{xs, s}
	})
	n := len(cv.xs)
	series := cv.series
	if n == 0 {
		series = nil
	}
	tops, bases := make([][]float64, len(series)), make([][]float64, len(series))
	for i, s := range series {
		tops[i], bases[i] = s.Values, make([]float64, n)
	}
	var ax chartX
	if n > 0 {
		ax = chartMemo(ch, "x", ids+fmt.Sprint(opts.Bandwidth, points), func() chartX { return newChartX(n, cv.xs, nil) })
	}
	lo, hi, ok := chartYRange(ch, series, true, nil)
	if !ok {
		n = 0
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return fmt.Sprintf("%.3g", v) }
	}
	return chartLines(ch, chartLineData{series: series, tops: tops, bases: bases, ax: ax, lo: lo, hi: hi, n: n, format: format})
}
