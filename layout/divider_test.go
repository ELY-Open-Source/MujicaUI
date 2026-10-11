package layout

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestDividerSizes(t *testing.T) {
	var h, v ui.Rect
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(300).Children(func() {
			h = Divider(c, DividerOptions{}).Bounds()
			ui.Row(c).Height(80).Children(func() {
				v = Divider(c, DividerOptions{Vertical: true}).Bounds()
			})
		})
	}, 400, 300).Frame()
	if h.W != 300 || h.H != 1 {
		t.Fatalf("horizontal %vx%v, want 300x1", h.W, h.H)
	}
	if v.W != 1 || v.H != 80 {
		t.Fatalf("vertical %vx%v, want 1x80", v.W, v.H)
	}
}

func TestDividerNeverFocused(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			Divider(c, DividerOptions{Label: "Section", Ornament: true})
			Divider(c, DividerOptions{Ornament: true})
			input.Button(c, "OK", input.ButtonOptions{})
		})
	}, 400, 300)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("OK") {
		t.Fatal("first Tab did not skip the dividers to the button")
	}
	if tt.Focused("Section") {
		t.Fatal("divider label took the focus")
	}
}

func TestDividerLabelCentered(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(400).Children(func() {
			Divider(c, DividerOptions{Label: "Chapter II"})
		})
	}, 400, 100)
	r, ok := tt.Find("Chapter II")
	if !ok {
		t.Fatal("label missing")
	}
	if mid := r.X + r.W/2; mid < 199 || mid > 201 {
		t.Fatalf("label centered at %v, want 200", mid)
	}
}

func TestDividerOrnamentColor(t *testing.T) {
	var want ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		want = core.Tokens(c).Ornament
		ui.Column(c).Width(201).Padding(10, 0).Children(func() {
			Divider(c, DividerOptions{Ornament: true})
		})
	}, 201, 40)
	// The diamond's center sits on the line's middle.
	got := tt.Image().At(100, 13)
	if got != (color.RGBA{want.R, want.G, want.B, 255}) {
		t.Fatalf("diamond center %v, want ornament %v", got, want)
	}
}

func TestDividerVerticalLabelPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("vertical label did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Divider(c, DividerOptions{Vertical: true, Label: "x"})
	}, 100, 100)
}

func TestDividerSnapshot(t *testing.T) {
	snapshot(t, "divider", 420, 170, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(20).Children(func() {
			Divider(c, DividerOptions{})
			Divider(c, DividerOptions{Ornament: true})
			Divider(c, DividerOptions{Label: "Chapter II"})
			Divider(c, DividerOptions{Label: "第二章", Ornament: true})
			ui.Row(c).Height(28).Gap(12).Children(func() {
				ui.Text(c, "Left")
				Divider(c, DividerOptions{Vertical: true})
				ui.Text(c, "Middle")
				Divider(c, DividerOptions{Vertical: true, Ornament: true})
				ui.Text(c, "Right")
			})
		})
	})
}

func TestDividerLabelTextScale(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(300).Children(func() { Divider(c, DividerOptions{Label: "Chapter II"}) })
	}, 400, 200)
	plain, _ := tt.Find("Chapter II")
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	if big, _ := tt.Find("Chapter II"); big.H < 1.8*plain.H {
		t.Fatalf("label height %v at 200%%, was %v", big.H, plain.H)
	}
}
