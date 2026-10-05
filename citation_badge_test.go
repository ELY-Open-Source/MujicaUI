package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var sourcesSample = []Source{
	{Title: "The Palace of Versailles", URL: "https://example.com/versailles", Site: "example.com",
		Snippet: "Seat of the court from 1682.", Published: time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC)},
	{Title: "Hall of Mirrors", URL: "https://example.org/hall", Snippet: "Seventeen arches of mirrors face the garden."},
	{Title: "Court etiquette", URL: "https://example.net/etiquette", Site: "example.net"},
}

// citationApp is an answer with badges beside a sources panel, sharing one state.
type citationApp struct {
	st     CitationState
	opened []int
}

func (a *citationApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Row(c).Fill().Padding(16).Gap(16).AlignItems(ui.Start).Children(func() {
		ui.Row(c).Gap(4).Children(func() {
			ui.Text(c, "The court moved in 1682")
			CitationBadge(c, 1, sourcesSample[0], &a.st)
			CitationBadge(c, 3, sourcesSample[2], &a.st)
		})
		ui.Column(c).Width(300).Height(400).Children(func() {
			if n, ok := SourcesPanel(c, sourcesSample, &a.st).Opened(); ok {
				a.opened = append(a.opened, n)
			}
		})
	})
}

func TestCitationBadgeActivatesAndPreviews(t *testing.T) {
	a := &citationApp{}
	tt := ui.NewTester(a.view, 900, 500)
	r, ok := tt.Find("Source 3: Court etiquette")
	if !ok {
		t.Fatal("badge not named")
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if a.st.hover != 3 {
		t.Fatalf("hover %d, want 3", a.st.hover)
	}
	time.Sleep(600 * time.Millisecond)
	tt.Frame()
	if n := len(tt.Texts()); n == 0 || !tt.HasText("example.net") {
		t.Fatalf("no preview card: %q", tt.Texts())
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if a.st.Active != 3 {
		t.Fatalf("active %d", a.st.Active)
	}
	tt.Move(5, 480)
	tt.Frame()
	if a.st.hover != 0 {
		t.Fatalf("hover %d after leaving", a.st.hover)
	}
}

func TestCitationBadgeKeyboard(t *testing.T) {
	a := &citationApp{}
	tt := ui.NewTester(a.view, 900, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if a.st.Active != 1 {
		t.Fatalf("Enter on the first badge: active %d", a.st.Active)
	}
}

func TestCitationBadgeBadNumber(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for number 0")
		}
	}()
	var st CitationState
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		CitationBadge(c, 0, sourcesSample[0], &st)
	}, 100, 100)
}

func TestCitationBadgeSnapshot(t *testing.T) {
	snapshot(t, "citation_badge", 360, 60, func(c *ui.Context) {
		st := ui.Local(c.Root(), "st", func() CitationState { return CitationState{Active: 2} })
		ui.Row(c).Padding(16).Gap(4).Children(func() {
			ui.Text(c, "Seventeen arches")
			CitationBadge(c, 1, sourcesSample[0], st)
			CitationBadge(c, 2, sourcesSample[1], st)
		})
	})
}
