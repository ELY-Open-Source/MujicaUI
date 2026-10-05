package main

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// dataDemo runs the demo of page id and returns its tester and event log.
func dataDemo(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(dataPages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("no page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(24).Children(func() {
			dataPages[i].Demo(c, func(s string) { *events = append(*events, s) })
		})
	}, 900, 700)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt, events
}

func TestDataDemoButtonsLog(t *testing.T) {
	tt, events := dataDemo(t, "074")
	for _, s := range []string{"Bells", "Ring"} {
		if err := tt.Click(s); err != nil {
			t.Fatal(err)
		}
	}
	if !slices.Contains(*events, "ring the bells") {
		t.Fatalf("Ring not logged: %q", *events)
	}
	tt, events = dataDemo(t, "075")
	for _, s := range []string{"Open from outside", "Light candles"} {
		if err := tt.Click(s); err != nil {
			t.Fatal(err)
		}
	}
	if !slices.Contains(*events, "light candles") {
		t.Fatalf("Light candles not logged: %q", *events)
	}
}

func TestDataAccordionDemoLogsChange(t *testing.T) {
	tt, events := dataDemo(t, "074")
	if err := tt.Click("Bells"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(*events, "open [bells]") {
		t.Fatalf("Accordion change not logged: %q", *events)
	}
}
