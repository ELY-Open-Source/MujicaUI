package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

func TestTogglePressAndRelease(t *testing.T) {
	on, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Toggle(c, &on, "Bold", ToggleOptions{}).Changed() {
			changes++
		}
	}, 300, 120)
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("click did not press it")
	}
	tt.Key(0, ui.KeySpace)
	if on || changes != 2 {
		t.Fatalf("on=%v changes=%d", on, changes)
	}
}

func TestToggleIconOnlyNamed(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Toggle(c, &on, "", ToggleOptions{Icon: icons.Must("bold"), Label: "Bold text"})
	}, 300, 120)
	if err := tt.Click("Bold text"); err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("icon toggle not pressed")
	}
}

func TestToggleDisabled(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Toggle(c, &on, "Italic", ToggleOptions{Disabled: true})
	}, 300, 120)
	_ = tt.Click("Italic")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if on {
		t.Fatal("disabled toggle changed")
	}
}

func TestToggleNeedsName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a nameless toggle")
		}
	}()
	on := false
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		Toggle(c, &on, "", ToggleOptions{})
	}, 100, 100, 1)
}

func TestToggleSnapshot(t *testing.T) {
	snapshot(t, "toggle", 480, 90, func(c *ui.Context) {
		on, off := true, false
		ui.Row(c).Padding(20).Gap(12).Children(func() {
			Toggle(c, &on, "Bold", ToggleOptions{Icon: icons.Must("bold")})
			Toggle(c, &off, "Italic", ToggleOptions{Icon: icons.Must("italic")})
			Toggle(c, &on, "", ToggleOptions{Icon: icons.Must("underline"), Label: "Underline"})
			Toggle(c, &on, "Locked", ToggleOptions{Disabled: true})
		})
	})
}
