package project

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// DueDatePickerOptions configure a DueDatePicker.
type DueDatePickerOptions struct {
	// Done stops a past date from showing as overdue.
	Done      bool
	WeekStart datetime.WeekStart
	// Today anchors shortcuts and overdue marks; zero uses today.
	Today    core.Date
	Disabled bool
	Label    string
}

// DueDatePicker combines date shortcuts, a calendar and overdue marks.
func DueDatePicker(c *ui.Context, value *core.Date, opts DueDatePickerOptions) datetime.DateTimeResult {
	core.DatetimeCheckDate("DueDatePicker value", *value)
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	first := datetime.CalendarFirstWeekday(c, opts.WeekStart)
	nextWeek := today.AddDays(7 - (int(today.Weekday())-int(first)+7)%7)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "due_date_picker.label")
	}
	old := *value
	col := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	var r datetime.DateTimeResult
	col.Children(func() {
		ui.Row(c).Gap(d.Unit()).Wrap().Children(func() {
			for _, q := range []struct {
				key string
				day core.Date
			}{{"due_date_picker.today", today}, {"due_date_picker.tomorrow", today.AddDays(1)}, {"due_date_picker.next_week", nextWeek}} {
				variant := input.Secondary
				text := core.Msg(c, q.key)
				if *value == q.day {
					variant = input.Outline
					text = "✓ " + text
				}
				b := input.Button(c, text, input.ButtonOptions{Variant: variant}).Label(core.Msg(c, q.key))
				if *value == q.day {
					b.Background(k.Selection).TextColor(k.Text)
				}
				if b.Clicked() && !b.IsDisabled() {
					*value = q.day
				}
			}
		})
		r = datetime.DatePicker(c, value, datetime.DatePickerOptions{Label: label, WeekStart: opts.WeekStart, Today: today, Clearable: true, Disabled: col.IsDisabled()})
		if text, late := projectDueText(c, *value, today, opts.Done); text != "" {
			fg, icon := k.TextMuted, "calendar"
			if late {
				fg, icon = k.Danger, "triangle-alert"
			}
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(fg).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize))
			})
		}
	})
	changed := *value != old
	if changed {
		slog.Debug("mujica: due_date_picker change", "from", old, "to", *value)
	}
	return datetime.NewDateTimeResult(col, r.Input, changed, false)
}
