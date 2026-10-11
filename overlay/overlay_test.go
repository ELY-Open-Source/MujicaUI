package overlay

import (
	"runtime"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type overlayApp struct {
	open   bool
	opts   OverlayOptions
	behind int
	panel  ui.Rect
}

func (a *overlayApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Row(c).Padding(20).Gap(8).Children(func() {
		if input.Button(c, "Show", input.ButtonOptions{}).Clicked() {
			a.open = true
		}
		if input.Button(c, "Behind", input.ButtonOptions{}).Clicked() {
			a.behind++
		}
	})
	if p := Overlay(c, &a.open, a.opts, func() {
		ui.Text(c, "Layer")
		if input.Button(c, "Done", input.ButtonOptions{}).Clicked() {
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

func TestOverlayEscapeClosesTopNonModalLayerFirst(t *testing.T) {
	bottom, top := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Show", input.ButtonOptions{}).Absolute().Right(10).Top(10).Clicked() {
			if bottom {
				top = true
			}
			bottom = true
		}
		Overlay(c, &top, OverlayOptions{Place: OverlayTopRight}, func() { ui.Text(c, "Top layer") })
		Drawer(c, &bottom, DrawerOptions{Title: "Bottom layer", NonModal: true, Side: DrawerLeft}, nil)
	}, 600, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	_ = tt.Click("Show")
	_ = tt.Click("Show") // the overlay opens after the drawer, though built before it
	if !bottom || !top {
		t.Fatalf("bottom %v top %v", bottom, top)
	}
	tt.Key(0, ui.KeyEscape)
	if top || !bottom {
		t.Fatalf("first Escape: top %v bottom %v", top, bottom)
	}
	tt.Key(0, ui.KeyEscape)
	if bottom {
		t.Fatal("second Escape did not close the drawer")
	}
}

func TestOverlayNonModalRestoresFocusToItsPlace(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			input.Button(c, "First", input.ButtonOptions{})
			Overlay(c, &open, OverlayOptions{}, func() { input.Button(c, "Inside", input.ButtonOptions{}) })
			input.Button(c, "Next", input.ButtonOptions{})
		})
	}, 600, 400)
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Next") {
		t.Fatal("focus lost on close instead of staying at the overlay's place")
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
		core.Use(c, core.Settings{})
		Overlay(c, &open, OverlayOptions{Dim: true}, func() {})
	}, 200, 200)
}

// overlayCase drives one component open and closed through real input.
type overlayCase struct {
	name        string
	view        func(c *ui.Context, open *bool)
	open, close func(tt *ui.Tester)
	marker      string
}

