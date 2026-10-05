package main

import (
	"fmt"
	"image"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

type multiFileDemoState struct {
	decisions []mujica.FileDecision
	disabled  bool
}

var agentDemoFiles = []mujica.FileChange{
	agentDemoChange,
	{Path: "court/menu.md", Old: "soup\nroast\n", New: "soup\nroast\ntart\n"},
	{Path: "court/old-plan.md", Old: "Hold the ball in the cellar.\n"},
}

func multiFileDiffDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "271", func() multiFileDemoState {
		return multiFileDemoState{decisions: make([]mujica.FileDecision, len(agentDemoFiles))}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if mujica.Button(c, "Reset", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			clear(s.decisions)
			log("reset")
		}
	})
	if mujica.MultiFileDiffReview(c, &s.decisions, agentDemoFiles, mujica.MultiFileDiffReviewOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("decisions = %v", s.decisions))
	}
}

type commandDemoState struct {
	run      mujica.CommandRun
	disabled bool
	lines    int
}

func commandExecutionDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "272", func() commandDemoState {
		return commandDemoState{run: mujica.CommandRun{Command: "go test ./court/...", Dir: "~/projects/court", Running: true}}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if mujica.Button(c, "Print a line", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.lines++
			color := []string{"32mok  ", "33mwarn", "31mFAIL"}[s.lines%3]
			s.run.Output += fmt.Sprintf("\x1b[%s\x1b[0m court/case%d  0.%02ds\n", color, s.lines, s.lines)
			log("printed a line")
		}
		if mujica.Button(c, "Restart", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.run.Running, s.run.Output, s.lines = true, "", 0
			log("restarted")
		}
	})
	if mujica.CommandExecutionCard(c, s.run, mujica.CommandExecutionCardOptions{Disabled: s.disabled}).Changed() {
		s.run.Running, s.run.ExitCode, s.run.Duration = false, 130, time.Duration(s.lines)*time.Second
		s.run.Output += "\x1b[31maborted by user\x1b[0m\n"
		log("abort requested")
	}
}

type screenshotDemoState struct {
	st       mujica.ScreenshotStreamState
	frames   []mujica.ScreenshotFrame
	disabled bool
}

// agentDemoFrame paints a frame with a moving block.
func agentDemoFrame(i int) mujica.ScreenshotFrame {
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
	return mujica.ScreenshotFrame{Image: ui.NewBitmap(img), Caption: fmt.Sprintf("Clicked button %d", i+1), At: time.Duration(i) * 2 * time.Second}
}

func screenshotStreamDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "273", func() screenshotDemoState {
		return screenshotDemoState{frames: []mujica.ScreenshotFrame{agentDemoFrame(0), agentDemoFrame(1)}}
	})
	overlayControls(c, func() {
		overlayCheck(c, log, &s.disabled, "Disabled")
		if mujica.Button(c, "Capture a frame", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.frames = append(s.frames, agentDemoFrame(len(s.frames)))
			log(fmt.Sprintf("frame %d arrived", len(s.frames)))
		}
	})
	if mujica.ScreenshotStream(c, &s.st, s.frames, mujica.ScreenshotStreamOptions{Height: 220, Disabled: s.disabled}).Changed() {
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
	art := mujica.Artifact{Title: "seating.go", Language: "go",
		Content: "package court\n\n// Seats maps guests to tables.\nfunc Seats() map[string]int {\n\treturn map[string]int{\"Sakiko\": 1, \"Uika\": 2}\n}\n"}
	opts := mujica.ArtifactPanelOptions{Height: 200, Disabled: s.disabled}
	if s.kind == "Document" {
		art = mujica.Artifact{Title: "Invitation", Kind: mujica.ArtifactDocument, Content: "# You are invited\n\nThe masquerade begins at nine.\n"}
		opts.Preview = func() {
			mujica.Heading(c, "You are invited", mujica.HeadingOptions{})
			ui.Text(c, "The masquerade begins at nine.")
		}
	}
	r := mujica.ArtifactPanel(c, art, opts)
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
	versions := []mujica.ArtifactVersion{{Label: "v1", Note: "First draft"}, {Label: "v2", Note: "Added the quartet"}, {Label: "v3", Note: "Moved the wine table"}}
	r := mujica.ArtifactVersionSwitcher(c, &s.current, versions, mujica.ArtifactVersionSwitcherOptions{Disabled: s.disabled})
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
	state := mujica.SandboxStarting
	overlayControls(c, func() {
		if ui.Select(c, &s.state, names).Changed() {
			log("state = " + s.state)
		}
		overlaySlider(c, log, &s.cpu, "CPU", "cores", 0, 2, 0.1)
		if mujica.Button(c, "Extend 5 min", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.expires = s.expires.Add(5 * time.Minute)
			log("extended")
		}
	})
	for i, n := range names {
		if n == s.state {
			state = mujica.SandboxState(i)
		}
	}
	mujica.SandboxStatus(c, mujica.SandboxInfo{Name: "sandbox-7", State: state, Expires: s.expires, Resources: []mujica.SandboxResource{
		{Name: "CPU", Used: s.cpu, Limit: 2, Unit: "cores"}, {Name: "Memory", Used: 3.1, Limit: 4, Unit: "GB"},
	}}, mujica.SandboxStatusOptions{})
}

