package chart

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestTimeRangeSelectorPresetsCustomAndDisabled(t *testing.T) {
	value := MarketRange{}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TimeRangeSelector(c, &value, TimeRangeSelectorOptions{Disabled: disabled, Today: core.Date{Year: 2026, Month: time.October, Day: 5}})
	}, 850, 160)
	tt.Key(0, ui.KeyTab)
	for _, want := range []RangePreset{Range5D, Range1M, Range3M, RangeYTD, Range1Y, Range5Y, RangeCustom} {
		tt.Key(0, ui.KeyRight)
		tt.Key(0, ui.KeySpace)
		if value.Preset != want {
			t.Fatalf("preset %v", value.Preset)
		}
	}
	if value.Custom.Start != (core.Date{Year: 2026, Month: time.September, Day: 5}) || value.Custom.End != (core.Date{Year: 2026, Month: time.October, Day: 5}) {
		t.Fatalf("custom %+v", value.Custom)
	}
	now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
	for p := Range1D; p <= RangeCustom; p++ {
		r := value
		r.Preset = p
		from, to := r.Bounds(now, time.UTC)
		if !from.Before(to) {
			t.Fatalf("range %v: %v %v", p, from, to)
		}
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	if value.Preset != RangeCustom {
		t.Fatal("disabled preset changed")
	}
}
func TestTimeRangeSelectorSnapshot(t *testing.T) {
	value := MarketRange{Preset: Range1M}
	snapshot(t, "time_range_selector", 680, 90, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() { TimeRangeSelector(c, &value, TimeRangeSelectorOptions{}) })
	})
}
