package mujica

import (
	"fmt"
	"log/slog"
	"maps"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ParallelCoordinatesOptions configure a ParallelCoordinates chart.
type ParallelCoordinatesOptions struct {
	Frame ChartFrame
	// Names labels rows; nil uses row numbers.
	Names  []string
	Format func(v float64) string
}

// pcDrag is a brush being dragged on one axis.
type pcDrag struct {
	axis   int
	y0, y1 float32
}

// ParallelCoordinates draws rows across parallel axes; dragging an axis filters.
func ParallelCoordinates(c *ui.Context, filters *map[string][2]float64, axes []string, rows [][]float64, opts ParallelCoordinatesOptions) ChartResult {
	if filters == nil || len(axes) < 2 || opts.Names != nil && len(opts.Names) != len(rows) {
		panic("mujica: ParallelCoordinates needs filters, two or more axes and a name per row")
	}
	seen := map[string]bool{}
	for _, name := range axes {
		if name == "" || seen[name] {
			panic("mujica: ParallelCoordinates axes need unique names")
		}
		seen[name] = true
	}
	for name, f := range *filters {
		if !seen[name] || math.IsNaN(f[0]) || math.IsNaN(f[1]) || math.IsInf(f[0], 0) || math.IsInf(f[1], 0) || f[0] > f[1] {
			panic("mujica: ParallelCoordinates filters need known axes and finite ascending ranges")
		}
	}
	ch := newChart(c, opts.Frame)
	k := Tokens(c)
	size := FontSize(c, theme.CaptionSize)
	ink := internal.Palette(k)[0]
	type scales struct{ lo, hi []float64 }
	sc := chartMemo(ch, "scales", [2]string{chartID(rows), chartID(axes)}, func() scales {
		for _, r := range rows {
			if len(r) != len(axes) {
				panic("mujica: ParallelCoordinates rows need one value per axis")
			}
			for _, v := range r {
				if math.IsInf(v, 0) {
					panic("mujica: ParallelCoordinates values must be finite or NaN")
				}
			}
		}
		s := scales{lo: make([]float64, len(axes)), hi: make([]float64, len(axes))}
		for a := range axes {
			lo, hi := math.Inf(1), math.Inf(-1)
			for _, r := range rows {
				if !math.IsNaN(r[a]) {
					lo, hi = min(lo, r[a]), max(hi, r[a])
				}
			}
			if math.IsInf(lo, 1) {
				lo, hi = 0, 1
			}
			s.lo[a], s.hi[a], _ = internal.Nice(lo, hi, 4)
		}
		return s
	})
	keep := func(r []float64) bool {
		for a, name := range axes {
			if f, ok := (*filters)[name]; ok && (math.IsNaN(r[a]) || r[a] < f[0] || r[a] > f[1]) {
				return false
			}
		}
		return true
	}
	passing := chartMemo(ch, "pc_filter", [2]string{chartID(rows), fmt.Sprint(*filters)}, func() []bool {
		out := make([]bool, len(rows))
		for i, row := range rows {
			out[i] = keep(row)
		}
		return out
	})
	visible := ui.Local(ch.res.Element, "mujica.pc.visible", func() []int { return nil })
	geom := func(r ui.Rect) (ui.Rect, func(a int) float32, func(a int, v float64) float32) {
		plot := ui.Rect{X: r.X + 36, Y: r.Y + 24, W: r.W - 72, H: r.H - 32}
		ax := func(a int) float32 { return plot.X + plot.W*float32(a)/float32(len(axes)-1) }
		y := func(a int, v float64) float32 {
			return plot.Y + plot.H*float32(1-(v-sc.lo[a])/(sc.hi[a]-sc.lo[a]))
		}
		return plot, ax, y
	}
	name := func(i int) string {
		if opts.Names != nil {
			return opts.Names[i]
		}
		return "#" + chartValue(c, float64(i+1), nil)
	}
	spec := ChartSpec{Marks: len(rows)}
	spec.Readout = func(i int) string {
		parts := []string{name(i)}
		for a, n := range axes {
			parts = append(parts, n+" "+chartValue(c, rows[i][a], opts.Format))
		}
		if !passing[i] {
			parts = append(parts, msg(c, "parallel_coordinates.filtered"))
		}
		return strings.Join(parts, " · ")
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		_, ax, ys := geom(r)
		best, dist := -1, float32(8)
		for _, i := range *visible {
			row := rows[i]
			for a := 0; a+1 < len(axes); a++ {
				x0, x1 := ax(a), ax(a+1)
				if x < x0 || x > x1 || math.IsNaN(row[a]) || math.IsNaN(row[a+1]) {
					continue
				}
				t := (x - x0) / (x1 - x0)
				if d := abs32(ys(a, row[a]) + t*(ys(a+1, row[a+1])-ys(a, row[a])) - y); d < dist {
					best, dist = i, d
				}
			}
		}
		return best
	}
	drag := ui.Local(ch.res.Element, "mujica.pc.drag", func() pcDrag { return pcDrag{axis: -1} })
	keyAxis := ui.Local(ch.res.Element, "mujica.pc.axis", func() int { return 0 })
	spec.Input = func(plot *ui.Element, r ui.Rect) {
		if plot.IsDisabled() {
			*drag = pcDrag{axis: -1}
			return
		}
		if plot.Shortcut(ui.Alt, ui.KeyPageDown) {
			*keyAxis = (*keyAxis + 1) % len(axes)
		}
		if plot.Shortcut(ui.Alt, ui.KeyPageUp) {
			*keyAxis = (*keyAxis + len(axes) - 1) % len(axes)
		}
		*keyAxis = min(*keyAxis, len(axes)-1)
		a := *keyAxis
		f, ok := (*filters)[axes[a]]
		if !ok {
			f = [2]float64{sc.lo[a], sc.hi[a]}
		}
		before := f
		step := (sc.hi[a] - sc.lo[a]) / 100
		if plot.Shortcut(ui.Alt, ui.KeyUp) {
			f[0] = min(f[1], f[0]+step)
		}
		if plot.Shortcut(ui.Alt, ui.KeyDown) {
			f[0] = max(sc.lo[a], f[0]-step)
		}
		if plot.Shortcut(ui.Alt|ui.Shift, ui.KeyUp) {
			f[1] = min(sc.hi[a], f[1]+step)
		}
		if plot.Shortcut(ui.Alt|ui.Shift, ui.KeyDown) {
			f[1] = max(f[0], f[1]-step)
		}
		clear := plot.Shortcut(ui.Alt, ui.KeyDelete)
		if f != before || clear && ok {
			m := maps.Clone(*filters)
			if m == nil {
				m = map[string][2]float64{}
			}
			if clear {
				delete(m, axes[a])
			} else {
				m[axes[a]] = f
			}
			*filters, ch.res.value = m, true
			c.Invalidate()
			c.Announce(fmt.Sprintf(msg(c, "parallel_coordinates.brush"), axes[a], chartValue(c, f[0], opts.Format), chartValue(c, f[1], opts.Format)))
			slog.Debug("mujica: parallel_coordinates keyboard", "axis", axes[a])
		}
		box, ax, _ := geom(r)
		px, py, _ := plot.PointerPosition()
		_, _, held := plot.Dragged()
		switch {
		case held && drag.axis < 0:
			for a := range axes {
				if abs32(px-ax(a)) <= 12 {
					*drag = pcDrag{axis: a, y0: py, y1: py}
				}
			}
		case held:
			drag.y1 = max(box.Y, min(box.Y+box.H, py))
		case drag.axis >= 0:
			a := drag.axis
			m := maps.Clone(*filters)
			if m == nil {
				m = map[string][2]float64{}
			}
			if abs32(drag.y1-drag.y0) < 4 {
				delete(m, axes[a])
			} else {
				v := func(y float32) float64 {
					return sc.lo[a] + float64((box.Y+box.H-y)/box.H)*(sc.hi[a]-sc.lo[a])
				}
				m[axes[a]] = [2]float64{min(v(drag.y0), v(drag.y1)), max(v(drag.y0), v(drag.y1))}
			}
			*filters, ch.res.value = m, true
			c.Invalidate()
			*drag = pcDrag{axis: -1}
			slog.Debug("mujica: parallel_coordinates brush", "axis", axes[a], "filters", len(m))
		}
	}
	shownDrag := *drag
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, ax, ys := geom(r)
		*visible = chartVisibleIndices(len(rows), max(1, int(plot.H*2)), active)
		line := func(row []float64, w float32, col ui.Color) {
			var path ui.Path
			started := false
			for a := range axes {
				if math.IsNaN(row[a]) {
					started = false
					continue
				}
				x, y := internal.Snap(p, ax(a)), internal.Snap(p, ys(a, row[a]))
				if started {
					path.LineTo(x, y)
				} else {
					path.MoveTo(x, y)
				}
				started = true
			}
			p.StrokePath(&path, w, col)
		}
		for _, i := range *visible {
			if !passing[i] {
				line(rows[i], 1, k.Border)
			}
		}
		for _, i := range *visible {
			if passing[i] {
				line(rows[i], 1.25, ink)
			}
		}
		if active >= 0 {
			line(rows[active], 3, k.Text)
		}
		for a, n := range axes {
			x := internal.Snap(p, ax(a))
			p.Line(x, plot.Y, x, plot.Y+plot.H, 1.5, k.ControlBorder)
			if f, ok := (*filters)[n]; ok {
				y0, y1 := max(plot.Y, ys(a, f[1])), min(plot.Y+plot.H, ys(a, f[0]))
				p.Fill(ui.Rect{X: x - 6, Y: y0, W: 12, H: max(y1-y0, 1)}, k.Selection, 2)
				p.Stroke(ui.Rect{X: x - 6, Y: y0, W: 12, H: max(y1-y0, 1)}, k.AccentText, 2, 1.5)
			}
			chartText(p, n, x, plot.Y-8, size, k.Text, 0.5, 1)
			chartText(p, internal.FormatTick(sc.hi[a], (sc.hi[a]-sc.lo[a])/4), x+10, plot.Y, size, k.TextMuted, 0, 0)
			chartText(p, internal.FormatTick(sc.lo[a], (sc.hi[a]-sc.lo[a])/4), x+10, plot.Y+plot.H, size, k.TextMuted, 0, 1)
		}
		if shownDrag.axis >= 0 {
			x := ax(shownDrag.axis)
			p.Stroke(ui.Rect{X: x - 6, Y: min(shownDrag.y0, shownDrag.y1), W: 12, H: abs32(shownDrag.y1 - shownDrag.y0)}, k.AccentText, 2, 1.5)
		}
	}
	res, _ := ch.build(spec)
	if res.Element.IsDisabled() {
		*drag = pcDrag{axis: -1}
	}
	return res
}
