package chart

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChartNoteKind is what a ChartNote draws.
type ChartNoteKind int

const (
	// NoteHLine is a horizontal line at Y value Value.
	NoteHLine ChartNoteKind = iota
	// NoteVLine is a vertical line at X value Value.
	NoteVLine
	// NoteBand shades X values Value to End.
	NoteBand
	// NoteText writes Text at X Value, Y End.
	NoteText
)

// ChartNote is one annotation; Draggable lets the pointer move it.
type ChartNote struct {
	Kind       ChartNoteKind
	Value, End float64
	Text       string
	Draggable  bool
}

// ChartAnnotationOptions configure a ChartAnnotation; both axes are required.
type ChartAnnotationOptions struct {
	X, Y *ChartAxisPart
}

// ChartAnnotationPart is notes ready to draw, and to drag.
type ChartAnnotationPart struct {
	c     *ui.Context
	notes []ChartNote
	opts  ChartAnnotationOptions
	k     theme.Tokens
	size  float32
}

// ChartAnnotation draws notes and enables optional movement.
func ChartAnnotation(c *ui.Context, notes []ChartNote, opts ChartAnnotationOptions) ChartAnnotationPart {
	if opts.X == nil || opts.Y == nil {
		panic("mujica: ChartAnnotation needs X and Y axes")
	}
	for _, n := range notes {
		if n.Kind < NoteHLine || n.Kind > NoteText || math.IsNaN(n.Value) || math.IsInf(n.Value, 0) || math.IsNaN(n.End) || math.IsInf(n.End, 0) || n.Kind == NoteBand && n.End <= n.Value || n.Kind == NoteText && n.Text == "" {
			panic("mujica: ChartAnnotation note needs a known kind, a band End above Value, and text")
		}
	}
	return ChartAnnotationPart{c: c, notes: notes, opts: opts, k: core.Tokens(c), size: core.FontSize(c, theme.CaptionSize)}
}

// Draw paints bands first, then lines and text, inside plot.
func (a ChartAnnotationPart) Draw(p *ui.Painter, plot ui.Rect) {
	k, x, y := a.k, a.opts.X, a.opts.Y
	for _, n := range a.notes {
		if n.Kind == NoteBand {
			x0, x1 := x.Map(plot, n.Value), x.Map(plot, n.End)
			p.Fill(ui.Rect{X: x0, Y: plot.Y, W: x1 - x0, H: plot.H}, k.Info.Alpha(0.12), 0)
			ChartText(p, n.Text, x0+4, plot.Y+4, a.size, k.Text, 0, 0)
		}
	}
	for _, n := range a.notes {
		switch n.Kind {
		case NoteHLine:
			at := internal.Snap(p, y.Map(plot, n.Value))
			chartDash(p, plot.X, at, plot.X+plot.W, at, k.Text)
			ChartText(p, chartNoteLabel(n, y), plot.X+plot.W-4, at-3, a.size, k.Text, 1, 1)
			if n.Draggable {
				p.FillPath(internal.Diamond(internal.Snap(p, plot.X+8), at, 11), k.AccentText)
			}
		case NoteVLine:
			at := internal.Snap(p, x.Map(plot, n.Value))
			chartDash(p, at, plot.Y, at, plot.Y+plot.H, k.Text)
			ChartText(p, chartNoteLabel(n, x), at+4, plot.Y+4, a.size, k.Text, 0, 0)
			if n.Draggable {
				p.FillPath(internal.Diamond(at, internal.Snap(p, plot.Y+plot.H-8), 11), k.AccentText)
			}
		case NoteText:
			px, py := x.Map(plot, n.Value), y.Map(plot, n.End)
			p.FillPath(internal.Diamond(internal.Snap(p, px), internal.Snap(p, py), 8), k.Text)
			ChartText(p, n.Text, px+8, py, a.size, k.Text, 0, 0.5)
		}
	}
}

// chartDash draws a dashed horizontal or vertical line.
func chartDash(p *ui.Painter, x0, y0, x1, y1 float32, col ui.Color) {
	const on, off = 6, 4
	if y0 == y1 {
		for x := x0; x < x1; x += on + off {
			p.Line(x, y0, min(x+on, x1), y0, 1.5, col)
		}
		return
	}
	for y := y0; y < y1; y += on + off {
		p.Line(x0, y, x0, min(y+on, y1), 1.5, col)
	}
}

// chartNoteLabel formats a line's label and value.
func chartNoteLabel(n ChartNote, a *ChartAxisPart) string {
	if n.Text != "" {
		return n.Text + " " + a.Label(n.Value)
	}
	return a.Label(n.Value)
}

type noteDrag struct {
	on   bool
	note int
}

