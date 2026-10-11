package code

import (
	"fmt"
	"hash/fnv"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// FlameNode is a profile frame; Value includes children.
type FlameNode struct {
	Name     string
	Value    float64
	Children []FlameNode
}

// FlamegraphState holds zoom and focus paths and the search.
type FlamegraphState struct {
	Zoom, Focus []int
	Query       string
	said        string
}

// FlamegraphOptions configure a Flamegraph.
type FlamegraphOptions struct {
	// Unit follows values in the readout, as "ms".
	Unit string
	// Status shows loading or failure in place of the graph.
	Status   data.DataStatus
	Error    string
	Disabled bool
}

// flameRow is a frame row's height, a full hit target.
const flameRow = theme.MinHitSize

// FlamegraphPalette colors frames at 3:1 on Surface.
func FlamegraphPalette(k theme.Tokens) []ui.Color {
	// Ornament alone is too faint on ivory.
	return []ui.Color{k.AccentText, k.Info, k.Success, k.Warning, k.Ornament.Mix(k.Text, 0.2), k.TextMuted}
}

// Flamegraph draws a zoomable, searchable profile.
func Flamegraph(c *ui.Context, root *FlameNode, s *FlamegraphState, opts FlamegraphOptions) core.ChoiceResult {
	if s == nil {
		panic("mujica: Flamegraph needs a state")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	var res core.ChoiceResult
	res.Element = CodePanel(c).Disabled(opts.Disabled)
	res.Element.Children(func() {
		CodeHeader(c, func() {
			input.SearchInput(c, &s.Query, input.SearchInputOptions{Placeholder: core.Msg(c, "flamegraph.search"), Disabled: opts.Disabled}).Element.Grow(1)
		})
		if root == nil || opts.Status != data.DataReady {
			ui.Column(c).Height(flameRow * 6).Children(func() {
				res = core.NewChoiceResult(res.Element, data.DataPlaceholder(c, opts.Status, opts.Error, core.Msg(c, "flamegraph.empty"), core.Msg(c, "flamegraph.empty")))
			})
			return
		}
		flameCheck(root)
		if flameAt(root, s.Zoom) == nil {
			s.Zoom = nil
		}
		if flameAt(root, s.Focus) == nil {
			s.Focus = nil
		}
		q := strings.ToLower(strings.TrimSpace(s.Query))
		depth := flameDepth(root)
		g := ui.Box(c).Height(float32(depth)*flameRow).Shrink(0).Margin(d.Unit()*2, d.Unit()*3, 0, d.Unit()*3).Focusable().Label(core.Msg(c, "flamegraph.label"))
		internal.Ring(g, k.Focus, theme.ControlRadius)
		hover := ui.Local(g, "flamegraph.hover", func() []int { return nil })
		zoom, focus := s.Zoom, s.Focus
		if !g.IsDisabled() {
			if flameInput(c, g, root, s, hover) {
				res = core.NewChoiceResult(res.Element, true)
			}
		}
		b := g.Bounds()
		g.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W != b.W {
				c.Invalidate()
			}
			flamePaint(p, r, k, root, zoom, focus, q, g.Focused(), core.FontSize(c, theme.CaptionSize))
		})
		readout := s.Focus
		if *hover != nil {
			readout = *hover
		}
		text := flameReadout(c, root, readout, opts.Unit)
		ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(d.Unit(), d.Unit()*3).
			Role(ui.RoleStatus).Shrink(0).SingleLine()
		if g.Focused() && text != s.said {
			s.said = text
			c.Announce(text)
		}
	})
	return res
}

// flameInput zooms, moves focus and tracks hover.
func flameInput(c *ui.Context, g *ui.Element, root *FlameNode, s *FlamegraphState, hover *[]int) bool {
	b := g.Bounds()
	x, y, over := g.PointerPosition()
	*hover = nil
	if over && b.W > 0 {
		*hover = flameHit(root, s.Zoom, x/b.W, int(y/flameRow))
	}
	changed := false
	zoomTo := func(path []int) {
		s.Zoom, s.Focus, changed = slices.Clone(path), slices.Clone(path), true
		slog.Debug("mujica: flamegraph zoom", "depth", len(path))
	}
	if g.Clicked() && *hover != nil {
		zoomTo(*hover)
	}
	f := s.Focus
	switch {
	case g.Shortcut(0, ui.KeyEnter):
		zoomTo(f)
	case g.Shortcut(0, ui.KeyEscape) && len(s.Zoom) > 0:
		zoomTo(s.Zoom[:len(s.Zoom)-1])
	case g.Shortcut(0, ui.KeyUp) && len(f) > 0:
		s.Focus, changed = f[:len(f)-1], true
	case g.Shortcut(0, ui.KeyDown):
		if n := flameAt(root, f); len(n.Children) > 0 {
			s.Focus, changed = append(slices.Clone(f), 0), true
		}
	case g.Shortcut(0, ui.KeyLeft) && len(f) > 0 && f[len(f)-1] > 0:
		s.Focus, changed = append(slices.Clone(f[:len(f)-1]), f[len(f)-1]-1), true
	case g.Shortcut(0, ui.KeyRight) && len(f) > 0 && f[len(f)-1] < len(flameAt(root, f[:len(f)-1]).Children)-1:
		s.Focus, changed = append(slices.Clone(f[:len(f)-1]), f[len(f)-1]+1), true
	}
	return changed
}

