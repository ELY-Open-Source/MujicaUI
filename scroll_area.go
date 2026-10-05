package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ScrollDirection is the way a ScrollArea scrolls. The zero value is vertical.
type ScrollDirection int

const (
	ScrollVertical ScrollDirection = iota
	ScrollHorizontal
	ScrollBoth
)

// ScrollAreaOptions configure a ScrollArea.
type ScrollAreaOptions struct {
	Direction ScrollDirection
	// Label names the region for assistive technology.
	Label string
}

// ScrollArea scrolls its children with the wheel, its scroll bars and the keyboard.
// state keeps the offset across frames: set X or Y to scroll there, math.MaxFloat32 for the end.
// A nil state keeps the offset by the area's place in the view.
// The area takes the focus while its content overflows, so the keys reach it.
func ScrollArea(c *ui.Context, state *ui.ScrollState, opts ScrollAreaOptions, children func()) *ui.Element {
	if children == nil {
		panic("mujica: ScrollArea needs children")
	}
	var e *ui.Element
	switch opts.Direction {
	case ScrollVertical:
		e = ui.Scroll(c)
	case ScrollHorizontal:
		e = ui.ScrollHorizontal(c)
	case ScrollBoth:
		e = ui.ScrollBoth(c)
	default:
		panic("mujica: invalid ScrollDirection")
	}
	if state == nil {
		state = ui.Local(e, "scroll_area.state", func() ui.ScrollState { return ui.ScrollState{} })
	}
	e.TrackScroll(state)
	if state.MaxX > 0 || state.MaxY > 0 {
		focus := Tokens(c).Focus
		// The ring goes inside: areas often fill a clipped pane edge to edge.
		e.Focusable().FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
			if e.FocusVisible() {
				in := theme.FocusWidth / 2
				p.Stroke(ui.Rect{X: r.X + in, Y: r.Y + in, W: r.W - 2*in, H: r.H - 2*in}, focus, theme.ControlRadius, theme.FocusWidth)
			}
		})
	}
	if opts.Label != "" {
		e.Label(opts.Label)
	}
	return e.Children(children)
}
