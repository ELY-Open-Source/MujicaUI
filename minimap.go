package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MinimapView is the editor's first line and line count.
type MinimapView struct{ Top, Visible int }

// MinimapOptions configure a Minimap.
type MinimapOptions struct {
	// Language colors the strokes as CodeViewer does.
	Language string
	// Width is the map's width in DIPs; zero is 96.
	Width    float32
	Disabled bool
}

// Line height and character width on the map.
const (
	minimapLine float32 = 2
	minimapChar float32 = 1
)

// Minimap draws code small and moves the visible frame.
func Minimap(c *ui.Context, lines []string, view *MinimapView, opts MinimapOptions) ChoiceResult {
	if view == nil {
		panic("mujica: Minimap needs a view")
	}
	if opts.Width < 0 || view.Visible < 0 || view.Top < 0 {
		panic("mujica: Minimap width, Top and Visible must not be negative")
	}
	k := Tokens(c)
	w := opts.Width
	if w == 0 {
		w = 96
	}
	n := len(lines)
	var res ChoiceResult
	move := func(top int) {
		top = max(0, min(top, n-view.Visible))
		if top != view.Top {
			view.Top = top
			res.changed = true
			slog.Debug("mujica: minimap move", "top", top)
		}
	}
	e := ui.Box(c).Width(w).Shrink(0).MinHeight(theme.MinHitSize*2).Focusable().Clip().Cursor(ui.CursorPointer).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Disabled(opts.Disabled).
		Label(msg(c, "minimap.label")).Description(fmt.Sprintf(msg(c, "minimap.range"), min(view.Top+1, n), min(view.Top+view.Visible, n), n))
	res.Element = e
	h := e.Bounds().H
	offset := minimapOffset(n, *view, h)
	if !e.IsDisabled() {
		if _, _, held := e.Dragged(); held {
			_, y, _ := e.PointerPosition()
			move(int((y+offset)/minimapLine) - view.Visible/2)
		}
		switch {
		case e.Shortcut(0, ui.KeyDown):
			move(view.Top + 1)
		case e.Shortcut(0, ui.KeyUp):
			move(view.Top - 1)
		case e.Shortcut(0, ui.KeyPageDown):
			move(view.Top + max(view.Visible, 1))
		case e.Shortcut(0, ui.KeyPageUp):
			move(view.Top - max(view.Visible, 1))
		case e.Shortcut(0, ui.KeyHome):
			move(0)
		case e.Shortcut(0, ui.KeyEnd):
			move(n)
		}
	}
	internal.Ring(e, k.Focus, theme.ControlRadius)
	v := *view
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.H != h {
			c.Invalidate()
		}
		p.Clip(r, theme.ControlRadius, func() {
			off := minimapOffset(n, v, r.H)
			first := int(off / minimapLine)
			for i := first; i < n && float32(i-first)*minimapLine < r.H+minimapLine; i++ {
				y := internal.Snap(p, r.Y+float32(i)*minimapLine-off)
				minimapStrokes(p, k, opts.Language, lines[i], r.X+6, y, r.X+r.W-6)
			}
			fy := internal.Snap(p, r.Y+float32(v.Top)*minimapLine-off)
			frame := ui.Rect{X: r.X + 1, Y: fy, W: r.W - 2, H: max(internal.Snap(p, float32(v.Visible)*minimapLine), 4)}
			p.Fill(frame, k.Selection.Alpha(0.55), 0)
			p.Stroke(frame, k.AccentText, 0, 1)
		})
	})
	return res
}

// minimapOffset scrolls a tall map with the view.
func minimapOffset(n int, v MinimapView, h float32) float32 {
	total := float32(n) * minimapLine
	if total <= h || n <= v.Visible {
		return 0
	}
	return float32(v.Top) / float32(n-v.Visible) * (total - h)
}

// minimapStrokes draws a line's tokens as short bars, skipping blanks.
func minimapStrokes(p *ui.Painter, k theme.Tokens, lang, line string, x0, y, x1 float32) {
	line = strings.ReplaceAll(line, "\t", "    ")
	toks, _ := internal.Highlight(lang, line, false)
	col := 0
	for _, t := range toks {
		runes := []rune(t.Text)
		start := -1
		for i, r := range append(runes, ' ') {
			if r != ' ' && start < 0 {
				start = i
			}
			if r == ' ' && start >= 0 {
				x := x0 + float32(col+start)*minimapChar
				if x < x1 {
					p.Fill(ui.Rect{X: x, Y: y, W: min(float32(i-start)*minimapChar, x1-x), H: minimapLine * 0.7}, codeTokenColor(k, t.Kind).Alpha(0.75), 0)
				}
				start = -1
			}
		}
		col += len(runes)
	}
}
