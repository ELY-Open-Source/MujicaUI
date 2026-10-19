package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SankeyNode is a stage that flow passes through.
type SankeyNode struct{ ID, Label string }

// SankeyLink carries Value from node From to node To.
type SankeyLink struct {
	From, To string
	Value    float64
}

// SankeyChartOptions configure a SankeyChart.
type SankeyChartOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// sankeyLayout puts nodes in columns by longest path.
type sankeyLayout struct {
	col      []int
	in, out  []float64
	cols     int
	from, to []int
}

func sankeyPlace(nodes []SankeyNode, links []SankeyLink) sankeyLayout {
	at := map[string]int{}
	for i, n := range nodes {
		if n.ID == "" || n.Label == "" {
			panic("mujica: SankeyChart nodes need an ID and a Label")
		}
		if _, dup := at[n.ID]; dup {
			panic("mujica: SankeyChart has two nodes " + n.ID)
		}
		at[n.ID] = i
	}
	l := sankeyLayout{col: make([]int, len(nodes)), in: make([]float64, len(nodes)), out: make([]float64, len(nodes))}
	for _, e := range links {
		f, okf := at[e.From]
		t, okt := at[e.To]
		if !okf || !okt || f == t || !(e.Value > 0) || math.IsInf(e.Value, 0) {
			panic("mujica: SankeyChart links need known ends and a positive value")
		}
		l.from, l.to = append(l.from, f), append(l.to, t)
		l.out[f] += e.Value
		l.in[t] += e.Value
	}
	// Relax columns; more passes than nodes means a cycle.
	for pass := 0; ; pass++ {
		if pass > len(nodes) {
			panic("mujica: SankeyChart links form a cycle")
		}
		moved := false
		for i := range links {
			if l.col[l.to[i]] < l.col[l.from[i]]+1 {
				l.col[l.to[i]], moved = l.col[l.from[i]]+1, true
			}
		}
		if !moved {
			break
		}
	}
	for i := range nodes {
		l.cols = max(l.cols, l.col[i]+1)
	}
	return l
}

