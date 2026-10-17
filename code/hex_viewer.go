package code

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// hexPerRow is how many bytes a row shows.
const hexPerRow = 16

// HexViewerState is the cursor and selection anchor.
type HexViewerState struct {
	// Anchor is where a Shift selection started.
	Cursor, Anchor int
	rows           ui.ListState
}

// Selection returns the selected byte range, start inclusive, end exclusive.
func (s *HexViewerState) Selection() (start, end int) {
	return min(s.Cursor, s.Anchor), max(s.Cursor, s.Anchor) + 1
}

// HexViewerOptions configure a HexViewer.
type HexViewerOptions struct {
	Disabled bool
}

// HexViewer shows bytes as synced hex and ASCII columns.
func HexViewer(c *ui.Context, data []byte, s *HexViewerState, opts HexViewerOptions) core.ChoiceResult {
	if s == nil {
		panic("mujica: HexViewer needs a state")
	}
	n := len(data)
	if n > 0 && (s.Cursor < 0 || s.Cursor >= n || s.Anchor < 0 || s.Anchor >= n) {
		panic("mujica: HexViewer cursor outside the data")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	cw := codeTextWidth(c, strings.Repeat("0", hexPerRow*3)) / (hexPerRow * 3)
	var res core.ChoiceResult
	frame := CodePanel(c).Focusable().Disabled(opts.Disabled).Label(core.Msg(c, "hex_viewer.label"))
	res.Element = frame
	internal.Ring(frame, k.Focus, theme.ControlRadius)
	move := func(to int, extend bool) {
		to = max(0, min(to, n-1))
		if to == s.Cursor && (extend || s.Anchor == to) {
			return
		}
		s.Cursor = to
		if !extend {
			s.Anchor = to
		}
		s.rows.ScrollIntoView(to / hexPerRow)
		res = core.NewChoiceResult(res.Element, true)
		slog.Debug("mujica: hex_viewer select", "cursor", s.Cursor, "anchor", s.Anchor)
	}
	if n > 0 && !frame.IsDisabled() {
		hexKeys(c, frame, data, s, move)
	}
	a, b := s.Selection()
	frame.Description(fmt.Sprintf(core.Msg(c, "hex_viewer.at"), s.Cursor, b-a))
	frame.Children(func() {
		if n == 0 {
			CodeCaption(c, core.Msg(c, "hex_viewer.empty")).Padding(d.Unit() * 4)
			return
		}
		rows := (n + hexPerRow - 1) / hexPerRow
		ui.List(c, &s.rows, rows, func(r int) {
			start := r * hexPerRow
			end := min(start+hexPerRow, n)
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(cw * 2).AlignItems(ui.Center).Children(func() {
				CodeMono(c, fmt.Sprintf("%08x", start)).TextColor(k.TextMuted).Shrink(0)
				hexColumn(c, s, data, start, end, cw, true, move)
				hexColumn(c, s, data, start, end, cw, false, move)
			})
		}).Grow(1).MinHeight(d.ControlHeight() * 3)
	})
	return res
}

// hexColumn draws a row's hex or ASCII with marks.
func hexColumn(c *ui.Context, s *HexViewerState, data []byte, start, end int, cw float32, hex bool, move func(int, bool)) {
	k := core.Tokens(c)
	var b strings.Builder
	step := cw
	if hex {
		step = 3 * cw
	}
	for i := start; i < end; i++ {
		if hex {
			fmt.Fprintf(&b, "%02x ", data[i])
		} else if data[i] >= 0x20 && data[i] < 0x7f {
			b.WriteByte(data[i])
		} else {
			b.WriteByte('.')
		}
	}
	box := ui.Box(c).Shrink(0).Width(step * hexPerRow).Cursor(ui.CursorText)
	lo, hi := s.Selection()
	// A hex mark covers the two digits.
	mark := step
	if hex {
		mark = 2 * cw
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		for i := max(lo, start); i < min(hi, end); i++ {
			x := r.X + float32(i-start)*step
			p.Fill(ui.Rect{X: x - 1, Y: r.Y, W: mark + 2, H: r.H}, k.Selection, 2)
		}
		if s.Cursor >= start && s.Cursor < end {
			x := r.X + float32(s.Cursor-start)*step
			p.Stroke(ui.Rect{X: x - 1, Y: r.Y, W: mark + 2, H: r.H}, k.AccentText, 2, 1.5)
		}
	})
	if box.Clicked() && !box.IsDisabled() {
		x, _, _ := box.PointerPosition()
		move(start+min(int(x/step), end-start-1), box.ClickModifiers()&ui.Shift != 0)
	}
	box.Children(func() { CodeMono(c, strings.TrimRight(b.String(), " ")) })
}

// hexKeys moves, extends and copies the selection.
func hexKeys(c *ui.Context, frame ui.Element, data []byte, s *HexViewerState, move func(int, bool)) {
	copyMods := ui.Cmd
	if runtime.GOOS != "darwin" {
		copyMods = ui.Ctrl
	}
	for _, mods := range []ui.Modifiers{0, ui.Shift} {
		ext := mods == ui.Shift
		switch {
		case frame.Shortcut(mods, ui.KeyRight):
			move(s.Cursor+1, ext)
		case frame.Shortcut(mods, ui.KeyLeft):
			move(s.Cursor-1, ext)
		case frame.Shortcut(mods, ui.KeyDown):
			move(s.Cursor+hexPerRow, ext)
		case frame.Shortcut(mods, ui.KeyUp):
			move(s.Cursor-hexPerRow, ext)
		case frame.Shortcut(mods, ui.KeyHome):
			move(s.Cursor/hexPerRow*hexPerRow, ext)
		case frame.Shortcut(mods, ui.KeyEnd):
			move(s.Cursor/hexPerRow*hexPerRow+hexPerRow-1, ext)
		}
	}
	if frame.Shortcut(copyMods, ui.KeyC) {
		a, b := s.Selection()
		c.WriteClipboard(fmt.Sprintf("% x", data[a:b]))
		slog.Debug("mujica: hex_viewer copy", "bytes", b-a)
	}
}
