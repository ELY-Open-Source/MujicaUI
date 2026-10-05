package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// CountdownFormat chooses how the time left reads.
type CountdownFormat int

const (
	// CountdownClock reads 01:02:03, with days as "2d" first.
	CountdownClock CountdownFormat = iota
	// CountdownWords reads "1 h 2 min 3 s".
	CountdownWords
)

// CountdownOptions configure a Countdown.
type CountdownOptions struct {
	Format CountdownFormat
	// Now fixes the clock and stops updates.
	Now   time.Time
	Label string
}

// CountdownResult is the countdown and its finish.
type CountdownResult struct {
	Element  *ui.Element
	finished bool
}

// Finished reports reaching zero, once per deadline.
func (r CountdownResult) Finished() bool { return r.finished }

// Countdown shows the time left, ticking each second while shown.
func Countdown(c *ui.Context, deadline time.Time, opts CountdownOptions) CountdownResult {
	if opts.Format < CountdownClock || opts.Format > CountdownWords {
		panic(fmt.Sprintf("mujica: Countdown Format %d out of range", opts.Format))
	}
	k := Tokens(c)
	now := opts.Now
	live := now.IsZero()
	if live {
		now = c.Now()
	}
	left := max(deadline.Sub(now), 0)
	shown := left.Truncate(time.Second)
	if left%time.Second != 0 {
		shown += time.Second
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "countdown.label")
	}
	text := CountdownText(c, shown, opts.Format)
	e := ui.Text(c, text).Font(theme.MonoFont).FontSize(FontSize(c, theme.H2Size)).TextColor(k.Text).FontFeatures("tnum").
		Label(label + ", " + CountdownText(c, shown, CountdownWords))
	r := CountdownResult{Element: e}
	fired := ui.Local(e, "countdown-fired", func() time.Time { return time.Time{} })
	switch {
	case left == 0 && !fired.Equal(deadline):
		*fired, r.finished = deadline, true
		slog.Debug("mujica: countdown finished", "deadline", deadline)
		c.Announce(fmt.Sprintf(msg(c, "countdown.done"), label))
	case left > 0 && live:
		c.After(left % time.Second)
		if left%time.Second == 0 {
			c.After(time.Second)
		}
	}
	if left == 0 {
		e.TextColor(k.AccentText)
	}
	return r
}

// CountdownText formats a non-negative duration in whole seconds.
func CountdownText(c *ui.Context, d time.Duration, f CountdownFormat) string {
	s := int(d / time.Second)
	days, h, m, sec := s/86400, s%86400/3600, s%3600/60, s%60
	if f == CountdownClock {
		t := fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
		if days > 0 {
			t = fmt.Sprintf(msg(c, "countdown.days_clock"), days) + " " + t
		}
		return t
	}
	out := ""
	for i, v := range []int{days, h, m, sec} {
		if v > 0 || i == 3 && out == "" {
			if out != "" {
				out += " "
			}
			out += fmt.Sprintf(msg(c, []string{"countdown.d", "countdown.h", "countdown.m", "countdown.s"}[i]), v)
		}
	}
	return out
}
