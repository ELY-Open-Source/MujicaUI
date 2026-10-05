package feedback

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCountUpCountsAndSettles(t *testing.T) {
	value := 100.0
	tt, n := motionFrames(func(c *ui.Context) {
		CountUp(c, value, CountUpOptions{Duration: 300 * time.Millisecond})
	}, 200, 60)
	if !tt.HasText("100") || motionBusy(tt, n) {
		t.Fatalf("first value shown at once: %q", tt.Texts())
	}
	value = 200
	tt.Frame()
	time.Sleep(100 * time.Millisecond)
	tt.Frame()
	mid := ""
	for _, s := range tt.Texts() {
		if s != "200" {
			mid = s
		}
	}
	if mid == "" || mid == "100" {
		t.Fatalf("no value between: %q", tt.Texts())
	}
	time.Sleep(300 * time.Millisecond)
	tt.Frame()
	if motionCount(tt, "200") != 2 || motionBusy(tt, n) {
		t.Fatalf("did not settle: %q", tt.Texts())
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	value = 5
	tt.Frame()
	if motionCount(tt, "5") != 2 {
		t.Fatalf("reduced motion did not jump: %q", tt.Texts())
	}
}

func TestCountUpFormatAndInvalid(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CountUp(c, 12.345, CountUpOptions{Decimals: 1})
		CountUp(c, 3, CountUpOptions{Format: func(v float64) string { return fmt.Sprintf("%.0f items", v) }})
	}, 300, 80)
	if !tt.HasText("12.3") || !tt.HasText("3 items") {
		t.Fatalf("format: %q", tt.Texts())
	}
	for _, o := range []CountUpOptions{{Duration: -1}, {Decimals: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", o)
				}
			}()
			ui.Render(func(c *ui.Context) { core.Use(c, core.Settings{}); CountUp(c, 1, o) }, 100, 50, 1)
		}()
	}
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), "finite") {
			t.Fatal("NaN did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) { core.Use(c, core.Settings{}); CountUp(c, nanValue(), CountUpOptions{}) }, 100, 50, 1)
}

func nanValue() float64 {
	zero := 0.0
	return zero / zero
}

func TestCountUpSnapshot(t *testing.T) {
	snapshot(t, "count_up", 300, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			CountUp(c, 12840, CountUpOptions{}).FontSize(core.FontSize(c, 28)).Font("Georgia, serif")
			CountUp(c, 0.875, CountUpOptions{Format: func(v float64) string { return fmt.Sprintf("%.1f%%", v*100) }})
		})
	})
}

func TestCountUpFiniteExtremeInterpolation(t *testing.T) {
	value := -math.MaxFloat64
	tt, _ := motionFrames(func(c *ui.Context) {
		CountUp(c, value, CountUpOptions{Duration: time.Second, Format: func(v float64) string {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				t.Fatal("finite endpoints produced non-finite text")
			}
			return fmt.Sprintf("%g", v)
		}})
	}, 300, 80)
	value = math.MaxFloat64
	tt.Frame()
}

func TestCountUpOffscreenStopsFrames(t *testing.T) {
	value := 1.0
	tt, n := motionFrames(func(c *ui.Context) {
		ui.Column(c).Height(80).Clip().Children(func() {
			ui.Box(c).Size(10, 200).Shrink(0)
			CountUp(c, value, CountUpOptions{Duration: time.Second})
		})
	}, 300, 80)
	value = 2
	tt.Frame()
	if motionBusy(tt, n) {
		t.Fatal("offscreen CountUp requested continuous rebuilds")
	}
}

func TestCountUpRestoredMotionKeepsFinalValue(t *testing.T) {
	value := 0.0
	shown := 0.0
	tt, n := motionFrames(func(c *ui.Context) {
		CountUp(c, value, CountUpOptions{Duration: 5 * time.Second, Format: func(v float64) string { shown = v; return fmt.Sprintf("%g", v) }})
	}, 300, 100)
	value = 100
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	if shown != 100 || motionBusy(tt, n) {
		t.Fatalf("restored value %v", shown)
	}
}
