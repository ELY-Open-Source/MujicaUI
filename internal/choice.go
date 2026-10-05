package internal

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CheckState is what a drawn check box shows.
type CheckState int

const (
	Unchecked CheckState = iota
	Checked
	Mixed
)

// CheckBox paints a square check box filling r: a tick, a dash for mixed.
func CheckBox(p *ui.Painter, r ui.Rect, s CheckState, k theme.Tokens, hover bool) {
	if s == Unchecked {
		bg := k.Surface
		if hover {
			bg = k.SurfaceHover
		}
		p.Fill(r, bg, theme.ControlRadius)
		p.Stroke(r, k.ControlBorder, theme.ControlRadius, theme.BorderWidth)
		return
	}
	bg := k.Accent
	if hover {
		bg = k.AccentHover
	}
	p.Fill(r, bg, theme.ControlRadius)
	mark := "check"
	if s == Mixed {
		mark = "minus"
	}
	in := r.W * 0.12
	Mark(p, ui.Rect{X: r.X + in, Y: r.Y + in, W: r.W - 2*in, H: r.H - 2*in}, mark, k.OnAccent)
}

// Mark paints the Lucide icon name in r, on the pixel grid so masks never vary.
func Mark(p *ui.Painter, r ui.Rect, name string, c ui.Color) {
	p.Icon(icons.Must(name), r, c)
}
