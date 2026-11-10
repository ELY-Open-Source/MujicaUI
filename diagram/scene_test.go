package diagram

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// inks gives every token a color of its own.
var inks = palette{bg: ui.Hex("#010101"), fg: ui.Hex("#020202"), muted: ui.Hex("#030303"), line: ui.Hex("#040404"),
	surface: ui.Hex("#050505"), alt: ui.Hex("#060606"), border: ui.Hex("#070707"), accent: ui.Hex("#080808"),
	critical: ui.Hex("#090909"), radius: 6, series: [8]ui.Color{ui.Hex("#101010"), ui.Hex("#111111"), ui.Hex("#121212"),
		ui.Hex("#131313"), ui.Hex("#141414"), ui.Hex("#151515"), ui.Hex("#161616"), ui.Hex("#171717")}}

func draw(t *testing.T, source string) *scene {
	t.Helper()
	sc, err := render(source, inks, 14, "")
	if err != nil {
		t.Fatalf("%q: %v", source, err)
	}
	return sc
}

func marks(sc *scene) (paths, texts []mark) {
	for _, r := range sc.runs {
		for _, m := range r.marks {
			if m.text != nil {
				texts = append(texts, m)
			} else {
				paths = append(paths, m)
			}
		}
	}
	return paths, texts
}

