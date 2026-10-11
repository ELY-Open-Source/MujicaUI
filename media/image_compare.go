package media

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ImageCompareOptions configure an ImageCompare.
type ImageCompareOptions struct {
	// BeforeLabel and AfterLabel caption the sides; empty uses defaults.
	BeforeLabel, AfterLabel string
	Disabled                bool
}

// ImageCompare shows before left of *split, after right.
func ImageCompare(c *ui.Context, before, after ui.ImageSource, split *float64, opts ImageCompareOptions) core.ChoiceResult {
	if *split < 0 || *split > 1 || math.IsNaN(*split) {
		panic("mujica: ImageCompare split outside 0..1")
	}
	bw, bh := data.ViewerSize(before)
	data.ViewerSize(after)
	k := core.Tokens(c)
	bl, al := opts.BeforeLabel, opts.AfterLabel
	if bl == "" {
		bl = core.Msg(c, "image_compare.before")
	}
	if al == "" {
		al = core.Msg(c, "image_compare.after")
	}
	g, _ := input.SliderSetup(0, 1, 0.01, nil, false, c)
	old := *split
	var s *ui.Element
	var next float64
	ui.Box(c).Disabled(opts.Disabled).Children(func() {
		s, next, _ = input.SliderControl(c, old, g, 0, 1)
	})
	s.FillWidth().AspectRatio(bw / bh).FocusRing(false).Cursor(ui.CursorResizeEW).Clip().Radius(theme.ControlRadius).
		Label(bl + " / " + al)
	internal.Ring(s, k.Focus, theme.ControlRadius)
	if !s.IsDisabled() {
		*split = next
	}
	changed := *split != old
	if changed {
		slog.Debug("mujica: image_compare split", "split", *split)
	}
	v, size := float32(*split), core.FontSize(c, theme.CaptionSize)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		x := internal.Snap(p, r.X+r.W*v)
		p.Image(after, r, ui.Cover)
		p.Clip(ui.Rect{X: r.X, Y: r.Y, W: x - r.X, H: r.H}, 0, func() { p.Image(before, r, ui.Cover) })
		compareTag(p, k, r.X+8, r.Y+8, bl, size)
		w, _ := p.MeasureText(1000, ui.Span{Text: al, Size: size})
		compareTag(p, k, r.X+r.W-w-20, r.Y+8, al, size)
		p.Fill(ui.Rect{X: x - 1, Y: r.Y, W: 2, H: r.H}, k.Surface, 0)
		cy := internal.Snap(p, r.Y+r.H/2)
		p.FillPath(internal.Diamond(x, cy, 22), k.Surface)
		p.FillPath(internal.Diamond(x, cy, 16), k.AccentText)
	})
	return core.NewChoiceResult(s, changed)
}

// compareTag draws a caption chip at (x, y).
func compareTag(p *ui.Painter, k theme.Tokens, x, y float32, s string, size float32) {
	w, h := p.MeasureText(1000, ui.Span{Text: s, Size: size})
	b := ui.Rect{X: x, Y: y, W: w + 12, H: h + 4}
	p.Fill(b, k.Surface.Alpha(0.88), theme.ControlRadius)
	p.Text(b.X+6, b.Y+2, s, size, k.Text)
}
