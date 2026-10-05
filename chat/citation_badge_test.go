package chat

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var sourcesSample = []Source{
	{Title: "The Palace of Versailles", URL: "https://example.com/versailles", Site: "example.com",
		Snippet: "Seat of the court from 1682.", Published: time.Date(2025, 5, 4, 0, 0, 0, 0, time.UTC)},
	{Title: "Hall of Mirrors", URL: "https://example.org/hall", Snippet: "Seventeen arches of mirrors face the garden."},
	{Title: "Court etiquette", URL: "https://example.net/etiquette", Site: "example.net"},
}

// citationApp shares state between inline citations and sources.
type citationApp struct {
	st     CitationState
	opened []int
}

func (a *citationApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
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
		core.Use(c, core.Settings{})
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

func TestCitationBadgeHitAreaAndDisabledPreview(t *testing.T) {
	var st CitationState
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(16).Disabled(off).Children(func() {
			CitationBadge(c, 1, sourcesSample[0], &st)
		})
	}, 240, 160)
	b, _ := tt.Find("Source 1: The Palace of Versailles")
	if b.W < 28 || b.H < 28 {
		t.Errorf("hit area %v, want at least 28x28", b)
	}
	tt.Key(0, ui.KeyTab)
	time.Sleep(550 * time.Millisecond)
	tt.Frame()
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("Seat of the court from 1682.") || st.Active != 0 {
		t.Errorf("disabled preview or activation: active %d texts %q", st.Active, tt.Texts())
	}
}

func TestCitationBadgePreviewFitsNarrowWindow(t *testing.T) {
	var st CitationState
	var preview ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		src := sourcesSample[0]
		ui.Row(c).Padding(8).Children(func() { CitationBadge(c, 1, src, &st) })
	}, 180, 120)
	tt.Key(0, ui.KeyTab)
	time.Sleep(550 * time.Millisecond)
	tt.Frame()
	preview, ok := tt.Find("Source preview")
	if !ok || preview.X < 0 || preview.X+preview.W > 180 || preview.Y+preview.H > 120 {
		t.Fatalf("preview content outside window: %v found %v", preview, ok)
	}
}

func TestCitationBadgeRemovedOwnerClearsHover(t *testing.T) {
	show := true
	st := CitationState{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			ui.Row(c).Key("badge").Height(36).Children(func() {
				if show {
					CitationBadge(c, 1, sourcesSample[0], &st)
				}
			})
			ui.Row(c).Key("card").Width(300).Children(func() { SourceCard(c, 1, sourcesSample[0], &st) })
		})
	}, 400, 300)
	b, _ := tt.Find("Source 1: The Palace of Versailles")
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	tt.Frame()
	if st.hover != 1 {
		t.Fatal("hover did not activate")
	}
	show = false
	tt.Move(390, 290)
	tt.Frame()
	if st.hover != 0 {
		t.Fatalf("removed owner retained hover %d", st.hover)
	}
}
