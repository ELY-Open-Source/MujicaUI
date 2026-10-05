package mujica

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Candle is one period of trading: open, high, low, close and volume.
type Candle struct {
	Time                   time.Time
	Open, High, Low, Close float64
	Volume                 float64
}

// marketCheckCandles panics on candles out of time order or with impossible prices.
func marketCheckCandles(candles []Candle) {
	for i, c := range candles {
		if c.Time.IsZero() || i > 0 && !c.Time.After(candles[i-1].Time) {
			panic(fmt.Sprintf("mujica: candle %d time is zero or not after the one before", i))
		}
		if !(c.Low <= min(c.Open, c.Close)) || !(c.High >= max(c.Open, c.Close)) || c.Volume < 0 || math.IsNaN(c.Volume) {
			panic(fmt.Sprintf("mujica: candle %d has impossible prices or volume", i))
		}
	}
}

// MarketView is the caller-held window of a market chart: zoom and pan.
type MarketView struct {
	// Count is how many candles show; zero shows 120.
	Count int
	// End is one past the last shown candle; zero follows the newest.
	End int
}

// marketWindow resolves the view to [from, to) over n candles, at most room wide.
func marketWindow(v MarketView, n, room int) (from, to int) {
	if v.Count < 0 || v.End < 0 {
		panic("mujica: MarketView Count and End must not be negative")
	}
	count := v.Count
	if count == 0 {
		count = 120
	}
	count = max(1, min(count, room, n))
	to = n
	if v.End > 0 && v.End < n {
		to = max(v.End, count)
	}
	return to - count, to
}

// marketPan moves the view by delta candles; reaching the newest follows it again.
func marketPan(v *MarketView, n, room, delta int) {
	from, to := marketWindow(*v, n, room)
	v.Count = to - from
	to = max(v.Count, min(n, to+delta))
	if to >= n {
		v.End = 0
	} else {
		v.End = to
	}
}

// marketZoom scales the shown count by f, keeping the newest shown candle, within 10 and room.
func marketZoom(v *MarketView, n, room int, f float64) {
	from, to := marketWindow(*v, n, room)
	count := int(math.Round(float64(to-from) * f))
	if f > 1 && count == to-from {
		count++
	}
	v.Count = max(min(10, n), min(count, room, n))
	if to < n {
		v.End = max(to, v.Count)
	}
}

// marketIndexAt finds the candle at or before t, or -1 before the first.
func marketIndexAt(candles []Candle, t time.Time) int {
	return sort.Search(len(candles), func(i int) bool { return candles[i].Time.After(t) }) - 1
}

// Interval names a candle period: 1m 3m 5m 15m 30m 1h 2h 4h 6h 12h 1D 3D 1W 1M.
type Interval string

// Intervals lists every Interval, shortest first.
var Intervals = []Interval{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "12h", "1D", "3D", "1W", "1M"}

// Close returns when the candle holding t closes: periods align to the clock in loc, weeks start Monday.
func (iv Interval) Close(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	step := func(d time.Duration) time.Time {
		since := t.Sub(day)
		return day.Add((since/d + 1) * d)
	}
	switch iv {
	case "1m", "3m", "5m", "15m", "30m":
		m, _ := time.ParseDuration(string(iv))
		return step(m)
	case "1h", "2h", "4h", "6h", "12h":
		h, _ := time.ParseDuration(string(iv))
		return step(h)
	case "1D":
		return day.AddDate(0, 0, 1)
	case "3D":
		days := int(day.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, loc)).Hours()/24 + 0.5)
		return day.AddDate(0, 0, 3-days%3)
	case "1W":
		back := (int(day.Weekday()) + 6) % 7
		return day.AddDate(0, 0, 7-back)
	case "1M":
		return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
	}
	panic(fmt.Sprintf("mujica: unknown Interval %q", string(iv)))
}
