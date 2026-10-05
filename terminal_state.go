package mujica

import (
	"log/slog"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
)

// terminalCell is one grid character; zero is blank.
type terminalCell struct {
	r  rune
	st internal.AnsiStyle
}

// TerminalPos is a cell; Line counts from the oldest line.
type TerminalPos struct{ Line, Col int }

func (a TerminalPos) before(b TerminalPos) bool {
	return a.Line < b.Line || a.Line == b.Line && a.Col < b.Col
}

// TerminalState is the screen; zero is an empty 80×24 grid.
type TerminalState struct {
	cols, rows int
	screen     [][]terminalCell
	back       [][]terminalCell
	cx, cy     int
	style      internal.AnsiStyle
	pending    string
	limit      int
	// offset is how far the view scrolled back.
	offset   int
	scrolled float32
	sel      [2]TerminalPos
	selOn    bool
	matches  []terminalMatch
	current  int
	version  int
	shaped   map[string][]ui.Glyph
	preedit  string
}

type terminalMatch struct {
	pos TerminalPos
	n   int
}

func (s *TerminalState) init() {
	if s.screen == nil {
		s.Resize(80, 24)
	}
	if s.limit == 0 {
		s.limit = 1000
	}
}

// Size returns the grid's columns and rows.
func (s *TerminalState) Size() (cols, rows int) {
	s.init()
	return s.cols, s.rows
}

// Cursor returns the cursor's column and row on the screen.
func (s *TerminalState) Cursor() (col, row int) { return s.cx, s.cy }

// Resize sets the grid size, keeping the cursor in view.
func (s *TerminalState) Resize(cols, rows int) {
	if cols < 1 || rows < 1 {
		panic("mujica: TerminalState size must be at least 1×1")
	}
	if s.limit == 0 {
		s.limit = 1000
	}
	old := s.screen
	// Lines above a falling-off cursor go to the scrollback.
	for drop := s.cy - (rows - 1); drop > 0 && len(old) > 0; drop-- {
		s.pushBack(old[0])
		old = old[1:]
		s.cy--
		if s.offset > 0 {
			s.offset++
		}
	}
	s.screen = make([][]terminalCell, rows)
	for y := range rows {
		s.screen[y] = make([]terminalCell, cols)
		if y < len(old) {
			copy(s.screen[y], old[y])
		}
	}
	s.cols, s.rows = cols, rows
	s.cx, s.cy = min(s.cx, cols-1), min(max(s.cy, 0), rows-1)
	s.version++
}

// Write feeds PTY output; call it on the UI thread.
func (s *TerminalState) Write(data string) {
	s.init()
	data = s.pending + data
	s.pending = ""
	for i := 0; i < len(data); {
		b := data[i]
		switch {
		case b == 0x1b:
			n, ok := s.escape(data[i:])
			if !ok {
				s.pending = data[i:]
				s.version++
				return
			}
			i += n
		case b == '\r':
			s.cx = 0
			i++
		case b == '\n':
			s.lineFeed()
			i++
		case b == '\b':
			s.cx = max(s.cx-1, 0)
			i++
		case b == '\t':
			s.cx = min((s.cx/8+1)*8, s.cols-1)
			i++
		case b < 0x20 || b == 0x7f:
			i++
		default:
			r, n := utf8.DecodeRuneInString(data[i:])
			s.put(r)
			i += n
		}
	}
	s.version++
}

func (s *TerminalState) put(r rune) {
	if s.cx >= s.cols {
		s.cx = 0
		s.lineFeed()
	}
	s.screen[s.cy][s.cx] = terminalCell{r, s.style}
	s.cx++
}

func (s *TerminalState) lineFeed() {
	if s.cy < s.rows-1 {
		s.cy++
		return
	}
	s.pushBack(s.screen[0])
	copy(s.screen, s.screen[1:])
	s.screen[s.rows-1] = make([]terminalCell, s.cols)
	if s.offset > 0 {
		s.offset = min(s.offset+1, len(s.back))
	}
}

func (s *TerminalState) pushBack(line []terminalCell) {
	s.back = append(s.back, line)
	if over := len(s.back) - s.limit; over > 0 {
		s.back = s.back[over:]
		s.offset = min(s.offset, len(s.back))
		s.selOn = false
	}
}

// escape runs the sequence at d; false means incomplete.
func (s *TerminalState) escape(d string) (int, bool) {
	if len(d) < 2 {
		return 0, false
	}
	switch d[1] {
	case '[':
	case ']':
		// OSC ends with BEL or ST; nothing to draw.
		if end := strings.IndexAny(d, "\a"); end >= 0 {
			return end + 1, true
		}
		if end := strings.Index(d, "\x1b\\"); end >= 0 {
			return end + 2, true
		}
		return 0, false
	default:
		return 2, true
	}
	j := 2
	for j < len(d) && (d[j] < 0x40 || d[j] > 0x7e) {
		j++
	}
	if j >= len(d) {
		return 0, false
	}
	params, final := d[2:j], d[j]
	if final == 'm' {
		_, s.style = internal.ParseAnsi(d[:j+1], s.style)
		return j + 1, true
	}
	s.control(params, final)
	return j + 1, true
}

