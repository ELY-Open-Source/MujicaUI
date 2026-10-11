package project

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type taskItemProbe struct {
	task             Task
	disabled         bool
	changes, openeds []string
}

func (p *taskItemProbe) tester() *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(8).Children(func() {
			r := TaskItem(c, &p.task, TaskItemOptions{Members: projMembers, Today: calendarToday, Disabled: p.disabled})
			if id, ok := r.Changed(); ok {
				p.changes = append(p.changes, id)
			}
			if id, ok := r.Opened(); ok {
				p.openeds = append(p.openeds, id)
			}
		})
	}, 600, 120)
}

func TestTaskItemFieldsAndCheck(t *testing.T) {
	p := &taskItemProbe{task: projTasks()[1]}
	tt := p.tester()
	for _, s := range []string{"Tune the organ", "Urgent", "Overdue by 2 days", "Mutsumi Wakaba"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !p.task.Done || len(p.changes) != 1 || tt.HasText("Overdue") {
		t.Fatalf("check: %+v", p.task)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(p.openeds) != 1 || p.openeds[0] != "AVE-2" {
		t.Fatalf("opened %v", p.openeds)
	}
}

func TestTaskItemInlineRename(t *testing.T) {
	p := &taskItemProbe{task: projTasks()[0]}
	tt := p.tester()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	tt.Command("selectAll")
	tt.Type("Write the finale")
	tt.Key(0, ui.KeyEnter)
	if p.task.Title != "Write the finale" || len(p.changes) != 1 || !tt.HasText("Write the finale") {
		t.Fatalf("rename %q", p.task.Title)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	tt.Type(" draft")
	tt.Key(0, ui.KeyEscape)
	if p.task.Title != "Write the finale" || !tt.HasText("Write the finale") {
		t.Fatalf("Escape kept %q", p.task.Title)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	tt.Command("selectAll")
	tt.Type("  ")
	tt.Key(0, ui.KeyEnter)
	if p.task.Title != "Write the finale" {
		t.Fatal("an empty title was saved")
	}
}

func TestTaskItemDisabled(t *testing.T) {
	p := &taskItemProbe{task: projTasks()[0], disabled: true}
	tt := p.tester()
	_ = tt.Click("Write the overture")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyF2)
	if p.task.Done || len(p.changes)+len(p.openeds) != 0 {
		t.Fatal("disabled item reacted")
	}
}

func TestTaskItemSnapshot(t *testing.T) {
	tasks := projTasks()
	snapshot(t, "task_item", 560, 220, func(c *ui.Context) {
		ui.Column(c).Padding(8).Children(func() {
			for i := range tasks {
				TaskItem(c, &tasks[i], TaskItemOptions{Members: projMembers, Today: calendarToday})
			}
		})
	})
}
