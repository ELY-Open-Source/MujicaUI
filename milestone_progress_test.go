package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func milestoneTester(m Milestone) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Box(c).Width(360).Children(func() { MilestoneProgress(c, m, MilestoneProgressOptions{Today: calendarToday}) })
	}, 400, 200)
}

func TestMilestoneProgressStates(t *testing.T) {
	for _, c := range []struct {
		m    Milestone
		want []string
	}{
		{Milestone{Title: "Act one", Due: oct(12), Done: 3, Total: 8}, []string{"37%", "3 of 8 tasks done", "Due October 12"}},
		{Milestone{Title: "Act two", Due: oct(1), Done: 5, Total: 8}, []string{"62%", "Overdue by 4 days"}},
		{Milestone{Title: "Premiere", Due: oct(1), Done: 4, Total: 4}, []string{"100%", "Reached"}},
		{Milestone{Title: "Encore", Total: 0}, []string{"0%", "0 of 0 tasks done"}},
	} {
		tt := milestoneTester(c.m)
		for _, s := range c.want {
			if !tt.HasText(s) {
				t.Errorf("%s: missing %q in %v", c.m.Title, s, tt.Texts())
			}
		}
	}
	if tt := milestoneTester(Milestone{Title: "Premiere", Due: oct(1), Done: 4, Total: 4}); tt.HasText("Overdue") {
		t.Error("a reached milestone is not overdue")
	}
}

func TestMilestoneProgressInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	milestoneTester(Milestone{Title: "x", Done: 5, Total: 3})
}

func TestMilestoneProgressSnapshot(t *testing.T) {
	snapshot(t, "milestone_progress", 760, 140, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			for _, m := range []Milestone{{Title: "Act one", Due: oct(12), Done: 3, Total: 8}, {Title: "Act two", Due: oct(1), Done: 5, Total: 8}} {
				ui.Box(c).Width(350).Children(func() { MilestoneProgress(c, m, MilestoneProgressOptions{Today: calendarToday}) })
			}
		})
	})
}
