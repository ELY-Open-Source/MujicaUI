package datetime

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
)

// WeekPickerOptions configure a WeekPicker.
type WeekPickerOptions struct {
	WeekStart WeekStart
	// Min and Max bound the chosen week's days.
	Min, Max core.Date
	// Today is marked; zero uses the system clock.
	Today     core.Date
	Clearable bool
	Disabled  bool
	Label     string
}

// WeekPicker chooses a week; *value is its first day.
func WeekPicker(c *ui.Context, value *core.Date, opts WeekPickerOptions) core.ChoiceResult {
	core.DatetimeCheckDate("WeekPicker value", *value)
	core.DatetimeCheckBounds("WeekPicker", opts.Min, opts.Max)
	first := CalendarFirstWeekday(c, opts.WeekStart)
	weekOf := func(d core.Date) core.Date { return d.AddDays(-((int(d.Weekday()) - int(first) + 7) % 7)) }
	if !value.IsZero() && weekOf(*value) != *value {
		panic(fmt.Sprintf("mujica: WeekPicker value %v does not start a week", *value))
	}
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "week_picker.label")
	}
	enabled := func(x core.Date) bool {
		w := weekOf(x)
		return w.AddDays(6).Within(opts.Min, core.Date{}) && w.Within(core.Date{}, opts.Max)
	}
	col := ui.Column(c).Width(d.ControlHeight() * 8).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "week-picker-open", func() bool { return false })
	cursor := ui.Local(col, "week-picker-cursor", func() core.Date { return core.Date{} })
	hover := ui.Local(col, "week-picker-hover", func() core.Date { return core.Date{} })
	old := *value
	text := ""
	if !value.IsZero() {
		text = WeekViewTitle(c, *value, value.AddDays(6))
	}
	col.Children(func() {
		field := core.DatetimeField(c, ui.Row(c), false, false)
		field.Children(func() {
			t := core.DatetimeTrigger(c, text, core.Msg(c, "week_picker.placeholder"), "calendar").Label(label).Description(text)
			if !t.IsDisabled() && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open, *hover = !*open, core.Date{}
				*cursor = *value
				if cursor.IsZero() {
					*cursor = today
				}
			}
			if opts.Clearable && !value.IsZero() {
				if b := core.IconAction(c, icons.Must("x"), core.Msg(c, "calendar.clear")); b.Clicked() && !b.IsDisabled() {
					*value = core.Date{}
				}
			}
		})
		core.DatetimePopup(c, field, open, false, func() {
			band := *hover
			if band.IsZero() {
				band = *cursor
			}
			picked, hovered := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x core.Date) bool { return !value.IsZero() && weekOf(x) == *value },
				inRange:   func(x core.Date) bool { return weekOf(x) == weekOf(band) },
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
	return core.NewChoiceResult(col, changed)
}
