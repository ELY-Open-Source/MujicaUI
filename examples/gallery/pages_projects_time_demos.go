package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
	"time"
)

func projectsDemoMilestoneProgress(c *ui.Context, log func(string)) {
	done := ui.Local(c.Root(), "377", func() float64 { return 3 })
	if mujica.NumberInput(c, done, mujica.NumberInputOptions{Min: 0, Max: 8, Step: 1, Label: "Completed tasks"}).Changed() {
		log(fmt.Sprintf("completed = %g", *done))
	}
	mujica.MilestoneProgress(c, mujica.Milestone{Title: "Release", Due: mujica.Date{Year: 2026, Month: 10, Day: 14}, Done: int(*done), Total: 8}, mujica.MilestoneProgressOptions{})
}

func projectsDemoWorkloadView(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "378", func() struct {
		assigned float64
		loads    []mujica.Workload
		version  int
	} {
		return struct {
			assigned float64
			loads    []mujica.Workload
			version  int
		}{assigned: 46, loads: []mujica.Workload{{Member: projectsMembers[0], Assigned: 32, Capacity: 40}, {Member: projectsMembers[1], Assigned: 46, Capacity: 40}, {Member: projectsMembers[2], Assigned: 12, Capacity: 20}}}
	})
	if mujica.NumberInput(c, &s.assigned, mujica.NumberInputOptions{Min: 0, Max: 80, Step: 1, Label: "Ben assigned hours"}).Changed() {
		s.loads[1].Assigned = s.assigned
		s.version++
		log(fmt.Sprintf("assigned = %g", s.assigned))
	}
	ui.Box(c).Height(240).FillWidth().Children(func() { mujica.WorkloadView(c, s.loads, mujica.WorkloadViewOptions{Version: s.version}) })
}

func projectsDemoTimeTracker(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "379", func() struct {
		timer    mujica.TimeTrackerState
		disabled bool
	} {
		return struct {
			timer    mujica.TimeTrackerState
			disabled bool
		}{}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	ui.ScrollHorizontal(c).Height(120).FillWidth().Children(func() {
		if mujica.TimeTracker(c, &s.timer, mujica.TimeTrackerOptions{Disabled: s.disabled}).Changed() {
			log("tracked " + s.timer.Total(c.Now()).String())
		}
	})
}

func projectsDemoPomodoroTimer(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "380", func() struct {
		timer            mujica.PomodoroState
		disabled, hidden bool
		seconds          float64
	} {
		return struct {
			timer            mujica.PomodoroState
			disabled, hidden bool
			seconds          float64
		}{seconds: 10}
	})
	projectsCheck(c, log, &s.disabled, "Disabled")
	projectsCheck(c, log, &s.hidden, "Hidden")
	if mujica.NumberInput(c, &s.seconds, mujica.NumberInputOptions{Min: 1, Max: 1500, Step: 1, Label: "Focus seconds", Disabled: s.timer.Running()}).Changed() {
		log(fmt.Sprintf("focus = %g seconds", s.seconds))
	}
	if s.hidden {
		return
	}
	r := mujica.PomodoroTimer(c, &s.timer, mujica.PomodoroTimerOptions{Focus: time.Duration(s.seconds) * time.Second, Break: 5 * time.Second, LongBreak: 10 * time.Second, LongEvery: 2, Disabled: s.disabled})
	if r.Changed() {
		log(fmt.Sprintf("timer phase = %d, running = %v", s.timer.Phase(), s.timer.Running()))
	}
	if phase, ok := r.Finished(); ok {
		log(fmt.Sprintf("finished phase %d", phase))
	}
}
