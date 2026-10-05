package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SunburstOptions configure a Sunburst.
type SunburstOptions struct {
	Frame ChartFrame
	// Depth is the number of rings; 0 means 3.
	Depth  int
	Format func(v float64) string
}

// sunArc is one drawn segment of a ring.
type sunArc struct {
	path   []string
	node   ChartNode
	ring   int
	a0, a1 float64
	color  int
}

// Sunburst draws descendants of *path as rings; branches drill in.
func Sunburst(c *ui.Context, path *[]string, root ChartNode, opts SunburstOptions) ChartResult {
	if opts.Depth < 0 {
		panic("mujica: Sunburst Depth must be ≥ 0")
	}
	depth := opts.Depth
	if depth == 0 {
		depth = 3
	}
	chartCheckTree(root)
	chartDescend(root, *path)
	k := Tokens(c)
	u := Density(c).Unit()
	size := FontSize(c, theme.CaptionSize)
	pal := internal.Palette(k)
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		res.value = chartPath(c, root.Name, path, wrap.IsDisabled())
		node := chartDescend(root, *path)
		sum := node.total()
		var arcs []sunArc
		var walk func(n ChartNode, at []string, ring int, a0, a1 float64, color int)
		walk = func(n ChartNode, at []string, ring int, a0, a1 float64, color int) {
			total, width := n.total(), a1-a0
			for i, ch := range n.Children {
				v := ch.total()
				if v <= 0 {
					continue
				}
				span := width * v / total
				col := color
				if ring == 0 {
					col = i
				}
				arc := sunArc{path: append(slices.Clone(at), ch.Name), node: ch, ring: ring, a0: a0, a1: a0 + span, color: col}
				arcs = append(arcs, arc)
				if ring+1 < depth {
					walk(ch, arc.path, ring+1, arc.a0, arc.a1, col)
				}
				a0 += span
			}
		}
		if sum > 0 {
			walk(node, nil, 0, 0, 2*math.Pi, 0)
		}
		geom := func(r ui.Rect) (cx, cy, hole, ring float32) {
			rad := min(r.W, r.H)/2 - 6
			hole = rad * 0.3
			return r.X + r.W/2, r.Y + r.H/2, hole, (rad - hole) / float32(depth)
		}
		ch := newChart(c, opts.Frame)
		spec := chartSpec{n: len(arcs)}
		spec.readout = func(i int) string {
			a := arcs[i]
			s := fmt.Sprintf("%s · %s · %.1f%%", strings.Join(a.path, " › "), chartValue(c, a.node.total(), opts.Format), 100*a.node.total()/sum)
			if len(a.node.Children) > 0 {
				s += " · " + msg(c, "chart.opens")
			}
			return s
		}
		spec.pick = func(r ui.Rect, x, y float32) int {
			cx, cy, hole, ring := geom(r)
			d := float32(math.Hypot(float64(x-cx), float64(y-cy)))
			a := math.Atan2(float64(x-cx), float64(cy-y))
			if a < 0 {
				a += 2 * math.Pi
			}
			for i, arc := range arcs {
				r0 := hole + ring*float32(arc.ring)
				if d >= r0 && d < r0+ring && a >= arc.a0 && a < arc.a1 {
					return i
				}
			}
			return -1
		}
		spec.activate = func(i int) {
			if len(arcs[i].node.Children) > 0 {
				*path = append(slices.Clone(*path), arcs[i].path...)
				res.value = true
				slog.Debug("mujica: sunburst drill", "path", strings.Join(*path, "/"))
			}
		}
		// color fades outer rings, keeping 3:1 against Surface.
		color := func(a sunArc) ui.Color {
			col := pal[a.color%len(pal)]
			if faded := col.Mix(k.Surface, 0.18*float32(a.ring)); theme.Contrast(faded, k.Surface) >= 3 {
				return faded
			}
			return col
		}
		spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
			cx, cy, hole, ring := geom(r)
			for i, a := range arcs {
				r0 := hole + ring*float32(a.ring)
				col := color(a)
				p.FillPath(internal.Sector(cx, cy, r0, r0+ring, a.a0, a.a1), col)
				p.StrokePath(internal.Sector(cx, cy, r0, r0+ring, a.a0, a.a1), 1.5, k.Surface)
				if i == active {
					p.StrokePath(internal.Sector(cx, cy, r0+2, r0+ring-2, a.a0, a.a1), 2, chartInk(k, col))
				}
			}
			// Labels go on top where arcs are long enough.
			for _, a := range arcs {
				r0 := hole + ring*float32(a.ring)
				col := color(a)
				if mid := float64(r0+ring/2) * (a.a1 - a.a0); mid > 44 {
					x, y := internal.Polar(cx, cy, r0+ring/2, (a.a0+a.a1)/2)
					chartText(p, chartFit(p, a.node.Name, min(float32(mid)*0.8, ring*1.25), size), x, y, size, chartInk(k, col), 0.5, 0.5)
				}
			}
			chartText(p, chartFit(p, node.Name, hole*2-8, size), cx, cy, size, k.TextMuted, 0.5, 1)
			chartText(p, chartValue(c, sum, opts.Format), cx, cy+2, FontSize(c, theme.BodySize), k.Text, 0.5, 0)
		}
		r, plot := ch.build(spec)
		if plot != nil && !plot.IsDisabled() && plot.Shortcut(0, ui.KeyBackspace) && len(*path) > 0 {
			*path, res.value = (*path)[:len(*path)-1], true
		}
		res.Point, res.changed, res.retry = r.Point, r.changed, r.retry
	})
	res.Element = wrap
	return res
}
