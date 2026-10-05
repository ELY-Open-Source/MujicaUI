package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestPriorityIndicatorNamesAndShapes(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light, Locale: ZhCN})
		ui.Column(c).Children(func() {
			for p := PriorityNone; p <= PriorityUrgent; p++ {
				PriorityIndicator(c, p, PriorityIndicatorOptions{})
			}
			PriorityIndicator(c, PriorityHigh, PriorityIndicatorOptions{IconOnly: true})
		})
	}, 200, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for _, s := range []string{"无优先级", "低", "中", "高", "紧急"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q: %v", s, tt.Texts())
		}
	}
}

func TestPriorityIndicatorBarsShowLevel(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Column(c).Children(func() {
			for p := PriorityLow; p <= PriorityHigh; p++ {
				ui.Box(c).Height(20).Children(func() { PriorityIndicator(c, p, PriorityIndicatorOptions{IconOnly: true}) })
			}
		})
	}, 100, 100)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("Medium") {
		t.Fatal("icon-only indicator lost its name")
	}
	img := tt.Image()
	k := theme.Light()
	for i, want := range []struct {
		n   int
		col ui.Color
	}{{1, k.TextMuted}, {2, k.TextMuted}, {3, k.Warning}} {
		runs, on := 0, false
		for x := 0; x < 16; x++ {
			hit := datetimeNear(img.RGBAAt(x, i*20+13), want.col)
			if hit && !on {
				runs++
			}
			on = hit
		}
		if runs != want.n {
			t.Errorf("priority %d: %d filled bars, want %d", i+1, runs, want.n)
		}
	}
}

func TestPriorityIndicatorInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); PriorityIndicator(c, 9, PriorityIndicatorOptions{}) }, 10, 10)
}

func TestPriorityIndicatorSnapshot(t *testing.T) {
	snapshot(t, "priority_indicator", 520, 60, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			for p := PriorityNone; p <= PriorityUrgent; p++ {
				PriorityIndicator(c, p, PriorityIndicatorOptions{})
			}
		})
	})
}
