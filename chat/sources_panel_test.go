package chat

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestSourcesPanelOrderOpenAndScroll(t *testing.T) {
	var many []Source
	for i := range 20 {
		many = append(many, Source{Title: fmt.Sprintf("Source title %d", i+1), URL: fmt.Sprintf("https://example.com/%d", i+1)})
	}
	var st CitationState
	opened := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
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
		core.Use(c, core.Settings{})
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

func TestSourcesPanelLargeData(t *testing.T) {
	sources := make([]Source, 100000)
	for i := range sources {
		sources[i] = Source{Title: fmt.Sprintf("Archive %d", i), URL: "https://example.com"}
	}
	st := CitationState{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(360).Height(300).Children(func() { SourcesPanel(c, sources, &st) })
	}, 400, 300)
	if n := len(tt.Texts()); n > 250 {
		t.Fatalf("built %d labels for 100000 sources, want <= 250", n)
	}
	st.Active = 99999
	tt.Frame()
	r, ok := tt.Find("Archive 99998")
	if !ok || r.Y < 0 || r.Y+r.H > 300 {
		t.Fatalf("active source bounds %v found %v", r, ok)
	}
}

func TestSourcesPanelRepeatsCitationLocate(t *testing.T) {
	sources := make([]Source, 100)
	for i := range sources {
		sources[i] = Source{Title: fmt.Sprintf("Archive %d", i+1), URL: "https://example.com"}
	}
	var st CitationState
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Fill().Children(func() {
			if CitationBadge(c, 90, sources[89], &st).Changed() {
				actions++
			}
			ui.Column(c).Height(250).Children(func() { SourcesPanel(c, sources, &st) })
		})
	}, 400, 300)
	for i := range 2 {
		if err := tt.Click("Source 90: Archive 90"); err != nil {
			t.Fatal(err)
		}
		b, ok := tt.Find("Archive 90")
		if !ok || b.H == 0 || b.Y < 0 || b.Y+b.H > 300 || actions != i+1 {
			t.Fatalf("locate %d: bounds %v found %v events %d", i+1, b, ok, actions)
		}
		tt.Scroll(100, 150, 0, -100000)
	}
}