type mcpDemoState struct {
	servers  []mujica.MCPServer
	disabled bool
}

func mcpServerListDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "277", func() mcpDemoState {
		return mcpDemoState{servers: []mujica.MCPServer{
			{ID: "fs", Name: "filesystem", Transport: "stdio", State: mujica.MCPConnected, Tools: 12},
			{ID: "gh", Name: "github", Transport: "http", State: mujica.MCPError, Error: "401 Unauthorized"},
			{ID: "db", Name: "postgres", Transport: "stdio", State: mujica.MCPStopped},
		}}
	})
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	r := mujica.MCPServerList(c, s.servers, mujica.MCPServerListOptions{Disabled: s.disabled})
	if r.Action == mujica.MCPNone {
		return
	}
	log(fmt.Sprintf("%s: %s", r.Server, [...]string{"", "start", "stop", "reconnect"}[r.Action]))
	for i := range s.servers {
		if s.servers[i].ID != r.Server {
			continue
		}
		switch r.Action {
		case mujica.MCPStop:
			s.servers[i].State = mujica.MCPStopped
		default:
			s.servers[i].State, s.servers[i].Tools, s.servers[i].Error = mujica.MCPConnected, 5, ""
		}
	}
}

type memoryDemoState struct {
	entries  []mujica.MemoryEntry
	disabled bool
}

func memoryPanelDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "278", func() memoryDemoState {
		return memoryDemoState{entries: []mujica.MemoryEntry{
			{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat, 3 Oct"},
			{ID: "2", Text: "The quartet needs two days' notice.", Source: "Plan, 1 Oct"},
			{ID: "3", Text: "Claret is served after nine.", Source: "Ledger"},
		}}
	})
	overlayControls(c, func() { overlayCheck(c, log, &s.disabled, "Disabled") })
	r := mujica.MemoryPanel(c, &s.entries, mujica.MemoryPanelOptions{Disabled: s.disabled})
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
	points := []mujica.Checkpoint{{ID: "c3", Label: "Seats assigned", Note: "12 files changed"},
		{ID: "c2", Label: "Menu drafted", Note: "3 files changed"}, {ID: "c1", Label: "Start"}}
	if mujica.CheckpointList(c, &s.current, points, mujica.CheckpointListOptions{Disabled: s.disabled}).Changed() {
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
		if mujica.Button(c, "Start a 30 s timer", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.deadline = time.Now().Add(30 * time.Second)
			log("timer started")
		}
	})
	r := mujica.HumanInputRequest(c, &s.answer, mujica.HumanInputRequestOptions{Question: "Which hall should host the masquerade?",
		Options: []string{"East wing", "Great hall"}, AllowFree: s.free, Deadline: s.deadline, Disabled: s.disabled})
	if r.Submitted {
		log("answer sent: " + s.answer)
	}
	if r.TimedOut {
		log("timed out")
	}
}
