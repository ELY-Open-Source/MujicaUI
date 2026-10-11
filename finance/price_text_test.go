package finance

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// tradingTester runs view in the default theme with trend t.
func tradingTester(t core.Trend, w, h int, view func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		core.UseTrend(c, t)
		ui.Column(c).Padding(10).AlignItems(ui.Start).Children(func() { view(c) })
	}, w, h)
}

func TestPriceTextFormatAndTrend(t *testing.T) {
	price := 12345.678
	tt := tradingTester(core.GreenUp, 300, 100, func(c *ui.Context) {
		PriceText(c, price, PriceTextOptions{Decimals: 2, Compare: true, Reference: 12000})
	})
	if !tt.HasText("12,345.68") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if _, ok := tt.Find("12,345.68, up"); !ok {
		t.Fatal("trend not spoken")
	}
	price = 11000
	tt.Frame()
	if _, ok := tt.Find("11,000.00, down"); !ok {
		t.Fatalf("down not spoken: %q", tt.Texts())
	}
	tt = tradingTester(core.GreenUp, 300, 100, func(c *ui.Context) {
		PriceText(c, 1234.5, PriceTextOptions{Decimals: 3, NoGrouping: true})
	})
	if !tt.HasText("1234.500") {
		t.Fatalf("ungrouped %q", tt.Texts())
	}
}

func TestTrendColorsFollowUseTrend(t *testing.T) {
	var up, down [2]ui.Color
	for i, tr := range []core.Trend{core.GreenUp, core.RedUp} {
		tradingTester(tr, 100, 100, func(c *ui.Context) {
			up[i], down[i] = core.TrendColor(c, 1), core.TrendColor(c, -1)
		})
	}
	k := theme.Light()
	if up[0] != k.Success || down[0] != k.Danger || up[1] != k.Danger || down[1] != k.Success {
		t.Fatalf("up %v down %v", up, down)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown trend accepted")
		}
	}()
	tradingTester(core.Trend(7), 100, 100, func(*ui.Context) {})
}

// A change tints once; reduced motion never tints.
func TestPriceTextFlashesOnceAndNotUnderReducedMotion(t *testing.T) {
	for _, still := range []bool{false, true} {
		prefs := ui.Preferences{ReduceMotion: still, TextScale: 1}
		price := 10.0
		tt := tradingTester(core.GreenUp, 200, 60, func(c *ui.Context) { PriceText(c, price, PriceTextOptions{}) })
		tt.SetPreferences(prefs)
		price = 10.5
		tt.Frame()
		fresh := tradingTester(core.GreenUp, 200, 60, func(c *ui.Context) { PriceText(c, price, PriceTextOptions{}) })
		fresh.SetPreferences(prefs)
		fresh.Frame()
		same := bytes.Equal(tt.Image().Pix, fresh.Image().Pix)
		if same == !still {
			t.Fatalf("reduced motion %v: flashed %v", still, !same)
		}
		time.Sleep(priceTextFlash + 20*time.Millisecond)
		tt.Frame()
		if !bytes.Equal(tt.Image().Pix, fresh.Image().Pix) {
			t.Fatalf("reduced motion %v: flash did not end", still)
		}
	}
}

func TestPriceTextRejectsNegativeDecimals(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("negative decimals accepted")
		}
	}()
	tradingTester(core.GreenUp, 100, 100, func(c *ui.Context) { PriceText(c, 1, PriceTextOptions{Decimals: -1}) })
}

func TestPriceTextSnapshot(t *testing.T) {
	snapshot(t, "price_text", 320, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			PriceText(c, 18234.5, PriceTextOptions{Size: theme.H2Size, Compare: true, Reference: 18000})
			PriceText(c, 0.004521, PriceTextOptions{Decimals: 6, Compare: true, Reference: 0.0046})
			PriceText(c, 100, PriceTextOptions{Compare: true, Reference: 100})
		})
	})
}
