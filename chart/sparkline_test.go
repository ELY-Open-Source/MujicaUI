package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

var sparkValues = []float64{3, 5, 4, math.NaN(), 7, 6, 9}

func TestSparklineEndValueAndKeys(t *testing.T) {
	kind := SparklineLine
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return Sparkline(c, sparkValues, SparklineOptions{Frame: ChartFrame{Label: "Bells"}, Kind: kind})
	}}
	tt := ui.NewTester(h.view, 300, 200)
	if got := chartReadout(tt); got != "9" {
		t.Fatalf("end value %q", got)
	}
	r, _ := tt.Find("Bells")
	if r.H < 28 || r.W != 96 {
		t.Fatalf("sparkline box %v", r)
	}
	tabTo(t, tt, "Bells")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "5" || h.res.Point != 1 {
		t.Fatalf("read %q at %d", got, h.res.Point)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "no value" {
		t.Fatalf("missing %q", got)
	}
	for _, kind = range []SparklineKind{SparklineBar, SparklineArea} {
		tt.Frame()
		tt.Move(r.X+r.W-2, r.Y+r.H/2)
		if got := chartReadout(tt); got != "9" || h.res.Point != 6 {
			t.Fatalf("kind %d hover %q", kind, got)
		}
	}
}

func TestSparklineStates(t *testing.T) {
	status := data.DataLoading
	var vals []float64
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return Sparkline(c, vals, SparklineOptions{Frame: ChartFrame{Label: "Bells", Status: status}})
	}}
	tt := ui.NewTester(h.view, 300, 200)
	status = data.DataFailed
	tt.Frame()
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	status = data.DataReady
	tt.Frame()
	if !tt.HasText("No data") {
		t.Fatalf("empty %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a label")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); Sparkline(c, vals, SparklineOptions{}) }, 100, 100)
}

func TestSparklineSnapshot(t *testing.T) {
	snapshot(t, "sparkline", 420, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(24).Children(func() {
			for _, k := range []SparklineKind{SparklineLine, SparklineBar, SparklineArea} {
				Sparkline(c, sparkValues, SparklineOptions{Frame: ChartFrame{Label: "Bells"}, Kind: k})
			}
		})
	})
}
