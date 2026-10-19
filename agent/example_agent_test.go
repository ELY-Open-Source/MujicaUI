package agent_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/core"
)

func ExampleToolCallCard() {
	call := agent.ToolCall{ID: "1", Name: "read_file", Args: "{\"path\":\"ledger.md\"}",
		Result: "{\"lines\":40}", State: agent.AgentDone, Duration: 420 * time.Millisecond}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.ToolCallCard(c, call, agent.ToolCallCardOptions{OpenArgs: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolCallGroup() {
	open := false
	calls := []agent.ToolCall{
		{ID: "a", Name: "list_dir", State: agent.AgentDone, Duration: 300 * time.Millisecond},
		{ID: "b", Name: "read_file", State: agent.AgentFailed, Error: "not found"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.ToolCallGroup(c, &open, calls, agent.ToolCallGroupOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolApprovalDialog() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := agent.ToolApprovalDialog(c, &open, agent.ToolApprovalOptions{
			Tool: "delete_files", Args: "{\"glob\":\"*.tmp\"}", Reason: "Clean the build folder.",
		})
		if r.Allowed {
			fmt.Println("run the tool")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePermissionPrompt() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := agent.PermissionPrompt(c, &open, agent.PermissionPromptOptions{
			Title: "Write to the ledger?", Scope: "Write access to ~/court/ledger",
		})
		if r.Choice == agent.PermissionAllowAlways {
			fmt.Println("remember the grant")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgentStepList() {
	steps := []agent.AgentStep{
		{ID: "read", Title: "Read the guest list", State: agent.AgentDone, Duration: time.Second},
		{ID: "seat", Title: "Assign the seats", State: agent.AgentFailed, Detail: "Table 4 is full."},
		{ID: "send", Title: "Send invitations", State: agent.AgentPending},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := agent.AgentStepList(c, steps, agent.AgentStepListOptions{}); r.Retry != "" {
			fmt.Println("retry", r.Retry)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgentPlan() {
	items := []agent.PlanItem{
		{ID: "a", Text: "Read the guest list", State: agent.AgentDone},
		{ID: "b", Text: "Assign the seats", State: agent.AgentRunning},
		{ID: "c", Text: "Order the flowers"},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.AgentPlan(c, &items, agent.AgentPlanOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgentStatus() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.AgentStatus(c, agent.AgentWaiting, agent.AgentStatusOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgentProgress() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.AgentProgress(c, 0.4, agent.AgentProgressOptions{Action: "Reading ledger.md"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSubAgentTree() {
	agents := []agent.SubAgent{{ID: "lead", Name: "Steward", State: agent.AgentRunning, Children: []agent.SubAgent{
		{ID: "cook", Name: "Cook", Task: "Plan the menu", State: agent.AgentWaiting},
	}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.SubAgentTree(c, agents, agent.SubAgentTreeOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFileChangeCard() {
	decision := agent.FilePending
	change := agent.FileChange{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if agent.FileChangeCard(c, &decision, change, agent.FileChangeCardOptions{}).OpenDiff {
			fmt.Println("open the full diff")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMultiFileDiffReview() {
	files := []agent.FileChange{
		{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"},
		{Path: "menu.md", Old: "soup\n", New: "soup\ntart\n"},
	}
	decisions := make([]agent.FileDecision, len(files))
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.MultiFileDiffReview(c, &decisions, files, agent.MultiFileDiffReviewOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandExecutionCard() {
	run := agent.CommandRun{Command: "go test ./...", Dir: "~/court", Running: true, Output: "\x1b[32mok\x1b[0m court\n"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if agent.CommandExecutionCard(c, run, agent.CommandExecutionCardOptions{}).Changed() {
			fmt.Println("abort the process")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScreenshotStream() {
	var st agent.ScreenshotStreamState
	var frames []agent.ScreenshotFrame // appended as the agent captures them
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.ScreenshotStream(c, &st, frames, agent.ScreenshotStreamOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleArtifactPanel() {
	art := agent.Artifact{Title: "seating.go", Language: "go", Content: "package court\n"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if agent.ArtifactPanel(c, art, agent.ArtifactPanelOptions{}).Exported {
			fmt.Println("save the file")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleArtifactVersionSwitcher() {
	current := 1
	versions := []agent.ArtifactVersion{{Label: "v1", Note: "First draft"}, {Label: "v2", Note: "Added music"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if agent.ArtifactVersionSwitcher(c, &current, versions, agent.ArtifactVersionSwitcherOptions{}).Compare {
			fmt.Println("diff v1 and v2")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSandboxStatus() {
	info := agent.SandboxInfo{Name: "sandbox-7", State: agent.SandboxRunning,
		Resources: []agent.SandboxResource{{Name: "CPU", Used: 1.5, Limit: 2, Unit: "cores"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		agent.SandboxStatus(c, info, agent.SandboxStatusOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMCPServerList() {
	servers := []agent.MCPServer{{ID: "fs", Name: "filesystem", Transport: "stdio", State: agent.MCPConnected, Tools: 12}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := agent.MCPServerList(c, servers, agent.MCPServerListOptions{}); r.Action == agent.MCPStop {
			fmt.Println("stop", r.Server)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMemoryPanel() {
	entries := []agent.MemoryEntry{{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := agent.MemoryPanel(c, &entries, agent.MemoryPanelOptions{}); r.Deleted != "" {
			fmt.Println("forget", r.Deleted)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCheckpointList() {
	current := "c2"
	points := []agent.Checkpoint{{ID: "c2", Label: "Menu drafted"}, {ID: "c1", Label: "Start"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if agent.CheckpointList(c, &current, points, agent.CheckpointListOptions{}).Changed() {
			fmt.Println("restore", current)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHumanInputRequest() {
	answer := ""
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := agent.HumanInputRequest(c, &answer, agent.HumanInputRequestOptions{
			Question: "Which hall?", Options: []string{"East wing", "Great hall"}, AllowFree: true,
		})
		if r.Submitted {
			fmt.Println("reply", answer)
		}
	}
	ui.Render(view, 400, 300, 1)
}
