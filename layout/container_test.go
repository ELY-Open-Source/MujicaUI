package layout

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestContainerPaddingByDensity(t *testing.T) {
	for _, tc := range []struct {
		density theme.Density
		pad     Space
		want    float32
	}{{theme.Default, SpaceDefault, 16}, {theme.Compact, SpaceDefault, 12}, {theme.Comfortable, SpaceS, 10}, {theme.Default, SpaceNone, 0}} {
		var outer, inner ui.Rect
		ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Density: tc.density})
			outer = Container(c, ContainerOptions{Padding: tc.pad}, func() {
				inner = ui.Box(c).Size(10, 10).Bounds()
			}).Bounds()
		}, 400, 300).Frame()
		if inner.X-outer.X != tc.want || inner.Y-outer.Y != tc.want {
			t.Errorf("%v %v: inset %v,%v, want %v", tc.density, tc.pad, inner.X-outer.X, inner.Y-outer.Y, tc.want)
		}
	}
}

func TestContainerMaxWidthCentered(t *testing.T) {
	var r ui.Rect
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(600).Children(func() {
			r = Container(c, ContainerOptions{MaxWidth: 300, Center: true}, func() { ui.Text(c, "Body") }).Bounds()
		})
	}, 600, 300).Frame()
	if r.X != 150 || r.W != 300 {
		t.Fatalf("container at x=%v w=%v, want 150 and 300", r.X, r.W)
	}
}

func TestContainerFill(t *testing.T) {
	var r ui.Rect
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Size(320, 200).Children(func() {
			r = Container(c, ContainerOptions{Fill: true}, func() { ui.Text(c, "x") }).Bounds()
		})
	}, 400, 300).Frame()
	if r.W != 320 || r.H != 200 {
		t.Fatalf("filled %vx%v, want 320x200", r.W, r.H)
	}
}

func TestContainerClip(t *testing.T) {
	for _, clip := range []bool{false, true} {
		var bg ui.Color
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			bg = core.Tokens(c).Background
			Container(c, ContainerOptions{Padding: SpaceNone, Clip: clip}, func() {
				ui.Box(c).Size(200, 200).Background(ui.RGB(255, 0, 0))
			}).Size(100, 100)
		}, 300, 300)
		got := tt.Image().At(150, 150)
		clipped := got == color.RGBA{bg.R, bg.G, bg.B, 255}
		if clipped != clip {
			t.Errorf("Clip %v: pixel outside %v", clip, got)
		}
	}
}

func TestContainerInvalidSpacePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("invalid Space did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Container(c, ContainerOptions{Padding: Space(99)}, func() {})
	}, 100, 100)
}

func TestContainerSnapshot(t *testing.T) {
	snapshot(t, "container", 420, 150, func(c *ui.Context) {
		k := core.Tokens(c)
		ui.Column(c).Padding(16).Children(func() {
			Container(c, ContainerOptions{MaxWidth: 280, Center: true}, func() {
				ui.Text(c, "Of the Court").Font(theme.SerifFont).FontSize(theme.H3Size)
				ui.Text(c, "A centered container, 280 DIPs at most, padded by four units.").TextColor(k.TextMuted)
			}).Background(k.Surface).Border(1, k.Border).Radius(theme.CardRadius)
		})
	})
}
