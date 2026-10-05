package internal

import (
	"github.com/egoist/mygo/ui"

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
	var path ui.Path
	if s == Mixed {
		path.MoveTo(r.X+r.W*0.25, r.Y+r.H/2).LineTo(r.X+r.W*0.75, r.Y+r.H/2)
	} else {
		path.MoveTo(r.X+r.W*0.22, r.Y+r.H*0.52).LineTo(r.X+r.W*0.42, r.Y+r.H*0.72).LineTo(r.X+r.W*0.78, r.Y+r.H*0.3)
	}
	p.StrokePath(&path, 1.75, k.OnAccent)
}

// Chevron strokes an arrow in r pointing down, or right when right.
func Chevron(p *ui.Painter, r ui.Rect, right bool, c ui.Color) {
	var path ui.Path
	if right {
		path.MoveTo(r.X+r.W*0.38, r.Y+r.H*0.22).LineTo(r.X+r.W*0.66, r.Y+r.H/2).LineTo(r.X+r.W*0.38, r.Y+r.H*0.78)
	} else {
		path.MoveTo(r.X+r.W*0.22, r.Y+r.H*0.38).LineTo(r.X+r.W/2, r.Y+r.H*0.66).LineTo(r.X+r.W*0.78, r.Y+r.H*0.38)
	}
	p.StrokePath(&path, 1.5, c)
}

// Cross strokes an x in r.
func Cross(p *ui.Painter, r ui.Rect, c ui.Color) {
	var path ui.Path
	k := r.W * 0.28
	path.MoveTo(r.X+k, r.Y+k).LineTo(r.X+r.W-k, r.Y+r.H-k).MoveTo(r.X+r.W-k, r.Y+k).LineTo(r.X+k, r.Y+r.H-k)
	p.StrokePath(&path, 1.5, c)
}
