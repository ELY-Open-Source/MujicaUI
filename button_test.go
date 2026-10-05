package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestButtonClickOnce(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		if Button(c, "Save", ButtonOptions{}).Clicked() {
			n++
		}
	}, 300, 200)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("clicked %d times", n)
	}
}

func TestButtonKeyboard(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Button(c, "Save", ButtonOptions{}).Clicked() {
			n++
		}
	}, 300, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if n != 2 {
		t.Fatalf("keyboard clicks %d, want 2", n)
	}
}

func TestButtonBlockedStates(t *testing.T) {
	for _, opts := range []ButtonOptions{{Disabled: true}, {Loading: true}} {
		n := 0
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			if Button(c, "Pay", opts).Clicked() {
				n++
			}
		}, 300, 200)
		_ = tt.Click("Pay")
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyEnter)
		if n != 0 {
			t.Fatalf("%+v: clicked %d times", opts, n)
		}
		if !tt.HasText("Pay") {
			t.Fatalf("%+v: label lost", opts)
		}
	}
}

func TestButtonMinHeight(t *testing.T) {
	for _, tc := range []struct {
		density theme.Density
		min     float32
	}{{theme.Compact, 28}, {theme.Default, 36}, {theme.Comfortable, 44}} {
		var h float32
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Density: tc.density})
			ui.Column(c).AlignItems(ui.Start).Children(func() {
				h = Button(c, "Go", ButtonOptions{}).Bounds().H
			})
		}, 300, 200)
		tt.Frame()
		if h < tc.min {
			t.Errorf("%v: height %v < %v", tc.density, h, tc.min)
		}
	}
}

func TestButtonSnapshot(t *testing.T) {
	snapshot(t, "button", 640, 140, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(12).Wrap().Children(func() {
			Button(c, "Primary", ButtonOptions{})
			Button(c, "Secondary", ButtonOptions{Variant: Secondary})
			Button(c, "Outline", ButtonOptions{Variant: Outline})
			Button(c, "Ghost", ButtonOptions{Variant: Ghost})
			Button(c, "Delete", ButtonOptions{Variant: Danger})
			Button(c, "Saving", ButtonOptions{Loading: true})
			Button(c, "Disabled", ButtonOptions{Disabled: true})
		})
	})
}
