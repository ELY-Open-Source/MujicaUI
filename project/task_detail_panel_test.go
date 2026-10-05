package project

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type detailProbe struct {
	task             Task
	disabled         bool
	version          int
	changes, closes  int
	activityRendered bool
}

func (p *detailProbe) tester() *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Children(func() {
			ui.Box(c).Width(420).Children(func() {
				r := TaskDetailPanel(c, &p.task, TaskDetailPanelOptions{Statuses: projStatuses, Members: projMembers, Labels: projLabels,
					Today: calendarToday, Disabled: p.disabled, Version: p.version, Activity: func() {
						p.activityRendered = true
						ui.Text(c, "Sakiko opened this task")
					}})
				if r.Changed() {
					p.changes++
				}
				if r.Closed() {
					p.closes++
				}
			})
		})
	}, 640, 1200)
}

func TestTaskDetailPanelEditsFields(t *testing.T) {
	p := &detailProbe{task: projTasks()[0]}
	tt := p.tester()
	if !tt.HasText("Sakiko opened this task") || !tt.HasText("Activity") {
		t.Fatal("activity slot missing")
	}
	_ = tt.Click("Todo")
	_ = tt.Click("In Review")
	_ = tt.Click("Tomorrow")
	_ = tt.Click("Add a subtask")
	tt.Type("Copy parts")
	tt.Key(0, ui.KeyEnter)
	if p.task.Status != "review" || p.task.Due != oct(6) || len(p.task.Subtasks) != 1 || p.changes != 3 {
		t.Fatalf("edited %+v (%d changes)", p.task, p.changes)
	}
	_ = tt.Click("Title")
	tt.Command("selectAll")
	tt.Type("Write the requiem")
	if p.task.Title != "Write the requiem" {
		t.Fatalf("title %q", p.task.Title)
	}
	_ = tt.Click("Close details")
	if p.closes != 1 {
		t.Fatal("close not reported")
	}
}

func TestTaskDetailPanelDisabled(t *testing.T) {
	p := &detailProbe{task: projTasks()[0], disabled: true}
	tt := p.tester()
	_ = tt.Click("Tomorrow")
	_ = tt.Click("Title")
	tt.Type("x")
	_ = tt.Click("Close details")
	if p.changes != 0 || p.closes != 0 || p.task.Title != "Write the overture" {
		t.Fatalf("disabled panel changed %+v", p.task)
	}
}

func TestTaskDetailPanelSnapshot(t *testing.T) {
	task := projTasks()[1]
	task.Subtasks = []Subtask{{ID: "s1", Title: "Replace the pipes", Done: true}, {ID: "s2", Title: "Tune by ear"}}
	snapshot(t, "task_detail_panel", 460, 900, func(c *ui.Context) {
		ui.Box(c).Width(420).Margin(16).Children(func() {
			TaskDetailPanel(c, &task, TaskDetailPanelOptions{Statuses: projStatuses, Members: projMembers, Labels: projLabels,
				Today: calendarToday, Activity: func() { ui.Text(c, "Mutsumi moved this to In Progress") }})
		})
	})
}

func TestTaskDetailPanelAncestorDisableClosesLabels(t *testing.T) {
	task, disabled := projTasks()[0], false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			TaskDetailPanel(c, &task, TaskDetailPanelOptions{Statuses: projStatuses, Members: projMembers, Labels: projLabels})
		})
	}, 600, 1000)
	if err := tt.Click("Docs"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Select all") {
		t.Fatal("labels popup did not open")
	}
	disabled = true
	tt.Frame()
	if tt.HasText("Select all") {
		t.Fatal("disabled panel retained its labels popup")
	}
	if taskSnapshot(task) != taskSnapshot(projTasks()[0]) {
		t.Fatal("disabled panel changed task fields")
	}
}

func TestTaskDetailPanelEscapeClosesPopupFirst(t *testing.T) {
	p := &detailProbe{task: projTasks()[0]}
	tt := p.tester()
	if err := tt.Click("Todo"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("In Review") {
		t.Fatal("status popup did not open")
	}
	tt.Key(0, ui.KeyEscape)
	if p.closes != 0 || tt.HasText("In Review") {
		t.Fatalf("popup escape closes = %d, popup visible = %v", p.closes, tt.HasText("In Review"))
	}
	tt.Key(0, ui.KeyEscape)
	if p.closes != 1 {
		t.Fatalf("panel escape closes = %d", p.closes)
	}
}

// taskSnapshot compares editable task fields.
func taskSnapshot(t Task) string {
	return fmt.Sprintf("%q|%q|%d|%v|%q|%q|%v", t.Title, t.Status, t.Priority, t.Due, t.Assignees, t.Labels, t.Subtasks)
}

func TestTaskDetailPanelVersionUpdatesSubtasks(t *testing.T) {
	p := &detailProbe{task: projTasks()[0]}
	p.task.Subtasks = []Subtask{{ID: "one", Title: "One"}}
	tt := p.tester()
	p.task.Subtasks[0].Done = true
	p.version++
	tt.Frame()
	if !tt.HasText("1 of 1 done") || p.changes != 0 {
		t.Fatalf("external update: changes %d, labels %v", p.changes, tt.Texts())
	}
}
