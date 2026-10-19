package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

type breakpointProbe struct {
	bps                     []Breakpoint
	disabled                bool
	toggled, edited, opened []string
}

func (p *breakpointProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := BreakpointList(c, &p.bps, BreakpointListOptions{Disabled: p.disabled})
	r.Element.Height(300)
	for _, e := range []struct {
		id  string
		out *[]string
	}{{r.Toggled, &p.toggled}, {r.Edited, &p.edited}, {r.Opened, &p.opened}} {
		if e.id != "" {
			*e.out = append(*e.out, e.id)
		}
	}
}

func newBreakpointProbe() *breakpointProbe {
	return &breakpointProbe{bps: []Breakpoint{
		{ID: "a", File: "court/seal.go", Line: 12, Enabled: true},
		{ID: "b", File: "court/ledger.go", Line: 40, Enabled: false, Condition: "n > 3"},
	}}
}

func TestBreakpointListToggleAndOpen(t *testing.T) {
	p := newBreakpointProbe()
	tt := ui.NewTester(p.view, 520, 320)
	if err := tt.Click("Enable ledger.go:40"); err != nil {
		t.Fatal(err)
	}
	if !p.bps[1].Enabled || len(p.toggled) != 1 || p.toggled[0] != "b" {
		t.Fatalf("toggle: %+v %v", p.bps[1], p.toggled)
	}
	if err := tt.Click("court/seal.go:12"); err != nil {
		t.Fatal(err)
	}
	if len(p.opened) != 1 || p.opened[0] != "a" {
		t.Fatalf("opened %v", p.opened)
	}
	if !tt.HasText("when n > 3") {
		t.Fatalf("condition not shown: %q", tt.Texts())
	}
}

func TestBreakpointListEditCondition(t *testing.T) {
	p := newBreakpointProbe()
	tt := ui.NewTester(p.view, 520, 320)
	if err := tt.Click("Edit condition of seal.go:12"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Condition") {
		t.Fatal("editor did not take the focus")
	}
	tt.Type("i == 7")
	tt.Key(0, ui.KeyEnter)
	if p.bps[0].Condition != "i == 7" || len(p.edited) != 1 {
		t.Fatalf("condition %q edited %v", p.bps[0].Condition, p.edited)
	}
	_ = tt.Click("Edit condition of ledger.go:40")
	tt.Type("x")
	tt.Key(0, ui.KeyEscape)
	if p.bps[1].Condition != "n > 3" {
		t.Fatalf("Escape kept %q", p.bps[1].Condition)
	}
}

func TestBreakpointListKeyboardToggle(t *testing.T) {
	p := newBreakpointProbe()
	tt := ui.NewTester(p.view, 520, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if p.bps[0].Enabled {
		t.Fatal("Space did not disable the first breakpoint")
	}
}

func TestBreakpointListDisabledAndDuplicate(t *testing.T) {
	p := newBreakpointProbe()
	p.disabled = true
	tt := ui.NewTester(p.view, 520, 320)
	_ = tt.Click("Enable ledger.go:40")
	_ = tt.Click("court/seal.go:12")
	if p.bps[1].Enabled || len(p.opened) != 0 {
		t.Fatal("disabled list changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate IDs did not panic")
		}
	}()
	p = newBreakpointProbe()
	p.bps[1].ID = "a"
	ui.NewTester(p.view, 520, 320)
}

func TestBreakpointListSnapshot(t *testing.T) {
	snapshot(t, "breakpoint_list", 460, 170, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", newBreakpointProbe)
		ui.Column(c).Padding(12).Children(func() {
			BreakpointList(c, &(*p).bps, BreakpointListOptions{}).Element.Height(146)
		})
	})
}
