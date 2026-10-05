package main

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/project"
)

func projectsDemoBurndownChart(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "369", func() struct {
		disabled bool
		status   data.DataStatus
		total    float64
		values   []float64
	} {
		return struct {
			disabled bool
			status   data.DataStatus
			total    float64
			values   []float64
		}{total: 40, values: []float64{40, 38, 32, 34, 26, math.NaN(), 20, 18}}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsChoice(c, log, &s.status, []input.SelectOption[data.DataStatus]{{Value: data.DataReady, Label: "Ready"}, {Value: data.DataLoading, Label: "Loading"}, {Value: data.DataFailed, Label: "Error"}}, "Data state")
	if input.NumberInput(c, &s.total, input.NumberInputOptions{Min: 1, Max: 100, Step: 1, Label: "Total work"}).Changed() {
		log(fmt.Sprintf("total = %g", s.total))
	}
	ui.Box(c).Disabled(s.disabled).FillWidth().Children(func() {
		r := project.BurndownChart(c, s.values, project.BurndownChartOptions{Frame: chart.ChartFrame{Label: "Release burndown", Status: s.status, Error: "Data unavailable"}, Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}, Total: s.total})
		if r.Changed() {
			log(fmt.Sprintf("point = %d", r.Point))
		}
		if name, ok := r.Toggled(); ok {
			log("toggle " + name)
		}
		if r.Retried() {
			s.status = data.DataReady
			log("retry")
		}
	})
}

func projectsDemoPriorityIndicator(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "370", func() struct {
		priority project.Priority
		iconOnly bool
	} {
		return struct {
			priority project.Priority
			iconOnly bool
		}{priority: project.PriorityHigh}
	})
	projectsCheck(c, log, &s.iconOnly, "Icon only")
	projectsChoice(c, log, &s.priority, []input.SelectOption[project.Priority]{{Value: project.PriorityNone, Label: "None"}, {Value: project.PriorityLow, Label: "Low"}, {Value: project.PriorityMedium, Label: "Medium"}, {Value: project.PriorityHigh, Label: "High"}, {Value: project.PriorityUrgent, Label: "Urgent"}}, "Priority")
	project.PriorityIndicator(c, s.priority, project.PriorityIndicatorOptions{IconOnly: s.iconOnly})
}

func projectsDemoStatusSelect(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "371", func() struct {
		value    string
		disabled bool
	} {
		return struct {
			value    string
			disabled bool
		}{value: "todo"}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	if project.StatusSelect(c, &s.value, projectsStatuses, project.StatusSelectOptions{Disabled: s.disabled}).Changed() {
		log("status " + s.value)
	}
}

func projectsDemoAssigneePicker(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "372", func() struct {
		ids      []string
		disabled bool
		max      float64
	} {
		return struct {
			ids      []string
			disabled bool
			max      float64
		}{ids: []string{"ann"}, max: 2}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	if input.NumberInput(c, &s.max, input.NumberInputOptions{Min: 0, Max: 3, Step: 1, Label: "Selection limit"}).Changed() {
		log(fmt.Sprintf("limit = %g", s.max))
	}
	if project.AssigneePicker(c, &s.ids, projectsMembers, project.AssigneePickerOptions{Max: int(s.max), Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("assignees %v", s.ids))
	}
}

func projectsDemoDueDatePicker(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "373", func() struct {
		due            core.Date
		disabled, done bool
	} {
		return struct {
			due            core.Date
			disabled, done bool
		}{due: core.Date{Year: 2026, Month: 10, Day: 3}}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, &s.done, "Completed")
	if project.DueDatePicker(c, &s.due, project.DueDatePickerOptions{Done: s.done, Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("due %v", s.due))
	}
}

func projectsDemoIssueCard(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "375")
	selected := ui.Local(c.Root(), "375-selected", func() bool { return false })
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, selected, "Selected")
	ui.Box(c).Width(300).Children(func() {
		if project.IssueCard(c, s.tasks[0], project.IssueCardOptions{Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Selected: *selected, Disabled: s.disabled}).Clicked() {
			*selected = !*selected
			log("open " + s.tasks[0].ID)
		}
	})
}

func projectsDemoIssueIdBadge(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "376", func() struct {
		disabled, copyable bool
		id                 string
	} {
		return struct {
			disabled, copyable bool
			id                 string
		}{copyable: true, id: "ENG-123"}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, &s.copyable, "Copyable")
	projectsChoice(c, log, &s.id, []input.SelectOption[string]{{Value: "ENG-123", Label: "ENG-123"}, {Value: "OPS-7", Label: "OPS-7"}}, "Issue key")
	if project.IssueIdBadge(c, s.id, project.IssueIdBadgeOptions{Copyable: s.copyable, Disabled: s.disabled}).Changed() {
		log("copied " + s.id)
	}
}
