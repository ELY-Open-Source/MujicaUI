package agent

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var checkpointsSample = []Checkpoint{
	{ID: "c3", Label: "Seats assigned", Note: "12 files changed"},
	{ID: "c2", Label: "Menu drafted", Note: "3 files changed"},
	{ID: "c1", Label: "Start"},
}

func checkpointTester(cur *string, points []Checkpoint, opts CheckpointListOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if CheckpointList(c, cur, points, opts).Changed() {
			*changes++
		}
	}, 560, 420)
}

func TestCheckpointListRollbackConfirm(t *testing.T) {
	cur, changes := "c3", 0
	tt := checkpointTester(&cur, checkpointsSample, CheckpointListOptions{}, &changes)
	if !tt.HasText("Current") || !tt.HasText("12 files changed") {
		t.Fatalf("texts %q", tt.Texts())
	}
	rolls, _ := tt.Find("Roll back")
	tt.ClickAt(rolls.X+4, rolls.Y+4) // Menu drafted
	if !tt.HasText("Roll back to “Menu drafted”?") || !tt.Focused("Cancel") {
		t.Fatalf("confirm: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter) // Enter on Cancel cancels
	if cur != "c3" || changes != 0 {
		t.Fatalf("Enter rolled back: %s", cur)
	}
	tt.ClickAt(rolls.X+4, rolls.Y+4)
	_ = tt.Click("Roll back to “Menu drafted”?")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if cur != "c2" || changes != 1 {
		t.Fatalf("rollback: cur %s changes %d", cur, changes)
	}
}

func TestCheckpointListDisabledEmptyAndMany(t *testing.T) {
	cur, n := "c3", 0
	tt := checkpointTester(&cur, checkpointsSample, CheckpointListOptions{Disabled: true}, &n)
	_ = tt.Click("Roll back")
	if tt.HasText("Roll back to “Menu drafted”?") || cur != "c3" {
		t.Fatal("disabled list asked")
	}
	tt = checkpointTester(&cur, nil, CheckpointListOptions{}, &n)
	if !tt.HasText("No checkpoints yet") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	many := make([]Checkpoint, 100000)
	for i := range many {
		many[i] = Checkpoint{ID: fmt.Sprint(i), Label: fmt.Sprintf("point %d", i)}
	}
	tt = checkpointTester(&cur, many, CheckpointListOptions{}, &n)
	if got := len(tt.Texts()); got > 300 {
		t.Fatalf("%d texts built for 100000 checkpoints", got)
	}
}

func TestCheckpointListSnapshot(t *testing.T) {
	snapshot(t, "checkpoint_list", 520, 240, func(c *ui.Context) {
		cur := "c2"
		ui.Column(c).Padding(16).Children(func() {
			CheckpointList(c, &cur, checkpointsSample, CheckpointListOptions{})
		})
	})
}
