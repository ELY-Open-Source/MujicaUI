package mujica

import (
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// BlinkHighlightOptions configure a BlinkHighlight.
type BlinkHighlightOptions struct {
	// Format writes the value; nil uses two decimals.
	Format func(v float64) string
	// Duration of the fade; zero means 800 ms.
	Duration time.Duration
}

type blinkState struct {
	value float64
	dir   int
	at    time.Time
	set   bool
}

// BlinkHighlight shows value and tints it once when it changes, up in Success, down in Danger.
// An arrow marks the last direction, and stays as the only mark under reduced motion.
func BlinkHighlight(c *ui.Context, value float64, opts BlinkHighlightOptions) *ui.Element {
	if math.IsNaN(value) || math.IsInf(value, 0) || opts.Duration < 0 {
		panic("mujica: BlinkHighlight needs a finite value and a non-negative Duration")
	}
	d := opts.Duration
	if d == 0 {
		d = 800 * time.Millisecond
	}
	format := opts.Format
	if format == nil {
		format = func(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }
	}
	k := Tokens(c)
	u := Density(c).Unit()
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
	col := k.Success
	icon, word := "arrow-up", "blink_highlight.up"
	if st.dir < 0 {
		col, icon, word = k.Danger, "arrow-down", "blink_highlight.down"
	}
	name := format(value)
	if st.dir != 0 {
		name += ", " + msg(c, word)
	}
	row.Label(name)
	if dur := Motion(c, d); st.dir != 0 && dur > 0 {
		if p := float32(now.Sub(st.at)) / float32(dur); p < 1 {
			row.Background(col.Alpha(0.22 * (1 - p)))
			c.AnimationFrame()
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
