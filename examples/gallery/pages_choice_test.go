package main

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// choiceDemoTester runs the demo of page id at text scale and returns its tester and event log.
func choiceDemoTester(t *testing.T, id string, scale float32) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(choicePages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("no page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Scroll(c).Fill().Padding(24).Children(func() {
			choicePages[i].Demo(c, func(s string) { *events = append(*events, s) })
		})
	}, 900, 700)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: scale})
	return tt, events
}

func TestChoiceSelectDemoKeepsPlaceholderChoice(t *testing.T) {
	tt, events := choiceDemoTester(t, "036", 1)
	if err := tt.Click("Pick a fruit"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Fig"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Fig") || !slices.Contains(*events, `Picked: "fig"`) {
		t.Fatalf("choice lost: texts %q, events %q", tt.Texts(), *events)
	}
}

func TestChoiceDemoHeadingsFollowTextScale(t *testing.T) {
	height := func(scale float32) float32 {
		tt, _ := choiceDemoTester(t, "031", scale)
		tt.Frame()
		r, ok := tt.Find("Parameters")
		if !ok {
			t.Fatal("heading missing")
		}
		return r.H
	}
	if one, two := height(1), height(2); two < one*1.5 {
		t.Fatalf("heading %v high at 2x, %v at 1x", two, one)
	}
}
