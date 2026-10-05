package main

// Code snippets of the Agent pages, the same as their compiled Examples.

const toolCallCardCode = `call := mujica.ToolCall{ID: "1", Name: "read_file", Args: "{\"path\":\"ledger.md\"}",
	Result: "{\"lines\":40}", State: mujica.AgentDone, Duration: 420 * time.Millisecond}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ToolCallCard(c, call, mujica.ToolCallCardOptions{OpenArgs: true})
}
ui.Render(view, 400, 300, 1)`

const toolCallGroupCode = `open := false
calls := []mujica.ToolCall{
	{ID: "a", Name: "list_dir", State: mujica.AgentDone, Duration: 300 * time.Millisecond},
	{ID: "b", Name: "read_file", State: mujica.AgentFailed, Error: "not found"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ToolCallGroup(c, &open, calls, mujica.ToolCallGroupOptions{})
}
ui.Render(view, 400, 300, 1)`

const toolApprovalDialogCode = `open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.ToolApprovalDialog(c, &open, mujica.ToolApprovalOptions{
		Tool: "delete_files", Args: "{\"glob\":\"*.tmp\"}", Reason: "Clean the build folder.",
	})
	if r.Allowed {
		fmt.Println("run the tool")
	}
}
ui.Render(view, 400, 300, 1)`

const permissionPromptCode = `open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.PermissionPrompt(c, &open, mujica.PermissionPromptOptions{
		Title: "Write to the ledger?", Scope: "Write access to ~/court/ledger",
	})
	if r.Choice == mujica.PermissionAllowAlways {
		fmt.Println("remember the grant")
	}
}
ui.Render(view, 400, 300, 1)`

const agentStepListCode = `steps := []mujica.AgentStep{
	{ID: "read", Title: "Read the guest list", State: mujica.AgentDone, Duration: time.Second},
	{ID: "seat", Title: "Assign the seats", State: mujica.AgentFailed, Detail: "Table 4 is full."},
	{ID: "send", Title: "Send invitations", State: mujica.AgentPending},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.AgentStepList(c, steps, mujica.AgentStepListOptions{}); r.Retry != "" {
		fmt.Println("retry", r.Retry)
	}
}
ui.Render(view, 400, 300, 1)`

const agentPlanCode = `items := []mujica.PlanItem{
	{ID: "a", Text: "Read the guest list", State: mujica.AgentDone},
	{ID: "b", Text: "Assign the seats", State: mujica.AgentRunning},
	{ID: "c", Text: "Order the flowers"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AgentPlan(c, &items, mujica.AgentPlanOptions{})
}
ui.Render(view, 400, 300, 1)`

const agentStatusCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AgentStatus(c, mujica.AgentWaiting, mujica.AgentStatusOptions{})
}
ui.Render(view, 400, 300, 1)`

const agentProgressCode = `view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AgentProgress(c, 0.4, mujica.AgentProgressOptions{Action: "Reading ledger.md"})
}
ui.Render(view, 400, 300, 1)`

const subAgentTreeCode = `agents := []mujica.SubAgent{{ID: "lead", Name: "Steward", State: mujica.AgentRunning, Children: []mujica.SubAgent{
	{ID: "cook", Name: "Cook", Task: "Plan the menu", State: mujica.AgentWaiting},
}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SubAgentTree(c, agents, mujica.SubAgentTreeOptions{})
}
ui.Render(view, 400, 300, 1)`

const fileChangeCardCode = `decision := mujica.FilePending
change := mujica.FileChange{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.FileChangeCard(c, &decision, change, mujica.FileChangeCardOptions{}).OpenDiff {
		fmt.Println("open the full diff")
	}
}
ui.Render(view, 400, 300, 1)`

const multiFileDiffReviewCode = `files := []mujica.FileChange{
	{Path: "ledger.md", Old: "Guests: 40\n", New: "Guests: 42\n"},
	{Path: "menu.md", Old: "soup\n", New: "soup\ntart\n"},
}
decisions := make([]mujica.FileDecision, len(files))
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MultiFileDiffReview(c, &decisions, files, mujica.MultiFileDiffReviewOptions{})
}
ui.Render(view, 400, 300, 1)`

const commandExecutionCardCode = `run := mujica.CommandRun{Command: "go test ./...", Dir: "~/court", Running: true, Output: "\x1b[32mok\x1b[0m court\n"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.CommandExecutionCard(c, run, mujica.CommandExecutionCardOptions{}).Changed() {
		fmt.Println("abort the process")
	}
}
ui.Render(view, 400, 300, 1)`

const screenshotStreamCode = `var st mujica.ScreenshotStreamState
var frames []mujica.ScreenshotFrame // appended as the agent captures them
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ScreenshotStream(c, &st, frames, mujica.ScreenshotStreamOptions{})
}
ui.Render(view, 400, 300, 1)`

const artifactPanelCode = `art := mujica.Artifact{Title: "seating.go", Language: "go", Content: "package court\n"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ArtifactPanel(c, art, mujica.ArtifactPanelOptions{}).Exported {
		fmt.Println("save the file")
	}
}
ui.Render(view, 400, 300, 1)`

const artifactVersionSwitcherCode = `current := 1
versions := []mujica.ArtifactVersion{{Label: "v1", Note: "First draft"}, {Label: "v2", Note: "Added music"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ArtifactVersionSwitcher(c, &current, versions, mujica.ArtifactVersionSwitcherOptions{}).Compare {
		fmt.Println("diff v1 and v2")
	}
}
ui.Render(view, 400, 300, 1)`

const sandboxStatusCode = `info := mujica.SandboxInfo{Name: "sandbox-7", State: mujica.SandboxRunning,
	Resources: []mujica.SandboxResource{{Name: "CPU", Used: 1.5, Limit: 2, Unit: "cores"}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SandboxStatus(c, info, mujica.SandboxStatusOptions{})
}
ui.Render(view, 400, 300, 1)`

const mcpServerListCode = `servers := []mujica.MCPServer{{ID: "fs", Name: "filesystem", Transport: "stdio", State: mujica.MCPConnected, Tools: 12}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.MCPServerList(c, servers, mujica.MCPServerListOptions{}); r.Action == mujica.MCPStop {
		fmt.Println("stop", r.Server)
	}
}
ui.Render(view, 400, 300, 1)`

const memoryPanelCode = `entries := []mujica.MemoryEntry{{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.MemoryPanel(c, &entries, mujica.MemoryPanelOptions{}); r.Deleted != "" {
		fmt.Println("forget", r.Deleted)
	}
}
ui.Render(view, 400, 300, 1)`

const checkpointListCode = `current := "c2"
points := []mujica.Checkpoint{{ID: "c2", Label: "Menu drafted"}, {ID: "c1", Label: "Start"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.CheckpointList(c, &current, points, mujica.CheckpointListOptions{}).Changed() {
		fmt.Println("restore", current)
	}
}
ui.Render(view, 400, 300, 1)`

const humanInputRequestCode = `answer := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.HumanInputRequest(c, &answer, mujica.HumanInputRequestOptions{
		Question: "Which hall?", Options: []string{"East wing", "Great hall"}, AllowFree: true,
	})
	if r.Submitted {
		fmt.Println("reply", answer)
	}
}
ui.Render(view, 400, 300, 1)`
