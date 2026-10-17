package datetime

import (
	"fmt"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// EventColor picks an event's color from the theme's event palette.
type EventColor int

const (
	EventAccent EventColor = iota
	EventInfo
	EventSuccess
	EventWarning
	EventDanger
	EventNeutral
)

// CalendarEvent is one entry of the views; End is exclusive.
type CalendarEvent struct {
	ID       string
	Title    string
	Start    time.Time
	End      time.Time
	AllDay   bool
	Color    EventColor
	Location string
	// Attendees, Reminders and Repeat are edited by EventEditor.
	Attendees []Attendee
	Reminders []time.Duration
	Repeat    Recurrence
}

// CalendarChangeKind says what a view gesture did.
type CalendarChangeKind int

const (
	CalendarMove CalendarChangeKind = iota
	CalendarResize
	CalendarCreate
)

// CalendarChange is a move, resize or creation to apply.
type CalendarChange struct {
	Kind       CalendarChangeKind
	ID         string
	Start, End time.Time
	AllDay     bool
}

// CalendarViewResult is a view and what the user asked.
type CalendarViewResult struct {
	Element ui.Element
	opened  string
	change  *CalendarChange
	day     core.Date
}

// Opened returns the event the user activated.
func (r CalendarViewResult) Opened() (string, bool) { return r.opened, r.opened != "" }

// Changed returns the move, resize or creation the user made.
func (r CalendarViewResult) Changed() (CalendarChange, bool) {
	if r.change == nil {
		return CalendarChange{}, false
	}
	return *r.change, true
}

// DayChosen returns the day the user activated.
func (r CalendarViewResult) DayChosen() (core.Date, bool) { return r.day, !r.day.IsZero() }

// eventPalette maps EventColor to tokens at 3:1 on Surface.
func eventPalette(k theme.Tokens) [6]ui.Color {
	return [6]ui.Color{k.AccentText, k.Info, k.Success, k.Warning, k.Danger, k.TextMuted}
}

func EventTone(k theme.Tokens, c EventColor) ui.Color {
	if c < EventAccent || c > EventNeutral {
		panic(fmt.Sprintf("mujica: EventColor %d out of range", c))
	}
	return eventPalette(k)[c]
}

// eventsCheck panics on events no view can place.
func eventsCheck(name string, events []CalendarEvent) {
	for _, e := range events {
		if e.ID == "" {
			panic("mujica: " + name + " event without ID")
		}
		if !e.End.After(e.Start) {
			panic(fmt.Sprintf("mujica: %s event %q ends before it starts", name, e.ID))
		}
	}
}

// eventsCache holds data derived from events overlapping a span.
type eventsCache[T any] struct {
	first    *CalendarEvent
	n        int
	version  int
	from, to time.Time
	value    T
	valid    bool
}

// eventsDerived rebuilds build(hits) only when events or span change.
func eventsDerived[T any](e ui.Element, key string, events []CalendarEvent, version int, from, to time.Time, build func(hits []int) T) T {
	cache := ui.Local(e, key, func() *eventsCache[T] { return &eventsCache[T]{} })
	var first *CalendarEvent
	if len(events) > 0 {
		first = &events[0]
	}
	if x := *cache; x.valid && x.first == first && x.n == len(events) && x.version == version && x.from.Equal(from) && x.to.Equal(to) {
		return x.value
	}
	hits := []int{}
	for i, ev := range events {
		if ev.Start.Before(to) && ev.End.After(from) {
			hits = append(hits, i)
		}
	}
	sort.SliceStable(hits, func(a, b int) bool {
		x, y := events[hits[a]], events[hits[b]]
		if !x.Start.Equal(y.Start) {
			return x.Start.Before(y.Start)
		}
		return x.End.After(y.End)
	})
	v := build(hits)
	**cache = eventsCache[T]{first, len(events), version, from, to, v, true}
	return v
}

// eventSpan is a placed interval, in any unit.
type eventSpan struct{ from, to float64 }

// eventColumns gives overlapping spans side-by-side columns and cluster widths.
func eventColumns(spans []eventSpan) (col, cols []int) {
	col, cols = make([]int, len(spans)), make([]int, len(spans))
	order := make([]int, len(spans))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool { return spans[order[a]].from < spans[order[b]].from })
	var ends []float64
	var cluster []int
	clusterEnd := -1.0
	flush := func() {
		for _, i := range cluster {
			cols[i] = len(ends)
		}
		cluster, ends = cluster[:0], ends[:0]
	}
	for _, i := range order {
		s := spans[i]
		if len(cluster) > 0 && s.from >= clusterEnd {
			flush()
		}
		placed := false
		for c, end := range ends {
			if end <= s.from {
				ends[c], col[i], placed = s.to, c, true
				break
			}
		}
		if !placed {
			col[i] = len(ends)
			ends = append(ends, s.to)
		}
		cluster = append(cluster, i)
		clusterEnd = max(clusterEnd, s.to)
	}
	flush()
	return col, cols
}

// eventDayStart is d's first instant in loc.
func eventDayStart(d core.Date, loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// eventTimeText formats an event's wall-clock span in loc.
func eventTimeText(c *ui.Context, ev CalendarEvent, loc *time.Location, hour12 bool) string {
	if ev.AllDay {
		return core.Msg(c, "event_chip.all_day")
	}
	cfg := timeConfig{hour12: hour12}
	s, e := ev.Start.In(loc), ev.End.In(loc)
	text := timeText(c, TimeOfDay{s.Hour(), s.Minute(), 0}, cfg) + " – " + timeText(c, TimeOfDay{e.Hour(), e.Minute(), 0}, cfg)
	return text
}

// eventShift moves ev by days, keeping wall time.
func eventShift(c *ui.Context, ev CalendarEvent, days int, loc *time.Location) (CalendarChange, string) {
	s, e := ev.Start.In(loc), ev.End.In(loc)
	ch := CalendarChange{Kind: CalendarMove, ID: ev.ID, AllDay: ev.AllDay}
	if ev.AllDay {
		ch.Start, ch.End = eventDayStart(core.DateOf(s).AddDays(days), loc), eventDayStart(core.DateOf(e).AddDays(days), loc)
		return ch, ""
	}
	day := core.DateOf(s).AddDays(days)
	wall := TimeOfDay{s.Hour(), s.Minute(), s.Second()}
	at, n := dateTimeResolve(day, wall, loc)
	if p := eventWallProblem(c, day, wall, loc, n); p != "" {
		return ch, p
	}
	ch.Start, ch.End = at, at.Add(ev.End.Sub(ev.Start))
	return ch, ""
}

// eventWallProblem reports a wall time DST skips or repeats.
func eventWallProblem(c *ui.Context, day core.Date, wall TimeOfDay, loc *time.Location, n int) string {
	key := ""
	switch {
	case n == 0:
		key = "date_time_picker.gap"
	case n > 1:
		key = "date_time_picker.ambiguous"
	default:
		return ""
	}
	return fmt.Sprintf(core.Msg(c, key), day.String()+" "+fmt.Sprintf("%02d:%02d", wall.Hour, wall.Minute), loc.String())
}

// eventLocation panics without a caller zone; views never assume UTC.
func eventLocation(name string, loc *time.Location) *time.Location {
	if loc == nil {
		panic("mujica: " + name + " needs a Location")
	}
	return loc
}
