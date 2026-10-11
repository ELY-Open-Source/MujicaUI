package layout

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func scrollAreaView(st *ui.ScrollState, opts ScrollAreaOptions, rows int) func(c *ui.Context) {
	return func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ScrollArea(c, st, opts, func() {
			for i := range rows {
				ui.Text(c, fmt.Sprintf("Row %d", i)).Height(20).Width(float32(100 + 40*i))
			}
		}).Size(200, 100)
	}
}

func TestScrollAreaScrollToEnd(t *testing.T) {
	st := &ui.ScrollState{}
	tt := ui.NewTester(scrollAreaView(st, ScrollAreaOptions{}, 20), 300, 300)
	if st.MaxY != 300 {
		t.Fatalf("MaxY %v, want 300", st.MaxY)
	}
	st.Y = math.MaxFloat32
	tt.Frame()
	if st.Y != 300 {
		t.Fatalf("Y %v after scrolling to the end, want 300", st.Y)
	}
	st.Y = 40
	tt.Frame()
	if r, _ := tt.Find("Row 2"); r.Y != 0 {
		t.Fatalf("Row 2 at %v after scrolling to 40, want 0", r.Y)
	}
}

func TestScrollAreaWheelAndKeys(t *testing.T) {
	st := &ui.ScrollState{}
	tt := ui.NewTester(scrollAreaView(st, ScrollAreaOptions{Label: "Log"}, 20), 300, 300)
	tt.Scroll(50, 50, 0, 60)
	if st.Y != 60 {
		t.Fatalf("wheel moved to %v, want 60", st.Y)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Log") {
		t.Fatal("overflowing area did not take the focus")
	}
	tt.Key(0, ui.KeyDown)
	if st.Y != 100 {
		t.Fatalf("Down moved to %v, want 100", st.Y)
	}
	tt.Key(0, ui.KeyEnd)
	if st.Y != 300 {
		t.Fatalf("End moved to %v, want 300", st.Y)
	}
	tt.Key(0, ui.KeyHome)
	if st.Y != 0 {
		t.Fatalf("Home moved to %v, want 0", st.Y)
	}
}

func TestScrollAreaNoOverflowSkipsFocus(t *testing.T) {
	tt := ui.NewTester(scrollAreaView(&ui.ScrollState{}, ScrollAreaOptions{Label: "Log"}, 2), 300, 300)
	tt.Key(0, ui.KeyTab)
	if tt.Focused("Log") {
		t.Fatal("area without overflow took the focus")
	}
}

func TestScrollAreaHorizontal(t *testing.T) {
	st := &ui.ScrollState{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ScrollArea(c, st, ScrollAreaOptions{Direction: ScrollHorizontal}, func() {
			ui.Box(c).Size(500, 50)
		}).Size(200, 60)
	}, 300, 300)
	tt.Scroll(50, 30, 120, 0)
	if st.X != 120 || st.Y != 0 || st.MaxX != 300 {
		t.Fatalf("offset %v,%v max %v, want 120,0 max 300", st.X, st.Y, st.MaxX)
	}
}

func TestScrollAreaKeepsPlaceWithoutState(t *testing.T) {
	tt := ui.NewTester(scrollAreaView(nil, ScrollAreaOptions{Direction: ScrollBoth}, 20), 300, 300)
	tt.Scroll(50, 50, 0, 80)
	tt.Frame()
	tt.Frame()
	if r, _ := tt.Find("Row 4"); r.Y != 0 {
		t.Fatalf("Row 4 at %v after frames, want 0", r.Y)
	}
}

func TestScrollAreaInvalidDirectionPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("invalid direction did not panic")
		}
	}()
	ui.NewTester(scrollAreaView(nil, ScrollAreaOptions{Direction: 7}, 1), 100, 100)
}

func TestScrollAreaSnapshot(t *testing.T) {
	snapshot(t, "scroll_area", 300, 160, func(c *ui.Context) {
		k := core.Tokens(c)
		st := ui.Local(c.Root(), "snap", func() ui.ScrollState { return ui.ScrollState{Y: 30} })
		ui.Column(c).Padding(16).Children(func() {
			ScrollArea(c, st, ScrollAreaOptions{}, func() {
				for i := range 12 {
					ui.Text(c, fmt.Sprintf("Verse %d", i+1)).Padding(4, 12)
				}
			}).Height(128).Background(k.Surface).Border(1, k.Border).Radius(4)
		})
	})
}

func TestScrollAreaFocusRingKeepsGap(t *testing.T) {
	var focus, bg ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		focus, bg = core.Tokens(c).Focus, core.Tokens(c).Background
		ScrollArea(c, nil, ScrollAreaOptions{Label: "Log"}, func() { ui.Box(c).Height(400).Shrink(0) }).Size(100, 100)
	}, 200, 200)
	tt.Key(0, ui.KeyTab)
	img := tt.Image()
	// The edge's first 2 DIPs are the gap; the ring follows inside it.
	if got := img.RGBAAt(50, 1); got != (color.RGBA{bg.R, bg.G, bg.B, 255}) {
		t.Fatalf("gap pixel %v, want background %v", got, bg)
	}
	if got := img.RGBAAt(50, 3); got != (color.RGBA{focus.R, focus.G, focus.B, 255}) {
		t.Fatalf("ring pixel %v, want focus %v", got, focus)
	}
}
