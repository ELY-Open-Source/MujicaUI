package code

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func terminalSearchTester(term *TerminalState, s *TerminalSearchState, opts TerminalSearchOptions, moves *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			if TerminalSearch(c, term, s, opts).Changed() {
				*moves++
			}
			Terminal(c, term, TerminalOptions{}).Size(400, 120)
		})
	}, 420, 220)
}

func terminalSearchFill() *TerminalState {
	term := &TerminalState{}
	term.Resize(40, 5)
	for i := range 30 {
		if i%10 == 3 {
			term.Write(fmt.Sprintf("step %d: Error found\r\n", i))
		} else {
			term.Write(fmt.Sprintf("step %d: ok\r\n", i))
		}
	}
	return term
}

func TestTerminalSearchCountsAndJumps(t *testing.T) {
	term := terminalSearchFill()
	s := &TerminalSearchState{}
	moves := 0
	tt := terminalSearchTester(term, s, TerminalSearchOptions{}, &moves)
	tt.Key(0, ui.KeyTab)
	tt.Type("error")
	if len(term.matches) != 3 || !tt.HasText("1 of 3") {
		t.Fatalf("matches %d, texts %q", len(term.matches), tt.Texts())
	}
	first := term.matches[0].pos.Line
	top := len(term.back) - term.Offset()
	if first < top || first >= top+5 {
		t.Fatalf("match line %d not on screen from %d", first, top)
	}
	tt.Key(0, ui.KeyEnter)
	if s.Current != 1 || !tt.HasText("2 of 3") {
		t.Fatalf("Enter moved to %d", s.Current)
	}
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Key(ui.Shift, ui.KeyEnter)
	if s.Current != 2 {
		t.Fatalf("Shift+Enter twice from 1 gave %d, want 2 (wraps)", s.Current)
	}
	if line := term.matches[2].pos.Line; line < len(term.back)-term.Offset() {
		t.Fatalf("current match line %d above the view", line)
	}
	if moves < 4 {
		t.Fatalf("Changed reported %d times", moves)
	}
}

func TestTerminalSearchFollowsNewOutput(t *testing.T) {
	term := terminalSearchFill()
	s := &TerminalSearchState{Query: "error"}
	moves := 0
	tt := terminalSearchTester(term, s, TerminalSearchOptions{}, &moves)
	term.Write("late error\r\n")
	tt.Frame()
	if len(term.matches) != 4 {
		t.Fatalf("%d matches after new output, want 4", len(term.matches))
	}
}

func TestTerminalSearchNoMatchAndDisabled(t *testing.T) {
	term := terminalSearchFill()
	s := &TerminalSearchState{Query: "zebra"}
	moves := 0
	tt := terminalSearchTester(term, s, TerminalSearchOptions{Disabled: true}, &moves)
	if !tt.HasText("No matches") || len(term.matches) != 0 {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("ok")
	if s.Query != "zebra" || moves != 0 {
		t.Fatalf("disabled search changed to %q", s.Query)
	}
}

func TestTerminalSearchSnapshot(t *testing.T) {
	snapshot(t, "terminal_search", 420, 170, func(c *ui.Context) {
		st := ui.Local(c.Root(), "snap", func() [2]any {
			return [2]any{terminalSearchFill(), &TerminalSearchState{Query: "error", Current: 1}}
		})
		term, s := (*st)[0].(*TerminalState), (*st)[1].(*TerminalSearchState)
		ui.Column(c).Padding(12).Gap(8).Children(func() {
			TerminalSearch(c, term, s, TerminalSearchOptions{})
			Terminal(c, term, TerminalOptions{}).Height(96)
		})
	})
}
