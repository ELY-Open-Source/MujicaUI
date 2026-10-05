package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func marketsDemo(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(marketsPages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("missing page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(24).Children(func() { marketsPages[i].Demo(c, func(s string) { *events = append(*events, s) }) })
	}, 1000, 800)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	return tt, events
}

func TestMarketsGalleryControls(t *testing.T) {
	if len(marketsPages) != 20 {
		t.Fatalf("pages %d", len(marketsPages))
	}
	for _, p := range marketsPages {
		t.Run(p.ID, func(t *testing.T) {
			if p.Category != "Markets" || p.Code == "" || p.Keys == "" {
				t.Fatalf("incomplete page %+v", p)
			}
			tt, events := marketsDemo(t, p.ID)
			if err := tt.Click("Disabled"); err != nil {
				t.Fatal(err)
			}
			if len(*events) == 0 || !strings.Contains((*events)[len(*events)-1], "disabled true") {
				t.Fatalf("toggle events %v", *events)
			}
			_ = tt.Click("Disabled")
			if tt.HasText("Error") {
				_ = tt.Click("Error")
				if !tt.HasText("Sample failed") {
					t.Fatalf("error state %q", tt.Texts())
				}
				_ = tt.Click("Empty")
				_ = tt.Click("Ready")
			}
		})
	}
}

func TestMarketsGalleryInteractionLogs(t *testing.T) {
	for _, tc := range []struct{ id, label string }{{"186", "AAPL"}, {"188", "Add condition"}, {"189", "S&P 500"}, {"190", "S&P 500"}, {"195", "Mark all read"}, {"197", "Retirement"}, {"200", "Swap currencies"}} {
		tt, events := marketsDemo(t, tc.id)
		before := len(*events)
		if err := tt.Click(tc.label); err != nil {
			t.Fatal(err)
		}
		if len(*events) == before {
			t.Fatalf("%s %s not logged", tc.id, tc.label)
		}
	}
}
