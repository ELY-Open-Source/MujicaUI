package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

type pomodoroProbe struct {
	s        PomodoroState
	now      time.Time
	opts     PomodoroTimerOptions
	finished []PomodoroPhase
	changes  int
}

func (p *pomodoroProbe) tester() *ui.Tester {
	p.now = time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		o := p.opts
		o.Now = p.now
		r := PomodoroTimer(c, &p.s, o)
		if ph, ok := r.Finished(); ok {
			p.finished = append(p.finished, ph)
		}
		if r.Changed() {
			p.changes++
		}
	}, 400, 400)
}

func TestPomodoroTimerCycle(t *testing.T) {
	p := &pomodoroProbe{opts: PomodoroTimerOptions{Focus: 2 * time.Second, Break: time.Second, LongBreak: 3 * time.Second, LongEvery: 2}}
	tt := p.tester()
	if !tt.HasText("Focus, 2 s left") {
		t.Fatalf("start %v", tt.Texts())
	}
	_ = tt.Click("Start")
	p.now = p.now.Add(900 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Focus, 2 s left") {
		t.Fatal("seconds round up")
	}
	_ = tt.Click("Pause")
	p.now = p.now.Add(time.Hour)
	tt.Frame()
	if !tt.HasText("Focus, 2 s left") || p.s.Running() {
		t.Fatal("pause should hold the clock")
	}
	_ = tt.Click("Start")
	p.now = p.now.Add(1100 * time.Millisecond)
	tt.Frame()
	if len(p.finished) != 1 || p.finished[0] != PomodoroFocus || p.s.Phase() != PomodoroBreak || p.s.Running() || p.s.Sessions() != 1 {
		t.Fatalf("focus end: %v %v", p.finished, p.s.Phase())
	}
	if a := tt.Announcements(); len(a) == 0 || a[len(a)-1] != "Focus finished." {
		t.Fatalf("announce %q", a)
	}
	_ = tt.Click("Skip")
	_ = tt.Click("Skip")
	if p.s.Phase() != PomodoroLongBreak || p.s.Sessions() != 2 {
		t.Fatalf("second focus should earn a long break: %v", p.s.Phase())
	}
	_ = tt.Click("Reset")
	if p.s.Phase() != PomodoroFocus || p.s.Sessions() != 0 || p.changes != 6 {
		t.Fatalf("reset %v (%d)", p.s.Phase(), p.changes)
	}
}

func TestPomodoroTimerKeyboardDisabledInvalid(t *testing.T) {
	p := &pomodoroProbe{opts: PomodoroTimerOptions{Disabled: true}}
	tt := p.tester()
	_ = tt.Click("Start")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if p.s.Running() || p.changes != 0 || !tt.HasText("Focus, 25 min left") {
		t.Fatal("disabled timer started")
	}
	q := &pomodoroProbe{}
	tt = q.tester()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !q.s.Running() {
		t.Fatal("Enter should start")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	bad := &pomodoroProbe{opts: PomodoroTimerOptions{Focus: 1500 * time.Millisecond}}
	bad.tester()
}

func TestPomodoroTimerSnapshot(t *testing.T) {
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	s := PomodoroState{spent: 10 * time.Minute}
	snapshot(t, "pomodoro_timer", 360, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { PomodoroTimer(c, &s, PomodoroTimerOptions{Now: now}) })
	})
}

func TestPomodoroTimerIdleAndHidden(t *testing.T) {
	s, shown, builds := PomodoroState{}, true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		builds++
		Use(c, Settings{})
		if shown {
			PomodoroTimer(c, &s, PomodoroTimerOptions{})
		}
	}, 400, 400)
	if builds > 2 {
		t.Fatalf("idle built %d frames", builds)
	}
	if err := tt.Click("Start"); err != nil {
		t.Fatal(err)
	}
	shown = false
	tt.Frame()
	before := builds
	time.Sleep(1100 * time.Millisecond)
	if builds != before || !s.Running() {
		t.Fatalf("hidden timer: builds %d, running %v", builds, s.Running())
	}
	shown = true
	tt.Frame()
	if !tt.HasText("Pause") {
		t.Fatal("visible timer did not retain its clock")
	}
}
