package main

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func layoutPageTester(t *testing.T, id string, events *[]string) *ui.Tester {
	t.Helper()
	for _, p := range layoutPages {
		if p.ID == id {
			return ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Scroll(c).Fill().Children(func() {
					p.Demo(c, func(s string) { *events = append(*events, s) })
				})
			}, 900, 700)
		}
	}
	t.Fatalf("no layout page %s", id)
	return nil
}

func TestLayoutStackDemoVerticalHasNoWrap(t *testing.T) {
	var events []string
	tt := layoutPageTester(t, "005", &events)
	if err := tt.Click("Wrap"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Horizontal"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Wrap") {
		t.Fatal("Wrap control offered for a vertical stack")
	}
	if err := tt.Click("Horizontal"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Wrap") {
		t.Fatal("Wrap control missing for a horizontal stack")
	}
	if len(events) != 3 || events[1] != "Horizontal: false" {
		t.Fatalf("events %q", events)
	}
}
