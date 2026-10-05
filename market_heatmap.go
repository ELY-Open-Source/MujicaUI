package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

type HeatNode struct {
	Name     string
	Cap      float64
	Change   float64
	Children []HeatNode
}

// MarketHeatmapOptions configure a MarketHeatmap.
type MarketHeatmapOptions struct {
	Frame   ChartFrame
	Version int
	Scale   float64
}

type marketHeatMemo struct {
	id            string
	tiles         []marketHeatTile
	caps          []float64
	width, height float32
	rects         []ui.Rect
	visible       []int
	small         []marketHeatCell
}

type marketHeatTile struct {
	name        string
	cap, change float64
	branch      bool
}

func MarketHeatmap(c *ui.Context, path *[]string, root HeatNode, opts MarketHeatmapOptions) ChartResult {
	if opts.Scale < 0 || math.IsNaN(opts.Scale) || math.IsInf(opts.Scale, 0) || path == nil {
		panic("mujica: MarketHeatmap Scale is negative")
	}
	scale := cmp.Or(opts.Scale, 3)
	k := Tokens(c)
	u := Density(c).Unit()
	size := FontSize(c, theme.CaptionSize)
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		res.value = chartPath(c, root.Name, path, wrap.IsDisabled())
		node := marketHeatDescend(root, *path)
		memo := ui.Local(wrap, "mujica.heat", func() marketHeatMemo { return marketHeatMemo{} })
		if id := strings.Join(*path, "/") + "|" + chartID(node.Children) + fmt.Sprint(opts.Version); id != memo.id {
			marketHeatCheck(root)
			memo.id, memo.tiles, memo.caps = id, nil, nil
			memo.rects = nil
			for _, ch := range node.Children {
				t := marketHeatTile{ch.Name, ch.total(), ch.change(), len(ch.Children) > 0}
				if t.cap > 0 {
					memo.tiles = append(memo.tiles, t)
				}
			}
			slices.SortStableFunc(memo.tiles, func(a, b marketHeatTile) int { return cmp.Compare(b.cap, a.cap) })
			for _, t := range memo.tiles {
				memo.caps = append(memo.caps, t.cap)
			}
		}
		tiles, caps := memo.tiles, memo.caps
		color := func(ch float64) ui.Color {
			return k.Surface.Mix(trendColor(c, trendSign(ch)), float32(0.15+0.85*min(math.Abs(ch)/scale, 1)))
		}
		ch := newChart(c, opts.Frame)
		spec := chartSpec{n: len(tiles)}
		spec.readout = func(i int) string {
			t := tiles[i]
			s := fmt.Sprintf("%s · %s %s%% · %s", t.name, trendWord(c, trendSign(t.change)), statisticSigned(t.change, 2), statisticFormat(t.cap, 0))
			if t.branch {
				s += " · " + msg(c, "chart.opens")
			}
			return s
		}
		spec.pick = func(r ui.Rect, x, y float32) int {
			memo.layout(caps, r)
			for i, t := range memo.rects {
				if t.Contains(x, y) {
					return i
				}
			}
			return -1
		}
		spec.activate = func(i int) {
			if tiles[i].branch {
				*path = append(slices.Clone(*path), tiles[i].name)
				res.value = true
				slog.Debug("mujica: market_heatmap drill", "path", strings.Join(*path, "/"))
			}
		}
		spec.draw = func(p *ui.Painter, r ui.Rect, active int) {
			memo.layout(caps, r)
			for _, cell := range memo.small {
				p.Fill(ui.Rect{X: r.X + cell.x, Y: r.Y + cell.y, W: 4, H: 4}, color(cell.change), 0)
			}
			for _, i := range memo.visible {
				t := memo.rects[i]
				t.X, t.Y = t.X+r.X, t.Y+r.Y
				bg := color(tiles[i].change)
				p.Fill(t, bg, 2)
				p.Stroke(t, k.Surface, 2, 2)
				ink := chartInk(k, bg)
				if i == active {
					p.Stroke(ui.Rect{X: t.X + 3, Y: t.Y + 3, W: t.W - 6, H: t.H - 6}, ink, 1, 2)
				}
				// Signed percents say the direction without color.
				if t.W >= 40 && t.H >= 34 {
					chartText(p, chartFit(p, tiles[i].name, t.W-12, size), t.X+t.W/2, t.Y+t.H/2-1, size, ink, 0.5, 1)
					chartText(p, statisticSigned(tiles[i].change, 2)+"%", t.X+t.W/2, t.Y+t.H/2+1, size, ink, 0.5, 0)
				}
			}
		}
		r, plot := ch.build(spec)
		if plot != nil && !plot.IsDisabled() && plot.Shortcut(0, ui.KeyBackspace) && len(*path) > 0 {
			*path, res.value = (*path)[:len(*path)-1], true
		}
		res.Element, res.Point, res.changed, res.retry = r.Element, r.Point, r.changed, r.retry
		marketHeatLegend(c, scale, color)
	})
	res.Element = wrap
	return res
}

