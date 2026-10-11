package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

type leverageProbe struct {
	value   float64
	opts    LeverageSliderOptions
	applied int
}

func (p *leverageProbe) tester(density theme.Density) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Density: density})
		ui.Column(c).Padding(10).Children(func() {
			if LeverageSlider(c, &p.value, p.opts).Changed() {
				p.applied++
			}
		})
	}, 420, 200)
}

func TestLeverageSliderAppliesOnlyOnConfirm(t *testing.T) {
	p := &leverageProbe{value: 5, opts: LeverageSliderOptions{Warning: 10, Danger: 50}}
	tt := p.tester(theme.Default)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Leverage") {
		t.Fatal("slider not focusable")
	}
	tt.Key(0, ui.KeyRight)
	if p.value != 5 || !tt.HasText("10×") || !tt.HasText("High risk") {
		t.Fatalf("draft: value %v texts %q", p.value, tt.Texts())
	}
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("Very high risk") {
		t.Fatalf("danger zone: %q", tt.Texts())
	}
	if err := tt.Click("Apply"); err != nil {
		t.Fatal(err)
	}
	if p.value != 100 || p.applied != 1 {
		t.Fatalf("apply: %v %d", p.value, p.applied)
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if p.value != 100 || !tt.HasText("1×") {
		t.Fatalf("home moved value: %v", p.value)
	}
}

func TestLeverageSliderPointerAndHitArea(t *testing.T) {
	p := &leverageProbe{value: 1, opts: LeverageSliderOptions{Levels: []float64{1, 2, 3}}}
	tt := p.tester(theme.Compact)
	// The track is at least 28 DIPs tall, even Compact.
	tt.ClickAt(10+400-10, 70)
	if !tt.HasText("3×") {
		t.Fatalf("pointer: %q", tt.Texts())
	}
	tt.ClickAt(30, 70)
	if !tt.HasText("1×") || p.applied != 0 {
		t.Fatalf("pointer back: %q", tt.Texts())
	}
}

func TestLeverageSliderDisabledAndBadValue(t *testing.T) {
	p := &leverageProbe{value: 2, opts: LeverageSliderOptions{Disabled: true}}
	tt := p.tester(theme.Default)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	_ = tt.Click("Apply")
	if p.value != 2 || tt.HasText("3×") || !tt.HasText("2×") {
		t.Fatal("disabled slider moved")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("value outside Levels accepted")
		}
	}()
	(&leverageProbe{value: 4}).tester(theme.Default)
}

func TestLeverageSliderSnapshot(t *testing.T) {
	v := 20.0
	snapshot(t, "leverage_slider", 360, 160, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			LeverageSlider(c, &v, LeverageSliderOptions{Warning: 10, Danger: 50})
		})
	})
}
