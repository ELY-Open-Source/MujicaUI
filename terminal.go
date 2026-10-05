package mujica

import (
	"log/slog"
	"math"
	"runtime"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TerminalOptions configure a Terminal.
type TerminalOptions struct {
	// Input receives what the user types, encoded for the PTY.
	Input func(data string)
	// Resized reports a new grid size for the PTY.
	Resized func(cols, rows int)
	// Scrollback caps lines kept above; zero keeps 1000.
	Scrollback int
	Label      string
	Disabled   bool
}

// terminalPad is the space around the grid in DIPs.
const terminalPad float32 = 6

// terminalCtrl maps letter keys to their Ctrl codes.
var terminalCtrl = map[ui.Key]byte{ui.KeyA: 1, ui.KeyB: 2, ui.KeyC: 3, ui.KeyD: 4, ui.KeyE: 5, ui.KeyF: 6, ui.KeyG: 7,
	ui.KeyH: 8, ui.KeyI: 9, ui.KeyJ: 10, ui.KeyK: 11, ui.KeyL: 12, ui.KeyM: 13, ui.KeyN: 14, ui.KeyO: 15, ui.KeyP: 16,
	ui.KeyQ: 17, ui.KeyR: 18, ui.KeyS: 19, ui.KeyT: 20, ui.KeyU: 21, ui.KeyV: 22, ui.KeyW: 23, ui.KeyX: 24, ui.KeyY: 25, ui.KeyZ: 26}

// terminalKeys are the sequences of keys that type no text.
var terminalKeys = map[ui.Key]string{ui.KeyEnter: "\r", ui.KeyBackspace: "\x7f", ui.KeyTab: "\t", ui.KeyEscape: "\x1b",
	ui.KeyUp: "\x1b[A", ui.KeyDown: "\x1b[B", ui.KeyRight: "\x1b[C", ui.KeyLeft: "\x1b[D", ui.KeyHome: "\x1b[H", ui.KeyEnd: "\x1b[F",
	ui.KeyDelete: "\x1b[3~", ui.KeyPageUp: "\x1b[5~", ui.KeyPageDown: "\x1b[6~"}

// terminalCell size and baseline for the code font.
func terminalMetrics(f ui.Font) (cw, ch, base float32) {
	m := f.Metrics()
	g := ui.Shape("M", f)
	cw = f.Size * 0.6
	if len(g) > 0 {
		cw = g[0].Advance
	}
	ch = float32(math.Ceil(float64((m.Ascent + m.Descent + m.LineGap) * 1.15)))
	return cw, ch, (ch-m.Ascent-m.Descent)/2 + m.Ascent
}

// Terminal draws a TerminalState and sends input to the PTY.
func Terminal(c *ui.Context, s *TerminalState, opts TerminalOptions) *ui.Element {
	if s == nil {
		panic("mujica: Terminal needs a TerminalState")
	}
	if opts.Scrollback < 0 {
		panic("mujica: Terminal Scrollback is negative")
	}
	s.init()
	if opts.Scrollback > 0 {
		s.limit = opts.Scrollback
	}
	k := Tokens(c)
	font := codeFont(c)
	cw, ch, base := terminalMetrics(font)
	label := opts.Label
	if label == "" {
		label = msg(c, "terminal.label")
	}
	e := ui.Box(c).Focusable().Cursor(ui.CursorText).Clip().Label(label).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Disabled(opts.Disabled).MinHeight(ch * 4)
	b := e.Bounds()
	if b.W > 0 {
		cols := max(int((b.W-2*terminalPad)/cw), 1)
		rows := max(int((b.H-2*terminalPad)/ch), 1)
		if cols != s.cols || rows != s.rows {
			s.Resize(cols, rows)
			slog.Debug("mujica: terminal resize", "cols", cols, "rows", rows)
			if opts.Resized != nil {
				opts.Resized(cols, rows)
			}
		}
	}
	send := func(data string) {
		s.offset, s.selOn, s.preedit = 0, false, ""
		if opts.Input != nil {
			opts.Input(data)
		}
	}
	cell := func(x, y float32) TerminalPos {
		col := max(0, min(int((x-terminalPad)/cw), s.cols-1))
		row := max(0, min(int((y-terminalPad)/ch), s.rows-1))
		return TerminalPos{Line: len(s.back) - s.offset + row, Col: col}
	}
	anchor := ui.Local(e, "terminal.anchor", func() TerminalPos { return TerminalPos{} })
	e.HandleInput(func(ev ui.InputEvent) bool {
		if e.IsDisabled() {
			return false
		}
		switch ev.Kind {
		case ui.InputKeyDown:
			return terminalKey(c, s, ev, send)
		case ui.InputText:
			send(ev.Text)
			return true
		case ui.InputCompose:
			s.preedit = ev.Text
			return true
		case ui.InputScroll:
			s.scrolled -= ev.DY / ch
			n := int(s.scrolled)
			s.scrolled -= float32(n)
			s.ScrollBack(n)
			return true
		case ui.InputPointerDown:
			if ev.Button != 0 {
				return false
			}
			e.Focus()
			*anchor = cell(ev.X, ev.Y)
			s.selOn = false
			return true
		case ui.InputPointerMove:
			if e.Pressed() {
				s.Select(*anchor, cell(ev.X, ev.Y))
			}
			return true
		}
		return false
	})
	cy := float32(s.cy+s.offset) * ch
	e.TextCaret(ui.Rect{X: terminalPad + float32(s.cx)*cw, Y: terminalPad + cy, W: cw, H: ch})
	focused := e.Focused()
	internal.Ring(e, k.Focus, theme.ControlRadius)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W != b.W || r.H != b.H {
			// Laid out at another size: rebuild to fit the grid.
			c.Invalidate()
		}
		terminalPaint(p, r, s, k, font, cw, ch, base, focused)
	})
	return e
}

