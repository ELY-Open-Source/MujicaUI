package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTypewriterRevealsFinishesAndSkips(t *testing.T) {
	text := "The court assembles at dusk."
	finished := 0
	tt, n := motionFrames(func(c *ui.Context) {
		if Typewriter(c, text, TypewriterOptions{Speed: 100, Skippable: true}).Finished() {
			finished++
		}
	}, 400, 120)
	if tt.HasText(text) && motionCount(tt, text) == 2 {
		t.Fatal("whole text shown at once")
	}
	time.Sleep(120 * time.Millisecond)
	tt.Frame()
	shown := 0
	for _, s := range tt.Texts() {
		if s != text && s != "Skip" && len(s) > shown {
			shown = len(s)
		}
	}
	if shown < 5 || shown >= len(text) {
		t.Fatalf("partial reveal %d of %d: %q", shown, len(text), tt.Texts())
	}
	if err := tt.Click("Skip"); err != nil {
		t.Fatal(err)
	}
	if motionCount(tt, text) != 2 || finished != 1 || tt.HasText("Skip") || motionBusy(tt, n) {
		t.Fatalf("skip: finished %d texts %q", finished, tt.Texts())
	}
	text += " The music begins."
	tt.Frame()
	if motionCount(tt, "The court assembles at dusk.") != 1 {
		t.Fatalf("appended text restarted: %q", tt.Texts())
	}
	time.Sleep(300 * time.Millisecond)
	tt.Frame()
	if motionCount(tt, text) != 2 || finished != 2 {
		t.Fatalf("append: finished %d %q", finished, tt.Texts())
	}
}

func TestTypewriterReducedMotionAndInvalid(t *testing.T) {
	finished := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		finished = Typewriter(c, "All at once", TypewriterOptions{}).Finished() || finished
	}, 300, 80)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt
	tt2 := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		finished = Typewriter(c, "Shown now", TypewriterOptions{}).Finished() || finished
	}, 300, 80)
	tt2.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionCount(tt2, "Shown now") != 2 || !finished {
		t.Fatalf("reduced motion: %q", tt2.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative speed did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Use(c, Settings{}); Typewriter(c, "x", TypewriterOptions{Speed: -1}) }, 100, 50, 1)
}

func TestTypewriterSnapshot(t *testing.T) {
	snapshot(t, "typewriter", 420, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Typewriter(c, "The court assembles at dusk; the music begins.", TypewriterOptions{Skippable: true})
		})
	})
}
