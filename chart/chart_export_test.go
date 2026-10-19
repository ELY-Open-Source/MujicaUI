package chart

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestChartExportPNGInTheme(t *testing.T) {
	for _, mode := range []core.Mode{core.Light, core.Dark} {
		var out []byte
		ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: mode})
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
		if mode == core.Dark {
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
		core.Use(c, core.Settings{})
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

func TestChartExportPreservesPreferences(t *testing.T) {
	var data []byte
	var got ui.Preferences
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Dark, Locale: core.ZhCN})
		if data != nil {
			return
		}
		var err error
		data, err = ChartExport(c, func(c *ui.Context) { got = c.Preferences(); ui.Text(c, "Scale") }, ChartExportOptions{Width: 100, Height: 80, Scale: 1})
		if err != nil {
			t.Fatal(err)
		}
	}, 100, 80)
	data = nil
	tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true, HighContrast: true})
	if got.TextScale != 2 || !got.ReduceMotion || !got.HighContrast {
		t.Fatalf("export preferences %+v", got)
	}
}
