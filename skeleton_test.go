package mujica

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestSkeletonSizes(t *testing.T) {
	var text, circle, rect, fixed ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Density: theme.Compact})
		ui.Column(c).Width(300).AlignItems(ui.Start).Children(func() {
			text = Skeleton(c, SkeletonOptions{Lines: 3}).Bounds()
			circle = Skeleton(c, SkeletonOptions{Shape: SkeletonCircle}).Bounds()
			rect = Skeleton(c, SkeletonOptions{Shape: SkeletonRect}).Bounds()
			fixed = Skeleton(c, SkeletonOptions{Shape: SkeletonRect, Width: 120, Height: 80}).Bounds()
		})
	}, 400, 400)
	tt.Frame()
	if text.W != 300 || text.H != 63 {
		t.Errorf("text %v", text)
	}
	if circle.W != 40 || circle.H != 40 {
		t.Errorf("circle %v", circle)
	}
	if rect.W != 300 || rect.H != 28 {
		t.Errorf("rect %v", rect)
	}
	if fixed.W != 120 || fixed.H != 80 {
		t.Errorf("fixed %v", fixed)
	}
}

func TestSkeletonStaticUnderReducedMotion(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Skeleton(c, SkeletonOptions{Shape: SkeletonRect})
	}, 200, 100)
	a := encode(t, tt.Image())
	time.Sleep(300 * time.Millisecond)
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("skeleton did not pulse")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	a = encode(t, tt.Image())
	time.Sleep(300 * time.Millisecond)
	tt.Frame()
	if !bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("skeleton moved under reduced motion")
	}
}

func TestSkeletonPanicsOnNegativeSize(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Skeleton(c, SkeletonOptions{Lines: -1})
	}, 200, 100)
}

func TestSkeletonSnapshot(t *testing.T) {
	snapshot(t, "skeleton", 420, 200, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				Skeleton(c, SkeletonOptions{Shape: SkeletonCircle})
				ui.Column(c).Grow(1).Children(func() {
					Skeleton(c, SkeletonOptions{Lines: 2})
				})
			})
			Skeleton(c, SkeletonOptions{Shape: SkeletonRect, Height: 72})
		})
	})
}
