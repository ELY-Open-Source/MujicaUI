package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSwitchToggles(t *testing.T) {
	on, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Switch(c, &on, "Wi-Fi", SwitchOptions{}).Changed() {
			changes++
		}
	}, 300, 120)
	if err := tt.Click("Wi-Fi"); err != nil {
		t.Fatal(err)
	}
	if !on || changes != 1 {
		t.Fatalf("on=%v changes=%d after click", on, changes)
	}
	tt.Key(0, ui.KeySpace)
	if on || changes != 2 {
		t.Fatalf("on=%v changes=%d after Space", on, changes)
	}
}

func TestSwitchKeyboardFocus(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Switch(c, &on, "Sync", SwitchOptions{})
	}, 300, 120)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Sync") {
		t.Fatal("Tab did not focus the switch")
	}
	tt.Key(0, ui.KeySpace)
	if !on {
		t.Fatal("Space did not turn it on")
	}
}

func TestSwitchBlockedStates(t *testing.T) {
	for _, opts := range []SwitchOptions{{Disabled: true}, {Loading: true}} {
		on := true
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			Switch(c, &on, "Backup", opts)
		}, 300, 120)
		_ = tt.Click("Backup")
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		if !on {
			t.Fatalf("%+v: value changed", opts)
		}
		if !tt.HasText("Backup") {
			t.Fatalf("%+v: label lost", opts)
		}
	}
}

func TestSwitchSnapshot(t *testing.T) {
	snapshot(t, "switch", 480, 120, func(c *ui.Context) {
		on, off := true, false
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			ui.Row(c).Gap(24).Children(func() {
				Switch(c, &on, "On", SwitchOptions{})
				Switch(c, &off, "Off", SwitchOptions{})
				Switch(c, &on, "Saving", SwitchOptions{Loading: true})
			})
			ui.Row(c).Gap(24).Children(func() {
				Switch(c, &on, "Disabled on", SwitchOptions{Disabled: true})
				Switch(c, &off, "Disabled off", SwitchOptions{Disabled: true})
			})
		})
	})
}