func corpus(t *testing.T) map[string][]string {
	t.Helper()
	raw, err := os.ReadFile("testdata/corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var docs map[string][]string
	if err := json.Unmarshal(raw, &docs); err != nil {
		t.Fatal(err)
	}
	return docs
}

func TestCorpusDraws(t *testing.T) {
	docs := corpus(t)
	if len(docs) != 33 {
		t.Fatalf("%d diagram types, want 33", len(docs))
	}
	for kind, sources := range docs {
		for i, source := range sources {
			sc, err := render(source, inks, 14, "")
			if err != nil {
				t.Errorf("%s %d: %v", kind, i, err)
				continue
			}
			paths, texts := marks(sc)
			// Gantt 6 holds settings and no task.
			if blank := kind == "gantt" && i == 6; sc.w <= 0 || sc.h <= 0 || (len(paths)+len(texts) == 0) != blank {
				t.Errorf("%s %d: scene %gx%g with %d paths and %d texts", kind, i, sc.w, sc.h, len(paths), len(texts))
			}
		}
	}
}

func TestFlowchartScene(t *testing.T) {
	sc := draw(t, "flowchart LR\n  A[Kiln] --> B[Glaze]")
	paths, texts := marks(sc)
	if len(paths) != 4 || len(texts) != 2 {
		t.Fatalf("%d paths and %d texts, want 2 boxes, an edge, its head and 2 labels", len(paths), len(texts))
	}
	edge, head, box := paths[0], paths[1], paths[2]
	if edge.stroke != inks.line || edge.fill.A != 0 || edge.width != 1 {
		t.Errorf("edge %+v, want a 1 DIP line in the line color", edge)
	}
	if head.fill != inks.line || head.stroke.A != 0 {
		t.Errorf("arrowhead %+v, want it filled in the line color", head)
	}
	if box.fill != inks.surface || box.stroke != inks.border {
		t.Errorf("box %+v, want surface fill and border stroke", box)
	}
	b := box.path.Bounds()
	if abs(b.X+b.W/2) > 0.01 || abs(b.Y+b.H/2) > 0.01 {
		t.Errorf("box bounds %v, want it centered on its group's origin", b)
	}
	label := texts[0]
	w, h := label.text.Size()
	if label.fill != inks.fg || label.halo != 0 || abs(label.x+w/2) > 0.01 || label.y > 0 || label.y+h < 0 {
		t.Errorf("label %+v sized %gx%g, want foreground text centered in its box", label, w, h)
	}
	if len(sc.runs) != 3 {
		t.Errorf("%d runs, want the edges and each node under its own transform", len(sc.runs))
	}
}

func abs(v float32) float32 { return max(v, -v) }

func TestStylesFromTheSource(t *testing.T) {
	sc := draw(t, "flowchart LR\n  A --> B\n  style A fill:#f9f,stroke:#00f8,stroke-width:4px,stroke-dasharray:5 5\n"+
		"  style B fill:tomato,stroke:green,color:banana")
	paths, texts := marks(sc)
	a, b := paths[2], paths[3]
	if a.fill != ui.Hex("#ff99ff") || a.stroke != ui.Hex("#0000ff88") || a.width != 4 {
		t.Errorf("A %+v, want its own fill, stroke and width", a)
	}
	if b.fill != ui.Hex("#ff6347") || b.stroke != ui.Hex("#008000") {
		t.Errorf("B %+v, want tomato and green", b)
	}
	if texts[1].fill != inks.fg {
		t.Errorf("label color %v, want a color CSS drops to leave the foreground", texts[1].fill)
	}
}

func TestTitleAndHalo(t *testing.T) {
	if sc := draw(t, "flowchart LR\n  accTitle: Firing plan\n  A --> B"); sc.title != "Firing plan" {
		t.Errorf("title %q", sc.title)
	}
	_, texts := marks(draw(t, corpus(t)["sankey"][5]))
	halos := 0
	for _, m := range texts {
		if m.halo > 0 {
			halos++
			if m.fill != inks.bg {
				t.Errorf("halo in %v, want the background", m.fill)
			}
		}
	}
	if halos == 0 {
		t.Error("no label has a halo")
	}
}

func TestColors(t *testing.T) {
	for in, want := range map[string]ui.Color{
		"#abc": ui.Hex("#aabbcc"), "#ABCD": ui.Hex("#aabbccdd"), "#102030": ui.Hex("#102030"), "#10203040": ui.Hex("#10203040"),
		"rgb(255, 0, 128)": ui.RGB(255, 0, 128), "rgba(0,0,0,0.5)": ui.RGBA(0, 0, 0, 0.5), "rgb(100% 0% 0%)": ui.RGB(255, 0, 0),
		"hsl(0, 100%, 50%)": ui.RGB(255, 0, 0), "hsla(240 100% 50% / 0.5)": ui.RGBA(0, 0, 255, 0.5), "RebeccaPurple": ui.Hex("#663399"),
		"none": ui.Transparent, "transparent": ui.Transparent, "currentColor": inks.accent,
		"var(--_s)": inks.surface, "var(--pele-series-3,#54a24b)": inks.series[2], "var(--pele-critical,#e45756)": inks.critical,
		"var(--mine,#123)": ui.Hex("#112233"), "var(--pele-accent,#5b7bd5)": inks.accent, "color-mix(in srgb,var(--_a) 50%,var(--_bg))": inks.alt.Mix(inks.bg, 0.5),
	} {
		got, ok, err := inks.color(in, inks.accent)
		if !ok || err != nil || got != want {
			t.Errorf("%s = %v %v %v, want %v", in, got, ok, err, want)
		}
	}
	for _, in := range []string{"banana", "#12", "#ggg", "rgb(1,2)", "hsl(0, 1, 1)", "rgb(a,b,c)", "url(#x)",
		"color-mix(in oklab,red,blue)", "color-mix(red,blue)", "inherit"} {
		if got, ok, err := inks.color(in, inks.accent); ok || err != nil {
			t.Errorf("%s = %v %v %v, want it dropped", in, got, ok, err)
		}
	}
	if _, _, err := inks.color("var(--_zz)", inks.accent); err == nil {
		t.Error("var(--_zz): no error")
	}
	if len(named) != 148 {
		t.Errorf("%d named colors, want 148", len(named))
	}
}

func TestUnknownMarkupFails(t *testing.T) {
	for in, want := range map[string]string{
		`<svg width="10" height="10"><line/></svg>`:                                       "unsupported element <line>",
		`<svg width="10" height="10"><g clip-path="url(#a)"/></svg>`:                      `unsupported attribute "clip-path"`,
		`<svg width="10" height="10"><g style="cursor:pointer"/></svg>`:                   `unsupported attribute "cursor"`,
		`<svg width="10" height="10"><g transform="scale(2)"/></svg>`:                     `unsupported transform "scale(2)"`,
		`<svg width="10" height="10"><rect fill="var(--_zz)"/></svg>`:                     `unknown token "--_zz"`,
		`<svg width="10" height="10"><rect fill="var(--pele-new,#123)"/></svg>`:           `unknown token "--pele-new"`,
		`<svg width="10" height="10"><g stroke-width="var(--_new,2)"/></svg>`:             `unknown token "--_new"`,
		`<svg width="10" height="10"><g font-weight="var(--pele-new-weight,400)"/></svg>`: `unknown token "--pele-new-weight"`,
		`<svg width="10" height="10"><g style="rx:calc(var(--_q)*2)"/></svg>`:             `unknown token "--_q"`,
		`<svg width="10" height="10"><rect width="wide"/></svg>`:                          `width="wide" is not a number`,
		`<svg width="10" height="10"><svg width="4" height="4"/></svg>`:                   "unsupported viewBox",
		`<svg width="10" height="10"><path d="M0 0 L"/></svg>`:                            "<path>",
		`<svg width="10" height="10"><g>`:                                                 "diagram:",
	} {
		if _, err := parse(in, inks, 14, ""); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", in, err, want)
		}
	}
}

