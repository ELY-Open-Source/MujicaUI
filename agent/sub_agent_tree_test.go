package agent

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var subAgentSample = []SubAgent{
	{ID: "lead", Name: "Steward", Task: "Prepare the ball", State: AgentRunning, Children: []SubAgent{
		{ID: "music", Name: "Musician", Task: "Hire the quartet", State: AgentDone},
		{ID: "kitchen", Name: "Cook", Task: "Plan the menu", State: AgentWaiting, Children: []SubAgent{
			{ID: "wine", Name: "Sommelier", Task: "Choose the wine", State: AgentFailed},
		}},
	}},
}

func subAgentTester(opts SubAgentTreeOptions, changes *int) *ui.Tester {
	opts.Still = true
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SubAgentTree(c, subAgentSample, opts).Changed() {
			*changes++
		}
	}, 560, 300)
}

func TestSubAgentTreeExpandByPointerAndKeys(t *testing.T) {
	changes := 0
	tt := subAgentTester(SubAgentTreeOptions{}, &changes)
	if tt.HasText("Musician") || !tt.HasText("Running") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Steward")
	if changes != 1 || !tt.HasText("Musician") || !tt.HasText("Waiting") || tt.HasText("Sommelier") {
		t.Fatalf("click: changes %d texts %q", changes, tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown) // Cook
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Sommelier") || !tt.HasText("Failed") {
		t.Fatal("Right did not expand Cook")
	}
	tt.Key(0, ui.KeyRight) // already open: nothing
	tt.Key(0, ui.KeyLeft)
	if tt.HasText("Sommelier") || changes != 3 {
		t.Fatalf("Left: changes %d", changes)
	}
}

func TestSubAgentTreeDisabled(t *testing.T) {
	changes := 0
	tt := subAgentTester(SubAgentTreeOptions{Disabled: true}, &changes)
	_ = tt.Click("Steward")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if changes != 0 || tt.HasText("Musician") {
		t.Fatal("disabled tree expanded")
	}
}

func TestSubAgentTreeRejectsDuplicateIDs(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SubAgentTree(c, []SubAgent{{ID: "a", Name: "x", Children: []SubAgent{{ID: "a", Name: "y"}}}}, SubAgentTreeOptions{})
	}, 200, 200)
}

func TestSubAgentTreeSnapshot(t *testing.T) {
	snapshotAfter(t, "sub_agent_tree", 540, 200, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() { SubAgentTree(c, subAgentSample, SubAgentTreeOptions{Still: true}) })
	}, func(tt *ui.Tester) {
		_ = tt.Click("Steward")
		_ = tt.Click("Cook")
	})
}
