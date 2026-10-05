package input

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ResizablePanel is one panel of a ResizablePanelGroup.
type ResizablePanel struct {
	// Label names the panel; dividers read it out.
	Label string
	// Min is the smallest size in DIPs; zero is 10 spacing units.
	Min         float32
	Collapsible bool
	Build       func()
}

// ResizablePanelGroupOptions configure a ResizablePanelGroup.
type ResizablePanelGroupOptions struct {
	// Vertical stacks the panels top to bottom.
	Vertical bool
	Disabled bool
}

// ResizablePanelGroupResult reports a group's changes: Resized while moving, Committed when a change is done, to save.
type ResizablePanelGroupResult struct {
	Element   *ui.Element
	Resized   bool
	Committed bool
}

// resizablePanelState keeps drags and the sizes collapsed panels return to.
type resizablePanelState struct {
	held     int // divider being dragged, -1 for none
	restored map[int]float32
}

// ResizablePanelGroup lays panels in a row with dividers; *sizes holds their fractions, summing to 1, for the caller to save.
func ResizablePanelGroup(c *ui.Context, sizes *[]float32, panels []ResizablePanel, opts ResizablePanelGroupOptions) ResizablePanelGroupResult {
	if sizes == nil || len(panels) < 2 {
		panic("mujica: ResizablePanelGroup needs sizes and two panels or more")
	}
	n := len(panels)
	if len(*sizes) == 0 {
		*sizes = make([]float32, n)
		for i := range *sizes {
			(*sizes)[i] = 1 / float32(n)
		}
	}
	if len(*sizes) != n {
		panic(fmt.Sprintf("mujica: ResizablePanelGroup has %d sizes for %d panels", len(*sizes), n))
	}
	var sum float32
	for i, f := range *sizes {
		if !(f >= 0) || f == 0 && !panels[i].Collapsible || panels[i].Label == "" || panels[i].Build == nil || panels[i].Min < 0 {
			panic(fmt.Sprintf("mujica: ResizablePanelGroup panel %d is invalid", i))
		}
		sum += f
	}
	if math.Abs(float64(sum-1)) > 1e-3 {
		panic(fmt.Sprintf("mujica: ResizablePanelGroup sizes sum to %v, not 1", sum))
	}
	u := core.Density(c).Unit()
	k := core.Tokens(c)
	var res ResizablePanelGroupResult
	box := ui.Row(c).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	if opts.Vertical {
		box.Column()
	}
	disabled := box.IsDisabled()
	st := *ui.Local(box, "resizable_panel_group", func() *resizablePanelState { return &resizablePanelState{held: -1, restored: map[int]float32{}} })
	was := box.Bounds()
	internal.Remeasure(c, box, was.W, was.H)
	total := was.W
	if opts.Vertical {
		total = was.H
	}
	avail := max(0, total-float32(n-1)*theme.BorderWidth)
	mins := make([]float32, n)
	for i, p := range panels {
		mins[i] = cmp.Or(p.Min, 10*u)
	}
	along := func(e *ui.Element, v float32) *ui.Element {
		if opts.Vertical {
			return e.Height(v).FillWidth()
		}
		return e.Width(v).FillHeight()
	}
	box.Children(func() {
		var grips []func()
		defer func() {
			// Grips go last, above the panels they overlap.
			for _, g := range grips {
				g()
			}
		}()
		for i, p := range panels {
			if i > 0 {
				grips = append(grips, resizablePanelDivider(c, k, u, st, sizes, panels, mins, avail, i-1, opts.Vertical, disabled, &res))
			}
			if (*sizes)[i] == 0 {
				continue
			}
			pane := ui.Box(c).Key(i).Shrink(0).Clip().Role(ui.RoleGroup).Label(p.Label)
			if avail > 0 {
				along(pane, (*sizes)[i]*avail)
			} else {
				pane.Grow((*sizes)[i])
			}
			pane.Children(p.Build)
		}
	})
	if res.Resized || res.Committed {
		c.Invalidate()
	}
	res.Element = box
	return res
}

