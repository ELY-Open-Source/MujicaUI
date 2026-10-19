package chart

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestHeikinAshiChartSwitchesToRaw(t *testing.T) {
	candles := sampleCandles(30)
	view := MarketView{Count: 30}
	raw := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		HeikinAshiChart(c, candles, &view, &raw, CandlestickOptions{Frame: ChartFrame{Label: "HA"}})
	}, 600, 360)
	ha := HeikinAshi(candles)
	readout := func() string {
		for _, s := range tt.Texts() {
			if strings.HasPrefix(s, "2026-01-01 00:00") {
				return s
			}
		}
		return ""
	}
	tt.Key(0, ui.KeyTab) // the switch
	tt.Key(0, ui.KeyTab) // the plot
	tt.Key(0, ui.KeyHome)
	if want := "O " + core.PriceUngrouped(ha[0].Open, 2); !strings.Contains(readout(), want) {
		t.Fatalf("Heikin-Ashi readout %q, want %q", readout(), want)
	}
	if err := tt.Click("Candles"); err != nil || !raw {
		t.Fatalf("switch: %v %v", err, raw)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if want := "O " + core.PriceUngrouped(candles[0].Open, 2); !strings.Contains(readout(), want) {
		t.Fatalf("raw readout %q, want %q", readout(), want)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("nil raw accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		HeikinAshiChart(c, candles, &view, nil, CandlestickOptions{})
	}, 200, 100)
}

func TestHeikinAshiChartSnapshot(t *testing.T) {
	candles := sampleCandles(70)
	view := MarketView{Count: 60}
	raw := false
	snapshot(t, "heikin_ashi_chart", 640, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			HeikinAshiChart(c, candles, &view, &raw, CandlestickOptions{Frame: ChartFrame{Label: "HA", Height: 230}})
		})
	})
}
