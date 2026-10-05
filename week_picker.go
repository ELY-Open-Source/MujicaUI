package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// WeekPickerOptions configure a WeekPicker.
type WeekPickerOptions struct {
	WeekStart WeekStart
	// Min and Max bound the chosen week's days.
	Min, Max Date
	// Today is marked; zero uses the system clock.
	Today     Date
	Clearable bool
	Disabled  bool
	Label     string
}

// WeekPicker chooses a week; *value is its first day.
func WeekPicker(c *ui.Context, value *Date, opts WeekPickerOptions) ChoiceResult {
	datetimeCheckDate("WeekPicker value", *value)
	datetimeCheckBounds("WeekPicker", opts.Min, opts.Max)
	first := calendarFirstWeekday(c, opts.WeekStart)
	weekOf := func(d Date) Date { return d.AddDays(-((int(d.Weekday()) - int(first) + 7) % 7)) }
	if !value.IsZero() && weekOf(*value) != *value {
		panic(fmt.Sprintf("mujica: WeekPicker value %v does not start a week", *value))
	}
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	label := opts.Label
	if label == "" {
		label = msg(c, "week_picker.label")
	}
	enabled := func(x Date) bool {
		w := weekOf(x)
		return w.AddDays(6).within(opts.Min, Date{}) && w.within(Date{}, opts.Max)
	}
	col := ui.Column(c).Width(d.ControlHeight() * 8).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "week-picker-open", func() bool { return false })
	cursor := ui.Local(col, "week-picker-cursor", func() Date { return Date{} })
	hover := ui.Local(col, "week-picker-hover", func() Date { return Date{} })
	old := *value
	text := ""
	if !value.IsZero() {
		text = weekViewTitle(c, *value, value.AddDays(6))
	}
	col.Children(func() {
		field := datetimeField(c, ui.Row(c), false, false)
		field.Children(func() {
			t := datetimeTrigger(c, text, msg(c, "week_picker.placeholder"), "calendar").Label(label).Description(text)
			if !t.IsDisabled() && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open, *hover = !*open, Date{}
				*cursor = *value
				if cursor.IsZero() {
					*cursor = today
				}
			}
			if opts.Clearable && !value.IsZero() {
				if b := iconAction(c, icons.Must("x"), msg(c, "calendar.clear")); b.Clicked() && !b.IsDisabled() {
					*value = Date{}
				}
			}
		})
		datetimePopup(c, field, open, false, func() {
			band := *hover
			if band.IsZero() {
				band = *cursor
			}
			picked, hovered := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x Date) bool { return !value.IsZero() && weekOf(x) == *value },
				inRange:   func(x Date) bool { return weekOf(x) == weekOf(band) },
				enabled:   enabled,
				today:     today,
				weekStart: opts.WeekStart,
				autoFocus: true,
			})
			if hovered != *hover {
				*hover = hovered
				c.Invalidate()
			}
			if !picked.IsZero() {
				*value, *open = weekOf(picked), false
			}
		})
	})
	changed := *value != old
	if changed {
		slog.Debug("mujica: week_picker change", "from", old, "to", *value)
	}
	return ChoiceResult{Element: col, changed: changed}
}
