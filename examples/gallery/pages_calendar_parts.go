package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

var calColors = []string{"Accent", "Info", "Success", "Warning", "Danger", "Neutral"}

func calChipDemo(c *ui.Context, log func(string)) {
	type state struct {
		color                      string
		inline, selected, disabled bool
	}
	s := ui.Local(c.Root(), "426", func() state { return state{color: "Accent"} })
	dtSection(c, "EventChip", func() {
		dtControls(c, func() {
			mujica.Select(c, &s.color, calOptions(calColors), mujica.SelectOptions{Label: "Color"})
			mujica.Checkbox(c, &s.inline, "Inline", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.selected, "Selected", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		ev := calSample()[0]
		for i, n := range calColors {
			if n == s.color {
				ev.Color = mujica.EventColor(i)
			}
		}
		ui.Box(c).Width(240).Children(func() {
			if mujica.EventChip(c, ev, mujica.EventChipOptions{Location: time.Local, Inline: s.inline, Selected: s.selected, Disabled: s.disabled}).Clicked() {
				log("clicked " + ev.Title)
			}
		})
	})
}

func calEditorDemo(c *ui.Context, log func(string)) {
	type state struct {
		ev       mujica.CalendarEvent
		disabled bool
	}
	s := ui.Local(c.Root(), "427", func() state { return state{ev: calSample()[0]} })
	dtSection(c, "EventEditor", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		r := mujica.EventEditor(c, &s.ev, mujica.EventEditorOptions{Location: time.Local, Disabled: s.disabled})
		if r.Saved() {
			log(fmt.Sprintf("saved %q %s, %d guests", s.ev.Title, s.ev.Start.Format("Jan 2 15:04"), len(s.ev.Attendees)))
		}
		if r.Canceled() {
			log("canceled")
		}
	})
}

func calPopoverDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, readOnly, disabled bool
	}
	s := ui.Local(c.Root(), "428", func() state { return state{} })
	dtSection(c, "EventPopover", func() {
		dtControls(c, func() {
			mujica.Checkbox(c, &s.readOnly, "Read-only", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		ev := calSample()[0]
		ev.Attendees = []mujica.Attendee{{Name: "Sakiko", Status: mujica.AttendeeAccepted, Organizer: true}, {Name: "Mutsumi"}}
		ev.Reminders = []time.Duration{15 * time.Minute}
		ui.Box(c).Width(240).Children(func() {
			chip := mujica.EventChip(c, ev, mujica.EventChipOptions{Location: time.Local, Disabled: s.disabled})
			if chip.Clicked() {
				s.open = true
				log("opened popover")
			}
			r := mujica.EventPopover(c, chip, &s.open, ev, mujica.EventPopoverOptions{Location: time.Local, ReadOnly: s.readOnly, Disabled: s.disabled})
			if r.Edit() {
				log("edit " + ev.ID)
			}
			if r.Delete() {
				log("delete " + ev.ID)
			}
		})
	})
}

func calNowDemo(c *ui.Context, log func(string)) {
	dtSection(c, "CurrentTimeIndicator", func() {
		k := mujica.Tokens(c)
		ui.Box(c).Size(320, 240).Border(1, k.Border).Children(func() {
			mujica.CurrentTimeIndicator(c, mujica.CurrentTimeIndicatorOptions{Location: time.Local})
		})
	})
}

func calRecurrenceDemo(c *ui.Context, log func(string)) {
	type state struct {
		rule     mujica.Recurrence
		disabled bool
	}
	s := ui.Local(c.Root(), "430", func() state { return state{} })
	dtSection(c, "RecurrenceEditor", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		r := mujica.RecurrenceEditor(c, &s.rule, mujica.RecurrenceEditorOptions{Start: mujica.DateOf(time.Now()), Disabled: s.disabled})
		if r.Changed() {
			log(mujica.RecurrenceSummary(c, s.rule) + strings.TrimSpace(" "+r.Problem()))
		}
	})
}

func calAvailabilityDemo(c *ui.Context, log func(string)) {
	type state struct {
		mine, theirs   [][]bool
		days           []mujica.Date
		others, frozen bool
	}
	s := ui.Local(c.Root(), "431", func() state {
		today := mujica.DateOf(time.Now())
		grid := func(seed int) [][]bool {
			g := make([][]bool, 5)
			for i := range g {
				g[i] = make([]bool, 18)
				for j := range g[i] {
					g[i][j] = (i*7+j*3+seed)%4 == 0
				}
			}
			return g
		}
		days := []mujica.Date{today, today.AddDays(1), today.AddDays(2), today.AddDays(3), today.AddDays(4)}
		return state{mine: grid(2), theirs: grid(1), days: days, others: true}
	})
	dtSection(c, "AvailabilityPicker", func() {
		dtControls(c, func() {
			mujica.Checkbox(c, &s.others, "Show Mutsumi", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.frozen, "Disabled", mujica.CheckboxOptions{})
		})
		var others []mujica.AvailabilityPerson
		if s.others {
			others = []mujica.AvailabilityPerson{{Name: "Mutsumi", Free: s.theirs}}
		}
		if mujica.AvailabilityPicker(c, &s.mine, mujica.AvailabilityPickerOptions{Days: s.days, Others: others, Disabled: s.frozen}).Changed() {
			n := 0
			for _, d := range s.mine {
				for _, v := range d {
					if v {
						n++
					}
				}
			}
			log(fmt.Sprintf("%d free slots", n))
		}
	})
}

func calTimezoneDemo(c *ui.Context, log func(string)) {
	type state struct {
		zone     string
		disabled bool
	}
	s := ui.Local(c.Root(), "432", func() state { return state{zone: "Asia/Shanghai"} })
	dtSection(c, "TimezoneSelect", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		if mujica.TimezoneSelect(c, &s.zone, mujica.TimezoneSelectOptions{Disabled: s.disabled}).Changed() {
			log("zone " + s.zone)
		}
	})
}

func calReminderDemo(c *ui.Context, log func(string)) {
	type state struct {
		v        []time.Duration
		max      string
		disabled bool
	}
	s := ui.Local(c.Root(), "433", func() state { return state{v: []time.Duration{15 * time.Minute}, max: "0"} })
	dtSection(c, "ReminderPicker", func() {
		dtControls(c, func() {
			mujica.Select(c, &s.max, calOptions([]string{"0", "2", "3"}), mujica.SelectOptions{Label: "Max (0 = any)"})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		if limit := calAtoi(s.max); limit > 0 && len(s.v) > limit {
			s.v = s.v[:limit]
		}
		if mujica.ReminderPicker(c, &s.v, mujica.ReminderPickerOptions{Max: calAtoi(s.max), Disabled: s.disabled}).Changed() {
			log(fmt.Sprint("reminders ", s.v))
		}
	})
}

func calAttendeeDemo(c *ui.Context, log func(string)) {
	type state struct {
		list      []mujica.Attendee
		removable bool
	}
	people := func() []mujica.Attendee {
		return []mujica.Attendee{{Name: "Sakiko Togawa", Status: mujica.AttendeeAccepted, Organizer: true},
			{Name: "Mutsumi Wakaba", Status: mujica.AttendeeTentative}, {Name: "Uika Misumi", Status: mujica.AttendeeDeclined},
			{Name: "Nyamu Yuutenji", Email: "nyamu@example.com"}}
	}
	s := ui.Local(c.Root(), "434", func() state { return state{list: people(), removable: true} })
	dtSection(c, "AttendeeList", func() {
		dtControls(c, func() {
			mujica.Checkbox(c, &s.removable, "Removable", mujica.CheckboxOptions{})
			if mujica.Button(c, "Reset", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				s.list = people()
				log("reset")
			}
		})
		if i, ok := mujica.AttendeeList(c, s.list, mujica.AttendeeListOptions{Removable: s.removable}).Removed(); ok {
			log("removed " + s.list[i].Name)
			s.list = append(s.list[:i:i], s.list[i+1:]...)
		}
	})
}
