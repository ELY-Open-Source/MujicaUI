package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSourcesPanelOrderOpenAndScroll(t *testing.T) {
	var many []Source
	for i := range 20 {
		many = append(many, Source{Title: fmt.Sprintf("Source title %d", i+1), URL: fmt.Sprintf("https://example.com/%d", i+1)})
	}
	var st CitationState
	opened := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(360).Height(300).Children(func() {
			if n, ok := SourcesPanel(c, many, &st).Opened(); ok {
				opened = n
			}
		})
	}, 400, 300)
	first, _ := tt.Find("Source title 1")
	second, _ := tt.Find("Source title 2")
	if !tt.HasText("Sources (20)") || second.Y <= first.Y {
		t.Fatalf("order or header: %q", tt.Texts())
	}
	_ = tt.Click("Open Source title 2")
	if opened != 2 {
		t.Fatalf("opened %d", opened)
	}
	st.Active = 18
	tt.Frame()
	tt.Frame()
	r, ok := tt.Find("Source title 18")
	if !ok || r.Y < 0 || r.Y+r.H > 300 {
		t.Fatalf("active source not scrolled into view: %v %v", r, ok)
	}
}

func TestSourcesPanelEmptyAndBadActive(t *testing.T) {
	var st CitationState
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		SourcesPanel(c, nil, &st)
	}, 300, 200)
	if !tt.HasText("No sources") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an Active past the list")
		}
	}()
	st.Active = 4
	tt.Frame()
}

func TestSourcesPanelSnapshot(t *testing.T) {
	snapshot(t, "sources_panel", 380, 420, func(c *ui.Context) {
		st := ui.Local(c.Root(), "st", func() CitationState { return CitationState{Active: 2} })
		ui.Column(c).Padding(16).Height(420).Children(func() { SourcesPanel(c, sourcesSample, st) })
	})
}
