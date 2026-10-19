package feedback

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestFeedbackTextFollowsTextScale(t *testing.T) {
	items := []Notification{{ID: "n", Title: "Backup complete", Read: true}, {ID: "u", Title: "Unread one"}}
	var empty, result float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			Progress(c, 0.5, ProgressOptions{ShowValue: true})
			empty = EmptyState(c, EmptyStateOptions{Title: "Empty title", Description: "Short", NoIcon: true}).Element.Bounds().H
			result = Result(c, ResultOptions{Title: "Result title"}).Element.Bounds().H
			StatusIndicator(c, StatusSuccess, StatusIndicatorOptions{})
			NotificationCenter(c, &items, NotificationCenterOptions{})
		})
	}, 900, 2400)
	tt.Frame()
	empty1, result1 := empty, result
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame() // Bounds report the previous layout
	for text, min := range map[string]float32{
		"50%":      theme.CaptionSize * 2,
		"Healthy":  theme.BodySize * 2,
		"1 unread": theme.CaptionSize * 2,
	} {
		if r, ok := tt.Find(text); !ok || r.H < min {
			t.Errorf("%q at 200%%: %v %v, want height >= %v", text, r, ok, min)
		}
	}
	// Titles are found only through their labeled groups: the groups must grow by the titles too.
	if grow := empty - empty1; grow < theme.BodySize*theme.LineHeight+theme.H3Size {
		t.Errorf("EmptyState grew %v at 200%%; its title did not scale", grow)
	}
	if grow := result - result1; grow < theme.H2Size {
		t.Errorf("Result grew %v at 200%%; its title did not scale", grow)
	}
}

func TestFeedbackRejectsInvalidInput(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	for name, f := range map[string]func(c *ui.Context){
		"progress shape":       func(c *ui.Context) { Progress(c, 0.5, ProgressOptions{Shape: 7}) },
		"progress NaN size":    func(c *ui.Context) { Progress(c, 0.5, ProgressOptions{Size: float32(nan)}) },
		"meter infinite range": func(c *ui.Context) { Meter(c, 1, MeterOptions{Max: inf}) },
		"meter NaN min":        func(c *ui.Context) { Meter(c, 1, MeterOptions{Min: nan, Max: 2}) },
		"meter NaN threshold": func(c *ui.Context) {
			Meter(c, 1, MeterOptions{Max: 2, Levels: &ui.MeterLevels{Warning: nan, Critical: 2}})
		},
		"skeleton shape":         func(c *ui.Context) { Skeleton(c, SkeletonOptions{Shape: 5}) },
		"skeleton NaN width":     func(c *ui.Context) { Skeleton(c, SkeletonOptions{Width: float32(nan)}) },
		"spinner negative size":  func(c *ui.Context) { core.Spinner(c, core.SpinnerOptions{Size: -4}) },
		"status unknown labeled": func(c *ui.Context) { StatusIndicator(c, 42, StatusIndicatorOptions{Label: "Custom"}) },
		"alert severity":         func(c *ui.Context) { Alert(c, nil, "x", AlertOptions{Severity: 9}) },
		"toast severity":         func(c *ui.Context) { ShowToast(c, ToastMessage{Message: "x", Severity: -1}) },
		"meter negative width":   func(c *ui.Context) { Meter(c, 1, MeterOptions{Max: 2, Width: -1}) },
		"meter NaN width":        func(c *ui.Context) { Meter(c, 1, MeterOptions{Max: 2, Width: float32(nan)}) },
		"meter infinite width":   func(c *ui.Context) { Meter(c, 1, MeterOptions{Max: 2, Width: float32(inf)}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); f(c) }, 300, 200)
		}()
	}
}

func TestNotificationTitleOnlyHitHeight(t *testing.T) {
	items := []Notification{{ID: "n", Title: "Backup complete", Read: true}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NotificationCenter(c, &items, NotificationCenterOptions{})
	}, 600, 400)
	if r, ok := tt.Find("Backup complete"); !ok || r.H < theme.MinHitSize {
		t.Fatalf("title-only notification hit box %v %v", r, ok)
	}
}

func TestNotificationHeaderWrapsAtDoubleText(t *testing.T) {
	items := []Notification{{ID: "a", Title: "Build failed"}, {ID: "b", Title: "Backup complete", Read: true}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(520).Height(700).Children(func() {
			NotificationCenter(c, &items, NotificationCenterOptions{})
		})
	}, 600, 760)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	for _, s := range []string{"Notifications", "1 unread", "Mark all read", "Clear all"} {
		if r, ok := tt.Find(s); !ok || r.W == 0 || r.H == 0 || r.X < 0 || r.X+r.W > 520 {
			t.Errorf("%q at 200%%: %v %v", s, r, ok)
		}
	}
	if err := tt.Click("Clear all"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("Clear all at 200%% left %d items", len(items))
	}
}

func TestProgressRingCaptionFitsAtDoubleText(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Children(func() {
			Progress(c, 1, ProgressOptions{Shape: ProgressRing, ShowValue: true})
		})
	}, 400, 200)
	for _, scale := range []float32{1, 2} {
		tt.SetPreferences(ui.Preferences{TextScale: scale})
		tt.Frame()
		ring, ok := tt.Find("Loading")
		text, ok2 := tt.Find("100%")
		if !ok || !ok2 {
			t.Fatalf("scale %v: ring %v text %v", scale, ok, ok2)
		}
		stroke := max(float32(3), ring.W/10)
		inside := text.X >= ring.X+stroke && text.X+text.W <= ring.X+ring.W-stroke
		outside := text.X >= ring.X+ring.W
		if !inside && !outside {
			t.Errorf("scale %v: caption %v clipped by ring %v", scale, text, ring)
		}
		if scale == 1 && !inside {
			t.Errorf("caption left the ring at 100%%: %v", text)
		}
	}
}

func TestProgressResumesFromDisplayedValueAfterReducedMotion(t *testing.T) {
	v := 0.0
	var bar ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Width(300).Children(func() {
			bar = Progress(c, v, ProgressOptions{}).Bounds()
		})
	}, 400, 100)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	v = 0.5
	tt.Frame()
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	if got, want := progressPixel(tt, bar.X+bar.W*0.45, bar.Y+bar.H/2), theme.Light().AccentText; got != want {
		t.Fatalf("after motion returned the fill at 45%% is %v, want %v", got, want)
	}
	v = 1
	tt.Frame()
	if got := progressPixel(tt, bar.X+bar.W*0.95, bar.Y+bar.H/2); got == theme.Light().AccentText {
		t.Fatal("restored motion did not animate the next change")
	}
}
