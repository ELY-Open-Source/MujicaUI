package mujica

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// AreaChartOptions configure an AreaChart.
type AreaChartOptions struct {
	Frame ChartFrame
	// X places point i at ascending X[i]; nil means i.
	X []float64
	// Times makes a time axis; set X or Times.
	Times []time.Time
	// Stacked piles visible series on each other.
	Stacked bool
	// Baseline is where areas start; zero by default.
	Baseline float64
	Format   func(v float64) string
}

// AreaChart fills series to the baseline, or stacked.
func AreaChart(c *ui.Context, series []ChartSeries, opts AreaChartOptions) ChartResult {
	n := chartCheckSeries(series)
	ch := newChart(c, opts.Frame)
	ax := chartMemo(ch, "x", [3]string{chartID(opts.X), chartID(opts.Times), chartID(series)}, func() chartX { return newChartX(n, opts.X, opts.Times) })
	deps := struct {
		id      string
		stacked bool
		base    float64
	}{chartSeriesID(ch, series), opts.Stacked, opts.Baseline}
	bands := chartMemo(ch, "bands", deps, func() [2][][]float64 { return areaBands(series, ch.hidden, opts.Stacked, opts.Baseline) })
	lo, hi, ok := chartYRange(ch, series, false, bands[0])
	if !ok {
		n = 0
	}
	lo, hi = min(lo, opts.Baseline), max(hi, opts.Baseline)
	return chartLines(ch, chartLineData{series: series, tops: bands[0], bases: bands[1], ax: ax, lo: lo, hi: hi, n: n, format: opts.Format})
}

// areaBands returns each series' top and bottom edges.
func areaBands(series []ChartSeries, hidden func(string) bool, stacked bool, base float64) [2][][]float64 {
	var tops, bases [][]float64
	var run []float64
	if len(series) > 0 {
		run = make([]float64, len(series[0].Values))
		for i := range run {
			run[i] = base
		}
	}
	for _, s := range series {
		top, bot := make([]float64, len(s.Values)), make([]float64, len(s.Values))
		for i, v := range s.Values {
			bot[i], top[i] = base, v
			if stacked {
				bot[i], top[i] = run[i], run[i]+v
			}
			if stacked && !math.IsNaN(v) && !hidden(s.Name) {
				run[i] += v
			}
		}
		tops, bases = append(tops, top), append(bases, bot)
	}
	return [2][][]float64{tops, bases}
}

func chartBool(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
