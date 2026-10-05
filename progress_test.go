package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestProgressReducedMotionJumpsMidTransition(t *testing.T) {
	v := 0.0
	var bar ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Column(c).Width(300).Children(func() {
			bar = Progress(c, v, ProgressOptions{}).Bounds()
		})
	}, 400, 100)
	v = 1
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if got, want := progressPixel(tt, bar.X+bar.W*0.95, bar.Y+bar.H/2), theme.Light().AccentText; got != want {
		t.Fatalf("fill at 95%% is %v, want %v at once", got, want)
	}
}

func progressPixel(tt *ui.Tester, x, y float32) ui.Color {
	r, g, b, _ := tt.Image().At(int(x), int(y)).RGBA()
	return ui.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func TestProgressStaticFillContrast(t *testing.T) {
	for _, m := range []Mode{Light, Dark} {
		var bar ui.Rect
		var k theme.Tokens
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Mode: m})
			k = Tokens(c)
			ui.Column(c).Width(300).Children(func() {
				bar = Progress(c, 0, ProgressOptions{Indeterminate: true}).Bounds()
			})
		}, 400, 100)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		tt.Frame()                                              // Bounds report the previous layout
		fill := progressPixel(tt, bar.X+bar.W/2, bar.Y+bar.H-3) // the bar is the last row
		for _, against := range []ui.Color{k.Border, k.Background} {
			if cr := theme.Contrast(fill, against); cr < 3 {
				t.Errorf("%v: static fill %v vs %v is %.2f:1", m, fill, against, cr)
			}
		}
	}
}

func TestProgressValueText(t *testing.T) {
	v := 0.424
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Progress(c, v, ProgressOptions{Label: "Uploading", ShowValue: true})
		Progress(c, 1, ProgressOptions{Shape: ProgressRing, Label: "Indexing", ShowValue: true})
	}, 400, 300)
	if !tt.HasText("42%") || !tt.HasText("100%") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if r, ok := tt.Find("Indexing"); !ok || r.H < 48 {
		t.Fatalf("ring %v %v", r, ok)
	}
	v = 0
	tt.Frame()
	if !tt.HasText("0%") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestProgressIndeterminateReducedMotion(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Progress(c, math.NaN(), ProgressOptions{Indeterminate: true, Label: "Connecting"})
	}, 400, 300)
	if tt.HasText("In progress") {
		t.Fatal("static note shown while moving")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("In progress") {
		t.Fatal("no static note under reduced motion")
	}
	a := encode(t, tt.Image())
	tt.Frame()
	if string(a) != string(encode(t, tt.Image())) {
		t.Fatal("indeterminate bar moved under reduced motion")
	}
}

func TestProgressRingSize(t *testing.T) {
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Children(func() {
			w = Progress(c, 0.5, ProgressOptions{Shape: ProgressRing, Size: 24}).Bounds().W
		})
	}, 400, 300)
	tt.Frame()
	if w != 24 {
		t.Fatalf("ring width %v", w)
	}
}

func TestMeterLevels(t *testing.T) {
	disk := &ui.MeterLevels{Warning: 80, Critical: 95}
	battery := &ui.MeterLevels{Warning: 20, Critical: 10}
	for _, tc := range []struct {
		v      float64
		levels *ui.MeterLevels
		want   int
	}{
		{50, disk, 0}, {80, disk, 1}, {94.9, disk, 1}, {95, disk, 2}, {100, disk, 2},
		{50, battery, 0}, {20, battery, 1}, {10, battery, 2}, {0, battery, 2}, {99, nil, 0},
	} {
		if got := meterLevel(tc.v, tc.levels); got != tc.want {
			t.Errorf("meterLevel(%v, %+v) = %d, want %d", tc.v, tc.levels, got, tc.want)
		}
	}
}

func TestMeterWords(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		Meter(c, 90, MeterOptions{Max: 100, Levels: &ui.MeterLevels{Warning: 80, Critical: 95}, Label: "磁盘", ShowValue: true})
		Meter(c, 15, MeterOptions{Max: 100, Levels: &ui.MeterLevels{Warning: 20, Critical: 10}, Label: "电量"})
	}, 400, 300)
	if !tt.HasText("90 / 100 · 偏高") || !tt.HasText("偏低") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if _, ok := tt.Find("磁盘 (偏高)"); !ok {
		t.Fatal("meter name lacks its level")
	}
}

func TestProgressAndMeterPanics(t *testing.T) {
	for name, f := range map[string]func(c *ui.Context){
		"progress above 1":  func(c *ui.Context) { Progress(c, 1.2, ProgressOptions{}) },
		"progress NaN":      func(c *ui.Context) { Progress(c, math.NaN(), ProgressOptions{}) },
		"meter empty range": func(c *ui.Context) { Meter(c, 5, MeterOptions{Min: 10, Max: 10}) },
		"meter outside":     func(c *ui.Context) { Meter(c, 2, MeterOptions{}) },
		"meter levels": func(c *ui.Context) {
			Meter(c, 5, MeterOptions{Max: 10, Levels: &ui.MeterLevels{Warning: 8, Critical: 12}})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); f(c) }, 300, 200)
		}()
	}
}

func TestProgressSnapshot(t *testing.T) {
	snapshot(t, "progress", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			Progress(c, 0.62, ProgressOptions{Label: "Uploading 3 files", ShowValue: true})
			Progress(c, 0, ProgressOptions{Indeterminate: true, Label: "Connecting"})
			ui.Row(c).Gap(24).AlignItems(ui.Center).Children(func() {
				Progress(c, 0.35, ProgressOptions{Shape: ProgressRing, ShowValue: true})
				Progress(c, 0.8, ProgressOptions{Shape: ProgressRing, Size: 24, Label: "Indexing", ShowValue: true})
				Progress(c, 0, ProgressOptions{Shape: ProgressRing, Indeterminate: true, Label: "Waiting"})
			})
			Meter(c, 42, MeterOptions{Max: 100, Label: "Memory", ShowValue: true, Levels: &ui.MeterLevels{Warning: 70, Critical: 90}})
			Meter(c, 93, MeterOptions{Max: 100, Label: "Disk", ShowValue: true, Levels: &ui.MeterLevels{Warning: 70, Critical: 90}})
		})
	})
}
