package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// ReminderPickerOptions configure a ReminderPicker.
type ReminderPickerOptions struct {
	// Choices are the lead times offered; empty uses defaults.
	Choices []time.Duration
	// Max limits how many are set; zero means any.
	Max      int
	Disabled bool
	Label    string
}

var reminderDefaults = []time.Duration{0, 5 * time.Minute, 15 * time.Minute, 30 * time.Minute, time.Hour, 24 * time.Hour, 7 * 24 * time.Hour}

// ReminderPicker toggles reminder lead times, kept in choice order.
func ReminderPicker(c *ui.Context, value *[]time.Duration, opts ReminderPickerOptions) ChoiceResult {
	choices := opts.Choices
	if len(choices) == 0 {
		choices = reminderDefaults
	}
	for i, d := range choices {
		if d < 0 || d%time.Minute != 0 || slices.Contains(choices[:i], d) {
			panic(fmt.Sprintf("mujica: ReminderPicker choice %v is negative, sub-minute or repeated", d))
		}
	}
	for _, d := range *value {
		if !slices.Contains(choices, d) {
			panic(fmt.Sprintf("mujica: ReminderPicker value %v is not a choice", d))
		}
	}
	if opts.Max < 0 || opts.Max > 0 && len(*value) > opts.Max {
		panic(fmt.Sprintf("mujica: ReminderPicker Max %d is negative or below the value", opts.Max))
	}
	full := opts.Max > 0 && len(*value) >= opts.Max
	items := make([]ToggleItem[time.Duration], len(choices))
	for i, d := range choices {
		items[i] = ToggleItem[time.Duration]{Value: d, Label: ReminderText(c, d), Disabled: full && !slices.Contains(*value, d)}
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "reminder_picker.label")
	}
	r := ToggleGroup(c, value, items, ToggleGroupOptions{Multiple: true, AllowEmpty: true, Vertical: true, Label: label, Disabled: opts.Disabled})
	if r.Changed() {
		slices.SortFunc(*value, func(a, b time.Duration) int { return slices.Index(choices, a) - slices.Index(choices, b) })
		slog.Debug("mujica: reminder_picker change", "value", *value)
	}
	return r
}

// ReminderText words a lead time, as "15 minutes before".
func ReminderText(c *ui.Context, d time.Duration) string {
	switch {
	case d == 0:
		return msg(c, "reminder_picker.at_start")
	case d%(7*24*time.Hour) == 0:
		return calendarCount(c, "reminder_picker.weeks", int(d/(7*24*time.Hour)))
	case d%(24*time.Hour) == 0:
		return calendarCount(c, "reminder_picker.days", int(d/(24*time.Hour)))
	case d%time.Hour == 0:
		return calendarCount(c, "reminder_picker.hours", int(d/time.Hour))
	}
	return calendarCount(c, "reminder_picker.minutes", int(d/time.Minute))
}
