package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSignaturePath(t *testing.T) {
	got := SignaturePath([][]SignaturePoint{{{1, 2}, {3.5, 4}}, {{10, 10}}})
	if want := "M 1 2 L 3.5 4 M 10 10"; got != want {
		t.Fatalf("path %q, want %q", got, want)
	}
	if SignaturePath(nil) != "" {
		t.Fatal("empty strokes made a path")
	}
}

func TestSignaturePadStrokesAndClear(t *testing.T) {
	var strokes [][]SignaturePoint
	dis := false
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(400).Children(func() {
			if SignaturePad(c, &strokes, SignaturePadOptions{Disabled: dis}).Changed() {
				n++
			}
		})
	}, 500, 300)
	tt.Press(50, 60)
	tt.Move(80, 70)
	tt.Move(80.5, 70) // under a DIP: not kept
	tt.Move(120, 90)
	tt.Release(120, 90)
	if len(strokes) != 1 || len(strokes[0]) != 3 || strokes[0][0] != (SignaturePoint{30, 40}) || n != 1 {
		t.Fatalf("stroke %+v, changes %d", strokes, n)
	}
	tt.Press(200, 60)
	tt.Release(200, 60)
	if len(strokes) != 2 || n != 2 {
		t.Fatalf("dot: %+v", strokes)
	}
	tt.Key(0, ui.KeyDelete) // drawing focused the pad
	if len(strokes) != 0 || n != 3 {
		t.Fatalf("Delete kept %d strokes", len(strokes))
	}
	tt.Press(50, 60)
	tt.Release(50, 60)
	if err := tt.Click("Clear"); err != nil || len(strokes) != 0 {
		t.Fatalf("Clear button: %v %d", err, len(strokes))
	}
	dis = true
	tt.Frame()
	tt.Press(50, 60)
	tt.Move(90, 60)
	tt.Release(90, 60)
	if len(strokes) != 0 {
		t.Fatal("disabled pad drew")
	}
}

func TestSignaturePadSnapshot(t *testing.T) {
	strokes := [][]SignaturePoint{{{30, 90}, {50, 60}, {60, 95}, {80, 55}, {95, 100}, {120, 70}, {150, 85}}, {{170, 60}, {175, 100}}}
	var empty [][]SignaturePoint
	snapshot(t, "signature_pad", 560, 260, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ui.Column(c).Width(250).Children(func() { SignaturePad(c, &strokes, SignaturePadOptions{Height: 180}) })
			ui.Column(c).Width(250).Children(func() { SignaturePad(c, &empty, SignaturePadOptions{Height: 180}) })
		})
	})
}
