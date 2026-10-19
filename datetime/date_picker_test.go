package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func datePickerTester(value *core.Date, opts DatePickerOptions, changes *int) *ui.Tester {
	opts.Today = calendarToday
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			if DatePicker(c, value, opts).Changed() {
				*changes++
			}
			input.Button(c, "Elsewhere", input.ButtonOptions{})
		})
	}, 420, 460)
}

func datePickerType(tt *ui.Tester, label, text string) {
	_ = tt.Click(label)
	tt.Command("selectAll")
	tt.Type(text)
}

func TestDatePickerTyping(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due"}, &n)
	datePickerType(tt, "Due", "2026-11-03")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.November, Day: 3}) || n != 1 {
		t.Fatalf("Enter committed %v (%d)", v, n)
	}
	datePickerType(tt, "Due", "2026-12-24")
	if err := tt.Click("Elsewhere"); err != nil {
		t.Fatal(err)
	}
	if v != (core.Date{Year: 2026, Month: time.December, Day: 24}) || n != 2 {
		t.Fatalf("blur committed %v (%d)", v, n)
	}
}

func TestDatePickerInvalidInputKept(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due", Max: core.Date{Year: 2026, Month: time.December, Day: 31}}, &n)
	datePickerType(tt, "Due", "2026-02-30")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.October, Day: 5}) || n != 0 {
		t.Fatalf("invalid day was corrected into %v", v)
	}
	if !tt.HasText("Enter a date like 2026-10-05.") {
		t.Fatalf("expected an error: %v", tt.Texts())
	}
	// The typed text stays: fixing its day makes it valid.
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	tt.Type("28")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.February, Day: 28}) || n != 1 || tt.HasText("Enter a date like") {
		t.Fatalf("edited text committed %v (%d)", v, n)
	}
	datePickerType(tt, "Due", "2027-01-02")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.February, Day: 28}) || !tt.HasText("Choose 2026-12-31 or earlier.") {
		t.Fatalf("out of range accepted: %v %v", v, tt.Texts())
	}
	datePickerType(tt, "Due", "")
	tt.Key(0, ui.KeyEnter)
	if v.IsZero() {
		t.Fatal("empty text cleared a non-clearable picker")
	}
}

func TestDatePickerCalendarPopup(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due"}, &n)
	if err := tt.Click("Open calendar"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Monday, October 5, 2026") {
		t.Fatal("the popup should focus the chosen day")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.October, Day: 6}) || n != 1 || tt.HasText("October 2026") {
		t.Fatalf("keyboard pick %v, popup open %v", v, tt.HasText("October 2026"))
	}
	_ = tt.Click("Due")
	tt.Key(0, ui.KeyDown)
	if !tt.HasText("October 2026") {
		t.Fatal("Down should open the calendar")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("October 2026") || v != (core.Date{Year: 2026, Month: time.October, Day: 6}) {
		t.Fatal("Escape should close without change")
	}
	_ = tt.Click("Open calendar")
	if err := tt.Click("Thursday, October 22, 2026"); err != nil {
		t.Fatal(err)
	}
	if v != (core.Date{Year: 2026, Month: time.October, Day: 22}) || n != 2 {
		t.Fatalf("click pick %v", v)
	}
}

func TestDatePickerClearAndLayout(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due", Clearable: true, Layout: "02/01/2006"}, &n)
	_ = tt.Click("Due")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.October, Day: 5}) || n != 0 || tt.HasText("Enter a date like") {
		t.Fatalf("shown text did not parse back with the layout: %v", tt.Texts())
	}
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if !v.IsZero() || n != 1 {
		t.Fatalf("clear left %v", v)
	}
	datePickerType(tt, "Due", "31/12/2026")
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2026, Month: time.December, Day: 31}) {
		t.Fatalf("custom layout parsed %v", v)
	}
}

func TestDatePickerDisabled(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 5}, 0
	tt := datePickerTester(&v, DatePickerOptions{Label: "Due", Disabled: true}, &n)
	_ = tt.Click("Open calendar")
	if tt.HasText("October 2026") {
		t.Fatal("disabled picker opened")
	}
}

func TestDatePickerSnapshot(t *testing.T) {
	a, b, e := core.Date{Year: 2026, Month: time.October, Day: 14}, core.Date{}, core.Date{Year: 2026, Month: time.October, Day: 5}
	snapshot(t, "date_picker", 640, 130, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Gap(12).Width(260).Children(func() {
				DatePicker(c, &b, DatePickerOptions{Today: calendarToday})
				DatePicker(c, &e, DatePickerOptions{Today: calendarToday, Disabled: true})
			})
			ui.Column(c).Width(300).Children(func() {
				DatePicker(c, &a, DatePickerOptions{Today: calendarToday, Clearable: true, Label: "Due"})
			})
		})
	})
}
