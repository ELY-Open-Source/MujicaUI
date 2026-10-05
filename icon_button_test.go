package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

func TestIconButtonClickByNameAndKeyboard(t *testing.T) {
	n := 0
	var w, h float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Children(func() {
			b := IconButton(c, icons.Must("settings"), "Settings", IconButtonOptions{Variant: Ghost})
			w, h = b.Bounds().W, b.Bounds().H
			if b.Clicked() {
				n++
			}
		})
	}, 300, 200)
	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if n != 2 {
		t.Fatalf("clicks %d", n)
	}
	if w != 36 || h != 36 {
		t.Fatalf("size %vx%v, want 36 square", w, h)
	}
}

func TestIconButtonBlockedStates(t *testing.T) {
	for _, opts := range []IconButtonOptions{{Disabled: true}, {Loading: true}} {
		n := 0
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			if IconButton(c, icons.Must("trash-2"), "Delete", opts).Clicked() {
				n++
			}
		}, 300, 200)
		_ = tt.Click("Delete")
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		if n != 0 {
			t.Fatalf("%+v: clicked %d", opts, n)
		}
	}
}

func TestIconButtonNeedsName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without an action name")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		IconButton(c, icons.Must("x"), "", IconButtonOptions{})
	}, 100, 100)
}

func TestIconButtonSnapshot(t *testing.T) {
	snapshot(t, "icon_button", 360, 90, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(12).Children(func() {
			IconButton(c, icons.Must("plus"), "Add", IconButtonOptions{})
			IconButton(c, icons.Must("settings"), "Settings", IconButtonOptions{Variant: Secondary})
			IconButton(c, icons.Must("star"), "Favorite", IconButtonOptions{Variant: Outline})
			IconButton(c, icons.Must("copy"), "Copy", IconButtonOptions{Variant: Ghost})
			IconButton(c, icons.Must("trash-2"), "Delete", IconButtonOptions{Variant: Danger})
			IconButton(c, icons.Must("refresh-cw"), "Refresh", IconButtonOptions{Loading: true})
			IconButton(c, icons.Must("lock"), "Lock", IconButtonOptions{Disabled: true})
		})
	})
}
