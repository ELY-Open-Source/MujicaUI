package mujica

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"
)

type overlayApp struct {
	open   bool
	opts   OverlayOptions
	behind int
	panel  ui.Rect
}

func (a *overlayApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Row(c).Padding(20).Gap(8).Children(func() {
		if Button(c, "Show", ButtonOptions{}).Clicked() {
			a.open = true
		}
		if Button(c, "Behind", ButtonOptions{}).Clicked() {
			a.behind++
		}
	})
	if p := Overlay(c, &a.open, a.opts, func() {
		ui.Text(c, "Layer")
		if Button(c, "Done", ButtonOptions{}).Clicked() {
			a.open = false
		}
	}); p != nil {
		a.panel = p.Bounds()
	}
}

func newOverlayApp(opts OverlayOptions) (*overlayApp, *ui.Tester) {
	a := &overlayApp{opts: opts}
	tt := ui.NewTester(a.view, 600, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt.Click("Show")
	tt.Frame()
	return a, tt
}

func TestOverlayPlaces(t *testing.T) {
	for _, modal := range []bool{false, true} {
		for _, tc := range []struct {
			place  OverlayPlace
			fx, fy float32
		}{{OverlayCenter, 0.5, 0.5}, {OverlayTopLeft, 0, 0}, {OverlayBottomRight, 1, 1}, {OverlayTop, 0.5, 0}} {
			a, _ := newOverlayApp(OverlayOptions{Place: tc.place, Modal: modal})
			p := a.panel
			wantX := 16 + (600-32-p.W)*tc.fx
			wantY := 16 + (400-32-p.H)*tc.fy
			if abs(p.X-wantX) > 1 || abs(p.Y-wantY) > 1 {
				t.Errorf("modal %v place %d: panel %+v, want at (%v, %v)", modal, tc.place, p, wantX, wantY)
			}
		}
	}
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func TestOverlayModal(t *testing.T) {
	a, tt := newOverlayApp(OverlayOptions{Modal: true, Dim: true})
	if !tt.Focused("Done") {
		t.Fatal("modal overlay did not take the focus")
	}
	_ = tt.Click("Behind")
	if a.behind != 0 || a.open {
		t.Fatalf("backdrop: behind %d open %v", a.behind, a.open)
	}
	if !tt.Focused("Show") {
		t.Fatal("focus not restored")
	}
}

func TestOverlayNonModal(t *testing.T) {
	a, tt := newOverlayApp(OverlayOptions{Place: OverlayBottomRight})
	_ = tt.Click("Behind")
	if a.behind != 1 || !a.open {
		t.Fatalf("behind %d open %v", a.behind, a.open)
	}
	tt.Key(0, ui.KeyEscape)
	if a.open {
		t.Fatal("Escape did not close")
	}
}

func TestOverlayDimNeedsModal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := true
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Overlay(c, &open, OverlayOptions{Dim: true}, func() {})
	}, 200, 200)
}

// TestOverlayLifecycle opens and closes every overlay 100 times and checks
// that nothing grows: goroutines, heap objects or the elements left behind.
func TestOverlayLifecycle(t *testing.T) {
	var open bool
	views := map[string]func(c *ui.Context, b *ui.Element){
		"Dialog": func(c *ui.Context, b *ui.Element) { Dialog(c, &open, DialogOptions{Title: "T"}, func() {}) },
		"Alert": func(c *ui.Context, b *ui.Element) {
			AlertDialog(c, &open, AlertDialogOptions{Title: "T", ConfirmText: "x"})
		},
		"Drawer":  func(c *ui.Context, b *ui.Element) { Drawer(c, &open, DrawerOptions{Title: "T"}, nil) },
		"Drawer2": func(c *ui.Context, b *ui.Element) { Drawer(c, &open, DrawerOptions{Title: "T", NonModal: true}, nil) },
		"Popover": func(c *ui.Context, b *ui.Element) { Popover(c, b, &open, PopoverOptions{}, func() {}) },
		"Confirm": func(c *ui.Context, b *ui.Element) { Popconfirm(c, b, &open, PopconfirmOptions{Title: "T"}) },
		"Overlay": func(c *ui.Context, b *ui.Element) { Overlay(c, &open, OverlayOptions{Modal: true}, func() {}) },
		"Menu": func(c *ui.Context, b *ui.Element) {
			popupMenuOpen(c, b, &open, (&menuApp{}).build)
		},
		"Context": nil,
	}
	for name, view := range views {
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			row := ui.Row(c)
			var b *ui.Element
			row.Children(func() { b = Button(c, "Anchor", ButtonOptions{}) })
			if view == nil {
				ContextMenu(c, row, ContextMenuOptions{}, (&menuApp{}).build)
			} else {
				view(c, b)
			}
		}, 600, 400)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		cycle := func() {
			open = true
			if view == nil {
				_ = tt.Click("Anchor")
				tt.Key(ui.Shift, ui.KeyF10)
				if !tt.HasText("Rename") {
					t.Fatalf("%s: did not open", name)
				}
			}
			tt.Frame()
			if name != "Drawer2" {
				tt.Key(0, ui.KeyEscape)
				if view == nil {
					open = false
				}
			} else {
				open = false
				tt.Frame()
			}
			if open || tt.HasText("Rename") {
				t.Fatalf("%s: Escape did not close", name)
			}
		}
		for range 20 {
			cycle()
		}
		texts := len(tt.Texts())
		before := heapObjects()
		goroutines := runtime.NumGoroutine()
		for range 100 {
			cycle()
		}
		if got := len(tt.Texts()); got != texts {
			t.Errorf("%s: texts %d -> %d", name, texts, got)
		}
		if g := runtime.NumGoroutine(); g > goroutines {
			t.Errorf("%s: goroutines %d -> %d", name, goroutines, g)
		}
		if after := heapObjects(); after > before+2000 {
			t.Errorf("%s: heap objects %d -> %d", name, before, after)
		}
	}
}

func heapObjects() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapObjects
}

func TestOverlaySnapshot(t *testing.T) {
	snapshot(t, "overlay", 480, 260, func(c *ui.Context) {
		open := true
		ui.Column(c).Padding(24).Children(func() { ui.Text(c, "Score of the second act") })
		Overlay(c, &open, OverlayOptions{Place: OverlayBottomRight, Label: "Now playing"}, func() {
			ui.Text(c, "Now playing").FontWeight(600)
			ui.Text(c, "Symbol I: △").TextColor(Tokens(c).TextMuted)
		})
	})
}
