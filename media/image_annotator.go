package media

import (
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AnnotationKind is the shape of an Annotation.
type AnnotationKind int

const (
	AnnotateRect AnnotationKind = iota
	AnnotateArrow
	AnnotateText
)

// Annotation is a mark in image pixels; text at From.
type Annotation struct {
	Kind     AnnotationKind
	From, To image.Point
	Text     string
}

// annotationJSON is the exported form of an Annotation.
type annotationJSON struct {
	Kind string `json:"kind"`
	X0   int    `json:"x0"`
	Y0   int    `json:"y0"`
	X1   int    `json:"x1"`
	Y1   int    `json:"y1"`
	Text string `json:"text,omitempty"`
}

// ExportAnnotations writes notes as JSON: kind, corners and text.
func ExportAnnotations(notes []Annotation) ([]byte, error) {
	out := make([]annotationJSON, len(notes))
	for i, n := range notes {
		out[i] = annotationJSON{[]string{"rect", "arrow", "text"}[n.Kind], n.From.X, n.From.Y, n.To.X, n.To.Y, n.Text}
	}
	return json.Marshal(out)
}

// ImageAnnotatorOptions configure an ImageAnnotator.
type ImageAnnotatorOptions struct {
	// Height is the stage's height; zero means 320.
	Height   float32
	Disabled bool
	Label    string
}

type annotatorState struct {
	tool       int // 0 select, 1 rect, 2 arrow, 3 text
	selected   int
	draft      string
	pressing   bool
	start, now image.Point
	moved      image.Point
}

type annotatorKey struct{}

// ImageAnnotator draws, selects, moves and deletes marks over src.
func ImageAnnotator(c *ui.Context, src ui.ImageSource, notes *[]Annotation, opts ImageAnnotatorOptions) core.ChoiceResult {
	iw, ih := data.ViewerSize(src)
	size := image.Pt(int(iw), int(ih))
	for _, n := range *notes {
		if n.Kind < AnnotateRect || n.Kind > AnnotateText || (n.Kind == AnnotateText && n.Text == "") {
			panic(fmt.Sprintf("mujica: ImageAnnotator bad annotation %+v", n))
		}
	}
	if opts.Height < 0 {
		panic("mujica: ImageAnnotator Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	before := slices.Clone(*notes)
	col := ui.Column(c).Gap(u * 2).Disabled(opts.Disabled)
	st := ui.Local(col, annotatorKey{}, func() *annotatorState {
		return &annotatorState{selected: -1, draft: core.Msg(c, "image_annotator.note")}
	})
	s := *st
	if s.selected >= len(*notes) {
		s.selected = -1
	}
	col.Children(func() {
		annotatorToolbar(c, k, s, notes)
		stage := ui.Box(c).FillWidth().Height(h).Focusable().Clip().Radius(theme.ControlRadius).
			Background(VideoBlack).Label(annotatorLabel(c, opts.Label)).Cursor(ui.CursorCrosshair)
		internal.Ring(stage, k.Focus, theme.ControlRadius)
		r := stage.Bounds()
		scale, ox, oy := cropperPlace(r, size)
		if !stage.IsDisabled() && scale > 0 {
			annotatorPointer(stage, s, notes, r, size, scale, ox, oy)
			annotatorKeys(stage, s, notes, size)
		}
		if s.selected >= 0 {
			n := (*notes)[s.selected]
			stage.Description(fmt.Sprintf(core.Msg(c, "image_annotator.selected"), s.selected+1, len(*notes), annotatorKind(c, n.Kind)))
		}
		list, sel, preview := slices.Clone(*notes), s.selected, annotatorPreview(s)
		size := core.FontSize(c, theme.BodySize)
		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			scale, ox, oy := cropperPlace(r, image.Pt(int(iw), int(ih)))
			p.Image(src, ui.Rect{X: ox, Y: oy, W: iw * scale, H: ih * scale}, ui.FillBox)
			at := func(pt image.Point) (float32, float32) {
				return internal.Snap(p, ox+float32(pt.X)*scale), internal.Snap(p, oy+float32(pt.Y)*scale)
			}
			for i, n := range append(list, preview...) {
				annotatorPaint(p, k, n, at, size, i == sel)
			}
		})
	})
	changed := !slices.Equal(before, *notes)
	if changed {
		slog.Debug("mujica: image_annotator changed", "count", len(*notes))
	}
	return core.NewChoiceResult(col, changed)
}

func annotatorLabel(c *ui.Context, l string) string {
	if l != "" {
		return l
	}
	return core.Msg(c, "image_annotator.label")
}

func annotatorKind(c *ui.Context, k AnnotationKind) string {
	return core.Msg(c, []string{"image_annotator.rect", "image_annotator.arrow", "image_annotator.text"}[k])
}

// annotatorToolbar picks tools, edits text and deletes the selection.
func annotatorToolbar(c *ui.Context, k theme.Tokens, s *annotatorState, notes *[]Annotation) {
	u := core.Density(c).Unit()
	ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Wrap().Children(func() {
		navigation.SegmentedControl(c, &s.tool, []navigation.Segment{
			{Label: core.Msg(c, "image_annotator.select"), Icon: core.MediaPointer, HideLabel: true},
			{Label: core.Msg(c, "image_annotator.rect"), Icon: core.MediaRect, HideLabel: true},
			{Label: core.Msg(c, "image_annotator.arrow"), Icon: core.MediaArrow, HideLabel: true},
			{Label: core.Msg(c, "image_annotator.text"), Icon: core.MediaText, HideLabel: true},
		}, navigation.SegmentedControlOptions{Label: core.Msg(c, "image_annotator.tool")})
		editing := s.selected >= 0 && (*notes)[s.selected].Kind == AnnotateText
		if editing {
			s.draft = (*notes)[s.selected].Text
		}
		if s.tool == 3 || editing {
			in := input.TextInput(c, &s.draft, input.TextInputOptions{Label: core.Msg(c, "image_annotator.textLabel")})
			in.Element.Width(180)
			if editing && in.Changed() && s.draft != "" {
				(*notes)[s.selected].Text = s.draft
			}
		}
		del := core.IconAction(c, icons.Must("trash-2"), core.Msg(c, "image_annotator.delete")).Disabled(s.selected < 0)
		if del.Clicked() && !del.IsDisabled() {
			*notes = slices.Delete(*notes, s.selected, s.selected+1)
			s.selected = -1
		}
		ui.Textf(c, core.Msg(c, "image_annotator.count"), len(*notes)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
}

// annotatorPreview is the shape being drawn, if any.
func annotatorPreview(s *annotatorState) []Annotation {
	if !s.pressing || s.tool == 0 || s.tool == 3 {
		return nil
	}
	return []Annotation{{Kind: AnnotationKind(s.tool - 1), From: s.start, To: s.now}}
}

// annotatorPointer draws shapes, selects by click, drags the selection.
func annotatorPointer(stage *ui.Element, s *annotatorState, notes *[]Annotation, r ui.Rect, size image.Point, scale, ox, oy float32) {
	_, _, _ = stage.Dragged()
	px, py, _ := stage.PointerPosition()
	pt := image.Pt(int(math.Round(float64((px+r.X-ox)/scale))), int(math.Round(float64((py+r.Y-oy)/scale))))
	pt = image.Pt(max(0, min(size.X, pt.X)), max(0, min(size.Y, pt.Y)))
	switch {
	case stage.Pressed() && !s.pressing:
		s.pressing, s.start, s.now, s.moved = true, pt, pt, image.Point{}
		if s.tool == 0 {
			s.selected = annotatorHit(*notes, pt, int(8/scale)+1)
		}
	case stage.Pressed():
		if s.tool == 0 && s.selected >= 0 {
			d := pt.Sub(s.now)
			n := &(*notes)[s.selected]
			n.From, n.To = n.From.Add(d), n.To.Add(d)
		}
		s.now = pt
	case s.pressing:
		s.pressing = false
		big := absInt(s.now.X-s.start.X) >= 4 || absInt(s.now.Y-s.start.Y) >= 4
		switch {
		case s.tool == 3 && s.draft != "":
			*notes = append(*notes, Annotation{Kind: AnnotateText, From: s.start, To: s.start, Text: s.draft})
			s.selected = len(*notes) - 1
		case (s.tool == 1 || s.tool == 2) && big:
			*notes = append(*notes, Annotation{Kind: AnnotationKind(s.tool - 1), From: s.start, To: s.now})
			s.selected = len(*notes) - 1
		}
	}
}

// annotatorKeys: [ ] select, arrows move, Delete removes, Escape clears.
func annotatorKeys(stage *ui.Element, s *annotatorState, notes *[]Annotation, size image.Point) {
	n := len(*notes)
	step := max(1, size.X/100)
	move := func(dx, dy int) {
		if s.selected >= 0 {
			a := &(*notes)[s.selected]
			a.From, a.To = a.From.Add(image.Pt(dx, dy)), a.To.Add(image.Pt(dx, dy))
		}
	}
	switch {
	case n > 0 && stage.Shortcut(0, ui.KeyBracketRight):
		s.selected = (s.selected + 1) % n
	case n > 0 && stage.Shortcut(0, ui.KeyBracketLeft):
		s.selected = (s.selected - 1 + n) % n
	case s.selected >= 0 && (stage.Shortcut(0, ui.KeyDelete) || stage.Shortcut(0, ui.KeyBackspace)):
		*notes = slices.Delete(*notes, s.selected, s.selected+1)
		s.selected = -1
	case stage.Shortcut(0, ui.KeyEscape):
		s.selected = -1
	case stage.Shortcut(0, ui.KeyLeft):
		move(-step, 0)
	case stage.Shortcut(0, ui.KeyRight):
		move(step, 0)
	case stage.Shortcut(0, ui.KeyUp):
		move(0, -step)
	case stage.Shortcut(0, ui.KeyDown):
		move(0, step)
	}
}

// annotatorHit is the topmost mark near pt, or -1.
func annotatorHit(notes []Annotation, pt image.Point, slop int) int {
	for i := len(notes) - 1; i >= 0; i-- {
		n := notes[i]
		b := image.Rectangle{Min: n.From, Max: n.To}.Canon()
		if n.Kind == AnnotateText {
			b = image.Rectangle{Min: n.From, Max: n.From.Add(image.Pt(len(n.Text)*8, 18))}
		}
		if pt.In(b.Inset(-slop)) {
			return i
		}
	}
	return -1
}

func absInt(v int) int { return max(v, -v) }

// annotatorPaint draws a mark with a halo for any picture.
func annotatorPaint(p *ui.Painter, k theme.Tokens, n Annotation, at func(image.Point) (float32, float32), size float32, selected bool) {
	x0, y0 := at(n.From)
	x1, y1 := at(n.To)
	halo, ink := ui.Hex("#FBF7EF"), ui.Hex("#A12635")
	switch n.Kind {
	case AnnotateRect:
		b := ui.Rect{X: min(x0, x1), Y: min(y0, y1), W: core.Abs32(x1 - x0), H: core.Abs32(y1 - y0)}
		p.Stroke(b, halo, 0, 4)
		p.Stroke(b, ink, 0, 2)
	case AnnotateArrow:
		var path ui.Path
		path.MoveTo(x0, y0).LineTo(x1, y1)
		a := math.Atan2(float64(y1-y0), float64(x1-x0))
		for _, d := range []float64{2.6, -2.6} {
			path.MoveTo(x1, y1).LineTo(x1+12*float32(math.Cos(a+d)), y1+12*float32(math.Sin(a+d)))
		}
		p.StrokePath(&path, 4, halo)
		p.StrokePath(&path, 2, ink)
	case AnnotateText:
		w, h := p.MeasureText(1000, ui.Span{Text: n.Text, Size: size})
		p.Fill(ui.Rect{X: x0, Y: y0, W: w + 8, H: h + 4}, halo, theme.ControlRadius)
		p.Text(x0+4, y0+2, n.Text, size, ink)
		x1, y1 = x0+w+8, y0+h+4
	}
	if selected {
		b := ui.Rect{X: min(x0, x1) - 4, Y: min(y0, y1) - 4, W: core.Abs32(x1-x0) + 8, H: core.Abs32(y1-y0) + 8}
		p.StrokeDashed(b, k.Focus, 0, 1)
		p.FillPath(internal.Diamond(internal.Snap(p, b.X), internal.Snap(p, b.Y), 8), k.Focus)
	}
}