// flameCheck panics on impossible values.
func flameCheck(n *FlameNode) {
	sum := 0.0
	for i := range n.Children {
		flameCheck(&n.Children[i])
		sum += n.Children[i].Value
	}
	if n.Value < 0 || sum > n.Value*(1+1e-9) {
		panic("mujica: FlameNode " + n.Name + " has a negative value or children larger than itself")
	}
}

func flameAt(root *FlameNode, path []int) *FlameNode {
	n := root
	for _, i := range path {
		if i < 0 || i >= len(n.Children) {
			return nil
		}
		n = &n.Children[i]
	}
	return n
}

func flameDepth(n *FlameNode) int {
	d := 0
	for i := range n.Children {
		d = max(d, flameDepth(&n.Children[i]))
	}
	return d + 1
}

// flameSpan returns a frame's fractional left and width.
func flameSpan(root *FlameNode, zoom, path []int) (x, w float64) {
	if len(path) <= len(zoom) {
		return 0, 1
	}
	z := flameAt(root, zoom)
	x, w = 0, 1
	n := z
	for _, i := range path[len(zoom):] {
		for j := range i {
			x += w * n.Children[j].Value / max(n.Value, 1e-12)
		}
		w *= n.Children[i].Value / max(n.Value, 1e-12)
		n = &n.Children[i]
	}
	return x, w
}

// flameHit finds the frame at fx in row.
func flameHit(root *FlameNode, zoom []int, fx float32, row int) []int {
	if row < len(zoom) {
		return slices.Clone(zoom[:row])
	}
	path := slices.Clone(zoom)
	for len(path) < row {
		n := flameAt(root, path)
		next := -1
		for i := range n.Children {
			x, w := flameSpan(root, zoom, append(path, i))
			if float64(fx) >= x && float64(fx) < x+w {
				next = i
			}
		}
		if next < 0 {
			return nil
		}
		path = append(path, next)
	}
	if x, w := flameSpan(root, zoom, path); float64(fx) < x || float64(fx) >= x+w {
		return nil
	}
	return path
}

func flameColor(k theme.Tokens, name string) ui.Color {
	h := fnv.New32a()
	h.Write([]byte(name))
	pal := FlamegraphPalette(k)
	return pal[h.Sum32()%uint32(len(pal))]
}

// flamePaint draws visible frames with search and focus marks.
func flamePaint(p *ui.Painter, r ui.Rect, k theme.Tokens, root *FlameNode, zoom, focus []int, q string, focused bool, size float32) {
	var walk func(n *FlameNode, path []int)
	walk = func(n *FlameNode, path []int) {
		fx, fw := flameSpan(root, zoom, path)
		w := float32(fw) * r.W
		if w < 1 {
			return
		}
		box := ui.Rect{X: internal.Snap(p, r.X+float32(fx)*r.W), Y: r.Y + float32(len(path))*flameRow, W: internal.Snap(p, w) - 1, H: flameRow - 2}
		col := flameColor(k, n.Name)
		face := col.Mix(k.Surface, 0.78)
		match := q != "" && strings.Contains(strings.ToLower(n.Name), q)
		above := len(path) < len(zoom)
		if q != "" && !match || above {
			face, col = k.Surface, k.Border
		}
		p.Fill(box, face, 2)
		p.Stroke(box, col, 2, 1)
		if match {
			p.Stroke(box, k.AccentText, 2, 2)
		}
		if focused && slices.Equal(path, focus) {
			internal.FocusRing(p, box, 2, k.Focus)
		}
		if box.W > 24 {
			fg := k.Text
			if above || q != "" && !match {
				fg = k.TextMuted
			}
			p.Clip(box, 2, func() { p.Text(box.X+6, box.Y+(box.H-size*1.3)/2, n.Name, size, fg) })
		}
		for i := range n.Children {
			if len(path) < len(zoom) && i != zoom[len(path)] {
				continue
			}
			walk(&n.Children[i], append(slices.Clone(path), i))
		}
	}
	walk(root, nil)
}

// flameReadout names a frame's time and share.
func flameReadout(c *ui.Context, root *FlameNode, path []int, unit string) string {
	n := flameAt(root, path)
	if n == nil {
		return core.Msg(c, "flamegraph.hint")
	}
	share := 100 * n.Value / max(root.Value, 1e-12)
	return fmt.Sprintf(core.Msg(c, "flamegraph.readout"), n.Name, n.Value, unit, share)
}
