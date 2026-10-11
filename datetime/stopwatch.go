package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// StopwatchState is the caller's run and laps.
type StopwatchState struct {
	running bool
	since   time.Time
	base    time.Duration
	// Laps are the elapsed times at each lap, oldest first.
	Laps []time.Duration
}

// Elapsed returns the time measured at now.
func (s *StopwatchState) Elapsed(now time.Time) time.Duration {
	if s.running {
		return s.base + now.Sub(s.since)
	}
	return s.base
}

// Running reports whether the stopwatch is counting.
func (s *StopwatchState) Running() bool { return s.running }

// StopwatchOptions configure a Stopwatch.
type StopwatchOptions struct {
	// Now fixes the clock and stops ticking.
	Now      time.Time
	Disabled bool
	Label    string
}

// Stopwatch measures time with laps, ticking while running and shown.
func Stopwatch(c *ui.Context, s *StopwatchState, opts StopwatchOptions) core.ChoiceResult {
	k := core.Tokens(c)
	d := core.Density(c)
	now := opts.Now
	live := now.IsZero()
	if live {
		now = c.Now()
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "stopwatch.label")
	}
	changed := false
	col := ui.Column(c).Gap(d.Unit() * 3).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	act := func(what string, fn func()) {
		if col.IsDisabled() {
			return
		}
		fn()
		changed = true
		slog.Debug("mujica: stopwatch "+what, "elapsed", s.Elapsed(now))
	}
	col.Children(func() {
		el := s.Elapsed(now)
		ui.Text(c, StopwatchText(el)).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.H1Size)).TextColor(k.Text).FontFeatures("tnum").
			Label(label + ", " + StopwatchText(el))
		ui.Row(c).Gap(d.Unit() * 2).Children(func() {
			start, variant := "stopwatch.start", input.Primary
			if s.running {
				start, variant = "stopwatch.pause", input.Secondary
			} else if el > 0 {
				start = "stopwatch.resume"
			}
			if input.Button(c, core.Msg(c, start), input.ButtonOptions{Variant: variant}).Clicked() {
				act("toggle", func() {
					if s.running {
						s.base, s.running = s.Elapsed(now), false
					} else {
						s.since, s.running = now, true
					}
				})
			}
			if input.Button(c, core.Msg(c, "stopwatch.lap"), input.ButtonOptions{Variant: input.Secondary, Disabled: !s.running}).Clicked() {
				act("lap", func() { s.Laps = append(s.Laps, s.Elapsed(now)) })
			}
			if input.Button(c, core.Msg(c, "stopwatch.reset"), input.ButtonOptions{Variant: input.Ghost, Disabled: s.running || el == 0 && len(s.Laps) == 0}).Clicked() {
				act("reset", func() { *s = StopwatchState{} })
			}
		})
		if len(s.Laps) > 0 {
			ui.Column(c).Role(ui.RoleList).Children(func() {
				for i := len(s.Laps) - 1; i >= 0; i-- {
					split := s.Laps[i]
					if i > 0 {
						split -= s.Laps[i-1]
					}
					row := ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Center).Gap(d.Unit() * 4).Role(ui.RoleListItem)
					row.DrawOver(func(p *ui.Painter, r ui.Rect) { internal.Rule(p, r.X, r.X+r.W, r.Y+r.H-1, k.Border) })
					row.Children(func() {
						ui.Text(c, fmt.Sprintf(core.Msg(c, "stopwatch.lap_n"), i+1)).Width(d.ControlHeight() * 2).TextColor(k.TextMuted)
						ui.Text(c, StopwatchText(split)).Font(theme.MonoFont).FontFeatures("tnum").TextColor(k.Text).Grow(1)
						ui.Text(c, StopwatchText(s.Laps[i])).Font(theme.MonoFont).FontFeatures("tnum").TextColor(k.TextMuted)
					})
				}
			})
		}
	})
	if s.running && live {
		c.After(100*time.Millisecond - s.Elapsed(now)%(100*time.Millisecond))
	}
	return core.NewChoiceResult(col, changed)
}

// StopwatchText formats elapsed time as 01:02.3, with hours when needed.
func StopwatchText(d time.Duration) string {
	t := int(d / (100 * time.Millisecond))
	h, m, s, tenth := t/36000, t%36000/600, t%600/10, t%10
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d.%d", h, m, s, tenth)
	}
	return fmt.Sprintf("%02d:%02d.%d", m, s, tenth)
}