// control applies a CSI cursor or erase command.
func (s *TerminalState) control(params string, final byte) {
	args := strings.Split(params, ";")
	arg := func(i, def int) int {
		if i >= len(args) || args[i] == "" {
			return def
		}
		v, err := strconv.Atoi(args[i])
		if err != nil {
			return def
		}
		return v
	}
	switch final {
	case 'A':
		s.cy = max(s.cy-max(arg(0, 1), 1), 0)
	case 'B':
		s.cy = min(s.cy+max(arg(0, 1), 1), s.rows-1)
	case 'C':
		s.cx = min(s.cx+max(arg(0, 1), 1), s.cols-1)
	case 'D':
		s.cx = max(min(s.cx, s.cols-1)-max(arg(0, 1), 1), 0)
	case 'G':
		s.cx = min(max(arg(0, 1), 1), s.cols) - 1
	case 'H', 'f':
		s.cy = min(max(arg(0, 1), 1), s.rows) - 1
		s.cx = min(max(arg(1, 1), 1), s.cols) - 1
	case 'J':
		s.eraseDisplay(arg(0, 0))
	case 'K':
		s.eraseLine(s.cy, arg(0, 0))
	default:
		slog.Debug("mujica: terminal unsupported sequence", "params", params, "final", string(final))
	}
}

func (s *TerminalState) eraseLine(y, mode int) {
	line := s.screen[y]
	from, to := 0, s.cols
	switch mode {
	case 0:
		from = min(s.cx, s.cols)
	case 1:
		to = min(s.cx+1, s.cols)
	}
	clear(line[from:to])
}

func (s *TerminalState) eraseDisplay(mode int) {
	switch mode {
	case 0:
		s.eraseLine(s.cy, 0)
		for y := s.cy + 1; y < s.rows; y++ {
			clear(s.screen[y])
		}
	case 1:
		s.eraseLine(s.cy, 1)
		for y := range s.cy {
			clear(s.screen[y])
		}
	default:
		for y := range s.rows {
			clear(s.screen[y])
		}
		if mode == 3 {
			s.back, s.offset = nil, 0
		}
	}
}

// Lines counts scrollback and screen lines.
func (s *TerminalState) Lines() int {
	s.init()
	return len(s.back) + s.rows
}

func (s *TerminalState) line(i int) []terminalCell {
	if i < len(s.back) {
		return s.back[i]
	}
	return s.screen[i-len(s.back)]
}

// LineText returns line i without trailing blanks.
func (s *TerminalState) LineText(i int) string {
	var b strings.Builder
	for _, cell := range s.line(i) {
		r := cell.r
		if r == 0 {
			r = ' '
		}
		b.WriteRune(r)
	}
	return strings.TrimRight(b.String(), " ")
}

// Select marks the cells from a to b.
func (s *TerminalState) Select(a, b TerminalPos) {
	if b.before(a) {
		a, b = b, a
	}
	s.sel, s.selOn = [2]TerminalPos{a, b}, true
}

// Selection returns the selected text, lines joined by newlines.
func (s *TerminalState) Selection() string {
	if !s.selOn {
		return ""
	}
	a, b := s.sel[0], s.sel[1]
	var out []string
	for y := a.Line; y <= b.Line && y < s.Lines(); y++ {
		cells := s.line(y)
		from, to := 0, len(cells)
		if y == a.Line {
			from = min(a.Col, len(cells))
		}
		if y == b.Line {
			to = min(b.Col+1, len(cells))
		}
		var sb strings.Builder
		for _, cell := range cells[from:max(to, from)] {
			r := cell.r
			if r == 0 {
				r = ' '
			}
			sb.WriteRune(r)
		}
		out = append(out, strings.TrimRight(sb.String(), " "))
	}
	return strings.Join(out, "\n")
}

func (s *TerminalState) selected(p TerminalPos) bool {
	return s.selOn && !p.before(s.sel[0]) && !s.sel[1].before(p)
}

// ScrollBack moves the view n lines up; negative goes down.
func (s *TerminalState) ScrollBack(n int) {
	s.init()
	s.offset = max(0, min(s.offset+n, len(s.back)))
}

// Offset returns how many lines the view is scrolled back.
func (s *TerminalState) Offset() int { return s.offset }

// reveal scrolls so that line shows on screen.
func (s *TerminalState) reveal(line int) {
	top := len(s.back) - s.offset
	switch {
	case line < top:
		s.offset = len(s.back) - line
	case line >= top+s.rows:
		s.offset = max(0, len(s.back)-(line-s.rows+1))
	}
}
