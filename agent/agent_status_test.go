package agent

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestAgentStatusWords(t *testing.T) {
	states := []AgentState{AgentPlanning, AgentRunning, AgentWaiting, AgentDone, AgentFailed}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		for _, s := range states {
			AgentStatus(c, s, AgentStatusOptions{Still: true})
		}
		AgentStatus(c, AgentRunning, AgentStatusOptions{Label: "Writing the report", Still: true})
	}, 400, 400)
	for _, w := range []string{"Planning", "Running", "Waiting", "Done", "Failed", "Writing the report"} {
		if !tt.HasText(w) {
			t.Fatalf("missing %q in %q", w, tt.Texts())
		}
	}
}

func TestAgentStatusAnimationStopsWhenStill(t *testing.T) {
	still := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AgentStatus(c, AgentRunning, AgentStatusOptions{Still: still})
	}, 300, 100)
	tt.Frame()
	moving := string(tt.Image().Pix)
	time.Sleep(80 * time.Millisecond)
	tt.Frame()
	if moving == string(tt.Image().Pix) {
		t.Fatal("running status did not animate")
	}
	still = true
	tt.Frame()
	a := string(tt.Image().Pix)
	time.Sleep(80 * time.Millisecond)
	tt.Frame()
	if a != string(tt.Image().Pix) {
		t.Fatal("still status kept moving")
	}
}

func TestAgentStatusRejectsUnknownState(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AgentStatus(c, AgentState(42), AgentStatusOptions{})
	}, 200, 100)
}

func TestAgentStatusSnapshot(t *testing.T) {
	snapshot(t, "agent_status", 480, 120, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(8).Wrap().Children(func() {
			for _, s := range []AgentState{AgentPlanning, AgentRunning, AgentWaiting, AgentDone, AgentFailed} {
				AgentStatus(c, s, AgentStatusOptions{})
			}
		})
	})
}
