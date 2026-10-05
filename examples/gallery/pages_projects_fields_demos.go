package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"math"
	mujica "mujicaui"
)

func projectsDemoBurndownChart(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "369", func() struct {
		disabled bool
		status   mujica.DataStatus
		total    float64
		values   []float64
	} {
		return struct {
			disabled bool
			status   mujica.DataStatus
			total    float64
			values   []float64
		}{total: 40, values: []float64{40, 38, 32, 34, 26, math.NaN(), 20, 18}}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsChoice(c, log, &s.status, []mujica.SelectOption[mujica.DataStatus]{{Value: mujica.DataReady, Label: "Ready"}, {Value: mujica.DataLoading, Label: "Loading"}, {Value: mujica.DataFailed, Label: "Error"}}, "Data state")
	if mujica.NumberInput(c, &s.total, mujica.NumberInputOptions{Min: 1, Max: 100, Step: 1, Label: "Total work"}).Changed() {
		log(fmt.Sprintf("total = %g", s.total))
	}
	ui.Box(c).Disabled(s.disabled).FillWidth().Children(func() {
		r := mujica.BurndownChart(c, s.values, mujica.BurndownChartOptions{Frame: mujica.ChartFrame{Label: "Release burndown", Status: s.status, Error: "Data unavailable"}, Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}, Total: s.total})
		if r.Changed() {
			log(fmt.Sprintf("point = %d", r.Point))
		}
		if name, ok := r.Toggled(); ok {
			log("toggle " + name)
		}
		if r.Retried() {
			s.status = mujica.DataReady
			log("retry")
		}
	})
}

func projectsDemoPriorityIndicator(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "370", func() struct {
		priority mujica.Priority
		iconOnly bool
	} {
		return struct {
			priority mujica.Priority
			iconOnly bool
		}{priority: mujica.PriorityHigh}
	})
	projectsCheck(c, log, &s.iconOnly, "Icon only")
	projectsChoice(c, log, &s.priority, []mujica.SelectOption[mujica.Priority]{{Value: mujica.PriorityNone, Label: "None"}, {Value: mujica.PriorityLow, Label: "Low"}, {Value: mujica.PriorityMedium, Label: "Medium"}, {Value: mujica.PriorityHigh, Label: "High"}, {Value: mujica.PriorityUrgent, Label: "Urgent"}}, "Priority")
	mujica.PriorityIndicator(c, s.priority, mujica.PriorityIndicatorOptions{IconOnly: s.iconOnly})
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
	if mujica.StatusSelect(c, &s.value, projectsStatuses, mujica.StatusSelectOptions{Disabled: s.disabled}).Changed() {
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
	if mujica.NumberInput(c, &s.max, mujica.NumberInputOptions{Min: 0, Max: 3, Step: 1, Label: "Selection limit"}).Changed() {
		log(fmt.Sprintf("limit = %g", s.max))
	}
	if mujica.AssigneePicker(c, &s.ids, projectsMembers, mujica.AssigneePickerOptions{Max: int(s.max), Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("assignees %v", s.ids))
	}
}

func projectsDemoDueDatePicker(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "373", func() struct {
		due            mujica.Date
		disabled, done bool
	} {
		return struct {
			due            mujica.Date
			disabled, done bool
		}{due: mujica.Date{Year: 2026, Month: 10, Day: 3}}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, &s.done, "Completed")
	if mujica.DueDatePicker(c, &s.due, mujica.DueDatePickerOptions{Done: s.done, Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("due %v", s.due))
	}
}

func projectsDemoIssueCard(c *ui.Context, log func(string)) {
	s := projectsTasks(c, "375")
	selected := ui.Local(c.Root(), "375-selected", func() bool { return false })
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, selected, "Selected")
	ui.Box(c).Width(300).Children(func() {
		if mujica.IssueCard(c, s.tasks[0], mujica.IssueCardOptions{Statuses: projectsStatuses, Members: projectsMembers, Labels: projectsLabels, Selected: *selected, Disabled: s.disabled}).Clicked() {
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
	projectsChoice(c, log, &s.id, []mujica.SelectOption[string]{{Value: "ENG-123", Label: "ENG-123"}, {Value: "OPS-7", Label: "OPS-7"}}, "Issue key")
	if mujica.IssueIdBadge(c, s.id, mujica.IssueIdBadgeOptions{Copyable: s.copyable, Disabled: s.disabled}).Changed() {
		log("copied " + s.id)
	}
}