func TestShapes(t *testing.T) {
	sc, err := parse(`<svg width="100" height="100" fill="var(--_fg)">`+
		`<rect x="10" y="20" width="30" height="40" style="rx:var(--_r)"/><rect width="8" height="8" rx="9"/>`+
		`<circle cx="50" cy="50" r="5"/><ellipse rx="6" ry="3"/><polygon points="0,0 10,0 5,8"/>`+
		`<svg x="20" y="30" width="48" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M0 0H24"/></svg>`+
		`<rect width="5" height="5" fill="none"/><rect width="5" height="5" fill="none" stroke="red" stroke-width="0"/></svg>`, inks, 14, "")
	if err != nil {
		t.Fatal(err)
	}
	paths, _ := marks(sc)
	want := []ui.Rect{{X: 10, Y: 20, W: 30, H: 40}, {W: 8, H: 8}, {X: 45, Y: 45, W: 10, H: 10}, {X: -6, Y: -3, W: 12, H: 6}, {W: 10, H: 8}, {W: 24}}
	if len(paths) != len(want) {
		t.Fatalf("%d paths, want %d: unpainted shapes leave no mark", len(paths), len(want))
	}
	for i, m := range paths {
		if b := m.path.Bounds(); abs(b.X-want[i].X)+abs(b.Y-want[i].Y)+abs(b.W-want[i].W)+abs(b.H-want[i].H) > 0.01 {
			t.Errorf("shape %d bounds %v, want %v", i, b, want[i])
		}
	}
	if icon := paths[5]; icon.stroke != inks.fg || icon.fill.A != 0 || len(sc.runs) != 2 {
		t.Errorf("icon %+v in %d runs, want a foreground stroke in a run of its own", icon, len(sc.runs))
	}
	if sc.w != 100 || sc.h != 100 {
		t.Errorf("scene %gx%g", sc.w, sc.h)
	}
}

func TestColorMixIsPremultiplied(t *testing.T) {
	for in, want := range map[string]ui.Color{
		"color-mix(in srgb, red 50%, transparent)": ui.RGBA(255, 0, 0, 0.5),
		"color-mix(in srgb, red, blue 25%)":        ui.RGB(191, 0, 64),
		"color-mix(in srgb, 25% red, blue)":        ui.RGB(64, 0, 191),
		"color-mix(in srgb, red 20%, blue 20%)":    ui.RGBA(128, 0, 128, 0.4),
		"color-mix(in srgb, red 150%, blue 50%)":   ui.Transparent,
	} {
		got, ok, err := inks.color(in, inks.fg)
		if dropped := want == ui.Transparent; err != nil || ok == dropped || !dropped && got != want {
			t.Errorf("%s = %v %v %v, want %v", in, got, ok, err, want)
		}
	}
	paths, _ := sketch(t, `<rect width="9" height="9" fill="red" style="fill:color-mix(in oklab,red,blue)"/>`)
	if paths[0].fill != ui.RGB(255, 0, 0) {
		t.Errorf("fill %v, want red: a mix outside srgb drops", paths[0].fill)
	}
}

// sketch parses a hand-written SVG body.
func sketch(t *testing.T, body string) (paths, texts []mark) {
	t.Helper()
	sc, err := parse(`<svg width="100" height="100">`+body+`</svg>`, inks, 14, "")
	if err != nil {
		t.Fatalf("%s: %v", body, err)
	}
	return marks(sc)
}

func TestStyleBeatsAttribute(t *testing.T) {
	paths, _ := sketch(t, `<g opacity="0.5"><rect width="9" height="9" opacity="0.5" style="opacity:0.8"/>`+
		`<rect width="9" height="9" opacity="0.5" style="opacity:much"/></g>`)
	if paths[0].fill.A != 102 || paths[1].fill.A != 64 {
		t.Errorf("alphas %d and %d, want 102 from the style and 64 from the attribute", paths[0].fill.A, paths[1].fill.A)
	}
}

func TestAnchorIsTheLinesOwn(t *testing.T) {
	_, texts := sketch(t, `<text text-anchor="middle"><tspan x="20" y="30" text-anchor="start">abc</tspan>`+
		`<tspan x="20" y="50">abc</tspan></text>`)
	w, _ := texts[0].text.Size()
	if len(texts) != 2 || texts[0].x != 20 || texts[1].x != 20-w/2 {
		t.Errorf("lines at %v and %v, want 20 and %v", texts[0].x, texts[1].x, 20-w/2)
	}
}