// terminalKey encodes a key, handling copy, paste and scrollback.
func terminalKey(c *ui.Context, s *TerminalState, ev ui.InputEvent, send func(string)) bool {
	clip := ui.Super
	if runtime.GOOS != "darwin" {
		clip = ui.Ctrl | ui.Shift
	}
	switch {
	case ev.Mods == clip && ev.Key == ui.KeyC:
		if t := s.Selection(); t != "" {
			c.WriteClipboard(t)
			slog.Debug("mujica: terminal copy", "runes", len([]rune(t)))
		}
		return true
	case ev.Mods == clip && ev.Key == ui.KeyV:
		if t := c.ReadClipboard(); t != "" {
			send(t)
		}
		return true
	case ev.Mods == ui.Shift && (ev.Key == ui.KeyPageUp || ev.Key == ui.KeyPageDown):
		n := s.rows - 1
		if ev.Key == ui.KeyPageDown {
			n = -n
		}
		s.ScrollBack(n)
		return true
	case ev.Mods == ui.Ctrl && terminalCtrl[ev.Key] != 0:
		send(string([]byte{terminalCtrl[ev.Key]}))
		return true
	}
	if seq, ok := terminalKeys[ev.Key]; ok && ev.Mods&^ui.Shift == 0 {
		if ev.Key == ui.KeyTab && ev.Mods == ui.Shift {
			seq = "\x1b[Z"
		}
		send(seq)
		return true
	}
	return false
}

