package main

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func mktDemoTester(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(12).Children(func() { mktDemo(c, func(s string) { *events = append(*events, s) }, id) })
	}, 900, 800)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt, events
}

func TestMarketChartPagesHaveInteractiveControls(t *testing.T) {
	if len(mktchartPages) != 20 {
		t.Fatalf("pages %d", len(mktchartPages))
	}
	for _, p := range mktchartPages {
		t.Run(p.ID, func(t *testing.T) {
			if p.Category != "Market Charts" || p.Code == "" || p.Keys == "" {
				t.Fatal("incomplete page")
			}
			tt, events := mktDemoTester(t, p.ID)
			if err := tt.Click("Red up"); err != nil {
				t.Fatal(err)
			}
			if len(*events) == 0 || (*events)[len(*events)-1] != "Red up: true" {
				t.Fatalf("control events %q", *events)
			}
			if err := tt.Click("Disabled"); err != nil {
				t.Fatal(err)
			}
			if (*events)[len(*events)-1] != "Disabled: true" {
				t.Fatalf("disabled event %q", *events)
			}
		})
	}
}

func TestMarketChartDemosLogKeyboardAndAlerts(t *testing.T) {
	for _, id := range []string{"141", "142", "143", "144", "145", "146", "147", "153", "155", "156", "157", "158", "159"} {
		t.Run(id, func(t *testing.T) {
			tt, events := mktDemoTester(t, id)
			tabUntil(t, tt, "BTC")
			tt.Key(0, ui.KeyRight)
			want := id + ": point 0"
			found := false
			for _, e := range *events {
				found = found || e == want
			}
			if !found {
				t.Fatalf("keyboard event %q, want %s", *events, want)
			}
		})
	}
	tt, events := mktDemoTester(t, "159")
	if err := tt.Click("Append crossing candle"); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range *events {
		found = found || strings.HasPrefix(e, "Alert target crossed ")
	}
	if !found {
		t.Fatalf("alert events %q", *events)
	}
}
