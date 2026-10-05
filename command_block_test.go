package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func commandBlockSample() CommandRun {
	return CommandRun{Command: "go test ./...", Output: "\x1b[32mok\x1b[0m  mujicaui  3.2s", ExitCode: 0, Duration: 3200 * time.Millisecond}
}

func TestCommandBlockToggleAndCopy(t *testing.T) {
	open, toggles := false, 0
	var copied []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := CommandBlock(c, commandBlockSample(), &open, CommandBlockOptions{})
		if r.Toggled() {
			toggles++
		}
		if w, ok := r.Copied(); ok {
			copied = append(copied, w)
		}
	}, 600, 300)
	if tt.HasText("mujicaui  3.2s") {
		t.Fatal("closed block shows output")
	}
	if !tt.HasText("exit 0") || !tt.HasText("3.2 s") {
		t.Fatalf("status: %q", tt.Texts())
	}
	if err := tt.Click("$ go test ./..."); err != nil {
		t.Fatal(err)
	}
	if !open || toggles != 1 || !tt.HasText("ok  mujicaui  3.2s") {
		t.Fatalf("open %v toggles %d texts %q", open, toggles, tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if open || toggles != 2 {
		t.Fatalf("Enter on the focused header: open %v toggles %d", open, toggles)
	}
	_ = tt.Click("Copy command")
	if tt.Clipboard() != "go test ./..." {
		t.Fatalf("clipboard %q", tt.Clipboard())
	}
	_ = tt.Click("Copy output")
	if tt.Clipboard() != "ok  mujicaui  3.2s" || len(copied) != 2 || copied[1] != "output" {
		t.Fatalf("output copy %q %q", tt.Clipboard(), copied)
	}
}

func TestCommandBlockFailureRunningDisabled(t *testing.T) {
	open := true
	run := CommandRun{Command: "make", Output: "boom", ExitCode: 2, Duration: 90 * time.Second}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			CommandBlock(c, run, &open, CommandBlockOptions{Prompt: "❯"})
		})
	}, 600, 300)
	if !tt.HasText("exit 2") || !tt.HasText("1m 30s") {
		t.Fatalf("failure: %q", tt.Texts())
	}
	run.Running = true
	tt.Frame()
	if !tt.HasText("Running") || tt.HasText("exit 2") {
		t.Fatalf("running: %q", tt.Texts())
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("❯ make")
	_ = tt.Click("Copy command")
	if !open || tt.Clipboard() != "" {
		t.Fatalf("disabled block acted: open %v clipboard %q", open, tt.Clipboard())
	}
}

func TestCommandBlockInvalid(t *testing.T) {
	for _, run := range []CommandRun{{}, {Command: "x", Duration: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", run)
				}
			}()
			open := false
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				CommandBlock(c, run, &open, CommandBlockOptions{})
			}, 100, 50, 1)
		}()
	}
}

func TestCommandBlockSnapshot(t *testing.T) {
	snapshot(t, "command_block", 620, 300, func(c *ui.Context) {
		a, b := true, false
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			CommandBlock(c, CommandRun{Command: "go test -run Git ./...", Output: "\x1b[32mok\x1b[0m  \tmujicaui\t3.214s\n\x1b[33mwarn\x1b[0m slow test: TestCommitGraphVirtualized", Duration: 3214 * time.Millisecond}, &a, CommandBlockOptions{})
			CommandBlock(c, CommandRun{Command: "git push origin main", ExitCode: 1, Duration: 640 * time.Millisecond, Output: "rejected"}, &b, CommandBlockOptions{})
		})
	})
}
