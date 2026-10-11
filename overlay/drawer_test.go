package overlay

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type drawerApp struct {
	open    bool
	opts    DrawerOptions
	clicked int
	panel   ui.Rect
}

func (a *drawerApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Row(c).Gap(8).Children(func() {
		if input.Button(c, "Filters", input.ButtonOptions{}).Clicked() {
			a.open = true
		}
		if input.Button(c, "Behind", input.ButtonOptions{}).Clicked() {
			a.clicked++
		}
	})
	if p := Drawer(c, &a.open, a.opts, func() {
		input.Button(c, "Apply", input.ButtonOptions{})
	}); p != nil {
		a.panel = p.Bounds()
	}
}

func newDrawerApp(opts DrawerOptions) (*drawerApp, *ui.Tester) {
	opts.Title = "Filters panel"
	a := &drawerApp{opts: opts}
	tt := ui.NewTester(a.view, 800, 500)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt.Click("Filters")
	tt.Frame()
	return a, tt
}

func TestDrawerSides(t *testing.T) {
	for _, tc := range []struct {
		side       DrawerSide
		x, y, w, h float32
	}{
		{DrawerRight, 440, 0, 360, 500},
		{DrawerLeft, 0, 0, 360, 500},
		{DrawerTop, 0, 0, 800, 280},
		{DrawerBottom, 0, 220, 800, 280},
	} {
		for _, nonModal := range []bool{false, true} {
			a, _ := newDrawerApp(DrawerOptions{Side: tc.side, NonModal: nonModal})
			got := a.panel
			if got.X != tc.x || got.Y != tc.y || got.W != tc.w || got.H != tc.h {
				t.Errorf("side %d nonModal %v: panel %+v", tc.side, nonModal, got)
			}
		}
	}
}

func TestDrawerSizeLimits(t *testing.T) {
	a, _ := newDrawerApp(DrawerOptions{Size: 5000})
	if a.panel.W != 720 {
		t.Errorf("oversize width %v, want 720 (90%%)", a.panel.W)
	}
	a, _ = newDrawerApp(DrawerOptions{Size: 40})
	if a.panel.W != 160 {
		t.Errorf("undersize width %v, want 160", a.panel.W)
	}
}

func TestDrawerModalTrapsAndRestores(t *testing.T) {
	a, tt := newDrawerApp(DrawerOptions{})
	_ = tt.Click("Behind")
	if a.clicked != 0 {
		t.Fatal("modal drawer let a click through")
	}
	a.open = true
	tt.Frame()
	if !tt.Focused("Close") {
		t.Fatal("focus not in the drawer")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Close") {
		t.Fatal("tab left the modal drawer")
	}
	tt.Key(0, ui.KeyEscape)
	if a.open || !tt.Focused("Filters") {
		t.Fatalf("Escape: open %v, opener focused %v", a.open, tt.Focused("Filters"))
	}
}

func TestDrawerNonModalLeavesWindowUsable(t *testing.T) {
	a, tt := newDrawerApp(DrawerOptions{NonModal: true})
	if err := tt.Click("Behind"); err != nil {
		t.Fatal(err)
	}
	if a.clicked != 1 || !a.open {
		t.Fatalf("clicked %d open %v", a.clicked, a.open)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if a.open {
		t.Fatal("close button did not close")
	}
	a.open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if a.open {
		t.Fatal("Escape did not close the non-modal drawer")
	}
}

func TestDrawerSnapshot(t *testing.T) {
	snapshot(t, "drawer", 640, 360, func(c *ui.Context) {
		open := true
		ui.Column(c).Padding(24).Children(func() { ui.Text(c, "The gallery behind") })
		Drawer(c, &open, DrawerOptions{Title: "Seating chart", Size: 300, Actions: func() {
			input.Button(c, "Reset", input.ButtonOptions{Variant: input.Secondary})
			input.Button(c, "Apply", input.ButtonOptions{})
		}}, func() {
			ui.Text(c, "Guests are placed by rank, then by arrival.").TextColor(core.Tokens(c).TextMuted)
		})
	})
}
