package chart

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// TreemapOptions configure a Treemap.
type TreemapOptions struct {
	Frame  ChartFrame
	Format func(v float64) string
}

// Treemap tiles the children at *path; branches drill in.
func Treemap(c *ui.Context, path *[]string, root ChartNode, opts TreemapOptions) ChartResult {
	chartCheckTree(root)
	node := chartDescend(root, *path)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	size := core.FontSize(c, theme.CaptionSize)
	pal := internal.Palette(k)
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		res.value = chartPath(c, root.Name, path, wrap.IsDisabled())
		node = chartDescend(root, *path)
		var kids []ChartNode
		for _, ch := range node.Children {
			if ch.total() > 0 {
				kids = append(kids, ch)
			}
		}
		slices.SortStableFunc(kids, func(a, b ChartNode) int { return cmp.Compare(b.total(), a.total()) })
		sum := node.total()
		values := make([]float64, len(kids))
		for i, ch := range kids {
			values[i] = ch.total()
		}
		ch := NewChart(c, opts.Frame)
		tiles := func(r ui.Rect) []ui.Rect { return squarify(values, r) }
		spec := ChartSpec{Marks: len(kids)}
		spec.Readout = func(i int) string {
			s := fmt.Sprintf("%s · %s · %.1f%%", kids[i].Name, ChartValue(c, values[i], opts.Format), 100*values[i]/sum)
			if len(kids[i].Children) > 0 {
				s += " · " + core.Msg(c, "chart.opens")
			}
			return s
		}
		spec.Pick = func(r ui.Rect, x, y float32) int {
			for i, t := range tiles(r) {
				if t.Contains(x, y) {
					return i
				}
			}
			return -1
		}
		spec.Activate = func(i int) {
			if len(kids[i].Children) > 0 {
				*path = append(slices.Clone(*path), kids[i].Name)
				res.value = true
				slog.Debug("mujica: treemap drill", "path", strings.Join(*path, "/"))
			}
		}
		spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
			for i, t := range tiles(r) {
				bg := pal[i%len(pal)]
				p.Fill(t, bg, 2)
				p.Stroke(t, k.Surface, 2, 2)
				if i == active {
					p.Stroke(ui.Rect{X: t.X + 3, Y: t.Y + 3, W: t.W - 6, H: t.H - 6}, ChartInk(k, bg), 1, 2)
				}
				// Labels show only where tiles have room.
				if t.W >= 36 && t.H >= 22 {
					ink := ChartInk(k, bg)
					ChartText(p, chartFit(p, kids[i].Name, t.W-12, size), t.X+6, t.Y+5, size, ink, 0, 0)
					if t.H >= 40 {
						ChartText(p, chartFit(p, ChartValue(c, values[i], opts.Format), t.W-12, size), t.X+6, t.Y+22, size, ink, 0, 0)
					}
				}
			}
		}
		r, plot := ch.Build(spec)
		if plot.Valid() && !plot.IsDisabled() && plot.Shortcut(0, ui.KeyBackspace) && len(*path) > 0 {
			*path, res.value = (*path)[:len(*path)-1], true
		}
		res.Element, res.Point, res.changed, res.retry = r.Element, r.Point, r.changed, r.retry
	})
	res.Element = wrap
	return res
}