// resizablePanelDivider builds the divider after panel i, stepped with arrows, Enter collapsing; it returns the grip's builder.
func resizablePanelDivider(c *ui.Context, k theme.Tokens, u float32, st *resizablePanelState, sizes *[]float32, panels []ResizablePanel, mins []float32, avail float32, i int, vertical, disabled bool, res *ResizablePanelGroupResult) func() {
	s := *sizes
	div := ui.Box(c).Key(fmt.Sprint("div", i)).Shrink(0).Focusable().FocusRing(false).Role(ui.RoleSplitter).
		Label(fmt.Sprintf(core.Msg(c, "resizable_panel_group.divider"), panels[i].Label, panels[i+1].Label))
	if vertical {
		div.Height(theme.BorderWidth).FillWidth()
	} else {
		div.Width(theme.BorderWidth).FillHeight()
	}
	// Moving the boundary by d DIPs trades size between panels i and i+1 within their minimums.
	move := func(d float32) bool {
		if avail <= 0 {
			return false
		}
		a, b := s[i]*avail, s[i+1]*avail
		pair := a + b
		na := max(0, min(pair, a+d))
		collapse := func(size, min float32, ok bool) bool { return ok && size < min/2 }
		switch {
		case collapse(na, mins[i], panels[i].Collapsible):
			na = 0
		case collapse(pair-na, mins[i+1], panels[i+1].Collapsible):
			na = pair
		default:
			na = max(min(mins[i], pair-mins[i+1]), min(na, pair-mins[i+1]))
		}
		if math.Abs(float64(na-a)) < 0.01 {
			return false
		}
		s[i], s[i+1] = na/avail, (pair-na)/avail
		return true
	}
	toggle := func() {
		j := i + 1
		if !panels[j].Collapsible {
			j = i
		}
		if !panels[j].Collapsible {
			return
		}
		other := i + i + 1 - j
		if s[j] > 0 {
			st.restored[j] = s[j]
			s[other] += s[j]
			s[j] = 0
		} else {
			back := min(cmp.Or(st.restored[j], mins[j]/max(avail, 1)), s[other]-mins[other]/max(avail, 1))
			s[j], s[other] = back, s[other]-back
		}
		res.Committed = true
		slog.Debug("mujica: resizable_panel_group toggled", "panel", panels[j].Label, "size", s[j])
	}
	if !disabled {
		back, forth := ui.KeyLeft, ui.KeyRight
		if vertical {
			back, forth = ui.KeyUp, ui.KeyDown
		}
		switch {
		case div.Shortcut(0, back):
			res.Committed = move(-4 * u)
		case div.Shortcut(0, forth):
			res.Committed = move(4 * u)
		case div.Shortcut(0, ui.KeyHome):
			res.Committed = move(-s[i] * avail)
		case div.Shortcut(0, ui.KeyEnd):
			res.Committed = move(s[i+1] * avail)
		case div.Shortcut(0, ui.KeyEnter):
			toggle()
		}
	}
	if s[i] > 0 && s[i+1] > 0 {
		div.Description(fmt.Sprintf("%.0f%%", 100*s[i]/(s[i]+s[i+1])))
	} else {
		div.Description(core.Msg(c, "split_pane.collapsed"))
	}
	return func() {
		grip := ui.Box(c).Key(fmt.Sprint("grip", i)).Absolute().Cursor(ui.CursorResizeEW)
		if vertical {
			grip.Cursor(ui.CursorResizeNS)
		}
		dx, dy, held := grip.Dragged()
		held = held && !disabled
		switch {
		case held:
			if st.held != i {
				st.held = i
				div.Focus()
			}
			d := dx
			if vertical {
				d = dy
			}
			if d != 0 && move(d) {
				res.Resized = true
			}
		case st.held == i:
			st.held = -1
			res.Committed = true
			slog.Debug("mujica: resizable_panel_group resized", "divider", i, "sizes", s)
		}
		resizablePanelGrip(c, k, u, s, avail, i, vertical, grip, div, st.held == i)
	}
}

// resizablePanelGrip places the grip over divider i and draws the divider.
func resizablePanelGrip(c *ui.Context, k theme.Tokens, u float32, s []float32, avail float32, i int, vertical bool, grip, div *ui.Element, held bool) {
	// The grip is wider than the 1 DIP line, centered on the boundary.
	var at float32
	for _, f := range s[:i+1] {
		at += f * avail
	}
	at += float32(i) * theme.BorderWidth
	if vertical {
		grip.Left(0).Right(0).Top(at - 1.5*u).Height(3 * u)
	} else {
		grip.Top(0).Bottom(0).Left(at - 1.5*u).Width(3 * u)
	}
	active := held || grip.Hovered()
	div.Draw(func(p *ui.Painter, r ui.Rect) {
		line := k.ControlBorder
		if active {
			line = k.AccentText
		}
		p.Fill(r, line, 0)
		if div.FocusVisible() {
			internal.FocusRing(p, r, 0, k.Focus)
		}
	})
}
