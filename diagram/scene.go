package diagram

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// scene is a diagram ready to draw.
type scene struct {
	w, h  float32
	title string
	runs  []run
}

// run is ink under one transform.
type run struct {
	t     ui.Transform
	marks []mark
}

// mark is a path, or a line of text.
type mark struct {
	path         *ui.Path
	fill, stroke ui.Color
	width        float32
	text         *ui.TextLayout
	x, y, halo   float32
}

// frame is an open element.
type frame struct {
	name string
	st   style
	geo  map[string]string
}

// line is a line of text being read.
type line struct {
	x, y   float32
	anchor string
	spans  []ui.Span
}

// walker turns Pele's SVG into a scene.
type walker struct {
	pal   palette
	font  string
	out   scene
	stack []frame
	runs  int
	last  int
	lines []line
}

// parse reads pele-go's SVG; size is its font size.
func parse(svg string, pal palette, size float32, font string) (*scene, error) {
	w := &walker{pal: pal, font: font}
	root := style{fill: ui.RGB(0, 0, 0), fillA: 1, strokeA: 1, alpha: 1, width: 1, cap: ui.ButtCap, join: ui.MiterJoin,
		size: size, root: size, family: font, anchor: "start"}
	dec := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := dec.Token()
		if err == io.EOF && (w.out.w <= 0 || w.out.h <= 0) {
			return nil, fmt.Errorf("diagram: svg has no size")
		}
		if err == io.EOF {
			return &w.out, nil
		}
		if err == nil {
			switch t := tok.(type) {
			case xml.StartElement:
				err = w.open(t, root)
			case xml.EndElement:
				err = w.close()
			case xml.CharData:
				w.chars(string(t))
			}
		}
		if err != nil {
			return nil, fmt.Errorf("diagram: %w", err)
		}
	}
}

func (w *walker) top() *frame { return &w.stack[len(w.stack)-1] }

