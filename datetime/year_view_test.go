package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func yearViewTester(year *int, events []CalendarEvent, opts YearViewOptions, days *[]core.Date) *ui.Tester {
	opts.Location, opts.Today = calTokyo, calendarToday
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		if d, ok := YearView(c, year, events, opts).DayChosen(); ok {
			*days = append(*days, d)
		}
	}, 900, 1200)
}

func TestYearViewDensityMarks(t *testing.T) {
	var evs []CalendarEvent
	for i := range 4 {
		evs = append(evs, calEvent(fmt.Sprint(i), 14, 9+i, 0, 30))
	}
	evs = append(evs, calEvent("x", 20, 9, 0, 30))
	year := 2026
	var days []core.Date
	tt := yearViewTester(&year, evs, YearViewOptions{}, &days)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("Wednesday, October 14, 2026, 4 events") || !tt.HasText("Tuesday, October 20, 2026, 1 event") {
		t.Fatalf("labels lack counts: %v", tt.Texts()[:5])
	}
	img := tt.Image()
	dotsAt := func(label string) int {
		r, _ := tt.Find(label)
		n, prev := 0, false
		for x := int(r.X); x < int(r.X+r.W); x++ {
			on := datetimeNear(img.RGBAAt(x, int(r.Y+r.H-5)), theme.Light().AccentText)
			if on && !prev {
				n++
			}
			prev = on
		}
		return n
	}
	if a, b, z := dotsAt("Wednesday, October 14, 2026, 4 events"), dotsAt("Tuesday, October 20, 2026, 1 event"), dotsAt(datetimeTestLabel(oct(21))); a != 2 || b != 1 || z != 0 {
		t.Fatalf("dots %d %d %d", a, b, z)
	}
	if yearDots(6) != 3 {
		t.Fatal("6 events should show three dots")
	}
}

func TestYearViewClickKeysAndYearChange(t *testing.T) {
	year := 2026
	var days []core.Date
	tt := yearViewTester(&year, nil, YearViewOptions{}, &days)
	_ = tt.Click(datetimeTestLabel(core.Date{Year: 2026, Month: time.March, Day: 3}))
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(days) != 2 || days[0] != (core.Date{Year: 2026, Month: time.March, Day: 3}) || days[1] != (core.Date{Year: 2026, Month: time.March, Day: 11}) {
		t.Fatalf("days %v", days)
	}
	_ = tt.Click(datetimeTestLabel(core.Date{Year: 2026, Month: time.December, Day: 31}))
	tt.Key(0, ui.KeyRight)
	if year != 2027 || !tt.Focused(datetimeTestLabel(core.Date{Year: 2027, Month: time.January, Day: 1})) {
		t.Fatalf("crossing the year: %d", year)
	}
	_ = tt.Click("Previous year")
	if year != 2026 {
		t.Fatalf("year %d", year)
	}
}

func TestYearViewDisabledAndInvalid(t *testing.T) {
	year := 2026
	var days []core.Date
	tt := yearViewTester(&year, nil, YearViewOptions{Disabled: true}, &days)
	_ = tt.Click(datetimeTestLabel(oct(3)))
	_ = tt.Click("Next year")
	if len(days) != 0 || year != 2026 {
		t.Fatal("disabled year view reacted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for year 0")
		}
	}()
	zero := 0
	yearViewTester(&zero, nil, YearViewOptions{}, &days)
}

func TestYearViewSnapshot(t *testing.T) {
	var evs []CalendarEvent
	for i := range 40 {
		evs = append(evs, calEvent(fmt.Sprint(i), 1+i%28, 9+i%5, 0, 30))
	}
	year := 2026
	snapshot(t, "year_view", 680, 960, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			YearView(c, &year, evs, YearViewOptions{Location: calTokyo, Today: calendarToday})
		})
	})
}
