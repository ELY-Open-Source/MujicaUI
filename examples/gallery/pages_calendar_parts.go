package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/datetime"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
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
			input.Select(c, &s.color, calOptions(calColors), input.SelectOptions{Label: "Color"})
			input.Checkbox(c, &s.inline, "Inline", input.CheckboxOptions{})
			input.Checkbox(c, &s.selected, "Selected", input.CheckboxOptions{})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		ev := calSample()[0]
		for i, n := range calColors {
			if n == s.color {
				ev.Color = datetime.EventColor(i)
			}
		}
		ui.Box(c).Width(240).Children(func() {
			if datetime.EventChip(c, ev, datetime.EventChipOptions{Location: time.Local, Inline: s.inline, Selected: s.selected, Disabled: s.disabled}).Clicked() {
				log("clicked " + ev.Title)
			}
		})
	})
}

func calEditorDemo(c *ui.Context, log func(string)) {
	type state struct {
		ev       datetime.CalendarEvent
		disabled bool
	}
	s := ui.Local(c.Root(), "427", func() state { return state{ev: calSample()[0]} })
	dtSection(c, "EventEditor", func() {
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		r := datetime.EventEditor(c, &s.ev, datetime.EventEditorOptions{Location: time.Local, Disabled: s.disabled})
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
			input.Checkbox(c, &s.readOnly, "Read-only", input.CheckboxOptions{})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		ev := calSample()[0]
		ev.Attendees = []datetime.Attendee{{Name: "Sakiko", Status: datetime.AttendeeAccepted, Organizer: true}, {Name: "Mutsumi"}}
		ev.Reminders = []time.Duration{15 * time.Minute}
		ui.Box(c).Width(240).Children(func() {
			chip := datetime.EventChip(c, ev, datetime.EventChipOptions{Location: time.Local, Disabled: s.disabled})
			if chip.Clicked() {
				s.open = true
				log("opened popover")
			}
			r := datetime.EventPopover(c, chip, &s.open, ev, datetime.EventPopoverOptions{Location: time.Local, ReadOnly: s.readOnly, Disabled: s.disabled})
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
		k := core.Tokens(c)
		ui.Box(c).Size(320, 240).Border(1, k.Border).Children(func() {
			datetime.CurrentTimeIndicator(c, datetime.CurrentTimeIndicatorOptions{Location: time.Local})
		})
	})
}

func calRecurrenceDemo(c *ui.Context, log func(string)) {
	type state struct {
		rule     datetime.Recurrence
		disabled bool
	}
	s := ui.Local(c.Root(), "430", func() state { return state{} })
	dtSection(c, "RecurrenceEditor", func() {
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		r := datetime.RecurrenceEditor(c, &s.rule, datetime.RecurrenceEditorOptions{Start: core.DateOf(time.Now()), Disabled: s.disabled})
		if r.Changed() {
			log(datetime.RecurrenceSummary(c, s.rule) + strings.TrimSpace(" "+r.Problem()))
		}
	})
}

func calAvailabilityDemo(c *ui.Context, log func(string)) {
	type state struct {
		mine, theirs   [][]bool
		days           []core.Date
		others, frozen bool
	}
	s := ui.Local(c.Root(), "431", func() state {
		today := core.DateOf(time.Now())
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
		days := []core.Date{today, today.AddDays(1), today.AddDays(2), today.AddDays(3), today.AddDays(4)}
		return state{mine: grid(2), theirs: grid(1), days: days, others: true}
	})
	dtSection(c, "AvailabilityPicker", func() {
		dtControls(c, func() {
			input.Checkbox(c, &s.others, "Show Mutsumi", input.CheckboxOptions{})
			input.Checkbox(c, &s.frozen, "Disabled", input.CheckboxOptions{})
		})
		var others []datetime.AvailabilityPerson
		if s.others {
			others = []datetime.AvailabilityPerson{{Name: "Mutsumi", Free: s.theirs}}
		}
		if datetime.AvailabilityPicker(c, &s.mine, datetime.AvailabilityPickerOptions{Days: s.days, Others: others, Disabled: s.frozen}).Changed() {
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
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		if datetime.TimezoneSelect(c, &s.zone, datetime.TimezoneSelectOptions{Disabled: s.disabled}).Changed() {
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
			input.Select(c, &s.max, calOptions([]string{"0", "2", "3"}), input.SelectOptions{Label: "Max (0 = any)"})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		if limit := calAtoi(s.max); limit > 0 && len(s.v) > limit {
			s.v = s.v[:limit]
		}
		if datetime.ReminderPicker(c, &s.v, datetime.ReminderPickerOptions{Max: calAtoi(s.max), Disabled: s.disabled}).Changed() {
			log(fmt.Sprint("reminders ", s.v))
		}
	})
}

func calAttendeeDemo(c *ui.Context, log func(string)) {
	type state struct {
		list      []datetime.Attendee
		removable bool
	}
	people := func() []datetime.Attendee {
		return []datetime.Attendee{{Name: "Sakiko Togawa", Status: datetime.AttendeeAccepted, Organizer: true},
			{Name: "Mutsumi Wakaba", Status: datetime.AttendeeTentative}, {Name: "Uika Misumi", Status: datetime.AttendeeDeclined},
			{Name: "Nyamu Yuutenji", Email: "nyamu@example.com"}}
	}
	s := ui.Local(c.Root(), "434", func() state { return state{list: people(), removable: true} })
	dtSection(c, "AttendeeList", func() {
		dtControls(c, func() {
			input.Checkbox(c, &s.removable, "Removable", input.CheckboxOptions{})
			if input.Button(c, "Reset", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				s.list = people()
				log("reset")
			}
		})
		if i, ok := datetime.AttendeeList(c, s.list, datetime.AttendeeListOptions{Removable: s.removable}).Removed(); ok {
			log("removed " + s.list[i].Name)
			s.list = append(s.list[:i:i], s.list[i+1:]...)
		}
	})
}
