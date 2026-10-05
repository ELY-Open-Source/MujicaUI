package mujica

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestChartExportPNGInTheme(t *testing.T) {
	for _, mode := range []Mode{Light, Dark} {
		var out []byte
		ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Mode: mode})
			if out != nil {
				return
			}
			b, err := ChartExport(c, func(c *ui.Context) {
				BarChart(c, barCats, barSeries(), BarChartOptions{Frame: ChartFrame{Label: "Tithes", Height: 160}})
			}, ChartExportOptions{Width: 300, Height: 220, Scale: 1.5})
			if err != nil {
				t.Fatal(err)
			}
			out = b
		}, 100, 100)
		img, err := png.Decode(bytes.NewReader(out))
		if err != nil {
			t.Fatal(err)
		}
		if b := img.Bounds(); b.Dx() != 450 || b.Dy() != 330 {
			t.Fatalf("%v: size %v", mode, b)
		}
		k := theme.Light()
		if mode == Dark {
			k = theme.Dark()
		}
		if got := color.RGBAModel.Convert(img.At(3, 3)).(color.RGBA); got != toRGBA(k.Surface) {
			t.Fatalf("%v: corner %v, want the theme's Surface", mode, got)
		}
	}
}

func TestChartExportBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a zero size")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		_, _ = ChartExport(c, func(*ui.Context) {}, ChartExportOptions{})
	}, 100, 100)
}

func TestChartExportSnapshot(t *testing.T) {
	snapshot(t, "chart_export", 340, 300, func(c *ui.Context) {
		b, err := ChartExport(c, func(c *ui.Context) {
			PieChart(c, pieSlices, PieChartOptions{Frame: ChartFrame{Label: "Stores", Height: 150}, Donut: true, MinAngle: 20})
		}, ChartExportOptions{Width: 300, Height: 260, Scale: 1})
		if err != nil {
			t.Fatal(err)
		}
		bmp, err := ui.DecodeBitmap(b)
		if err != nil {
			t.Fatal(err)
		}
		ui.Column(c).Padding(20).Children(func() { ui.Image(c, bmp).Size(300, 260) })
	})
}
