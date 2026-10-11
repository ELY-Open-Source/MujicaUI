package chart

import (
	"bytes"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMarketIndicatorDomainChangesWithParams(t *testing.T) {
	cs := sampleCandles(300)
	v := MarketView{Count: 80}
	panes := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorRSI}}}
	build := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IndicatorPane(c, cs, &v, &panes, IndicatorPaneOptions{Label: "BTC"})
	}
	tt := ui.NewTester(build, 640, 400)
	panes.Panes[0].Kind = IndicatorOBV
	tt.Frame()
	fresh := ui.NewTester(build, 640, 400)
	if !bytes.Equal(encode(t, tt.Image()), encode(t, fresh.Image())) {
		t.Fatal("changed indicator retained its previous scale")
	}

}

func TestMarketRejectsInvalidLayers(t *testing.T) {
	for _, kind := range []string{"compare zero", "compare infinity", "alert infinity", "alert duplicate", "profile NaN", "depth NaN", "indicator infinity"} {
		t.Run(kind, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("invalid data accepted")
				}
			}()
			cs := sampleCandles(2)
			v := MarketView{}
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				o := CandlestickOptions{Frame: ChartFrame{Label: "Price"}}
				switch kind {
				case "compare zero":
					x := []CompareSymbol{{Name: "X", Closes: []float64{0, 1}}}
					o.Compare = &x
				case "compare infinity":
					x := []CompareSymbol{{Name: "X", Closes: []float64{1, math.Inf(1)}}}
					o.Compare = &x
				case "alert infinity":
					x := []PriceAlertLine{{ID: "a", Price: math.Inf(1)}}
					o.Alerts = &x
				case "alert duplicate":
					x := []PriceAlertLine{{ID: "a", Price: 1}, {ID: "a", Price: 2}}
					o.Alerts = &x
				case "profile NaN":
					VolumeProfile(c, []ProfileBin{{1, 2, math.NaN()}}, VolumeProfileOptions{Frame: o.Frame})
					return
				case "depth NaN":
					DepthChart(c, []DepthLevel{{math.NaN(), 1}}, []DepthLevel{{2, 1}}, DepthChartOptions{Frame: o.Frame})
					return
				case "indicator infinity":
					o.Overlays = []IndicatorOverlay{{Kind: IndicatorBOLL, Params: []float64{2, math.Inf(1)}}}
				}
				CandlestickChart(c, cs, &v, o)
			}, 500, 320)
		})
	}
}
