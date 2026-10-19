package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

const agentCategory = "Agent"

var agentPages = []page{
	{ID: "261", Name: "ToolCallCard", ZhName: "工具调用卡片", Category: agentCategory, Demo: toolCallCardDemo, Code: toolCallCardCode,
		Keys: "Tab 到参数、结果折叠按钮，Enter / Space 展开或折叠。"},
	{ID: "262", Name: "ToolCallGroup", ZhName: "工具调用组", Category: agentCategory, Demo: toolCallGroupDemo, Code: toolCallGroupCode,
		Keys: "Tab 到组标题，Enter / Space 展开或折叠。"},
	{ID: "263", Name: "ToolApprovalDialog", ZhName: "工具审批对话框", Category: agentCategory, Demo: toolApprovalDemo, Code: toolApprovalDialogCode,
		Keys: "打开时焦点在拒绝；Enter 只触发聚焦按钮，Tab 到允许才能批准；Escape 拒绝。"},
	{ID: "264", Name: "PermissionPrompt", ZhName: "权限请求", Category: agentCategory, Demo: permissionDemo, Code: permissionPromptCode,
		Keys: "打开时焦点在拒绝；Tab 到始终允许、允许一次；Escape 拒绝。"},
	{ID: "265", Name: "AgentStepList", ZhName: "执行步骤列表", Category: agentCategory, Demo: agentStepListDemo, Code: agentStepListCode,
		Keys: "↑ ↓ 在步骤与重试间移动；Enter 展开详情或重试。"},
	{ID: "266", Name: "AgentPlan", ZhName: "执行计划", Category: agentCategory, Demo: agentPlanDemo, Code: agentPlanCode,
		Keys: "Space 勾选；Tab 到跳过或恢复；在输入框中按 Enter 新增。"},
	{ID: "267", Name: "AgentStatus", ZhName: "Agent 状态", Category: agentCategory, Demo: agentStatusDemo, Code: agentStatusCode,
		Keys: "不获取焦点；屏幕阅读器读出状态文字。"},
	{ID: "268", Name: "AgentProgress", ZhName: "Agent 进度", Category: agentCategory, Demo: agentProgressDemo, Code: agentProgressCode,
		Keys: "不获取焦点；当前动作以状态读出。"},
	{ID: "269", Name: "SubAgentTree", ZhName: "子 Agent 树", Category: agentCategory, Demo: subAgentTreeDemo, Code: subAgentTreeCode,
		Keys: "↑ ↓ 移动；→ 展开，← 折叠；Enter / Space 切换。"},
	{ID: "270", Name: "FileChangeCard", ZhName: "文件变更卡片", Category: agentCategory, Demo: fileChangeCardDemo, Code: fileChangeCardCode,
		Keys: "Tab 到打开差异、拒绝、接受或撤销，Enter / Space 触发。"},
	{ID: "271", Name: "MultiFileDiffReview", ZhName: "多文件差异审阅", Category: agentCategory, Demo: multiFileDiffDemo, Code: multiFileDiffReviewCode,
		Keys: "文件列表中 ↑ ↓ 移动，Enter 选择；Tab 到接受、拒绝、全部接受。"},
	{ID: "272", Name: "CommandExecutionCard", ZhName: "命令执行卡片", Category: agentCategory, Demo: commandExecutionDemo, Code: commandExecutionCardCode,
		Keys: "Tab 到中止，Enter / Space 请求中止；输出区可滚动。"},
	{ID: "273", Name: "ScreenshotStream", ZhName: "Agent 画面流", Category: agentCategory, Demo: screenshotStreamDemo, Code: screenshotStreamCode,
		Keys: "画面区中 ← → 逐帧，Space 暂停或恢复。"},
	{ID: "274", Name: "ArtifactPanel", ZhName: "产物面板", Category: agentCategory, Demo: artifactPanelDemo, Code: artifactPanelCode,
		Keys: "切换组中 ← → 移动、Enter 选择；Tab 到复制、导出。"},
	{ID: "275", Name: "ArtifactVersionSwitcher", ZhName: "产物版本切换", Category: agentCategory, Demo: artifactVersionDemo, Code: artifactVersionSwitcherCode,
		Keys: "焦点在组件内时 ← → 切换版本；Tab 到比较按钮。"},
	{ID: "276", Name: "SandboxStatus", ZhName: "沙箱状态", Category: agentCategory, Demo: sandboxStatusDemo, Code: sandboxStatusCode,
		Keys: "不获取焦点；剩余时间每秒刷新。"},
	{ID: "277", Name: "MCPServerList", ZhName: "MCP 服务器列表", Category: agentCategory, Demo: mcpServerListDemo, Code: mcpServerListCode,
		Keys: "Tab 到各行的启动、停止、重新连接，Enter / Space 触发。"},
	{ID: "278", Name: "MemoryPanel", ZhName: "记忆面板", Category: agentCategory, Demo: memoryPanelDemo, Code: memoryPanelCode,
		Keys: "Tab 到搜索、编辑、删除；删除确认默认聚焦取消，Escape 取消。"},
	{ID: "279", Name: "CheckpointList", ZhName: "检查点列表", Category: agentCategory, Demo: checkpointListDemo, Code: checkpointListCode,
		Keys: "Tab 到回滚，Enter 打开确认；确认框默认聚焦取消。"},
	{ID: "280", Name: "HumanInputRequest", ZhName: "人工输入请求", Category: agentCategory, Demo: humanInputDemo, Code: humanInputRequestCode,
		Keys: "方向键选择选项；输入框中 Enter 发送；Tab 到发送回答。"},
}

