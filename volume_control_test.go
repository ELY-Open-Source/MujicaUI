package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func volumeTester(v *float64, muted *bool, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(40).AlignItems(ui.Start).Children(func() {
			if VolumeControl(c, v, muted, VolumeControlOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 400, 200)
}

func TestVolumeControlArrowsAndMute(t *testing.T) {
	v, muted, n := 0.5, false, 0
	tt := volumeTester(&v, &muted, false, &n)
	if err := tt.Click("Mute"); err != nil {
		t.Fatal(err)
	}
	if !muted {
		t.Fatal("mute button did not mute")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	if v != 0.45 || !muted {
		t.Fatalf("Left: %v muted %v", v, muted)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if v != 0.55 || muted {
		t.Fatalf("raising should unmute: %v muted %v", v, muted)
	}
	tt.Key(0, ui.KeyM)
	if !muted || n != 5 {
		t.Fatalf("M: muted %v changes %d", muted, n)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 1 {
		t.Fatalf("End: %v", v)
	}
}

func TestVolumeControlDisabled(t *testing.T) {
	v, muted, n := 0.5, false, 0
	tt := volumeTester(&v, &muted, true, &n)
	_ = tt.Click("Mute")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyM)
	if v != 0.5 || muted || n != 0 {
		t.Fatalf("disabled changed %v %v", v, muted)
	}
}

func TestVolumeControlRange(t *testing.T) {
	loud, m := 1.5, false
	mediaPanics(t, "over 1", func(c *ui.Context) { VolumeControl(c, &loud, &m, VolumeControlOptions{}) })
}

func TestVolumeControlSnapshot(t *testing.T) {
	a, b, ma, mb := 0.7, 0.3, false, true
	snapshot(t, "volume_control", 260, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).AlignItems(ui.Start).Children(func() {
			VolumeControl(c, &a, &ma, VolumeControlOptions{})
			VolumeControl(c, &b, &mb, VolumeControlOptions{})
		})
	})
}
