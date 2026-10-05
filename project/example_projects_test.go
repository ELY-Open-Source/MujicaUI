package project_test

import (
	"log"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/datetime"
	"github.com/ZacharyZhang-NY/MujicaUI/project"
)

func ExampleKanbanBoard() {
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.KanbanBoard(c, &tasks, project.KanbanBoardOptions{Columns: []project.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses})
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
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	filter := project.TaskShowAll
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.TaskList(c, tasks, &filter, project.TaskListOptions{GroupBy: project.TaskGroupStatus, Statuses: statuses})
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
	task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.TaskItem(c, &task, project.TaskItemOptions{})
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
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.TaskDetailPanel(c, &task, project.TaskDetailPanelOptions{Statuses: statuses, Activity: func() { ui.Text(c, "Release planned") }})
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
	subtasks := []project.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.SubtaskList(c, &subtasks, project.SubtaskListOptions{}).Changed() {
			log.Println("subtasks changed", len(subtasks))
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleSprintBoard() {
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
	sprint := project.Sprint{Name: "Release sprint", Goal: "Ship", Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.SprintBoard(c, sprint, &tasks, project.SprintBoardOptions{Board: project.KanbanBoardOptions{Columns: []project.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses}})
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
	opts := project.BurndownChartOptions{Frame: chart.ChartFrame{Label: "Release burndown"}, Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}, Total: 40}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.BurndownChart(c, remaining, opts)
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
	priority := project.PriorityHigh
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		project.PriorityIndicator(c, priority, project.PriorityIndicatorOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleStatusSelect() {
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	value := "todo"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.StatusSelect(c, &value, statuses, project.StatusSelectOptions{}).Changed() {
			log.Println("status", value)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleAssigneePicker() {
	members := []project.ProjectMember{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}}
	ids := []string{"ann"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.AssigneePicker(c, &ids, members, project.AssigneePickerOptions{Max: 2}).Changed() {
			log.Println("assignees", ids)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleDueDatePicker() {
	var due core.Date
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.DueDatePicker(c, &due, project.DueDatePickerOptions{}).Changed() {
			log.Println("due", due)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleLabelManager() {
	labels := []project.TaskLabel{{ID: "ui", Name: "UI", Color: datetime.EventInfo}}
	selected := []string{}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.LabelManager(c, &labels, &selected, project.LabelManagerOptions{}).Changed() {
			log.Println("labels changed", len(labels), selected)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleIssueCard() {
	statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
	task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.IssueCard(c, task, project.IssueCardOptions{Statuses: statuses}).Clicked() {
			log.Println("open", task.ID)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleIssueIdBadge() {
	id := "ENG-123"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.IssueIdBadge(c, id, project.IssueIdBadgeOptions{Copyable: true}).Changed() {
			log.Println("copied", id)
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleMilestoneProgress() {
	milestone := project.Milestone{Title: "Release", Due: core.Date{Year: 2026, Month: 10, Day: 14}, Done: 3, Total: 8}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		project.MilestoneProgress(c, milestone, project.MilestoneProgressOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleWorkloadView() {
	loads := []project.Workload{{Member: project.ProjectMember{ID: "ann", Name: "Ann"}, Assigned: 46, Capacity: 40}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		project.WorkloadView(c, loads, project.WorkloadViewOptions{})
	}
	ui.Render(view, 1000, 700, 1)
}

func ExampleTimeTracker() {
	var state project.TimeTrackerState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if project.TimeTracker(c, &state, project.TimeTrackerOptions{}).Changed() {
			log.Println("tracked", state.Total(c.Now()))
		}
	}
	ui.Render(view, 1000, 700, 1)
}

func ExamplePomodoroTimer() {
	var state project.PomodoroState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := project.PomodoroTimer(c, &state, project.PomodoroTimerOptions{Focus: 25 * time.Minute})
		if r.Changed() {
			log.Println("timer", state.Phase(), state.Running())
		}
		if phase, ok := r.Finished(); ok {
			log.Println("finished", phase)
		}
	}
	ui.Render(view, 1000, 700, 1)
}
