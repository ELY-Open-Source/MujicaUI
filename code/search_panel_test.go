package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func searchPanelSample() []SearchFileResult {
	return []SearchFileResult{
		{Path: "court/seal.go", Matches: []SearchMatch{{Line: 12, Text: "func Seal(d Decree) {", Start: 5, End: 9}, {Line: 30, Text: "\treturn Seal(next)", Start: 8, End: 12}}},
		{Path: "docs/ledger.md", Matches: []SearchMatch{{Line: 4, Text: "The seal is kept by the chancellor.", Start: 4, End: 8}}},
	}
}

type searchPanelProbe struct {
	results []SearchFileResult
	s       SearchPanelState
	opts    SearchPanelOptions
	changed int
	opened  []SearchOpen
}

func (p *searchPanelProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := SearchPanel(c, p.results, &p.s, p.opts)
	r.Element.Height(400)
	if r.Changed {
		p.changed++
	}
	if r.Opened != nil {
		p.opened = append(p.opened, *r.Opened)
	}
}

func TestSearchPanelGroupsAndOpens(t *testing.T) {
	p := &searchPanelProbe{results: searchPanelSample()}
	tt := ui.NewTester(p.view, 560, 420)
	if !tt.HasText("3 results in 2 files") || !tt.HasText("court/seal.go") {
		t.Fatalf("texts %q", tt.Texts())
	}
	r, _ := tt.Find("12")
	tt.ClickAt(r.X+60, r.Y+r.H/2)
	tt.ClickAt(r.X+60, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0] != (SearchOpen{"court/seal.go", 12}) {
		t.Fatalf("opened %v", p.opened)
	}
	if err := tt.Click("Fold court/seal.go"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("12") || !tt.HasText("docs/ledger.md") {
		t.Fatalf("fold kept hits: %q", tt.Texts())
	}
	_ = tt.Click("Unfold court/seal.go")
	if !tt.HasText("30") {
		t.Fatal("unfold did not show the hits")
	}
}

func TestSearchPanelQueryAndPatternsReport(t *testing.T) {
	p := &searchPanelProbe{}
	tt := ui.NewTester(p.view, 560, 420)
	if !tt.HasText("Type to search in files") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("seal")
	tt.Key(0, ui.KeyTab)
	tt.Type("*.go")
	if p.changed != 2 || p.s.Query != "seal" || p.s.Include != "*.go" {
		t.Fatalf("changed %d query %q include %q", p.changed, p.s.Query, p.s.Include)
	}
	if !tt.HasText("No results") {
		t.Fatal("empty search with a query says nothing")
	}
}

func TestSearchPanelKeyboard(t *testing.T) {
	p := &searchPanelProbe{results: searchPanelSample()}
	tt := ui.NewTester(p.view, 560, 420)
	for range 5 {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyLeft)
	if !p.s.Collapsed.Has("court/seal.go") {
		t.Fatal("Left on a file did not fold it")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 1 || p.opened[0].Line != 30 {
		t.Fatalf("Enter opened %v", p.opened)
	}
}

func TestSearchPanelStatusAndDisabled(t *testing.T) {
	p := &searchPanelProbe{opts: SearchPanelOptions{Status: data.DataFailed, Error: "index unavailable"}}
	tt := ui.NewTester(p.view, 560, 420)
	if !tt.HasText("index unavailable") {
		t.Fatalf("texts %q", tt.Texts())
	}
	p = &searchPanelProbe{results: searchPanelSample(), opts: SearchPanelOptions{Disabled: true}}
	tt = ui.NewTester(p.view, 560, 420)
	_ = tt.Click("Fold court/seal.go")
	if p.s.Collapsed.Len() != 0 {
		t.Fatal("disabled panel folded")
	}
}

func TestSearchPanelBadRangePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("match outside its text did not panic")
		}
	}()
	p := &searchPanelProbe{results: []SearchFileResult{{Path: "a", Matches: []SearchMatch{{Text: "ab", Start: 1, End: 9}}}}}
	ui.NewTester(p.view, 560, 420)
}

func TestSearchPanelSnapshot(t *testing.T) {
	snapshot(t, "search_panel", 520, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *SearchPanelState { return &SearchPanelState{Query: "seal", Include: "*.go, *.md"} })
		ui.Column(c).Padding(12).Children(func() {
			SearchPanel(c, searchPanelSample(), *s, SearchPanelOptions{}).Element.Height(276)
		})
	})
}
