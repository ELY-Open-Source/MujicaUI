package media

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func annotatorTester(notes *[]Annotation, disabled bool, n *int) *ui.Tester {
	src := mediaTestBitmap(400, 200)
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(440).Children(func() {
			if ImageAnnotator(c, src, notes, ImageAnnotatorOptions{Height: 200, Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 520, 360)
}

// annotatorAt maps image pixel (x, y) to window DIPs on the stage.
func annotatorAt(tt *ui.Tester, x, y float32) (float32, float32) {
	r, _ := tt.Find("Annotations")
	scale := min(r.W/400, r.H/200)
	return r.X + (r.W-400*scale)/2 + x*scale, r.Y + (r.H-200*scale)/2 + y*scale
}

func annotatorDrag(tt *ui.Tester, x0, y0, x1, y1 float32) {
	ax, ay := annotatorAt(tt, x0, y0)
	bx, by := annotatorAt(tt, x1, y1)
	tt.Press(ax, ay)
	tt.Move((ax+bx)/2, (ay+by)/2)
	tt.Move(bx, by)
	tt.Release(bx, by)
}

func TestImageAnnotatorDrawShapes(t *testing.T) {
	var notes []Annotation
	n := 0
	tt := annotatorTester(&notes, false, &n)
	if err := tt.Click("Rectangle"); err != nil {
		t.Fatal(err)
	}
	annotatorDrag(tt, 40, 40, 140, 100)
	if err := tt.Click("Arrow"); err != nil {
		t.Fatal(err)
	}
	annotatorDrag(tt, 200, 150, 300, 60)
	if err := tt.Click("Text"); err != nil {
		t.Fatal(err)
	}
	x, y := annotatorAt(tt, 250, 20)
	tt.ClickAt(x, y)
	if len(notes) != 3 || n != 3 {
		t.Fatalf("notes %+v, changes %d", notes, n)
	}
	near := func(p, q image.Point) bool { return absInt(p.X-q.X) <= 2 && absInt(p.Y-q.Y) <= 2 }
	if notes[0].Kind != AnnotateRect || !near(notes[0].From, image.Pt(40, 40)) || !near(notes[0].To, image.Pt(140, 100)) {
		t.Fatalf("rect %+v", notes[0])
	}
	if notes[1].Kind != AnnotateArrow || !near(notes[1].To, image.Pt(300, 60)) {
		t.Fatalf("arrow %+v", notes[1])
	}
	if notes[2].Kind != AnnotateText || notes[2].Text != "Note" {
		t.Fatalf("text %+v", notes[2])
	}
}

func TestImageAnnotatorSelectMoveDelete(t *testing.T) {
	notes := []Annotation{{Kind: AnnotateRect, From: image.Pt(40, 40), To: image.Pt(140, 100)}, {Kind: AnnotateArrow, From: image.Pt(200, 150), To: image.Pt(300, 60)}}
	n := 0
	tt := annotatorTester(&notes, false, &n)
	x, y := annotatorAt(tt, 90, 70)
	tt.ClickAt(x, y)
	if !tt.Focused("Annotations") {
		t.Fatal("stage not focused")
	}
	tt.Key(0, ui.KeyRight)
	if notes[0].From != image.Pt(44, 40) {
		t.Fatalf("Right moved %+v", notes[0])
	}
	tt.Key(0, ui.KeyBracketRight)
	tt.Key(0, ui.KeyDelete)
	if len(notes) != 1 || notes[0].Kind != AnnotateRect {
		t.Fatalf("delete left %+v", notes)
	}
	data, err := ExportAnnotations(notes)
	if err != nil || string(data) != `[{"kind":"rect","x0":44,"y0":40,"x1":144,"y1":100}]` {
		t.Fatalf("export %s %v", data, err)
	}
}

func TestImageAnnotatorDisabledAndBad(t *testing.T) {
	var notes []Annotation
	n := 0
	tt := annotatorTester(&notes, true, &n)
	_ = tt.Click("Rectangle")
	annotatorDrag(tt, 40, 40, 140, 100)
	if len(notes) != 0 || n != 0 {
		t.Fatalf("disabled annotator drew %+v", notes)
	}
	bad := []Annotation{{Kind: AnnotateText}}
	mediaPanics(t, "empty text", func(c *ui.Context) {
		ImageAnnotator(c, mediaTestBitmap(10, 10), &bad, ImageAnnotatorOptions{})
	})
}

func TestImageAnnotatorSnapshot(t *testing.T) {
	src := mediaTestBitmap(400, 200)
	notes := []Annotation{
		{Kind: AnnotateRect, From: image.Pt(30, 40), To: image.Pt(150, 120)},
		{Kind: AnnotateArrow, From: image.Pt(330, 160), To: image.Pt(220, 80)},
		{Kind: AnnotateText, From: image.Pt(230, 30), To: image.Pt(230, 30), Text: "Chandelier"},
	}
	snapshot(t, "image_annotator", 460, 320, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(420).Children(func() {
			ImageAnnotator(c, src, &notes, ImageAnnotatorOptions{Height: 220})
		})
	})
}
