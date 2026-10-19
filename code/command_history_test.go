package code

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func commandHistorySample() []HistoryEntry {
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	return []HistoryEntry{
		{Command: "git status", Time: at},
		{Command: "go test ./...", Time: at.Add(time.Minute), ExitCode: 1},
		{Command: "go build", Time: at.Add(2 * time.Minute)},
	}
}

func commandHistoryTester(entries []HistoryEntry, s *CommandHistoryState, disabled bool, got *[]string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := CommandHistory(c, entries, s, CommandHistoryOptions{Disabled: disabled})
		r.Element.Grow(1)
		if r.Replayed != nil {
			*got = append(*got, r.Replayed.Command)
		}
	}, 420, 320)
}

func TestCommandHistoryNewestFirstAndReplay(t *testing.T) {
	s := &CommandHistoryState{}
	var got []string
	tt := commandHistoryTester(commandHistorySample(), s, false, &got)
	first, _ := tt.Find("go build")
	last, _ := tt.Find("git status")
	if first.Y >= last.Y {
		t.Fatal("newest command is not on top")
	}
	if !tt.HasText("exit 1") {
		t.Fatal("failed command shows no exit code")
	}
	if err := tt.Click("Run git status again"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(first.X, first.Y+2)
	tt.ClickAt(first.X, first.Y+2)
	if len(got) != 2 || got[0] != "git status" || got[1] != "go build" {
		t.Fatalf("replayed %q", got)
	}
}

func TestCommandHistorySearchAndKeys(t *testing.T) {
	s := &CommandHistoryState{}
	var got []string
	tt := commandHistoryTester(commandHistorySample(), s, false, &got)
	tt.Key(0, ui.KeyTab)
	tt.Type("go ")
	if tt.HasText("git status") || !tt.HasText("go build") || !tt.HasText("go test ./...") {
		t.Fatalf("filter shows %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 1 || got[0] != "go test ./..." {
		t.Fatalf("keyboard replay %q", got)
	}
	s.Query = "zzz"
	tt.Frame()
	if !tt.HasText("No matching commands") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestCommandHistoryDisabledAndEmpty(t *testing.T) {
	var got []string
	tt := commandHistoryTester(commandHistorySample(), &CommandHistoryState{}, true, &got)
	_ = tt.Click("Run git status again")
	if len(got) != 0 {
		t.Fatalf("disabled replayed %q", got)
	}
	tt = commandHistoryTester(nil, &CommandHistoryState{}, false, &got)
	if !tt.HasText("No commands yet") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestCommandHistoryLargeIsVirtual(t *testing.T) {
	entries := make([]HistoryEntry, 100_000)
	for i := range entries {
		entries[i].Command = fmt.Sprintf("echo %d", i)
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CommandHistory(c, entries, &CommandHistoryState{}, CommandHistoryOptions{})
	}, 420, 320)
	rows := 0
	for _, s := range tt.Texts() {
		if strings.HasPrefix(s, "echo ") {
			rows++
		}
	}
	if !tt.HasText("echo 99999") || rows == 0 || rows > 60 {
		t.Fatalf("built %d rows of 100000", rows)
	}
}

func TestCommandHistorySnapshot(t *testing.T) {
	snapshot(t, "command_history", 420, 240, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *CommandHistoryState { return &CommandHistoryState{} })
		ui.Column(c).Padding(12).Height(240).Children(func() {
			CommandHistory(c, commandHistorySample(), *s, CommandHistoryOptions{}).Element.Grow(1)
		})
	})
}
