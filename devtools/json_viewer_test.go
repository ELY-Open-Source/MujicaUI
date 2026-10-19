package devtools

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

const jsonSample = `{"ball":{"hall":"east wing","guests":[{"name":"Sakiko","seat":1},{"name":"Uika","seat":2}],"wine":null},"open":true,"two words":3.5}`

func jsonViewerTester(data *string, s *JsonViewerState, opts JsonViewerOptions, copied *[]string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := JsonViewer(c, *data, s, opts); r.Copied != "" {
			*copied = append(*copied, r.Copied)
		}
	}, 640, 460)
}

func TestJsonTreePaths(t *testing.T) {
	tr, err := parseJSONTree(jsonSample)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, n := range tr.nodes {
		paths = append(paths, n.path)
	}
	want := `$ $.ball $.ball.hall $.ball.guests $.ball.guests[0] $.ball.guests[0].name $.ball.guests[0].seat $.ball.guests[1] $.ball.guests[1].name $.ball.guests[1].seat $.ball.wine $.open $["two words"]`
	if strings.Join(paths, " ") != want {
		t.Fatalf("paths %v", paths)
	}
	if tr.nodes[3].count != 2 || tr.nodes[3].end != 10 {
		t.Fatalf("guests node %+v", tr.nodes[3])
	}
	if _, err := parseJSONTree(`{"a":1} x`); err == nil {
		t.Fatal("trailing data accepted")
	}
}

func TestJsonViewerFoldKeyboardCopy(t *testing.T) {
	data := jsonSample
	var s JsonViewerState
	var copied []string
	tt := jsonViewerTester(&data, &s, JsonViewerOptions{}, &copied)
	if !tt.HasText("ball: { 3 items }") || tt.HasText(`hall: "east wing"`) {
		t.Fatalf("closed: %q", tt.Texts())
	}
	_ = tt.Click("$.ball") // its fold mark
	if !tt.HasText(`hall: "east wing"`) || !tt.HasText("guests: [ 2 items ]") {
		t.Fatalf("open: %q", tt.Texts())
	}
	_ = tt.Click(`hall: "east wing"`)
	_ = tt.Click("Copy path")
	if len(copied) != 1 || copied[0] != "$.ball.hall" || tt.Clipboard() != "$.ball.hall" {
		t.Fatalf("copy: %v", copied)
	}
	_ = tt.Click(`hall: "east wing"`) // back on the list
	tt.Key(0, ui.KeyDown)             // guests
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("[0]: { 2 items }") && !tt.HasText("{ 2 items }") {
		t.Fatalf("Right did not open guests: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyLeft)
	if tt.HasText("{ 2 items }") {
		t.Fatal("Left did not close guests")
	}
}

func TestJsonViewerSearchOpensMatches(t *testing.T) {
	data := jsonSample
	var s JsonViewerState
	var copied []string
	tt := jsonViewerTester(&data, &s, JsonViewerOptions{}, &copied)
	_ = tt.Click("Find")
	tt.Type("uika")
	if !tt.HasText("1 of 1") || !tt.HasText("$.ball.guests[1].name") {
		t.Fatalf("search: %q", tt.Texts())
	}
	data = `{"broken": }`
	tt.Frame()
	found := false
	for _, x := range tt.Texts() {
		found = found || strings.HasPrefix(x, "Invalid JSON")
	}
	if !found {
		t.Fatalf("parse error: %q", tt.Texts())
	}
}

func TestJsonViewerLargeDocumentVirtualized(t *testing.T) {
	var b strings.Builder
	b.WriteString("[")
	for i := range 100000 {
		if i > 0 {
			b.WriteString(",")
		}
		fmt.Fprintf(&b, "%d", i)
	}
	b.WriteString("]")
	data := b.String()
	var s JsonViewerState
	var copied []string
	tt := jsonViewerTester(&data, &s, JsonViewerOptions{Expanded: true}, &copied)
	if n := len(tt.Texts()); n > 300 {
		t.Fatalf("%d texts for 100000 items", n)
	}
}

func TestJsonViewerSnapshot(t *testing.T) {
	snapshot(t, "json_viewer", 600, 360, func(c *ui.Context) {
		var s JsonViewerState
		ui.Column(c).Padding(16).Children(func() {
			JsonViewer(c, jsonSample, &s, JsonViewerOptions{Expanded: true, Height: 320})
		})
	})
}

func TestJsonViewerAncestorDisabledKeepsTreeAndClipboard(t *testing.T) {
	data := jsonSample
	var s JsonViewerState
	disabled := false
	copied := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if JsonViewer(c, data, &s, JsonViewerOptions{}).Copied != "" {
				copied++
			}
		})
	}, 640, 460)
	if err := tt.Click("ball: { 3 items }"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyRight)
	_ = tt.Click("$.ball")
	_ = tt.Click("Copy path")
	if s.open["$.ball"] || copied != 0 || tt.Clipboard() != "" {
		t.Fatal("disabled tree input")
	}
}

func TestJsonViewerAncestorDisabledSelection(t *testing.T) {
	var state JsonViewerState
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			JsonViewer(c, `{"a":1,"b":2}`, &state, JsonViewerOptions{})
		})
	}, 640, 460)
	if err := tt.Click("a: 1"); err != nil {
		t.Fatal(err)
	}
	before := state.sel
	disabled = true
	tt.Key(0, ui.KeyDown)
	if state.sel != before {
		t.Fatal("disabled JSON selection changed")
	}
}
