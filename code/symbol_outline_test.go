package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func symbolSample() []Symbol {
	return []Symbol{
		{Name: "Court", Kind: SymbolType, Line: 3, EndLine: 40, Detail: "struct", Children: []Symbol{
			{Name: "Seal", Kind: SymbolMethod, Line: 10, EndLine: 20, Detail: "(d Decree) error"},
			{Name: "ledger", Kind: SymbolField, Line: 5, EndLine: 5},
		}},
		{Name: "Herald", Kind: SymbolFunction, Line: 44, EndLine: 60},
		{Name: "MaxSeals", Kind: SymbolConstant, Line: 1, EndLine: 1},
	}
}

type symbolProbe struct {
	symbols []Symbol
	s       SymbolOutlineState
	opts    SymbolOutlineOptions
	jumped  []string
}

func (p *symbolProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := SymbolOutline(c, p.symbols, &p.s, p.opts)
	r.Element.Height(320)
	if r.Jumped != nil {
		p.jumped = append(p.jumped, r.Jumped.Name)
	}
}

func TestSymbolOutlineFollowsCursor(t *testing.T) {
	p := &symbolProbe{symbols: symbolSample(), opts: SymbolOutlineOptions{CursorLine: 12}}
	tt := ui.NewTester(p.view, 400, 340)
	if !tt.HasText("Seal") || p.s.cursor < 0 || p.s.nodes[p.s.outline.Item(p.s.cursor)].sym.Name != "Seal" {
		t.Fatalf("cursor at line 12 did not choose Seal: %q", tt.Texts())
	}
	p.opts.CursorLine = 50
	tt.Frame()
	if name := p.s.nodes[p.s.outline.Item(p.s.cursor)].sym.Name; name != "Herald" {
		t.Fatalf("line 50 chose %q", name)
	}
}

func TestSymbolOutlineFilter(t *testing.T) {
	p := &symbolProbe{symbols: symbolSample()}
	tt := ui.NewTester(p.view, 400, 340)
	tt.Key(0, ui.KeyTab)
	tt.Type("led")
	if !tt.HasText("Court") || !tt.HasText("ledger") || tt.HasText("Herald") || tt.HasText("Seal") {
		t.Fatalf("filter shows %q", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("zz")
	if !tt.HasText("No symbols") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestSymbolOutlineJumpByKeyboard(t *testing.T) {
	p := &symbolProbe{symbols: symbolSample()}
	tt := ui.NewTester(p.view, 400, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.jumped) != 1 || p.jumped[0] != "Herald" {
		t.Fatalf("jumped %v", p.jumped)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Seal") {
		t.Fatal("Right did not open Court")
	}
}

func TestSymbolOutlineDisabledAndBadKind(t *testing.T) {
	p := &symbolProbe{symbols: symbolSample(), opts: SymbolOutlineOptions{Disabled: true}}
	tt := ui.NewTester(p.view, 400, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(p.jumped) != 0 {
		t.Fatal("disabled outline jumped")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("bad kind did not panic")
		}
	}()
	ui.NewTester((&symbolProbe{symbols: []Symbol{{Name: "x", Kind: 99}}}).view, 400, 340)
}

func TestSymbolOutlineSnapshot(t *testing.T) {
	snapshot(t, "symbol_outline", 360, 260, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *symbolProbe {
			return &symbolProbe{symbols: symbolSample(), opts: SymbolOutlineOptions{CursorLine: 12}}
		})
		ui.Column(c).Padding(12).Children(func() {
			SymbolOutline(c, (*p).symbols, &(*p).s, (*p).opts).Element.Height(236)
		})
	})
}
