package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func paramSample(temp, topP *float64) []ModelParam {
	return []ModelParam{
		{Label: "Temperature", Value: temp, Min: 0, Max: 2, Step: 0.1, Default: 1},
		{Label: "Top-p", Value: topP, Min: 0, Max: 1, Step: 0.05, Default: 1},
	}
}

func TestParameterPanelAdjustAndReset(t *testing.T) {
	temp, topP, n, off := 1.0, 0.9, 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Disabled(off).Children(func() {
			if ParameterPanel(c, paramSample(&temp, &topP)).Changed() {
				n++
			}
		})
	}, 500, 400)
	for _, s := range []string{"1.0", "0.90", "Default 1.0", "2.0"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if temp != 1.1 || n != 1 {
		t.Fatalf("Right: temp %v changes %d", temp, n)
	}
	_ = tt.Click("Reset Temperature")
	if temp != 1 {
		t.Fatalf("reset temp %v", temp)
	}
	_ = tt.Click("Reset all")
	if topP != 1 || n != 3 {
		t.Fatalf("reset all: topP %v changes %d", topP, n)
	}
	off = true
	tt.Key(0, ui.KeyRight)
	if temp != 1 {
		t.Fatal("disabled panel changed")
	}
}

func TestParameterPanelBadParam(t *testing.T) {
	v := 0.0
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a default outside the range")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ParameterPanel(c, []ModelParam{{Label: "x", Value: &v, Min: 0, Max: 1, Step: 0.1, Default: 3}})
	}, 100, 100)
}

func TestParameterPanelSnapshot(t *testing.T) {
	temp, topP := 0.7, 1.0
	snapshot(t, "parameter_panel", 420, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ParameterPanel(c, paramSample(&temp, &topP)) })
	})
}
