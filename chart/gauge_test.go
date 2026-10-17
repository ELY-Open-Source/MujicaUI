package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var gaugeSegments = []GaugeSegment{{60, "Calm", GaugeGood}, {85, "Busy", GaugeWarning}, {100, "Crowded", GaugeDanger}}

func TestGaugeReadoutNamesSegment(t *testing.T) {
	full := false
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return Gauge(c, 72, GaugeOptions{Frame: ChartFrame{Label: "Nave"}, Max: 100, Unit: "%", Segments: gaugeSegments, Full: full})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Nave")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "72 % · Busy" {
		t.Fatalf("value %q", got)
	}
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Crowded: 85–100" {
		t.Fatalf("last segment %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	for _, full = range []bool{false, true} {
		tt.Frame()
		b := h.box(tt)
		cx := b.X + b.W/2
		var cy, rad float32
		if full {
			rad = 200/2 - 12
			cy = b.Y + 100
		} else {
			rad = min(b.W/2-12, 200-36)
			cy = b.Y + 8 + rad
		}
		tt.Move(cx-rad*0.7, cy-rad*0.6) // the left of the dial: the first segment
		if got := chartReadout(tt); got != "Calm: 0–60" {
			t.Fatalf("full=%v hover %q", full, got)
		}
	}
}

func TestGaugeBad(t *testing.T) {
	for name, o := range map[string]GaugeOptions{
		"range":    {Max: 0},
		"segments": {Max: 100, Segments: []GaugeSegment{{50, "a", 0}, {40, "b", 0}}},
		"end":      {Max: 100, Segments: []GaugeSegment{{50, "a", 0}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			o.Frame.Label = "x"
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); Gauge(c, 0, o) }, 300, 300)
		}()
	}
}

func TestGaugeSnapshot(t *testing.T) {
	snapshot(t, "gauge", 560, 300, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ui.Column(c).Width(256).Children(func() {
				Gauge(c, 72, GaugeOptions{Frame: ChartFrame{Label: "Half"}, Max: 100, Unit: "%", Segments: gaugeSegments})
			})
			ui.Column(c).Width(256).Children(func() {
				Gauge(c, 30, GaugeOptions{Frame: ChartFrame{Label: "Full"}, Max: 100, Unit: "%", Segments: gaugeSegments, Full: true})
			})
		})
	})
}
