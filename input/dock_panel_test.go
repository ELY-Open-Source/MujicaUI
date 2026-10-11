package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func dockItems(c *ui.Context) []DockItem {
	var out []DockItem
	for _, id := range []string{"files", "editor", "terminal"} {
		out = append(out, DockItem{ID: id, Title: "Panel " + id, Build: func() { ui.Text(c, "Body of "+id) }})
	}
	return out
}

func dockTester(l *DockLayout, dis *bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := DockPanel(c, l, dockItems(c), DockPanelOptions{Disabled: *dis})
		r.Element.Size(800, 500)
		if r.Changed() {
			*n++
		}
	}, 820, 520)
}

func TestDockPanelTabsAndKeyboardDock(t *testing.T) {
	var l DockLayout
	dis, n := false, 0
	tt := dockTester(&l, &dis, &n)
	if !slices.Equal(l.Zones[DockCenter], []string{"files", "editor", "terminal"}) || l.Active[DockCenter] != "files" || !tt.HasText("Body of files") {
		t.Fatalf("default layout %+v", l)
	}
	if err := tt.Click("Panel editor"); err != nil {
		t.Fatal(err)
	}
	if l.Active[DockCenter] != "editor" || !tt.HasText("Body of editor") {
		t.Fatalf("tab click: %+v", l)
	}
	tt.Key(ui.Shift|ui.Alt, ui.KeyLeft)
	if !slices.Equal(l.Zones[DockLeft], []string{"editor"}) || slices.Contains(l.Zones[DockCenter], "editor") {
		t.Fatalf("Shift+Alt+Left: %+v", l)
	}
	if _, ok := tt.Find("Left dock"); !ok {
		t.Fatal("left zone not shown")
	}
	tt.Key(ui.Shift|ui.Alt, ui.KeyUp)
	if len(l.Zones[DockLeft]) != 0 || l.Active[DockCenter] != "editor" {
		t.Fatalf("back to center: %+v", l)
	}
	before := n
	dis = true
	tt.Frame()
	tt.Key(ui.Shift|ui.Alt, ui.KeyDown)
	if len(l.Zones[DockBottom]) != 0 || n != before {
		t.Fatalf("disabled docked: %+v", l)
	}
}

func TestDockPanelDragToEdge(t *testing.T) {
	var l DockLayout
	dis, n := false, 0
	tt := dockTester(&l, &dis, &n)
	r, _ := tt.Find("Panel terminal")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	for _, px := range []float32{x + 10, x + 200, 600, 790} {
		tt.Move(px, 250)
	}
	tt.Release(790, 250)
	tt.Frame()
	if !slices.Equal(l.Zones[DockRight], []string{"terminal"}) || n == 0 {
		t.Fatalf("drag to the right edge: %+v", l)
	}
}

func TestDockPanelRejects(t *testing.T) {
	for name, l := range map[string]DockLayout{
		"unknown":  {Zones: [4][]string{{"nope"}}},
		"repeated": {Zones: [4][]string{{"files"}, {"files"}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); DockPanel(c, &l, dockItems(c), DockPanelOptions{}) }, 300, 200)
		}()
	}
}

func TestDockPanelSnapshot(t *testing.T) {
	l := DockLayout{Zones: [4][]string{{"editor"}, {"files"}, nil, {"terminal"}}}
	snapshot(t, "dock_panel", 640, 360, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			DockPanel(c, &l, dockItems(c), DockPanelOptions{}).Element.Size(616, 336)
		})
	})
}
