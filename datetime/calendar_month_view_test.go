package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type monthProbe struct {
	month    YearMonth
	events   []CalendarEvent
	opts     CalendarMonthViewOptions
	opened   []string
	changes  []CalendarChange
	days     []core.Date
	disabled bool
}

func (p *monthProbe) tester() *ui.Tester {
	if p.month.IsZero() {
		p.month = YearMonth{2026, time.October}
	}
	p.opts.Location, p.opts.Today = calTokyo, calendarToday
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := p.opts
		o.Disabled = p.disabled
		r := CalendarMonthView(c, &p.month, p.events, o)
		if id, ok := r.Opened(); ok {
			p.opened = append(p.opened, id)
		}
		if ch, ok := r.Changed(); ok {
			p.changes = append(p.changes, ch)
		}
		if d, ok := r.DayChosen(); ok {
			p.days = append(p.days, d)
		}
	}, 1000, 900)
}

func TestCalendarMonthViewOverflowAndPopover(t *testing.T) {
	p := &monthProbe{}
	for i := range 5 {
		p.events = append(p.events, calEvent(fmt.Sprint(i), 14, 8+2*i, 0, 60))
	}
	tt := p.tester()
	if !tt.HasText("+2") || tt.HasText("Event 3") {
		t.Fatalf("3 rows then +2 expected: %v", tt.Texts())
	}
	_ = tt.Click("2 more events")
	if !tt.HasText("Event 4") {
		t.Fatal("+N should list the hidden events")
	}
	_ = tt.Click("Event 4, 16:00 – 17:00")
	if len(p.opened) != 1 || p.opened[0] != "4" {
		t.Fatalf("opened %v", p.opened)
	}
}

func TestCalendarMonthViewOverlapSideBySide(t *testing.T) {
	p := &monthProbe{events: []CalendarEvent{calEvent("a", 14, 9, 0, 120), calEvent("b", 14, 10, 0, 60), calEvent("c", 14, 15, 0, 60)}}
	tt := p.tester()
	a, _ := tt.Find("Event a, 09:00 – 11:00")
	b, _ := tt.Find("Event b, 10:00 – 11:00")
	cc, _ := tt.Find("Event c, 15:00 – 16:00")
	if a.Y != b.Y || b.X <= a.X || cc.Y <= a.Y {
		t.Fatalf("overlapping events should share a row: %v %v %v", a, b, cc)
	}
}

func TestCalendarMonthViewDragAndKeys(t *testing.T) {
	p := &monthProbe{events: []CalendarEvent{calEvent("a", 14, 9, 30, 60)}}
	tt := p.tester()
	chip, _ := tt.Find("Event a, 09:30 – 10:30")
	target, _ := tt.Find(datetimeTestLabel(oct(16)))
	tt.Press(chip.X+10, chip.Y+5)
	for i := 1; i <= 4; i++ {
		f := float32(i) / 4
		tt.Move(chip.X+10+(target.X+20-chip.X-10)*f, chip.Y+5+(target.Y+40-chip.Y-5)*f)
	}
	tt.Release(target.X+20, target.Y+40)
	want := time.Date(2026, time.October, 16, 9, 30, 0, 0, calTokyo)
	if len(p.changes) != 1 || !p.changes[0].Start.Equal(want) || p.changes[0].End.Sub(p.changes[0].Start) != time.Hour {
		t.Fatalf("drop gave %+v", p.changes)
	}
	_ = tt.Click("Event a, 09:30 – 10:30")
	tt.Key(ui.Alt, ui.KeyDown)
	if len(p.changes) != 2 || !p.changes[1].Start.Equal(time.Date(2026, time.October, 21, 9, 30, 0, 0, calTokyo)) {
		t.Fatalf("Alt+Down gave %+v", p.changes)
	}
}

func TestCalendarMonthViewDSTGapReported(t *testing.T) {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	ev := CalendarEvent{ID: "g", Title: "Early", Start: time.Date(2026, time.March, 7, 2, 30, 0, 0, ny)}
	ev.End = ev.Start.Add(time.Hour)
	p := &monthProbe{month: YearMonth{2026, time.March}, events: []CalendarEvent{ev}}
	p.opts.Location = ny
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ch, ok := CalendarMonthView(c, &p.month, p.events, CalendarMonthViewOptions{Location: ny, Today: core.Date{Year: 2026, Month: time.March, Day: 1}}).Changed(); ok {
			p.changes = append(p.changes, ch)
		}
	}, 1000, 900)
	_ = tt.Click("Early, 02:30 – 03:30")
	tt.Key(ui.Alt, ui.KeyRight)
	if len(p.changes) != 0 || !tt.HasText("2026-03-08 02:30 does not exist in America/New_York") {
		t.Fatalf("gap not reported: %v", tt.Texts())
	}
}

func TestCalendarMonthViewNavigationDayAndDisabled(t *testing.T) {
	p := &monthProbe{events: []CalendarEvent{calEvent("a", 14, 9, 0, 60)}}
	tt := p.tester()
	_ = tt.Click(datetimeTestLabel(oct(20)))
	_ = tt.Click("Next month")
	if len(p.days) != 1 || p.days[0] != oct(20) || p.month != (YearMonth{2026, time.November}) {
		t.Fatalf("day %v month %v", p.days, p.month)
	}
	p.month, p.disabled = YearMonth{2026, time.October}, true
	tt.Frame()
	_ = tt.Click("Next month")
	_ = tt.Click("Event a, 09:00 – 10:00")
	tt.Key(ui.Alt, ui.KeyRight)
	if p.month != (YearMonth{2026, time.October}) || len(p.opened) != 0 || len(p.changes) != 0 {
		t.Fatal("disabled view reacted")
	}
}

func TestCalendarMonthViewLargeDataAndInvalid(t *testing.T) {
	p := &monthProbe{}
	base := time.Date(2000, 1, 1, 9, 0, 0, 0, calTokyo)
	for i := range 100000 {
		s := base.Add(time.Duration(i) * 3 * time.Hour)
		p.events = append(p.events, CalendarEvent{ID: fmt.Sprint(i), Title: "E", Start: s, End: s.Add(time.Hour)})
	}
	tt := p.tester()
	if n := len(tt.Texts()); n > 800 {
		t.Fatalf("built %d labels for one month", n)
	}
	for _, bad := range []func(){
		func() { q := &monthProbe{events: []CalendarEvent{{ID: "x", Start: base, End: base}}}; q.tester() },
		func() { q := &monthProbe{opts: CalendarMonthViewOptions{MaxPerDay: -1}}; q.tester() },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			bad()
		}()
	}
}

func TestCalendarMonthViewSnapshot(t *testing.T) {
	m := YearMonth{2026, time.October}
	evs := []CalendarEvent{calEvent("1", 6, 9, 0, 60), calEvent("2", 14, 9, 0, 120), calEvent("3", 14, 10, 0, 30),
		calEvent("4", 14, 13, 0, 30), calEvent("5", 14, 15, 0, 30), calEvent("6", 14, 17, 0, 30), calEvent("7", 22, 18, 0, 60)}
	evs[1].Color, evs[2].Color, evs[6].Color = EventInfo, EventSuccess, EventWarning
	snapshot(t, "calendar_month_view", 900, 860, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CalendarMonthView(c, &m, evs, CalendarMonthViewOptions{Location: calTokyo, Today: calendarToday})
		})
	})
}
