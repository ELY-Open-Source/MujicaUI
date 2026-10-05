package agent

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestAgentProgressValueAndAction(t *testing.T) {
	value, action := 0.25, "Reading ledger.md"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AgentProgress(c, value, AgentProgressOptions{Action: action, Still: true})
	}, 400, 200)
	if !tt.HasText("25%") || !tt.HasText("Reading ledger.md") || !tt.HasText("Progress") {
		t.Fatalf("texts %q", tt.Texts())
	}
	value, action = 0.8, "Writing summary.md"
	tt.Frame()
	if !tt.HasText("80%") || !tt.HasText("Writing summary.md") || tt.HasText("Reading ledger.md") {
		t.Fatalf("update: %q", tt.Texts())
	}
}

func TestAgentProgressIndeterminate(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AgentProgress(c, 0, AgentProgressOptions{Label: "Research", Indeterminate: true, Still: true})
	}, 400, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if tt.HasText("0%") || !tt.HasText("Research") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestAgentProgressRejectsOutOfRange(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AgentProgress(c, 1.5, AgentProgressOptions{})
	}, 200, 100)
}

func TestAgentProgressSnapshot(t *testing.T) {
	snapshot(t, "agent_progress", 420, 110, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			AgentProgress(c, 0.62, AgentProgressOptions{Label: "Preparing the ball", Action: "Writing the seating chart"})
		})
	})
}
