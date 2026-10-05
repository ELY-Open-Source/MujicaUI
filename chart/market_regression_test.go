package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMarketEmptyWindow(t *testing.T) {
	from, to := marketWindow(MarketView{}, 0, 200)
	if from != 0 || to != 0 {
		t.Fatalf("empty window %d..%d", from, to)
	}
}

func TestMarketDisabledAlert(t *testing.T) {
	cs := sampleCandles(2)
	v := MarketView{}
	alerts := []PriceAlertLine{{ID: "a", Price: cs[1].Close + 1}}
	disabled, fired := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			CandlestickChart(c, cs, &v, CandlestickOptions{Frame: ChartFrame{Label: "Price"}, Alerts: &alerts,
				OnAlert: func(PriceAlertLine, Candle) { fired++ }})
		})
	}, 500, 320)
	disabled = true
	next := cs[1]
	next.Time = next.Time.Add(3600e9)
	next.Close += 2
	next.High = max(next.High, next.Close)
	cs = append(cs, next)
	tt.Frame()
	if fired != 0 {
		t.Fatalf("disabled alert fired %d times", fired)
	}
}

func TestMarketDisabledLayoutCopy(t *testing.T) {
	state := MultiChartState{Count: 2}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			MultiChartLayout(c, &state, MultiChartLayoutOptions{}, func(_ int, cell *ChartCell) {
				cell.Symbol = "changed"
			})
		})
	}, 600, 600)
	if state.Cells[0].Symbol != "" || state.Cells[1].Symbol != "" {
		t.Fatalf("disabled layout wrote cells: %+v", state.Cells)
	}
}

func TestMarketDisabledParamsCopy(t *testing.T) {
	var params []float64
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() { IndicatorParamsField(c, IndicatorMA, &params) })
	}, 500, 100)
	if params != nil {
		t.Fatalf("disabled field initialized caller parameters: %v", params)
	}
}

func TestMarketInvalidFiniteData(t *testing.T) {
	for _, v := range []float64{math.Inf(1), math.Inf(-1), math.NaN()} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("nonfinite candle accepted: %v", v)
				}
			}()
			cs := sampleCandles(1)
			cs[0].High = v
			marketCheckCandles(cs)
		}()
	}
}

func TestMarketLockedUndo(t *testing.T) {
	d := DrawingTools{Items: []Drawing{{Tool: ToolHorizontal}}}
	d.Clear()
	d.Locked = true
	if d.Undo() || len(d.Items) != 0 {
		t.Fatal("locked drawings accepted undo")
	}
}
