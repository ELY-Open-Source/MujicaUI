package feedback

import (
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CountUpOptions configure a CountUp.
type CountUpOptions struct {
	// Duration of a change; zero means StateDuration.
	Duration time.Duration
	// Decimals shown; Format overrides it.
	Decimals int
	// Format writes the shown value; nil uses Decimals.
	Format func(v float64) string
}

type countUpState struct {
	from, to float64
	start    time.Time
}

// CountUp counts from the last value to the new one.
func CountUp(c *ui.Context, value float64, opts CountUpOptions) *ui.Element {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic("mujica: CountUp value must be finite")
	}
	if opts.Duration < 0 || opts.Decimals < 0 {
		panic("mujica: CountUp Duration and Decimals must not be negative")
	}
	d := opts.Duration
	if d == 0 {
		d = theme.StateDuration
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', opts.Decimals, 64) }
	}
	dur := core.Motion(c, d)
	now := c.Now()
	row := ui.Row(c).Shrink(0).Label(format(value))
	st := ui.Local(row, "count", func() countUpState { return countUpState{from: value, to: value} })
	at := func() float64 {
		if dur == 0 {
			return st.to
		}
		p := min(1, float64(now.Sub(st.start))/float64(dur))
		eased := float64(ui.EaseOut(float32(p)))
		return st.from*(1-eased) + st.to*eased
	}
	if value != st.to {
		from := at()
		slog.Debug("mujica: count up", "from", from, "to", value)
		st.from, st.to, st.start = from, value, now
	}
	if dur == 0 {
		st.from, st.start = st.to, now
	}
	shown := at()
	if shown != st.to {
		row.DrawOver(func(*ui.Painter, ui.Rect) { c.AnimationFrame() })
	}
	return row.Children(func() { ui.Text(c, format(shown)).FontFeatures("tnum").SingleLine() })
}
