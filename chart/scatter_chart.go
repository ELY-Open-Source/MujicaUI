package chart

import (
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ScatterPoint is one dot; Size sets a bubble's area.
type ScatterPoint struct{ X, Y, Size float64 }

// ScatterGroup is named points sharing color and shape.
type ScatterGroup struct {
	Name   string
	Points []ScatterPoint
}

// ScatterKey names point Point of group Group.
type ScatterKey struct{ Group, Point int }

// ScatterChartOptions configure a ScatterChart.
type ScatterChartOptions struct {
	Frame ChartFrame
	// Bubbles maps Size to dot area, 4–18 DIPs across.
	Bubbles bool
	Format  func(v float64) string
}

type scatterDrag struct {
	on             bool
	x0, y0, x1, y1 float32
}

// ScatterChart plots groups; box drag, Enter or Space selects.
func ScatterChart(c *ui.Context, selected *[]ScatterKey, groups []ScatterGroup, opts ScatterChartOptions) ChartResult {
	ch := NewChart(c, opts.Frame)
	k := core.Tokens(c)
	pal := internal.Palette(k)
	size := core.FontSize(c, theme.CaptionSize)
	keys := make([]ChartLegendItem, len(groups))
	var marks []ScatterKey
	b := [6]float64{math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)}
	for gi, g := range groups {
		keys[gi] = ChartLegendItem{Name: g.Name, Color: pal[gi%len(pal)], Shape: gi}
		if g.Name == "" || ch.Hidden(g.Name) {
			continue
		}
		for pi, pt := range g.Points {
			if math.IsNaN(pt.X) || math.IsNaN(pt.Y) || pt.Size < 0 {
				panic("mujica: ScatterChart points need finite X, Y and Size ≥ 0")
			}
			marks = append(marks, ScatterKey{gi, pi})
			b = [6]float64{min(b[0], pt.X), max(b[1], pt.X), min(b[2], pt.Y), max(b[3], pt.Y), min(b[4], pt.Size), max(b[5], pt.Size)}
		}
	}
	var xs, ys [3]float64
	if len(marks) > 0 {
		// A margin keeps dots off the plot edges.
		px, py := (b[1]-b[0])*0.06, (b[3]-b[2])*0.06
		xs[0], xs[1], xs[2] = internal.Nice(b[0]-px, b[1]+px, 4)
		ys[0], ys[1], ys[2] = internal.Nice(b[2]-py, b[3]+py, 4)
	}
	pt := func(m ScatterKey) ScatterPoint { return groups[m.Group].Points[m.Point] }
	geom := func(r ui.Rect) (ui.Rect, internal.Linear, internal.Linear) {
		plot := ChartInset(r, 52, 26)
		return plot, internal.Linear{D0: xs[0], D1: xs[1], R0: float64(plot.X), R1: float64(plot.X + plot.W)},
			internal.Linear{D0: ys[0], D1: ys[1], R0: float64(plot.Y + plot.H), R1: float64(plot.Y)}
	}
	radius := func(m ScatterKey) float32 {
		if !opts.Bubbles || b[5] <= 0 {
			return 4
		}
		return 2 + 7*float32(math.Sqrt(pt(m).Size/b[5]))
	}
	toggle := func(m ScatterKey) {
		if i := slices.Index(*selected, m); i >= 0 {
			*selected = slices.Delete(slices.Clone(*selected), i, i+1)
		} else {
			*selected = append(slices.Clone(*selected), m)
		}
		ch.Res.value = true
		slog.Debug("mujica: scatter_chart select", "count", len(*selected))
	}
	spec := ChartSpec{Marks: len(marks), Legend: keys, Hidden: ch.Toggles(len(groups) > 1)}
	spec.Readout = func(i int) string {
		m, p := marks[i], pt(marks[i])
		s := fmt.Sprintf("%s · x %s · y %s", groups[m.Group].Name, ChartValue(c, p.X, opts.Format), ChartValue(c, p.Y, opts.Format))
		if opts.Bubbles {
			s += " · " + core.Msg(c, "scatter_chart.size") + " " + ChartValue(c, p.Size, opts.Format)
		}
		if slices.Contains(*selected, m) {
			s += " · " + core.Msg(c, "scatter_chart.selected")
		}
		return s
	}
	spec.Pick = func(r ui.Rect, x, y float32) int {
		_, sx, sy := geom(r)
		return internal.Nearest(len(marks), func(i int) (float32, float32) { return sx.Map(pt(marks[i]).X), sy.Map(pt(marks[i]).Y) }, x, y, 14, false)
	}
	spec.Activate = func(i int) { toggle(marks[i]) }
	drag := ui.Local(ch.Res.Element, "mujica.scatter.drag", func() scatterDrag { return scatterDrag{} })
	box := *drag
	spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
		plot, sx, sy := geom(r)
		internal.AxisY(p, plot, sy, internal.Ticks(ys[0], ys[1], ys[2]), func(v float64) string { return internal.FormatTick(v, ys[2]) }, size, k)
		var tx []float32
		var tl []string
		for _, t := range internal.Ticks(xs[0], xs[1], xs[2]) {
			tx, tl = append(tx, sx.Map(t)), append(tl, internal.FormatTick(t, xs[2]))
		}
		internal.AxisX(p, plot, tx, tl, size, k)
		for i, m := range marks {
			x, y, rad := sx.Map(pt(m).X), sy.Map(pt(m).Y), radius(m)
			internal.Marker(p, m.Group, x, y, 2*rad, keys[m.Group].Color.Alpha(0.85))
			if slices.Contains(*selected, m) || i == active {
				var ring ui.Path
				ring.Circle(internal.Snap(p, x), internal.Snap(p, y), rad+3+2*float32(chartBool(i == active)))
				p.StrokePath(&ring, 1.5, k.Text)
			}
		}
		if box.on {
			p.StrokeDashed(scatterRect(box, r.X, r.Y), k.Text, 0, 1)
		}
	}
	res, plot := ch.Build(spec)
	if plot != nil && !plot.IsDisabled() {
		x, y, _ := plot.PointerPosition()
		if _, _, held := plot.Dragged(); held {
			if !drag.on {
				*drag = scatterDrag{on: true, x0: x, y0: y}
			}
			drag.x1, drag.y1 = x, y
		} else if drag.on {
			if pb := plot.Bounds(); max(core.Abs32(drag.x1-drag.x0), core.Abs32(drag.y1-drag.y0)) > 4 {
				_, sx, sy := geom(ui.Rect{W: pb.W, H: pb.H})
				in := scatterRect(*drag, 0, 0)
				var picked []ScatterKey
				for _, m := range marks {
					if in.Contains(sx.Map(pt(m).X), sy.Map(pt(m).Y)) {
						picked = append(picked, m)
					}
				}
				*selected, res.value = picked, true
				slog.Debug("mujica: scatter_chart box", "count", len(picked))
			}
			*drag = scatterDrag{}
		}
	}
	return res
}

func scatterRect(d scatterDrag, ox, oy float32) ui.Rect {
	return ui.Rect{X: ox + min(d.x0, d.x1), Y: oy + min(d.y0, d.y1), W: core.Abs32(d.x1 - d.x0), H: core.Abs32(d.y1 - d.y0)}
}
