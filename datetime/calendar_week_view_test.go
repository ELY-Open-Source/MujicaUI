package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type weekProbe struct {
	week     core.Date
	events   []CalendarEvent
	opts     CalendarWeekViewOptions
	changes  []CalendarChange
	opened   []string
	disabled bool
}

func (p *weekProbe) tester() *ui.Tester {
	if p.week.IsZero() {
		p.week = oct(14)
	}
	p.opts.Location, p.opts.Today = calTokyo, calendarToday
	p.opts.Now = time.Date(2026, time.October, 14, 10, 30, 0, 0, calTokyo)
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := p.opts
		o.Disabled = p.disabled
		r := CalendarWeekView(c, &p.week, p.events, o)
		if ch, ok := r.Changed(); ok {
			p.changes = append(p.changes, ch)
		}
		if id, ok := r.Opened(); ok {
			p.opened = append(p.opened, id)
		}
	}, 1100, 900)
}

func calDrag(tt *ui.Tester, x0, y0, x1, y1 float32) {
	tt.Press(x0, y0)
	for i := 1; i <= 4; i++ {
		f := float32(i) / 4
		tt.Move(x0+(x1-x0)*f, y0+(y1-y0)*f)
	}
	tt.Release(x1, y1)
}

func calAt(t *testing.T, h, m int, day int) time.Time {
	t.Helper()
	return time.Date(2026, time.October, day, h, m, 0, 0, calTokyo)
}

func TestCalendarWeekViewDragCreateAndKeyboardCreate(t *testing.T) {
	p := &weekProbe{}
	tt := p.tester()
	col, ok := tt.Find("Schedule, " + datetimeTestLabel(oct(15)))
	if !ok {
		t.Fatal("no column")
	}
	hourH := float32(72)
	// Drag 9:00 to 10:30; the visible top is 8:00.
	y9 := col.Y + hourH + 5
	calDrag(tt, col.X+20, y9, col.X+20, y9+hourH*1.5)
	if len(p.changes) != 1 || p.changes[0].Kind != CalendarCreate || !p.changes[0].Start.Equal(calAt(t, 9, 0, 15)) || !p.changes[0].End.Equal(calAt(t, 11, 0, 15)) {
		t.Fatalf("drag create %+v", p.changes)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.changes) != 2 || !p.changes[1].Start.Equal(calAt(t, 9, 30, 15)) || !p.changes[1].End.Equal(calAt(t, 10, 0, 15)) {
		t.Fatalf("keyboard create %+v", p.changes)
	}
}

func TestCalendarWeekViewMoveResizeAndOverlap(t *testing.T) {
	p := &weekProbe{events: []CalendarEvent{calEvent("a", 14, 9, 0, 60), calEvent("b", 14, 9, 30, 60)}}
	tt := p.tester()
	a, _ := tt.Find("Event a, 09:00 – 10:00")
	b, _ := tt.Find("Event b, 09:30 – 10:30")
	if b.X <= a.X || b.X-a.X < a.W/2 || b.Y <= a.Y {
		t.Fatalf("overlapping events should sit side by side: %v %v", a, b)
	}
	calDrag(tt, a.X+10, a.Y+20, a.X+10, a.Y+20+72)
	if len(p.changes) != 1 || p.changes[0].Kind != CalendarMove || !p.changes[0].Start.Equal(calAt(t, 10, 0, 14)) || !p.changes[0].End.Equal(calAt(t, 11, 0, 14)) {
		t.Fatalf("move %+v", p.changes)
	}
	h, _ := tt.Find("Resize Event b")
	calDrag(tt, h.X+10, h.Y+3, h.X+10, h.Y+3+36)
	if len(p.changes) != 2 || p.changes[1].Kind != CalendarResize || !p.changes[1].End.Equal(calAt(t, 11, 0, 14)) {
		t.Fatalf("resize %+v", p.changes)
	}
	_ = tt.Click("Event a, 09:00 – 10:00")
	if len(p.opened) != 1 {
		t.Fatalf("click should open: %v", p.opened)
	}
	tt.Key(ui.Alt, ui.KeyRight)
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt|ui.Shift, ui.KeyUp)
	if len(p.changes) != 5 || !p.changes[2].Start.Equal(calAt(t, 9, 0, 15)) || !p.changes[3].Start.Equal(calAt(t, 9, 30, 14)) ||
		!p.changes[4].End.Equal(calAt(t, 9, 30, 14)) {
		t.Fatalf("keyboard changes %+v", p.changes[2:])
	}
}

