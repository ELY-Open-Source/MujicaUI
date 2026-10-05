package layout

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SplitPaneSide names a pane of a SplitPane.
type SplitPaneSide int

const (
	SplitPaneNeither SplitPaneSide = iota
	SplitPaneFirst
	SplitPaneSecond
)

// SplitPaneState is the app-held layout of a SplitPane.
type SplitPaneState struct {
	// Size is the first pane's width, or height when vertical, in DIPs; zero starts at half.
	Size float32
	// Collapsed hides the Collapsible pane; Size is kept for when it opens again.
	Collapsed bool
}

// SplitPaneOptions configure a SplitPane.
type SplitPaneOptions struct {
	// Vertical puts the first pane above the second.
	Vertical bool
	// MinFirst and MinSecond bound the panes in DIPs; zero is 10 spacing units.
	MinFirst, MinSecond float32
	// Collapsible names the pane that may collapse.
	Collapsible SplitPaneSide
}

// SplitPaneResult is a SplitPane and its events in this frame.
type SplitPaneResult struct {
	Element *ui.Element
	Resized bool
	Toggled bool
}

type splitPaneDrag struct {
	held bool
	// want is where the pointer asks the divider to be, before clamping.
	want float32
}

// SplitPane lays two panes side by side with a divider the pointer drags and the arrow keys move.
func SplitPane(c *ui.Context, state *SplitPaneState, opts SplitPaneOptions, first, second func()) SplitPaneResult {
	if state == nil || first == nil || second == nil {
		panic("mujica: SplitPane needs a state and two panes")
	}
	if opts.MinFirst < 0 || opts.MinSecond < 0 {
		panic("mujica: SplitPane minimums must not be negative")
	}
	if opts.Collapsible < SplitPaneNeither || opts.Collapsible > SplitPaneSecond {
		panic("mujica: invalid SplitPane Collapsible")
	}
	if state.Collapsed && opts.Collapsible == SplitPaneNeither {
		panic("mujica: SplitPane collapsed without a Collapsible pane")
	}
	u := core.Density(c).Unit()
	lo, minSecond := opts.MinFirst, opts.MinSecond
	if lo == 0 {
		lo = 10 * u
	}
	if minSecond == 0 {
		minSecond = 10 * u
	}
	e := ui.Row(c).AlignItems(ui.Stretch)
	if opts.Vertical {
		e.Column()
	}
	var res SplitPaneResult
	e.Children(func() {
		// Sizes come from this inner box: the caller's border and padding stay outside it.
		in := ui.Row(c).Key("in").Grow(1).MinWidth(0).AlignItems(ui.Stretch)
		if opts.Vertical {
			in.Column().MinHeight(0)
		}
		splitPaneBuild(c, in, state, opts, lo, minSecond, &res, first, second)
	})
	if res.Resized || res.Toggled {
		c.Invalidate()
	}
	res.Element = e
	return res
}

