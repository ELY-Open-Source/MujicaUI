package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestFeedbackTextFollowsTextScale(t *testing.T) {
	items := []Notification{{ID: "n", Title: "Backup complete", Read: true}, {ID: "u", Title: "Unread one"}}
	var empty, result float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
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
		"spinner negative size":  func(c *ui.Context) { Spinner(c, SpinnerOptions{Size: -4}) },
		"status unknown labeled": func(c *ui.Context) { StatusIndicator(c, 42, StatusIndicatorOptions{Label: "Custom"}) },
		"alert severity":         func(c *ui.Context) { Alert(c, nil, "x", AlertOptions{Severity: 9}) },
		"toast severity":         func(c *ui.Context) { ShowToast(c, ToastMessage{Message: "x", Severity: -1}) },
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

func TestNotificationTitleOnlyHitHeight(t *testing.T) {
	items := []Notification{{ID: "n", Title: "Backup complete", Read: true}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		NotificationCenter(c, &items, NotificationCenterOptions{})
	}, 600, 400)
	if r, ok := tt.Find("Backup complete"); !ok || r.H < theme.MinHitSize {
		t.Fatalf("title-only notification hit box %v %v", r, ok)
	}
}
