package devtools

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func uptimeSample(n int) []UptimeDay {
	days := make([]UptimeDay, n)
	start := core.Date{Year: 2026, Month: 9, Day: 1}
	for i := range days {
		days[i] = UptimeDay{Date: start.AddDays(i), Uptime: 1}
	}
	days[3] = UptimeDay{Date: start.AddDays(3), Uptime: 0.9975, Note: "slow search"}
	days[7] = UptimeDay{Date: start.AddDays(7), Uptime: 0.91, Note: "2 h outage"}
	days[9].Uptime = math.NaN()
	return days
}

func uptimeTester(days []UptimeDay, frame chart.ChartFrame, points *[]int) *ui.Tester {
	if frame.Label == "" {
		frame.Label = "API uptime"
	}
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			r := UptimeBar(c, days, UptimeBarOptions{Frame: frame})
			*points = append(*points, r.Point)
		})
	}, 520, 140)
}

func TestUptimeBarKeyboardReadout(t *testing.T) {
	var pts []int
	tt := uptimeTester(uptimeSample(30), chart.ChartFrame{}, &pts)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	for range 7 {
		tt.Key(0, ui.KeyRight)
	}
	if !tt.HasText("2026-09-08 · Outage · 91.00% · 2 h outage") {
		t.Fatalf("readout: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("2026-09-10 · No data") {
		t.Fatalf("no data readout: %q", tt.Texts())
	}
}

func TestUptimeBarHover(t *testing.T) {
	var pts []int
	tt := uptimeTester(uptimeSample(30), chart.ChartFrame{}, &pts)
	// The 4th of 30 bars across 500 DIPs starting at x = 10.
	tt.Move(10+500*3.5/30, 30)
	if pts[len(pts)-1] != 3 || !tt.HasText("2026-09-04 · Degraded · 99.75% · slow search") {
		t.Fatalf("hover point %d: %q", pts[len(pts)-1], tt.Texts())
	}
}

func TestUptimeBarStatesAndRules(t *testing.T) {
	var pts []int
	tt := uptimeTester(nil, chart.ChartFrame{}, &pts)
	if !tt.HasText("No data") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	tt = uptimeTester(uptimeSample(12), chart.ChartFrame{Status: data.DataLoading}, &pts)
	if !tt.HasText("Loading") {
		t.Fatalf("loading: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("uptime above 1 did not panic")
		}
	}()
	uptimeTester([]UptimeDay{{Uptime: 1.2}}, chart.ChartFrame{}, &pts)
}

func TestUptimeBarSnapshot(t *testing.T) {
	snapshot(t, "uptime_bar", 520, 110, func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			UptimeBar(c, uptimeSample(45), UptimeBarOptions{Frame: chart.ChartFrame{Label: "API uptime"}})
		})
	})
}

func TestUptimeBarManyDaysBounded(t *testing.T) {
	days := uptimeSample(100000)
	var pts []int
	tt := uptimeTester(days, chart.ChartFrame{}, &pts)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if pts[len(pts)-1] != 99999 {
		t.Fatal("end boundary")
	}
	for i := range 20 {
		tt.Key(0, ui.KeyRight)
		_ = i
	}
	if pts[len(pts)-1] != 99999 {
		t.Fatal("right boundary")
	}
}

func TestUptimeBarPixelBucketsPreserveOutages(t *testing.T) {
	days := uptimeSample(100000)
	var cache uptimeCache
	cache.prepare(days, 500)
	if len(cache.slots) > 500 || cache.slots[0].index != 7 {
		t.Fatalf("%d buckets, first worst %d", len(cache.slots), cache.slots[0].index)
	}
	cache.prepare(days, 1000)
	if len(cache.slots) > 1000 || cache.slots[len(cache.slots)-1].last != 99999 {
		t.Fatal("resized buckets")
	}
	for _, n := range []int{1, 3, 12} {
		cache.width = -1
		cache.prepare(days[:n], 500)
		if len(cache.slots) != n {
			t.Fatalf("%d days rendered as %d buckets", n, len(cache.slots))
		}
	}
}
