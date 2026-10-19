package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func speedTester(rate *float64, opts PlaybackSpeedControlOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			if PlaybackSpeedControl(c, rate, opts).Changed() {
				*n++
			}
		})
	}, 300, 400)
}

func TestPlaybackSpeedChoose(t *testing.T) {
	rate, n := 1.0, 0
	tt := speedTester(&rate, PlaybackSpeedControlOptions{}, &n)
	if !tt.HasText("1×") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Playback speed"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"0.5×", "0.75×", "1.25×", "1.5×", "2×"} {
		if !tt.HasText(want) {
			t.Fatalf("no %s in %v", want, tt.Texts())
		}
	}
	if err := tt.Click("2×"); err != nil {
		t.Fatal(err)
	}
	if rate != 2 || n != 1 {
		t.Fatalf("rate %v changes %d", rate, n)
	}
}

func TestPlaybackSpeedKeyboardAndDisabled(t *testing.T) {
	rate, n := 1.0, 0
	tt := speedTester(&rate, PlaybackSpeedControlOptions{}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if rate != 1.25 {
		t.Fatalf("keyboard chose %v", rate)
	}
	rate2, n2 := 1.0, 0
	off := speedTester(&rate2, PlaybackSpeedControlOptions{Disabled: true}, &n2)
	_ = off.Click("Playback speed")
	_ = off.Click("2×")
	if rate2 != 1 || n2 != 0 {
		t.Fatalf("disabled chose %v", rate2)
	}
}

func TestPlaybackSpeedBadOptions(t *testing.T) {
	odd := 3.0
	mediaPanics(t, "rate not offered", func(c *ui.Context) { PlaybackSpeedControl(c, &odd, PlaybackSpeedControlOptions{}) })
	one := 1.0
	mediaPanics(t, "unsorted", func(c *ui.Context) {
		PlaybackSpeedControl(c, &one, PlaybackSpeedControlOptions{Rates: []float64{2, 1}})
	})
}

func TestPlaybackSpeedSnapshot(t *testing.T) {
	rate := 1.25
	snapshot(t, "playback_speed_control", 200, 80, func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			PlaybackSpeedControl(c, &rate, PlaybackSpeedControlOptions{})
		})
	})
}
