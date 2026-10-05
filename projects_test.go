package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var projStatuses = []WorkflowStatus{
	{ID: "backlog", Name: "Backlog", Group: StatusBacklog},
	{ID: "todo", Name: "Todo", Group: StatusUnstarted},
	{ID: "doing", Name: "In Progress", Group: StatusStarted},
	{ID: "review", Name: "In Review", Group: StatusStarted},
	{ID: "done", Name: "Done", Group: StatusCompleted},
	{ID: "dropped", Name: "Dropped", Group: StatusCanceled},
}

var projMembers = []ProjectMember{{ID: "saki", Name: "Sakiko Togawa"}, {ID: "mutsu", Name: "Mutsumi Wakaba"},
	{ID: "uika", Name: "Uika Misumi"}, {ID: "nyamu", Name: "Nyamu Yuutenji"}, {ID: "umiri", Name: "Umiri Yahata"}}

var projLabels = []TaskLabel{{ID: "bug", Name: "Bug", Color: EventDanger}, {ID: "ui", Name: "UI", Color: EventInfo},
	{ID: "docs", Name: "Docs", Color: EventSuccess}}

// projTasks returns fresh tasks for one test.
func projTasks() []Task {
	return []Task{
		{ID: "AVE-1", Title: "Write the overture", Status: "todo", Priority: PriorityHigh, Due: oct(7), Assignees: []string{"saki"}, Labels: []string{"docs"}},
		{ID: "AVE-2", Title: "Tune the organ", Status: "doing", Priority: PriorityUrgent, Due: oct(3), Assignees: []string{"mutsu", "uika"}, Labels: []string{"bug"}},
		{ID: "AVE-3", Title: "Masks for the ball", Status: "doing", Priority: PriorityMedium, Assignees: []string{"nyamu"}},
		{ID: "AVE-4", Title: "Rehearse act two", Status: "review", Priority: PriorityLow, Due: oct(12)},
		{ID: "AVE-5", Title: "Book the hall", Status: "done", Done: true, Priority: PriorityNone, Due: oct(1), Labels: []string{"ui"}},
	}
}

func TestProjectDueText(t *testing.T) {
	var got []string
	var late []bool
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		for _, c2 := range []struct {
			due  Date
			done bool
		}{{oct(2), false}, {oct(4), false}, {oct(5), false}, {oct(6), false}, {oct(14), false}, {oct(1), true}} {
			s, l := projectDueText(c, c2.due, calendarToday, c2.done)
			got, late = append(got, s), append(late, l)
		}
	}, 10, 10)
	want := []string{"Overdue by 3 days", "Overdue by 1 day", "Due today", "Due tomorrow", "Due October 14", "Due October 1"}
	for i := range want {
		if got[i] != want[i] || late[i] != (i < 2) {
			t.Fatalf("%d: %q %v", i, got[i], late[i])
		}
	}
}