// terminalPaint draws the visible rows, marks and cursor.
func terminalPaint(p *ui.Painter, r ui.Rect, s *TerminalState, k theme.Tokens, font ui.Font, cw, ch, base float32, focused bool) {
	ox, oy := r.X+terminalPad, r.Y+terminalPad
	top := len(s.back) - s.offset
	for row := range s.rows {
		li := top + row
		cells := s.line(li)
		y := oy + float32(row)*ch
		for col := 0; col < len(cells); {
			st := cells[col].st
			end := col + 1
			for end < len(cells) && cells[end].st == st {
				end++
			}
			fg, bg := codeAnsiColor(k, st.FG), codeAnsiColor(k, st.BG)
			if fg.A == 0 {
				fg = k.Text
			}
			if st.Inverse {
				if bg.A == 0 {
					bg = k.Surface
				}
				fg, bg = bg, fg
			}
			if st.Dim {
				fg = fg.Alpha(0.6)
			}
			x0 := ox + float32(col)*cw
			if bg.A > 0 {
				p.Fill(ui.Rect{X: x0, Y: y, W: float32(end-col) * cw, H: ch}, bg, 0)
			}
			for i := col; i < end; i++ {
				if s.selected(TerminalPos{li, i}) {
					p.Fill(ui.Rect{X: ox + float32(i)*cw, Y: y, W: cw, H: ch}, k.Selection, 0)
				}
			}
			terminalRun(p, s, cells[col:end], st, font, x0, y+base, cw, fg)
			if st.Underline || st.Strike {
				ly := y + base + 2
				if st.Strike {
					ly = y + ch/2
				}
				internal.Rule(p, x0, x0+float32(end-col)*cw, ly, fg)
			}
			col = end
		}
		for mi, m := range s.matches {
			if m.pos.Line != li {
				continue
			}
			mr := ui.Rect{X: ox + float32(m.pos.Col)*cw, Y: y, W: float32(m.n) * cw, H: ch}
			internal.Rule(p, mr.X, mr.X+mr.W, y+ch-1, k.AccentText)
			if mi == s.current {
				p.Stroke(mr, k.AccentText, 1, 1.5)
			}
		}
	}
	if s.offset > 0 || s.cy >= s.rows {
		return
	}
	cr := ui.Rect{X: ox + float32(s.cx)*cw, Y: oy + float32(s.cy)*ch, W: cw, H: ch}
	if s.cx >= s.cols {
		cr.X = ox + float32(s.cols-1)*cw
	}
	if s.preedit != "" {
		g := terminalShape(s, s.preedit, font)
		p.Glyphs(g, cr.X, cr.Y+base, k.AccentText)
		internal.Rule(p, cr.X, cr.X+float32(len([]rune(s.preedit)))*cw, cr.Y+ch-1, k.AccentText)
		return
	}
	if !focused {
		p.Stroke(cr, k.AccentText, 0, 1)
		return
	}
	p.Fill(cr, k.AccentText, 0)
	if c := min(s.cx, s.cols-1); s.screen[s.cy][c].r > ' ' {
		p.Glyphs(terminalShape(s, string(s.screen[s.cy][c].r), font), cr.X, cr.Y+base, k.Surface)
	}
}

// terminalRun draws same-style cells, one glyph per cell.
func terminalRun(p *ui.Painter, s *TerminalState, cells []terminalCell, st internal.AnsiStyle, font ui.Font, x, y, cw float32, fg ui.Color) {
	runes := make([]rune, len(cells))
	blank := true
	for i, cell := range cells {
		runes[i] = cell.r
		if cell.r == 0 {
			runes[i] = ' '
		} else if cell.r != ' ' {
			blank = false
		}
	}
	if blank {
		return
	}
	if st.Bold {
		font.Weight = 700
	}
	font.Italic = st.Italic
	p.Glyphs(terminalShape(s, string(runes), font), x, y, fg)
}

// terminalShape shapes text into cells, cached across frames.
func terminalShape(s *TerminalState, text string, f ui.Font) []ui.Glyph {
	key := string(rune(f.Weight)) + string(rune(int(f.Size*4))) + text
	if f.Italic {
		key = "i" + key
	}
	if g, ok := s.shaped[key]; ok {
		return g
	}
	if s.shaped == nil || len(s.shaped) > 4096 {
		s.shaped = map[string][]ui.Glyph{}
	}
	g := ui.Shape(text, f)
	cw := f.Size * 0.6
	if m := ui.Shape("M", ui.Font{Family: f.Family, Size: f.Size}); len(m) > 0 {
		cw = m[0].Advance
	}
	for i := range g {
		g[i].X = float32(g[i].Cluster) * cw
	}
	s.shaped[key] = g
	return g
}