// splitPaneBuild lays the panes, divider and drag handle out in the inner box in.
func splitPaneBuild(c *ui.Context, in *ui.Element, state *SplitPaneState, opts SplitPaneOptions, lo, minSecond float32, res *SplitPaneResult, first, second func()) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	was := in.Bounds()
	internal.Remeasure(c, in, was.W, was.H)
	total := was.W
	if opts.Vertical {
		total = was.H
	}
	hi := max(lo, total-theme.BorderWidth-minSecond)
	clamp := func(v float32) float32 { return max(lo, min(v, hi)) }
	if state.Size <= 0 && total > 0 {
		state.Size = clamp((total - theme.BorderWidth) / 2)
	} else if total > 0 {
		state.Size = clamp(state.Size)
	}
	resize := func(v float32) {
		if v = clamp(v); v != state.Size {
			state.Size = v
			res.Resized = true
		}
	}
	toggle := func() {
		state.Collapsed = !state.Collapsed
		res.Toggled = true
		slog.Debug("mujica: split_pane toggled", "collapsed", state.Collapsed)
	}
	firstShown := !(state.Collapsed && opts.Collapsible == SplitPaneFirst)
	secondShown := !(state.Collapsed && opts.Collapsible == SplitPaneSecond)
	along := func(el *ui.Element, v float32) {
		if opts.Vertical {
			el.Height(v)
		} else {
			el.Width(v)
		}
	}
	in.Children(func() {
		if firstShown {
			a := ui.Box(c).Key("first").Shrink(0).Clip()
			switch {
			case !secondShown:
				a.Grow(1).Shrink(1)
			case state.Size <= 0 && opts.Vertical:
				a.HeightPercent(50)
			case state.Size <= 0:
				a.WidthPercent(50)
			default:
				along(a, state.Size)
			}
			a.Children(first)
		}
		div := ui.Box(c).Key("divider").Shrink(0).Focusable().FocusRing(false).Role(ui.RoleSplitter).Label(core.Msg(c, "split_pane.divider"))
		along(div, theme.BorderWidth)
		if state.Collapsed {
			div.Description(core.Msg(c, "split_pane.collapsed"))
		} else if total > 0 {
			div.Description(fmt.Sprintf("%.0f%%", 100*state.Size/total))
		}
		splitPaneKeys(c, div, state, opts, lo, hi, u, resize, toggle)
		if secondShown {
			b := ui.Box(c).Key("second").Grow(1).Shrink(1).Clip()
			if opts.Vertical {
				b.MinHeight(0)
			} else {
				b.MinWidth(0)
			}
			b.Children(second)
		}
		at := state.Size
		switch {
		case !firstShown:
			at = 0
		case !secondShown:
			at = total - theme.BorderWidth
		}
		grip := 1.5 * u
		handle := ui.Box(c).Key("handle").Absolute()
		if opts.Vertical {
			handle.Left(0).Right(0).Top(at - (grip-1)/2).Height(grip).Cursor(ui.CursorResizeNS)
		} else {
			handle.Top(0).Bottom(0).Left(at - (grip-1)/2).Width(grip).Cursor(ui.CursorResizeEW)
		}
		drag := ui.Local(handle, "split_pane.drag", func() splitPaneDrag { return splitPaneDrag{} })
		dx, dy, held := handle.Dragged()
		if held {
			div.Focus()
			if !drag.held {
				drag.want = at
			}
			delta := dx
			if opts.Vertical {
				delta = dy
			}
			if delta != 0 {
				drag.want += delta
				splitPaneFollow(state, opts, drag.want, lo, minSecond, total, resize, toggle)
			}
		}
		drag.held = held
		if opts.Collapsible != SplitPaneNeither && handle.DoubleClicked() {
			toggle()
		}
		active := held || handle.Hovered()
		div.Draw(func(p *ui.Painter, r ui.Rect) {
			line := k.ControlBorder
			if active {
				line = k.AccentText
			}
			p.Fill(r, line, 0)
			if state.Collapsed {
				// A short tab marks the edge a collapsed pane opens from.
				if opts.Vertical {
					p.Fill(ui.Rect{X: r.X + r.W/2 - 12, Y: r.Y - 1, W: 24, H: 3}, line, 1.5)
				} else {
					p.Fill(ui.Rect{X: r.X - 1, Y: r.Y + r.H/2 - 12, W: 3, H: 24}, line, 1.5)
				}
			}
		})
		handle.Draw(func(p *ui.Painter, r ui.Rect) {
			if div.FocusVisible() {
				internal.FocusRing(p, r, 0, k.Focus)
			}
		})
	})
}

// splitPaneFollow moves the divider where the pointer wants it; past half its minimum, a collapsible pane closes.
func splitPaneFollow(state *SplitPaneState, opts SplitPaneOptions, want, lo, minSecond, total float32, resize func(float32), toggle func()) {
	var past bool
	switch opts.Collapsible {
	case SplitPaneFirst:
		past = want < lo/2
	case SplitPaneSecond:
		past = total-theme.BorderWidth-want < minSecond/2
	}
	if past != state.Collapsed {
		toggle()
	}
	if !state.Collapsed {
		resize(want)
	}
}

// splitPaneKeys moves the divider with the arrows, Home and End, and toggles the collapsible pane with Enter.
func splitPaneKeys(c *ui.Context, div *ui.Element, state *SplitPaneState, opts SplitPaneOptions, lo, hi, u float32, resize func(float32), toggle func()) {
	back, forth := ui.KeyLeft, ui.KeyRight
	if opts.Vertical {
		back, forth = ui.KeyUp, ui.KeyDown
	}
	step := 4 * u
	switch {
	case div.Shortcut(0, back):
		switch {
		case state.Collapsed && opts.Collapsible == SplitPaneSecond:
			toggle()
		case state.Collapsed:
		case opts.Collapsible == SplitPaneFirst && state.Size <= lo:
			toggle()
		default:
			resize(state.Size - step)
		}
	case div.Shortcut(0, forth):
		switch {
		case state.Collapsed && opts.Collapsible == SplitPaneFirst:
			toggle()
		case state.Collapsed:
		case opts.Collapsible == SplitPaneSecond && state.Size >= hi:
			toggle()
		default:
			resize(state.Size + step)
		}
	case div.Shortcut(0, ui.KeyHome) && !state.Collapsed:
		resize(lo)
	case div.Shortcut(0, ui.KeyEnd) && !state.Collapsed:
		resize(hi)
	case opts.Collapsible != SplitPaneNeither && div.Shortcut(0, ui.KeyEnter):
		toggle()
	}
}
