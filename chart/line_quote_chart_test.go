package chart

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func samplePoints() ([]QuotePoint, time.Time, time.Time) {
	day := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	var pts []QuotePoint
	for m := 8 * 60; m < 18*60; m += 5 {
		price := 100 + 2*math.Sin(float64(m)/50) + float64(m-480)/200
		pts = append(pts, QuotePoint{Time: day.Add(time.Duration(m) * time.Minute), Price: math.Round(price*100) / 100, Volume: 100 + float64(m%7)*10})
	}
	return pts, day.Add(9*time.Hour + 30*time.Minute), day.Add(16 * time.Hour)
}

func TestLineQuoteChartReadoutAverageSessions(t *testing.T) {
	pts, open, closeAt := samplePoints()
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = LineQuoteChart(c, pts, LineQuoteChartOptions{Frame: ChartFrame{Label: "AAPL"}, PrevClose: 100, RegularOpen: open, RegularClose: closeAt})
	}, 800, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	var pv, v float64
	pv, v = pts[0].Price*pts[0].Volume, pts[0].Volume
	want := "08:00 · " + core.PriceUngrouped(pts[0].Price, 2)
	found := ""
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, want) {
			found = s
		}
	}
	if found == "" || !strings.Contains(found, "Average "+core.PriceUngrouped(pv/v, 2)) || !strings.HasSuffix(found, "Pre-market") {
		t.Fatalf("first readout %q in %q", found, tt.Texts())
	}
	tt.Key(0, ui.KeyEnd)
	last := ""
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "17:55") {
			last = s
		}
	}
	if !strings.HasSuffix(last, "After hours") || res.Point != len(pts)-1 {
		t.Fatalf("last readout %q", last)
	}
	for name, o := range map[string]LineQuoteChartOptions{
		"prev close": {Frame: ChartFrame{Label: "x"}},
		"half bound": {Frame: ChartFrame{Label: "x"}, PrevClose: 1, RegularOpen: open},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); LineQuoteChart(c, pts, o) }, 300, 200)
		}()
	}
}

func TestLineQuoteChartSnapshot(t *testing.T) {
	pts, open, closeAt := samplePoints()
	snapshot(t, "line_quote_chart", 680, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			LineQuoteChart(c, pts, LineQuoteChartOptions{Frame: ChartFrame{Label: "AAPL", Height: 260}, PrevClose: 100, RegularOpen: open, RegularClose: closeAt})
		})
	})
}
