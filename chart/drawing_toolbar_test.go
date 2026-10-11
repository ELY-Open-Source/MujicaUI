package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestDrawingToolbarShortcutsLockAndDisabled(t *testing.T) {
	d := DrawingTools{}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Disabled(disabled).Children(func() { DrawingToolbar(c, &d, DrawingToolbarOptions{Horizontal: true}) })
	}, 550, 90)
	tt.Key(ui.Alt, ui.Key4)
	if d.Tool != ToolFibonacci {
		t.Fatalf("tool %v", d.Tool)
	}
	if err := tt.Click("Lock drawings"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Alt, ui.Key3)
	if !d.Locked || d.Tool != ToolSelect {
		t.Fatal("locked tool changed")
	}
	if err := tt.Click("Unlock drawings"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	tt.Key(ui.Alt, ui.Key2)
	if d.Tool != ToolSelect {
		t.Fatal("disabled tool changed")
	}
}
func TestDrawingToolbarSnapshot(t *testing.T) {
	d := DrawingTools{Tool: ToolTrendLine}
	snapshot(t, "drawing_toolbar", 440, 90, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() { DrawingToolbar(c, &d, DrawingToolbarOptions{Horizontal: true}) })
	})
}
