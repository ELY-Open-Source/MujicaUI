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

// HeatNode is a sector or a stock: market cap sizes it, Change (percent) colors it.
type HeatNode struct {
	Name     string
	Cap      float64
	Change   float64
	Children []HeatNode
}

// MarketHeatmapOptions configure a MarketHeatmap.
type MarketHeatmapOptions struct {
	Frame ChartFrame
	// Scale is the change in percent at full color; zero means 3.
	Scale float64
}

type marketHeatTile struct {
	name        string
	cap, change float64
	branch      bool
}

// MarketHeatmap tiles sectors and stocks by market cap, colored by change; branches drill in.
func MarketHeatmap(c *ui.Context, path *[]string, root HeatNode, opts MarketHeatmapOptions) ChartResult {
	if opts.Scale < 0 {
		panic("mujica: MarketHeatmap Scale is negative")
	}
	scale := cmp.Or(opts.Scale, 3)
	marketHeatCheck(root)
	k := Tokens(c)
	u := Density(c).Unit()
	size := FontSize(c, theme.CaptionSize)
	var res ChartResult
	wrap := ui.Column(c).Gap(u * 2).MinWidth(0)
	wrap.Children(func() {
		res.value = chartPath(c, root.Name, path, wrap.IsDisabled())
		node := marketHeatDescend(root, *path)
		var tiles []marketHeatTile
		for _, ch := range node.Children {
			t := marketHeatTile{ch.Name, ch.total(), ch.change(), len(ch.Children) > 0}
			if t.cap > 0 {
				tiles = append(tiles, t)
			}
		}
		slices.SortStableFunc(tiles, func(a, b marketHeatTile) int { return cmp.Compare(b.cap, a.cap) })
		caps := make([]float64, len(tiles))
		for i, t := range tiles {
			caps[i] = t.cap
		}
		color := func(ch float64) ui.Color {
			return k.Surface.Mix(trendColor(c, trendSign(ch)), float32(0.15+0.85*min(math.Abs(ch)/scale, 1)))
		}
		ch := newChart(c, opts.Frame)
		spec := ChartSpec{Marks: len(tiles)}
		spec.Readout = func(i int) string {
			t := tiles[i]
			s := fmt.Sprintf("%s · %s %s%% · %s", t.name, trendWord(c, trendSign(t.change)), statisticSigned(t.change, 2), statisticFormat(t.cap, 0))
			if t.branch {
				s += " · " + msg(c, "chart.opens")
			}
			return s
		}
		spec.Pick = func(r ui.Rect, x, y float32) int {
			for i, t := range squarify(caps, r) {
				if t.Contains(x, y) {
					return i
				}
			}
			return -1
		}
		spec.Activate = func(i int) {
			if tiles[i].branch {
				*path = append(slices.Clone(*path), tiles[i].name)
				res.value = true
				slog.Debug("mujica: market_heatmap drill", "path", strings.Join(*path, "/"))
			}
		}
		spec.Draw = func(p *ui.Painter, r ui.Rect, active int) {
			for i, t := range squarify(caps, r) {
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

// marketHeatLegend is the color scale from −Scale to +Scale with signed labels.
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

// change is the node's own change, or its children's cap-weighted mean.
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
