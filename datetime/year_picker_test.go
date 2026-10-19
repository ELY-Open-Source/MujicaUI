package datetime

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func yearPickerTester(value *int, opts YearPickerOptions, changes *int) *ui.Tester {
	opts.Today = calendarToday
	opts.Label = "Vintage"
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Children(func() {
			if YearPicker(c, value, opts).Changed() {
				*changes++
			}
		})
	}, 420, 420)
}

func TestYearPickerClickAndDecade(t *testing.T) {
	v, n := 0, 0
	tt := yearPickerTester(&v, YearPickerOptions{}, &n)
	_ = tt.Click("Vintage")
	if !tt.HasText("2020 – 2029") || !tt.HasText("2019") || !tt.HasText("2030") {
		t.Fatalf("decade page: %v", tt.Texts())
	}
	_ = tt.Click("Previous decade")
	if err := tt.Click("2014"); err != nil {
		t.Fatal(err)
	}
	if v != 2014 || n != 1 || tt.HasText("2015") {
		t.Fatalf("picked %d", v)
	}
}

func TestYearPickerKeyboard(t *testing.T) {
	v, n := 2026, 0
	tt := yearPickerTester(&v, YearPickerOptions{}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("2026") {
		t.Fatal("popup should focus the chosen year")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("2030 – 2039") || !tt.Focused("2030") {
		t.Fatalf("crossing the decade should turn the page: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyPageUp)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeySpace)
	if v != 2017 || n != 1 {
		t.Fatalf("Space chose %d", v)
	}
}

func TestYearPickerLimitsAndClear(t *testing.T) {
	v, n := 2024, 0
	opts := YearPickerOptions{Min: 2022, Max: 2027, Clearable: true, DisableYear: func(y int) bool { return y == 2025 }}
	tt := yearPickerTester(&v, opts, &n)
	_ = tt.Click("Vintage")
	for _, y := range []string{"2021", "2025", "2028"} {
		_ = tt.Click(y)
	}
	if v != 2024 || n != 0 {
		t.Fatalf("disabled year chosen: %d", v)
	}
	_ = tt.Click("2026")
	_ = tt.Click("Clear")
	if v != 0 || n != 2 {
		t.Fatalf("clear left %d (%d)", v, n)
	}
}

func TestYearPickerSnapshot(t *testing.T) {
	a, b := 2026, 0
	snapshot(t, "year_picker", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			YearPicker(c, &b, YearPickerOptions{Today: calendarToday})
			YearPicker(c, &a, YearPickerOptions{Today: calendarToday, Clearable: true})
		})
	})
}
