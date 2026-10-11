package project

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var projSprint = Sprint{Name: "Sprint 12", Goal: "Ready the masquerade", Start: oct(1), End: oct(14)}

func sprintTester(sprint Sprint, tasks *[]Task, today core.Date, disabled bool, moves *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Height(640).Padding(8).Children(func() {
			r := SprintBoard(c, sprint, tasks, SprintBoardOptions{Today: today, Board: KanbanBoardOptions{Columns: kanbanColumns3,
				Statuses: projStatuses, Members: projMembers, Labels: projLabels, Disabled: disabled}})
			if _, ok := r.Moved(); ok {
				*moves++
			}
		})
	}, 1240, 660)
}

func TestSprintBoardHeaderAndMoves(t *testing.T) {
	tasks, n := projTasks(), 0
	tt := sprintTester(projSprint, &tasks, calendarToday, false, &n)
	for _, s := range []string{"Sprint 12", "Ready the masquerade", "October 1 – October 14, 2026", "10 days left", "1 of 5 done"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	_ = tt.Click("AVE-4 Rehearse act two")
	tt.Key(ui.Alt, ui.KeyRight)
	if n != 1 || !tt.HasText("2 of 5 done") {
		t.Fatal("moving to Done should count")
	}
}

func TestSprintBoardDaysAndDisabled(t *testing.T) {
	tasks, n := projTasks(), 0
	for _, c := range []struct {
		today core.Date
		want  string
	}{{core.Date{Year: 2026, Month: 9, Day: 30}, "Starts in 1 day"}, {oct(14), "1 day left"}, {oct(15), "Ended"}} {
		if tt := sprintTester(projSprint, &tasks, c.today, false, &n); !tt.HasText(c.want) {
			t.Errorf("%v: want %q", c.today, c.want)
		}
	}
	tt := sprintTester(projSprint, &tasks, calendarToday, true, &n)
	_ = tt.Click("AVE-4 Rehearse act two")
	tt.Key(ui.Alt, ui.KeyRight)
	if n != 0 {
		t.Fatal("disabled board moved")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an inverted sprint")
		}
	}()
	sprintTester(Sprint{Name: "x", Start: oct(5), End: oct(1)}, &tasks, calendarToday, false, &n)
}

func TestSprintBoardSnapshot(t *testing.T) {
	tasks := projTasks()
	snapshot(t, "sprint_board", 1240, 560, func(c *ui.Context) {
		ui.Column(c).Height(560).Padding(16).Children(func() {
			SprintBoard(c, projSprint, &tasks, SprintBoardOptions{Today: calendarToday, Board: KanbanBoardOptions{Columns: kanbanColumns3,
				Statuses: projStatuses, Members: projMembers, Labels: projLabels}})
		})
	})
}

func TestSprintBoardLargeData(t *testing.T) {
	tasks := make([]Task, 100000)
	for i := range tasks {
		tasks[i] = Task{ID: fmt.Sprintf("BIG-%d", i+1), Title: "Task", Status: "todo"}
	}
	moves := 0
	tt := sprintTester(projSprint, &tasks, calendarToday, false, &moves)
	if built := len(tt.Texts()); built > 600 {
		t.Fatalf("built %d labels", built)
	}
	if !tt.HasText("0 of 100000 done") {
		t.Fatal("incorrect sprint summary")
	}
}
