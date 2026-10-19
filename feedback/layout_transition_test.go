package feedback

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

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
	tt.Advance(260 * time.Millisecond)
	if motionBusy(tt, n) || w != 300 {
		t.Fatalf("settled width %v", w)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative duration did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		LayoutTransition(c, LayoutTransitionOptions{Duration: -1}, nil)
	}, 100, 50, 1)
}

func TestLayoutTransitionSnapshot(t *testing.T) {
	snapshot(t, "layout_transition", 300, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
			LayoutTransition(c, LayoutTransitionOptions{}, func() {
				ui.Text(c, "A panel that grows").Padding(12).Background(core.Tokens(c).Surface).Radius(4)
			})
		})
	})
}
