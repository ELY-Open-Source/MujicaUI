package datetime

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// datePickerLive builds a DatePicker whose options a test changes between frames.
func datePickerLive(value *core.Date, opts *DatePickerOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			o := *opts
			o.Today, o.Label = calendarToday, "Due"
			if DatePicker(c, value, o).Changed() {
				*changes++
			}
		})
	}, 420, 460)
}

func TestDatePickerDisablingClosesPopup(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	opts := DatePickerOptions{}
	tt := datePickerLive(&v, &opts, &n)
	_ = tt.Click("Open calendar")
	day := datetimeTestLabel(oct(5))
	if !tt.Focused(day) {
		t.Fatal("setup: the popup should focus the chosen day")
	}
	opts.Disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click(datetimeTestLabel(oct(20)))
	if v != oct(5) || n != 0 {
		t.Fatalf("disabled picker committed %v (%d events)", v, n)
	}
	if tt.HasText("October 2026") || tt.Focused(datetimeTestLabel(oct(6))) || tt.Focused("Due") {
		t.Fatal("the popup should close and nothing in the disabled picker keep focus")
	}
}

func TestDatePickerDisabledBlocksInput(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due", Disabled: true}, &n)
	_ = tt.Click("Due")
	tt.Type("2027-01-01")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	_ = tt.Click("Open calendar")
	if v != oct(5) || n != 0 || tt.Focused("Due") || tt.HasText("October 2026") {
		t.Fatalf("disabled picker took input: %v (%d)", v, n)
	}
}

func TestTimePickerKeepsInvalidHour(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	tt := timePickerTester(&v, TimePickerOptions{}, &n)
	_ = tt.Click("Alarm: Hours")
	tt.Key(0, ui.Key2)
	tt.Key(0, ui.Key9)
	if v != (TimeOfDay{2, 30, 0}) || n != 1 {
		t.Fatalf("first digit should apply: %v", v)
	}
	if !tt.HasText("29") || !tt.HasText("Hours must be from 0 to 23.") || !tt.Focused("Alarm: Hours") {
		t.Fatalf("29 should stay shown with an error: %v", tt.Texts())
	}
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key5)
	if v != (TimeOfDay{15, 30, 0}) || n != 3 || tt.HasText("Hours must be") {
		t.Fatalf("fixing the hour: %v", v)
	}
}

func TestTimeSegmentDescription(t *testing.T) {
	var got []string
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		got = []string{timeSegmentDescription(c, timeHourPart, " 9", 1, 12), timeSegmentDescription(c, timePeriodPart, "PM", 0, 59)}
	}, 100, 100)
	if got[0] != "9, from 1 to 12" || got[1] != "PM" {
		t.Fatalf("descriptions %q", got)
	}
}

func TestTimePickerSegmentFocusRing(t *testing.T) {
	v := TimeOfDay{9, 30, 0}
	tt := timePickerTester(&v, TimePickerOptions{}, new(int))
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	r, ok := tt.Find("Alarm: Hours")
	if !ok || !tt.Focused("Alarm: Hours") {
		t.Fatal("setup: hours not focused")
	}
	want := theme.Light().Focus
	img := tt.Image()
	found := false
	for dy := float32(-4); dy <= -2 && !found; dy++ {
		c := img.RGBAAt(int(r.X+r.W/2), int(r.Y+dy))
		found = datetimeNear(c, want)
	}
	if !found {
		t.Fatal("no focus ring above the focused segment")
	}
}

func datetimeNear(c color.RGBA, w ui.Color) bool {
	d := func(a, b uint8) int {
		if a > b {
			return int(a - b)
		}
		return int(b - a)
	}
	return d(c.R, w.R)+d(c.G, w.G)+d(c.B, w.B) < 60
}

func TestDatePickerLayoutChangeReformats(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	opts := DatePickerOptions{}
	tt := datePickerLive(&v, &opts, &n)
	opts.Layout = "02/01/2006"
	tt.Frame()
	_ = tt.Click("Due")
	tt.Key(0, ui.KeyEnter)
	if v != oct(5) || n != 0 || tt.HasText("Enter a date like") {
		t.Fatalf("committed text kept the old layout: %v", tt.Texts())
	}
	_ = tt.Click("Due")
	tt.Command("selectAll")
	tt.Type("12/10")
	opts.Layout = "2006-01-02"
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Enter a date like 2026-10-05.") || v != oct(5) {
		t.Fatal("a draft being edited should be kept across a layout change")
	}
}

func TestYearPickerHomeEndStayOnPage(t *testing.T) {
	v := 2026
	tt := yearPickerTester(&v, YearPickerOptions{}, new(int))
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("2020 – 2029") || !tt.Focused("2020") {
		t.Fatalf("Home left the page: %v", tt.Texts()[:4])
	}
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("2020 – 2029") || !tt.Focused("2029") {
		t.Fatal("End left the page")
	}
}

func TestDatePickerReadOnly(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due", ReadOnly: true, Clearable: true}, &n)
	_ = tt.Click("Due")
	if !tt.Focused("Due") {
		t.Fatal("read-only field should take focus")
	}
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "2026-10-05" {
		t.Fatalf("copied %q", tt.Clipboard())
	}
	tt.Type("2027-01-01")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	if v != oct(5) || n != 0 || tt.HasText("October 2026") || tt.HasText("Open calendar") || tt.HasText("Clear") {
		t.Fatalf("read-only picker changed: %v %v", v, tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "2026-10-05" {
		t.Fatalf("typed text leaked into the field: %q", tt.Clipboard())
	}
}

func TestTimePickerReadOnly(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	tt := timePickerTester(&v, TimePickerOptions{ReadOnly: true}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.Key1)
	tt.Key(ui.Alt, ui.KeyDown)
	if !tt.Focused("Alarm: Hours") || v != (TimeOfDay{9, 30, 0}) || n != 0 || tt.HasText("Choose time") {
		t.Fatalf("read-only time changed: %v", v)
	}
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "09:30" {
		t.Fatalf("copied %q", tt.Clipboard())
	}
}

func TestMonthPickerReadOnly(t *testing.T) {
	v, n := YearMonth{2026, time.May}, 0
	tt := monthPickerTester(&v, MonthPickerOptions{ReadOnly: true, Clearable: true}, &n)
	_ = tt.Click("Billing month")
	tt.Key(0, ui.KeyDown)
	if tt.HasText("January") || tt.HasText("Clear") || !tt.Focused("Billing month") {
		t.Fatalf("read-only month opened: %v", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "May 2026" || v != (YearMonth{2026, time.May}) || n != 0 {
		t.Fatalf("copied %q, value %v", tt.Clipboard(), v)
	}
}

func TestDateTimePickerLayout(t *testing.T) {
	v := time.Time{}
	n, x := 0, 0
	maxT := time.Date(2026, time.October, 10, 12, 0, 0, 0, time.UTC)
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: time.UTC, Max: maxT, Layout: "02/01/2006"}, &n, &x)
	_ = tt.Click("Meeting")
	_ = tt.Click(datetimeTestLabel(oct(10)))
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key3)
	_ = tt.Click("Confirm")
	if n != 0 || !tt.HasText("Choose 10/10/2026 12:00 UTC or earlier.") {
		t.Fatalf("error ignores the layout: %v", tt.Texts())
	}
	_ = tt.Click(datetimeTestLabel(oct(9)))
	_ = tt.Click("Confirm")
	if n != 1 || !tt.HasText("09/10/2026 13:00 UTC") {
		t.Fatalf("display ignores the layout: %v", tt.Texts())
	}
}
