package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

func legendItems() []ChartLegendItem {
	pal := internal.Palette(theme.Light())
	var out []ChartLegendItem
	for i, n := range []string{"Matins", "Lauds", "Prime", "Terce", "Sext"} {
		out = append(out, ChartLegendItem{Name: n, Color: pal[i], Shape: i})
	}
	return out
}

func TestChartLegendToggleKeyboard(t *testing.T) {
	hidden := map[string]bool{}
	var toggled []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		if s, ok := ChartLegend(c, &hidden, legendItems(), ChartLegendOptions{}).Toggled(); ok {
			toggled = append(toggled, s)
		}
	}, 500, 200)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Matins") {
		t.Fatal("Tab did not enter the legend")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if !hidden["Prime"] || len(toggled) != 1 || toggled[0] != "Prime" {
		t.Fatalf("arrows and Space: hidden %v toggled %v", hidden, toggled)
	}
	if err := tt.Click("Prime"); err != nil || hidden["Prime"] {
		t.Fatalf("click did not show Prime again: %v", hidden)
	}
	tt.Key(0, ui.KeyTab)
	if tt.Focused("Lauds") || tt.Focused("Sext") {
		t.Fatal("legend entries are more than one Tab stop")
	}
}

func TestChartLegendFoldAndStatic(t *testing.T) {
	vertical := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ChartLegend(c, nil, legendItems(), ChartLegendOptions{Vertical: vertical, MaxItems: 2})
	}, 500, 400)
	if tt.HasText("Prime") || !tt.HasText("+3 more") {
		t.Fatalf("fold: %q", tt.Texts())
	}
	a, _ := tt.Find("Matins")
	b, _ := tt.Find("Lauds")
	if b.Y <= a.Y || b.X != a.X {
		t.Fatalf("vertical layout %v %v", a, b)
	}
	if err := tt.Click("+3 more"); err != nil || !tt.HasText("Sext") || !tt.HasText("Show fewer") {
		t.Fatalf("unfold: %v %q", err, tt.Texts())
	}
	if _, ok := tt.Find("Matins"); !ok || tt.Focused("Matins") {
		t.Fatal("static legend took focus")
	}
}

func TestChartLegendBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for duplicate names")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ChartLegend(c, nil, []ChartLegendItem{{Name: "a"}, {Name: "a"}}, ChartLegendOptions{})
	}, 100, 100)
}

func TestChartLegendSnapshot(t *testing.T) {
	snapshot(t, "chart_legend", 520, 120, func(c *ui.Context) {
		hidden := map[string]bool{"Lauds": true}
		items := legendItems()
		for i := range items {
			items[i].Color = internal.Palette(Tokens(c))[i]
		}
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			ChartLegend(c, &hidden, items, ChartLegendOptions{})
			ChartLegend(c, nil, items, ChartLegendOptions{MaxItems: 3})
		})
	})
}
