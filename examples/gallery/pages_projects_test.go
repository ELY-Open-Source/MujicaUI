package main

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func projectsGalleryTester(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	var events []string
	for _, p := range projectsPages {
		if p.ID == id {
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Column(c).Fill().Gap(12).Padding(16).Children(func() { p.Demo(c, func(s string) { events = append(events, s) }) })
			}, 1200, 1000)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			return tt, &events
		}
	}
	t.Fatal("unknown project page " + id)
	return nil, nil
}

func TestProjectsGalleryCatalog(t *testing.T) {
	if len(projectsPages) != 18 {
		t.Fatalf("pages = %d", len(projectsPages))
	}
	seen := map[string]bool{}
	for _, p := range projectsPages {
		if p.Category != "Projects" || p.Keys == "" || p.Code == "" || seen[p.ID] {
			t.Fatalf("invalid page %s", p.ID)
		}
		seen[p.ID] = true
		projectsGalleryTester(t, p.ID)
	}
	if seen["366"] || seen["367"] {
		t.Fatal("dropped components have gallery pages")
	}
}

func TestProjectsGalleryEvents(t *testing.T) {
	for _, tc := range []struct{ id, action, event string }{
		{"361", "ENG-1 Plan the release", "open ENG-1"},
		{"362", "Done 1", "filter = 2"},
		{"363", "Plan the release", "open ENG-1"},
		{"364", "Close details", "close details"},
		{"365", "Clear subtasks", "clear subtasks"},
		{"368", "ENG-1 Plan the release", "open ENG-1"},
		{"369", "Ideal", "toggle Ideal"},
		{"370", "Icon only", "Icon only = true"},
		{"371", "Status", ""},
		{"372", "Disabled", "Disabled = true"},
		{"373", "Today", "due "},
		{"374", "UI", "labels changed:"},
		{"375", "ENG-1 Plan the release", "open ENG-1"},
		{"376", "Copy ENG-123", "copied ENG-123"},
		{"377", "Completed tasks", ""},
		{"378", "Ben assigned hours", ""},
		{"379", "Start", "tracked "},
		{"380", "Start", "timer phase = 0, running = true"},
	} {
		t.Run(tc.id, func(t *testing.T) {
			tt, events := projectsGalleryTester(t, tc.id)
			if err := tt.Click(tc.action); err != nil {
				t.Fatal(err)
			}
			switch tc.id {
			case "371":
				tt.Key(0, ui.KeyEnd)
				tt.Key(0, ui.KeyEnter)
				tc.event = "status done"
			case "377", "378":
				tt.Command("selectAll")
				tt.Type("5")
				tt.Key(0, ui.KeyEnter)
				if tc.id == "377" {
					tc.event = "completed = 5"
				} else {
					tc.event = "assigned = 5"
				}
			}
			found := false
			for _, e := range *events {
				if strings.HasPrefix(e, tc.event) {
					found = true
				}
			}
			if !found {
				t.Fatalf("event %q missing in %v", tc.event, *events)
			}
		})
	}
}
