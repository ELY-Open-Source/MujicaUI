package main

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestMermaidDemoSwitchesSamplesAndLogs(t *testing.T) {
	var events []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(24).Children(func() {
			diagramPages[0].Demo(c, func(s string) { events = append(events, s) })
		})
	}, 900, 700)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if _, ok := tt.AccessByRole(ui.RoleImage, "Flowchart"); !ok || len(events) != 0 {
		t.Fatalf("first frame: flowchart %v, events %q", ok, events)
	}
	for _, name := range []string{"Sequence", "Pie"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if _, ok := tt.AccessByRole(ui.RoleImage, name); !ok {
			t.Errorf("%s is not drawn", name)
		}
	}
	if err := tt.Click("Broken"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if _, ok := tt.AccessByRole(ui.RoleImage, "Broken"); ok || !tt.HasText("Diagram failed") {
		t.Errorf("the broken sample: texts %q", tt.Texts())
	}
	if want := []string{"drew Sequence", "drew Pie", "drew Broken"}; !slices.Equal(events, want) {
		t.Errorf("events %q, want %q", events, want)
	}
	if n := len(catalog()); n != 488 || catalog()[n-1].ID != "501" {
		t.Errorf("catalog holds %d pages ending at %s, want 488 ending at 501", n, catalog()[n-1].ID)
	}
}
