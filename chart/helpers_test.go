package chart

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal/snaptest"
)

// tabTo presses Tab until the element labeled label has the focus.
func tabTo(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 60 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("Tab never reached %q", label)
}

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

func toRGBA(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, c.A} }

func encode(t *testing.T, img image.Image) []byte {
	t.Helper()
	return snaptest.Encode(t, img)
}

func mktHeat() HeatNode {
	return HeatNode{Name: "Market", Children: []HeatNode{
		{Name: "Technology", Children: []HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}, {Name: "NVDA", Cap: 2200, Change: 3.4}}},
		{Name: "Energy", Children: []HeatNode{{Name: "XOM", Cap: 450, Change: -2.1}, {Name: "CVX", Cap: 300, Change: -1.4}}},
		{Name: "Health", Children: []HeatNode{{Name: "LLY", Cap: 700, Change: 0.4}, {Name: "JNJ", Cap: 380, Change: 0}}},
	}}
}

func chatTester(w, h int, view func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Fill().Children(func() { view(c) })
	}, w, h)
}

func marketsPanic(t *testing.T, view func(*ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	ui.Render(func(c *ui.Context) { core.Use(c, core.Settings{}); view(c) }, 960, 500, 1)
}

func snapshotTheme(c *ui.Context, mode string) {
	m := core.Light
	if mode == "dark" {
		m = core.Dark
	}
	core.Use(c, core.Settings{Mode: m})
}
