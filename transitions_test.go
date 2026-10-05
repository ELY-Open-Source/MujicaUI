package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestPresenceEntersExitsAndRemoves(t *testing.T) {
	for _, kind := range []PresenceKind{PresenceFade, PresenceSlide, PresenceCollapse} {
		show := false
		tt, n := motionFrames(func(c *ui.Context) {
			ui.Column(c).Children(func() {
				Presence(c, show, PresenceOptions{Kind: kind}, func() { ui.Text(c, "Notice") })
				ui.Text(c, "Below")
			})
		}, 300, 120)
		if tt.HasText("Notice") {
			t.Fatal("hidden content built")
		}
		show = true
		before := *n
		tt.Frame()
		if !tt.HasText("Notice") || *n-before < 3 {
			t.Fatalf("kind %d: no entry motion (%d frames)", kind, *n-before)
		}
		time.Sleep(250 * time.Millisecond)
		tt.Frame()
		show = false
		before = *n
		tt.Frame()
		if *n-before < 3 {
			t.Fatalf("kind %d: no exit motion", kind)
		}
		time.Sleep(250 * time.Millisecond)
		tt.Frame()
		if tt.HasText("Notice") || motionBusy(tt, n) {
			t.Fatalf("kind %d: not removed after exit", kind)
		}
	}
}

func TestPresenceReducedMotionAndInvalid(t *testing.T) {
	show := false
	tt, n := motionFrames(func(c *ui.Context) {
		ui.Column(c).Children(func() { Presence(c, show, PresenceOptions{Kind: PresenceSlide}, func() { ui.Text(c, "Notice") }) })
	}, 300, 120)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	show = true
	if motionBusy(tt, n) || !tt.HasText("Notice") {
		t.Fatal("reduced motion animated")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown kind did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Use(c, Settings{}); Presence(c, true, PresenceOptions{Kind: 9}, nil) }, 100, 50, 1)
}

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

func TestReorderTransitionSlides(t *testing.T) {
	keys := []string{"a", "b", "c"}
	tt, n := motionFrames(func(c *ui.Context) {
		ReorderTransition(c, keys, ReorderTransitionOptions{Gap: 4}, func(k string) {
			ui.Text(c, "Item "+k).Padding(8)
		})
	}, 300, 200)
	ra, _ := tt.Find("Item a")
	keys = []string{"c", "b", "a"}
	before := *n
	tt.Frame()
	if *n-before < 3 {
		t.Fatal("reorder did not animate")
	}
	mid, _ := tt.Find("Item a")
	time.Sleep(260 * time.Millisecond)
	tt.Frame()
	end, _ := tt.Find("Item a")
	if !(mid.Y > ra.Y && mid.Y < end.Y) || motionBusy(tt, n) {
		t.Fatalf("a moved %v → %v → %v", ra.Y, mid.Y, end.Y)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate key did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		ReorderTransition(c, []int{1, 1}, ReorderTransitionOptions{}, func(int) {})
	}, 100, 50, 1)
}

func TestLayoutTransitionResizes(t *testing.T) {
	wide := false
	var w float32
	tt, n := motionFrames(func(c *ui.Context) {
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			box := LayoutTransition(c, LayoutTransitionOptions{Clip: true}, func() {
				width := float32(100)
				if wide {
					width = 300
				}
				ui.Box(c).Size(width, 40)
			})
			w = box.Bounds().W
		})
	}, 400, 120)
	wide = true
	before := *n
	tt.Frame()
	if *n-before < 3 {
		t.Fatal("resize did not animate")
	}
	time.Sleep(260 * time.Millisecond)
	tt.Frame()
	if motionBusy(tt, n) || w != 300 {
		t.Fatalf("settled width %v", w)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative duration did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) { Use(c, Settings{}); LayoutTransition(c, LayoutTransitionOptions{Duration: -1}, nil) }, 100, 50, 1)
}

func TestPresenceSnapshot(t *testing.T) {
	snapshot(t, "presence", 300, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Presence(c, true, PresenceOptions{}, func() { ui.Text(c, "The doors open.") })
		})
	})
}

func TestStaggerSnapshot(t *testing.T) {
	snapshot(t, "stagger", 300, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Stagger(c, 3, StaggerOptions{}, func(i int) { ui.Textf(c, "Guest %d arrives", i+1) })
		})
	})
}

func TestReorderTransitionSnapshot(t *testing.T) {
	snapshot(t, "reorder_transition", 300, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ReorderTransition(c, []string{"Sakiko", "Uika", "Mutsumi"}, ReorderTransitionOptions{Gap: 4}, func(k string) {
				ui.Text(c, k).Padding(4, 8).Background(Tokens(c).Surface).Radius(4)
			})
		})
	})
}

func TestLayoutTransitionSnapshot(t *testing.T) {
	snapshot(t, "layout_transition", 300, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
			LayoutTransition(c, LayoutTransitionOptions{}, func() {
				ui.Text(c, "A panel that grows").Padding(12).Background(Tokens(c).Surface).Radius(4)
			})
		})
	})
}
