package data

import (
	"image/color"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func cardView(decorated bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			Card(c, CardOptions{Title: "Archive", Decorated: decorated, Footer: func() {
				input.Button(c, "Open", input.ButtonOptions{})
			}}, func() {
				ui.Text(c, "Letters of the court")
			}).Width(260)
		})
	}
}

func TestCardSlots(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		cardView(false)(c)
	}, 400, 300)
	for _, s := range []string{"Archive", "Letters of the court", "Open"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	title, _ := tt.Find("Archive")
	body, _ := tt.Find("Letters of the court")
	foot, _ := tt.Find("Open")
	if !(title.Y < body.Y && body.Y < foot.Y) {
		t.Fatalf("slot order: title %v body %v footer %v", title.Y, body.Y, foot.Y)
	}
}

func TestCardNoTitleNoFooter(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Card(c, CardOptions{}, func() { ui.Text(c, "Only body") })
	}, 300, 200)
	if got := tt.Texts(); len(got) != 1 || got[0] != "Only body" {
		t.Fatalf("texts %q", got)
	}
}

func TestCardDecorationColor(t *testing.T) {
	for _, dec := range []bool{false, true} {
		var bounds ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
				bounds = Card(c, CardOptions{Title: "T", Decorated: dec}, func() {}).Width(200).Bounds()
			})
		}, 300, 200)
		tt.Frame()
		img := tt.Image()
		// The inner corner diamond sits at the inset corner.
		px := img.At(int(bounds.X+cardInset), int(bounds.Y+cardInset))
		want := toRGBA(theme.Light().Ornament)
		if got := color.RGBAModel.Convert(px).(color.RGBA); (got == want) != dec {
			t.Fatalf("decorated=%v: corner pixel %v, ornament %v", dec, got, want)
		}
	}
}

func toRGBA(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, c.A} }

func TestCardSnapshot(t *testing.T) {
	snapshot(t, "card", 640, 240, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(20).AlignItems(ui.Start).Children(func() {
			Card(c, CardOptions{Title: "Archive", Footer: func() {
				input.Button(c, "Open", input.ButtonOptions{Variant: input.Secondary})
			}}, func() { ui.Text(c, "Letters of the court, sealed in wine-red wax.") }).Width(280)
			Card(c, CardOptions{Title: "Chapel", Decorated: true, Footer: func() {
				input.Button(c, "Enter", input.ButtonOptions{})
			}}, func() { ui.Text(c, "Vespers at seven; candles by the north door.") }).Width(280)
		})
	})
	for _, d := range []theme.Density{theme.Compact, theme.Comfortable} {
		snapshot(t, "card-"+strings.ToLower(d.String()), 640, 260, func(c *ui.Context) {
			mode := core.Light
			if core.Tokens(c).Dark {
				mode = core.Dark
			}
			core.Use(c, core.Settings{Mode: mode, Density: d})
			cardView(true)(c)
		})
	}
}
