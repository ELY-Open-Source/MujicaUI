package internal

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
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
	Mark(p, r, mark, k.OnAccent)
}

// Mark paints the Lucide icon name inset in r. Icons sit on the pixel grid,
// so their cached masks never differ between frames or tests.
func Mark(p *ui.Painter, r ui.Rect, name string, c ui.Color) {
	in := r.W * 0.12
	p.Icon(icons.Must(name), ui.Rect{X: r.X + in, Y: r.Y + in, W: r.W - 2*in, H: r.H - 2*in}, c)
}
