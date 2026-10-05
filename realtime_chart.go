package mujica

import (
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

// RealtimeBuffer keeps the latest samples of named series.
type RealtimeBuffer struct {
	names []string
	limit int
	times []time.Time
	vals  [][]float64
}

// NewRealtimeBuffer keeps up to capacity samples of the named series.
func NewRealtimeBuffer(capacity int, names ...string) *RealtimeBuffer {
	if capacity < 2 || len(names) == 0 {
		panic("mujica: RealtimeBuffer needs a capacity of 2 or more and a series name")
	}
	return &RealtimeBuffer{names: names, limit: capacity, vals: make([][]float64, len(names))}
}

// Append adds one sample per series at a non-decreasing t.
func (b *RealtimeBuffer) Append(t time.Time, values ...float64) {
	if len(values) != len(b.names) || len(b.times) > 0 && t.Before(b.times[len(b.times)-1]) {
		panic("mujica: RealtimeBuffer.Append needs one value per series and a time not before the last")
	}
	b.times = append(b.times, t)
	for i, v := range values {
		b.vals[i] = append(b.vals[i], v)
	}
	// Old samples are dropped in blocks, so appends stay cheap.
	if len(b.times) >= 2*b.limit {
		drop := len(b.times) - b.limit
		b.times = append([]time.Time(nil), b.times[drop:]...)
		for i := range b.vals {
			b.vals[i] = append([]float64(nil), b.vals[i][drop:]...)
		}
	}
}

// Len is how many samples the buffer keeps.
func (b *RealtimeBuffer) Len() int { return min(len(b.times), b.limit) }

// RealtimeChartOptions configure a RealtimeChart.
type RealtimeChartOptions struct {
	Frame ChartFrame
	// Window is the time shown up to the newest sample.
	Window time.Duration
	Format func(v float64) string
}

// RealtimeChartState holds whether the view is paused, and where.
type RealtimeChartState struct {
	Paused bool
	end    time.Time
}

// RealtimeChart shows a fixed window ending at the newest sample.
func RealtimeChart(c *ui.Context, s *RealtimeChartState, buf *RealtimeBuffer, opts RealtimeChartOptions) ChartResult {
	if opts.Window <= 0 {
		panic("mujica: RealtimeChart needs a positive Window")
	}
	u := Density(c).Unit()
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		label := msg(c, "realtime_chart.pause")
		if s.Paused {
			label = msg(c, "realtime_chart.resume")
		}
		if b := Button(c, label, ButtonOptions{Variant: Secondary, Disabled: wrap.IsDisabled()}).AlignSelf(ui.Start); b.Clicked() && !b.IsDisabled() {
			s.Paused, res.value = !s.Paused, true
			slog.Debug("mujica: realtime_chart pause", "paused", s.Paused)
		}
		times := buf.times[max(0, len(buf.times)-buf.limit):]
		last := time.Time{}
		if len(times) > 0 {
			last = times[len(times)-1]
		}
		if !s.Paused || s.end.IsZero() {
			s.end = last
		}
		from := s.end.Add(-opts.Window)
		i0 := sort.Search(len(times), func(i int) bool { return !times[i].Before(from) })
		i1 := sort.Search(len(times), func(i int) bool { return times[i].After(s.end) })
		off := len(buf.times) - len(times)
		ch := newChart(c, opts.Frame)
		series := make([]ChartSeries, len(buf.names))
		tops := make([][]float64, len(buf.names))
		for i, name := range buf.names {
			series[i] = ChartSeries{Name: name, Values: buf.vals[i][off+i0 : off+i1]}
			tops[i] = series[i].Values
		}
		n := i1 - i0
		var ax chartX
		if n > 0 {
			ax = chartMemo(ch, "x", [2]any{chartID(times[i0:i1]), s.end}, func() chartX {
				a := newChartX(n, nil, times[i0:i1])
				a.d0, a.d1 = float64(from.UnixNano())/1e9, float64(s.end.UnixNano())/1e9
				if a.d1 == a.d0 {
					a.d0 -= 1
				}
				return a
			})
		}
		lo, hi, ok := chartYRange(ch, series, false, nil)
		if !ok || math.IsNaN(lo) {
			n = 0
		}
		r := chartLines(ch, chartLineData{series: series, tops: tops, ax: ax, lo: lo, hi: hi, n: n, format: opts.Format})
		res.Point, res.changed, res.retry, res.toggled = r.Point, r.changed, r.retry, r.toggled
	})
	res.Element = wrap
	return res
}
