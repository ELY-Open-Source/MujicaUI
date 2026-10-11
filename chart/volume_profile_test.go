package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestVolumeProfileOfAndValueArea(t *testing.T) {
	cs := candlesOf([5]float64{10, 14, 10, 12, 400}, [5]float64{12, 12, 12, 12, 50}, [5]float64{11, 12, 10, 11, 100})
	bins := VolumeProfileOf(cs, 4) // bins of 1: 10-11, 11-12, 12-13, 13-14
	want := []float64{100 + 50, 100 + 50, 100 + 50, 100}
	for i, b := range bins {
		if math.Abs(b.Volume-want[i]) > 1e-9 || b.Low != 10+float64(i) {
			t.Fatalf("bin %d %+v, want volume %v", i, b, want[i])
		}
	}
	vb := []ProfileBin{{0, 1, 5}, {1, 2, 10}, {2, 3, 40}, {3, 4, 30}, {4, 5, 15}}
	poc, lo, hi := ValueArea(vb, 0.7)
	// 40, then 30 above: 70 of 100 is enough.
	if poc != 2 || lo != 2 || hi != 3 {
		t.Fatalf("value area %d %d..%d", poc, lo, hi)
	}
	if _, _, hi := ValueArea(vb, 1); hi != 4 {
		t.Fatal("full value area")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("zero bins accepted")
		}
	}()
	VolumeProfileOf(cs, 0)
}

func TestVolumeProfileKeyboard(t *testing.T) {
	vb := []ProfileBin{{0, 1, 5}, {1, 2, 10}, {2, 3, 40}, {3, 4, 30}, {4, 5, 15}}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = VolumeProfile(c, vb, VolumeProfileOptions{Frame: ChartFrame{Label: "Profile"}})
	}, 500, 360)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if res.Point != 2 || !tt.HasText("2.00–3.00 · Vol 40 · POC · Value area") {
		t.Fatalf("POC readout %d %q", res.Point, tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	if !tt.HasText("1.00–2.00 · Vol 10") || tt.HasText("1.00–2.00 · Vol 10 · Value area") {
		t.Fatalf("outside readout %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("overlapping bins accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		VolumeProfile(c, []ProfileBin{{0, 2, 1}, {1, 3, 1}}, VolumeProfileOptions{Frame: ChartFrame{Label: "x"}})
	}, 200, 100)
}

func TestVolumeProfileSnapshot(t *testing.T) {
	bins := VolumeProfileOf(sampleCandles(200), 18)
	snapshot(t, "volume_profile", 480, 380, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			VolumeProfile(c, bins, VolumeProfileOptions{Frame: ChartFrame{Label: "Profile", Height: 290}})
		})
	})
}
