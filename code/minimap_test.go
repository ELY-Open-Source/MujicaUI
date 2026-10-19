package code

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func minimapLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("\tvalue%d := compute(%d) // note", i, i)
	}
	return lines
}

func minimapTester(lines []string, v *MinimapView, disabled bool, moves *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := Minimap(c, lines, v, MinimapOptions{Language: "go", Disabled: disabled})
		r.Element.Height(200)
		if r.Changed() {
			*moves++
		}
	}, 200, 220)
}

func TestMinimapClickAndDragJump(t *testing.T) {
	v := &MinimapView{Visible: 20}
	moves := 0
	tt := minimapTester(minimapLines(80), v, false, &moves)
	// y=100 is line 50; centered, top is 40.
	tt.Press(40, 100)
	if v.Top != 40 {
		t.Fatalf("press moved top to %d, want 40", v.Top)
	}
	tt.Move(40, 150)
	tt.Release(40, 150)
	if v.Top != 60 {
		t.Fatalf("drag clamped top to %d, want 60 (80-20)", v.Top)
	}
	if moves < 2 {
		t.Fatalf("Changed reported %d times", moves)
	}
}

func TestMinimapKeys(t *testing.T) {
	v := &MinimapView{Visible: 20}
	moves := 0
	tt := minimapTester(minimapLines(500), v, false, &moves)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Minimap") {
		t.Fatal("Tab did not focus the minimap")
	}
	tt.Key(0, ui.KeyPageDown)
	tt.Key(0, ui.KeyDown)
	if v.Top != 21 {
		t.Fatalf("PageDown+Down gave %d", v.Top)
	}
	tt.Key(0, ui.KeyEnd)
	if v.Top != 480 {
		t.Fatalf("End gave %d", v.Top)
	}
	tt.Key(0, ui.KeyHome)
	if v.Top != 0 {
		t.Fatalf("Home gave %d", v.Top)
	}
}

func TestMinimapScrollsWhenTaller(t *testing.T) {
	if off := minimapOffset(1000, MinimapView{Top: 980, Visible: 20}, 200); off != 1800 {
		t.Fatalf("offset at the end %v, want 1800", off)
	}
	if off := minimapOffset(50, MinimapView{Top: 10, Visible: 20}, 200); off != 0 {
		t.Fatalf("short map offset %v", off)
	}
}

func TestMinimapDisabled(t *testing.T) {
	v := &MinimapView{Visible: 20}
	moves := 0
	tt := minimapTester(minimapLines(80), v, true, &moves)
	tt.Press(40, 100)
	tt.Release(40, 100)
	if v.Top != 0 || moves != 0 {
		t.Fatalf("disabled minimap moved to %d", v.Top)
	}
}

func TestMinimapNegativePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("negative Visible did not panic")
		}
	}()
	minimapTester(nil, &MinimapView{Visible: -1}, false, new(int))
}

func TestMinimapSnapshot(t *testing.T) {
	snapshot(t, "minimap", 140, 220, func(c *ui.Context) {
		v := ui.Local(c.Root(), "snap", func() *MinimapView { return &MinimapView{Top: 12, Visible: 24} })
		ui.Column(c).Padding(12).Children(func() {
			Minimap(c, codeMinimapSnapshotLines(), *v, MinimapOptions{Language: "go"}).Element.Height(196)
		})
	})
}

func codeMinimapSnapshotLines() []string {
	var lines []string
	for i := range 12 {
		lines = append(lines, codeViewerSample...)
		lines = append(lines, fmt.Sprintf("// section %d", i), "")
	}
	return lines
}
