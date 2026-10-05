package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var toolCallSample = ToolCall{ID: "1", Name: "read_file", Args: `{"path":"court/ledger.md","limit":40}`,
	Result: `{"lines":40,"truncated":false}`, State: AgentDone, Duration: 420 * time.Millisecond}

func TestToolCallCardShowsAndFolds(t *testing.T) {
	call := toolCallSample
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ToolCallCard(c, call, ToolCallCardOptions{}).Changed() {
			changes++
		}
	}, 500, 500)
	for _, s := range []string{"read_file", "Done", "420 ms", "Arguments", "Result"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	if tt.HasText(`  "path": "court/ledger.md",`) {
		t.Fatal("arguments shown while folded")
	}
	_ = tt.Click("Arguments")
	if changes != 1 || !tt.HasText(`  "path": "court/ledger.md",`) {
		t.Fatalf("changes %d, pretty JSON %q", changes, tt.Texts())
	}
	tt.Key(0, ui.KeyTab) // keyboard reaches the Result fold
	tt.Key(0, ui.KeyEnter)
	if changes != 2 || !tt.HasText(`  "lines": 40,`) {
		t.Fatalf("keyboard fold: changes %d", changes)
	}
}

func TestToolCallCardErrorAndPlainText(t *testing.T) {
	call := ToolCall{Name: "shell", Args: "not json at all", Error: "permission denied", State: AgentFailed, Duration: 3 * time.Second}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolCallCard(c, call, ToolCallCardOptions{OpenArgs: true, OpenResult: true})
	}, 500, 400)
	for _, s := range []string{"Failed", "3.0 s", "Error", "permission denied", "not json at all"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
}

func TestToolCallCardRunningHidesDuration(t *testing.T) {
	call := ToolCall{Name: "search", State: AgentRunning, Duration: time.Second}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolCallCard(c, call, ToolCallCardOptions{Still: true})
	}, 400, 200)
	if !tt.HasText("Running") || tt.HasText("1.0 s") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestToolCallCardLargeResultIsVirtualized(t *testing.T) {
	var b strings.Builder
	for i := range 100000 {
		fmt.Fprintf(&b, "row %d\n", i)
	}
	call := ToolCall{Name: "dump", Result: b.String(), State: AgentDone}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolCallCard(c, call, ToolCallCardOptions{OpenResult: true})
	}, 500, 500)
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("%d texts built for 100000 rows", n)
	}
	if !tt.HasText("row 0") {
		t.Fatal("first row missing")
	}
}

func TestToolCallCardNeedsName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolCallCard(c, ToolCall{}, ToolCallCardOptions{})
	}, 200, 200)
}

func TestToolCallCardSnapshot(t *testing.T) {
	snapshot(t, "tool_call_card", 460, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ToolCallCard(c, toolCallSample, ToolCallCardOptions{OpenArgs: true})
		})
	})
}
