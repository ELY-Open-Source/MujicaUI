package feedback

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
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
	tt.Advance(150 * time.Millisecond)
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
		core.Use(c, core.Settings{})
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
