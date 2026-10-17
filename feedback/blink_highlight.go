package feedback

import (
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// BlinkHighlightOptions configure a BlinkHighlight.
type BlinkHighlightOptions struct {
	// Format writes the value; nil uses two decimals.
	Format func(v float64) string
	// Duration of the fade; zero means StateDuration.
	Duration time.Duration
}

type blinkState struct {
	value float64
	dir   int
	at    time.Time
	set   bool
}

// BlinkHighlight tints value once on change and marks its direction.
func BlinkHighlight(c *ui.Context, value float64, opts BlinkHighlightOptions) ui.Element {
	if math.IsNaN(value) || math.IsInf(value, 0) || opts.Duration < 0 {
		panic("mujica: BlinkHighlight needs a finite value and a non-negative Duration")
	}
	d := opts.Duration
	if d == 0 {
		d = theme.StateDuration
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.Row(c).AlignItems(ui.Center).Gap(u).Padding(0, u*1.5).Radius(theme.ControlRadius).Shrink(0)
	st := ui.Local(row, "blink", func() blinkState { return blinkState{} })
	now := c.Now()
	if st.set && value != st.value {
		st.dir = 1
		if value < st.value {
			st.dir = -1
		}
		st.at = now
		slog.Debug("mujica: blink highlight", "from", st.value, "to", value)
	}
	st.value, st.set = value, true
	if c.Preferences().ReduceMotion {
		st.at = time.Time{}
	}
	col := k.Success
	icon, word := "arrow-up", "blink_highlight.up"
	if st.dir < 0 {
		col, icon, word = k.Danger, "arrow-down", "blink_highlight.down"
	}
	name := format(value)
	if st.dir != 0 {
		name += ", " + core.Msg(c, word)
	}
	row.Label(name)
	if dur := core.Motion(c, d); st.dir != 0 && !st.at.IsZero() && dur > 0 {
		if p := float32(now.Sub(st.at)) / float32(dur); p < 1 {
			row.Background(col.Alpha(0.22 * (1 - p)))
			row.DrawOver(func(*ui.Painter, ui.Rect) { c.AnimationFrame() })
		}
	}
	row.Children(func() {
		ui.Text(c, format(value)).FontFeatures("tnum").SingleLine()
		if st.dir != 0 {
			ui.Icon(c, icons.Must(icon)).Size(12, 12).TextColor(col)
		}
	})
	return row
}
