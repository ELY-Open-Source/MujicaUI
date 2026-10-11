package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestThemeSelectorPicksByPointerAndKeys(t *testing.T) {
	mode, n := core.System, 0
	disabled := false
	tt := accountTester(600, 220, func(c *ui.Context) {
		if ThemeSelector(c, &mode, ThemeSelectorOptions{Disabled: disabled}).Changed() {
			n++
		}
	})
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	if mode != core.Dark || n != 1 {
		t.Fatalf("click: %v %d", mode, n)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if mode != core.Light || n != 2 {
		t.Fatalf("keys: %v %d", mode, n)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("System")
	if mode != core.Light {
		t.Fatal("disabled selector changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown mode accepted")
		}
	}()
	bad := core.Mode(9)
	accountTester(100, 100, func(c *ui.Context) { ThemeSelector(c, &bad, ThemeSelectorOptions{}) })
}

func TestThemeSelectorSnapshot(t *testing.T) {
	snapshot(t, "theme_selector", 520, 170, func(c *ui.Context) {
		m := ui.Local(c.Root(), "m", func() core.Mode { return core.Dark })
		ui.Column(c).Padding(16).Children(func() { ThemeSelector(c, m, ThemeSelectorOptions{}) })
	})
}
