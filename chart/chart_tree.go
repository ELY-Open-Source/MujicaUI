package chart

import (
	"log/slog"
	"math"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChartNode is a hierarchy node; branches sum their leaves.
type ChartNode struct {
	Name     string
	Value    float64
	Children []ChartNode
}

// total is the node's value, or its children's sum.
func (n ChartNode) total() float64 {
	if len(n.Children) == 0 {
		return n.Value
	}
	s := 0.0
	for _, ch := range n.Children {
		s += ch.total()
	}
	return s
}

// chartCheckTree panics on unnamed, duplicate-sibling or negative nodes.
func chartCheckTree(n ChartNode) {
	if n.Name == "" || n.Value < 0 || math.IsNaN(n.Value) || math.IsInf(n.Value, 0) {
		panic("mujica: chart nodes need a name and a finite value ≥ 0")
	}
	seen := map[string]bool{}
	for _, ch := range n.Children {
		if seen[ch.Name] {
			panic("mujica: chart node has two children named " + ch.Name)
		}
		seen[ch.Name] = true
		chartCheckTree(ch)
	}
}

// chartDescend returns the node at path, panicking if unknown.
func chartDescend(root ChartNode, path []string) ChartNode {
	n := root
	for _, name := range path {
		i := slices.IndexFunc(n.Children, func(ch ChartNode) bool { return ch.Name == name })
		if i < 0 {
			panic("mujica: chart path has no node " + name)
		}
		n = n.Children[i]
	}
	return n
}

// chartPath shows Back and the breadcrumb; reports a move up.
func chartPath(c *ui.Context, root string, path *[]string, disabled bool) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	moved := false
	ui.Row(c).Gap(u).AlignItems(ui.Center).Wrap().Children(func() {
		back := core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "chart.back")).Disabled(disabled || len(*path) == 0)
		if back.Clicked() && !back.IsDisabled() {
			*path, moved = (*path)[:len(*path)-1], true
		}
		names := append([]string{root}, *path...)
		for i, name := range names {
			if i > 0 {
				ui.Text(c, "›").TextColor(k.TextMuted)
			}
			if i == len(names)-1 {
				ui.Text(c, name).FontWeight(600).SingleLine()
				continue
			}
			crumb := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u).Radius(theme.ControlRadius).Disabled(disabled).TextColor(k.AccentText)
			internal.Ring(crumb, k.Focus, theme.ControlRadius)
			crumb.Children(func() { ui.Text(c, name).Underline().SingleLine() })
			if crumb.Clicked() && !crumb.IsDisabled() {
				*path, moved = slices.Clone((*path)[:i]), true
			}
		}
	})
	if moved {
		slog.Debug("mujica: chart path", "path", strings.Join(*path, "/"))
	}
	return moved
}

// chartFit shortens s with an ellipsis to fit w.
func chartFit(p *ui.Painter, s string, w, size float32) string {
	r := []rune(s)
	for len(r) > 0 {
		t := string(r)
		if len(r) < len([]rune(s)) {
			t += "…"
		}
		if tw, _ := p.MeasureText(0, ui.Span{Text: t, Size: size}); tw <= w {
			return t
		}
		r = r[:len(r)-1]
	}
	return ""
}

// squarify tiles descending values into near-square rectangles.
func squarify(values []float64, r ui.Rect) []ui.Rect {
	out := make([]ui.Rect, len(values))
	total := 0.0
	for _, v := range values {
		total += v
	}
	if total <= 0 {
		return out
	}
	scale := float64(r.W*r.H) / total
	for i := 0; i < len(values); {
		side := float64(min(r.W, r.H))
		j, best := i+1, math.Inf(1)
		for ; j <= len(values); j++ {
			worst := squarifyWorst(values[i:j], side, scale)
			if worst > best {
				break
			}
			best = worst
		}
		j--
		row := 0.0
		for _, v := range values[i:j] {
			row += v * scale
		}
		thick := float32(row / side)
		at := float32(0)
		for k := i; k < j; k++ {
			l := float32(values[k]*scale) / thick
			if r.W >= r.H {
				out[k] = ui.Rect{X: r.X, Y: r.Y + at, W: thick, H: l}
			} else {
				out[k] = ui.Rect{X: r.X + at, Y: r.Y, W: l, H: thick}
			}
			at += l
		}
		if r.W >= r.H {
			r.X, r.W = r.X+thick, r.W-thick
		} else {
			r.Y, r.H = r.Y+thick, r.H-thick
		}
		i = j
	}
	return out
}

// squarifyWorst is a row's worst aspect ratio along side.
func squarifyWorst(row []float64, side, scale float64) float64 {
	sum, lo, hi := 0.0, math.Inf(1), 0.0
	for _, v := range row {
		a := v * scale
		sum, lo, hi = sum+a, min(lo, a), max(hi, a)
	}
	if lo == 0 {
		return math.Inf(1)
	}
	s2 := side * side
	return max(s2*hi/(sum*sum), sum*sum/(s2*lo))
}