func marketHeatLegend(c *ui.Context, scale float64, color func(float64) ui.Color) {
	k := Tokens(c)
	size := FontSize(c, theme.CaptionSize)
	steps := []float64{-scale, -scale / 2, 0, scale / 2, scale}
	ui.Row(c).Gap(2).Children(func() {
		for _, s := range steps {
			ui.Column(c).Gap(2).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(44, 10).Background(color(s)).Radius(1)
				ui.Text(c, statisticSigned(s, 1)+"%").FontSize(size).TextColor(k.TextMuted).FontFeatures("tnum")
			})
		}
	})
}

func (n HeatNode) total() float64 {
	if len(n.Children) == 0 {
		return n.Cap
	}
	s := 0.0
	for _, ch := range n.Children {
		s += ch.total()
	}
	return s
}

func (n HeatNode) change() float64 {
	if len(n.Children) == 0 {
		return n.Change
	}
	sum, w := 0.0, 0.0
	for _, ch := range n.Children {
		sum += ch.change() * ch.total()
		w += ch.total()
	}
	if w == 0 {
		return 0
	}
	return sum / w
}

// marketHeatCheck panics on unnamed, duplicate-sibling, negative or non-finite nodes.
func marketHeatCheck(n HeatNode) {
	if n.Name == "" || n.Cap < 0 || math.IsNaN(n.Cap) || math.IsInf(n.Cap, 0) || math.IsNaN(n.Change) || math.IsInf(n.Change, 0) {
		panic("mujica: heat nodes need a name, a finite cap ≥ 0 and a finite change")
	}
	seen := map[string]bool{}
	for _, ch := range n.Children {
		if seen[ch.Name] {
			panic("mujica: heat node has two children named " + ch.Name)
		}
		seen[ch.Name] = true
		marketHeatCheck(ch)
	}
}

func marketHeatDescend(root HeatNode, path []string) HeatNode {
	n := root
	for _, name := range path {
		i := slices.IndexFunc(n.Children, func(ch HeatNode) bool { return ch.Name == name })
		if i < 0 {
			panic("mujica: heatmap path has no node " + name)
		}
		n = n.Children[i]
	}
	return n
}

type marketHeatCell struct {
	x, y   float32
	change float64
}

func (m *marketHeatMemo) layout(caps []float64, r ui.Rect) {
	if m.rects != nil && m.width == r.W && m.height == r.H {
		return
	}
	m.width, m.height = r.W, r.H
	m.rects = squarify(caps, ui.Rect{W: r.W, H: r.H})
	m.visible, m.small = nil, nil
	type aggregate struct{ weight, change float64 }
	bins := map[[2]int]aggregate{}
	for i, t := range m.rects {
		if t.W >= 4 && t.H >= 4 {
			m.visible = append(m.visible, i)
			continue
		}
		key := [2]int{int((t.X + t.W/2) / 4), int((t.Y + t.H/2) / 4)}
		a := bins[key]
		a.weight += caps[i]
		a.change += caps[i] * m.tiles[i].change
		bins[key] = a
	}
	for key, a := range bins {
		m.small = append(m.small, marketHeatCell{float32(key[0]) * 4, float32(key[1]) * 4, a.change / a.weight})
	}
	slices.SortFunc(m.small, func(a, b marketHeatCell) int {
		if a.y != b.y {
			return cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})
}
