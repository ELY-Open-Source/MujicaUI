package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestBadgeText(t *testing.T) {
	for _, tc := range []struct {
		opts        BadgeOptions
		shown, full string
	}{
		{BadgeOptions{Count: 7}, "7", "7"},
		{BadgeOptions{Count: 99}, "99", "99"},
		{BadgeOptions{Count: 100}, "99+", "100"},
		{BadgeOptions{Count: 12, Max: 9}, "9+", "12"},
		{BadgeOptions{Count: 0}, "", ""},
		{BadgeOptions{Text: "New", Count: 5}, "New", "New"},
	} {
		shown, full := badgeText(tc.opts)
		if shown != tc.shown || full != tc.full {
			t.Errorf("%+v: %q %q", tc.opts, shown, full)
		}
	}
}

func TestBadgeReadsFullValue(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Badge(c, BadgeOptions{Count: 1204, Label: "Unread"}, nil)
	}, 300, 100)
	tt.Frame()
	if !tt.HasText("99+") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if _, ok := tt.Find("Unread 1204"); !ok {
		t.Fatal("full value not read")
	}
}

func TestBadgeCornerPosition(t *testing.T) {
	for _, tc := range []struct {
		pos         BadgePosition
		right, down bool
	}{{BadgeTopRight, true, false}, {BadgeTopLeft, false, false}, {BadgeBottomRight, true, true}, {BadgeBottomLeft, false, true}} {
		var host ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			ui.Column(c).Padding(40).Children(func() {
				Badge(c, BadgeOptions{Count: 3, Position: tc.pos}, func() {
					host = ui.Box(c).Size(60, 60).Bounds()
				})
			})
		}, 300, 300)
		tt.Frame()
		r, ok := tt.Find("3")
		if !ok {
			t.Fatalf("%v: badge missing", tc.pos)
		}
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		wantX, wantY := host.X, host.Y
		if tc.right {
			wantX += host.W
		}
		if tc.down {
			wantY += host.H
		}
		if d := cx - wantX; d < -2 || d > 2 {
			t.Errorf("%v: x %v, want %v", tc.pos, cx, wantX)
		}
		if d := cy - wantY; d < -2 || d > 2 {
			t.Errorf("%v: y %v, want %v", tc.pos, cy, wantY)
		}
	}
}

func TestBadgeBadOptionsPanic(t *testing.T) {
	for i, build := range []func(c *ui.Context){
		func(c *ui.Context) { Badge(c, BadgeOptions{Count: -1}, nil) },
		func(c *ui.Context) { Badge(c, BadgeOptions{Count: 1, Position: BadgeTopRight}, nil) },
		func(c *ui.Context) { Badge(c, BadgeOptions{Count: 1}, func() {}) },
		func(c *ui.Context) { Badge(c, BadgeOptions{Count: 1, Tone: BadgeTone(42)}, nil) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("case %d: no panic", i)
				}
			}()
			tt := ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				build(c)
			}, 100, 100)
			tt.Frame()
		}()
	}
}

func TestBadgeSnapshot(t *testing.T) {
	snapshot(t, "badge", 460, 100, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(14).AlignItems(ui.Center).Children(func() {
			Badge(c, BadgeOptions{Count: 7}, nil)
			Badge(c, BadgeOptions{Count: 250}, nil)
			Badge(c, BadgeOptions{Text: "Draft", Tone: BadgeNeutral}, nil)
			Badge(c, BadgeOptions{Text: "Live", Tone: BadgeSuccess}, nil)
			Badge(c, BadgeOptions{Text: "Late", Tone: BadgeWarning}, nil)
			Badge(c, BadgeOptions{Text: "Info", Tone: BadgeInfo}, nil)
			Badge(c, BadgeOptions{Count: 3, Tone: BadgeDanger, Position: BadgeTopRight}, func() {
				Avatar(c, "Sakiko Togawa", AvatarOptions{})
			})
		})
	})
}

func TestBadgeTextScale(t *testing.T) {
	var w [2]float32
	for i, s := range []float32{1, 2} {
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			ui.Row(c).AlignItems(ui.Start).Children(func() {
				w[i] = Badge(c, BadgeOptions{Text: "Draft"}, nil).Bounds().W
			})
		}, 400, 200)
		tt.SetPreferences(ui.Preferences{TextScale: s})
		tt.Frame()
	}
	if w[1] < w[0]*1.5 {
		t.Fatalf("width %v at 2x, %v at 1x", w[1], w[0])
	}
}
