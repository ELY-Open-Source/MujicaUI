package chart

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DrawingToolbarOptions configure a DrawingToolbar.
type DrawingToolbarOptions struct {
	// Horizontal uses a row; the default uses a column.
	Horizontal bool
}

// DrawingToolbar selects tools and manages lock, clear, and undo.
func DrawingToolbar(c *ui.Context, d *DrawingTools, opts DrawingToolbarOptions) core.ChoiceResult {
	if d == nil {
		panic("mujica: DrawingToolbar needs drawings")
	}
	if d.Tool < ToolSelect || d.Tool > ToolText {
		panic(fmt.Sprintf("mujica: unknown DrawingTool %d", d.Tool))
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	changed := false
	bar := ui.Column(c).Gap(u).Padding(u).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Role(ui.RoleToolbar).Label(core.Msg(c, "drawing_toolbar.label"))
	if opts.Horizontal {
		bar.Row().FocusGroup(ui.Horizontal)
	} else {
		bar.FocusGroup(ui.Vertical)
	}
	pick := func(t DrawingTool) {
		if !d.Locked && d.Tool != t && !bar.IsDisabled() {
			d.Tool, changed = t, true
			d.pending, d.moving = nil, false
			slog.Debug("mujica: drawing_toolbar tool", "tool", t)
		}
	}
	digits := []ui.Key{ui.Key1, ui.Key2, ui.Key3, ui.Key4, ui.Key5, ui.Key6}
	for t := ToolSelect; t <= ToolText; t++ {
		if c.Shortcut(ui.Alt, digits[t]) {
			pick(t)
		}
	}
	bar.Children(func() {
		for t := ToolSelect; t <= ToolText; t++ {
			name := fmt.Sprintf("%s (Alt+%d)", core.Msg(c, fmt.Sprintf("drawing_toolbar.tool%d", t)), t+1)
			on := d.Tool == t
			b := ui.ButtonBase(c.Key(t)).Size(theme.MinHitSize+4, theme.MinHitSize+4).Radius(theme.ControlRadius).Label(name).Tooltip(name).Disabled(d.Locked)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			fg := k.Text
			switch {
			case on:
				b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText)
				fg = k.AccentText
			case b.Hovered():
				b.Background(k.SurfaceHover)
			}
			b.Draw(func(p *ui.Painter, r ui.Rect) { drawingToolGlyph(p, r, t, fg) })
			if b.Clicked() {
				pick(t)
			}
		}
		div := ui.Box(c).Background(k.Border).Role(ui.RoleNone)
		if opts.Horizontal {
			div.Size(theme.BorderWidth, u*6)
		} else {
			div.Size(u*6, theme.BorderWidth)
		}
		lockName := core.Msg(c, "drawing_toolbar.lock")
		if d.Locked {
			lockName = core.Msg(c, "drawing_toolbar.unlock")
		}
		lock := core.IconAction(c, icons.Must("lock"), lockName).Tooltip(lockName)
		if d.Locked {
			lock.Background(k.Selection).TextColor(k.AccentText)
		}
		if lock.Clicked() && !bar.IsDisabled() {
			d.Locked, changed = !d.Locked, true
			d.Tool = ToolSelect
			slog.Debug("mujica: drawing_toolbar lock", "locked", d.Locked)
		}
		undo := core.IconAction(c, icons.Must("rotate-ccw"), core.Msg(c, "drawing_toolbar.undo")).Tooltip(core.Msg(c, "drawing_toolbar.undo")).Disabled(d.Locked || !d.CanUndo())
		if undo.Clicked() && !bar.IsDisabled() && d.Undo() {
			changed = true
		}
		clear := core.IconAction(c, icons.Must("trash-2"), core.Msg(c, "drawing_toolbar.clear")).Tooltip(core.Msg(c, "drawing_toolbar.clear")).Disabled(d.Locked || len(d.Items) == 0)
		if clear.Clicked() && !bar.IsDisabled() && d.Clear() {
			changed = true
		}
	})
	return core.NewChoiceResult(bar, changed)
}

// drawingToolGlyph draws a small picture of tool t in r.
func drawingToolGlyph(p *ui.Painter, r ui.Rect, t DrawingTool, col ui.Color) {
	cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
	s := float32(8)
	switch t {
	case ToolSelect:
		var path ui.Path
		path.MoveTo(cx-5, cy-7).LineTo(cx+5, cy+1).LineTo(cx, cy+1).LineTo(cx+3, cy+7).LineTo(cx+1, cy+8).LineTo(cx-2, cy+2).LineTo(cx-5, cy+5).Close()
		p.FillPath(&path, col)
	case ToolTrendLine:
		drawingStroke(p, cx-s, cy+s, cx+s, cy-s, 1.5, col)
		p.FillPath(internal.Diamond(cx-s, cy+s, 5), col)
		p.FillPath(internal.Diamond(cx+s, cy-s, 5), col)
	case ToolHorizontal:
		p.Line(cx-s, cy, cx+s, cy, 1.5, col)
		p.FillPath(internal.Diamond(cx, cy, 5), col)
	case ToolFibonacci:
		for i, f := range []float32{-1, -0.4, 0.2, 1} {
			w := []float32{1.5, 1, 1, 1.5}[i]
			p.Line(cx-s, cy+f*s, cx+s, cy+f*s, w, col)
		}
	case ToolRectangle:
		p.Stroke(ui.Rect{X: cx - s, Y: cy - s*0.7, W: 2 * s, H: 1.4 * s}, col, 1, 1.5)
	case ToolText:
		p.Line(cx-s*0.8, cy-s, cx+s*0.8, cy-s, 2, col)
		p.Line(cx, cy-s, cx, cy+s, 2, col)
	}
}
