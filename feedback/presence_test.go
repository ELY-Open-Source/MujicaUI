package feedback

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
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
		tt.Advance(250 * time.Millisecond)
		show = false
		before = *n
		tt.Frame()
		if *n-before < 3 {
			t.Fatalf("kind %d: no exit motion", kind)
		}
		tt.Advance(250 * time.Millisecond)
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
	ui.Render(func(c *ui.Context) { core.Use(c, core.Settings{}); Presence(c, true, PresenceOptions{Kind: 9}, nil) }, 100, 50, 1)
}

func TestPresenceSnapshot(t *testing.T) {
	snapshot(t, "presence", 300, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Presence(c, true, PresenceOptions{}, func() { ui.Text(c, "The doors open.") })
		})
	})
}

func TestPresenceReducedMotionCancelsActiveExit(t *testing.T) {
	show := true
	tt, n := motionFrames(func(c *ui.Context) {
		Presence(c, show, PresenceOptions{Duration: 5 * time.Second}, func() { ui.Text(c, "Leaving notice") })
	}, 300, 100)
	show = false
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionBusy(tt, n) || tt.HasText("Leaving notice") {
		t.Fatal("active exit continued under reduced motion")
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	if motionBusy(tt, n) {
		t.Fatal("exit replayed after restoring motion")
	}
}
