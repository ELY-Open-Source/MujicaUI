package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// DateSeparatorOptions configure a DateSeparator.
type DateSeparatorOptions struct {
	// Today decides "Today" and "Yesterday"; zero uses the window clock.
	Today Date
}

// DateSeparator marks where a day's messages begin, localized.
func DateSeparator(c *ui.Context, day Date, opts DateSeparatorOptions) *ui.Element {
	if !day.Valid() {
		panic("mujica: DateSeparator needs a valid date")
	}
	today := opts.Today
	if today.IsZero() {
		today = DateOf(c.Now())
	}
	k := Tokens(c)
	u := Density(c).Unit()
	label := dateSeparatorLabel(c, day, today)
	row := ui.Row(c).FillWidth().Gap(u*3).AlignItems(ui.Center).Padding(u*2, 0).Role(ui.RoleText).Label(label)
	line := func() {
		ui.Box(c).Grow(1).Height(theme.BorderWidth).Background(k.Border)
	}
	row.Children(func() {
		line()
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Shrink(0).Children(func() {
			ui.Box(c).Size(6, 6).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(r.X+3, r.Y+3, 6), k.Ornament)
			})
			ui.Text(c, label).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		})
		line()
	})
	return row
}

// dateSeparatorLabel names day relative to today, in the window's language.
func dateSeparatorLabel(c *ui.Context, day, today Date) string {
	switch day {
	case today:
		return msg(c, "date_separator.today")
	case today.AddDays(-1):
		return msg(c, "date_separator.yesterday")
	}
	return fmt.Sprintf(msg(c, "calendar.day_label"), msg(c, fmt.Sprintf("calendar.weekday_long_%d", day.Weekday())),
		msg(c, fmt.Sprintf("calendar.month_%d", day.Month)), day.Day, day.Year)
}

// dateSeparatorWhen is a short stamp: time, day or date.
func dateSeparatorWhen(c *ui.Context, t time.Time, today Date) string {
	day := DateOf(t)
	if day == today {
		return t.Format("15:04")
	}
	if day == today.AddDays(-1) {
		return msg(c, "date_separator.yesterday")
	}
	return fmt.Sprintf(msg(c, "date_separator.short"), msg(c, fmt.Sprintf("calendar.month_%d", day.Month)), day.Day)
}
