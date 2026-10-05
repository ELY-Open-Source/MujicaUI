package mujica

import (
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// EventEditorOptions configure an EventEditor.
type EventEditorOptions struct {
	// Location is the zone times are edited in; required.
	Location  *time.Location
	Hour12    bool
	WeekStart WeekStart
	Disabled  bool
	Label     string
}

// EventEditorResult is the editor and its save or cancel.
type EventEditorResult struct {
	Element  *ui.Element
	saved    bool
	canceled bool
}

// Saved reports a valid draft written to the event.
func (r EventEditorResult) Saved() bool { return r.saved }

// Canceled reports that the user discarded the draft.
func (r EventEditorResult) Canceled() bool { return r.canceled }

// eventEditorState is the draft and whether Save was tried.
type eventEditorState struct {
	draft  CalendarEvent
	tried  bool
	origin string
}

// EventEditor edits a copy of *ev; Save checks and writes.
func EventEditor(c *ui.Context, ev *CalendarEvent, opts EventEditorOptions) EventEditorResult {
	loc := eventLocation("EventEditor", opts.Location)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "event_editor.label")
	}
	r := EventEditorResult{}
	root := ui.Column(c).Gap(d.Unit() * 4).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	st := ui.Local(root, "event-editor", func() *eventEditorState { return nil })
	if *st == nil || (*st).origin != ev.ID {
		*st = &eventEditorState{draft: eventCopy(*ev), origin: ev.ID}
	}
	s := *st
	dr := &s.draft
	problems := eventEditorProblems(c, *dr, s.tried)
	root.Children(func() {
		FormField(c, msg(c, "event_editor.title"), FormFieldOptions{Required: true, Error: problems["title"]}, func() *ui.Element {
			return TextInput(c, &dr.Title, TextInputOptions{Label: msg(c, "event_editor.title"), Disabled: opts.Disabled, Error: problems["title"]}).Input
		})
		Checkbox(c, &dr.AllDay, msg(c, "event_editor.all_day"), CheckboxOptions{Disabled: opts.Disabled})
		ui.Row(c).Gap(d.Unit() * 3).Wrap().AlignItems(ui.Start).Children(func() {
			for _, f := range []struct {
				key string
				t   *time.Time
			}{{"event_editor.starts", &dr.Start}, {"event_editor.ends", &dr.End}} {
				FormField(c, msg(c, f.key), FormFieldOptions{Error: problems[f.key]}, func() *ui.Element {
					if dr.AllDay {
						day := DateOf(f.t.In(loc))
						r := DatePicker(c, &day, DatePickerOptions{Label: msg(c, f.key), WeekStart: opts.WeekStart, Disabled: opts.Disabled})
						if r.Changed() {
							*f.t = eventDayStart(day, loc)
						}
						return r.Input
					}
					return DateTimePicker(c, f.t, DateTimePickerOptions{Location: loc, Hour12: opts.Hour12, MinuteStep: 5,
						WeekStart: opts.WeekStart, Label: msg(c, f.key), Disabled: opts.Disabled}).Element
				})
			}
		})
		FormField(c, msg(c, "event_editor.where"), FormFieldOptions{}, func() *ui.Element {
			return TextInput(c, &dr.Location, TextInputOptions{Label: msg(c, "event_editor.where"), Disabled: opts.Disabled}).Input
		})
		FormField(c, msg(c, "event_editor.guests"), FormFieldOptions{Description: msg(c, "event_editor.guests_help")}, func() *ui.Element {
			names := make([]string, len(dr.Attendees))
			for i, a := range dr.Attendees {
				names[i] = a.Name
			}
			p := TagsInput(c, &names, TagsInputOptions{Label: msg(c, "event_editor.guests"), Disabled: opts.Disabled})
			if p.Changed() {
				dr.Attendees = eventAttendees(dr.Attendees, names)
			}
			return p.Input
		})
		if len(dr.Attendees) > 0 {
			AttendeeList(c, dr.Attendees, AttendeeListOptions{Disabled: opts.Disabled})
		}
		FormField(c, msg(c, "event_editor.reminders"), FormFieldOptions{}, func() *ui.Element {
			return ReminderPicker(c, &dr.Reminders, ReminderPickerOptions{Disabled: opts.Disabled}).Element
		})
		var rec RecurrenceResult
		FormField(c, msg(c, "recurrence_editor.repeat"), FormFieldOptions{}, func() *ui.Element {
			rec = RecurrenceEditor(c, &dr.Repeat, RecurrenceEditorOptions{Start: DateOf(dr.Start.In(loc)), WeekStart: opts.WeekStart, Disabled: opts.Disabled})
			return rec.Element
		})
		datetimeFooter(c, func() {
			if Button(c, msg(c, "calendar.cancel"), ButtonOptions{Variant: Secondary}).Clicked() && !root.IsDisabled() {
				slog.Debug("mujica: event_editor cancel", "id", ev.ID)
				*st, r.canceled = nil, true
			}
			if Button(c, msg(c, "event_editor.save"), ButtonOptions{}).Clicked() && !root.IsDisabled() {
				s.tried = true
				if len(eventEditorProblems(c, *dr, true)) == 0 && rec.Problem() == "" {
					slog.Debug("mujica: event_editor save", "id", ev.ID)
					out := eventCopy(*dr)
					if out.AllDay {
						first, last := DateOf(out.Start.In(loc)), DateOf(out.End.Add(-time.Nanosecond).In(loc))
						out.Start, out.End = eventDayStart(first, loc), eventDayStart(last.AddDays(1), loc)
					}
					*ev, r.saved = out, true
				} else {
					slog.Debug("mujica: event_editor invalid", "id", ev.ID)
				}
			}
		})
	})
	return r
}

// eventEditorProblems checks a draft once Save was tried.
func eventEditorProblems(c *ui.Context, ev CalendarEvent, tried bool) map[string]string {
	p := map[string]string{}
	if !tried {
		return p
	}
	if strings.TrimSpace(ev.Title) == "" {
		p["title"] = msg(c, "event_editor.need_title")
	}
	if ev.Start.IsZero() {
		p["event_editor.starts"] = msg(c, "event_editor.need_start")
	}
	if !ev.End.After(ev.Start) {
		p["event_editor.ends"] = msg(c, "event_editor.end_order")
	}
	return p
}

// eventAttendees keeps known attendees and adds new ones.
func eventAttendees(old []Attendee, names []string) []Attendee {
	out := make([]Attendee, 0, len(names))
	for _, n := range names {
		i := slices.IndexFunc(old, func(a Attendee) bool { return a.Name == n })
		if i >= 0 {
			out = append(out, old[i])
		} else {
			out = append(out, Attendee{Name: n})
		}
	}
	return out
}

// eventCopy clones an event's slices for a draft.
func eventCopy(ev CalendarEvent) CalendarEvent {
	ev.Attendees = slices.Clone(ev.Attendees)
	ev.Reminders = slices.Clone(ev.Reminders)
	ev.Repeat.Weekdays = slices.Clone(ev.Repeat.Weekdays)
	return ev
}