func TestParseNeedsASizedRoot(t *testing.T) {
	for _, svg := range []string{"", `<svg width="0" height="10"></svg>`, `<svg></svg>`} {
		if _, err := parse(svg, inks, 14, ""); err == nil {
			t.Errorf("%q parsed", svg)
		}
	}
}

func TestADiagramsOwnStyles(t *testing.T) {
	styled := func(decls string) (box mark, label mark) {
		paths, texts := marks(draw(t, "flowchart LR\n  A[Fire the kiln] --> B\n  style A "+decls))
		return paths[2], texts[0]
	}
	plain, _ := styled("fill:#fff")
	_, text := styled("fill:#fff")
	width, _ := text.text.Size()
	if box, label := styled("opacity:50%"); box.fill.A != 128 || label.fill.A != 128 {
		t.Errorf("opacity 50%%: alphas %d and %d, want 128", box.fill.A, label.fill.A)
	}
	if box, _ := styled("stroke-width:1.5pt,stroke-opacity:25%"); box.width != 2 || box.stroke.A != 64 {
		t.Errorf("stroke %v wide at alpha %d, want 2 at 64", box.width, box.stroke.A)
	}
	if box, _ := styled("rx:50%"); box.path.Bounds() != plain.path.Bounds() {
		t.Errorf("a rounded box moved: %v", box.path.Bounds())
	}
	_, spaced := styled("word-spacing:10px")
	if w, _ := spaced.text.Size(); abs(w-width-20) > 0.5 {
		t.Errorf("word spacing: %v wide from %v, want 20 more", w, width)
	}
	for decls, want := range map[string]float32{"font-size:21px": 1.5, "font-size:150%": 1.5, "font-size:1.5em": 1.5,
		"font-size:1.5rem": 1.5, "font-size:larger": 1.2, "font-size:x-large": 24.0 / 14} {
		_, label := styled(decls)
		if _, h := label.text.Size(); abs(h/14-want*17/14) > 0.15 {
			t.Errorf("%s: line %v high, want about %v", decls, h, want*17)
		}
	}
}

