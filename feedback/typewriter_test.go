package feedback

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestTypewriterRevealsFinishesAndSkips(t *testing.T) {
	text := "The court assembles at dusk."
	finished := 0
	tt, n := motionFrames(func(c *ui.Context) {
		if Typewriter(c, text, TypewriterOptions{Speed: 100, Skippable: true}).Changed() {
			finished++
		}
	}, 400, 120)
	if tt.HasText(text) && motionCount(tt, text) == 2 {
		t.Fatal("whole text shown at once")
	}
	tt.Advance(120 * time.Millisecond)
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
	tt.Advance(300 * time.Millisecond)
	if motionCount(tt, text) != 2 || finished != 2 {
		t.Fatalf("append: finished %d %q", finished, tt.Texts())
	}
}

func TestTypewriterReducedMotionAndInvalid(t *testing.T) {
	finished := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finished = Typewriter(c, "All at once", TypewriterOptions{}).Changed() || finished
	}, 300, 80)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt
	tt2 := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		finished = Typewriter(c, "Shown now", TypewriterOptions{}).Changed() || finished
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
	ui.Render(func(c *ui.Context) { core.Use(c, core.Settings{}); Typewriter(c, "x", TypewriterOptions{Speed: -1}) }, 100, 50, 1)
}

func TestTypewriterSnapshot(t *testing.T) {
	snapshot(t, "typewriter", 420, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Typewriter(c, "The court assembles at dusk; the music begins.", TypewriterOptions{Skippable: true})
		})
	})
}

func TestTypewriterReplacementRestartsClock(t *testing.T) {
	text := "A long opening message"
	tt, _ := motionFrames(func(c *ui.Context) {
		Typewriter(c, text, TypewriterOptions{Speed: 10})
	}, 300, 100)
	time.Sleep(220 * time.Millisecond)
	text = "New letter"
	tt.Frame()
	if motionCount(tt, "Ne") > 0 {
		t.Fatalf("replacement consumed old elapsed time: %q", tt.Texts())
	}
}

func TestTypewriterInvalidSpeedBoundary(t *testing.T) {
	for _, speed := range []float64{nanValue(), math.Inf(1), 2e9} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("invalid speed %v accepted", speed)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Typewriter(c, "text", TypewriterOptions{Speed: speed})
			}, 200, 100, 1)
		}()
	}
}

func TestTypewriterMinimumIntervalBoundary(t *testing.T) {
	minimum := float64(time.Second) / float64(math.MaxInt64)
	tt, n := motionFrames(func(c *ui.Context) {
		Typewriter(c, "letter", TypewriterOptions{Speed: minimum})
	}, 200, 100)
	if motionCount(tt, "letter") != 1 || motionBusy(tt, n) {
		t.Fatal("slow typing should wait with partial text")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unrepresentable interval accepted")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Typewriter(c, "letter", TypewriterOptions{Speed: math.Nextafter(minimum, 0)})
	}, 200, 100, 1)
}
