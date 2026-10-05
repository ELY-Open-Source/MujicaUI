package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

type hexProbe struct {
	data     []byte
	s        HexViewerState
	disabled bool
	changed  int
}

func (p *hexProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := HexViewer(c, p.data, &p.s, HexViewerOptions{Disabled: p.disabled})
	r.Element.Height(300)
	if r.Changed() {
		p.changed++
	}
}

func TestHexViewerColumns(t *testing.T) {
	p := &hexProbe{data: []byte("Court seal\x00\x01 ledger of the realm")}
	tt := ui.NewTester(p.view, 760, 320)
	for _, s := range []string{"00000000", "00000010", "43 6f 75 72 74 20 73 65 61 6c 00 01 20 6c 65 64", "Court seal.. led"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
}

func TestHexViewerKeysSelectAndCopy(t *testing.T) {
	p := &hexProbe{data: []byte("0123456789abcdefGHIJ")}
	tt := ui.NewTester(p.view, 760, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(ui.Shift, ui.KeyRight)
	tt.Key(ui.Shift, ui.KeyRight)
	if a, b := p.s.Selection(); a != 1 || b != 4 {
		t.Fatalf("selection %d..%d", a, b)
	}
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "31 32 33" {
		t.Fatalf("copied %q", tt.Clipboard())
	}
	tt.Key(0, ui.KeyDown)
	if p.s.Cursor != 19 || p.s.Anchor != 19 {
		t.Fatalf("Down clamped to %d", p.s.Cursor)
	}
	tt.Key(0, ui.KeyHome)
	if p.s.Cursor != 16 {
		t.Fatalf("Home gave %d", p.s.Cursor)
	}
}

func TestHexViewerClickSyncsColumns(t *testing.T) {
	p := &hexProbe{data: []byte("0123456789abcdef")}
	tt := ui.NewTester(p.view, 760, 320)
	r, _ := tt.Find("0123456789abcdef")
	cw := r.W / 16
	tt.ClickAt(r.X+cw*5.5, r.Y+r.H/2)
	if p.s.Cursor != 5 || p.changed != 1 {
		t.Fatalf("ASCII click chose %d", p.s.Cursor)
	}
	h, _ := tt.Find("30 31 32 33 34 35 36 37 38 39 61 62 63 64 65 66")
	tt.ClickAt(h.X+cw*3*9+cw, h.Y+h.H/2)
	if p.s.Cursor != 9 {
		t.Fatalf("hex click chose %d", p.s.Cursor)
	}
}

func TestHexViewerLargeIsVirtual(t *testing.T) {
	p := &hexProbe{data: make([]byte, 1_600_000)}
	tt := ui.NewTester(p.view, 760, 320)
	rows := 0
	for _, s := range tt.Texts() {
		if len(s) == 8 && strings.Trim(s, "0123456789abcdef") == "" {
			rows++
		}
	}
	if rows == 0 || rows > 40 {
		t.Fatalf("built %d of 100000 rows", rows)
	}
}

func TestHexViewerDisabledAndBadCursor(t *testing.T) {
	p := &hexProbe{data: []byte("abc"), disabled: true}
	tt := ui.NewTester(p.view, 760, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if p.s.Cursor != 0 {
		t.Fatal("disabled viewer moved")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("cursor outside data did not panic")
		}
	}()
	ui.NewTester((&hexProbe{data: []byte("a"), s: HexViewerState{Cursor: 4}}).view, 760, 320)
}

func TestHexViewerSnapshot(t *testing.T) {
	snapshot(t, "hex_viewer", 720, 170, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *hexProbe {
			return &hexProbe{data: []byte("MUJICA\x00\x01court ledger\nseal: 42\xff\xfe"), s: HexViewerState{Cursor: 9, Anchor: 6}}
		})
		ui.Column(c).Padding(12).Children(func() {
			HexViewer(c, (*p).data, &(*p).s, HexViewerOptions{}).Element.Height(146)
		})
	})
}
