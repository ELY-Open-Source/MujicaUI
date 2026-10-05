package mujica

import (
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
)

// CountUpOptions configure a CountUp.
type CountUpOptions struct {
	// Duration of a change; zero means 600 ms.
	Duration time.Duration
	// Decimals shown; Format overrides it.
	Decimals int
	// Format writes the shown value; nil uses Decimals.
	Format func(v float64) string
}

type countUpState struct {
	from, to float64
	start    time.Time
	set      bool
}

// CountUp counts from the last value to value; reduced motion shows it at once.
func CountUp(c *ui.Context, value float64, opts CountUpOptions) *ui.Element {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic("mujica: CountUp value must be finite")
	}
	if opts.Duration < 0 || opts.Decimals < 0 {
		panic("mujica: CountUp Duration and Decimals must not be negative")
	}
	d := opts.Duration
	if d == 0 {
		d = 600 * time.Millisecond
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', opts.Decimals, 64) }
	}
	dur := Motion(c, d)
	now := c.Now()
	row := ui.Row(c).Shrink(0).Label(format(value))
	st := ui.Local(row, "count", func() countUpState { return countUpState{from: value, to: value} })
	at := func() float64 {
		if dur == 0 {
			return st.to
		}
		p := min(1, float64(now.Sub(st.start))/float64(dur))
		return st.from + (st.to-st.from)*float64(ui.EaseOut(float32(p)))
	}
	if value != st.to {
		from := at()
		slog.Debug("mujica: count up", "from", from, "to", value)
		st.from, st.to, st.start = from, value, now
	}
	shown := at()
	if shown != st.to {
		c.AnimationFrame()
	}
	return row.Children(func() { ui.Text(c, format(shown)).FontFeatures("tnum").SingleLine() })
}