// Drag handles movement; plot uses element coordinates.
func (a ChartAnnotationPart) Drag(plotElement *ui.Element, plot ui.Rect) bool {
	if plotElement.IsDisabled() {
		return false
	}
	if a.keyboard(plotElement) {
		return true
	}
	st := ui.Local(plotElement, "mujica.notes.drag", func() noteDrag { return noteDrag{} })
	dx, dy, held := plotElement.Dragged()
	if !held || plotElement.IsDisabled() {
		*st = noteDrag{}
		return false
	}
	px, py, _ := plotElement.PointerPosition()
	x, y := a.opts.X, a.opts.Y
	if !st.on {
		*st = noteDrag{on: true, note: -1}
		for i, n := range a.notes {
			if !n.Draggable {
				continue
			}
			var d float32 = 1e9
			switch n.Kind {
			case NoteHLine:
				d = core.Abs32(y.Map(plot, n.Value) - (py - dy))
			case NoteVLine:
				d = core.Abs32(x.Map(plot, n.Value) - (px - dx))
			case NoteBand:
				if x0, x1 := x.Map(plot, n.Value), x.Map(plot, n.End); px-dx >= x0 && px-dx <= x1 {
					d = 0
				}
			case NoteText:
				d = max(core.Abs32(x.Map(plot, n.Value)-(px-dx)), core.Abs32(y.Map(plot, n.End)-(py-dy)))
			}
			if d <= 6 {
				st.note = i
			}
		}
	}
	if st.note < 0 || dx == 0 && dy == 0 {
		return false
	}
	n := &a.notes[st.note]
	shift := func(ax *ChartAxisPart, v float64, at, d float32) float64 {
		return v + ax.Invert(plot, at) - ax.Invert(plot, at-d)
	}
	switch n.Kind {
	case NoteHLine:
		n.Value = shift(y, n.Value, py, dy)
	case NoteVLine:
		n.Value = shift(x, n.Value, px, dx)
	case NoteBand:
		w := n.End - n.Value
		n.Value = shift(x, n.Value, px, dx)
		n.End = n.Value + w
	case NoteText:
		n.Value, n.End = shift(x, n.Value, px, dx), shift(y, n.End, py, dy)
	}
	slog.Debug("mujica: chart_annotation drag", "note", st.note, "value", n.Value)
	return true
}

func (a ChartAnnotationPart) keyboard(e *ui.Element) bool {
	index := ui.Local(e, "mujica.notes.key", func() int { return 0 })
	dir := 0
	if e.Shortcut(ui.Alt, ui.KeyPageDown) {
		dir = 1
	}
	if e.Shortcut(ui.Alt, ui.KeyPageUp) {
		dir = -1
	}
	var editable []int
	for i, n := range a.notes {
		if n.Draggable {
			editable = append(editable, i)
		}
	}
	if len(editable) == 0 {
		return false
	}
	*index = (*index + dir + len(editable)) % len(editable)
	n := &a.notes[editable[*index]]
	dx, dy := 0.0, 0.0
	if e.Shortcut(ui.Alt, ui.KeyRight) {
		dx = 1
	}
	if e.Shortcut(ui.Alt, ui.KeyLeft) {
		dx = -1
	}
	if e.Shortcut(ui.Alt, ui.KeyUp) {
		dy = 1
	}
	if e.Shortcut(ui.Alt, ui.KeyDown) {
		dy = -1
	}
	xstep := (a.opts.X.hi - a.opts.X.lo) / 100
	ystep := (a.opts.Y.hi - a.opts.Y.lo) / 100
	before := *n
	switch n.Kind {
	case NoteHLine:
		n.Value += dy * ystep
	case NoteVLine:
		n.Value += dx * xstep
	case NoteBand:
		n.Value, n.End = n.Value+dx*xstep, n.End+dx*xstep
	case NoteText:
		n.Value, n.End = n.Value+dx*xstep, n.End+dy*ystep
	}
	if dir != 0 || before != *n {
		key := "chart_annotation.x"
		value := a.opts.X.Label(n.Value)
		switch n.Kind {
		case NoteHLine:
			key, value = "chart_annotation.y", a.opts.Y.Label(n.Value)
		case NoteBand:
			value += "–" + a.opts.X.Label(n.End)
		case NoteText:
			value += " · " + fmt.Sprintf(core.Msg(a.c, "chart_annotation.y"), a.opts.Y.Label(n.End))
		}
		text := fmt.Sprintf(core.Msg(a.c, "chart_annotation.readout"), editable[*index]+1, fmt.Sprintf(core.Msg(a.c, key), value))
		e.Description(text)
		a.c.Announce(text)
		slog.Debug("mujica: chart_annotation keyboard", "note", editable[*index], "value", n.Value)
	}
	return before != *n
}
