package chart

import (
	"image/color"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestTrendIndicatorWordsAndColors(t *testing.T) {
	for _, tc := range []struct {
		change float64
		invert bool
		label  string
		good   bool
	}{
		{12.5, false, "up +12.5%", true},
		{-3, false, "down −3.0%", false},
		{-3, true, "down −3.0%", true},
		{0, false, "unchanged 0.0%", false},
	} {
		var r ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			ui.Column(c).Padding(10).AlignItems(ui.Start).Children(func() {
				r = TrendIndicator(c, tc.change, TrendIndicatorOptions{Precision: 1, Unit: "%", Invert: tc.invert, Note: "vs Lent"}).Bounds()
			})
		}, 300, 100)
		tt.Frame()
		if _, ok := tt.Find(tc.label); !ok || !tt.HasText("vs Lent") {
			t.Fatalf("%v: no %q in %q", tc.change, tc.label, tt.Texts())
		}
		k := theme.Light()
		want := map[bool]color.RGBA{true: toRGBA(k.Success), false: toRGBA(k.Danger)}[tc.good]
		if tc.change == 0 {
			continue
		}
		found := false
		img := tt.Image()
		for x := int(r.X); x < int(r.X+16) && !found; x++ {
			for y := int(r.Y); y < int(r.Y+r.H) && !found; y++ {
				found = color.RGBAModel.Convert(img.At(x, y)).(color.RGBA) == want
			}
		}
		if !found {
			t.Fatalf("%v invert=%v: arrow not in %v", tc.change, tc.invert, want)
		}
	}
}

func TestTrendIndicatorFormatAndBad(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TrendIndicator(c, 2, TrendIndicatorOptions{Format: func(v float64) string { return "two more" }})
	}, 200, 100)
	if !tt.HasText("two more") {
		t.Fatalf("format %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for NaN")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TrendIndicator(c, math.NaN(), TrendIndicatorOptions{})
	}, 100, 100)
}

func TestTrendIndicatorSnapshot(t *testing.T) {
	snapshot(t, "trend_indicator", 420, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(24).Children(func() {
			TrendIndicator(c, 12.5, TrendIndicatorOptions{Precision: 1, Unit: "%", Note: "vs Lent"})
			TrendIndicator(c, -3, TrendIndicatorOptions{Unit: " souls"})
			TrendIndicator(c, -4, TrendIndicatorOptions{Unit: "%", Invert: true})
			TrendIndicator(c, 0, TrendIndicatorOptions{})
		})
	})
}

func TestTrendIndicatorFiniteAndChinese(t *testing.T) {
	charts2MustPanic(t, func(c *ui.Context) { TrendIndicator(c, math.Inf(1), TrendIndicatorOptions{}) })
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		TrendIndicator(c, -2, TrendIndicatorOptions{Unit: "%"})
	}, 200, 100)
	if _, ok := tt.Find("下降 −2%"); !ok {
		t.Fatalf("Chinese readout %q", tt.Texts())
	}
}
