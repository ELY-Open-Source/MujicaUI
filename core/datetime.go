package core

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Date is a calendar day without clock or location; zero means none.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// DateOf returns the day t falls on in t's own location, without converting it.
func DateOf(t time.Time) Date {
	y, m, d := t.Date()
	return Date{y, m, d}
}

// IsZero reports whether d is the empty date.
func (d Date) IsZero() bool { return d == Date{} }

// Valid reports whether d names a real day.
func (d Date) Valid() bool { return !d.IsZero() && DateOf(d.Civil()) == d }

// String formats d as 2006-01-02, or "" when empty.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Format formats d with a Go time layout.
func (d Date) Format(layout string) string { return d.Civil().Format(layout) }

// Compare returns -1, 0 or 1 as d is before, equal to or after o.
func (d Date) Compare(o Date) int { return d.Civil().Compare(o.Civil()) }

// AddDays returns d moved by n days.
func (d Date) AddDays(n int) Date { return DateOf(d.Civil().AddDate(0, 0, n)) }

// AddMonths returns d moved by n months, keeping the day within the month.
func (d Date) AddMonths(n int) Date {
	first := time.Date(d.Year, d.Month+time.Month(n), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1).Day()
	return Date{first.Year(), first.Month(), min(d.Day, last)}
}

// Weekday returns the day of the week of d.
func (d Date) Weekday() time.Weekday { return d.Civil().Weekday() }

// Civil is d at midnight on a fixed clock, only for arithmetic.
func (d Date) Civil() time.Time { return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC) }

// Within reports whether d lies inside the optional bounds lo and hi.
func (d Date) Within(lo, hi Date) bool {
	return (lo.IsZero() || d.Compare(lo) >= 0) && (hi.IsZero() || d.Compare(hi) <= 0)
}

// DatetimeCheckDate panics when a non-empty date is not a real day.
func DatetimeCheckDate(name string, d Date) {
	if !d.IsZero() && !d.Valid() {
		panic(fmt.Sprintf("mujica: %s %v is not a calendar date", name, d))
	}
}

// DatetimeCheckBounds panics on invalid or inverted Min and Max.
func DatetimeCheckBounds(name string, lo, hi Date) {
	DatetimeCheckDate(name+" Min", lo)
	DatetimeCheckDate(name+" Max", hi)
	if !lo.IsZero() && !hi.IsZero() && lo.Compare(hi) > 0 {
		panic(fmt.Sprintf("mujica: %s Min %v is after Max %v", name, lo, hi))
	}
}

// DatetimeToday returns today, or the caller's fixed day when set.
func DatetimeToday(c *ui.Context, fixed Date) Date {
	if !fixed.IsZero() {
		DatetimeCheckDate("Today", fixed)
		return fixed
	}
	return DateOf(c.Now())
}

// DatetimeMonthName returns the localized name of m.
func DatetimeMonthName(c *ui.Context, m time.Month) string {
	return Msg(c, fmt.Sprintf("calendar.month_%d", int(m)))
}

// DatetimeMonthTitle returns the localized title of a month, as "October 2026".
func DatetimeMonthTitle(c *ui.Context, y int, m time.Month) string {
	return fmt.Sprintf(Msg(c, "calendar.title"), DatetimeMonthName(c, m), y)
}

// DatetimeDayLabel returns the localized spoken name of d.
func DatetimeDayLabel(c *ui.Context, d Date) string {
	wd := Msg(c, fmt.Sprintf("calendar.weekday_long_%d", int(d.Weekday())))
	return fmt.Sprintf(Msg(c, "calendar.day_label"), wd, DatetimeMonthName(c, d.Month), d.Day, d.Year)
}

// DatetimeField styles row as a single-line field: solid fill and fine border.
func DatetimeField(c *ui.Context, row *ui.Element, invalid, readOnly bool) *ui.Element {
	k := Tokens(c)
	d := Density(c)
	bg, border := k.Surface, k.ControlBorder
	switch {
	case invalid:
		border = k.Danger
	case readOnly:
		bg, border = k.Background, k.Border
	case row.FocusWithin():
		border = k.Focus
	case row.Hovered():
		border = k.Text
	}
	return row.MinHeight(d.ControlHeight()).Padding(0, d.Unit()).PaddingX(d.Unit()*3).Gap(d.Unit()).
		AlignItems(ui.Center).Radius(theme.ControlRadius).Background(bg).Border(theme.BorderWidth, border)
}

// DatetimeTrigger is the popup button inside a field; Cmd+C copies its text.
func DatetimeTrigger(c *ui.Context, text, placeholder, icon string) *ui.Element {
	k := Tokens(c)
	b := ui.ButtonBase(c).Grow(1).Shrink(1).Justify(ui.Start).Gap(Density(c).Unit() * 2).
		MinHeight(theme.MinHitSize).Radius(theme.ControlRadius).Role(ui.RolePopUpButton)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Shortcut(ui.Cmd, ui.KeyC) && text != "" {
		c.WriteClipboard(text)
	}
	b.Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(16, 16).TextColor(k.TextMuted)
		if text == "" {
			ui.Text(c, placeholder).TextColor(k.TextMuted).SingleLine()
		} else {
			ui.Text(c, text).TextColor(k.Text).SingleLine().FontFeatures("tnum")
		}
	})
	return b
}

// DatetimePopup shows fn below anchor while *open; a disabled or read-only anchor closes it.
func DatetimePopup(c *ui.Context, anchor *ui.Element, open *bool, readOnly bool, fn func()) *ui.Element {
	d := Density(c)
	if *open && (readOnly || anchor.IsDisabled()) {
		slog.Debug("mujica: datetime popup closed by state", "readOnly", readOnly)
		*open = false
	}
	return ui.PopoverBase(c, anchor, open, func(panel *ui.Element) {
		OverlayStyle(c, panel, OverlayEnter(c, panel))
		panel.MaxWidth(d.ControlHeight()*7+d.Unit()*6+2*theme.BorderWidth).Margin(d.Unit(), 0, 0, 0).Padding(d.Unit() * 3).Gap(d.Unit() * 3)
		fn()
	})
}

// DatetimeSelectedMark strokes the 1 DIP inset frame that marks a chosen cell.
func DatetimeSelectedMark(p *ui.Painter, r ui.Rect, c ui.Color) {
	const inset = 2
	p.Stroke(ui.Rect{X: r.X + inset, Y: r.Y + inset, W: r.W - 2*inset, H: r.H - 2*inset}, c, theme.ControlRadius/2, theme.BorderWidth)
}

// DatetimeError shows a field error in Danger beneath a control.
func DatetimeError(c *ui.Context, text string) {
	if text == "" {
		return
	}
	ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).TextColor(Tokens(c).Danger).Role(ui.RoleStatus)
}

// DatetimeFooter lays out popup actions at the end of a row.
func DatetimeFooter(c *ui.Context, fn func()) {
	ui.Row(c).Gap(Density(c).Unit() * 2).Justify(ui.End).AlignItems(ui.Center).FillWidth().Children(fn)
}
