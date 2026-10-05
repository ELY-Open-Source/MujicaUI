package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"mujicaui/theme"
)

type taskListProbe struct {
	tasks          []Task
	filter         TaskFilter
	opts           TaskListOptions
	changed, opens []string
}

func (p *taskListProbe) tester() *ui.Tester {
	p.opts.Members, p.opts.Statuses, p.opts.Today = projMembers, projStatuses, calendarToday
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(560).Padding(8).Children(func() {
			r := TaskList(c, p.tasks, &p.filter, p.opts)
			if id, ok := r.Changed(); ok {
				p.changed = append(p.changed, id)
			}
			if id, ok := r.Opened(); ok {
				p.opens = append(p.opens, id)
			}
		})
	}, 640, 580)
}

func textIndex(texts []string, s string) int {
	for i, x := range texts {
		if x == s {
			return i
		}
	}
	return -1
}

func TestTaskListGroupsAndFilter(t *testing.T) {
	p := &taskListProbe{tasks: projTasks(), opts: TaskListOptions{GroupBy: TaskGroupStatus}}
	tt := p.tester()
	tx := tt.Texts()
	todo, doing, review, done := textIndex(tx, "Todo"), textIndex(tx, "In Progress"), textIndex(tx, "In Review"), textIndex(tx, "Done")
	if !(todo >= 0 && todo < doing && doing < review && review < done) || textIndex(tx, "Backlog") >= 0 {
		t.Fatalf("groups %d %d %d %d in %v", todo, doing, review, done, tx)
	}
	_ = tt.Click("Open 4")
	if p.filter != TaskShowOpen || tt.HasText("Book the hall") || !tt.HasText("Tune the organ") {
		t.Fatalf("open filter %v", tt.Texts())
	}
	_ = tt.Click("Done 1")
	if tt.HasText("Tune the organ") || !tt.HasText("Book the hall") {
		t.Fatal("done filter")
	}
	_ = tt.Click("Book the hall")
	if len(p.opens) != 1 || p.opens[0] != "AVE-5" {
		t.Fatalf("opened %v", p.opens)
	}
}

func TestTaskListPriorityGroupsAndEdit(t *testing.T) {
	p := &taskListProbe{tasks: projTasks(), opts: TaskListOptions{GroupBy: TaskGroupPriority}}
	tt := p.tester()
	tx := tt.Texts()
	if u, h := textIndex(tx, "Urgent"), textIndex(tx, "High"); u < 0 || u > h {
		t.Fatalf("priority order %v", tx)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if len(p.changed) != 1 || !p.tasks[1].Done {
		t.Fatalf("keyboard check %v %+v", p.changed, p.tasks[1])
	}
}

func TestTaskListDisabledEmptyInvalid(t *testing.T) {
	p := &taskListProbe{tasks: projTasks(), opts: TaskListOptions{Disabled: true}}
	tt := p.tester()
	_ = tt.Click("Done 1")
	_ = tt.Click("Write the overture")
	if p.filter != TaskShowAll || len(p.opens) != 0 {
		t.Fatal("disabled list reacted")
	}
	q := &taskListProbe{}
	if tt := q.tester(); !tt.HasText("No tasks here.") {
		t.Fatal("no empty state")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	bad := &taskListProbe{filter: 7}
	bad.tester()
}

func TestTaskListLargeData(t *testing.T) {
	p := &taskListProbe{opts: TaskListOptions{GroupBy: TaskGroupPriority}}
	for i := range 100000 {
		p.tasks = append(p.tasks, Task{ID: fmt.Sprintf("BIG-%d", i+1), Title: "Task", Status: "todo", Priority: Priority(i % 5)})
	}
	tt := p.tester()
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("built %d labels", n)
	}
}

func TestTaskListSnapshot(t *testing.T) {
	tasks := projTasks()
	filter := TaskShowAll
	snapshot(t, "task_list", 600, 420, func(c *ui.Context) {
		ui.Column(c).Height(420).Padding(8).Children(func() {
			TaskList(c, tasks, &filter, TaskListOptions{GroupBy: TaskGroupStatus, Statuses: projStatuses, Members: projMembers, Today: calendarToday})
		})
	})
}

func TestTaskListCompletionRefreshesCounts(t *testing.T) {
	p := &taskListProbe{tasks: projTasks()}
	tt := p.tester()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !tt.HasText("Open 3") || !tt.HasText("Done 2") {
		t.Fatal("completion left cached filter counts")
	}
}

func TestTaskListFilterHitAreas(t *testing.T) {
	for _, density := range []theme.Density{theme.Compact, theme.Default, theme.Comfortable} {
		tasks, filter := projTasks(), TaskShowAll
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Density: density})
			ui.Box(c).Height(300).Children(func() { TaskList(c, tasks, &filter, TaskListOptions{Members: projMembers}) })
		}, 600, 320)
		for _, label := range []string{"All 5", "Open 4", "Done 1"} {
			r, ok := tt.Find(label)
			if !ok || r.W < 28 || r.H < 28 {
				t.Fatalf("density %v, %s: bounds %v", density, label, r)
			}
		}
	}
}
