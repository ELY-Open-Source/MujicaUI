package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// AgendaViewOptions configure an AgendaView.
type AgendaViewOptions struct {
	// Location is the zone days are cut in; required.
	Location *time.Location
	// Version must change when events change in place.
	Version int
	// From is the first day listed; zero means today.
	From core.Date
	// Days is how many days are listed; zero means 30.
	Days int
	// Today is marked; zero uses the system clock.
	Today    core.Date
	Hour12   bool
	Disabled bool
	Label    string
}

// agendaRow is a day heading (ev < 0) or event.
type agendaRow struct {
	day core.Date
	ev  int
}

// AgendaView lists events by day, building only visible rows.
func AgendaView(c *ui.Context, events []CalendarEvent, opts AgendaViewOptions) CalendarViewResult {
	loc := eventLocation("AgendaView", opts.Location)
	if opts.Days < 0 {
		panic("mujica: AgendaView Days is negative")
	}
	days := opts.Days
	if days == 0 {
		days = 30
	}
	eventsCheck("AgendaView", events)
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	from := opts.From
	if from.IsZero() {
		from = today
	}
	core.DatetimeCheckDate("AgendaView From", from)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "agenda_view.label")
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Grow(1).MinHeight(0).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	rows := eventsDerived(root, "agenda-rows", events, opts.Version, eventDayStart(from, loc), eventDayStart(from.AddDays(days), loc),
		func(hits []int) []agendaRow { return agendaRows(events, hits, from, days, loc) })
	state := ui.Local(root, "agenda-list", func() *ui.ListState { return &ui.ListState{} })
	st := *state
	st.Key = func(i int) any { return rows[i] }
	st.Header = func(i int) bool { return rows[i].ev < 0 }
	st.Label = func(i int) string {
		if rows[i].ev < 0 {
			return core.DatetimeDayLabel(c, rows[i].day)
		}
		return events[rows[i].ev].Title
	}
	root.Children(func() {
		list := ui.List(c, st, len(rows), func(i int) {
			row := rows[i]
			if row.ev < 0 {
				h := ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Center).Gap(d.Unit() * 2).PaddingX(d.Unit() * 3).Background(k.Background).Role(ui.RoleColumnHeader)
				h.Children(func() {
					t := ui.Text(c, core.DatetimeDayLabel(c, row.day)).Font(theme.SerifFont).TextColor(k.Text)
					if row.day == today {
						t.Bold().TextColor(k.AccentText)
						ui.Text(c, core.Msg(c, "agenda_view.today")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText)
					}
				})
				h.DrawOver(func(p *ui.Painter, rr ui.Rect) { internal.Rule(p, rr.X, rr.X+rr.W, rr.Y+rr.H-1, k.Ornament) })
				return
			}
			if agendaItem(c, events[row.ev], loc, opts.Hour12) && !root.IsDisabled() {
				r.opened = events[row.ev].ID
			}
		}).Grow(1).MinHeight(d.ControlHeight() * 4).FocusGroup(ui.Vertical)
		list.Children(func() {
			if len(rows) == 0 {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "agenda_view.empty"), days)).TextColor(k.TextMuted).Padding(d.Unit() * 3)
			}
		})
	})
	return r
}

// agendaRows flattens the span into headed days.
func agendaRows(events []CalendarEvent, hits []int, from core.Date, days int, loc *time.Location) []agendaRow {
	buckets := make([][]int, days)
	for _, h := range hits {
		ev := events[h]
		first := max(0, monthDays(ev.Start.In(loc), from)*-1)
		last := min(days-1, -monthDays(ev.End.Add(-time.Nanosecond).In(loc), from))
		for i := first; i <= last; i++ {
			buckets[i] = append(buckets[i], h)
		}
	}
	var out []agendaRow
	for i, b := range buckets {
		if len(b) > 0 {
			out = append(out, agendaRow{from.AddDays(i), -1})
			for _, h := range b {
				out = append(out, agendaRow{from.AddDays(i), h})
			}
		}
	}
	return out
}

// agendaItem is one event row; it reports a click.
func agendaItem(c *ui.Context, ev CalendarEvent, loc *time.Location, hour12 bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	when := eventTimeText(c, ev, loc, hour12)
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Justify(ui.Start).Gap(d.Unit() * 3).PaddingX(d.Unit() * 3).Label(ev.Title + ", " + when)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	col := EventTone(k, ev.Color)
	b.Children(func() {
		ui.Text(c, when).Width(d.ControlHeight() * 3.5).Shrink(0).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
		ui.Box(c).Size(8, 8).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, r.W), col)
		})
		ui.Column(c).Grow(1).Shrink(1).Children(func() {
			ui.Text(c, ev.Title).TextColor(k.Text).SingleLine().Ellipsis("…")
			if ev.Location != "" {
				ui.Text(c, ev.Location).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().Ellipsis("…")
			}
		})
	})
	return b.Clicked()
}
