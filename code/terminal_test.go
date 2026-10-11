package code

import (
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestTerminalStateText(t *testing.T) {
	var s TerminalState
	s.Resize(10, 3)
	s.Write("hello\r\nwor")
	s.Write("ld\x1b[31m!\x1b[0m")
	if got := s.LineText(0); got != "hello" {
		t.Fatalf("line 0 %q", got)
	}
	if got := s.LineText(1); got != "world!" {
		t.Fatalf("line 1 %q", got)
	}
	if st := s.screen[1][5].st; st.FG.Index != 1 || st.FG.Kind == 0 {
		t.Fatalf("'!' style %+v, want red", st)
	}
	if st := s.screen[1][4].st; st.FG.Kind != 0 {
		t.Fatalf("'d' style %+v, want default", st)
	}
}

func TestTerminalStateCursorAndErase(t *testing.T) {
	var s TerminalState
	s.Resize(10, 3)
	s.Write("abcdef\x1b[3D")
	if x, _ := s.Cursor(); x != 3 {
		t.Fatalf("cursor after 3D at %d, want 3", x)
	}
	s.Write("\x1b[K")
	if got := s.LineText(0); got != "abc" {
		t.Fatalf("after erase-to-end %q", got)
	}
	s.Write("\x1b[2;4HX")
	if got := s.LineText(1); got != "   X" {
		t.Fatalf("after CUP %q", got)
	}
	s.Write("\x1b[2J")
	if s.LineText(0) != "" || s.LineText(1) != "" {
		t.Fatal("2J left text")
	}
	// A sequence split across writes completes on the next one.
	s.Write("\x1b[")
	s.Write("1;1Hz")
	if got := s.LineText(0); got != "z" {
		t.Fatalf("split CUP gave %q", got)
	}
}

func TestTerminalStateScrollback(t *testing.T) {
	var s TerminalState
	s.Resize(8, 2)
	for _, l := range []string{"one", "two", "three", "four"} {
		s.Write(l + "\r\n")
	}
	if s.Lines() != 5 || s.LineText(0) != "one" || s.LineText(3) != "four" {
		t.Fatalf("lines %d: %q %q", s.Lines(), s.LineText(0), s.LineText(3))
	}
	s.ScrollBack(10)
	if s.Offset() != 3 {
		t.Fatalf("offset %d, want the 3 scrollback lines", s.Offset())
	}
	s.ScrollBack(-1)
	if s.Offset() != 2 {
		t.Fatalf("offset %d after scrolling down one", s.Offset())
	}
}

func TestTerminalStateWrapsLongLines(t *testing.T) {
	var s TerminalState
	s.Resize(4, 3)
	s.Write("abcdefg")
	if s.LineText(0) != "abcd" || s.LineText(1) != "efg" {
		t.Fatalf("wrapped %q %q", s.LineText(0), s.LineText(1))
	}
}

func terminalTester(s *TerminalState, opts *TerminalOptions, typed *[]string) *ui.Tester {
	opts.Input = func(d string) { *typed = append(*typed, d) }
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Terminal(c, s, *opts).Size(400, 200)
	}, 420, 220)
}

func TestTerminalInputAndResize(t *testing.T) {
	var s TerminalState
	var typed []string
	var size [2]int
	opts := TerminalOptions{Resized: func(cols, rows int) { size = [2]int{cols, rows} }}
	tt := terminalTester(&s, &opts, &typed)
	if size[0] < 30 || size[1] < 5 {
		t.Fatalf("resize reported %v for a 400×200 view", size)
	}
	if cols, rows := s.Size(); cols != size[0] || rows != size[1] {
		t.Fatalf("state %d×%d, reported %v", cols, rows, size)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Terminal") {
		t.Fatal("Tab did not focus the terminal")
	}
	tt.Type("ls")
	tt.Key(0, ui.KeyEnter)
	tt.Key(ui.Ctrl, ui.KeyC)
	tt.Key(0, ui.KeyUp)
	want := []string{"ls", "\r", "\x03", "\x1b[A"}
	if strings.Join(typed, "|") != strings.Join(want, "|") {
		t.Fatalf("typed %q, want %q", typed, want)
	}
}

func TestTerminalSelectCopyPaste(t *testing.T) {
	var s TerminalState
	var typed []string
	tt := terminalTester(&s, &TerminalOptions{}, &typed)
	s.Write("alpha beta")
	tt.Frame()
	cw, _, _ := terminalMetrics(ui.Font{Family: "monospace", Size: 13})
	y := float32(terminalPad + 5)
	tt.Press(terminalPad+1+cw*0.5, y)
	tt.Move(terminalPad+1+cw*4.5, y)
	tt.Release(terminalPad+1+cw*4.5, y)
	if got := s.Selection(); got != "alpha" {
		t.Fatalf("selection %q, want alpha", got)
	}
	mod := ui.Super
	if runtime.GOOS != "darwin" {
		mod = ui.Ctrl | ui.Shift
	}
	tt.Key(mod, ui.KeyC)
	if tt.Clipboard() != "alpha" {
		t.Fatalf("clipboard %q", tt.Clipboard())
	}
	tt.SetClipboard("echo hi")
	tt.Key(mod, ui.KeyV)
	if len(typed) != 1 || typed[0] != "echo hi" {
		t.Fatalf("paste sent %q", typed)
	}
}

func TestTerminalDisabledSendsNothing(t *testing.T) {
	var s TerminalState
	var typed []string
	tt := terminalTester(&s, &TerminalOptions{Disabled: true}, &typed)
	tt.Key(0, ui.KeyTab)
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if len(typed) != 0 {
		t.Fatalf("disabled terminal sent %q", typed)
	}
}

func TestTerminalScrollbackKeys(t *testing.T) {
	var s TerminalState
	var typed []string
	tt := terminalTester(&s, &TerminalOptions{}, &typed)
	for range 40 {
		s.Write("line\r\n")
	}
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyPageUp)
	if s.Offset() == 0 {
		t.Fatal("Shift+PageUp did not scroll back")
	}
	tt.Type("q")
	if s.Offset() != 0 {
		t.Fatal("typing did not return to the bottom")
	}
}

func TestTerminalSnapshot(t *testing.T) {
	snapshot(t, "terminal", 420, 160, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *TerminalState {
			s := &TerminalState{}
			s.Resize(50, 7)
			s.Write("\x1b[1m~/court\x1b[0m $ ls\r\n\x1b[34mchronicles\x1b[0m  decrees.txt  \x1b[32mrun.sh\x1b[0m\r\n")
			s.Write("\x1b[33mwarning:\x1b[0m ledger is \x1b[7mlocked\x1b[0m\r\n\x1b[31merror\x1b[0m: missing seal\r\n\x1b[1m~/court\x1b[0m $ ")
			return s
		})
		ui.Column(c).Padding(12).Children(func() {
			Terminal(c, *s, TerminalOptions{}).Height(136)
		})
	})
}
