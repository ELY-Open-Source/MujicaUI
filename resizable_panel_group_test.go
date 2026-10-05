package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func panelTester(sizes *[]float32, opts *ResizablePanelGroupOptions, commits *int) *ui.Tester {
	panels := []ResizablePanel{
		{Label: "Files", Min: 100, Collapsible: true, Build: func() {}},
		{Label: "Editor", Min: 150, Build: func() {}},
		{Label: "Outline", Min: 80, Build: func() {}},
	}
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := ResizablePanelGroup(c, sizes, panels, *opts)
		r.Element.Size(602, 300)
		if r.Committed {
			*commits++
		}
	}, 700, 400)
}

func near(a, b float32) bool { return math.Abs(float64(a-b)) < 0.002 }

func TestResizablePanelGroupKeyboard(t *testing.T) {
	var sizes []float32
	opts := ResizablePanelGroupOptions{}
	commits := 0
	tt := panelTester(&sizes, &opts, &commits)
	tt.Frame()
	if len(sizes) != 3 || !near(sizes[0], 1.0/3) {
		t.Fatalf("default sizes %v", sizes)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Resize Files and Editor") {
		t.Fatal("first divider not focused")
	}
	tt.Key(0, ui.KeyRight) // 16 DIPs of 600
	if !near(sizes[0], (200.0+16)/600) || !near(sizes[1], (200.0-16)/600) || commits != 1 {
		t.Fatalf("Right: %v (%d)", sizes, commits)
	}
	tt.Key(0, ui.KeyEnd) // Editor keeps its 150 minimum
	if !near(sizes[1], 150.0/600) {
		t.Fatalf("End: %v", sizes)
	}
	tt.Key(0, ui.KeyEnter)
	if sizes[0] != 0 || !near(sizes[1], 1-sizes[2]) {
		t.Fatalf("collapse: %v", sizes)
	}
	tt.Key(0, ui.KeyEnter)
	if !near(sizes[0], 250.0/600) {
		t.Fatalf("restore: %v", sizes)
	}
	opts.Disabled = true
	tt.Frame()
	before := sizes[0]
	tt.Key(0, ui.KeyLeft)
	if sizes[0] != before {
		t.Fatal("disabled group resized")
	}
}

func TestResizablePanelGroupDrag(t *testing.T) {
	sizes := []float32{0.5, 0.25, 0.25}
	opts := ResizablePanelGroupOptions{}
	commits := 0
	tt := panelTester(&sizes, &opts, &commits)
	tt.Frame()
	x := float32(300) // the first boundary
	tt.Press(x, 100)
	tt.Move(x-60, 100)
	if commits != 0 || !near(sizes[0], 240.0/600) {
		t.Fatalf("while dragging: %v (%d)", sizes, commits)
	}
	tt.Release(x-60, 100)
	if commits != 1 {
		t.Fatalf("release did not commit: %d", commits)
	}
	tt.Press(x-60, 100)
	tt.Move(x-60-200, 100) // past half of Files' minimum: it collapses
	tt.Release(x-60-200, 100)
	if sizes[0] != 0 {
		t.Fatalf("drag collapse: %v", sizes)
	}
}

func TestResizablePanelGroupRejects(t *testing.T) {
	ok := ResizablePanel{Label: "A", Build: func() {}}
	for name, tc := range map[string]struct {
		sizes  []float32
		panels []ResizablePanel
	}{
		"one panel":  {nil, []ResizablePanel{ok}},
		"count":      {[]float32{1}, []ResizablePanel{ok, ok}},
		"sum":        {[]float32{0.5, 0.2}, []ResizablePanel{ok, ok}},
		"zero fixed": {[]float32{0, 1}, []ResizablePanel{ok, ok}},
		"no label":   {nil, []ResizablePanel{ok, {Build: func() {}}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			s := tc.sizes
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); ResizablePanelGroup(c, &s, tc.panels, ResizablePanelGroupOptions{}) }, 300, 200)
		}()
	}
}

func TestResizablePanelGroupSnapshot(t *testing.T) {
	sizes := []float32{0.25, 0.5, 0.25}
	vs := []float32{0.6, 0.4}
	snapshot(t, "resizable_panel_group", 560, 260, func(c *ui.Context) {
		k := Tokens(c)
		pane := func(s string) func() {
			return func() {
				ui.Column(c).Fill().Padding(12).Background(k.Surface).Children(func() { ui.Text(c, s) })
			}
		}
		ui.Row(c).Padding(16).Children(func() {
			ResizablePanelGroup(c, &sizes, []ResizablePanel{
				{Label: "Files", Build: pane("Files")},
				{Label: "Editor", Build: func() {
					ResizablePanelGroup(c, &vs, []ResizablePanel{{Label: "Code", Build: pane("Code")}, {Label: "Terminal", Build: pane("Terminal")}},
						ResizablePanelGroupOptions{Vertical: true}).Element.Fill()
				}},
				{Label: "Outline", Build: pane("Outline")},
			}, ResizablePanelGroupOptions{}).Element.Size(528, 228)
		})
	})
}
