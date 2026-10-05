package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

type dialogApp struct {
	open      bool
	opts      DialogOptions
	saved     int
	name      string
	densities theme.Density
}

func (a *dialogApp) view(c *ui.Context) {
	Use(c, Settings{Mode: Light, Density: a.densities})
	if Button(c, "Rename", ButtonOptions{}).Clicked() {
		a.open = true
	}
	a.opts.Actions = func() {
		if Button(c, "Cancel", ButtonOptions{Variant: Secondary}).Clicked() {
			a.open = false
		}
		if Button(c, "Save", ButtonOptions{}).Clicked() {
			a.saved++
			a.open = false
		}
	}
	Dialog(c, &a.open, a.opts, func() {
		ui.TextInput(c, &a.name)
	})
}

func newDialogApp(opts DialogOptions) (*dialogApp, *ui.Tester) {
	opts.Title = "Rename"
	a := &dialogApp{opts: opts}
	tt := ui.NewTester(a.view, 600, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return a, tt
}

func TestDialogOpenFocusTrapRestore(t *testing.T) {
	a, tt := newDialogApp(DialogOptions{Description: "Pick a new name."})
	if err := tt.Click("Rename"); err != nil {
		t.Fatal(err)
	}
	if !a.open || !tt.HasText("Pick a new name.") {
		t.Fatalf("open %v, texts %q", a.open, tt.Texts())
	}
	if !tt.Focused("Close") {
		t.Fatalf("focus did not move into the dialog")
	}
	// Tab cycles Close, input, Cancel, Save and back to Close.
	for range 3 {
		tt.Key(0, ui.KeyTab)
	}
	if !tt.Focused("Save") {
		t.Fatalf("tab order left the dialog")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Close") {
		t.Fatalf("tab did not wrap inside the dialog")
	}
	tt.Key(0, ui.KeyEscape)
	if a.open {
		t.Fatal("Escape did not close")
	}
	if !tt.Focused("Rename") {
		t.Fatal("focus not restored to the opener")
	}
}

func TestDialogActionsAndClose(t *testing.T) {
	a, tt := newDialogApp(DialogOptions{})
	_ = tt.Click("Rename")
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if a.saved != 1 || a.open {
		t.Fatalf("saved %d open %v", a.saved, a.open)
	}
	_ = tt.Click("Rename")
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if a.open || a.saved != 1 {
		t.Fatalf("close button: open %v saved %d", a.open, a.saved)
	}
}

func TestDialogCloseRules(t *testing.T) {
	a, tt := newDialogApp(DialogOptions{})
	_ = tt.Click("Rename")
	tt.ClickAt(10, 390)
	if a.open {
		t.Fatal("backdrop click did not close")
	}
	a, tt = newDialogApp(DialogOptions{KeepOnBackdrop: true, KeepOnEscape: true, HideClose: true})
	_ = tt.Click("Rename")
	tt.ClickAt(10, 390)
	tt.Key(0, ui.KeyEscape)
	if !a.open {
		t.Fatal("kept dialog closed")
	}
	if tt.HasText("Close") {
		t.Fatal("close button shown with HideClose")
	}
	if _, ok := tt.Find("Close"); ok {
		t.Fatal("close button present with HideClose")
	}
}

func TestDialogEscapeGoesToNestedPopoverFirst(t *testing.T) {
	open, inner := true, false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Dialog(c, &open, DialogOptions{Title: "Outer"}, func() {
			b := Button(c, "More", ButtonOptions{})
			if b.Clicked() {
				inner = true
			}
			Popover(c, b, &inner, PopoverOptions{}, func() { ui.Text(c, "Inner") })
		})
	}, 600, 400)
	_ = tt.Click("More")
	if !inner {
		t.Fatal("popover did not open")
	}
	tt.Key(0, ui.KeyEscape)
	if inner || !open {
		t.Fatalf("after one Escape inner %v outer %v", inner, open)
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("second Escape did not close the dialog")
	}
}

func TestDialogNeedsTitle(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := true
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Dialog(c, &open, DialogOptions{}, nil)
	}, 300, 200)
}

func dialogSnapshotView(d theme.Density, decorated bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		mode := Light
		if Tokens(c).Dark {
			mode = Dark
		}
		Use(c, Settings{Mode: mode, Density: d})
		open := true
		name := "Ophelia"
		ui.Column(c).Padding(24).Children(func() { ui.Text(c, "Behind the dialog") })
		Dialog(c, &open, DialogOptions{
			Title: "Rename the chamber", Description: "The new name shows in every hall.", Decorated: decorated,
			Actions: func() {
				Button(c, "Cancel", ButtonOptions{Variant: Secondary})
				Button(c, "Rename", ButtonOptions{})
			},
		}, func() { ui.TextInput(c, &name) })
	}
}

func TestDialogSnapshot(t *testing.T) {
	snapshot(t, "dialog", 560, 360, dialogSnapshotView(theme.Default, true))
	snapshot(t, "dialog-compact", 560, 360, dialogSnapshotView(theme.Compact, false))
	snapshot(t, "dialog-comfortable", 560, 400, dialogSnapshotView(theme.Comfortable, true))
}
