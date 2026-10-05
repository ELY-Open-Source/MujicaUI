package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func traceSample() []TraceSpan {
	ms := time.Millisecond
	return []TraceSpan{
		{ID: "root", Name: "agent.run", Start: 0, End: 900 * ms, Attrs: []TraceAttr{{"model", "opus"}}},
		{ID: "plan", Parent: "root", Name: "llm.plan", Start: 10 * ms, End: 210 * ms, Attrs: []TraceAttr{{"tokens", "1200"}}},
		{ID: "search", Parent: "root", Name: "tool.search", Start: 220 * ms, End: 520 * ms, Failed: true},
		{ID: "fetch", Parent: "search", Name: "http.get", Start: 230 * ms, End: 500 * ms},
		{ID: "write", Parent: "root", Name: "llm.write", Start: 530 * ms, End: 880 * ms},
	}
}

func TestTraceBuild(t *testing.T) {
	ix := traceBuild(traceSample())
	if len(ix.roots) != 1 || ix.roots[0] != "root" || len(ix.kids["root"]) != 3 || ix.end != 900*time.Millisecond {
		t.Fatalf("index %+v", ix)
	}
	for name, spans := range map[string][]TraceSpan{
		"dup":   {{ID: "a"}, {ID: "a"}},
		"time":  {{ID: "a", Start: 2, End: 1}},
		"cycle": {{ID: "a", Parent: "b"}, {ID: "b", Parent: "a"}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s did not panic", name)
				}
			}()
			traceBuild(spans)
		}()
	}
}

func TestTraceViewerSelectShowsAttributes(t *testing.T) {
	var s TraceViewerState
	spans := traceSample()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(600).Children(func() {
			TraceViewer(c, &s, spans, TraceViewerOptions{}).Element.Grow(1)
		})
	}, 900, 620)
	if !tt.HasText("llm.plan") || !tt.HasText("200.0 ms") {
		t.Fatalf("tree: %q", tt.Texts())
	}
	if _, ok := tt.Find("tool.search, failed"); !ok {
		t.Fatal("failed span not named")
	}
	if tt.HasText("http.get") {
		t.Fatal("nested span shown before opening its parent")
	}
	if err := tt.Click("llm.plan"); err != nil {
		t.Fatal(err)
	}
	if id, _ := s.Tree.Selected(); id != "plan" || !tt.HasText("tokens") || !tt.HasText("1200") {
		t.Fatalf("attributes: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if id, _ := s.Tree.Selected(); id != "search" || !tt.HasText("http.get") {
		t.Fatalf("keyboard open: %q", tt.Texts())
	}
}

func TestTraceViewerSnapshot(t *testing.T) {
	spans := traceSample()
	snapshot(t, "trace_viewer", 820, 420, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *TraceViewerState {
			st := &TraceViewerState{}
			st.Tree.SetOpen("search", true)
			return st
		})
		ui.Column(c).Padding(16).Height(420).Children(func() {
			TraceViewer(c, *s, spans, TraceViewerOptions{}).Element.Grow(1)
		})
	})
}
