package finance

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func bullSpread() []OptionLeg {
	return []OptionLeg{{Kind: LegCall, Strike: 100, Premium: 3, Quantity: 1}, {Kind: LegCall, Strike: 110, Premium: 1, Quantity: -1}}
}

func TestPayoffMath(t *testing.T) {
	legs := bullSpread()
	for _, tc := range []struct{ price, want float64 }{{90, -2}, {100, -2}, {102, 0}, {110, 8}, {130, 8}} {
		if got := PayoffAt(legs, tc.price); got != tc.want {
			t.Errorf("payoff at %v = %v, want %v", tc.price, got, tc.want)
		}
	}
	if got := Breakevens(legs, 70, 140); !slices.Equal(got, []float64{102}) {
		t.Fatalf("breakevens %v", got)
	}
	straddle := []OptionLeg{{Kind: LegCall, Strike: 100, Premium: 4, Quantity: 1}, {Kind: LegPut, Strike: 100, Premium: 3, Quantity: 1}}
	if got := Breakevens(straddle, 50, 150); !slices.Equal(got, []float64{93, 107}) {
		t.Fatalf("straddle breakevens %v", got)
	}
	covered := []OptionLeg{{Kind: LegStock, Premium: 50, Quantity: 100}, {Kind: LegCall, Strike: 55, Premium: 2, Quantity: -100}}
	if got := PayoffAt(covered, 60); got != 700 {
		t.Fatalf("covered call at 60 = %v", got)
	}
}

func TestPayoffDiagramKeysAndText(t *testing.T) {
	var last chart.ChartResult
	tt := chatTester(640, 380, func(c *ui.Context) {
		last = PayoffDiagram(c, bullSpread(), PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Spread"}, Spot: 103, Low: 90, High: 120, Points: 31})
	})
	if !tt.HasText("Breakeven: 102") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if last.Point != 30 || !tt.HasText("Price 120 · profit 8") {
		t.Fatalf("End: point %d texts %q", last.Point, tt.Texts())
	}
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("Price 90 · loss -2") {
		t.Fatalf("Home: %q", tt.Texts())
	}
	ann := strings.Join(tt.Announcements(), "|")
	if !strings.Contains(ann, "Price 90") {
		t.Fatalf("announcements %q", ann)
	}
}

func TestPayoffDiagramEmptyAndBadRange(t *testing.T) {
	tt := chatTester(640, 380, func(c *ui.Context) {
		PayoffDiagram(c, nil, PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Empty"}})
	})
	if !tt.HasText("No data") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PayoffDiagram(c, bullSpread(), PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "x"}, Low: 10, High: 5})
	}, 100, 100, 1)
}

func TestPayoffDiagramSnapshot(t *testing.T) {
	snapshot(t, "payoff_diagram", 640, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			PayoffDiagram(c, bullSpread(), PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Bull call spread"}, Spot: 103})
		})
	})
}

func TestPayoffDiagramSparseGridKeepsPeak(t *testing.T) {
	legs := []OptionLeg{{Kind: LegCall, Strike: 90, Premium: 1, Quantity: 1}, {Kind: LegCall, Strike: 100, Quantity: -2}, {Kind: LegCall, Strike: 110, Quantity: 1}}
	tt := chatTester(640, 380, func(c *ui.Context) {
		PayoffDiagram(c, legs, PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Butterfly"}, Low: 80, High: 120, Points: 2})
	})
	r, ok := tt.Find("Butterfly")
	if !ok {
		t.Fatal("plot missing")
	}
	x := int(r.X + 56 + (r.W-68)/2)
	found := false
	ink := theme.Light().AccentText
	for y := int(r.Y + 8); y < int(r.Y+r.H/2); y++ {
		for dx := -2; dx <= 2; dx++ {
			pixel := tt.Image().RGBAAt(x+dx, y)
			if pixel.R == ink.R && pixel.G == ink.G && pixel.B == ink.B {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("profit peak falls outside the plot")
	}
}
