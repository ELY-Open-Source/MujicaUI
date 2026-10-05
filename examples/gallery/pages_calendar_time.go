package main

import (
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func calWeekPickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		week      mujica.Date
		weekStart string
		disabled  bool
	}
	s := ui.Local(c.Root(), "435", func() state { return state{weekStart: "Locale"} })
	dtSection(c, "WeekPicker", func() {
		dtControls(c, func() {
			if mujica.Select(c, &s.weekStart, calOptions(calWeekStarts), mujica.SelectOptions{Label: "Week start"}).Changed() {
				s.week = mujica.Date{}
				log("week start " + s.weekStart + ", value cleared")
			}
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		opts := mujica.WeekPickerOptions{WeekStart: calWeekStart(s.weekStart), Clearable: true, Disabled: s.disabled}
		if mujica.WeekPicker(c, &s.week, opts).Changed() {
			log("week of " + s.week.String())
		}
	})
}

func calDurationDemo(c *ui.Context, log func(string)) {
	type state struct {
		v                 time.Duration
		seconds, disabled bool
	}
	s := ui.Local(c.Root(), "436", func() state { return state{v: 90 * time.Minute} })
	dtSection(c, "DurationPicker", func() {
		dtControls(c, func() {
			mujica.Checkbox(c, &s.seconds, "Seconds", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		opts := mujica.DurationPickerOptions{Seconds: s.seconds, Min: 15 * time.Minute, Max: 8 * time.Hour, Disabled: s.disabled}
		if mujica.DurationPicker(c, &s.v, opts).Changed() {
			log("duration " + s.v.String())
		}
	})
}

func calCronDemo(c *ui.Context, log func(string)) {
	type state struct {
		expr     string
		disabled bool
	}
	s := ui.Local(c.Root(), "437", func() state { return state{expr: "*/15 9-17 * * MON-FRI"} })
	dtSection(c, "CronEditor", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		if mujica.CronEditor(c, &s.expr, mujica.CronEditorOptions{Location: time.Local, Disabled: s.disabled}).Changed() {
			if _, err := mujica.ParseCron(s.expr); err != nil {
				log("invalid: " + s.expr)
			} else {
				log("schedule " + s.expr)
			}
		}
	})
}

func calRelativeDemo(c *ui.Context, log func(string)) {
	type state struct{ at time.Time }
	s := ui.Local(c.Root(), "438", func() state { return state{at: time.Now()} })
	dtSection(c, "RelativeTime", func() {
		dtControls(c, func() {
			for _, b := range []struct {
				label string
				d     time.Duration
			}{{"Now", 0}, {"5 minutes ago", -5 * time.Minute}, {"In 3 hours", 3 * time.Hour}, {"40 days ago", -40 * 24 * time.Hour}} {
				if mujica.Button(c, b.label, mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
					s.at = time.Now().Add(b.d)
					log("time set: " + b.label)
				}
			}
		})
		mujica.RelativeTime(c, s.at, mujica.RelativeTimeOptions{Location: time.Local})
	})
}

func calCountdownDemo(c *ui.Context, log func(string)) {
	type state struct {
		deadline time.Time
		words    bool
	}
	s := ui.Local(c.Root(), "439", func() state { return state{deadline: time.Now().Add(10 * time.Second)} })
	dtSection(c, "Countdown", func() {
		dtControls(c, func() {
			if mujica.Button(c, "Restart 10 s", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				s.deadline = time.Now().Add(10 * time.Second)
				log("restarted")
			}
			mujica.Checkbox(c, &s.words, "Words", mujica.CheckboxOptions{})
		})
		format := mujica.CountdownClock
		if s.words {
			format = mujica.CountdownWords
		}
		if mujica.Countdown(c, s.deadline, mujica.CountdownOptions{Format: format}).Finished() {
			log("finished")
		}
	})
}

func calStopwatchDemo(c *ui.Context, log func(string)) {
	type state struct {
		sw       mujica.StopwatchState
		disabled bool
	}
	s := ui.Local(c.Root(), "440", func() state { return state{} })
	dtSection(c, "Stopwatch", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		if mujica.Stopwatch(c, &s.sw, mujica.StopwatchOptions{Disabled: s.disabled}).Changed() {
			state := "paused"
			if s.sw.Running() {
				state = "running"
			}
			log(state + ", " + mujica.StopwatchText(s.sw.Elapsed(time.Now())))
		}
	})
}
