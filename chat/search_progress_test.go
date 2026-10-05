package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func searchStages(done int) []SearchStage {
	labels := []string{"Planning queries", "Searching the web", "Reading sources", "Writing the answer"}
	out := make([]SearchStage, len(labels))
	for i, l := range labels {
		out[i] = SearchStage{Label: l, Done: i < done}
	}
	return out
}

func TestSearchProgressFoldsWhenDone(t *testing.T) {
	opts := SearchProgressOptions{Query: "Versailles 1682", Read: 3, Total: 8, Stages: searchStages(2)}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SearchProgress(c, opts)
	}, 500, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for _, s := range []string{"Searching for “Versailles 1682”", "Read 3 of 8 sources", "Reading sources"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	opts.Done, opts.Read, opts.Stages = true, 8, searchStages(4)
	tt.Frame()
	if !tt.HasText("Searched “Versailles 1682”, read 8 sources") || tt.HasText("Reading sources") {
		t.Fatalf("done: %q", tt.Texts())
	}
	_ = tt.Click("Searched “Versailles 1682”, read 8 sources")
	if !tt.HasText("Reading sources") {
		t.Fatal("finished search does not open again")
	}
}

func TestSearchProgressBadCounts(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for Read above Total")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SearchProgress(c, SearchProgressOptions{Query: "q", Read: 3, Total: 1})
	}, 100, 100)
}

func TestSearchProgressReopensForNextSearch(t *testing.T) {
	opts := SearchProgressOptions{Query: "first", Read: 1, Total: 2, Stages: []SearchStage{{Label: "First stage"}}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SearchProgress(c, opts)
	}, 400, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	opts.Done = true
	tt.Frame()
	opts = SearchProgressOptions{Query: "second", Total: 3, Stages: []SearchStage{{Label: "Next stage"}}}
	tt.Frame()
	if !tt.HasText("Next stage") || !tt.HasText("Read 0 of 3 sources") {
		t.Fatalf("next search remains folded: %q", tt.Texts())
	}
}

func TestSearchProgressSnapshot(t *testing.T) {
	snapshot(t, "search_progress", 440, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			SearchProgress(c, SearchProgressOptions{Query: "Versailles 1682", Read: 3, Total: 8, Stages: searchStages(2)})
			SearchProgress(c, SearchProgressOptions{Query: "Hall of Mirrors", Read: 5, Total: 5, Stages: searchStages(4), Done: true})
		})
	})
}
