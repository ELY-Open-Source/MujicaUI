package chat

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DateSeparatorOptions configure a DateSeparator.
type DateSeparatorOptions struct {
	// Today decides "Today" and "Yesterday"; zero uses the window clock.
	Today core.Date
}

// DateSeparator marks where a day's messages begin, localized.
func DateSeparator(c *ui.Context, day core.Date, opts DateSeparatorOptions) *ui.Element {
	if !day.Valid() {
		panic("mujica: DateSeparator needs a valid date")
	}
	today := opts.Today
	if today.IsZero() {
		today = core.DateOf(c.Now())
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
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
			ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		})
		line()
	})
	return row
}

// dateSeparatorLabel names day relative to today, in the window's language.
func dateSeparatorLabel(c *ui.Context, day, today core.Date) string {
	switch day {
	case today:
		return core.Msg(c, "date_separator.today")
	case today.AddDays(-1):
		return core.Msg(c, "date_separator.yesterday")
	}
	return fmt.Sprintf(core.Msg(c, "calendar.day_label"), core.Msg(c, fmt.Sprintf("calendar.weekday_long_%d", day.Weekday())),
		core.Msg(c, fmt.Sprintf("calendar.month_%d", day.Month)), day.Day, day.Year)
}

// DateSeparatorWhen is a short stamp: time, day or date.
func DateSeparatorWhen(c *ui.Context, t time.Time, today core.Date) string {
	day := core.DateOf(t)
	if day == today {
		return t.Format("15:04")
	}
	if day == today.AddDays(-1) {
		return core.Msg(c, "date_separator.yesterday")
	}
	return fmt.Sprintf(core.Msg(c, "date_separator.short"), core.Msg(c, fmt.Sprintf("calendar.month_%d", day.Month)), day.Day)
}
