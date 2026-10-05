package mujica_test

import (
	"github.com/egoist/mygo/ui"
	"log"
	mujica "mujicaui"
	"time"
)

func ExampleKanbanBoard() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	tasks := []mujica.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.KanbanBoard(c, &tasks, mujica.KanbanBoardOptions{Columns: []mujica.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses})
		if move, ok := r.Moved(); ok {
			log.Println("move", move.ID, move.To, move.Index)
		}
		if id, ok := r.Opened(); ok {
			log.Println("open", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleTaskList() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	tasks := []mujica.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	filter := mujica.TaskShowAll
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TaskList(c, tasks, &filter, mujica.TaskListOptions{GroupBy: mujica.TaskGroupStatus, Statuses: statuses})
		if id, ok := r.Changed(); ok {
			log.Println("edit", id)
		}
		if id, ok := r.Opened(); ok {
			log.Println("open", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleTaskItem() {
	task := mujica.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TaskItem(c, &task, mujica.TaskItemOptions{})
		if id, ok := r.Changed(); ok {
			log.Println("edit", id)
		}
		if id, ok := r.Opened(); ok {
			log.Println("open", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleTaskDetailPanel() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	task := mujica.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.TaskDetailPanel(c, &task, mujica.TaskDetailPanelOptions{Statuses: statuses, Activity: func() { ui.Text(c, "Release planned") }})
		if r.Changed() {
			log.Println("task edited", task.ID)
		}
		if r.Closed() {
			log.Println("close", task.ID)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleSubtaskList() {
	subtasks := []mujica.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SubtaskList(c, &subtasks, mujica.SubtaskListOptions{}).Changed() {
			log.Println("subtasks changed", len(subtasks))
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleSprintBoard() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	tasks := []mujica.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	sprint := mujica.Sprint{Name: "Release sprint", Goal: "Ship", Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.SprintBoard(c, sprint, &tasks, mujica.SprintBoardOptions{Board: mujica.KanbanBoardOptions{Columns: []mujica.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses}})
		if move, ok := r.Moved(); ok {
			log.Println("move", move.ID, move.To)
		}
		if id, ok := r.Opened(); ok {
			log.Println("open", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleBurndownChart() {
	remaining := []float64{40, 32, 25, 18}
	opts := mujica.BurndownChartOptions{Frame: mujica.ChartFrame{Label: "Release burndown"}, Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}, Total: 40}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.BurndownChart(c, remaining, opts)
		if r.Changed() {
			log.Println("point", r.Point)
		}
		if name, ok := r.Toggled(); ok {
			log.Println("series", name)
		}
		if r.Retried() {
			log.Println("retry")
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExamplePriorityIndicator() {
	priority := mujica.PriorityHigh
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.PriorityIndicator(c, priority, mujica.PriorityIndicatorOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleStatusSelect() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	value := "todo"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.StatusSelect(c, &value, statuses, mujica.StatusSelectOptions{}).Changed() {
			log.Println("status", value)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleAssigneePicker() {
	members := []mujica.Member{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}}
	ids := []string{"ann"}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.AssigneePicker(c, &ids, members, mujica.AssigneePickerOptions{Max: 2}).Changed() {
			log.Println("assignees", ids)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleDueDatePicker() {
	var due mujica.Date
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.DueDatePicker(c, &due, mujica.DueDatePickerOptions{}).Changed() {
			log.Println("due", due)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleLabelManager() {
	labels := []mujica.TaskLabel{{ID: "ui", Name: "UI", Color: mujica.EventInfo}}
	selected := []string{}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.LabelManager(c, &labels, &selected, mujica.LabelManagerOptions{}).Changed() {
			log.Println("labels changed", len(labels), selected)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleIssueCard() {
	statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
	task := mujica.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.IssueCard(c, task, mujica.IssueCardOptions{Statuses: statuses}).Clicked() {
			log.Println("open", task.ID)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleIssueIdBadge() {
	id := "ENG-123"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.IssueIdBadge(c, id, mujica.IssueIdBadgeOptions{Copyable: true}).Changed() {
			log.Println("copied", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleMilestoneProgress() {
	milestone := mujica.Milestone{Title: "Release", Due: mujica.Date{Year: 2026, Month: 10, Day: 14}, Done: 3, Total: 8}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.MilestoneProgress(c, milestone, mujica.MilestoneProgressOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleWorkloadView() {
	loads := []mujica.Workload{{Member: mujica.Member{ID: "ann", Name: "Ann"}, Assigned: 46, Capacity: 40}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.WorkloadView(c, loads, mujica.WorkloadViewOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleTimeTracker() {
	var state mujica.TimeTrackerState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.TimeTracker(c, &state, mujica.TimeTrackerOptions{}).Changed() {
			log.Println("tracked", state.Total(c.Now()))
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExamplePomodoroTimer() {
	var state mujica.PomodoroState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.PomodoroTimer(c, &state, mujica.PomodoroTimerOptions{Focus: 25 * time.Minute})
		if r.Changed() {
			log.Println("timer", state.Phase(), state.Running())
		}
		if phase, ok := r.Finished(); ok {
			log.Println("finished", phase)
		}
	}
	ui.Render(view, 1000, 700, 1)
}
