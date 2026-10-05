package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func agentPlanSample() []PlanItem {
	return []PlanItem{
		{ID: "a", Text: "Read the guest list", State: AgentDone},
		{ID: "b", Text: "Assign the seats", State: AgentRunning},
		{ID: "c", Text: "Order the flowers"},
	}
}

func planTester(items *[]PlanItem, opts AgentPlanOptions, changes *int) *ui.Tester {
	opts.Still = true
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if AgentPlan(c, items, opts).Changed() {
			*changes++
		}
	}, 560, 500)
}

func TestAgentPlanProgressCheckSkipAdd(t *testing.T) {
	items := agentPlanSample()
	changes := 0
	tt := planTester(&items, AgentPlanOptions{}, &changes)
	if !tt.HasText("1 / 3 done") || !tt.HasText("In progress") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Order the flowers") // the box is named by the item
	if items[2].State != AgentDone || changes != 1 || !tt.HasText("2 / 3 done") {
		t.Fatalf("check: %+v changes %d", items[2], changes)
	}
	skips, _ := tt.Find("Skip")
	tt.ClickAt(skips.X+4, skips.Y+4) // the running item's Skip
	if items[1].State != AgentSkipped || !tt.HasText("2 / 2 done") {
		t.Fatalf("skip: %+v texts %q", items[1], tt.Texts())
	}
	_ = tt.Click("Restore")
	if items[1].State != AgentPending {
		t.Fatalf("restore: %+v", items[1])
	}
	_ = tt.Click("Add a step")
	tt.Type("  ")
	tt.Key(0, ui.KeyEnter)
	if len(items) != 3 {
		t.Fatal("blank item added")
	}
	tt.Type("Light the candles")
	tt.Key(0, ui.KeyEnter)
	if len(items) != 4 || items[3].Text != "Light the candles" || items[3].ID == "" || !tt.HasText("Light the candles") {
		t.Fatalf("add: %+v", items)
	}
}

func TestAgentPlanLiveUpdateAndDisabled(t *testing.T) {
	items := agentPlanSample()
	changes := 0
	tt := planTester(&items, AgentPlanOptions{Disabled: true, Title: "Ball"}, &changes)
	_ = tt.Click("Order the flowers")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if changes != 0 || items[2].State != AgentPending || tt.HasText("Add a step") || tt.HasText("Skip") {
		t.Fatalf("disabled plan changed: %+v", items)
	}
	items[2].State = AgentDone // the agent checks it off
	tt.Frame()
	if !tt.HasText("2 / 3 done") || !tt.HasText("Ball") {
		t.Fatalf("live update: %q", tt.Texts())
	}
}

func TestAgentPlanRejectsBadState(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	items := []PlanItem{{ID: "a", Text: "x", State: AgentFailed}}
	n := 0
	planTester(&items, AgentPlanOptions{}, &n)
}

func TestAgentPlanSnapshot(t *testing.T) {
	snapshot(t, "agent_plan", 520, 330, func(c *ui.Context) {
		items := agentPlanSample()
		items = append(items, PlanItem{ID: "d", Text: "Hire the quartet", State: AgentSkipped})
		ui.Column(c).Padding(16).Children(func() { AgentPlan(c, &items, AgentPlanOptions{Still: true}) })
	})
}
