package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var graphNodes = []GraphNode{{"a", "Abbot", 0, 0}, {"b", "Prior", 10, 0}, {"c", "Cellarer", 10, 10}, {"d", "Novice", 0, 10}}
var graphEdges = []GraphEdge{{"a", "b"}, {"a", "c"}, {"c", "d"}}

func TestNetworkGraphReadoutAndKeys(t *testing.T) {
	var s NetworkGraphState
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return NetworkGraph(c, &s, graphNodes, graphEdges, NetworkGraphOptions{Frame: ChartFrame{Label: "House"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "House")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Abbot · 2 links: Prior, Cellarer" {
		t.Fatalf("readout %q", got)
	}
	x0, y0, _ := s.Position("a")
	tt.Key(ui.Shift, ui.KeyRight)
	if x, y, _ := s.Position("a"); x <= x0 || y != y0 {
		t.Fatalf("Shift+Right left the node at %v,%v", x, y)
	}
	tt.Key(0, ui.KeyEqual)
	if s.Zoom() < 1.19 {
		t.Fatalf("zoom %v", s.Zoom())
	}
	tt.Key(0, ui.Key0)
	if s.Zoom() != 1 {
		t.Fatalf("reset zoom %v", s.Zoom())
	}
}

func TestNetworkGraphDragNodeAndPan(t *testing.T) {
	var s NetworkGraphState
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return NetworkGraph(c, &s, graphNodes, graphEdges, NetworkGraphOptions{Frame: ChartFrame{Label: "House"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	// The 10×10 layout fits 240-48 DIPs, centered in the plot.
	cx, cy := b.X+b.W/2, b.Y+120
	ax, ay := cx-96, cy-96
	tt.Press(ax, ay)
	tt.Move(ax+20, ay)
	tt.Move(ax+48, ay)
	tt.Release(ax+48, ay)
	if x, y, _ := s.Position("a"); math.Abs(x-2.5) > 0.2 || math.Abs(y) > 0.01 {
		t.Fatalf("dragged node to %v,%v", x, y)
	}
	tt.Press(cx, cy) // the empty middle pans
	tt.Move(cx+30, cy+10)
	tt.Release(cx+30, cy+10)
	if x, _, _ := s.Position("b"); x != 10 {
		t.Fatalf("panning moved a node: %v", x)
	}
	if got := s.panX; got != 30 {
		t.Fatalf("pan %v", got)
	}
}

func TestNetworkGraphForceLayout(t *testing.T) {
	pos := networkLayout(graphNodes, [][2]int{{0, 1}, {0, 2}, {2, 3}}, true)
	again := networkLayout(graphNodes, [][2]int{{0, 1}, {0, 2}, {2, 3}}, true)
	for id, p := range pos {
		if again[id] != p || math.IsNaN(p[0]) {
			t.Fatalf("layout not deterministic for %s", id)
		}
	}
	d := func(a, b string) float64 { return math.Hypot(pos[a][0]-pos[b][0], pos[a][1]-pos[b][1]) }
	if d("c", "d") >= d("b", "d") {
		t.Fatalf("linked nodes not closer: %v vs %v", d("c", "d"), d("b", "d"))
	}
}

func TestNetworkGraphBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an unknown edge end")
		}
	}()
	var s NetworkGraphState
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NetworkGraph(c, &s, graphNodes, []GraphEdge{{"a", "z"}}, NetworkGraphOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestNetworkGraphSnapshot(t *testing.T) {
	snapshot(t, "network_graph", 520, 320, func(c *ui.Context) {
		var s NetworkGraphState
		ui.Column(c).Padding(16).Children(func() {
			NetworkGraph(c, &s, graphNodes, append(graphEdges, GraphEdge{"b", "d"}, GraphEdge{"b", "c"}), NetworkGraphOptions{Frame: ChartFrame{Label: "House"}, Force: true})
		})
	})
}
