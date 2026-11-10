package diagram

import (
	"bytes"
	"context"
	"image"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// look is the tokens and font size a view drew with.
type look struct {
	k    theme.Tokens
	size float32
}

const flow = "flowchart LR\n  accTitle: Order flow\n  Cart[Cart] --> Pay{Paid?}\n  Pay -- yes --> Ship([Ship order])\n" +
	"  Pay -- no --> Cart\n  style Ship fill:#0a7,stroke:#064,color:#fff"

// show opens a window holding one Mermaid of *src.
func show(src *string, o MermaidOptions, w int, seen *look) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if seen != nil {
			*seen = look{core.Tokens(c), core.FontSize(c, theme.BodySize)}
		}
		Mermaid(c, *src, o)
	}, w, 400)
	tt.Advance(time.Second)
	return tt
}

// logged counts log lines holding s.
func logged(t *testing.T) func(s string) int {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(old) })
	return func(s string) int { return strings.Count(buf.String(), s) }
}

func picture(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	info, ok := tt.AccessByRole(ui.RoleImage, name)
	if !ok {
		t.Fatalf("no image named %q", name)
	}
	return info.Bounds
}

func at(img *image.RGBA, x, y float32) ui.Color {
	c := img.RGBAAt(int(x), int(y))
	return ui.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// spots finds where Pele's SVG puts each node and arrowhead.
func spots(t *testing.T, seen *look) (nodes map[string][2]float32, heads [][2]float32, w float32) {
	t.Helper()
	off := false
	r, err := pele.Render(context.Background(), flow, pele.Options{Measurer: measurer{""},
		FontSize: pele.Ptr(float64(seen.size)), Links: &off, Images: &off})
	if err != nil {
		t.Fatal(err)
	}
	num := func(s string) float32 {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			t.Fatal(err)
		}
		return float32(f)
	}
	nodes = map[string][2]float32{}
	for _, m := range regexp.MustCompile(`data-id="(\w+)" transform="translate\(([\d.]+),([\d.]+)\)"`).FindAllStringSubmatch(r.SVG, -1) {
		nodes[m[1]] = [2]float32{num(m[2]), num(m[3])}
	}
	for _, m := range regexp.MustCompile(`class="pele-marker" d="M([\d.]+),([\d.]+)L([\d.]+),`).FindAllStringSubmatch(r.SVG, -1) {
		heads = append(heads, [2]float32{(num(m[1]) + 2*num(m[3])) / 3, num(m[2])})
	}
	if len(nodes) != 3 || len(heads) != 3 {
		t.Fatalf("%d nodes and %d arrowheads in %s", len(nodes), len(heads), r.SVG)
	}
	return nodes, heads, float32(r.Width)
}

func TestMermaidDrawsWherePeleDraws(t *testing.T) {
	src, seen := flow, look{}
	tt := show(&src, MermaidOptions{}, 700, &seen)
	nodes, heads, w := spots(t, &seen)
	if box := picture(t, tt, "Order flow"); box.X != 0 || box.Y != 0 || box.W != w {
		t.Fatalf("picture %+v, want Pele's width %v at the corner", box, w)
	}
	img := tt.Image()
	for id, fill := range map[string]ui.Color{"Cart": seen.k.Surface, "Pay": seen.k.Surface, "Ship": ui.Hex("#00aa77")} {
		spot := nodes[id]
		if got := at(img, spot[0], spot[1]-14); got != fill {
			t.Errorf("node %s draws %+v, want %+v", id, got, fill)
		}
		if ink(img, spot[0]-12, spot[1]-6, 24, 12) == 0 {
			t.Errorf("node %s has no label", id)
		}
	}
	for _, h := range heads {
		if got := at(img, h[0], h[1]); got != seen.k.Text {
			t.Errorf("arrowhead at %v draws %+v, want %+v", h, got, seen.k.Text)
		}
	}
	cart, pay := nodes["Cart"], nodes["Pay"]
	if got := at(img, (cart[0]+pay[0])/2-20, heads[0][1]); got == seen.k.Background {
		t.Error("no edge between Cart and Pay")
	}
	if got := at(img, w+20, 20); got != seen.k.Background {
		t.Errorf("beside the diagram %+v, want the background", got)
	}
}

// ink counts pixels unlike the box's corner.
func ink(img *image.RGBA, x, y, w, h float32) int {
	n, ground := 0, at(img, x, y)
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if at(img, px, py) != ground {
				n++
			}
		}
	}
	return n
}

func TestMermaidFitsItsWidth(t *testing.T) {
	src, seen := flow, look{}
	tt := show(&src, MermaidOptions{}, 260, &seen)
	nodes, _, w := spots(t, &seen)
	sc, err := render(flow, tokens(seen.k), seen.size, "")
	if err != nil {
		t.Fatal(err)
	}
	box := picture(t, tt, "Order flow")
	if want := 260 * sc.h / sc.w; w <= 260 || sc.w != w || box.W != 260 || box.H < want-0.5 || box.H > want+0.5 {
		t.Fatalf("picture %+v, want 260 x %v of a scene %v wide", box, want, w)
	}
	k := 260 / w
	if got := at(tt.Image(), nodes["Ship"][0]*k, (nodes["Ship"][1]-14)*k); got != ui.Hex("#00aa77") {
		t.Errorf("shrunk node draws %+v", got)
	}
}

