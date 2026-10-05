package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// DueDatePickerOptions configure a DueDatePicker.
type DueDatePickerOptions struct {
	// Done stops a past date from showing as overdue.
	Done      bool
	WeekStart WeekStart
	// Today anchors shortcuts and overdue marks; zero uses today.
	Today    Date
	Disabled bool
	Label    string
}

// DueDatePicker combines date shortcuts, a calendar and overdue marks.
func DueDatePicker(c *ui.Context, value *Date, opts DueDatePickerOptions) DateTimeResult {
	datetimeCheckDate("DueDatePicker value", *value)
	k := Tokens(c)
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	first := calendarFirstWeekday(c, opts.WeekStart)
	nextWeek := today.AddDays(7 - (int(today.Weekday())-int(first)+7)%7)
	label := opts.Label
	if label == "" {
		label = msg(c, "due_date_picker.label")
	}
	old := *value
	col := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	var r DateTimeResult
	col.Children(func() {
		ui.Row(c).Gap(d.Unit()).Wrap().Children(func() {
			for _, q := range []struct {
				key string
				day Date
			}{{"due_date_picker.today", today}, {"due_date_picker.tomorrow", today.AddDays(1)}, {"due_date_picker.next_week", nextWeek}} {
				variant := Secondary
				text := msg(c, q.key)
				if *value == q.day {
					variant = Outline
					text = "✓ " + text
				}
				b := Button(c, text, ButtonOptions{Variant: variant}).Label(msg(c, q.key))
				if *value == q.day {
					b.Background(k.Selection).TextColor(k.Text)
				}
				if b.Clicked() && !b.IsDisabled() {
					*value = q.day
				}
			}
		})
		r = DatePicker(c, value, DatePickerOptions{Label: label, WeekStart: opts.WeekStart, Today: today, Clearable: true, Disabled: col.IsDisabled()})
		if text, late := projectDueText(c, *value, today, opts.Done); text != "" {
			fg, icon := k.TextMuted, "calendar"
			if late {
				fg, icon = k.Danger, "triangle-alert"
			}
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(fg).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize))
			})
		}
	})
	changed := *value != old
	if changed {
		slog.Debug("mujica: due_date_picker change", "from", old, "to", *value)
	}
	return DateTimeResult{Element: col, Input: r.Input, changed: changed}
}
