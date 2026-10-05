package mujica

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestShimmerTextSweepsAndStops(t *testing.T) {
	active := true
	tt, n := motionFrames(func(c *ui.Context) {
		ShimmerText(c, "Thinking", ShimmerTextOptions{Active: active})
	}, 200, 60)
	if _, ok := tt.Find("Thinking"); !ok {
		t.Fatal("text not named")
	}
	a := encode(t, tt.Image())
	time.Sleep(200 * time.Millisecond)
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) || !motionBusy(tt, n) {
		t.Fatal("no sweep")
	}
	active = false
	tt.Frame()
	if motionBusy(tt, n) || !tt.HasText("Thinking") {
		t.Fatal("inactive text still sweeps")
	}
	active = true
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionBusy(tt, n) {
		t.Fatal("sweeps under reduced motion")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("empty text did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		ShimmerText(c, "", ShimmerTextOptions{})
	}, 100, 50, 1)
}

func TestShimmerTextSnapshot(t *testing.T) {
	snapshot(t, "shimmer_text", 300, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			ShimmerText(c, "Thinking…", ShimmerTextOptions{Active: true})
			ShimmerText(c, "思考中…", ShimmerTextOptions{})
		})
	})
}

func TestShimmerTextConstrainedWidth(t *testing.T) {
	var background ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		background = Tokens(c).Background
		ui.Column(c).Padding(10).Children(func() {
			ShimmerText(c, "Thinking about the complete program and planning the next steps", ShimmerTextOptions{Active: true}).Width(80)
		})
	}, 500, 100)
	for _, reduce := range []bool{false, true} {
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
		img := tt.Image()
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 92; x < img.Bounds().Dx(); x++ {
				pixel := img.RGBAAt(x, y)
				if pixel.R != background.R || pixel.G != background.G || pixel.B != background.B {
					t.Fatalf("text escaped 80 DIP width at %d,%d", x, y)
				}
			}
		}
	}
}
