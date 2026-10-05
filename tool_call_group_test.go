package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var toolCallGroupSample = []ToolCall{
	{ID: "a", Name: "list_dir", State: AgentDone, Duration: 300 * time.Millisecond},
	{ID: "b", Name: "read_file", State: AgentFailed, Error: "not found", Duration: 1200 * time.Millisecond},
	{ID: "c", Name: "grep", State: AgentDone, Duration: 500 * time.Millisecond},
}

func TestToolCallGroupSummaryAndFold(t *testing.T) {
	open, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if ToolCallGroup(c, &open, toolCallGroupSample, ToolCallGroupOptions{}).Changed() {
			changes++
		}
	}, 500, 600)
	for _, s := range []string{"3 tool calls", "2.0 s", "1 failed"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	if tt.HasText("read_file") {
		t.Fatal("calls shown while folded")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || changes != 1 || !tt.HasText("read_file") || !tt.HasText("grep") {
		t.Fatalf("open %v changes %d texts %q", open, changes, tt.Texts())
	}
}

func TestToolCallGroupRunningAndEmpty(t *testing.T) {
	open := false
	calls := []ToolCall{{ID: "x", Name: "fetch", State: AgentRunning}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ToolCallGroup(c, &open, calls, ToolCallGroupOptions{Title: "Research", Still: true})
		ToolCallGroup(c, &open, nil, ToolCallGroupOptions{})
	}, 500, 300)
	if !tt.HasText("Research") || !tt.HasText("0 tool calls") || tt.HasText("1 failed") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestToolCallGroupRejectsDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := false
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ToolCallGroup(c, &open, []ToolCall{{ID: "a", Name: "x"}, {ID: "a", Name: "y"}}, ToolCallGroupOptions{})
	}, 200, 200)
}

func TestToolCallGroupSnapshot(t *testing.T) {
	snapshot(t, "tool_call_group", 460, 330, func(c *ui.Context) {
		open := true
		ui.Column(c).Padding(16).Children(func() {
			ToolCallGroup(c, &open, toolCallGroupSample, ToolCallGroupOptions{})
		})
	})
}
