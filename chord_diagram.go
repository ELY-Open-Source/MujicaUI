package mujica

import (
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChordDiagramOptions configure a ChordDiagram.
type ChordDiagramOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// chordLayout is each group's arc and each pair's chord ends, as angles.
type chordLayout struct {
	arcs  [][2]float64
	sub   [][][2]float64
	pairs [][2]int
	total float64
}

// chordGap is the angle left between group arcs.
const chordGap = 0.04

func chordPlace(matrix [][]float64) chordLayout {
	n := len(matrix)
	l := chordLayout{arcs: make([][2]float64, n), sub: make([][][2]float64, n)}
	sums := make([]float64, n)
	for i, row := range matrix {
		if len(row) != n {
			panic("mujica: ChordDiagram matrix must be square")
		}
		for j, v := range row {
			if v < 0 || math.IsNaN(v) || math.IsInf(v, 0) {
				panic("mujica: ChordDiagram values must be finite and ≥ 0")
			}
			sums[i] += v
			if j > i && v+matrix[j][i] > 0 {
				l.pairs = append(l.pairs, [2]int{i, j})
			}
		}
		l.total += sums[i]
	}
	if l.total == 0 {
		return l
	}
	span := 2*math.Pi - chordGap*float64(n)
	a := 0.0
	for i := range n {
		l.arcs[i] = [2]float64{a, a + span*sums[i]/l.total}
		l.sub[i] = make([][2]float64, n)
		at := a
		for j := range n {
			w := span * matrix[i][j] / l.total
			l.sub[i][j] = [2]float64{at, at + w}
			at += w
		}
		a = l.arcs[i][1] + chordGap
	}
	return l
}

// ChordDiagram joins group arcs by chords; reading a group highlights them.
func ChordDiagram(c *ui.Context, names []string, matrix [][]float64, opts ChordDiagramOptions) ChartResult {
	if len(names) != len(matrix) {
		panic("mujica: ChordDiagram needs a name per matrix row")
	}
	ch := newChart(c, opts.Frame)
	l := chartMemo(ch, "chords", chartID(matrix), func() chordLayout { return chordPlace(matrix) })
	k := Tokens(c)
	pal := internal.Palette(k)
	size := FontSize(c, theme.CaptionSize)
	val := func(v float64) string { return chartValue(c, v, opts.Format) }
	n := len(names) + len(l.pairs)
	if l.total == 0 {
		n = 0
	}
	geom := func(r ui.Rect) (cx, cy, rad float32) {
		return r.X + r.W/2, r.Y + r.H/2, min(r.W, r.H)/2 - 30
	}
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		if i < len(names) {
			out, in := 0.0, 0.0
			for j := range names {
				out += matrix[i][j]
				in += matrix[j][i]
			}
			return names[i] + " · " + msg(c, "chord_diagram.out") + " " + val(out) + " · " + msg(c, "chord_diagram.in") + " " + val(in)
		}
		p := l.pairs[i-len(names)]
		return names[p[0]] + " → " + names[p[1]] + " " + val(matrix[p[0]][p[1]]) + " · " + names[p[1]] + " → " + names[p[0]] + " " + val(matrix[p[1]][p[0]])
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		cx, cy, rad := geom(r)
		d := float32(math.Hypot(float64(x-cx), float64(y-cy)))
		a := math.Atan2(float64(x-cx), float64(cy-y))
		if a < 0 {
			a += 2 * math.Pi
		}
		if d >= rad-2 && d <= rad+14 {
			for i, arc := range l.arcs {
				if a >= arc[0] && a < arc[1] {
					return i
				}
			}
		}
		return -1
	}
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		cx, cy, rad := geom(r)
		lit := func(pair [2]int) bool {
			switch {
			case active < 0:
				return true
			case active < len(names):
				return pair[0] == active || pair[1] == active
			}
			return l.pairs[active-len(names)] == pair
		}
		for pi, pair := range l.pairs {
			i, j := pair[0], pair[1]
			s, t := l.sub[i][j], l.sub[j][i]
			var path ui.Path
			internal.Arc(&path, cx, cy, rad, s[0], s[1], true)
			tx, ty := internal.Polar(cx, cy, rad, t[0])
			path.QuadTo(cx, cy, tx, ty)
			internal.Arc(&path, cx, cy, rad, t[0], t[1], false)
			sx, sy := internal.Polar(cx, cy, rad, s[0])
			path.QuadTo(cx, cy, sx, sy).Close()
			alpha := float32(0.45)
			if !lit(pair) {
				alpha = 0.08
			}
			p.FillPath(&path, pal[i%len(pal)].Alpha(alpha))
			if active == len(names)+pi {
				p.StrokePath(&path, 1.5, k.Text)
			}
		}
		for i, arc := range l.arcs {
			grow := float32(0)
			if i == active {
				grow = 4
			}
			p.FillPath(internal.Sector(cx, cy, rad+2, rad+12+grow, arc[0], arc[1]), pal[i%len(pal)])
			if arc[1]-arc[0] > 0.15 {
				lx, ly := internal.Polar(cx, cy, rad+20, (arc[0]+arc[1])/2)
				ax := 0.5 - 0.5*float32(math.Sin((arc[0]+arc[1])/2))
				chartText(p, names[i], lx, ly, size, k.Text, ax, 0.5+0.5*float32(math.Cos((arc[0]+arc[1])/2)))
			}
		}
	}
	res, _ := ch.build(spec)
	return res
}
