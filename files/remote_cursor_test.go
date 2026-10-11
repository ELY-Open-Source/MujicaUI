package files

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// remoteMono maps offsets onto 40-column, 8×20 DIP cells.
func remoteMono() RemoteCursorOptions {
	at := func(o int) (float32, float32) { return float32(o%40) * 8, float32(o/40) * 20 }
	return RemoteCursorOptions{
		Caret: func(o int) ui.Rect { x, y := at(o); return ui.Rect{X: x, Y: y, W: 2, H: 20} },
		Selection: func(from, to int) []ui.Rect {
			var out []ui.Rect
			for o := from; o < to; o++ {
				x, y := at(o)
				out = append(out, ui.Rect{X: x, Y: y, W: 8, H: 20})
			}
			return out
		},
	}
}

func remoteTester(carets []RemoteCaret, button *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Box(c).Size(320, 120).Children(func() {
			if input.Button(c, "Underneath", input.ButtonOptions{}).Size(320, 120).Clicked() {
				*button++
			}
			RemoteCursor(c, carets, remoteMono())
		})
	}, 340, 140)
}

func TestRemoteCursorDrawCaretAndSelection(t *testing.T) {
	carets := []RemoteCaret{{ID: "ada", Name: "Ada", Offset: 45, Anchor: 42}}
	tt := remoteTester(carets, new(int))
	col := remoteColor(theme.Light(), "ada")
	img := tt.Image()
	// Offset 45 places the caret at x=40, y=20..40.
	if got := img.RGBAAt(40, 35); got != (color.RGBA{col.R, col.G, col.B, 255}) {
		t.Fatalf("caret pixel %v, want %v", got, col)
	}
	if got := img.RGBAAt(28, 35); got == (color.RGBA{col.R, col.G, col.B, 255}) || got == img.RGBAAt(100, 35) {
		t.Fatalf("selection pixel %v not tinted", got)
	}
	if !tt.HasText("Also editing: Ada") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestRemoteCursorPassInput(t *testing.T) {
	n := 0
	tt := remoteTester([]RemoteCaret{{ID: "b", Name: "Blaise", Offset: 3, Anchor: 3}}, &n)
	tt.ClickAt(30, 10)
	if n != 1 {
		t.Fatal("the overlay took the click")
	}
}

func TestRemoteCursorInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("missing ID did not panic")
		}
	}()
	remoteTester([]RemoteCaret{{Name: "x"}}, new(int))
}

func TestRemoteCursorSnapshot(t *testing.T) {
	snapshot(t, "remote_cursor", 360, 120, func(c *ui.Context) {
		k := core.Tokens(c)
		ui.Box(c).Margin(12).Size(320, 80).Background(k.Surface).Children(func() {
			ui.Column(c).Children(func() {
				for _, l := range []string{"The seal of the court is kept", "by the chancellor, who answers"} {
					ui.Text(c, l).Font(theme.MonoFont).FontSize(core.FontSize(c, 13)).Height(20)
				}
			})
			RemoteCursor(c, []RemoteCaret{{ID: "ada", Name: "Ada", Offset: 48, Anchor: 43}, {ID: "blaise", Name: "Blaise", Offset: 8, Anchor: 8}}, remoteMono())
		})
	})
}
