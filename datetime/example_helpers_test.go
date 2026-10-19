package datetime_test

import (
	"time"

	"github.com/ELY-Open-Source/MujicaUI/datetime"
)

func exampleEvents(loc *time.Location) []datetime.CalendarEvent {
	start := time.Date(2026, time.October, 14, 9, 0, 0, 0, loc)
	return []datetime.CalendarEvent{
		{ID: "rehearsal", Title: "Rehearsal", Start: start, End: start.Add(2 * time.Hour), Location: "Hall"},
		{ID: "fitting", Title: "Fitting", Start: start.Add(time.Hour), End: start.Add(90 * time.Minute), Color: datetime.EventInfo},
	}
}
