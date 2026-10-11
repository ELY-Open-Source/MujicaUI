package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChordDiagramOptions configure a ChordDiagram.
type ChordDiagramOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// chordLayout stores group arcs and relation endpoints.
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
	for _, row := range matrix {
		if len(row) != n {
			panic("mujica: ChordDiagram matrix must be square")
		}
	}
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
			if j >= i && v+matrix[j][i] > 0 {
				l.pairs = append(l.pairs, [2]int{i, j})
			}
		}
		l.total += sums[i]
	}
	if math.IsInf(l.total, 0) {
		panic("mujica: ChordDiagram total must be finite")
	}
	if l.total == 0 {
		return l
	}
	gap := min(chordGap, math.Pi/float64(n))
	span := 2*math.Pi - gap*float64(n)
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
		a = l.arcs[i][1] + gap
	}
	return l
}

// ChordDiagram connects group arcs and highlights read relationships.
func ChordDiagram(c *ui.Context, names []string, matrix [][]float64, opts ChordDiagramOptions) ChartResult {
	if len(names) != len(matrix) {
		panic("mujica: ChordDiagram needs a name per matrix row")
	}
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" || seen[name] {
			panic("mujica: ChordDiagram groups need unique names")
		}
		seen[name] = true
	}
	ch := NewChart(c, opts.Frame)
	l := ChartMemo(ch, "chords", ChartID(matrix), func() chordLayout { return chordPlace(matrix) })
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	val := func(v float64) string { return ChartValue(c, v, opts.Format) }
	n := len(names) + len(l.pairs)
	if l.total == 0 {
		n = 0
	}
	geom := func(r ui.Rect) (cx, cy, rad float32) {
		return r.X + r.W/2, r.Y + r.H/2, min(r.W, r.H)/2 - 30
	}
	visible := ui.Local(ch.Res.Element, "mujica.chord.visible", func() []int { return nil })
	spec := ChartSpec{Marks: n}
	spec.Readout = func(i int) string {
		if i < len(names) {
			out, in := 0.0, 0.0
			for j := range names {
				out += matrix[i][j]
				in += matrix[j][i]
			}
			return names[i] + " · " + core.Msg(c, "chord_diagram.out") + " " + val(out) + " · " + core.Msg(c, "chord_diagram.in") + " " + val(in)
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
		*visible = chartVisibleIndices(len(l.pairs), max(1, int(rad*4)), active-len(names))
		for _, pi := range *visible {
			pair := l.pairs[pi]
			i, j := pair[0], pair[1]
			s, t := l.sub[i][j], l.sub[j][i]
			var path ui.Path
			internal.Arc(&path, cx, cy, rad, s[0], s[1], true, p)
			tx, ty := internal.Polar(cx, cy, rad, t[0])
			path.QuadTo(internal.Snap(p, cx), internal.Snap(p, cy), internal.Snap(p, tx), internal.Snap(p, ty))
			internal.Arc(&path, cx, cy, rad, t[0], t[1], false, p)
			sx, sy := internal.Polar(cx, cy, rad, s[0])
			path.QuadTo(internal.Snap(p, cx), internal.Snap(p, cy), internal.Snap(p, sx), internal.Snap(p, sy)).Close()
			alpha := float32(0.45)
			if !lit(pair) {
				alpha = 0.08
			}
			p.FillPath(&path, pal[i%len(pal)].Alpha(alpha))
			if lit(pair) {
				p.StrokePath(&path, 1, pal[i%len(pal)])
			}
			if active == len(names)+pi {
				p.StrokePath(&path, 1.5, k.Text)
			}
		}
		for i, arc := range l.arcs {
			grow := float32(0)
			if i == active {
				grow = 4
			}
			p.FillPath(internal.Sector(cx, cy, rad+2, rad+12+grow, arc[0], arc[1], p), pal[i%len(pal)])
			if arc[1]-arc[0] > 0.15 {
				lx, ly := internal.Polar(cx, cy, rad+20, (arc[0]+arc[1])/2)
				ax := 0.5 - 0.5*float32(math.Sin((arc[0]+arc[1])/2))
				ChartText(p, names[i], lx, ly, size, k.Text, ax, 0.5+0.5*float32(math.Cos((arc[0]+arc[1])/2)))
			}
		}
	}
	res, _ := ch.Build(spec)
	return res
}