func (w *walker) open(el xml.StartElement, root style) error {
	st, outer := root, len(w.stack) == 0
	if !outer {
		st = w.top().st
	}
	from := st
	geo, css, decls := map[string]string{}, "", [][2]string{}
	// add keeps a value; CSS drops a bad corner.
	add := func(name, value string, styled bool) error {
		if name != "fill" && name != "stroke" {
			var err error
			if value, err = w.plain(value); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
		_, round := st.corner(value, 0)
		switch corner := name == "rx" || name == "ry"; {
		case styled && corner && !round:
		case corner || geometry[name] && !styled:
			geo[name] = value
		default:
			decls = append(decls, [2]string{name, value})
		}
		return nil
	}
	for _, a := range el.Attr {
		n := a.Name.Local
		if n == "style" {
			css = a.Value
		} else if !ignored[n] {
			if err := add(n, a.Value, false); err != nil {
				return err
			}
		}
	}
	// The outer style holds Pele's tokens alone.
	if !outer {
		err := declarations(css, func(name, value string) error { return add(name, value, true) })
		if err != nil {
			return err
		}
	}
	// Font size goes first: em lengths read it.
	for _, sized := range []bool{true, false} {
		for _, d := range decls {
			if (d[0] == "font-size") != sized {
				continue
			}
			if err := w.set(&st, &from, d[0], d[1]); err != nil {
				return err
			}
		}
	}
	w.stack = append(w.stack, frame{el.Name.Local, st, geo})
	switch el.Name.Local {
	case "svg":
		return w.svg(outer)
	case "g", "title", "desc":
		return nil
	case "rect", "circle", "ellipse", "polygon", "path":
		return w.shape()
	case "text", "tspan":
		return w.breakLine(el.Name.Local == "text")
	}
	return fmt.Errorf("unsupported element <%s>", el.Name.Local)
}

// svg sizes the scene, or places an icon's own box.
func (w *walker) svg(outer bool) error {
	f := w.top()
	n, err := f.numbers("x", "y", "width", "height")
	if err != nil {
		return err
	}
	if outer {
		w.out.w, w.out.h = n[2], n[3]
		return nil
	}
	box, ok := lengths(f.geo["viewBox"])
	if !ok || len(box) != 4 || box[2] <= 0 || box[3] <= 0 {
		return fmt.Errorf("unsupported viewBox %q", f.geo["viewBox"])
	}
	w.runs++
	f.st.run = w.runs
	f.st.t = ui.Translation(-box[0], -box[1]).Then(ui.Scaling(n[2]/box[2], n[3]/box[3])).Then(ui.Translation(n[0], n[1])).Then(f.st.t)
	return nil
}

// shape adds a path for the element on top.
func (w *walker) shape() error {
	f := w.top()
	d, err := f.data()
	if err != nil {
		return err
	}
	p, err := ui.ParsePath(d)
	if err != nil {
		return fmt.Errorf("<%s>: %w", f.name, err)
	}
	st := &f.st
	m := mark{path: p.Cap(st.cap).Join(st.join), fill: st.fill.Alpha(st.fillA * st.alpha), width: st.width}
	if st.width > 0 {
		m.stroke = st.stroke.Alpha(st.strokeA * st.alpha)
	}
	if len(st.dash) > 0 {
		p.Dash(st.dashOff, st.dash...)
	}
	if m.fill.A > 0 || m.stroke.A > 0 {
		w.emit(st, m)
	}
	return nil
}

// breakLine starts a line at a text or placed tspan.
func (w *walker) breakLine(text bool) error {
	f := w.top()
	if text {
		w.lines = w.lines[:0]
	} else if f.geo["x"] == "" && f.geo["y"] == "" {
		return nil
	}
	n, err := f.numbers("x", "y")
	if err != nil {
		return err
	}
	w.lines = append(w.lines, line{x: n[0], y: n[1], anchor: f.st.anchor})
	return nil
}

// chars keeps a title, or adds a span.
func (w *walker) chars(s string) {
	if len(w.stack) == 0 {
		return
	}
	f := w.top()
	switch {
	case f.name == "title" && len(w.stack) == 2:
		w.out.title += s
	// MyGo reads size zero as unset: draw nothing.
	case f.st.size > 0 && (f.name == "text" || f.name == "tspan"):
		st := &f.st
		sp := ui.Span{Font: st.family, Size: st.size, Weight: st.weight, Italic: st.italic, Underline: st.underline,
			Strikethrough: st.struck, BaselineShift: st.shift, Color: st.fill.Alpha(st.fillA * st.alpha)}
		if sp.Color.A == 0 {
			sp.Color = st.stroke.Alpha(st.strokeA * st.alpha)
		}
		l := &w.lines[len(w.lines)-1]
		for _, word := range words(s, st.word != 0) {
			sp.Text, sp.LetterSpacing = word, st.spacing
			if word == " " || word == string(nbsp) {
				sp.LetterSpacing += st.word
			}
			l.spans = append(l.spans, sp)
		}
	}
}

// nbsp is a word separator, as the space is.
const nbsp = '\u00a0'

// words cuts s around its spaces when they are spaced.
func words(s string, spaced bool) []string {
	if !spaced {
		return []string{s}
	}
	var out []string
	from := 0
	for i, r := range s {
		if r == ' ' || r == nbsp {
			to := i + utf8.RuneLen(r)
			out, from = append(out, s[from:i], s[i:to]), to
		}
	}
	return slices.DeleteFunc(append(out, s[from:]), func(w string) bool { return w == "" })
}

func (w *walker) close() error {
	f := w.top()
	if f.name == "text" {
		w.text(&f.st)
	}
	w.stack = w.stack[:len(w.stack)-1]
	return nil
}

// text shapes the lines read and places them by anchor.
func (w *walker) text(st *style) {
	font := ui.Font{Family: st.family, Size: st.size, Weight: st.weight, Italic: st.italic}
	for _, l := range w.lines {
		if len(l.spans) == 0 {
			continue
		}
		layout := ui.ShapeRichText(l.spans, font, 0)
		rows := layout.Lines()
		if len(rows) == 0 {
			continue
		}
		width, _ := layout.Size()
		m := mark{text: layout, x: l.x, y: l.y - rows[0].Baseline, fill: l.spans[0].Color}
		switch l.anchor {
		case "middle":
			m.x -= width / 2
		case "end":
			m.x -= width
		}
		if st.fill.A == 0 {
			m.halo = st.width / 2
		}
		w.emit(st, m)
	}
}

// emit adds a mark; a changed transform starts a run.
func (w *walker) emit(st *style, m mark) {
	if len(w.out.runs) == 0 || w.last != st.run {
		w.out.runs, w.last = append(w.out.runs, run{t: st.t}), st.run
	}
	r := &w.out.runs[len(w.out.runs)-1]
	r.marks = append(r.marks, m)
}
