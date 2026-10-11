package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/project"
)

func projectsDemoMilestoneProgress(c *ui.Context, log func(string)) {
	done := ui.Local(c.Root(), "377", func() float64 { return 3 })
	if input.NumberInput(c, done, input.NumberInputOptions{Min: 0, Max: 8, Step: 1, Label: "Completed tasks"}).Changed() {
		log(fmt.Sprintf("completed = %g", *done))
	}
	project.MilestoneProgress(c, project.Milestone{Title: "Release", Due: core.Date{Year: 2026, Month: 10, Day: 14}, Done: int(*done), Total: 8}, project.MilestoneProgressOptions{})
}

func projectsDemoWorkloadView(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "378", func() struct {
		assigned float64
		loads    []project.Workload
		version  int
	} {
		return struct {
			assigned float64
			loads    []project.Workload
			version  int
		}{assigned: 46, loads: []project.Workload{{Member: projectsMembers[0], Assigned: 32, Capacity: 40}, {Member: projectsMembers[1], Assigned: 46, Capacity: 40}, {Member: projectsMembers[2], Assigned: 12, Capacity: 20}}}
	})
	if input.NumberInput(c, &s.assigned, input.NumberInputOptions{Min: 0, Max: 80, Step: 1, Label: "Ben assigned hours"}).Changed() {
		s.loads[1].Assigned = s.assigned
		s.version++
		log(fmt.Sprintf("assigned = %g", s.assigned))
	}
	ui.Box(c).Height(240).FillWidth().Children(func() { project.WorkloadView(c, s.loads, project.WorkloadViewOptions{Version: s.version}) })
}

func projectsDemoTimeTracker(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "379", func() struct {
		timer    project.TimeTrackerState
		disabled bool
	} {
		return struct {
			timer    project.TimeTrackerState
			disabled bool
		}{}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	ui.ScrollHorizontal(c).Height(120).FillWidth().Children(func() {
		if project.TimeTracker(c, &s.timer, project.TimeTrackerOptions{Disabled: s.disabled}).Changed() {
			log("tracked " + s.timer.Total(c.Now()).String())
		}
	})
}

func projectsDemoPomodoroTimer(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "380", func() struct {
		timer            project.PomodoroState
		disabled, hidden bool
		seconds          float64
	} {
		return struct {
			timer            project.PomodoroState
			disabled, hidden bool
			seconds          float64
		}{seconds: 10}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, &s.hidden, "Hidden")
	if input.NumberInput(c, &s.seconds, input.NumberInputOptions{Min: 1, Max: 1500, Step: 1, Label: "Focus seconds", Disabled: s.timer.Running()}).Changed() {
		log(fmt.Sprintf("focus = %g seconds", s.seconds))
	}
	if s.hidden {
		return
	}
	r := project.PomodoroTimer(c, &s.timer, project.PomodoroTimerOptions{Focus: time.Duration(s.seconds) * time.Second, Break: 5 * time.Second, LongBreak: 10 * time.Second, LongEvery: 2, Disabled: s.disabled})
	if r.Changed() {
		log(fmt.Sprintf("timer phase = %d, running = %v", s.timer.Phase(), s.timer.Running()))
	}
	if phase, ok := r.Finished(); ok {
		log(fmt.Sprintf("finished phase %d", phase))
	}
}
