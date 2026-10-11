package chart

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
)

func axisIn(opts ChartAxisOptions) (a ChartAxisPart) {
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); a = ChartAxis(c, opts) }, 100, 100)
	return a
}

func TestChartAxisLinearAndLog(t *testing.T) {
	plot := ui.Rect{X: 10, Y: 0, W: 100, H: 50}
	lin := axisIn(ChartAxisOptions{Min: 0, Max: 97})
	if fmt.Sprint(lin.Ticks()) != "[0 20 40 60 80]" || lin.Map(plot, 48.5) != 60 || lin.Label(40) != "40" {
		t.Fatalf("linear %v %v", lin.Ticks(), lin.Map(plot, 48.5))
	}
	if v := lin.Invert(plot, 60); math.Abs(v-48.5) > 1e-9 {
		t.Fatalf("invert %v", v)
	}
	log := axisIn(ChartAxisOptions{Kind: AxisLog, Side: SideLeft, Min: 1, Max: 10000})
	if fmt.Sprint(log.Ticks()) != "[1 10 100 1000 10000]" || log.Map(plot, 100) != 25 || log.Label(1000) != "1000" {
		t.Fatalf("log %v %v", log.Ticks(), log.Map(plot, 100))
	}
	if short := axisIn(ChartAxisOptions{Kind: AxisLog, Min: 1, Max: 50}); fmt.Sprint(short.Ticks()) != "[1 2 5 10 20 50]" {
		t.Fatalf("short log %v", short.Ticks())
	}
	fixed := axisIn(ChartAxisOptions{Min: 0, Max: 1, Ticks: []float64{0, 0.25, 1}, Format: func(v float64) string { return fmt.Sprintf("%.0f%%", v*100) }})
	if fmt.Sprint(fixed.Ticks()) != "[0 0.25 1]" || fixed.Label(0.25) != "25%" {
		t.Fatalf("fixed %v", fixed.Ticks())
	}
}

func TestChartAxisTimeAndCategory(t *testing.T) {
	t0 := time.Date(2026, 10, 5, 9, 7, 0, 0, time.UTC)
	tm := axisIn(ChartAxisOptions{Kind: AxisTime, From: t0, To: t0.Add(2 * time.Hour)})
	if len(tm.Ticks()) == 0 || tm.Label(tm.Ticks()[0]) != "09:30" {
		t.Fatalf("time ticks %v", tm.Ticks())
	}
	cat := axisIn(ChartAxisOptions{Kind: AxisCategory, Categories: []string{"Nave", "Choir", "Crypt", "Bells"}})
	plot := ui.Rect{W: 400, H: 10}
	if cat.Map(plot, 0) != 50 || cat.Map(plot, 3) != 350 || cat.Label(cat.Invert(plot, 260)) != "Crypt" {
		t.Fatalf("category %v %v", cat.Map(plot, 0), cat.Label(cat.Invert(plot, 260)))
	}
}

func TestChartAxisBad(t *testing.T) {
	for name, o := range map[string]ChartAxisOptions{
		"empty":           {},
		"log zero":        {Kind: AxisLog, Min: 0, Max: 10},
		"time":            {Kind: AxisTime},
		"category":        {Kind: AxisCategory},
		"side":            {Min: 0, Max: 1, Side: 9},
		"negative side":   {Min: 0, Max: 1, Side: -1},
		"infinite domain": {Min: 0, Max: math.Inf(1)},
		"outside tick":    {Min: 0, Max: 1, Ticks: []float64{2}},
		"infinite tick":   {Min: 0, Max: 1, Ticks: []float64{math.Inf(1)}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			axisIn(o)
		}()
	}
}

// partsChart composes public parts over three points.
func partsChart(c *ui.Context, notes []ChartNote, crosshair bool) ChartResult {
	pts := [][2]float64{{2, 30}, {5, 80}, {8, 55}}
	x := ChartAxis(c, ChartAxisOptions{Min: 0, Max: 10})
	y := ChartAxis(c, ChartAxisOptions{Side: SideLeft, Min: 0, Max: 100})
	grid := ChartGrid(c, ChartGridOptions{X: &x, Y: &y, Minor: 2})
	hair := ChartCrosshair(c, ChartCrosshairOptions{X: &x, Y: &y})
	notesPart := ChartAnnotation(c, notes, ChartAnnotationOptions{X: &x, Y: &y})
	plotOf := func(r ui.Rect) ui.Rect { return ChartInset(r, 52, 26) }
	ink := core.Tokens(c).AccentText
	at := func(r ui.Rect, i int) (float32, float32) {
		plot := plotOf(r)
		return x.Map(plot, pts[i][0]), y.Map(plot, pts[i][1])
	}
	return Chart(c, ChartFrame{Label: "Parts"}, ChartSpec{
		Marks:   len(pts),
		Readout: func(i int) string { return fmt.Sprintf("x %s · y %s", x.Label(pts[i][0]), y.Label(pts[i][1])) },
		Pick: func(r ui.Rect, px, py float32) int {
			return hair.Snap(len(pts), func(i int) (float32, float32) { return at(r, i) }, px, py)
		},
		Anchor: at,
		Input:  func(plot *ui.Element, r ui.Rect) { notesPart.Drag(plot, plotOf(r)) },
		Draw: func(p *ui.Painter, r ui.Rect, active int) {
			plot := plotOf(r)
			grid.Draw(p, plot)
			x.Draw(p, plot)
			y.Draw(p, plot)
			notesPart.Draw(p, plot)
			for i := range pts {
				px, py := at(r, i)
				internal.Marker(p, 0, px, py, 10, ink)
			}
			if crosshair && active >= 0 {
				px, py := at(r, active)
				hair.Draw(p, plot, px, py)
			}
		},
	})
}

func TestChartAxisSnapshot(t *testing.T) {
	t0 := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	snapshot(t, "chart_axis", 560, 300, func(c *ui.Context) {
		axes := []ChartAxisPart{
			ChartAxis(c, ChartAxisOptions{Min: 0, Max: 97}),
			ChartAxis(c, ChartAxisOptions{Side: SideLeft, Kind: AxisLog, Min: 1, Max: 10000}),
			ChartAxis(c, ChartAxisOptions{Side: SideTop, Kind: AxisTime, From: t0, To: t0.Add(3 * time.Hour)}),
			ChartAxis(c, ChartAxisOptions{Side: SideRight, Kind: AxisCategory, Categories: []string{"Nave", "Choir", "Crypt"}}),
		}
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			plot := ui.Rect{X: r.X + 60, Y: r.Y + 40, W: r.W - 140, H: r.H - 80}
			for _, a := range axes {
				a.Draw(p, plot)
			}
		})
	})
}

func TestChartAxisLargeCategoriesBounded(t *testing.T) {
	categories := make([]string, 100_000)
	for i := range categories {
		categories[i] = fmt.Sprint(i)
	}
	a := axisIn(ChartAxisOptions{Kind: AxisCategory, Categories: categories})
	if n := len(a.Ticks()); n > 65 || n == 0 {
		t.Fatalf("built %d category ticks", n)
	}
	if a.Label(99999) != "99999" || a.Map(ui.Rect{W: 1000}, 99999) < 999 {
		t.Fatal("last category mapping")
	}
}

func charts2MustPanic(t *testing.T, view func(*ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("invalid options accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); view(c) }, 320, 240)
}
