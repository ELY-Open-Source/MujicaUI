package datetime

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// WeekStart is a month's first column; zero follows the locale.
type WeekStart int

const (
	WeekStartLocale WeekStart = iota
	WeekStartSunday
	WeekStartMonday
	WeekStartSaturday
)

// CalendarOptions configure a Calendar.
type CalendarOptions struct {
	// Min and Max bound choosable days; zero is open.
	Min, Max core.Date
	// DisableDate marks single days that cannot be chosen.
	DisableDate func(core.Date) bool
	WeekStart   WeekStart
	// Today is marked; zero uses the system clock.
	Today    core.Date
	Disabled bool
	// Label names the calendar for assistive technology.
	Label string
}

// Calendar shows one month and chooses a day into *value.
func Calendar(c *ui.Context, value *core.Date, opts CalendarOptions) DateTimeResult {
	core.DatetimeCheckDate("Calendar value", *value)
	core.DatetimeCheckBounds("Calendar", opts.Min, opts.Max)
	box := ui.Column(c)
	cursor := ui.Local(box, "calendar-cursor", func() core.Date { return core.Date{} })
	seen := ui.Local(box, "calendar-seen", func() core.Date { return core.Date{} })
	today := core.DatetimeToday(c, opts.Today)
	if *seen != *value || cursor.IsZero() {
		*seen = *value
		*cursor = *value
		if cursor.IsZero() {
			*cursor = today
		}
	}
	changed := false
	box.Disabled(opts.Disabled).Children(func() {
		picked, _ := calendarMonth(c, cursor, calendarConfig{
			chosen:    func(d core.Date) bool { return d == *value },
			enabled:   calendarEnabled(opts.Min, opts.Max, opts.DisableDate),
			today:     today,
			weekStart: opts.WeekStart,
			label:     opts.Label,
		})
		if !picked.IsZero() && picked != *value {
			slog.Debug("mujica: calendar change", "from", *value, "to", picked)
			*value, *seen = picked, picked
			changed = true
		}
	})
	return DateTimeResult{Element: box, changed: changed}
}

// calendarConfig is how calendarMonth marks and limits its days.
type calendarConfig struct {
	chosen    func(core.Date) bool
	inRange   func(core.Date) bool
	enabled   func(core.Date) bool
	today     core.Date
	weekStart WeekStart
	autoFocus bool
	label     string
}

// calendarEnabled combines bounds and a day filter.
func calendarEnabled(lo, hi core.Date, disable func(core.Date) bool) func(core.Date) bool {
	return func(d core.Date) bool { return d.Within(lo, hi) && (disable == nil || !disable(d)) }
}

// CalendarFirstWeekday resolves the first column for the window's locale.
func CalendarFirstWeekday(c *ui.Context, w WeekStart) time.Weekday {
	switch w {
	case WeekStartSunday:
		return time.Sunday
	case WeekStartMonday:
		return time.Monday
	case WeekStartSaturday:
		return time.Saturday
	case WeekStartLocale:
		if core.CurrentLocale(c) == core.ZhCN {
			return time.Monday
		}
		return time.Sunday
	}
	panic(fmt.Sprintf("mujica: unknown WeekStart %d", w))
}

// calendarSeek returns the first enabled day from d stepping by dir, or d.
func calendarSeek(d core.Date, dir int, enabled func(core.Date) bool) core.Date {
	for i, x := 0, d; i < 31; i, x = i+1, x.AddDays(dir) {
		if enabled(x) {
			return x
		}
	}
	return d
}

