package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGridFixedAndFrTracks(t *testing.T) {
	var a, b ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Grid(c, GridOptions{Columns: []ui.Track{ui.Fixed(100), ui.Fr(1)}}, func() {
			a = ui.Box(c).Height(20).Bounds()
			b = ui.Box(c).Height(20).Bounds()
		}).Width(400)
	}, 400, 300).Frame()
	if a.W != 100 || b.X != 112 || b.W != 288 {
		t.Fatalf("tracks %v and %v at %v, want 100, 288 at 112", a.W, b.W, b.X)
	}
}

func TestGridSpans(t *testing.T) {
	var head, tall, side ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Grid(c, GridOptions{Columns: []ui.Track{ui.Fr(1), ui.Fr(1), ui.Fr(1)}, Gap: SpaceS, RowGap: SpaceL}, func() {
			head = GridCell(c, GridCellOptions{ColumnSpan: -1}, func() { ui.Box(c).Height(20) }).Bounds()
			tall = GridCell(c, GridCellOptions{RowSpan: 2}, func() { ui.Box(c).Height(10) }).Bounds()
			side = GridCell(c, GridCellOptions{Column: 3, Row: 3}, func() { ui.Box(c).Height(30) }).Bounds()
			GridCell(c, GridCellOptions{Column: 2, Row: 2}, func() { ui.Box(c).Height(30) })
		}).Width(316)
	}, 400, 300).Frame()
	if head.W != 316 {
		t.Fatalf("header width %v, want 316", head.W)
	}
	// Rows of 30 with a 16 gap: the tall cell spans 30+16+30.
	if tall.H != 76 {
		t.Fatalf("row-spanning cell height %v, want 76", tall.H)
	}
	if side.X != 216 || side.Y != tall.Y+46 {
		t.Fatalf("placed cell at %v,%v, want 216,%v", side.X, side.Y, tall.Y+46)
	}
}

func TestGridRowTracksShareHeight(t *testing.T) {
	var a, b ui.Rect
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Grid(c, GridOptions{Rows: []ui.Track{ui.Fixed(40), ui.Fr(1)}, Gap: SpaceNone}, func() {
			a = ui.Box(c).Bounds()
			b = ui.Box(c).Bounds()
		}).Height(200)
	}, 300, 300).Frame()
	if a.H != 40 || b.H != 160 {
		t.Fatalf("rows %v and %v, want 40 and 160", a.H, b.H)
	}
}

func TestGridCellRejectsBadPlacement(t *testing.T) {
	for _, o := range []GridCellOptions{{Column: -1}, {RowSpan: -2}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", o)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				Grid(c, GridOptions{}, func() { GridCell(c, o, func() {}) })
			}, 100, 100)
		}()
	}
}

func TestGridSnapshot(t *testing.T) {
	snapshot(t, "grid", 420, 200, func(c *ui.Context) {
		k := Tokens(c)
		cell := func(s string) {
			ui.Box(c).Center().Padding(8).Background(k.Surface).Border(1, k.Border).Radius(4).
				Children(func() { ui.Text(c, s).TextColor(k.TextMuted) })
		}
		ui.Column(c).Padding(16).Children(func() {
			Grid(c, GridOptions{Columns: []ui.Track{ui.Fixed(96), ui.Fr(1), ui.Fr(1)}}, func() {
				GridCell(c, GridCellOptions{RowSpan: 2}, func() { cell("Fixed 96") })
				GridCell(c, GridCellOptions{ColumnSpan: 2}, func() { cell("Span 2") })
				cell("Fr 1")
				cell("Fr 1")
				GridCell(c, GridCellOptions{ColumnSpan: -1}, func() { cell("Across") })
			})
		})
	})
}
