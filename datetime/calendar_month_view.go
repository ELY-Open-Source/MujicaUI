package datetime

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CalendarMonthViewOptions configure a CalendarMonthView.
type CalendarMonthViewOptions struct {
	// Location is the zone days are cut in; required.
	Location *time.Location
	// Version must change when events change in place.
	Version   int
	WeekStart WeekStart
	// MaxPerDay is rows shown before "+N"; zero means 3.
	MaxPerDay int
	// Today is marked; zero uses the system clock.
	Today    core.Date
	Hour12   bool
	Disabled bool
	Label    string
}

// monthDrag is the event id a month view chip drags.
type monthDrag string

// monthDay is one cell's rows of side-by-side events.
type monthDay struct {
	rows [][]int
	n    int
}

// CalendarMonthView shows a month of events; drags reschedule them.
func CalendarMonthView(c *ui.Context, month *YearMonth, events []CalendarEvent, opts CalendarMonthViewOptions) CalendarViewResult {
	loc := eventLocation("CalendarMonthView", opts.Location)
	if month.IsZero() || month.Month < time.January || month.Month > time.December {
		panic(fmt.Sprintf("mujica: CalendarMonthView month %v is not valid", *month))
	}
	if opts.MaxPerDay < 0 {
		panic("mujica: CalendarMonthView MaxPerDay is negative")
	}
	maxRows := opts.MaxPerDay
	if maxRows == 0 {
		maxRows = 3
	}
	eventsCheck("CalendarMonthView", events)
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	first := CalendarFirstWeekday(c, opts.WeekStart)
	start := core.Date{Year: month.Year, Month: month.Month, Day: 1}
	start = start.AddDays(-((int(start.Weekday()) - int(first) + 7) % 7))
	label := opts.Label
	if label == "" {
		label = core.DatetimeMonthTitle(c, month.Year, month.Month)
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	problem := ui.Local(root, "month-view-problem", func() string { return "" })
	more := ui.Local(root, "month-view-more", func() core.Date { return core.Date{} })
	days := eventsDerived(root, "month-view-days", events, opts.Version, eventDayStart(start, loc), eventDayStart(start.AddDays(42), loc),
		func(hits []int) []monthDay { return monthLayout(events, hits, start, loc) })
	disabled := func() bool { return root.IsDisabled() }
	apply := func(ch CalendarChange, p string) {
		*problem = p
		if p == "" && !disabled() {
			slog.Debug("mujica: calendar_month_view move", "id", ch.ID, "start", ch.Start)
			r.change = &ch
		}
	}
	root.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(d.Unit()).Children(func() {
			if core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "calendar.prev_month")).Clicked() && !disabled() {
				*month = YearMonth{start.AddDays(7).AddMonths(-1).Year, start.AddDays(7).AddMonths(-1).Month}
			}
			ui.Text(c, core.DatetimeMonthTitle(c, month.Year, month.Month)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).
				TextColor(k.Text).Grow(1).TextAlign(ui.Center)
			if core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "calendar.next_month")).Clicked() && !disabled() {
				*month = YearMonth{start.AddDays(7).AddMonths(1).Year, start.AddDays(7).AddMonths(1).Month}
			}
		})
		ui.Row(c).Children(func() {
			for i := range 7 {
				wd := (int(first) + i) % 7
				ui.Text(c, core.Msg(c, "calendar.weekday_"+strconv.Itoa(wd))).Grow(1).Basis(0).TextAlign(ui.Center).
					FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Label(core.Msg(c, "calendar.weekday_long_"+strconv.Itoa(wd)))
			}
		})
		ui.Column(c).Border(theme.BorderWidth, k.Border).Dividers(theme.BorderWidth, k.Border).FocusGroup(ui.Horizontal | ui.Vertical).Children(func() {
			for w := range 6 {
				ui.Row(c).Dividers(theme.BorderWidth, k.Border).Children(func() {
					for i := range 7 {
						at := w*7 + i
						day := start.AddDays(at)
						monthCell(c, day, *month, today, days[at], events, maxRows, loc, opts, &r, more, apply)
					}
				})
			}
		})
		core.DatetimeError(c, *problem)
	})
	return r
}

// monthLayout buckets events per day, concurrent ones side by side.
func monthLayout(events []CalendarEvent, hits []int, start core.Date, loc *time.Location) []monthDay {
	out := make([]monthDay, 42)
	for i := range out {
		from, to := eventDayStart(start.AddDays(i), loc), eventDayStart(start.AddDays(i+1), loc)
		var idx []int
		var spans []eventSpan
		for _, h := range hits {
			ev := events[h]
			if ev.Start.Before(to) && ev.End.After(from) {
				idx = append(idx, h)
				spans = append(spans, eventSpan{float64(max(ev.Start.Unix(), from.Unix())), float64(min(ev.End.Unix(), to.Unix()))})
			}
		}
		col, _ := eventColumns(spans)
		for j, h := range idx {
			if col[j] == 0 || len(out[i].rows) == 0 {
				out[i].rows = append(out[i].rows, []int{h})
			} else {
				last := len(out[i].rows) - 1
				out[i].rows[last] = append(out[i].rows[last], h)
			}
		}
		out[i].n = len(idx)
	}
	return out
}

