package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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

// switchGaps measures the edge-to-thumb gaps of a switch's track.
func switchGaps(t *testing.T, on bool) (left, right, top, bottom int) {
	t.Helper()
	var k theme.Tokens
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		k = core.Tokens(c)
		ui.Column(c).Padding(10).Children(func() { Switch(c, &on, "Wi-Fi", SwitchOptions{}) })
	}, 300, 120)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	a, ok := tt.Access("Wi-Fi")
	if !ok {
		t.Fatal("no switch")
	}
	b := a.Bounds
	x0, y0 := int(b.X), int(b.Y+(b.H-switchH)/2)
	thumb, bg, from := k.ControlBorder, k.Surface, x0+1
	if on {
		thumb, bg, from = k.OnAccent, k.Accent, x0+12 // past the "on" bar
	}
	img := tt.Image()
	near := func(x, y int, c ui.Color) int {
		p := img.RGBAAt(x, y)
		return abs(int(p.R)-int(c.R)) + abs(int(p.G)-int(c.G)) + abs(int(p.B)-int(c.B))
	}
	isThumb := func(x, y int) bool { return near(x, y, thumb) < near(x, y, bg) }
	row, minX, maxX := y0+int(switchH)/2, -1, -1
	for x := from; x < x0+int(switchW)-1; x++ {
		if isThumb(x, row) {
			if minX < 0 {
				minX = x
			}
			maxX = x
		}
	}
	col, minY, maxY := (minX+maxX)/2, -1, -1
	for y := y0 + 1; y < y0+int(switchH)-1; y++ {
		if isThumb(col, y) {
			if minY < 0 {
				minY = y
			}
			maxY = y
		}
	}
	return minX - x0, x0 + int(switchW) - 1 - maxX, minY - y0, y0 + int(switchH) - 1 - maxY
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func TestSwitchThumbCenteredInsideBorder(t *testing.T) {
	for _, on := range []bool{false, true} {
		l, r, top, bottom := switchGaps(t, on)
		side := l
		if on {
			side = r
		}
		if top != bottom || side != top {
			t.Errorf("on %v: gaps left %d right %d top %d bottom %d", on, l, r, top, bottom)
		}
	}
}
