package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"
)

// RangePreset names a quick time range. The zero value is Range1D.
type RangePreset int

const (
	Range1D RangePreset = iota
	Range5D
	Range1M
	Range3M
	RangeYTD
	Range1Y
	Range5Y
	RangeCustom
)

var rangePresetNames = []string{"1D", "5D", "1M", "3M", "YTD", "1Y", "5Y"}

// MarketRange is a chosen range: a preset, or Custom with its dates.
type MarketRange struct {
	Preset RangePreset
	// Custom holds the dates when Preset is RangeCustom.
	Custom DateRange
}

// Bounds returns the range's start and end at now, in loc; custom ends include their last day.
func (r MarketRange) Bounds(now time.Time, loc *time.Location) (from, to time.Time) {
	now = now.In(loc)
	switch r.Preset {
	case Range1D:
		return now.AddDate(0, 0, -1), now
	case Range5D:
		return now.AddDate(0, 0, -5), now
	case Range1M:
		return now.AddDate(0, -1, 0), now
	case Range3M:
		return now.AddDate(0, -3, 0), now
	case RangeYTD:
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc), now
	case Range1Y:
		return now.AddDate(-1, 0, 0), now
	case Range5Y:
		return now.AddDate(-5, 0, 0), now
	case RangeCustom:
		if r.Custom.IsZero() {
			panic("mujica: custom MarketRange without dates")
		}
		s, e := r.Custom.Start, r.Custom.End
		return time.Date(s.Year, s.Month, s.Day, 0, 0, 0, 0, loc), time.Date(e.Year, e.Month, e.Day, 0, 0, 0, 0, loc).AddDate(0, 0, 1)
	}
	panic(fmt.Sprintf("mujica: unknown RangePreset %d", r.Preset))
}

// TimeRangeSelectorOptions configure a TimeRangeSelector.
type TimeRangeSelectorOptions struct {
	Disabled bool
	// Today marks the calendar's today; zero uses the clock.
	Today Date
}

// TimeRangeSelector chooses a preset range or a custom one; the presets are one Tab stop moved with the arrows.
func TimeRangeSelector(c *ui.Context, value *MarketRange, opts TimeRangeSelectorOptions) ChoiceResult {
	if value == nil {
		panic("mujica: TimeRangeSelector needs a value")
	}
	if value.Preset < Range1D || value.Preset > RangeCustom {
		panic(fmt.Sprintf("mujica: unknown RangePreset %d", value.Preset))
	}
	if value.Preset == RangeCustom && value.Custom.IsZero() {
		panic("mujica: custom MarketRange without dates")
	}
	items := make([]ToggleItem[RangePreset], len(rangePresetNames)+1)
	for i, n := range rangePresetNames {
		items[i] = ToggleItem[RangePreset]{Value: RangePreset(i), Label: n}
	}
	items[len(items)-1] = ToggleItem[RangePreset]{Value: RangeCustom, Label: msg(c, "time_range_selector.custom")}
	changed := false
	row := ui.Row(c).Gap(Density(c).Unit() * 2).AlignItems(ui.Center).Wrap().Disabled(opts.Disabled)
	row.Children(func() {
		sel := []RangePreset{value.Preset}
		if ToggleGroup(c, &sel, items, ToggleGroupOptions{Label: msg(c, "time_range_selector.label")}).Changed() && !row.IsDisabled() {
			if sel[0] == RangeCustom && value.Custom.IsZero() {
				today := datetimeToday(c, opts.Today)
				value.Custom = DateRange{Start: today.AddMonths(-1), End: today}
			}
			value.Preset, changed = sel[0], true
		}
		if value.Preset == RangeCustom {
			if DateRangePicker(c, &value.Custom, DateRangePickerOptions{Today: opts.Today, Max: datetimeToday(c, opts.Today)}).Changed() {
				changed = true
			}
		}
	})
	if changed {
		slog.Debug("mujica: time_range_selector changed", "preset", value.Preset, "custom", value.Custom)
	}
	return ChoiceResult{row, changed}
}