// agentStates are the states a demo can pick from.
var agentStates = []string{"Pending", "Planning", "Running", "Waiting", "Done", "Failed", "Skipped"}

// agentStatePick is a select over agentStates, logging changes.
func agentStatePick(c *ui.Context, log func(string), v *string) agent.AgentState {
	if ui.Select(c, v, agentStates).Changed() {
		log("state = " + *v)
	}
	for i, s := range agentStates {
		if s == *v {
			return agent.AgentState(i)
		}
	}
	panic("gallery: unknown state " + *v)
}

type toolCallDemoState struct {
	state  string
	failed bool
	still  bool
}

func toolCallCardDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "261", func() toolCallDemoState { return toolCallDemoState{state: "Done"} })
	var st agent.AgentState
	overlayControls(c, func() {
		st = agentStatePick(c, log, &s.state)
		overlayCheck(c, log, &s.failed, "Error")
		overlayCheck(c, log, &s.still, "Still")
	})
	call := agent.ToolCall{ID: "1", Name: "read_file", Args: `{"path":"court/ledger.md","limit":40}`,
		Result: `{"lines":40,"truncated":false}`, State: st, Duration: 420 * time.Millisecond}
	if s.failed {
		call.Error = "permission denied: court/ledger.md"
	}
	if agent.ToolCallCard(c, call, agent.ToolCallCardOptions{Still: s.still}).Changed() {
		log("folded or unfolded a section")
	}
}

func toolCallGroupDemo(c *ui.Context, log func(string)) {
	open := ui.Local(c.Root(), "262", func() bool { return true })
	calls := []agent.ToolCall{
		{ID: "a", Name: "list_dir", Args: `{"path":"court"}`, Result: `["ledger.md","menu.md"]`, State: agent.AgentDone, Duration: 300 * time.Millisecond},
		{ID: "b", Name: "read_file", Args: `{"path":"court/wine.md"}`, Error: "not found", State: agent.AgentFailed, Duration: 1200 * time.Millisecond},
		{ID: "c", Name: "grep", Args: `{"pattern":"Sakiko"}`, Result: `{"matches":3}`, State: agent.AgentDone, Duration: 500 * time.Millisecond},
	}
	if agent.ToolCallGroup(c, open, calls, agent.ToolCallGroupOptions{}).Changed() {
		log(fmt.Sprintf("group open = %v", *open))
	}
}

func toolApprovalDemo(c *ui.Context, log func(string)) {
	open := ui.Local(c.Root(), "263", func() bool { return false })
	if input.Button(c, "Run delete_files", input.ButtonOptions{}).Clicked() {
		*open = true
		log("asked for approval")
	}
	r := agent.ToolApprovalDialog(c, open, agent.ToolApprovalOptions{Tool: "delete_files",
		Args: `{"glob":"build/*.tmp","dry_run":false}`, Reason: "The agent wants to clean temporary files from the build folder."})
	if r.Allowed {
		log("allowed")
	}
	if r.Denied {
		log("denied")
	}
}

type permissionDemoState struct {
	open  bool
	grant string
}

func permissionDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "264", func() permissionDemoState { return permissionDemoState{grant: "none"} })
	ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
		if input.Button(c, "Ask to write", input.ButtonOptions{}).Clicked() {
			s.open = true
			log("asked for permission")
		}
		ui.Text(c, "Grant: "+s.grant)
	})
	r := agent.PermissionPrompt(c, &s.open, agent.PermissionPromptOptions{Title: "Write to the ledger?",
		Description: "The agent wants to record tonight's guests.", Scope: "Write access to ~/court/ledger and its subfolders"})
	if r.Choice != agent.PermissionNone {
		s.grant = [...]string{"none", "once", "always", "denied"}[r.Choice]
		log("choice = " + s.grant)
	}
}

type stepDemoState struct {
	steps    []agent.AgentStep
	disabled bool
}

func agentStepListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "265", func() stepDemoState {
		return stepDemoState{steps: []agent.AgentStep{
			{ID: "read", Title: "Read the guest list", State: agent.AgentDone, Duration: 800 * time.Millisecond, Detail: "42 guests, 3 arrive late."},
			{ID: "seat", Title: "Assign the seats", State: agent.AgentFailed, Duration: 2 * time.Second, Detail: "Table 4 is over capacity."},
			{ID: "send", Title: "Send the invitations", State: agent.AgentPending},
		}}
	})
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	r := agent.AgentStepList(c, s.steps, agent.AgentStepListOptions{Disabled: s.disabled})
	if r.Toggled != "" {
		log("toggled " + r.Toggled)
	}
	if r.Retry != "" {
		log("retry " + r.Retry)
		for i := range s.steps {
			if s.steps[i].ID == r.Retry {
				s.steps[i].State, s.steps[i].Detail = agent.AgentRunning, "Retrying with a larger table."
			}
		}
	}
}

type planDemoState struct {
	items    []agent.PlanItem
	disabled bool
}

func agentPlanDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "266", func() planDemoState {
		return planDemoState{items: []agent.PlanItem{
			{ID: "a", Text: "Read the guest list", State: agent.AgentDone},
			{ID: "b", Text: "Assign the seats", State: agent.AgentRunning},
			{ID: "c", Text: "Order the flowers"},
		}}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if input.Button(c, "Agent finishes a step", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			for i := range s.items {
				if s.items[i].State != agent.AgentDone && s.items[i].State != agent.AgentSkipped {
					s.items[i].State = agent.AgentDone
					log("agent checked " + s.items[i].Text)
					break
				}
			}
		}
	})
	if agent.AgentPlan(c, &s.items, agent.AgentPlanOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("plan changed (%d items)", len(s.items)))
	}
}

type statusDemoState struct {
	state string
	still bool
}

func agentStatusDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "267", func() statusDemoState { return statusDemoState{state: "Running"} })
	var st agent.AgentState
	overlayControls(c, func() {
		st = agentStatePick(c, log, &s.state)
		overlayCheck(c, log, &s.still, "Still")
	})
	ui.Row(c).Gap(8).Wrap().Children(func() {
		agent.AgentStatus(c, st, agent.AgentStatusOptions{Still: s.still})
		for _, v := range []agent.AgentState{agent.AgentPlanning, agent.AgentRunning, agent.AgentWaiting, agent.AgentDone, agent.AgentFailed} {
			agent.AgentStatus(c, v, agent.AgentStatusOptions{Still: s.still})
		}
	})
}

type progressDemoState struct {
	value         float64
	indeterminate bool
}

func agentProgressDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "268", func() progressDemoState { return progressDemoState{value: 0.4} })
	overlayControls(c, func() {
		if ui.Slider(c, &s.value, 0, 1).Width(200).Changed() {
			log(fmt.Sprintf("value = %.2f", s.value))
		}
		overlayCheck(c, log, &s.indeterminate, "Indeterminate")
	})
	action := "Reading ledger.md"
	if s.value > 0.5 {
		action = "Writing the seating chart"
	}
	agent.AgentProgress(c, s.value, agent.AgentProgressOptions{Label: "Preparing the ball", Action: action, Indeterminate: s.indeterminate})
}

func subAgentTreeDemo(c *ui.Context, log func(string)) {
	disabled := ui.Local(c.Root(), "269", func() bool { return false })
	overlayControls(c, func() { overlayCheck(c, log, disabled, "Disabled") })
	agents := []agent.SubAgent{{ID: "lead", Name: "Steward", Task: "Prepare the ball", State: agent.AgentRunning, Children: []agent.SubAgent{
		{ID: "music", Name: "Musician", Task: "Hire the quartet", State: agent.AgentDone},
		{ID: "kitchen", Name: "Cook", Task: "Plan the menu", State: agent.AgentWaiting, Children: []agent.SubAgent{
			{ID: "wine", Name: "Sommelier", Task: "Choose the wine", State: agent.AgentFailed},
		}},
	}}}
	if agent.SubAgentTree(c, agents, agent.SubAgentTreeOptions{Disabled: *disabled}).Changed() {
		log("expanded or collapsed an agent")
	}
}

type fileChangeDemoState struct {
	decision agent.FileDecision
	disabled bool
	preview  float64
}

func fileChangeCardDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "270", func() fileChangeDemoState { return fileChangeDemoState{preview: 8} })
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		overlaySlider(c, log, &s.preview, "Preview", "lines", 2, 12, 1)
	})
	r := agent.FileChangeCard(c, &s.decision, agentDemoChange, agent.FileChangeCardOptions{Preview: int(s.preview), Disabled: s.disabled})
	if r.Decided {
		log(fmt.Sprintf("decision = %d", s.decision))
	}
	if r.OpenDiff {
		log("open the full diff")
	}
}

var agentDemoChange = agent.FileChange{Path: "court/ledger.md",
	Old: "# Ledger\nGuests: 40\nWine: claret\nMusic: none\nHall: east wing\n",
	New: "# Ledger\nGuests: 42\nWine: claret\nMusic: quartet\nFlowers: lilies\nHall: east wing\n"}
