package mujica

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// ChartSeries is named values; NaN marks a missing one.
type ChartSeries struct {
	Name   string
	Values []float64
}

type chartMemoEntry struct {
	deps any
	val  any
}

// chartMemo recomputes f only when the comparable deps change.
func chartMemo[T any](ch *chart, name string, deps any, f func() T) T {
	m := ui.Local(ch.res.Element, "mujica.memo."+name, func() chartMemoEntry { return chartMemoEntry{} })
	if m.val == nil || m.deps != deps {
		m.deps, m.val = deps, f()
	}
	return m.val.(T)
}

// chartID identifies a slice by first element and length.
func chartID[T any](s []T) string {
	if len(s) == 0 {
		return "0"
	}
	return fmt.Sprintf("%p:%d", &s[0], len(s))
}

// chartSeriesID identifies series data and the hidden set together.
func chartSeriesID(ch *chart, series []ChartSeries) string {
	var b strings.Builder
	for _, s := range series {
		fmt.Fprintf(&b, "%s|%s|%v;", s.Name, chartID(s.Values), ch.hidden(s.Name))
	}
	return b.String()
}

// chartCheckSeries panics on duplicate names or unequal lengths.
func chartCheckSeries(series []ChartSeries) int {
	if len(series) == 0 {
		return 0
	}
	n := len(series[0].Values)
	seen := map[string]bool{}
	for _, s := range series {
		if s.Name == "" || seen[s.Name] || len(s.Values) != n {
			panic("mujica: chart series need unique names and equal lengths: " + s.Name)
		}
		seen[s.Name] = true
	}
	return n
}

// chartX is an ascending axis of numbers, times or indices.
type chartX struct {
	vals  []float64
	times bool
	// d0 and d1, when d1 > d0, fix the domain.
	d0, d1 float64
}

func newChartX(n int, x []float64, times []time.Time) chartX {
	switch {
	case x != nil && times != nil:
		panic("mujica: chart takes X or Times, not both")
	case x != nil && len(x) != n, times != nil && len(times) != n:
		panic("mujica: chart X or Times length differs from the series")
	}
	a := chartX{vals: x, times: times != nil}
	if times != nil {
		a.vals = make([]float64, n)
		for i, t := range times {
			a.vals[i] = float64(t.UnixNano()) / 1e9
		}
	}
	if a.vals == nil {
		a.vals = make([]float64, n)
		for i := range a.vals {
			a.vals[i] = float64(i)
		}
	}
	if !sort.Float64sAreSorted(a.vals) {
		panic("mujica: chart X or Times must ascend")
	}
	return a
}

func (a chartX) scale(plot ui.Rect) internal.Linear {
	d0, d1 := a.vals[0], a.vals[len(a.vals)-1]
	if a.d1 > a.d0 {
		d0, d1 = a.d0, a.d1
	}
	return internal.Linear{D0: d0, D1: d1, R0: float64(plot.X), R1: float64(plot.X + plot.W)}
}

func (a chartX) label(i int) string {
	if a.times {
		return chartTime(a.vals[i]).Format("2006-01-02 15:04:05")
	}
	return strconv.FormatFloat(a.vals[i], 'g', 6, 64)
}

// ticks returns tick positions and labels along plot.
func (a chartX) ticks(plot ui.Rect) ([]float32, []string) {
	s := a.scale(plot)
	n := max(2, int(plot.W/90))
	var xs []float32
	var labels []string
	if a.times {
		ts, layout := internal.TimeTicks(chartTime(s.D0), chartTime(s.D1), n)
		for _, t := range ts {
			xs = append(xs, s.Map(float64(t.UnixNano())/1e9))
			labels = append(labels, t.Format(layout))
		}
		return xs, labels
	}
	lo, hi, step := internal.Nice(s.D0, s.D1, n)
	for _, t := range internal.Ticks(lo, hi, step) {
		if t >= s.D0 && t <= s.D1 {
			xs = append(xs, s.Map(t))
			labels = append(labels, internal.FormatTick(t, step))
		}
	}
	return xs, labels
}

// nearest returns the index whose x is closest to v.
func (a chartX) nearest(v float64) int {
	i := sort.SearchFloat64s(a.vals, v)
	if i == len(a.vals) || i > 0 && v-a.vals[i-1] < a.vals[i]-v {
		i--
	}
	return i
}

func chartTime(sec float64) time.Time {
	return time.Unix(0, int64(math.Round(sec*1e9))).UTC()
}

// chartYDomain spans visible finite values, optionally with zero.
func chartYDomain(series []ChartSeries, hidden func(string) bool, withZero bool) (lo, hi float64, ok bool) {
	lo, hi = math.Inf(1), math.Inf(-1)
	for _, s := range series {
		if hidden(s.Name) {
			continue
		}
		for _, v := range s.Values {
			if !math.IsNaN(v) {
				lo, hi = min(lo, v), max(hi, v)
			}
		}
	}
	if math.IsInf(lo, 1) {
		return 0, 0, false
	}
	if withZero {
		lo, hi = min(lo, 0), max(hi, 0)
	}
	return lo, hi, true
}

// chartValue formats v, or names a missing value.
func chartValue(c *ui.Context, v float64, format func(float64) string) string {
	if math.IsNaN(v) {
		return msg(c, "chart.missing")
	}
	if format != nil {
		return format(v)
	}
	return strconv.FormatFloat(v, 'g', 6, 64)
}

// chartRuns calls f per finite run, thinned to pixel columns.
func chartRuns(vals []float64, x func(i int) float32, f func(idx []int)) {
	for i := 0; i < len(vals); {
		if math.IsNaN(vals[i]) {
			i++
			continue
		}
		j := i
		for j < len(vals) && !math.IsNaN(vals[j]) {
			j++
		}
		run := internal.Columns(j-i, func(k int) float32 { return x(i + k) }, func(k int) float64 { return vals[i+k] })
		for k := range run {
			run[k] += i
		}
		f(run)
		i = j
	}
}

func chartVisibleIndices(n, limit, active int) []int {
	step := max(1, (n+limit-1)/limit)
	out := make([]int, 0, min(n, limit+1))
	for i := 0; i < n; i += step {
		out = append(out, i)
	}
	if active >= 0 && active < n && active%step != 0 {
		out = append(out, active)
	}
	return out
}
