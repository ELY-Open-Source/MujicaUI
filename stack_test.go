package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestStackVerticalGap(t *testing.T) {
	var a, b ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Stack(c, StackOptions{}, func() {
			a = ui.Box(c).Height(20).Bounds()
			b = ui.Box(c).Height(20).Bounds()
		})
	}, 300, 300).Frame()
	if b.Y-(a.Y+a.H) != 8 {
		t.Fatalf("gap %v, want 8", b.Y-(a.Y+a.H))
	}
	if a.W != 300 {
		t.Fatalf("vertical stack did not stretch: width %v", a.W)
	}
}

func TestStackHiddenItemTakesNoGap(t *testing.T) {
	for _, hide := range []bool{false, true} {
		var a, last ui.Rect
		ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			Stack(c, StackOptions{Horizontal: true, Gap: SpaceL}, func() {
				a = ui.Box(c).Size(30, 10).Bounds()
				if !hide {
					ui.Box(c).Size(30, 10)
				}
				last = ui.Box(c).Size(30, 10).Bounds()
			})
		}, 300, 300).Frame()
		want := a.X + 30 + 16
		if !hide {
			want += 30 + 16
		}
		if last.X != want {
			t.Errorf("hidden=%v: last at %v, want %v", hide, last.X, want)
		}
	}
}

func TestStackAlignEnd(t *testing.T) {
	var r ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(200).Children(func() {
			Stack(c, StackOptions{Align: StackAlignEnd}, func() {
				r = ui.Box(c).Size(50, 10).Bounds()
			})
		})
	}, 300, 300).Frame()
	if r.X != 150 {
		t.Fatalf("end-aligned child at %v, want 150", r.X)
	}
}

func TestStackJustifyAndWrap(t *testing.T) {
	var first, third ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(250).Children(func() {
			Stack(c, StackOptions{Horizontal: true, Wrap: true, Gap: SpaceNone, Justify: ui.End}, func() {
				first = ui.Box(c).Size(100, 20).Bounds()
				ui.Box(c).Size(100, 20)
				third = ui.Box(c).Size(100, 20).Bounds()
			})
		})
	}, 300, 300).Frame()
	if third.Y != first.Y+20 {
		t.Fatalf("third at y=%v, want the next line at %v", third.Y, first.Y+20)
	}
	if first.X != 50 || third.X != 150 {
		t.Fatalf("end-justified lines start at %v and %v, want 50 and 150", first.X, third.X)
	}
}

func TestStackWrapNeedsHorizontal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("vertical Wrap did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Stack(c, StackOptions{Wrap: true}, func() {})
	}, 100, 100)
}

func TestStackSnapshot(t *testing.T) {
	snapshot(t, "stack", 420, 150, func(c *ui.Context) {
		k := Tokens(c)
		chip := func(s string) {
			ui.Box(c).Padding(6, 12).Background(k.Surface).Border(1, k.Border).Radius(2).Children(func() { ui.Text(c, s) })
		}
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			Stack(c, StackOptions{Horizontal: true}, func() {
				for _, s := range []string{"Overture", "Aria", "Recitative"} {
					chip(s)
				}
			})
			Stack(c, StackOptions{Horizontal: true, Justify: ui.SpaceBetween}, func() {
				chip("Start")
				chip("Between")
				chip("End")
			})
		})
	})
}
