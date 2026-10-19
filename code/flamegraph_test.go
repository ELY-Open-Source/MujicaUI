package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func flameSample() *FlameNode {
	return &FlameNode{Name: "main", Value: 100, Children: []FlameNode{
		{Name: "court.Seal", Value: 60, Children: []FlameNode{
			{Name: "wax.Melt", Value: 40},
			{Name: "ink.Dry", Value: 15},
		}},
		{Name: "ledger.Write", Value: 30},
	}}
}

type flameProbe struct {
	root    *FlameNode
	s       FlamegraphState
	opts    FlamegraphOptions
	changed int
}

func (p *flameProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	if Flamegraph(c, p.root, &p.s, p.opts).Changed() {
		p.changed++
	}
}

func TestFlamegraphPaletteContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for i, col := range FlamegraphPalette(k) {
			if r := theme.Contrast(col, k.Surface); r < 3 {
				t.Errorf("dark=%v color %d contrast %.2f < 3", k.Dark, i, r)
			}
		}
	}
}

func TestFlamegraphClickZoomAndEscape(t *testing.T) {
	p := &flameProbe{root: flameSample(), opts: FlamegraphOptions{Unit: "ms"}}
	tt := ui.NewTester(p.view, 600, 300)
	g, _ := tt.Find("Flame graph")
	// Row 1, at 30% of the width, is court.Seal.
	tt.ClickAt(g.X+g.W*0.3, g.Y+flameRow*1.5)
	if len(p.s.Zoom) != 1 || p.s.Zoom[0] != 0 || p.changed != 1 {
		t.Fatalf("zoom %v changed %d", p.s.Zoom, p.changed)
	}
	if x, w := flameSpan(p.root, p.s.Zoom, []int{0, 0}); x != 0 || w < 0.66 || w > 0.67 {
		t.Fatalf("wax.Melt spans %v+%v of the zoomed width", x, w)
	}
	if !tt.HasText("court.Seal: 60.0 ms, 60.0% of total") {
		t.Fatalf("readout %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEscape)
	if len(p.s.Zoom) != 0 {
		t.Fatalf("Escape left zoom %v", p.s.Zoom)
	}
}

func TestFlamegraphKeyboardReadout(t *testing.T) {
	p := &flameProbe{root: flameSample(), opts: FlamegraphOptions{Unit: "ms"}}
	tt := ui.NewTester(p.view, 600, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Flame graph") {
		t.Fatal("Tab did not reach the graph")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if len(p.s.Focus) != 2 || p.s.Focus[1] != 1 {
		t.Fatalf("focus %v", p.s.Focus)
	}
	want := "ink.Dry: 15.0 ms, 15.0% of total"
	if !tt.HasText(want) {
		t.Fatalf("readout %q", tt.Texts())
	}
	found := false
	for _, a := range tt.Announcements() {
		found = found || a == want
	}
	if !found {
		t.Fatalf("readout not announced: %q", tt.Announcements())
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if len(p.s.Zoom) != 1 {
		t.Fatalf("Enter zoomed to %v", p.s.Zoom)
	}
}

func TestFlamegraphSearchAndStates(t *testing.T) {
	p := &flameProbe{root: flameSample(), s: FlamegraphState{Query: "wax"}}
	tt := ui.NewTester(p.view, 600, 300)
	if hit := flameHit(p.root, nil, 0.2, 2); len(hit) != 2 || flameAt(p.root, hit).Name != "wax.Melt" {
		t.Fatalf("hit %v", hit)
	}
	_ = tt
	p = &flameProbe{opts: FlamegraphOptions{Status: data.DataFailed, Error: "profile unreadable"}}
	tt = ui.NewTester(p.view, 600, 300)
	if !tt.HasText("profile unreadable") {
		t.Fatalf("texts %q", tt.Texts())
	}
	p = &flameProbe{}
	tt = ui.NewTester(p.view, 600, 300)
	if !tt.HasText("No profile") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestFlamegraphInvalidAndDisabled(t *testing.T) {
	p := &flameProbe{root: flameSample(), opts: FlamegraphOptions{Disabled: true}}
	tt := ui.NewTester(p.view, 600, 300)
	g, _ := tt.Find("Flame graph")
	tt.ClickAt(g.X+g.W*0.3, g.Y+flameRow*1.5)
	if len(p.s.Zoom) != 0 {
		t.Fatal("disabled graph zoomed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("children above their parent did not panic")
		}
	}()
	bad := &FlameNode{Name: "x", Value: 1, Children: []FlameNode{{Name: "y", Value: 2}}}
	ui.NewTester((&flameProbe{root: bad}).view, 600, 300)
}

func TestFlamegraphSnapshot(t *testing.T) {
	snapshot(t, "flamegraph", 560, 200, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *flameProbe {
			return &flameProbe{root: flameSample(), s: FlamegraphState{Focus: []int{0, 1}}, opts: FlamegraphOptions{Unit: "ms"}}
		})
		ui.Column(c).Padding(12).Children(func() { (*p).view(c) })
	})
}
