package main

import (
	"fmt"
	"image"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/agent"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

type multiFileDemoState struct {
	decisions []agent.FileDecision
	disabled  bool
}

var agentDemoFiles = []agent.FileChange{
	agentDemoChange,
	{Path: "court/menu.md", Old: "soup\nroast\n", New: "soup\nroast\ntart\n"},
	{Path: "court/old-plan.md", Old: "Hold the ball in the cellar.\n"},
}

func multiFileDiffDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "271", func() multiFileDemoState {
		return multiFileDemoState{decisions: make([]agent.FileDecision, len(agentDemoFiles))}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if input.Button(c, "Reset", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			clear(s.decisions)
			log("reset")
		}
	})
	if agent.MultiFileDiffReview(c, &s.decisions, agentDemoFiles, agent.MultiFileDiffReviewOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("decisions = %v", s.decisions))
	}
}

type commandDemoState struct {
	run      agent.CommandRun
	disabled bool
	lines    int
}

func commandExecutionDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "272", func() commandDemoState {
		return commandDemoState{run: agent.CommandRun{Command: "go test ./court/...", Dir: "~/projects/court", Running: true}}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if input.Button(c, "Print a line", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.lines++
			color := []string{"32mok  ", "33mwarn", "31mFAIL"}[s.lines%3]
			s.run.Output += fmt.Sprintf("\x1b[%s\x1b[0m court/case%d  0.%02ds\n", color, s.lines, s.lines)
			log("printed a line")
		}
		if input.Button(c, "Restart", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.run.Running, s.run.Output, s.lines = true, "", 0
			log("restarted")
		}
	})
	if agent.CommandExecutionCard(c, s.run, agent.CommandExecutionCardOptions{Disabled: s.disabled}).Changed() {
		s.run.Running, s.run.ExitCode, s.run.Duration = false, 130, time.Duration(s.lines)*time.Second
		s.run.Output += "\x1b[31maborted by user\x1b[0m\n"
		log("abort requested")
	}
}

type screenshotDemoState struct {
	st       agent.ScreenshotStreamState
	frames   []agent.ScreenshotFrame
	disabled bool
}

// agentDemoFrame paints a frame with a moving block.
func agentDemoFrame(i int) agent.ScreenshotFrame {
	img := image.NewRGBA(image.Rect(0, 0, 320, 180))
	for y := range 180 {
		for x := range 320 {
			o := img.PixOffset(x, y)
			v := uint8(235)
			if x >= 20+i*24%260 && x < 60+i*24%260 && y > 60 && y < 120 {
				v = 120
			}
			img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3] = v, v-10, v-20, 255
		}
	}
	return agent.ScreenshotFrame{Image: ui.NewBitmap(img), Caption: fmt.Sprintf("Clicked button %d", i+1), At: time.Duration(i) * 2 * time.Second}
}

func screenshotStreamDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "273", func() screenshotDemoState {
		return screenshotDemoState{frames: []agent.ScreenshotFrame{agentDemoFrame(0), agentDemoFrame(1)}}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if input.Button(c, "Capture a frame", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.frames = append(s.frames, agentDemoFrame(len(s.frames)))
			log(fmt.Sprintf("frame %d arrived", len(s.frames)))
		}
	})
	if agent.ScreenshotStream(c, &s.st, s.frames, agent.ScreenshotStreamOptions{Height: 220, Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("frame %d, paused = %v", s.st.Index+1, s.st.Paused))
	}
}

type artifactDemoState struct {
	kind     string
	disabled bool
}

func artifactPanelDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "274", func() artifactDemoState { return artifactDemoState{kind: "Code"} })
	overlayControls(c, func() {
		if ui.Select(c, &s.kind, []string{"Code", "Document"}).Changed() {
			log("kind = " + s.kind)
		}
		overlayCheck(c, log, &s.disabled, "Disabled")
	})
	art := agent.Artifact{Title: "seating.go", Language: "go",
		Content: "package court\n\n// Seats maps guests to tables.\nfunc Seats() map[string]int {\n\treturn map[string]int{\"Sakiko\": 1, \"Uika\": 2}\n}\n"}
	opts := agent.ArtifactPanelOptions{Height: 200, Disabled: s.disabled}
	if s.kind == "Document" {
		art = agent.Artifact{Title: "Invitation", Kind: agent.ArtifactDocument, Content: "# You are invited\n\nThe masquerade begins at nine.\n"}
		opts.Preview = func() {
			display.Heading(c, "You are invited", display.HeadingOptions{})
			ui.Text(c, "The masquerade begins at nine.")
		}
	}
	r := agent.ArtifactPanel(c, art, opts)
	if r.Copied {
		log("copied to the clipboard")
	}
	if r.Exported {
		log("export requested")
	}
}

type versionDemoState struct {
	current  int
	disabled bool
}

func artifactVersionDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "275", func() versionDemoState { return versionDemoState{current: 2} })
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	versions := []agent.ArtifactVersion{{Label: "v1", Note: "First draft"}, {Label: "v2", Note: "Added the quartet"}, {Label: "v3", Note: "Moved the wine table"}}
	r := agent.ArtifactVersionSwitcher(c, &s.current, versions, agent.ArtifactVersionSwitcherOptions{Disabled: s.disabled})
	if r.Changed {
		log("version = " + versions[s.current].Label)
	}
	if r.Compare {
		log("compare " + versions[s.current-1].Label + " and " + versions[s.current].Label)
	}
}

type sandboxDemoState struct {
	state   string
	cpu     float64
	expires time.Time
}

func sandboxStatusDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "276", func() sandboxDemoState {
		return sandboxDemoState{state: "Running", cpu: 1.2, expires: time.Now().Add(10 * time.Minute)}
	})
	names := []string{"Starting", "Running", "Stopped", "Failed"}
	state := agent.SandboxStarting
	overlayControls(c, func() {
		if ui.Select(c, &s.state, names).Changed() {
			log("state = " + s.state)
		}
		overlaySlider(c, log, &s.cpu, "CPU", "cores", 0, 2, 0.1)
		if input.Button(c, "Extend 5 min", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.expires = s.expires.Add(5 * time.Minute)
			log("extended")
		}
	})
	for i, n := range names {
		if n == s.state {
			state = agent.SandboxState(i)
		}
	}
	agent.SandboxStatus(c, agent.SandboxInfo{Name: "sandbox-7", State: state, Expires: s.expires, Resources: []agent.SandboxResource{
		{Name: "CPU", Used: s.cpu, Limit: 2, Unit: "cores"}, {Name: "Memory", Used: 3.1, Limit: 4, Unit: "GB"},
	}}, agent.SandboxStatusOptions{})
}

type mcpDemoState struct {
	servers  []agent.MCPServer
	disabled bool
}

func mcpServerListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "277", func() mcpDemoState {
		return mcpDemoState{servers: []agent.MCPServer{
			{ID: "fs", Name: "filesystem", Transport: "stdio", State: agent.MCPConnected, Tools: 12},
			{ID: "gh", Name: "github", Transport: "http", State: agent.MCPError, Error: "401 Unauthorized"},
			{ID: "db", Name: "postgres", Transport: "stdio", State: agent.MCPStopped},
		}}
	})
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	r := agent.MCPServerList(c, s.servers, agent.MCPServerListOptions{Disabled: s.disabled})
	if r.Action == agent.MCPNone {
		return
	}
	log(fmt.Sprintf("%s: %s", r.Server, [...]string{"", "start", "stop", "reconnect"}[r.Action]))
	for i := range s.servers {
		if s.servers[i].ID != r.Server {
			continue
		}
		switch r.Action {
		case agent.MCPStop:
			s.servers[i].State = agent.MCPStopped
		default:
			s.servers[i].State, s.servers[i].Tools, s.servers[i].Error = agent.MCPConnected, 5, ""
		}
	}
}

type memoryDemoState struct {
	entries  []agent.MemoryEntry
	disabled bool
}

func memoryPanelDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "278", func() memoryDemoState {
		return memoryDemoState{entries: []agent.MemoryEntry{
			{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat, 3 Oct"},
			{ID: "2", Text: "The quartet needs two days' notice.", Source: "Plan, 1 Oct"},
			{ID: "3", Text: "Claret is served after nine.", Source: "Ledger"},
		}}
	})
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	r := agent.MemoryPanel(c, &s.entries, agent.MemoryPanelOptions{Disabled: s.disabled})
	if r.Edited != "" {
		log("edited " + r.Edited)
	}
	if r.Deleted != "" {
		log("deleted " + r.Deleted)
	}
}

type checkpointDemoState struct {
	current  string
	disabled bool
}

func checkpointListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "279", func() checkpointDemoState { return checkpointDemoState{current: "c3"} })
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	points := []agent.Checkpoint{{ID: "c3", Label: "Seats assigned", Note: "12 files changed"},
		{ID: "c2", Label: "Menu drafted", Note: "3 files changed"}, {ID: "c1", Label: "Start"}}
	if agent.CheckpointList(c, &s.current, points, agent.CheckpointListOptions{Disabled: s.disabled}).Changed() {
		log("rolled back to " + s.current)
	}
}

type humanDemoState struct {
	answer   string
	free     bool
	deadline time.Time
	disabled bool
}

func humanInputDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "280", func() humanDemoState { return humanDemoState{free: true} })
	overlayControls(c, func() {
		overlayCheck(c, log, &s.free, "AllowFree")
		overlayCheck(c, log, &s.disabled, "Disabled")
		if input.Button(c, "Start a 30 s timer", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.deadline = time.Now().Add(30 * time.Second)
			log("timer started")
		}
	})
	r := agent.HumanInputRequest(c, &s.answer, agent.HumanInputRequestOptions{Question: "Which hall should host the masquerade?",
		Options: []string{"East wing", "Great hall"}, AllowFree: s.free, Deadline: s.deadline, Disabled: s.disabled})
	if r.Submitted {
		log("answer sent: " + s.answer)
	}
	if r.TimedOut {
		log("timed out")
	}
}
