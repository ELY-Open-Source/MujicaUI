package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTitleBarDefaultHeight(t *testing.T) {
	var r ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r = TitleBar(c, "Archive", TitleBarOptions{}).Bounds()
	}, 400, 300).Frame()
	if r.H != 44 || r.W != 400 {
		t.Fatalf("bar %vx%v, want 400x44", r.W, r.H)
	}
}

func TestTitleBarClearsWindowControls(t *testing.T) {
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r = TitleBar(c, "Archive", TitleBarOptions{
			Trailing: func() { Button(c, "Share", ButtonOptions{Variant: Ghost}) },
		}).Bounds()
	}, 400, 300)
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 78, Right: 40})
	tt.Frame()
	if r.H != 52 {
		t.Fatalf("bar height %v, want the controls' 52", r.H)
	}
	title, _ := tt.Find("Archive")
	if title.X != 78+12 {
		t.Fatalf("title at %v, want 90 past the controls", title.X)
	}
	share, _ := tt.Find("Share")
	if share.X+share.W > 400-40-12+0.5 {
		t.Fatalf("trailing tool ends at %v, inside the right controls", share.X+share.W)
	}
}

func TestTitleBarSlotsStayClickable(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		TitleBar(c, "Archive", TitleBarOptions{
			Leading: func() {
				if Button(c, "Menu", ButtonOptions{Variant: Ghost}).Clicked() {
					n++
				}
			},
		})
	}, 400, 300)
	if err := tt.Click("Menu"); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("leading tool clicked %d times, want 1", n)
	}
}

func TestTitleBarLongTitleKeepsTools(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		TitleBar(c, "A very long window title that cannot possibly fit in the bar", TitleBarOptions{
			Trailing: func() { Button(c, "Share", ButtonOptions{Variant: Ghost}) },
		})
	}, 240, 100)
	share, ok := tt.Find("Share")
	if !ok || share.X+share.W > 240-12+0.5 {
		t.Fatalf("tool pushed out: %+v", share)
	}
}

func TestTitleBarSnapshot(t *testing.T) {
	snapshot(t, "title_bar", 420, 60, func(c *ui.Context) {
		TitleBar(c, "Chronicle of the Court", TitleBarOptions{
			Trailing: func() { Button(c, "Share", ButtonOptions{Variant: Ghost}) },
		})
	})
}

func TestTitleBarTextScale(t *testing.T) {
	var bar ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		bar = TitleBar(c, "Archive", TitleBarOptions{}).Bounds()
	}, 400, 300)
	plain, _ := tt.Find("Archive")
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	big, _ := tt.Find("Archive")
	if big.H < 1.8*plain.H {
		t.Fatalf("title height %v at 200%%, was %v", big.H, plain.H)
	}
	if bar.H < big.H {
		t.Fatalf("bar %v shorter than its title %v", bar.H, big.H)
	}
}