// monthCell draws one day: number, chips and "+N".
func monthCell(c *ui.Context, day core.Date, month YearMonth, today core.Date, lay monthDay, events []CalendarEvent, maxRows int,
	loc *time.Location, opts CalendarMonthViewOptions, r *CalendarViewResult, more *core.Date, apply func(CalendarChange, string)) {
	k := core.Tokens(c)
	d := core.Density(c)
	cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).MinHeight(d.ControlHeight() * 3.5).Padding(d.Unit()).Gap(d.Unit()).Clip()
	if id, ok := ui.Drop[monthDrag](cell); ok && !cell.IsDisabled() {
		for _, ev := range events {
			if ev.ID == string(id) {
				apply(eventShift(c, ev, monthDays(ev.Start.In(loc), day), loc))
			}
		}
	}
	if _, over := ui.DragOver[monthDrag](cell); over {
		cell.Background(k.Selection)
		cell.DrawOver(func(p *ui.Painter, rr ui.Rect) { core.DatetimeSelectedMark(p, rr, k.Focus) })
	}
	cell.Children(func() {
		fg := k.Text
		if day.Month != month.Month {
			fg = k.TextMuted
		}
		num := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.MinHitSize / 2).Label(core.DatetimeDayLabel(c, day))
		internal.Ring(num, k.Focus, theme.MinHitSize/2)
		if day == today {
			num.Background(k.Selection).Border(theme.BorderWidth, k.AccentText)
		} else if num.Hovered() {
			num.Background(k.SurfaceHover)
		}
		num.Children(func() {
			t := ui.Text(c, strconv.Itoa(day.Day)).TextColor(fg).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
			if day == today {
				t.Bold()
			}
		})
		if num.Clicked() && !num.IsDisabled() {
			r.day = day
		}
		shown := 0
		for i, row := range lay.rows {
			if i >= maxRows {
				break
			}
			ui.Row(c).Gap(d.Unit() / 2).Children(func() {
				for _, h := range row {
					monthChip(c, events[h], loc, opts, r, apply)
					shown++
				}
			})
		}
		if hidden := lay.n - shown; hidden > 0 {
			b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Justify(ui.Start).PaddingX(d.Unit()).Radius(theme.ControlRadius).
				Label(core.CalendarCount(c, "calendar_month_view.more_label", hidden))
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() {
				ui.Text(c, fmt.Sprintf(core.Msg(c, "calendar_month_view.more"), hidden)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText)
			})
			if b.Clicked() && !b.IsDisabled() {
				*more = day
			}
			open := *more == day
			ui.PopoverBase(c, b, &open, func(panel ui.Element) {
				core.OverlayStyle(c, panel, core.OverlayEnter(c, panel))
				_, wh := c.Size()
				panel.Padding(d.Unit() * 2).Gap(d.Unit()).MinWidth(d.ControlHeight() * 6).MaxHeight(wh / 2)
				ui.Text(c, core.DatetimeDayLabel(c, day)).Font(theme.SerifFont).TextColor(k.Text)
				ui.Scroll(c).Shrink(1).Children(func() {
					ui.Column(c).Gap(d.Unit()).Children(func() {
						for _, row := range lay.rows {
							for _, h := range row {
								monthChip(c, events[h], loc, opts, r, apply)
							}
						}
					})
				})
			})
			if !open && *more == day {
				*more = core.Date{}
			}
		}
	})
}

// monthChip is a draggable chip; Alt+arrows move it.
func monthChip(c *ui.Context, ev CalendarEvent, loc *time.Location, opts CalendarMonthViewOptions, r *CalendarViewResult, apply func(CalendarChange, string)) {
	chip := EventChip(c, ev, EventChipOptions{Location: loc, Hour12: opts.Hour12, Inline: true}).Grow(1).Basis(0).MinWidth(theme.MinHitSize)
	if chip.IsDisabled() {
		return
	}
	chip.Drag(monthDrag(ev.ID))
	if chip.Clicked() {
		r.opened = ev.ID
	}
	for key, days := range map[ui.Key]int{ui.KeyLeft: -1, ui.KeyRight: 1, ui.KeyUp: -7, ui.KeyDown: 7} {
		if chip.Shortcut(ui.Alt, key) {
			apply(eventShift(c, ev, days, loc))
		}
	}
}

// monthDays counts calendar days from t's day to day.
func monthDays(t time.Time, day core.Date) int {
	return int(day.Civil().Sub(core.DateOf(t).Civil()).Hours() / 24)
}
