package input

import (
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMasonryPlaceShortestColumn(t *testing.T) {
	hs := map[string]float32{"0": 100, "1": 40, "2": 40, "3": 50}
	l := masonryPlace(4, 2, 100, 10, 160, hs, strconv.Itoa)
	// 0 left; 1 right; 2 right again (50 < 110); 3 right (100 = 110? no: left 110, right 100).
	want := [][2]float32{{0, 0}, {110, 0}, {110, 50}, {110, 100}}
	for i, w := range want {
		if l.x[i] != w[0] || l.y[i] != w[1] {
			t.Errorf("item %d at %v,%v, want %v", i, l.x[i], l.y[i], w)
		}
	}
	if l.total != 150 {
		t.Errorf("total %v", l.total)
	}
	defer func() {
		if recover() == nil {
			t.Error("repeated key accepted")
		}
	}()
	masonryPlace(2, 1, 100, 0, 10, nil, func(int) string { return "same" })
}

func TestMasonryVirtualizesAndMeasures(t *testing.T) {
	built := map[int]bool{}
	n := 100000
	width := 700
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Masonry(c, n, MasonryOptions{Key: strconv.Itoa, ColumnWidth: 200, Gap: 10}, func(i int) {
			built[i] = true
			ui.Box(c).Height(float32(40 + i%3*30)).Label("tile " + strconv.Itoa(i))
		}).Size(float32(width), 500)
	}, 800, 600)
	tt.Frame()
	clear(built)
	tt.Frame()
	if len(built) == 0 || len(built) > 120 {
		t.Fatalf("built %d of %d", len(built), n)
	}
	if !built[0] || !built[2] || built[5000] {
		t.Fatalf("wrong items built: %v", built[5000])
	}
	a, _ := tt.Find("tile 0")
	b, _ := tt.Find("tile 1")
	d, _ := tt.Find("tile 3")
	if a.W < 226 || a.W > 227 || b.X < 236 || b.X > 237 {
		t.Fatalf("three columns in 700: tile 0 %v, tile 1 %v", a, b)
	}
	// Measured: tile 0 is 40 high, so tile 3 goes under it at 40 + 10.
	if d.X != 0 || d.Y != 50 {
		t.Fatalf("tile 3 at %v; texts %q", d, tt.Texts())
	}
	tt.Scroll(100, 100, 0, 600000)
	clear(built)
	tt.Frame()
	if len(built) == 0 || len(built) > 120 || built[0] {
		t.Fatalf("after scrolling built %d, includes first %v", len(built), built[0])
	}
}

func TestMasonrySnapshot(t *testing.T) {
	snapshot(t, "masonry", 520, 300, func(c *ui.Context) {
		k := core.Tokens(c)
		Masonry(c, 9, MasonryOptions{Key: strconv.Itoa, ColumnWidth: 150}, func(i int) {
			ui.Column(c).Height(float32(50+i%4*25)).Padding(8).Background(k.Surface).Border(1, k.Border).Radius(4).Children(func() {
				ui.Text(c, "Card "+strconv.Itoa(i+1))
			})
		}).Size(500, 280).Margin(10)
	})
}
