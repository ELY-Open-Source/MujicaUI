package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func filesDemoTester(t *testing.T, id string, mode core.Mode) (*ui.Tester, *[]string) {
	t.Helper()
	for _, p := range filesPages {
		if p.ID == id {
			events := &[]string{}
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{Mode: mode})
				ui.Scroll(c).Fill().Padding(24).Gap(12).Children(func() { p.Demo(c, func(s string) { *events = append(*events, s) }) })
			}, 900, 800)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			return tt, events
		}
	}
	t.Fatalf("missing files page %s", id)
	return nil, nil
}

func TestFilesGalleryExistingSet(t *testing.T) {
	want := []string{"381", "382", "383", "384", "385", "386", "387", "388", "389", "390", "391", "392", "394", "395", "396", "399", "400"}
	if len(filesPages) != len(want) {
		t.Fatalf("files pages %d, want %d", len(filesPages), len(want))
	}
	for i, p := range filesPages {
		if p.ID != want[i] || p.Category != "Files & Collaboration" || p.Code == "" || p.Keys == "" {
			t.Fatalf("page %s has incorrect metadata", p.ID)
		}
	}
}

func TestFilesGalleryActions(t *testing.T) {
	cases := []struct{ id, action, event string }{
		{"381", "Disabled", "Disabled = true"},
		{"382", "Empty", "Empty = true"},
		{"383", "Folder", "Folder = true"},
		{"384", "Close preview", "Preview closed"},
		{"385", "court", "Navigate: /court"},
		{"386", "Rename", "Rename started"},
		{"387", "Cancel", "Operation canceled"},
		{"388", "Pause ledger.zip", "Transfer a: 1"},
		{"390", "Clear", "Recent files cleared"},
		{"391", "Empty", "Empty = true"},
		{"392", "Follow Ada Lovelace", "Follow: a"},
		{"394", "Clear reactions", "Reactions cleared"},
		{"395", "Share Ledger", "Share dialog opened"},
		{"396", "Disabled", "Disabled = true"},
		{"399", "Last sync time", "Last sync time = false"},
		{"400", "Selection", "Selection = false"},
	}
	for _, tc := range cases {
		for _, mode := range []core.Mode{core.Light, core.Dark} {
			t.Run(fmt.Sprintf("%s-%s", tc.id, mode), func(t *testing.T) {
				tt, events := filesDemoTester(t, tc.id, mode)
				if err := tt.Click(tc.action); err != nil {
					t.Fatal(err)
				}
				found := false
				for _, event := range *events {
					if event == tc.event {
						found = true
					}
				}
				if !found {
					t.Fatalf("events %v, want %s", *events, tc.event)
				}
				if tc.id == "395" {
					if err := tt.Click("Copy link"); err != nil {
						t.Fatal(err)
					}
					if tt.Clipboard() != "https://court.org/ledger" {
						t.Fatal("share callback did not copy link")
					}
				}
			})
		}
	}
}

func TestFilesGalleryStorageParameter(t *testing.T) {
	tt, events := filesDemoTester(t, "389", core.Light)
	if err := tt.Click("Photos (GB)"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("41")
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("77.0 GB of 128.0 GB used") || len(*events) == 0 {
		t.Fatalf("storage parameter left %q; events %v", tt.Texts(), *events)
	}
}

func TestFilesGalleryGridRange(t *testing.T) {
	tt, events := filesDemoTester(t, "382", core.Light)
	if err := tt.Click("notes.md"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.ClickWith(ui.Shift, "crest.png"); err != nil {
		t.Fatal(err)
	}
	last := (*events)[len(*events)-1]
	if strings.Contains(last, "/court/chronicles") || !strings.Contains(last, "/court/notes.md") || !strings.Contains(last, "/court/decrees/crest.png") {
		t.Fatalf("range selection event %q", last)
	}
}
