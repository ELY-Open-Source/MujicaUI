package mujica

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestBlinkHighlightFlashesOnceWithDirection(t *testing.T) {
	value := 10.0
	tt, n := motionFrames(func(c *ui.Context) {
		BlinkHighlight(c, value, BlinkHighlightOptions{Duration: 200 * time.Millisecond})
	}, 200, 60)
	if motionBusy(tt, n) {
		t.Fatal("first value flashed")
	}
	if _, ok := tt.Find("10.00"); !ok {
		t.Fatal("value not named")
	}
	value = 12
	tt.Frame()
	a := encode(t, tt.Image())
	if _, ok := tt.Find("12.00, up"); !ok {
		t.Fatalf("direction not named: %q", tt.Texts())
	}
	time.Sleep(250 * time.Millisecond)
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) || motionBusy(tt, n) {
		t.Fatal("flash did not fade and stop")
	}
	value = 9
	tt.Frame()
	if _, ok := tt.Find("9.00, down"); !ok {
		t.Fatalf("down not named: %q", tt.Texts())
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	value = 11
	tt.Frame()
	if motionBusy(tt, n) {
		t.Fatal("flash under reduced motion")
	}
	if _, ok := tt.Find("11.00, up"); !ok {
		t.Fatal("static mark missing under reduced motion")
	}
}

func TestBlinkHighlightInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		BlinkHighlight(c, 1, BlinkHighlightOptions{Duration: -1})
	}, 100, 50, 1)
}

func TestBlinkHighlightSnapshot(t *testing.T) {
	snapshotAfter(t, "blink_highlight", 300, 100, func(c *ui.Context) {
		v := ui.Local(c.Root(), "v", func() float64 { return 101.5 })
		w := ui.Local(c.Root(), "w", func() float64 { return 88.25 })
		ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
			BlinkHighlight(c, *v, BlinkHighlightOptions{})
			BlinkHighlight(c, *w, BlinkHighlightOptions{})
		})
		*v, *w = 103.75, 86.5
	}, func(tt *ui.Tester) { tt.Frame() })
}

func TestBlinkHighlightRestoredMotionKeepsStaticMark(t *testing.T) {
	value := 1.0
	tt, n := motionFrames(func(c *ui.Context) { BlinkHighlight(c, value, BlinkHighlightOptions{Duration: 5 * time.Second}) }, 300, 100)
	value = 2
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	if motionBusy(tt, n) {
		t.Fatal("old highlight replayed after restoring motion")
	}
	if _, ok := tt.Find("2.00, up"); !ok {
		t.Fatal("direction mark lost")
	}
}
