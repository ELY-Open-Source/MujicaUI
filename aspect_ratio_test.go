package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// aspectRatioBounds returns the frame (via a filling child) and a fixed 40x20 child.
func aspectRatioBounds(t *testing.T, ratio float32, w, h int) (box, child ui.Rect) {
	t.Helper()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Size(float32(w), float32(h)).Children(func() {
			AspectRatio(c, AspectRatioOptions{Ratio: ratio}, func() {
				box = ui.Box(c).Fill().Absolute().Bounds()
				child = ui.Box(c).Size(40, 20).Bounds()
			})
		})
	}, 600, 600).Frame()
	return box, child
}

func TestAspectRatioWidthBound(t *testing.T) {
	box, _ := aspectRatioBounds(t, 16.0/9, 320, 400)
	if box.W != 320 || box.H != 180 {
		t.Fatalf("box %vx%v, want 320x180", box.W, box.H)
	}
}

func TestAspectRatioHeightBound(t *testing.T) {
	box, _ := aspectRatioBounds(t, 2, 400, 100)
	if box.W != 200 || box.H != 100 {
		t.Fatalf("box %vx%v, want 200x100", box.W, box.H)
	}
}

func TestAspectRatioDefaultSquareAndUnstretchedChild(t *testing.T) {
	box, child := aspectRatioBounds(t, 0, 200, 400)
	if box.W != 200 || box.H != 200 {
		t.Fatalf("box %vx%v, want 200x200", box.W, box.H)
	}
	if child.W != 40 || child.H != 20 {
		t.Fatalf("child stretched to %vx%v", child.W, child.H)
	}
	if child.X != box.X+80 || child.Y != box.Y+90 {
		t.Fatalf("child at %v,%v not centered in %+v", child.X, child.Y, box)
	}
}

func TestAspectRatioFollowsWindowResize(t *testing.T) {
	var box ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Fill().Children(func() {
			AspectRatio(c, AspectRatioOptions{Ratio: 2}, func() {
				box = ui.Box(c).Fill().Absolute().Bounds()
			})
		})
	}, 400, 100)
	if box.W != 200 || box.H != 100 {
		t.Fatalf("short window: %vx%v, want 200x100", box.W, box.H)
	}
	tt.SetSize(400, 400)
	tt.Frame()
	if box.W != 400 || box.H != 200 {
		t.Fatalf("tall window: %vx%v, want 400x200", box.W, box.H)
	}
	tt.SetSize(300, 400)
	tt.Frame()
	if box.W != 300 || box.H != 150 {
		t.Fatalf("narrow window: %vx%v, want 300x150", box.W, box.H)
	}
}

func TestAspectRatioRejectsNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("negative ratio did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		AspectRatio(c, AspectRatioOptions{Ratio: -1}, func() {})
	}, 100, 100)
}

func TestAspectRatioSnapshot(t *testing.T) {
	snapshot(t, "aspect_ratio", 420, 160, func(c *ui.Context) {
		k := Tokens(c)
		ui.Row(c).Padding(16).Gap(16).Height(160).AlignItems(ui.Stretch).Children(func() {
			for _, r := range []float32{1, 16.0 / 9, 3.0 / 4} {
				ui.Box(c).Width(120).Children(func() {
					AspectRatio(c, AspectRatioOptions{Ratio: r}, func() {
						ui.Box(c).Fill().Center().Background(k.Surface).Border(1, k.Border).Radius(4).Children(func() {
							ui.Text(c, aspectRatioLabel(r)).TextColor(k.TextMuted)
						})
					})
				})
			}
		})
	})
}

func aspectRatioLabel(r float32) string {
	switch r {
	case 1:
		return "1 : 1"
	case 3.0 / 4:
		return "3 : 4"
	}
	return "16 : 9"
}