// SankeyChart draws flows; reading a mark highlights its paths.
func SankeyChart(c *ui.Context, nodes []SankeyNode, links []SankeyLink, opts SankeyChartOptions) ChartResult {
	ch := NewChart(c, opts.Frame)
	l := ChartMemo(ch, "layout", [2]string{ChartID(nodes), ChartID(links)}, func() sankeyLayout { return sankeyPlace(nodes, links) })
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	val := func(v float64) string { return ChartValue(c, v, opts.Format) }
	size1 := func(i int) float64 { return max(l.in[i], l.out[i]) }
	// geom returns node boxes and link ends.
	geom := func(r ui.Rect) (boxes []ui.Rect, y0, y1 []float32, scale float32) {
		const gap, w = 10, 12
		plot := ui.Rect{X: r.X + 8, Y: r.Y + 6, W: r.W - 16, H: r.H - 12}
		most, count := 0.0, 0
		for col := range l.cols {
			sum, n := 0.0, 0
			for i := range nodes {
				if l.col[i] == col {
					sum, n = sum+size1(i), n+1
				}
			}
			if sum > most {
				most, count = sum, n
			}
		}
		scale = (plot.H - gap*float32(max(count-1, 0))) / float32(max(most, 1e-9))
		boxes = make([]ui.Rect, len(nodes))
		step := (plot.W - w) / float32(max(l.cols-1, 1))
		next := make([]float32, l.cols)
		for i := range nodes {
			h := max(float32(size1(i))*scale, 2)
			boxes[i] = ui.Rect{X: plot.X + step*float32(l.col[i]), Y: plot.Y + next[l.col[i]], W: w, H: h}
			next[l.col[i]] += h + gap
		}
		outAt, inAt := make([]float32, len(nodes)), make([]float32, len(nodes))
		for i := range links {
			y0, y1 = append(y0, boxes[l.from[i]].Y+outAt[l.from[i]]), append(y1, boxes[l.to[i]].Y+inAt[l.to[i]])
			outAt[l.from[i]] += float32(links[i].Value) * scale
			inAt[l.to[i]] += float32(links[i].Value) * scale
		}
		return
	}
	// ribbon returns link i's top edge at t and thickness.
	ribbon := func(boxes []ui.Rect, y0, y1 []float32, scale float32, i int, t float32) (x, y, h float32) {
		a, b := boxes[l.from[i]], boxes[l.to[i]]
		s := t * t * (3 - 2*t)
		return a.X + a.W + (b.X-a.X-a.W)*t, y0[i] + (y1[i]-y0[i])*s, float32(links[i].Value) * scale
	}
	n := len(nodes) + len(links)
	if len(links) == 0 {
		n = 0
	}
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		if i < len(nodes) {
			return nodes[i].Label + " · " + core.Msg(c, "sankey_chart.in") + " " + val(l.in[i]) + " · " + core.Msg(c, "sankey_chart.out") + " " + val(l.out[i])
		}
		e := links[i-len(nodes)]
		return nodes[l.from[i-len(nodes)]].Label + " → " + nodes[l.to[i-len(nodes)]].Label + " · " + val(e.Value)
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		boxes, y0, y1, scale := geom(r)
		for i, b := range boxes {
			if (ui.Rect{X: b.X - 4, Y: b.Y, W: b.W + 8, H: b.H}).Contains(x, y) {
				return i
			}
		}
		for i := range links {
			a, b := boxes[l.from[i]], boxes[l.to[i]]
			if x < a.X+a.W || x > b.X {
				continue
			}
			_, top, h := ribbon(boxes, y0, y1, scale, i, (x-a.X-a.W)/(b.X-a.X-a.W))
			if y >= top && y <= top+max(h, 4) {
				return len(nodes) + i
			}
		}
		return -1
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		boxes, y0, y1, scale := geom(r)
		lit := sankeyLit(l, len(nodes), active)
		for i := range links {
			var band ui.Path
			const steps = 24
			for s := 0; s <= steps; s++ {
				x, y, _ := ribbon(boxes, y0, y1, scale, i, float32(s)/steps)
				if s == 0 {
					band.MoveTo(x, y)
				} else {
					band.LineTo(x, y)
				}
			}
			for s := steps; s >= 0; s-- {
				x, y, h := ribbon(boxes, y0, y1, scale, i, float32(s)/steps)
				band.LineTo(x, y+h)
			}
			band.Close()
			alpha := float32(0.4)
			if lit != nil && !lit[len(nodes)+i] {
				alpha = 0.12
			}
			p.FillPath(&band, pal[l.from[i]%len(pal)].Alpha(alpha))
			if active == len(nodes)+i {
				p.StrokePath(&band, 1.5, k.Text)
			}
		}
		for i, b := range boxes {
			col := pal[i%len(pal)]
			if lit != nil && !lit[i] {
				col = col.Mix(k.Surface, 0.6)
			}
			p.Fill(b, col, 1)
			if i == active {
				p.Stroke(ui.Rect{X: b.X - 2, Y: b.Y - 2, W: b.W + 4, H: b.H + 4}, k.Text, 2, 1.5)
			}
			if l.col[i] == l.cols-1 {
				ChartText(p, nodes[i].Label, b.X-6, b.Y+b.H/2, size, k.Text, 1, 0.5)
			} else {
				ChartText(p, nodes[i].Label, b.X+b.W+6, b.Y+b.H/2, size, k.Text, 0, 0.5)
			}
		}
	}
	res, _ := ch.Build(spec)
	return res
}

// sankeyLit marks everything on paths through active; nil means all.
func sankeyLit(l sankeyLayout, nodes, active int) []bool {
	if active < 0 {
		return nil
	}
	lit := make([]bool, nodes+len(l.from))
	var walk func(node int, down bool)
	walk = func(node int, down bool) {
		lit[node] = true
		for i := range l.from {
			if down && l.from[i] == node && !lit[nodes+i] {
				lit[nodes+i] = true
				walk(l.to[i], true)
			}
			if !down && l.to[i] == node && !lit[nodes+i] {
				lit[nodes+i] = true
				walk(l.from[i], false)
			}
		}
	}
	if active < nodes {
		walk(active, true)
		walk(active, false)
		return lit
	}
	i := active - nodes
	lit[active] = true
	walk(l.to[i], true)
	walk(l.from[i], false)
	return lit
}
