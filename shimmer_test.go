package mujica

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestShimmerOnlyWhileLoading(t *testing.T) {
	loading := true
	tt, n := motionFrames(func(c *ui.Context) {
		Shimmer(c, loading, ShimmerOptions{Height: 40, Label: "Loading report"}, func() { ui.Text(c, "Report ready") })
	}, 300, 100)
	if _, ok := tt.Find("Loading report"); !ok || tt.HasText("Report ready") {
		t.Fatalf("loading: %q", tt.Texts())
	}
	a := encode(t, tt.Image())
	time.Sleep(150 * time.Millisecond)
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) || !motionBusy(tt, n) {
		t.Fatal("shimmer does not sweep")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionBusy(tt, n) {
		t.Fatal("shimmer repaints under reduced motion")
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	loading = false
	tt.Frame()
	if !tt.HasText("Report ready") || motionBusy(tt, n) {
		t.Fatalf("loaded: %q", tt.Texts())
	}
}

func TestShimmerInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		Shimmer(c, true, ShimmerOptions{Width: -1}, nil)
	}, 100, 50, 1)
}

func TestShimmerSnapshot(t *testing.T) {
	snapshot(t, "shimmer", 400, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			Shimmer(c, true, ShimmerOptions{}, nil)
			Shimmer(c, true, ShimmerOptions{Width: 160, Height: 24}, nil)
			Shimmer(c, false, ShimmerOptions{}, func() { ui.Text(c, "Loaded content") })
		})
	})
}

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
