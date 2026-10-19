package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func c2Demo(t *testing.T, id string) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(charts2Pages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("missing page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(24).Children(func() { charts2Pages[i].Demo(c, func(s string) { *events = append(*events, s) }) })
	}, 900, 760)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt, events
}

func TestCharts2ExistingSet(t *testing.T) {
	want := []string{"121", "122", "123", "124", "125", "126", "127", "128", "129", "130", "131", "132", "133", "138", "140"}
	var ids []string
	for _, p := range charts2Pages {
		ids = append(ids, p.ID)
		if p.Category != "Charts II" || p.Code == "" || p.Keys == "" {
			t.Fatalf("incomplete page %s", p.ID)
		}
		c2Demo(t, p.ID)
	}
	if !slices.Equal(ids, want) {
		t.Fatalf("page IDs %v", ids)
	}
}

func TestCharts2ReadoutsAndControlsLog(t *testing.T) {
	for _, id := range []string{"122", "123", "124", "125", "127", "129", "130", "131", "132", "138"} {
		tt, events := c2Demo(t, id)
		p := charts2Pages[slices.IndexFunc(charts2Pages, func(p page) bool { return p.ID == id })]
		label := p.Name
		if id == "132" {
			label = "Alms"
		}
		tabUntil(t, tt, label)
		tt.Key(0, ui.KeyRight)
		if !slices.Contains(*events, p.Name+": point 0") {
			t.Errorf("%s readout events %q", id, *events)
		}
		if err := tt.Click("Disabled"); err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(*events, "Disabled = true") {
			t.Errorf("%s disabled events %q", id, *events)
		}
	}
}

func TestCharts2LegendBrushAnnotationAndRetry(t *testing.T) {
	tt, events := c2Demo(t, "121")
	if err := tt.Click("+3 more"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Matins"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(*events, "ChartLegend: expanded true") || !slices.Contains(*events, "ChartLegend: toggled Matins") {
		t.Fatalf("legend events %q", *events)
	}
	tt, events = c2Demo(t, "126")
	tabUntil(t, tt, "Months: 20 to 60")
	tt.Key(0, ui.KeyRight)
	if !slices.Contains(*events, "ChartBrush: range 21–61") {
		t.Fatalf("brush events %q", *events)
	}
	tt, events = c2Demo(t, "127")
	tabUntil(t, tt, "ChartAnnotation")
	tt.Key(ui.Alt, ui.KeyUp)
	if !slices.Contains(*events, "ChartAnnotation: moved") || !tt.HasText("Goal 51.0 · Advent 3.0 · Lent 6.0–8.0 · Feast 1.0, 90.0") {
		t.Fatalf("annotation events %q text %q", *events, tt.Texts())
	}
	tt, events = c2Demo(t, "140")
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(*events, "ChartEmptyState: retry") || !tt.HasText("Loading") {
		t.Fatalf("retry events %q", *events)
	}
}

func TestCharts2ExportAndTrendControls(t *testing.T) {
	tt, events := c2Demo(t, "128")
	tabUntil(t, tt, "Export PNG")
	tt.Key(0, ui.KeyEnter)
	if len(*events) == 0 || !strings.HasPrefix((*events)[len(*events)-1], "ChartExport: ") || !strings.HasSuffix((*events)[len(*events)-1], " bytes") {
		t.Fatalf("export events %q", *events)
	}
	tt, events = c2Demo(t, "133")
	if err := tt.Click("Invert"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(*events, "Invert = true") {
		t.Fatalf("trend events %q", *events)
	}
}
