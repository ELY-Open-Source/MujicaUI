package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func cmdTester(run *CommandRun, opts CommandExecutionCardOptions, aborts *int) *ui.Tester {
	opts.Still = true
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if CommandExecutionCard(c, *run, opts).Changed() {
			*aborts++
		}
	}, 600, 420)
}

func TestCommandExecutionCardRunningAbortAndFinish(t *testing.T) {
	run := CommandRun{Command: "go test ./...", Dir: "~/court", Running: true, Output: "\x1b[32mok\x1b[0m  court\n"}
	aborts := 0
	tt := cmdTester(&run, CommandExecutionCardOptions{}, &aborts)
	for _, s := range []string{"go test ./...", "~/court", "ok", "  court", "Running"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	if tt.HasText("\x1b[32mok") {
		t.Fatal("escape codes shown as text")
	}
	_ = tt.Click("Abort")
	if aborts != 1 {
		t.Fatalf("aborts %d", aborts)
	}
	run.Output += "FAIL  ledger\n"
	run.Running, run.ExitCode, run.Duration = false, 1, 3200*time.Millisecond
	tt.Frame()
	if !tt.HasText("FAIL  ledger") || !tt.HasText("Exit 1") || !tt.HasText("3.2 s") || tt.HasText("Abort") {
		t.Fatalf("finished: %q", tt.Texts())
	}
}

func TestCommandExecutionCardAnsiAcrossLinesAndReset(t *testing.T) {
	var o commandOutput
	o.update("\x1b[31mred\nstill red\x1b[0m plain\n")
	if len(o.lines) != 2 || o.lines[1][0].Style.FG.Index != 1 || o.lines[1][1].Style.FG.Index != 0 {
		t.Fatalf("lines %+v", o.lines)
	}
	o.update("\x1b[31mred\nstill red\x1b[0m plain\nnext\n")
	if len(o.lines) != 3 || o.lines[2][0].Text != "next" {
		t.Fatalf("append: %+v", o.lines)
	}
	o.update("fresh\n")
	if len(o.lines) != 1 || o.lines[0][0].Text != "fresh" {
		t.Fatalf("replace: %+v", o.lines)
	}
}

func TestCommandExecutionCardDisabledAndEmpty(t *testing.T) {
	run := CommandRun{Command: "sleep 10", Running: true}
	aborts := 0
	tt := cmdTester(&run, CommandExecutionCardOptions{Disabled: true}, &aborts)
	_ = tt.Click("Abort")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if aborts != 0 || !tt.HasText("No output") {
		t.Fatalf("aborts %d texts %q", aborts, tt.Texts())
	}
}

func TestCommandExecutionCardLargeOutputVirtualized(t *testing.T) {
	var b strings.Builder
	for i := range 100000 {
		fmt.Fprintf(&b, "\x1b[2mline %d\x1b[0m\n", i)
	}
	run := CommandRun{Command: "yes", Output: b.String()}
	n := 0
	tt := cmdTester(&run, CommandExecutionCardOptions{}, &n)
	if got := len(tt.Texts()); got > 200 {
		t.Fatalf("%d texts built for 100000 lines", got)
	}
	if !tt.HasText("line 99999") {
		t.Fatal("output does not follow its end")
	}
}

func TestCommandExecutionCardSnapshot(t *testing.T) {
	snapshot(t, "command_execution_card", 560, 300, func(c *ui.Context) {
		run := CommandRun{Command: "go test ./court/...", Dir: "~/projects/court", ExitCode: 1, Duration: 4200 * time.Millisecond,
			Output: "\x1b[32mok\x1b[0m   court/guests  0.21s\n\x1b[1;31mFAIL\x1b[0m court/ledger  0.08s\n    ledger_test.go:14: want 42 guests, got 40\n\x1b[33mwarning:\x1b[0m 1 test skipped\n"}
		ui.Column(c).Padding(16).Children(func() {
			CommandExecutionCard(c, run, CommandExecutionCardOptions{Height: 140})
		})
	})
}
