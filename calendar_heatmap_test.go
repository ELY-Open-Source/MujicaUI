package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func calendarValues() map[Date]float64 {
	m := map[Date]float64{}
	d := Date{Year: 2026, Month: time.January, Day: 1}
	for i := range 365 {
		if i%3 == 0 {
			m[d.AddDays(i)] = float64(i % 11)
		}
	}
	return m
}

func TestCalendarHeatmapKeysAndYear(t *testing.T) {
	year := 2026
	values := 0
	start := WeekStartMonday
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := CalendarHeatmap(c, &year, calendarValues(), CalendarHeatmapOptions{Frame: ChartFrame{Label: "Prayers"}, WeekStart: start})
		if r.ValueChanged() {
			values++
		}
		return r
	}}
	tt := ui.NewTester(h.view, 560, 420)
	tabTo(t, tt, "Prayers")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "2026-01-01 · 0" {
		t.Fatalf("first day %q", got)
	}
	tt.Key(0, ui.KeyRight) // a week later
	tt.Key(0, ui.KeyDown)  // a day later
	if got := chartReadout(tt); got != "2026-01-09 · no value" {
		t.Fatalf("week and day %q", got)
	}
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "2026-12-31 · no value" {
		t.Fatalf("End %q", got)
	}
	if err := tt.Click("Next year"); err != nil {
		t.Fatal(err)
	}
	if year != 2027 || values != 1 || !tt.HasText("2027") {
		t.Fatalf("year %d events %d", year, values)
	}
}

func TestCalendarHeatmapWeekStartHover(t *testing.T) {
	year := 2026 // 1 January 2026 is a Thursday
	for _, tc := range []struct {
		start WeekStart
		row   int
	}{{WeekStartSunday, 4}, {WeekStartMonday, 3}} {
		h := &chartHarness{build: func(c *ui.Context) ChartResult {
			return CalendarHeatmap(c, &year, calendarValues(), CalendarHeatmapOptions{Frame: ChartFrame{Label: "Prayers"}, WeekStart: tc.start})
		}}
		tt := ui.NewTester(h.view, 560, 420)
		b := h.res.Element.Bounds()
		plotY := b.Y + b.H - 20 - 8 - 240 // above the readout line and gap
		step := min((b.W-36)/53, float32(240-48)/7)
		tt.Move(b.X+32+step/2, plotY+18+float32(tc.row)*step+step/2)
		if got := chartReadout(tt); got != "2026-01-01 · 0" {
			t.Fatalf("start %d: row %d reads %q", tc.start, tc.row, got)
		}
	}
}

func TestCalendarHeatmapBadYear(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	year := 0
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		CalendarHeatmap(c, &year, nil, CalendarHeatmapOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestCalendarHeatmapSnapshot(t *testing.T) {
	snapshot(t, "calendar_heatmap", 600, 260, func(c *ui.Context) {
		year := 2026
		ui.Column(c).Padding(16).Children(func() {
			CalendarHeatmap(c, &year, calendarValues(), CalendarHeatmapOptions{Frame: ChartFrame{Label: "Prayers", Height: 150}, WeekStart: WeekStartMonday})
		})
	})
}
