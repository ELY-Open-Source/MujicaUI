package project

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// TimeEntry records a span; zero End means running.
type TimeEntry struct {
	Start, End time.Time
}

// TimeTrackerState holds the caller's tracked spans and manual correction.
type TimeTrackerState struct {
	Entries []TimeEntry
	// Adjust is added to the tracked spans by manual correction.
	Adjust time.Duration
}

// Running reports whether the last entry is still open.
func (s *TimeTrackerState) Running() bool {
	return len(s.Entries) > 0 && s.Entries[len(s.Entries)-1].End.IsZero()
}

// Total returns the tracked time at now, including the correction.
func (s *TimeTrackerState) Total(now time.Time) time.Duration {
	t := s.Adjust
	for _, e := range s.Entries {
		end := e.End
		if end.IsZero() {
			end = now
		}
		t += end.Sub(e.Start)
	}
	return t
}

// TimeTrackerOptions configure a TimeTracker.
type TimeTrackerOptions struct {
	// Now fixes the clock; zero uses the system clock.
	Now      time.Time
	Disabled bool
	Label    string
}

// TimeTracker records spans and supports manual total corrections.
func TimeTracker(c *ui.Context, s *TimeTrackerState, opts TimeTrackerOptions) core.ChoiceResult {
	now := opts.Now
	live := now.IsZero()
	if live {
		now = c.Now()
	}
	for i, e := range s.Entries {
		if e.Start.IsZero() || e.Start.After(now) || !e.End.IsZero() && (e.End.Before(e.Start) || e.End.After(now)) || e.End.IsZero() && i != len(s.Entries)-1 || i > 0 && e.Start.Before(s.Entries[i-1].End) {
			panic("mujica: TimeTracker entries must be ordered spans with only the last open")
		}
	}
	if s.Total(now) < 0 {
		panic("mujica: TimeTracker total must be non-negative")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "time_tracker.label")
	}
	changed := false
	row := ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	editing := ui.Local(row, "time-tracker-edit", func() bool { return false })
	draft := ui.Local(row, "time-tracker-draft", func() time.Duration { return 0 })
	act := func(what string, fn func()) {
		if row.IsDisabled() {
			return
		}
		fn()
		changed = true
		slog.Debug("mujica: time_tracker "+what, "total", s.Total(now))
	}
	row.Children(func() {
		total := s.Total(now)
		running := s.Running()
		icon, key, variant := "play", "time_tracker.start", input.Secondary
		if running {
			icon, key, variant = "pause", "time_tracker.stop", input.Primary
		}
		if input.Button(c, core.Msg(c, key), input.ButtonOptions{Variant: variant, Icon: icons.Must(icon)}).Clicked() {
			act(key[len("time_tracker."):], func() {
				if running {
					s.Entries[len(s.Entries)-1].End = now
				} else {
					s.Entries = append(s.Entries, TimeEntry{Start: now})
				}
			})
		}
		if *editing && !row.IsDisabled() {
			datetime.DurationPicker(c, draft, datetime.DurationPickerOptions{Seconds: true, Label: core.Msg(c, "time_tracker.correct")})
			if input.Button(c, core.Msg(c, "time_tracker.apply"), input.ButtonOptions{}).Clicked() {
				act("correct", func() { s.Adjust += *draft - total })
				*editing = false
			}
			if input.Button(c, core.Msg(c, "calendar.cancel"), input.ButtonOptions{Variant: input.Ghost}).Clicked() {
				*editing = false
			}
			return
		}
		*editing = false
		fg := k.Text
		if running {
			fg = k.AccentText
		}
		ui.Text(c, datetime.StopwatchText(total)).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(fg).FontFeatures("tnum").
			Label(label + ", " + datetime.DurationText(c, total.Truncate(time.Second)))
		if running {
			ui.Text(c, core.Msg(c, "time_tracker.running")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText)
		}
		ui.Box(c).Disabled(running).Children(func() {
			if b := core.IconAction(c, icons.Must("pencil"), core.Msg(c, "time_tracker.correct")); b.Clicked() && !b.IsDisabled() {
				*editing, *draft = true, total.Truncate(time.Second)
			}
		})
	})
	if s.Running() && live {
		c.After(100*time.Millisecond - s.Total(now)%(100*time.Millisecond))
	}
	return core.NewChoiceResult(row, changed)
}