func TestCSSValues(t *testing.T) {
	st := style{size: 20, root: 10}
	for in, want := range map[string]float32{"3": 3, "3px": 3, "3pt": 4, "1in": 96, "2em": 40, "2rem": 20, "1e1px": 10, "-2": -2} {
		if got, ok := st.measure(in); !ok || abs(got-want) > 0.001 {
			t.Errorf("measure %q = %v %v, want %v", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "px", "3vw", "wide", "3 px"} {
		if got, ok := st.measure(in); ok {
			t.Errorf("measure %q = %v", in, got)
		}
	}
	for _, c := range [][4]int{{400, 700, 100, 0}, {300, 400, 100, 0}, {700, 900, 400, 0}, {900, 900, 700, 0}} {
		if b, l := weight("bolder", c[0], 0), weight("lighter", c[0], 0); b != c[1] || l != c[2] {
			t.Errorf("from %d: bolder %d lighter %d, want %d and %d", c[0], b, l, c[1], c[2])
		}
	}
	paths, texts := sketch(t, `<g style="letter-spacing:2px;text-decoration:underline line-through;font-family:'Fira Code', monospace">`+
		`<text x="0" y="20" style="letter-spacing:normal;text-anchor:sideways">ab</text></g>`+
		`<rect width="9" height="9" fill="red" style="fill:var(--mine);stroke-dasharray:4 x"/>`)
	if w, _ := texts[0].text.Size(); texts[0].x != 0 || w > 30 {
		t.Errorf("text at %v, %v wide: want the start anchor and no tracking", texts[0].x, w)
	}
	if paths[0].fill != ui.RGB(255, 0, 0) {
		t.Errorf("fill %v, want red: an unset variable of the diagram's own is dropped", paths[0].fill)
	}
}

func TestTokensOutsideColors(t *testing.T) {
	paths, texts := sketch(t, `<rect width="40" height="40" style="rx:var(--pele-radius,4px)"/>`+
		`<rect width="40" height="40" stroke="red" stroke-width="var(--mine, 3px)" style="rx:var(--mine)"/>`+
		`<text x="0" y="20" font-weight="var(--_tw)" font-family="var(--_fm)">mm</text>`+
		`<text x="0" y="40" font-weight="var(--pele-title-weight,bold)" font-family="var(--pele-font-mono,monospace)">mm</text>`+
		`<text x="0" y="60">mm</text>`)
	if !paths[0].path.Contains(20, 20) || paths[0].path.Contains(0.5, 0.5) || !paths[1].path.Contains(0.5, 0.5) {
		t.Error("want the theme's radius on the first box and square corners on the second")
	}
	if paths[1].width != 3 {
		t.Errorf("stroke %v wide, want the fallback's 3", paths[1].width)
	}
	private, _ := texts[0].text.Size()
	public, _ := texts[1].text.Size()
	plain, _ := texts[2].text.Size()
	if private != public || private == plain {
		t.Errorf("widths %v, %v and %v: want both token texts bold and monospace", private, public, plain)
	}
}

func TestEmReadsTheFinalFontSize(t *testing.T) {
	width := func(css string) float32 {
		_, texts := sketch(t, `<text x="0" y="40" style="`+css+`">a b</text>`)
		w, _ := texts[0].text.Size()
		return w
	}
	for _, prop := range []string{"letter-spacing", "word-spacing"} {
		first, last := width(prop+":1em;font-size:28px"), width("font-size:28px;"+prop+":1em")
		if first != last || first <= width("font-size:28px")+27 {
			t.Errorf("%s: %v wide before the font size, %v after", prop, first, last)
		}
	}
}

func TestABadValueKeepsTheGoodOne(t *testing.T) {
	paths, texts := sketch(t, `<rect width="40" height="40" style="rx:8px;rx:banana"/>`+
		`<text x="0" y="20" style="text-decoration:underline;text-decoration:banana;font-size:20px;font-size:-3px">mm</text>`+
		`<text x="0" y="40" style="font-size:20px">mm</text>`)
	if paths[0].path.Contains(0.5, 0.5) {
		t.Error("a bad rx squared the corner the good one rounded")
	}
	kept, _ := texts[0].text.Size()
	want, _ := texts[1].text.Size()
	if kept != want {
		t.Errorf("text %v wide after a bad font size, want %v", kept, want)
	}
	if _, none := sketch(t, `<text x="0" y="20" style="font-size:0">mm</text>`); len(none) != 0 {
		t.Error("text of no size left a mark")
	}
	styled, _ := marks(draw(t, "flowchart LR\n  A[Fire the kiln] --> B\n  style A rx:banana"))
	if b := styled[2].path.Bounds(); styled[2].path.Contains(b.X+0.5, b.Y+0.5) {
		t.Error("a diagram's bad rx removed the theme's radius")
	}
}

func TestOnlyCSSNumbersCount(t *testing.T) {
	for _, in := range []string{"inf", "nan", "Infinity", "0x10", "1_0", "1e99", ""} {
		if v, ok := number(in); ok {
			t.Errorf("number %q = %v", in, v)
		}
		if got, ok, err := inks.color("rgb("+in+" 0 0)", inks.fg); ok || err != nil {
			t.Errorf("rgb(%s 0 0) = %v %v", in, got, err)
		}
	}
	paths, texts := sketch(t, `<rect width="9" height="9" opacity="0.5" style="opacity:inf"/>`+
		`<text x="0" y="20" font-weight="650.4" style="word-spacing:10px">a&#160;b c</text><text x="0" y="40" font-weight="650">a&#160;b c</text>`)
	if paths[0].fill.A != 128 {
		t.Errorf("alpha %d, want 128: inf is no CSS number", paths[0].fill.A)
	}
	spaced, _ := texts[0].text.Size()
	plain, _ := texts[1].text.Size()
	if abs(spaced-plain-20) > 0.5 {
		t.Errorf("%v wide against %v, want 20 more for a space and a no-break space", spaced, plain)
	}
}

func TestCSSWideKeywordsDrop(t *testing.T) {
	paths, texts := sketch(t, `<g font-family="monospace" font-size="20" fill="red" opacity="0.5">`+
		`<text x="0" y="20" style="font-family:inherit;font-size:initial;fill:unset;opacity:revert">mmmm</text>`+
		`<text x="0" y="40">mmmm</text><rect width="9" height="9" style="stroke-width:inherit;rx:initial"/></g>`)
	kept, _ := texts[0].text.Size()
	want, _ := texts[1].text.Size()
	if kept != want || texts[0].fill != texts[1].fill || len(paths) != 1 {
		t.Errorf("text %v wide in %v, want %v in %v: a CSS-wide keyword leaves the earlier value", kept, texts[0].fill, want, texts[1].fill)
	}
	if _, err := parse(`<svg width="9" height="9"><g style="cursor:inherit"/></svg>`, inks, 14, ""); err == nil {
		t.Error("an unknown property passed behind a CSS-wide keyword")
	}
}
