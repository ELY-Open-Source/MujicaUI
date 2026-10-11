package agent

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var agentStepsSample = []AgentStep{
	{ID: "plan", Title: "Read the guest list", State: AgentDone, Duration: 800 * time.Millisecond, Detail: "42 guests, 3 late."},
	{ID: "seat", Title: "Assign the seats", State: AgentFailed, Duration: 2 * time.Second, Detail: "Table 4 is over capacity."},
	{ID: "send", Title: "Send the invitations", State: AgentRunning},
	{ID: "log", Title: "Record the ledger", State: AgentPending},
}

func stepTester(steps []AgentStep, opts AgentStepListOptions, got *[]AgentStepListResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := AgentStepList(c, steps, opts)
		if r.Retry != "" || r.Toggled != "" {
			*got = append(*got, r)
		}
	}, 600, 500)
}

func TestAgentStepListCurrentDetailRetry(t *testing.T) {
	var got []AgentStepListResult
	tt := stepTester(agentStepsSample, AgentStepListOptions{Still: true}, &got)
	if !tt.HasText("Current") || !tt.HasText("Running") || !tt.HasText("800 ms") {
		t.Fatalf("texts %q", tt.Texts())
	}
	if tt.HasText("Table 4 is over capacity.") {
		t.Fatal("detail shown while closed")
	}
	_ = tt.Click("Assign the seats")
	if len(got) != 1 || got[0].Toggled != "seat" || !tt.HasText("Table 4 is over capacity.") {
		t.Fatalf("toggle: %+v", got)
	}
	_ = tt.Click("Retry")
	if len(got) != 2 || got[1].Retry != "seat" {
		t.Fatalf("retry: %+v", got)
	}
}

func TestAgentStepListKeyboard(t *testing.T) {
	var got []AgentStepListResult
	tt := stepTester(agentStepsSample, AgentStepListOptions{Still: true}, &got)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 1 || got[0].Toggled != "plan" || !tt.HasText("42 guests, 3 late.") {
		t.Fatalf("Enter on the first step: %+v", got)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown) // past Assign the seats onto its Retry
	tt.Key(0, ui.KeyEnter)
	if len(got) != 2 || got[1].Retry != "seat" {
		t.Fatalf("arrows to Retry: %+v", got)
	}
}

func TestAgentStepListDisabled(t *testing.T) {
	var got []AgentStepListResult
	tt := stepTester(agentStepsSample, AgentStepListOptions{Disabled: true, Still: true}, &got)
	_ = tt.Click("Assign the seats")
	_ = tt.Click("Retry")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 0 {
		t.Fatalf("disabled list reported %+v", got)
	}
}

func TestAgentStepListRejectsDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	stepTester([]AgentStep{{ID: "a", Title: "x"}, {ID: "a", Title: "y"}}, AgentStepListOptions{}, nil)
}

func TestAgentStepListSnapshot(t *testing.T) {
	snapshot(t, "agent_step_list", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { AgentStepList(c, agentStepsSample, AgentStepListOptions{Still: true}) })
	})
}
