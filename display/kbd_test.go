package display

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestKbdKeyPlatforms(t *testing.T) {
	for _, tc := range []struct {
		key           string
		mac           bool
		shown, spoken string
	}{
		{"Mod", true, "⌘", "Command"},
		{"Mod", false, "Ctrl", "Ctrl"},
		{"shift", true, "⇧", "Shift"},
		{"Alt", false, "Alt", "Alt"},
		{"Alt", true, "⌥", "Option"},
		{"Up", false, "↑", "Up"},
		{"k", false, "K", "K"},
		{"F5", true, "F5", "F5"},
	} {
		shown, spoken := KbdKey(tc.key, tc.mac)
		if shown != tc.shown || spoken != tc.spoken {
			t.Errorf("%q mac=%v: %q %q", tc.key, tc.mac, shown, spoken)
		}
	}
}

func TestKbdShowsPlatformCombo(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Kbd(c, "Mod", "Shift", "p")
	}, 300, 100)
	tt.Frame()
	spoken, mod := "Ctrl+Shift+P", "Ctrl"
	if runtime.GOOS == "darwin" {
		spoken, mod = "Command+Shift+P", "⌘"
	}
	if _, ok := tt.Find(spoken); !ok {
		t.Fatalf("not read as %q", spoken)
	}
	if !tt.HasText(mod) || !tt.HasText("P") {
		t.Fatalf("texts %v", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	if tt.Focused(spoken) {
		t.Fatal("Kbd took focus")
	}
}

func TestKbdEmptyPanics(t *testing.T) {
	for _, keys := range [][]string{nil, {"Mod", ""}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%q: no panic", keys)
				}
			}()
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Kbd(c, keys...)
			}, 100, 100)
			tt.Frame()
		}()
	}
}

func TestKbdSnapshot(t *testing.T) {
	snapshot(t, "kbd", 360, 80, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(16).AlignItems(ui.Center).Children(func() {
			Kbd(c, "Mod", "K")
			Kbd(c, "Mod", "Shift", "P")
			Kbd(c, "Esc")
			Kbd(c, "Alt", "Up")
		})
	})
}

func TestKbdTextScale(t *testing.T) {
	var h [2]float32
	for i, s := range []float32{1, 2} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Row(c).AlignItems(ui.Start).Children(func() {
				h[i] = Kbd(c, "Mod", "K").Bounds().H
			})
		}, 400, 200)
		tt.SetPreferences(ui.Preferences{TextScale: s})
		tt.Frame()
	}
	if h[1] < h[0]*1.5 {
		t.Fatalf("height %v at 2x, %v at 1x", h[1], h[0])
	}
}
