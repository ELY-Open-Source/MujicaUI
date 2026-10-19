package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func variablesSample() []Variable {
	return []Variable{
		{Name: "decree", Type: "court.Decree", Value: "{…}", Children: []Variable{
			{Name: "id", Value: "42", Type: "int", Changed: true},
			{Name: "title", Value: `"Of seals"`, Type: "string"},
		}},
		{Name: "count", Value: "7", Type: "int"},
	}
}

type variablesProbe struct {
	vars     []Variable
	s        VariablesPanelState
	disabled bool
	edits    []VariableEdit
}

func (p *variablesProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := VariablesPanel(c, p.vars, &p.s, VariablesPanelOptions{Disabled: p.disabled})
	r.Element.Height(300)
	if r.Edited != nil {
		p.edits = append(p.edits, *r.Edited)
	}
}

func TestVariablesPanelExpandAndChanged(t *testing.T) {
	p := &variablesProbe{vars: variablesSample()}
	tt := ui.NewTester(p.view, 460, 320)
	if tt.HasText("title") {
		t.Fatal("children shown before expanding")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("title") || !tt.HasText("42") {
		t.Fatalf("Right did not expand: %q", tt.Texts())
	}
	if !p.s.byPath["decree.id"].Changed {
		t.Fatal("changed flag lost")
	}
}

func TestVariablesPanelEditLeaf(t *testing.T) {
	p := &variablesProbe{vars: variablesSample()}
	tt := ui.NewTester(p.view, 460, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyF2)
	if !tt.Focused("Value") {
		t.Fatal("F2 did not open the editor")
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("8")
	tt.Key(0, ui.KeyEnter)
	if len(p.edits) != 1 || p.edits[0] != (VariableEdit{"count", "8"}) {
		t.Fatalf("edits %v", p.edits)
	}
	// Enter on a branch does not edit it.
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyEnter)
	if tt.Focused("Value") {
		t.Fatal("a branch opened the editor")
	}
}

func TestVariablesPanelEscapeAndDisabled(t *testing.T) {
	p := &variablesProbe{vars: variablesSample()}
	tt := ui.NewTester(p.view, 460, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	tt.Type("x")
	tt.Key(0, ui.KeyEscape)
	if len(p.edits) != 0 || tt.Focused("Value") {
		t.Fatal("Escape committed")
	}
	p = &variablesProbe{vars: variablesSample(), disabled: true}
	tt = ui.NewTester(p.view, 460, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	if tt.Focused("Value") {
		t.Fatal("disabled panel opened the editor")
	}
}

func TestVariablesPanelDuplicatePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate names did not panic")
		}
	}()
	variablesIndex([]Variable{{Name: "a"}, {Name: "a"}}, &VariablesPanelState{})
}

func TestVariablesPanelSnapshot(t *testing.T) {
	snapshot(t, "variables_panel", 420, 180, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *variablesProbe {
			p := &variablesProbe{vars: variablesSample()}
			p.s.outline.Open.Add("decree")
			return p
		})
		ui.Column(c).Padding(12).Children(func() {
			VariablesPanel(c, (*p).vars, &(*p).s, VariablesPanelOptions{}).Element.Height(156)
		})
	})
}
