package chart

import (
	"fmt"
	"log/slog"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// GraphNode is a node placed at X, Y unless forced.
type GraphNode struct {
	ID, Label string
	X, Y      float64
}

// GraphEdge joins two nodes by ID.
type GraphEdge struct{ From, To string }

// NetworkGraphOptions configure a NetworkGraph.
type NetworkGraphOptions struct {
	Frame ChartFrame
	// Force lays nodes out by a simple force simulation.
	Force bool
}

// NetworkGraphState holds node positions, zoom and pan.
type NetworkGraphState struct {
	pos      map[string][2]float64
	lo, hi   [2]float64
	layout   string
	zoom     float32
	panX     float32
	panY     float32
	drag     string
	dragging bool
}

// Position returns where node id is, in layout units.
func (s *NetworkGraphState) Position(id string) (x, y float64, ok bool) {
	p, ok := s.pos[id]
	return p[0], p[1], ok
}

// Zoom returns the view's scale, 1 when fitted.
func (s *NetworkGraphState) Zoom() float32 { return max(s.zoom, 0) }

// ResetView fits the graph again.
func (s *NetworkGraphState) ResetView() { s.zoom, s.panX, s.panY = 1, 0, 0 }

// networkForceLimit bounds the nodes the O(n²) force layout takes.
const networkForceLimit = 500

// NetworkGraph draws nodes and edges; drag, pan and zoom.
func NetworkGraph(c *ui.Context, s *NetworkGraphState, nodes []GraphNode, edges []GraphEdge, opts NetworkGraphOptions) ChartResult {
	at := map[string]int{}
	for i, n := range nodes {
		if _, dup := at[n.ID]; dup || n.ID == "" || n.Label == "" {
			panic("mujica: NetworkGraph nodes need unique IDs and labels")
		}
		at[n.ID] = i
	}
	links := make([][2]int, len(edges))
	for i, e := range edges {
		f, okf := at[e.From]
		t, okt := at[e.To]
		if !okf || !okt {
			panic("mujica: NetworkGraph edge ends must be node IDs")
		}
		links[i] = [2]int{f, t}
	}
	if opts.Force && len(nodes) > networkForceLimit {
		panic("mujica: NetworkGraph Force takes at most 500 nodes")
	}
	if s.zoom == 0 {
		s.zoom = 1
	}
	if key := ChartID(nodes) + ChartID(edges) + fmt.Sprint(opts.Force); s.layout != key {
		s.layout, s.pos = key, networkLayout(nodes, links, opts.Force)
		// The view fits the layout; dragging does not refit.
		s.lo, s.hi = [2]float64{math.Inf(1), math.Inf(1)}, [2]float64{math.Inf(-1), math.Inf(-1)}
		for _, p := range s.pos {
			s.lo, s.hi = [2]float64{min(s.lo[0], p[0]), min(s.lo[1], p[1])}, [2]float64{max(s.hi[0], p[0]), max(s.hi[1], p[1])}
		}
		slog.Debug("mujica: network_graph layout", "nodes", len(nodes))
	}
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	lo, hi := s.lo, s.hi
	// screen maps a layout point into r, zoomed and panned.
	screen := func(r ui.Rect, p [2]float64) (float32, float32) {
		span := max(hi[0]-lo[0], hi[1]-lo[1], 1e-9)
		fit := float64(min(r.W, r.H)-48) / span
		cx, cy := (lo[0]+hi[0])/2, (lo[1]+hi[1])/2
		x := float32((p[0]-cx)*fit)*s.zoom + r.X + r.W/2 + s.panX
		y := float32((p[1]-cy)*fit)*s.zoom + r.Y + r.H/2 + s.panY
		return x, y
	}
	unit := func(r ui.Rect) float64 {
		return max(hi[0]-lo[0], hi[1]-lo[1], 1e-9) / float64(min(r.W, r.H)-48) / float64(s.zoom)
	}
	pt := func(r ui.Rect, i int) (float32, float32) { return screen(r, s.pos[nodes[i].ID]) }
	spec := ChartSpec{Marks: len(nodes)}
	spec.Readout = func(i int) string {
		var near []string
		for _, l := range links {
			switch i {
			case l[0]:
				near = append(near, nodes[l[1]].Label)
			case l[1]:
				near = append(near, nodes[l[0]].Label)
			}
		}
		return fmt.Sprintf(core.Msg(c, "network_graph.readout"), nodes[i].Label, len(near), strings.Join(near, ", "))
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		return internal.Nearest(len(nodes), func(i int) (float32, float32) { return pt(r, i) }, x, y, 12, false)
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		p.Clip(r, 0, func() {
			for _, l := range links {
				x0, y0 := pt(r, l[0])
				x1, y1 := pt(r, l[1])
				var line ui.Path
				line.MoveTo(internal.Snap(p, x0), internal.Snap(p, y0)).LineTo(internal.Snap(p, x1), internal.Snap(p, y1))
				w, col := float32(1.25), k.ControlBorder
				if active == l[0] || active == l[1] {
					w, col = 2.5, k.Text
				}
				p.StrokePath(&line, w, col)
			}
			for i, n := range nodes {
				x, y := pt(r, i)
				internal.Marker(p, 0, x, y, 14, k.Surface)
				internal.Marker(p, 0, x, y, 11, pal[i%len(pal)])
				if i == active {
					var ring ui.Path
					ring.Circle(internal.Snap(p, x), internal.Snap(p, y), 10)
					p.StrokePath(&ring, 2, k.Text)
				}
				ChartText(p, n.Label, x, y+9, size, k.Text, 0.5, 0)
			}
		})
	}
	res, plot := ch.Build(spec)
	if plot == nil || plot.IsDisabled() {
		return res
	}
	hit := func(x, y float32) string {
		b := plot.Bounds()
		if i := spec.Pick(ui.Rect{W: b.W, H: b.H}, x, y); i >= 0 {
			return nodes[i].ID
		}
		return ""
	}
	networkInput(plot, s, ch.st.focus, nodes, unit, hit)
	return res
}

// networkInput drags, pans, zooms and moves the focused node.
func networkInput(plot *ui.Element, s *NetworkGraphState, focus int, nodes []GraphNode, unit func(ui.Rect) float64, hit func(x, y float32) string) {
	b := plot.Bounds()
	r := ui.Rect{W: b.W, H: b.H}
	move := func(id string, dx, dy float32) {
		p := s.pos[id]
		s.pos[id] = [2]float64{p[0] + float64(dx)*unit(r), p[1] + float64(dy)*unit(r)}
	}
	zoom := func(f float32) { s.zoom = max(0.25, min(8, s.zoom*f)) }
	if dx, dy, held := plot.Dragged(); held {
		if !s.dragging {
			x, y, _ := plot.PointerPosition()
			s.dragging, s.drag = true, hit(x, y)
		}
		if s.drag != "" {
			move(s.drag, dx, dy)
		} else {
			s.panX, s.panY = s.panX+dx, s.panY+dy
		}
	} else if s.dragging {
		s.dragging = false
		slog.Debug("mujica: network_graph drag", "node", s.drag)
	}
	plot.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || ev.DY == 0 || plot.IsDisabled() {
			return false
		}
		if ev.DY < 0 {
			zoom(1.2)
		} else {
			zoom(1 / 1.2)
		}
		return true
	})
	switch {
	case plot.Shortcut(0, ui.KeyEqual), plot.Shortcut(ui.Shift, ui.KeyEqual):
		zoom(1.2)
	case plot.Shortcut(0, ui.KeyMinus):
		zoom(1 / 1.2)
	case plot.Shortcut(0, ui.Key0):
		s.ResetView()
	}
	if focus < 0 || focus >= len(nodes) {
		return
	}
	for key, d := range map[ui.Key][2]float32{ui.KeyLeft: {-12, 0}, ui.KeyRight: {12, 0}, ui.KeyUp: {0, -12}, ui.KeyDown: {0, 12}} {
		if plot.Shortcut(ui.Shift, key) {
			move(nodes[focus].ID, d[0], d[1])
		}
	}
}

