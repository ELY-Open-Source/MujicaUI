package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func sampleBook() ([]DepthLevel, []DepthLevel) {
	var bids, asks []DepthLevel
	for i := range 12 {
		bids = append(bids, DepthLevel{Price: 99.9 - float64(i)*0.1, Size: float64(10 + i*3)})
		asks = append(asks, DepthLevel{Price: 100.1 + float64(i)*0.1, Size: float64(8 + i*4)})
	}
	return bids, asks
}

func TestDepthChartCumulativeAndKeyboard(t *testing.T) {
	bids, asks := sampleBook()
	d := depthBuild(bids, asks)
	if d.mid != 100 || d.total[0] != 10 || d.total[1] != 23 || d.total[12] != 8 || d.total[13] != 20 {
		t.Fatalf("cumulative %+v", d.total[:14])
	}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = DepthChart(c, bids, asks, DepthChartOptions{Frame: ChartFrame{Label: "Depth"}})
	}, 700, 360)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight) // the best bid first
	if res.Point != 0 || !tt.HasText("Bid 99.90 · cumulative 10 · Mid 100.00") {
		t.Fatalf("first level %d %q", res.Point, tt.Texts())
	}
	tt.Key(0, ui.KeyLeft) // outward to lower bids
	if res.Point != 1 || !tt.HasText("Bid 99.80 · cumulative 23") {
		t.Fatalf("lower bid %d %q", res.Point, tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if res.Point != 12 || !tt.HasText("Ask 100.10 · cumulative 8") {
		t.Fatalf("crossing to asks %d %q", res.Point, tt.Texts())
	}
	for name, book := range map[string][2][]DepthLevel{
		"bid order": {{{Price: 1, Size: 1}, {Price: 2, Size: 1}}, nil},
		"crossed":   {{{Price: 5, Size: 1}}, {{Price: 4, Size: 1}}},
		"size":      {nil, {{Price: 4, Size: 0}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			depthBuild(book[0], book[1])
		}()
	}
}

func TestDepthChartSnapshot(t *testing.T) {
	bids, asks := sampleBook()
	snapshot(t, "depth_chart", 640, 320, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			DepthChart(c, bids, asks, DepthChartOptions{Frame: ChartFrame{Label: "Depth", Height: 240}})
		})
	})
}