// calendarMonth draws *cursor's month; returns picked and hovered days.
func calendarMonth(c *ui.Context, cursor *core.Date, cfg calendarConfig) (picked, hovered core.Date) {
	k := core.Tokens(c)
	d := core.Density(c)
	cur := *cursor
	title := core.DatetimeMonthTitle(c, cur.Year, cur.Month)
	label := cfg.label
	if label == "" {
		label = core.Msg(c, "calendar.label")
	}
	cell := d.ControlHeight()
	root := ui.Column(c).Gap(d.Unit() * 2).Shrink(0).Width(cell * 7).Role(ui.RoleGroup).Label(label)
	keyed := ui.Local(root, "calendar-keys", func() bool { return false })
	move := func(to core.Date, dir int, keys bool) {
		*cursor = calendarSeek(to, dir, cfg.enabled)
		*keyed = keys
	}
	root.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(d.Unit()).Children(func() {
			if core.IconAction(c, icons.Must("chevrons-left"), core.Msg(c, "calendar.prev_year")).Clicked() {
				move(cur.AddMonths(-12), -1, false)
			}
			if core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "calendar.prev_month")).Clicked() {
				move(cur.AddMonths(-1), -1, false)
			}
			ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.BodySize+2)).Bold().
				TextColor(k.Text).Grow(1).TextAlign(ui.Center).SingleLine()
			if core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "calendar.next_month")).Clicked() {
				move(cur.AddMonths(1), 1, false)
			}
			if core.IconAction(c, icons.Must("chevrons-right"), core.Msg(c, "calendar.next_year")).Clicked() {
				move(cur.AddMonths(12), 1, false)
			}
		})
		first := CalendarFirstWeekday(c, cfg.weekStart)
		ui.Row(c).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).PaddingY(d.Unit()).Children(func() {
			for i := range 7 {
				wd := (int(first) + i) % 7
				ui.Text(c, core.Msg(c, "calendar.weekday_"+strconv.Itoa(wd))).Width(cell).TextAlign(ui.Center).
					FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleColumnHeader).
					Label(core.Msg(c, "calendar.weekday_long_"+strconv.Itoa(wd)))
			}
		})
		days := ui.Column(c).Role(ui.RoleTable).Label(title).FocusGroup(ui.Horizontal | ui.Vertical)
		keysMoved := true
		switch {
		case days.Shortcut(0, ui.KeyLeft):
			move(cur.AddDays(-1), -1, true)
		case days.Shortcut(0, ui.KeyRight):
			move(cur.AddDays(1), 1, true)
		case days.Shortcut(0, ui.KeyUp):
			move(cur.AddDays(-7), -1, true)
		case days.Shortcut(0, ui.KeyDown):
			move(cur.AddDays(7), 1, true)
		case days.Shortcut(0, ui.KeyPageUp):
			move(cur.AddMonths(-1), -1, true)
		case days.Shortcut(0, ui.KeyPageDown):
			move(cur.AddMonths(1), 1, true)
		case days.Shortcut(ui.Shift, ui.KeyPageUp):
			move(cur.AddMonths(-12), -1, true)
		case days.Shortcut(ui.Shift, ui.KeyPageDown):
			move(cur.AddMonths(12), 1, true)
		case days.Shortcut(0, ui.KeyHome):
			move(core.Date{Year: cur.Year, Month: cur.Month, Day: 1}, 1, true)
		case days.Shortcut(0, ui.KeyEnd):
			move(core.Date{Year: cur.Year, Month: cur.Month, Day: 1}.AddMonths(1).AddDays(-1), -1, true)
		default:
			keysMoved = false
		}
		cur = *cursor
		start := core.Date{Year: cur.Year, Month: cur.Month, Day: 1}
		start = start.AddDays(-((int(start.Weekday()) - int(first) + 7) % 7))
		days.Children(func() {
			for w := range 6 {
				ui.Row(c).Role(ui.RoleRow).Children(func() {
					for i := range 7 {
						day := start.AddDays(w*7 + i)
						clicked, hover := calendarDay(c, day, cur, cfg, cell, keysMoved, *keyed)
						if clicked {
							picked = day
							move(day, 1, false)
						}
						if hover {
							hovered = day
						}
					}
				})
			}
		})
	})
	return picked, hovered
}

// calendarDay draws one day cell and reports whether it was picked or hovered.
func calendarDay(c *ui.Context, day, cur core.Date, cfg calendarConfig, size float32, focus, keyed bool) (picked, hovered bool) {
	k := core.Tokens(c)
	enabled := cfg.enabled(day)
	chosen := cfg.chosen(day)
	sel := core.Date{}
	if chosen {
		sel = day
	}
	cell := ui.RadioBase(c, &sel, day).Size(size, size-4).Center().Radius(theme.ControlRadius).
		Label(core.DatetimeDayLabel(c, day)).Disabled(!enabled).FocusRing(false)
	if day == cur && enabled {
		if focus {
			cell.Focus()
		} else if cfg.autoFocus {
			core.AutoFocus(c, cell)
		}
	}
	band := cfg.inRange != nil && cfg.inRange(day)
	fg := k.Text
	if day.Month != cur.Month {
		fg = k.TextMuted
	}
	switch {
	case chosen:
		cell.Background(k.Accent)
		fg = k.OnAccent
		if cell.Hovered() {
			cell.Background(k.AccentHover)
		}
	case band:
		cell.Background(k.Selection).Radius(0)
	case cell.Pressed():
		cell.Background(k.Selection)
	case cell.Hovered():
		cell.Background(k.SurfaceHover)
	}
	ring := cell.FocusVisible() || cell.Focused() && keyed
	cell.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if chosen {
			core.DatetimeSelectedMark(p, r, k.OnAccent)
		}
		if day == cfg.today {
			p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H-5, 5), fg)
		}
		if ring {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	cell.Children(func() {
		t := ui.Text(c, strconv.Itoa(day.Day)).TextColor(fg).FontFeatures("tnum")
		if day == cfg.today {
			t.Bold()
		}
		if !enabled {
			t.Strikethrough()
		}
	})
	return cell.Clicked() && enabled, cell.Hovered()
}
