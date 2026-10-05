package chart

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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

type MarketHeatMemo struct {
	id            string
	Tiles         []MarketHeatTile
	Caps          []float64
	width, height float32
	Rects         []ui.Rect
	Visible       []int
	Small         []marketHeatCell
}

type MarketHeatTile struct {
	Name        string
	Cap, Change float64
	branch      bool
}

func MarketHeatmap(c *ui.Context, path *[]string, root HeatNode, opts MarketHeatmapOptions) ChartResult {
	if opts.Scale < 0 || math.IsNaN(opts.Scale) || math.IsInf(opts.Scale, 0) || path == nil {
		panic("mujica: MarketHeatmap Scale is negative")
	}
	scale := cmp.Or(opts.Scale, 3)
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	size := core.FontSize(c, theme.CaptionSize)
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		res.value = chartPath(c, root.Name, path, wrap.IsDisabled())
		node := marketHeatDescend(root, *path)
		memo := ui.Local(wrap, "mujica.heat", func() MarketHeatMemo { return MarketHeatMemo{} })
		if id := strings.Join(*path, "/") + "|" + ChartID(node.Children) + fmt.Sprint(opts.Version); id != memo.id {
			marketHeatCheck(root)
			memo.id, memo.Tiles, memo.Caps = id, nil, nil
			memo.Rects = nil
			for _, ch := range node.Children {
				t := MarketHeatTile{ch.Name, ch.total(), ch.change(), len(ch.Children) > 0}
				if t.Cap > 0 {
					memo.Tiles = append(memo.Tiles, t)
				}
			}
			slices.SortStableFunc(memo.Tiles, func(a, b MarketHeatTile) int { return cmp.Compare(b.Cap, a.Cap) })
			for _, t := range memo.Tiles {
				memo.Caps = append(memo.Caps, t.Cap)
			}
		}
		tiles, caps := memo.Tiles, memo.Caps
		color := func(ch float64) ui.Color {
			return k.Surface.Mix(core.TrendColor(c, core.TrendSign(ch)), float32(0.15+0.85*min(math.Abs(ch)/scale, 1)))
		}
		ch := NewChart(c, opts.Frame)
		spec := ChartSpec{Marks: len(tiles)}
		spec.Readout = func(i int) string {
			t := tiles[i]
			s := fmt.Sprintf("%s · %s %s%% · %s", t.Name, core.TrendWord(c, core.TrendSign(t.Change)), data.StatisticSigned(t.Change, 2), data.StatisticFormat(t.Cap, 0))
			if t.branch {
				s += " · " + core.Msg(c, "chart.opens")
			}
			return s
		}
		spec.Pick = func(r ui.Rect, x, y float32) int {
			memo.Layout(caps, r)
			for i, t := range memo.Rects {
				if t.Contains(x, y) {
					return i
				}
			}
			return -1
		}
		spec.Activate = func(i int) {
			if tiles[i].branch {
				*path = append(slices.Clone(*path), tiles[i].Name)
				res.value = true
				slog.Debug("mujica: market_heatmap drill", "path", strings.Join(*path, "/"))
			}
		}
		spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
			memo.Layout(caps, r)
			for _, cell := range memo.Small {
				p.Fill(ui.Rect{X: r.X + cell.x, Y: r.Y + cell.y, W: 4, H: 4}, color(cell.change), 0)
			}
			for _, i := range memo.Visible {
				t := memo.Rects[i]
				t.X, t.Y = t.X+r.X, t.Y+r.Y
				bg := color(tiles[i].Change)
				p.Fill(t, bg, 2)
				p.Stroke(t, k.Surface, 2, 2)
				ink := ChartInk(k, bg)
				if i == active {
					p.Stroke(ui.Rect{X: t.X + 3, Y: t.Y + 3, W: t.W - 6, H: t.H - 6}, ink, 1, 2)
				}
				// Signed percents say the direction without color.
				if t.W >= 40 && t.H >= 34 {
					ChartText(p, chartFit(p, tiles[i].Name, t.W-12, size), t.X+t.W/2, t.Y+t.H/2-1, size, ink, 0.5, 1)
					ChartText(p, data.StatisticSigned(tiles[i].Change, 2)+"%", t.X+t.W/2, t.Y+t.H/2+1, size, ink, 0.5, 0)
				}
			}
		}
		r, plot := ch.Build(spec)
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
	k := core.Tokens(c)
	size := core.FontSize(c, theme.CaptionSize)
	steps := []float64{-scale, -scale / 2, 0, scale / 2, scale}
	ui.Row(c).Gap(2).Children(func() {
		for _, s := range steps {
			ui.Column(c).Gap(2).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(44, 10).Background(color(s)).Radius(1)
				ui.Text(c, data.StatisticSigned(s, 1)+"%").FontSize(size).TextColor(k.TextMuted).FontFeatures("tnum")
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

func (m *MarketHeatMemo) Layout(caps []float64, r ui.Rect) {
	if m.Rects != nil && m.width == r.W && m.height == r.H {
		return
	}
	m.width, m.height = r.W, r.H
	m.Rects = squarify(caps, ui.Rect{W: r.W, H: r.H})
	m.Visible, m.Small = nil, nil
	type aggregate struct{ weight, change float64 }
	bins := map[[2]int]aggregate{}
	for i, t := range m.Rects {
		if t.W >= 4 && t.H >= 4 {
			m.Visible = append(m.Visible, i)
			continue
		}
		key := [2]int{int((t.X + t.W/2) / 4), int((t.Y + t.H/2) / 4)}
		a := bins[key]
		a.weight += caps[i]
		a.change += caps[i] * m.Tiles[i].Change
		bins[key] = a
	}
	for key, a := range bins {
		m.Small = append(m.Small, marketHeatCell{float32(key[0]) * 4, float32(key[1]) * 4, a.change / a.weight})
	}
	slices.SortFunc(m.Small, func(a, b marketHeatCell) int {
		if a.y != b.y {
			return cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})
}