func TestMermaidFollowsThemeAndTextScale(t *testing.T) {
	count := logged(t)
	src, seen := flow, look{}
	tt := show(&src, MermaidOptions{}, 900, &seen)
	for range 5 {
		tt.Frame()
	}
	light, size := seen.k.Surface, picture(t, tt, "Order flow")
	if n := count("mermaid drew"); n != 1 {
		t.Fatalf("drew %d times over six frames", n)
	}
	tt.SetDark(true)
	tt.Advance(time.Second)
	if seen.k.Surface == light || count("mermaid drew") != 2 {
		t.Fatalf("dark theme: surface %+v, drew %d", seen.k.Surface, count("mermaid drew"))
	}
	nodes, _, _ := spots(t, &seen)
	if got := at(tt.Image(), nodes["Cart"][0], nodes["Cart"][1]-14); got != seen.k.Surface {
		t.Errorf("Cart draws %+v in the dark, want %+v", got, seen.k.Surface)
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1.5})
	tt.Advance(time.Second)
	if grown := picture(t, tt, "Order flow"); grown.W <= size.W || count("mermaid drew") != 3 {
		t.Errorf("text scale 1.5: width %v from %v, drew %d", grown.W, size.W, count("mermaid drew"))
	}
	if tt.FramePending() {
		t.Error("still moving")
	}
}

func TestMermaidShowsPelesError(t *testing.T) {
	count := logged(t)
	src := "flowchart LR\n  Cart -->"
	tt := show(&src, MermaidOptions{}, 500, nil)
	for range 3 {
		tt.Frame()
	}
	if !tt.HasText("Diagram failed") || count("mermaid failed") != 1 || count("line=3 column=1") != 1 {
		t.Fatalf("texts %q, failed logged %d", tt.Texts(), count("mermaid failed"))
	}
	if count("Cart") != 0 {
		t.Error("the log holds the source")
	}
	told := false
	for _, s := range tt.Texts() {
		told = told || strings.HasPrefix(s, "3:1 Parse error on line 3:") && strings.Contains(s, "Expecting 'NODE_STRING', got 'EOF'")
	}
	if !told {
		t.Errorf("texts %q lack Pele's message", tt.Texts())
	}
	if _, ok := tt.AccessByRole(ui.RoleImage, "Diagram"); ok {
		t.Error("a failed diagram still draws")
	}
	src = "flowchart LR\n  Cart --> Pay"
	tt.Frame()
	tt.Advance(time.Second)
	picture(t, tt, "Diagram")
	if tt.HasText("Diagram failed") || count("mermaid drew") != 1 {
		t.Errorf("after the fix: texts %q, drew %d", tt.Texts(), count("mermaid drew"))
	}
	src = "no such diagram"
	tt.Frame()
	if !tt.HasText("No diagram type detected for the given text.") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMermaidLabels(t *testing.T) {
	src := flow
	picture(t, show(&src, MermaidOptions{Label: "Checkout"}, 500, nil), "Checkout")
	src = "pie\n  \"Tea\" : 3\n  \"Coffee\" : 5"
	picture(t, show(&src, MermaidOptions{}, 500, nil), "Diagram")
	src = "not a diagram"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		Mermaid(c, src, MermaidOptions{})
	}, 500, 400)
	if !tt.HasText("图表出错") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMermaidSnapshot(t *testing.T) {
	snapshot(t, "mermaid", 520, 560, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			Mermaid(c, flow, MermaidOptions{})
			Mermaid(c, "pie title Set list by tempo\n  \"Slow\" : 4\n  \"Mid\" : 6\n  \"Fast\" : 3", MermaidOptions{})
			Mermaid(c, "flowchart LR\n  Score -->", MermaidOptions{})
		})
	})
}

func TestMermaidPaintsInsideItsBox(t *testing.T) {
	var k theme.Tokens
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		k = core.Tokens(c)
		ui.Column(c).Padding(60).Children(func() {
			Mermaid(c, "flowchart LR\n  Cart[Cart] --> Pay[Pay]\n  style Cart stroke:#f00,stroke-width:40px", MermaidOptions{Label: "Wide stroke"})
		})
	}, 700, 300)
	tt.Advance(time.Second)
	box := picture(t, tt, "Wide stroke")
	img := tt.Image()
	if got := at(img, box.X+2, box.Y+box.H/2); got != ui.Hex("#ff0000") {
		t.Fatalf("inside the left edge %+v, want the wide red stroke", got)
	}
	if got := at(img, box.X-3, box.Y+box.H/2); got != k.Background {
		t.Errorf("outside the left edge %+v, want the background", got)
	}
}