func TestCalendarWeekViewAllDayWeekStartAndNav(t *testing.T) {
	ev := calEvent("h", 16, 0, 0, 24*60)
	ev.AllDay, ev.Title = true, "Holiday"
	p := &weekProbe{events: []CalendarEvent{ev}, opts: CalendarWeekViewOptions{WeekStart: WeekStartMonday}}
	tt := p.tester()
	if !tt.HasText("October 12 – October 18, 2026") || !tt.HasText("All day") || !tt.HasText("Holiday") {
		t.Fatalf("texts %v", tt.Texts()[:12])
	}
	_ = tt.Click("Holiday, All day")
	tt.Key(ui.Alt, ui.KeyRight)
	if len(p.changes) != 1 || !p.changes[0].Start.Equal(calAt(t, 0, 0, 17)) || !p.changes[0].AllDay {
		t.Fatalf("all-day move %+v", p.changes)
	}
	_ = tt.Click("Next week")
	if p.week != oct(21) {
		t.Fatalf("week %v", p.week)
	}
}

func TestCalendarWeekViewDisabledAndInvalid(t *testing.T) {
	p := &weekProbe{events: []CalendarEvent{calEvent("a", 14, 9, 0, 60)}, disabled: true}
	tt := p.tester()
	a, _ := tt.Find("Event a, 09:00 – 10:00")
	calDrag(tt, a.X+10, a.Y+20, a.X+10, a.Y+92)
	_ = tt.Click("Event a, 09:00 – 10:00")
	tt.Key(ui.Alt, ui.KeyDown)
	_ = tt.Click("Next week")
	if len(p.changes) != 0 || len(p.opened) != 0 || p.week != oct(14) {
		t.Fatal("disabled week view reacted")
	}
	for _, opts := range []CalendarWeekViewOptions{{SlotMinutes: 7}, {SlotMinutes: -15}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: expected panic", opts)
				}
			}()
			q := &weekProbe{opts: opts}
			q.tester()
		}()
	}
}

func TestCalendarWeekViewLargeData(t *testing.T) {
	p := &weekProbe{}
	base := time.Date(2000, 1, 1, 9, 0, 0, 0, calTokyo)
	for i := range 100000 {
		s := base.Add(time.Duration(i) * 3 * time.Hour)
		p.events = append(p.events, CalendarEvent{ID: fmt.Sprint(i), Title: "E", Start: s, End: s.Add(time.Hour)})
	}
	tt := p.tester()
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("built %d labels for one week", n)
	}
}

func TestCalendarWeekViewSnapshot(t *testing.T) {
	w := oct(14)
	evs := []CalendarEvent{calEvent("1", 12, 9, 0, 60), calEvent("2", 14, 9, 0, 120), calEvent("3", 14, 9, 30, 60),
		calEvent("4", 15, 13, 0, 90), calEvent("5", 16, 10, 0, 45)}
	evs[1].Color, evs[2].Color, evs[3].Color = EventInfo, EventSuccess, EventWarning
	h := calEvent("h", 13, 0, 0, 24*60)
	h.AllDay, h.Title = true, "Holiday"
	evs = append(evs, h)
	snapshot(t, "calendar_week_view", 1000, 860, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CalendarWeekView(c, &w, evs, CalendarWeekViewOptions{Location: calTokyo, Today: calendarToday,
				Now: time.Date(2026, time.October, 14, 10, 30, 0, 0, calTokyo)})
		})
	})
}
