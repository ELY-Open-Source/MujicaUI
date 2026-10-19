package main

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func devPageTester(p page, mode core.Mode, events *[]string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: mode})
		ui.Column(c).Padding(16).Children(func() { p.Demo(c, func(s string) { *events = append(*events, s) }) })
	}, 1000, 800)
}

func TestDevtoolsPagesControlsAndThemes(t *testing.T) {
	if len(devtoolsPages) != 17 {
		t.Fatalf("pages=%d", len(devtoolsPages))
	}
	for _, p := range devtoolsPages {
		t.Run(p.Name, func(t *testing.T) {
			for _, mode := range []core.Mode{core.Light, core.Dark} {
				var events []string
				tt := devPageTester(p, mode, &events)
				if p.Category != "Dashboards & Dev Tools" || p.Code == "" || p.Keys == "" {
					t.Fatal("page metadata")
				}
				if err := tt.Click("Disabled"); err != nil {
					t.Fatal(err)
				}
				if len(events) != 1 || events[0] != "disabled changed" {
					t.Fatalf("control events %v", events)
				}
				if err := tt.Click("Disabled"); err != nil {
					t.Fatal(err)
				}
				if _, ok := tt.Find("Empty data"); ok {
					if err := tt.Click("Empty data"); err != nil {
						t.Fatal(err)
					}
					if events[len(events)-1] != "data state changed" {
						t.Fatal("empty event")
					}
					if err := tt.Click("Empty data"); err != nil {
						t.Fatal(err)
					}
				}
			}
		})
	}
}

func TestDevtoolsDemoEvents(t *testing.T) {
	cases := []struct{ id, label, event string }{
		{"443", "Acknowledge", "acknowledged cpu"},
		{"447", "Append line", "line appended"},
		{"450", "Append sample", "sample appended"},
		{"451", "Refresh", "refresh requested"},
		{"453", "Copy path", "path copied"},
		{"454", "Format", "JSON edited; valid=true"},
		{"456", "Send", "send requested"},
		{"458", "Add row", "pairs edited"},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			for _, p := range devtoolsPages {
				if p.ID != tc.id {
					continue
				}
				var events []string
				tt := devPageTester(p, core.Light, &events)
				if err := tt.Click(tc.label); err != nil {
					t.Fatal(err)
				}
				found := false
				for _, event := range events {
					found = found || event == tc.event
				}
				if !found {
					t.Fatalf("events %v want %s", events, tc.event)
				}
			}
		})
	}
}
