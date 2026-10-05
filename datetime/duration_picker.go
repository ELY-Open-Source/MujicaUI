package datetime

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// DurationPickerOptions configure a DurationPicker.
type DurationPickerOptions struct {
	// Seconds adds a seconds field.
	Seconds bool
	// Min and Max bound the duration; zero Max means 99:59:59.
	Min, Max time.Duration
	Disabled bool
	Label    string
}

// DurationPicker edits hours, minutes and seconds within a range.
func DurationPicker(c *ui.Context, value *time.Duration, opts DurationPickerOptions) core.ChoiceResult {
	hi := opts.Max
	if hi == 0 {
		hi = 99*time.Hour + 59*time.Minute + 59*time.Second
	}
	unit := time.Minute
	if opts.Seconds {
		unit = time.Second
	}
	switch {
	case opts.Min < 0 || hi < opts.Min:
		panic(fmt.Sprintf("mujica: DurationPicker range %v..%v is not valid", opts.Min, hi))
	case *value < 0 || *value%unit != 0:
		panic(fmt.Sprintf("mujica: DurationPicker value %v is negative or finer than the fields", *value))
	}
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "duration_picker.label")
	}
	col := ui.Column(c).Gap(d.Unit()).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	draft := ui.Local(col, "duration-draft", func() [3]float64 { return durationParts(*value) })
	seen := ui.Local(col, "duration-seen", func() time.Duration { return *value })
	if *seen != *value {
		*draft, *seen = durationParts(*value), *value
	}
	changed, edited := false, false
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			fields := []struct {
				key string
				max float64
			}{{"duration_picker.hours", math.Floor(hi.Hours())}, {"duration_picker.minutes", 59}, {"duration_picker.seconds", 59}}
			for i, f := range fields {
				if i == 2 && !opts.Seconds {
					break
				}
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Width(d.ControlHeight() * 3.2).Children(func() {
						r := input.NumberInput(c, &draft[i], input.NumberInputOptions{Min: 0, Max: f.max, Label: core.Msg(c, f.key), Disabled: opts.Disabled})
						edited = edited || r.Changed()
					})
					ui.Text(c, core.Msg(c, f.key+"_unit")).TextColor(core.Tokens(c).TextMuted)
				})
			}
		})
		total := time.Duration(draft[0])*time.Hour + time.Duration(draft[1])*time.Minute + time.Duration(draft[2])*time.Second
		problem := ""
		if total < opts.Min || total > hi {
			problem = fmt.Sprintf(core.Msg(c, "duration_picker.range"), DurationText(c, opts.Min), DurationText(c, hi))
		}
		if edited && problem == "" && total != *value && !col.IsDisabled() {
			slog.Debug("mujica: duration_picker change", "from", *value, "to", total)
			*value, *seen, changed = total, total, true
		}
		core.DatetimeError(c, problem)
	})
	return core.NewChoiceResult(col, changed)
}

func durationParts(v time.Duration) [3]float64 {
	return [3]float64{float64(v / time.Hour), float64(v % time.Hour / time.Minute), float64(v % time.Minute / time.Second)}
}

// DurationText words a duration, as "1 h 30 min".
func DurationText(c *ui.Context, v time.Duration) string {
	p := durationParts(v)
	s := ""
	for i, key := range []string{"duration_picker.hours_unit", "duration_picker.minutes_unit", "duration_picker.seconds_unit"} {
		if p[i] > 0 || i == 1 && s == "" && p[2] == 0 {
			if s != "" {
				s += " "
			}
			s += fmt.Sprintf("%d %s", int(p[i]), core.Msg(c, key))
		}
	}
	return s
}
