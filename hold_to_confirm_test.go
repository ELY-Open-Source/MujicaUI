package mujica

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHoldToConfirmStep(t *testing.T) {
	var s holdToConfirmState
	t0 := time.Unix(0, 0)
	d := time.Second
	if p, ok := s.step(true, t0, d); p != 0 || ok {
		t.Fatalf("start: %v %v", p, ok)
	}
	if p, ok := s.step(true, t0.Add(d/2), d); p != 0.5 || ok {
		t.Fatalf("half: %v %v", p, ok)
	}
	if _, ok := s.step(true, t0.Add(d), d); !ok {
		t.Fatal("full hold did not confirm")
	}
	if _, ok := s.step(true, t0.Add(2*d), d); ok {
		t.Fatal("one hold confirmed twice")
	}
	s.step(false, t0.Add(2*d), d)
	if p, _ := s.step(true, t0.Add(3*d), d); p != 0 {
		t.Fatalf("release did not restart: %v", p)
	}
}

func holdTester(dis *bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Padding(20).Children(func() {
			if HoldToConfirm(c, "Delete repo", HoldToConfirmOptions{Duration: 60 * time.Millisecond, Disabled: *dis}).Changed() {
				*n++
			}
		})
	}, 400, 200)
}

func TestHoldToConfirmPointerHold(t *testing.T) {
	dis, n := false, 0
	tt := holdTester(&dis, &n)
	r, _ := tt.Find("Delete repo")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	_ = tt.Click("Delete repo")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if n != 0 {
		t.Fatalf("a click or key tap confirmed: %d", n)
	}
	tt.Press(cx, cy)
	time.Sleep(20 * time.Millisecond)
	tt.Release(cx, cy)
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if n != 0 {
		t.Fatal("an early release confirmed")
	}
	tt.Press(cx, cy)
	time.Sleep(80 * time.Millisecond)
	tt.Frame()
	tt.Frame()
	tt.Release(cx, cy)
	if n != 1 {
		t.Fatalf("full hold confirmed %d times", n)
	}
	dis = true
	tt.Frame()
	tt.Press(cx, cy)
	time.Sleep(80 * time.Millisecond)
	tt.Frame()
	tt.Release(cx, cy)
	if n != 1 {
		t.Fatal("disabled hold confirmed")
	}
}

func TestHoldToConfirmSpaceHold(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		// A nanosecond hold completes in the frames between Space down and up.
		if HoldToConfirm(c, "Drop table", HoldToConfirmOptions{Duration: time.Nanosecond}).Changed() {
			n++
		}
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 0 {
		t.Fatal("Enter confirmed")
	}
	tt.Key(0, ui.KeySpace)
	if n != 1 {
		t.Fatalf("Space hold confirmed %d times", n)
	}
}

func TestHoldToConfirmReducedMotionCountdown(t *testing.T) {
	dis, n := false, 0
	tt := holdTester(&dis, &n)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	r, _ := tt.Find("Delete repo")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.HasPrefix(s, "Keep holding… 0.")
	}
	if !found {
		t.Fatalf("no countdown text: %q", tt.Texts())
	}
}

func TestHoldToConfirmSnapshot(t *testing.T) {
	snapshot(t, "hold_to_confirm", 400, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			HoldToConfirm(c, "Delete project", HoldToConfirmOptions{})
			HoldToConfirm(c, "Wipe", HoldToConfirmOptions{Disabled: true})
		})
	})
}
