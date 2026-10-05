package main

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/datetime"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/project"
)

var projectsStatuses = []project.WorkflowStatus{
	{ID: "todo", Name: "To do", Group: project.StatusUnstarted},
	{ID: "doing", Name: "In progress", Group: project.StatusStarted},
	{ID: "done", Name: "Done", Group: project.StatusCompleted},
}
var projectsMembers = []project.ProjectMember{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}, {ID: "cy", Name: "Cy"}}
var projectsLabels = []project.TaskLabel{{ID: "ui", Name: "UI", Color: datetime.EventInfo}, {ID: "release", Name: "Release", Color: datetime.EventSuccess}}

func projectsCheck(c *ui.Context, log func(string), value *bool, label string) {
	if input.Checkbox(c, value, label, input.CheckboxOptions{}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *value))
	}
}

func projectsChoice[T comparable](c *ui.Context, log func(string), value *T, options []input.SelectOption[T], label string) {
	if input.Select(c, value, options, input.SelectOptions{Label: label}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *value))
	}
}

type projectsTaskState struct {
	tasks    []project.Task
	disabled bool
	version  int
	group    project.TaskGrouping
	filter   project.TaskFilter
	wip      float64
	closed   bool
}

func projectsTasks(c *ui.Context, key string) *projectsTaskState {
	return *ui.Local(c.Root(), key, func() *projectsTaskState {
		return &projectsTaskState{wip: 2, tasks: []project.Task{
			{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh, Due: core.Date{Year: 2026, Month: 10, Day: 3}, Assignees: []string{"ann"}, Labels: []string{"ui"}, Subtasks: []project.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}},
			{ID: "ENG-2", Title: "Check keyboard paths", Status: "doing", Priority: project.PriorityUrgent, Assignees: []string{"ben"}, Labels: []string{"release"}},
			{ID: "ENG-3", Title: "Publish docs", Status: "done", Done: true, Priority: project.PriorityLow},
		}}
	})
}

func projectsBoardOptions(s *projectsTaskState) project.KanbanBoardOptions {
	return project.KanbanBoardOptions{Columns: []project.KanbanColumn{{Status: "todo"}, {Status: "doing", WIP: int(s.wip)}, {Status: "done"}}, Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Version: s.version, Disabled: s.disabled}
}

func projectsLogBoard(s *projectsTaskState, r project.KanbanResult, log func(string)) {
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
	if input.NumberInput(c, &s.wip, input.NumberInputOptions{Min: 0, Max: 10, Step: 1, Label: "WIP limit"}).Changed() {
		log(fmt.Sprintf("WIP = %g", s.wip))
	}
	ui.ScrollHorizontal(c).Height(420).FillWidth().Children(func() {
		ui.Box(c).Width(880).Height(420).Children(func() { projectsLogBoard(s, project.KanbanBoard(c, &s.tasks, projectsBoardOptions(s)), log) })
	})
}

func projectsDemoTaskList(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "362")
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsChoice(c, log, &s.group, []input.SelectOption[project.TaskGrouping]{{Value: project.TaskGroupNone, Label: "Ungrouped"}, {Value: project.TaskGroupStatus, Label: "By status"}, {Value: project.TaskGroupPriority, Label: "By priority"}}, "Grouping")
	ui.Box(c).Height(360).FillWidth().Children(func() {
		before := s.filter
		r := project.TaskList(c, s.tasks, &s.filter, project.TaskListOptions{Statuses: projectsStatuses, Members: projectsMembers, GroupBy: s.group, Version: s.version, Disabled: s.disabled})
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
	projectsChoice(c, log, &s.tasks[0].Priority, []input.SelectOption[project.Priority]{{Value: project.PriorityLow, Label: "Low"}, {Value: project.PriorityHigh, Label: "High"}, {Value: project.PriorityUrgent, Label: "Urgent"}}, "Priority")
	r := project.TaskItem(c, &s.tasks[0], project.TaskItemOptions{Members: projectsMembers, Disabled: s.disabled})
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
		if input.Button(c, "Open details", input.ButtonOptions{}).Clicked() {
			s.closed = false
			log("open details")
		}
		return
	}
	ui.Scroll(c).Height(620).FillWidth().Children(func() {
		ui.Box(c).MaxWidth(500).FillWidth().Children(func() {
			r := project.TaskDetailPanel(c, &s.tasks[0], project.TaskDetailPanelOptions{Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Disabled: s.disabled, Activity: func() { ui.Text(c, "Ann created this task") }})
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
	if input.Button(c, "Clear subtasks", input.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.tasks[0].Subtasks = nil
		log("clear subtasks")
	}
	if project.SubtaskList(c, &s.tasks[0].Subtasks, project.SubtaskListOptions{Disabled: s.disabled}).Changed() {
		log("subtasks changed")
	}
}

func projectsDemoSprintBoard(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "368")
	projectsCheck(c, log, &s.disabled, "Disabled")
	sprint := project.Sprint{Name: "Release sprint", Goal: "Ship the first release", Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}}
	ui.ScrollHorizontal(c).Height(540).FillWidth().Children(func() {
		ui.Box(c).Width(880).Height(540).Children(func() {
			projectsLogBoard(s, project.SprintBoard(c, sprint, &s.tasks, project.SprintBoardOptions{Board: projectsBoardOptions(s)}), log)
		})
	})
}

func projectsDemoLabelManager(c *ui.Context, log func(string)) {
	s := *ui.Local(c.Root(), "374", func() *struct {
		labels   []project.TaskLabel
		selected []string
		disabled bool
	} {
		return &struct {
			labels   []project.TaskLabel
			selected []string
			disabled bool
		}{labels: slices.Clone(projectsLabels)}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	if project.LabelManager(c, &s.labels, &s.selected, project.LabelManagerOptions{Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("labels changed: %d, selected %v", len(s.labels), s.selected))
	}
}
