package main

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

var projectsStatuses = []mujica.WorkflowStatus{
	{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted},
	{ID: "doing", Name: "In progress", Group: mujica.StatusStarted},
	{ID: "done", Name: "Done", Group: mujica.StatusCompleted},
}
var projectsMembers = []mujica.Member{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}, {ID: "cy", Name: "Cy"}}
var projectsLabels = []mujica.TaskLabel{{ID: "ui", Name: "UI", Color: mujica.EventInfo}, {ID: "release", Name: "Release", Color: mujica.EventSuccess}}

func projectsCheck(c *ui.Context, log func(string), value *bool, label string) {
	if mujica.Checkbox(c, value, label, mujica.CheckboxOptions{}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *value))
	}
}

func projectsChoice[T comparable](c *ui.Context, log func(string), value *T, options []mujica.SelectOption[T], label string) {
	if mujica.Select(c, value, options, mujica.SelectOptions{Label: label}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *value))
	}
}

type projectsTaskState struct {
	tasks    []mujica.Task
	disabled bool
	version  int
	group    mujica.TaskGrouping
	filter   mujica.TaskFilter
	wip      float64
	closed   bool
}

func projectsTasks(c *ui.Context, key string) *projectsTaskState {
	return *ui.Local(c.Root(), key, func() *projectsTaskState {
		return &projectsTaskState{wip: 2, tasks: []mujica.Task{
			{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh, Due: mujica.Date{Year: 2026, Month: 10, Day: 3}, Assignees: []string{"ann"}, Labels: []string{"ui"}, Subtasks: []mujica.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}},
			{ID: "ENG-2", Title: "Check keyboard paths", Status: "doing", Priority: mujica.PriorityUrgent, Assignees: []string{"ben"}, Labels: []string{"release"}},
			{ID: "ENG-3", Title: "Publish docs", Status: "done", Done: true, Priority: mujica.PriorityLow},
		}}
	})
}

func projectsBoardOptions(s *projectsTaskState) mujica.KanbanBoardOptions {
	return mujica.KanbanBoardOptions{Columns: []mujica.KanbanColumn{{Status: "todo"}, {Status: "doing", WIP: int(s.wip)}, {Status: "done"}}, Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Version: s.version, Disabled: s.disabled}
}

func projectsLogBoard(s *projectsTaskState, r mujica.KanbanResult, log func(string)) {
	if move, ok := r.Moved(); ok {
		s.version++
		log(fmt.Sprintf("move %s: %s -> %s at %d", move.ID, move.From, move.To, move.Index))
	}
	if id, ok := r.Opened(); ok {
		log("open " + id)
	}
}

func projectsDemoKanbanBoard(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "361")
	projectsCheck(c, log, &s.disabled, "Disabled")
	if mujica.NumberInput(c, &s.wip, mujica.NumberInputOptions{Min: 0, Max: 10, Step: 1, Label: "WIP limit"}).Changed() {
		log(fmt.Sprintf("WIP = %g", s.wip))
	}
	ui.ScrollHorizontal(c).Height(420).FillWidth().Children(func() {
		ui.Box(c).Width(880).Height(420).Children(func() { projectsLogBoard(s, mujica.KanbanBoard(c, &s.tasks, projectsBoardOptions(s)), log) })
	})
}

func projectsDemoTaskList(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "362")
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsChoice(c, log, &s.group, []mujica.SelectOption[mujica.TaskGrouping]{{Value: mujica.TaskGroupNone, Label: "Ungrouped"}, {Value: mujica.TaskGroupStatus, Label: "By status"}, {Value: mujica.TaskGroupPriority, Label: "By priority"}}, "Grouping")
	ui.Box(c).Height(360).FillWidth().Children(func() {
		before := s.filter
		r := mujica.TaskList(c, s.tasks, &s.filter, mujica.TaskListOptions{Statuses: projectsStatuses, Members: projectsMembers, GroupBy: s.group, Version: s.version, Disabled: s.disabled})
		if s.filter != before {
			log(fmt.Sprintf("filter = %d", s.filter))
		}
		if id, ok := r.Changed(); ok {
			s.version++
			log("edit " + id)
		}
		if id, ok := r.Opened(); ok {
			log("open " + id)
		}
	})
}

func projectsDemoTaskItem(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "363")
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsChoice(c, log, &s.tasks[0].Priority, []mujica.SelectOption[mujica.Priority]{{Value: mujica.PriorityLow, Label: "Low"}, {Value: mujica.PriorityHigh, Label: "High"}, {Value: mujica.PriorityUrgent, Label: "Urgent"}}, "Priority")
	r := mujica.TaskItem(c, &s.tasks[0], mujica.TaskItemOptions{Members: projectsMembers, Disabled: s.disabled})
	if id, ok := r.Changed(); ok {
		log("edit " + id)
	}
	if id, ok := r.Opened(); ok {
		log("open " + id)
	}
}

func projectsDemoTaskDetailPanel(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "364")
	projectsCheck(c, log, &s.disabled, "Disabled")
	if s.closed {
		if mujica.Button(c, "Open details", mujica.ButtonOptions{}).Clicked() {
			s.closed = false
			log("open details")
		}
		return
	}
	ui.Scroll(c).Height(620).FillWidth().Children(func() {
		ui.Box(c).MaxWidth(500).FillWidth().Children(func() {
			r := mujica.TaskDetailPanel(c, &s.tasks[0], mujica.TaskDetailPanelOptions{Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Disabled: s.disabled, Activity: func() { ui.Text(c, "Ann created this task") }})
			if r.Changed() {
				log("task fields changed")
			}
			if r.Closed() {
				s.closed = true
				log("close details")
			}
		})
	})
}

func projectsDemoSubtaskList(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "365")
	projectsCheck(c, log, &s.disabled, "Disabled")
	if mujica.Button(c, "Clear subtasks", mujica.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.tasks[0].Subtasks = nil
		log("clear subtasks")
	}
	if mujica.SubtaskList(c, &s.tasks[0].Subtasks, mujica.SubtaskListOptions{Disabled: s.disabled}).Changed() {
		log("subtasks changed")
	}
}

func projectsDemoSprintBoard(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "368")
	projectsCheck(c, log, &s.disabled, "Disabled")
	sprint := mujica.Sprint{Name: "Release sprint", Goal: "Ship the first release", Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}}
	ui.ScrollHorizontal(c).Height(540).FillWidth().Children(func() {
		ui.Box(c).Width(880).Height(540).Children(func() {
			projectsLogBoard(s, mujica.SprintBoard(c, sprint, &s.tasks, mujica.SprintBoardOptions{Board: projectsBoardOptions(s)}), log)
		})
	})
}

func projectsDemoLabelManager(c *ui.Context, log func(string)) {
	s := *ui.Local(c.Root(), "374", func() *struct {
		labels   []mujica.TaskLabel
		selected []string
		disabled bool
	} {
		return &struct {
			labels   []mujica.TaskLabel
			selected []string
			disabled bool
		}{labels: slices.Clone(projectsLabels)}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	if mujica.LabelManager(c, &s.labels, &s.selected, mujica.LabelManagerOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("labels changed: %d, selected %v", len(s.labels), s.selected))
	}
}
