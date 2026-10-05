package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestStaggerRevealsInTurnOnce(t *testing.T) {
	count := 5
	tt, n := motionFrames(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			Stagger(c, count, StaggerOptions{Interval: 50 * time.Millisecond}, func(i int) { ui.Textf(c, "Row %d", i) })
		})
	}, 300, 240)
	if !tt.HasText("Row 0") || tt.HasText("Row 4") {
		t.Fatalf("first frame: %q", tt.Texts())
	}
	time.Sleep(120 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Row 2") || tt.HasText("Row 4") {
		t.Fatalf("after 120 ms: %q", tt.Texts())
	}
	time.Sleep(300 * time.Millisecond)
	tt.Frame()
	time.Sleep(250 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Row 4") || motionBusy(tt, n) {
		t.Fatalf("did not finish and settle: %q", tt.Texts())
	}
	count = 6
	tt.Frame()
	if !tt.HasText("Row 5") {
		t.Fatal("item added after the sequence waited")
	}
}

func TestStaggerReducedMotionAndInvalid(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if c.Preferences().ReduceMotion {
			Stagger(c, 3, StaggerOptions{}, func(i int) { ui.Textf(c, "Row %d", i) })
		}
	}, 300, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("Row 2") {
		t.Fatalf("reduced motion staggered: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative n did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Use(c, Settings{}); Stagger(c, -1, StaggerOptions{}, nil) }, 100, 50, 1)
}

func TestStaggerSnapshot(t *testing.T) {
	snapshot(t, "stagger", 300, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Stagger(c, 3, StaggerOptions{}, func(i int) { ui.Textf(c, "Guest %d arrives", i+1) })
		})
	})
}

func TestStaggerEmptySettles(t *testing.T) {
	tt, n := motionFrames(func(c *ui.Context) { Stagger(c, 0, StaggerOptions{}, func(int) { t.Fatal("empty item built") }) }, 200, 80)
	if motionBusy(tt, n) {
		t.Fatal("empty stagger animated")
	}
}