// networkLayout uses X, Y, or fixed force steps.
func networkLayout(nodes []GraphNode, links [][2]int, force bool) map[string][2]float64 {
	pos := make([][2]float64, len(nodes))
	for i, n := range nodes {
		pos[i] = [2]float64{n.X, n.Y}
		if force {
			a := 2 * math.Pi * float64(i) / float64(len(nodes))
			pos[i] = [2]float64{math.Cos(a), math.Sin(a)}
		}
	}
	if force && len(nodes) > 1 {
		ideal := math.Sqrt(4.0 / float64(len(nodes)))
		for step := range 300 {
			heat := 0.1 * (1 - float64(step)/300)
			move := make([][2]float64, len(nodes))
			for i := range pos {
				for j := range pos {
					if i == j {
						continue
					}
					dx, dy := pos[i][0]-pos[j][0], pos[i][1]-pos[j][1]
					d := max(math.Hypot(dx, dy), 1e-6)
					f := ideal * ideal / d
					move[i][0], move[i][1] = move[i][0]+dx/d*f, move[i][1]+dy/d*f
				}
			}
			for _, l := range links {
				dx, dy := pos[l[0]][0]-pos[l[1]][0], pos[l[0]][1]-pos[l[1]][1]
				d := max(math.Hypot(dx, dy), 1e-6)
				f := d * d / ideal
				move[l[0]][0], move[l[0]][1] = move[l[0]][0]-dx/d*f, move[l[0]][1]-dy/d*f
				move[l[1]][0], move[l[1]][1] = move[l[1]][0]+dx/d*f, move[l[1]][1]+dy/d*f
			}
			for i, m := range move {
				d := max(math.Hypot(m[0], m[1]), 1e-9)
				pos[i][0] += m[0] / d * min(d, heat)
				pos[i][1] += m[1] / d * min(d, heat)
			}
		}
	}
	out := make(map[string][2]float64, len(nodes))
	for i, n := range nodes {
		out[n.ID] = pos[i]
	}
	return out
}
