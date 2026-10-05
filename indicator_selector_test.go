package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestIndicatorSelectorSearchAddRemoveAndDisabled(t *testing.T) {
	var overlays []IndicatorOverlay
	var panes IndicatorPanes
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		IndicatorSelector(c, &overlays, &panes, IndicatorSelectorOptions{Disabled: disabled})
	}, 640, 740)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Indicators (0)"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Search indicators"); err != nil {
		t.Fatal(err)
	}
	tt.Type("RSI")
	if tt.HasText("Add MA") || !tt.HasText("RSI") {
		t.Fatalf("search %q", tt.Texts())
	}
	if err := tt.Click("Add RSI"); err != nil {
		t.Fatal(err)
	}
	if len(panes.Panes) != 1 || panes.Panes[0].Kind != IndicatorRSI {
		t.Fatalf("panes %+v", panes.Panes)
	}
	if err := tt.Click("Remove RSI(14)"); err != nil {
		t.Fatal(err)
	}
	if len(panes.Panes) != 0 {
		t.Fatal("pane not removed")
	}
	disabled = true
	tt.Frame()
	if tt.HasText("Search indicators") {
		t.Fatal("disabled panel remained open")
	}
}
func TestIndicatorSelectorSnapshot(t *testing.T) {
	overlays := []IndicatorOverlay{{Kind: IndicatorMA}}
	panes := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorRSI}}}
	snapshot(t, "indicator_selector", 360, 100, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() { IndicatorSelector(c, &overlays, &panes, IndicatorSelectorOptions{}) })
	})
}

func TestIndicatorSelectorOverlayParams(t *testing.T) {
	var overlays []IndicatorOverlay
	var panes IndicatorPanes
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		IndicatorSelector(c, &overlays, &panes, IndicatorSelectorOptions{})
	}, 640, 740)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Indicators (0)"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Add MA"); err != nil {
		t.Fatal(err)
	}
	if len(overlays) != 1 || overlays[0].Kind != IndicatorMA || overlays[0].Params[0] != 20 {
		t.Fatalf("overlay %+v", overlays)
	}
	if err := tt.Click("MA 1"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyUp)
	if overlays[0].Params[0] != 21 {
		t.Fatalf("params %v", overlays[0].Params)
	}
}
