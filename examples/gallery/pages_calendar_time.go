package main

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/datetime"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

func calWeekPickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		week      core.Date
		weekStart string
		disabled  bool
	}
	s := ui.Local(c.Root(), "435", func() state { return state{weekStart: "Locale"} })
	dtSection(c, "WeekPicker", func() {
		dtControls(c, func() {
			if input.Select(c, &s.weekStart, calOptions(calWeekStarts), input.SelectOptions{Label: "Week start"}).Changed() {
				s.week = core.Date{}
				log("week start " + s.weekStart + ", value cleared")
			}
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		opts := datetime.WeekPickerOptions{WeekStart: calWeekStart(s.weekStart), Clearable: true, Disabled: s.disabled}
		if datetime.WeekPicker(c, &s.week, opts).Changed() {
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
			input.Checkbox(c, &s.seconds, "Seconds", input.CheckboxOptions{})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		opts := datetime.DurationPickerOptions{Seconds: s.seconds, Min: 15 * time.Minute, Max: 8 * time.Hour, Disabled: s.disabled}
		if datetime.DurationPicker(c, &s.v, opts).Changed() {
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
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		if datetime.CronEditor(c, &s.expr, datetime.CronEditorOptions{Location: time.Local, Disabled: s.disabled}).Changed() {
			if _, err := datetime.ParseCron(s.expr); err != nil {
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
				if input.Button(c, b.label, input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					s.at = time.Now().Add(b.d)
					log("time set: " + b.label)
				}
			}
		})
		datetime.RelativeTime(c, s.at, datetime.RelativeTimeOptions{Location: time.Local})
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
			if input.Button(c, "Restart 10 s", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				s.deadline = time.Now().Add(10 * time.Second)
				log("restarted")
			}
			input.Checkbox(c, &s.words, "Words", input.CheckboxOptions{})
		})
		format := datetime.CountdownClock
		if s.words {
			format = datetime.CountdownWords
		}
		if datetime.Countdown(c, s.deadline, datetime.CountdownOptions{Format: format}).Finished() {
			log("finished")
		}
	})
}

func calStopwatchDemo(c *ui.Context, log func(string)) {
	type state struct {
		sw       datetime.StopwatchState
		disabled bool
	}
	s := ui.Local(c.Root(), "440", func() state { return state{} })
	dtSection(c, "Stopwatch", func() {
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		if datetime.Stopwatch(c, &s.sw, datetime.StopwatchOptions{Disabled: s.disabled}).Changed() {
			state := "paused"
			if s.sw.Running() {
				state = "running"
			}
			log(state + ", " + datetime.StopwatchText(s.sw.Elapsed(time.Now())))
		}
	})
}
