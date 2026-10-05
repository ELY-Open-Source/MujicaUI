package mujica

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// YearViewOptions configure a YearView.
type YearViewOptions struct {
	// Location is the zone days are cut in; required.
	Location *time.Location
	// Version must change when events change in place.
	Version   int
	WeekStart WeekStart
	// Today is marked; zero uses the system clock.
	Today    Date
	Disabled bool
	Label    string
}

// YearView shows twelve months with busy-day dots.
func YearView(c *ui.Context, year *int, events []CalendarEvent, opts YearViewOptions) CalendarViewResult {
	loc := eventLocation("YearView", opts.Location)
	if *year < 1 {
		panic(fmt.Sprintf("mujica: YearView year %d is not valid", *year))
	}
	eventsCheck("YearView", events)
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	first := Date{*year, time.January, 1}
	label := opts.Label
	if label == "" {
		label = strconv.Itoa(*year)
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Gap(d.Unit() * 3).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	counts := eventsDerived(root, "year-view-counts", events, opts.Version, eventDayStart(first, loc), eventDayStart(first.AddMonths(12), loc),
		func(hits []int) map[Date]int {
			m := map[Date]int{}
			for _, h := range hits {
				ev := events[h]
				for day := DateOf(ev.Start.In(loc)); eventDayStart(day, loc).Before(ev.End); day = day.AddDays(1) {
					m[day]++
				}
			}
			return m
		})
	cursor := ui.Local(root, "year-view-cursor", func() Date { return Date{} })
	if cursor.Year != *year {
		*cursor = Date{*year, today.Month, min(today.Day, 28)}
	}
	root.Children(func() {
		calendarViewHeader(c, strconv.Itoa(*year), msg(c, "calendar.prev_year"), msg(c, "calendar.next_year"), func(dir int) {
			*year += dir
		})
		grid := ui.Row(c).Wrap().AlignItems(ui.Start).Gap(d.Unit() * 4).FocusGroup(ui.Horizontal | ui.Vertical)
		moved := false
		if !grid.IsDisabled() {
			for key, days := range map[ui.Key]int{ui.KeyLeft: -1, ui.KeyRight: 1, ui.KeyUp: -7, ui.KeyDown: 7} {
				if grid.Shortcut(0, key) {
					*cursor, moved = cursor.AddDays(days), true
				}
			}
			for key, months := range map[ui.Key]int{ui.KeyPageUp: -1, ui.KeyPageDown: 1} {
				if grid.Shortcut(0, key) {
					*cursor, moved = cursor.AddMonths(months), true
				}
			}
			if moved && cursor.Year != *year {
				*year = cursor.Year
			}
		}
		grid.Children(func() {
			for m := time.January; m <= time.December; m++ {
				if day, ok := yearMonth(c, *year, m, counts, today, *cursor, moved, opts.WeekStart); ok && !grid.IsDisabled() {
					r.day, *cursor = day, day
				}
			}
		})
	})
	return r
}

// yearMonth draws one small month and reports a day clicked.
func yearMonth(c *ui.Context, year int, m time.Month, counts map[Date]int, today, cursor Date, focus bool, ws WeekStart) (Date, bool) {
	k := Tokens(c)
	cell := theme.MinHitSize
	first := calendarFirstWeekday(c, ws)
	start := Date{year, m, 1}
	lead := (int(start.Weekday()) - int(first) + 7) % 7
	var picked Date
	ui.Column(c).Width(cell * 7).Gap(2).Children(func() {
		ui.Text(c, datetimeMonthName(c, m)).Font(theme.SerifFont).TextColor(k.Text).Bold()
		ui.Row(c).Children(func() {
			for i := range 7 {
				wd := (int(first) + i) % 7
				ui.Text(c, msg(c, "calendar.weekday_"+strconv.Itoa(wd))).Width(cell).TextAlign(ui.Center).
					FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			}
		})
		last := start.AddMonths(1).AddDays(-1).Day
		for w := 0; w*7 < lead+last; w++ {
			ui.Row(c).Children(func() {
				for i := range 7 {
					n := w*7 + i - lead + 1
					if n < 1 || n > last {
						ui.Box(c).Size(cell, cell)
						continue
					}
					day := Date{year, m, n}
					busy := counts[day]
					text := datetimeDayLabel(c, day)
					if busy > 0 {
						text += ", " + calendarCount(c, "year_view.events", busy)
					}
					b := ui.ButtonBase(c).Size(cell, cell).Radius(theme.ControlRadius).Label(text).Column()
					internal.Ring(b, k.Focus, theme.ControlRadius)
					if day == cursor && focus {
						b.Focus()
					}
					switch {
					case day == today:
						b.Background(k.Selection)
					case b.Hovered():
						b.Background(k.SurfaceHover)
					}
					dots := yearDots(busy)
					b.Draw(func(p *ui.Painter, r ui.Rect) {
						x0 := r.X + r.W/2 - float32(dots*5-2)/2
						for j := range dots {
							p.Fill(ui.Rect{X: x0 + float32(j)*5, Y: r.Y + r.H - 6, W: 3, H: 3}, k.AccentText, 1.5)
						}
					})
					b.Children(func() {
						t := ui.Text(c, strconv.Itoa(n)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Text).FontFeatures("tnum")
						if day == today {
							t.Bold().TextColor(k.AccentText)
						}
					})
					if b.Clicked() {
						picked = day
					}
				}
			})
		}
	})
	return picked, !picked.IsZero()
}

// yearDots maps an event count to one to three dots.
func yearDots(n int) int {
	switch {
	case n == 0:
		return 0
	case n <= 2:
		return 1
	case n <= 5:
		return 2
	}
	return 3
}