func overlayCases() []overlayCase {
	clickAnchor := func(tt *ui.Tester) { _ = tt.Click("Anchor") }
	escape := func(tt *ui.Tester) { tt.Key(0, ui.KeyEscape) }
	hover := func(tt *ui.Tester) {
		r, _ := tt.Find("Anchor")
		tt.Move(r.X+4, r.Y+4)
	}
	leave := func(tt *ui.Tester) { tt.Move(590, 390) }
	anchor := func(c *ui.Context, open *bool) *ui.Element {
		b := input.Button(c, "Anchor", input.ButtonOptions{})
		if b.Clicked() {
			*open = true
		}
		return b
	}
	body := func(c *ui.Context) func() { return func() { ui.Text(c, "Layer body") } }
	return []overlayCase{
		{"Dialog", func(c *ui.Context, open *bool) {
			anchor(c, open)
			Dialog(c, open, DialogOptions{Title: "T"}, body(c))
		}, clickAnchor, escape, "Layer body"},
		{"AlertDialog", func(c *ui.Context, open *bool) {
			anchor(c, open)
			AlertDialog(c, open, AlertDialogOptions{Title: "Layer body", ConfirmText: "x"})
		}, clickAnchor, escape, "Layer body"},
		{"Drawer", func(c *ui.Context, open *bool) {
			anchor(c, open)
			Drawer(c, open, DrawerOptions{Title: "T"}, body(c))
		}, clickAnchor, escape, "Layer body"},
		{"DrawerNonModal", func(c *ui.Context, open *bool) {
			anchor(c, open)
			Drawer(c, open, DrawerOptions{Title: "T", NonModal: true}, body(c))
		}, clickAnchor, escape, "Layer body"},
		{"Popover", func(c *ui.Context, open *bool) {
			Popover(c, anchor(c, open), open, PopoverOptions{}, body(c))
		}, clickAnchor, escape, "Layer body"},
		{"Tooltip", func(c *ui.Context, open *bool) {
			Tooltip(c, anchor(c, open), "Layer body", TooltipOptions{Delay: time.Millisecond})
		}, hover, leave, "Layer body"},
		{"HoverCard", func(c *ui.Context, open *bool) {
			HoverCard(c, anchor(c, open), HoverCardOptions{OpenDelay: time.Millisecond, CloseDelay: time.Millisecond}, body(c))
		}, hover, leave, "Layer body"},
		{"DropdownMenu", func(c *ui.Context, open *bool) {
			DropdownMenu(c, "Anchor", DropdownMenuOptions{}, (&menuApp{}).build)
		}, clickAnchor, escape, "Rename"},
		{"ContextMenu", func(c *ui.Context, open *bool) {
			row := ui.Row(c)
			row.Children(func() { input.Button(c, "Anchor", input.ButtonOptions{}) })
			ContextMenu(c, row, ContextMenuOptions{}, (&menuApp{}).build)
		}, func(tt *ui.Tester) {
			clickAnchor(tt)
			tt.Key(ui.Shift, ui.KeyF10)
		}, escape, "Rename"},
		{"Popconfirm", func(c *ui.Context, open *bool) {
			Popconfirm(c, anchor(c, open), open, PopconfirmOptions{Title: "Layer body"})
		}, clickAnchor, escape, "Layer body"},
		{"Overlay", func(c *ui.Context, open *bool) {
			anchor(c, open)
			Overlay(c, open, OverlayOptions{Modal: true}, body(c))
		}, clickAnchor, escape, "Layer body"},
	}
}

// TestOverlayLifecycle opens and closes every overlay 100 times; goroutines, heap objects and elements must not grow.
func TestOverlayLifecycle(t *testing.T) {
	for _, oc := range overlayCases() {
		open := false
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			oc.view(c, &open)
		}, 600, 400)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		shown := func() bool { return tt.HasText(oc.marker) }
		cycle := func(i int) {
			oc.open(tt)
			if !overlayUntil(tt, shown) {
				t.Fatalf("%s: cycle %d did not open", oc.name, i)
			}
			oc.close(tt)
			if !overlayUntil(tt, func() bool { return !shown() }) {
				t.Fatalf("%s: cycle %d did not close", oc.name, i)
			}
		}
		for i := range 20 {
			cycle(i)
		}
		texts := len(tt.Texts())
		before := heapObjects()
		goroutines := runtime.NumGoroutine()
		for i := range 100 {
			cycle(i)
		}
		if got := len(tt.Texts()); got != texts {
			t.Errorf("%s: texts %d -> %d", oc.name, texts, got)
		}
		// A fired c.After timer (ui/runtime.go:596) runs briefly in its own goroutine; only lasting growth counts.
		if g := overlayGoroutines(goroutines); g > goroutines {
			t.Errorf("%s: goroutines %d -> %d", oc.name, goroutines, g)
		}
		if after := heapObjects(); after > before+2000 {
			t.Errorf("%s: heap objects %d -> %d", oc.name, before, after)
		}
	}
}

// overlayGoroutines waits up to 2 s for the goroutine count to fall to want, and returns it.
func overlayGoroutines(want int) int {
	end := time.Now().Add(2 * time.Second)
	g := runtime.NumGoroutine()
	for g > want && time.Now().Before(end) {
		time.Sleep(10 * time.Millisecond)
		g = runtime.NumGoroutine()
	}
	return g
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
			ui.Text(c, "Symbol I: △").TextColor(core.Tokens(c).TextMuted)
		})
	})
}
