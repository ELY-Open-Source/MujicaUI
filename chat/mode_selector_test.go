package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestModeSelectorClickAndArrows(t *testing.T) {
	mode, n, off := ModeChat, 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ModeSelector(c, &mode, ModeSelectorOptions{Disabled: off}).Changed() {
			n++
		}
	}, 600, 100)
	_ = tt.Click("Deep research")
	if mode != ModeResearch || n != 1 {
		t.Fatalf("click: mode %d changes %d", mode, n)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if mode != ModeAgent {
		t.Fatalf("Right: mode %d", mode)
	}
	off = true
	_ = tt.Click("Chat")
	if mode != ModeAgent {
		t.Fatal("disabled selector changed")
	}
}

func TestModeSelectorIconsOnlyAndBadMode(t *testing.T) {
	mode := ModeSearch
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ModeSelector(c, &mode, ModeSelectorOptions{IconsOnly: true})
	}, 300, 100)
	if _, ok := tt.Find("Search"); !ok {
		t.Fatal("icon segments unnamed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an unknown mode")
		}
	}()
	bad := ChatMode(9)
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ModeSelector(c, &bad, ModeSelectorOptions{})
	}, 100, 100)
}

func TestModeSelectorSnapshot(t *testing.T) {
	mode := ModeResearch
	snapshot(t, "mode_selector", 520, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
			ModeSelector(c, &mode, ModeSelectorOptions{})
			ModeSelector(c, &mode, ModeSelectorOptions{IconsOnly: true})
		})
	})
}

func TestModeSelectorSelectionSurface(t *testing.T) {
	mode := ModeResearch
	tt := ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); ModeSelector(c, &mode, ModeSelectorOptions{}) }, 600, 100)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	b, _ := tt.Find("Deep research")
	pixel := tt.Image().RGBAAt(int(b.X+4), int(b.Y+4))
	k := theme.Light()
	if pixel.R != k.Selection.R || pixel.G != k.Selection.G || pixel.B != k.Selection.B {
		t.Fatalf("selected surface %v, want Selection %v", pixel, k.Selection)
	}
}
