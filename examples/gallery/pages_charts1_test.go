package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// c1Demo runs page id's demo with an event log.
func c1Demo(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(charts1Pages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("no page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(24).Children(func() {
			charts1Pages[i].Demo(c, func(s string) { *events = append(*events, s) })
		})
	}, 900, 760)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt, events
}

// c1Labels names each demo's chart for the keys.
var c1Labels = map[string]string{
	"101": "Attendance", "102": "Attendance", "103": "Tithes", "104": "Stores", "105": "Choir", "106": "Skills",
	"107": "Hours", "108": "Prayers", "109": "Ages", "110": "Hymns", "111": "Ledger", "112": "Order", "113": "Alms",
	"114": "Bells 1", "115": "Nave", "116": "Rooms", "117": "Rooms", "118": "Money", "119": "House", "120": "Sound",
}

func TestCharts1DemosLogReadouts(t *testing.T) {
	for _, p := range charts1Pages {
		if p.Category != "Charts I" || p.Code == "" || p.Keys == "" {
			t.Fatalf("%s: incomplete page", p.ID)
		}
		tt, events := c1Demo(t, p.ID)
		tabUntil(t, tt, c1Labels[p.ID])
		tt.Key(0, ui.KeyRight)
		if len(*events) == 0 || !strings.HasSuffix((*events)[len(*events)-1], "point 0") {
			t.Errorf("%s: keyboard readout not logged: %q", p.ID, *events)
		}
	}
}

func TestCharts1TreemapDrillLogs(t *testing.T) {
	tt, events := c1Demo(t, "116")
	tabUntil(t, tt, "Rooms")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if !slices.Contains(*events, "Treemap: path /Library") {
		t.Fatalf("drill not logged: %q", *events)
	}
}

func tabUntil(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 40 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("Tab never reached %q", label)
}
