package project

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var kanbanColumns3 = []KanbanColumn{{Status: "todo"}, {Status: "doing", WIP: 3}, {Status: "review"}, {Status: "done"}}

type kanbanProbe struct {
	tasks    []Task
	cols     []KanbanColumn
	disabled bool
	moves    []KanbanMove
	opened   []string
}

func (p *kanbanProbe) tester() *ui.Tester {
	if p.cols == nil {
		p.cols = kanbanColumns3
	}
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Height(560).Padding(8).Children(func() {
			r := KanbanBoard(c, &p.tasks, KanbanBoardOptions{Columns: p.cols, Statuses: projStatuses, Members: projMembers,
				Labels: projLabels, Disabled: p.disabled})
			if m, ok := r.Moved(); ok {
				p.moves = append(p.moves, m)
			}
			if id, ok := r.Opened(); ok {
				p.opened = append(p.opened, id)
			}
		})
	}, 1240, 600)
}

// kanbanOrder lists IDs per status in slice order.
func kanbanOrder(tasks []Task, status string) string {
	var out []string
	for _, t := range tasks {
		if t.Status == status {
			out = append(out, t.ID)
		}
	}
	return strings.Join(out, ",")
}

func TestKanbanBoardDragAcrossAndWithin(t *testing.T) {
	p := &kanbanProbe{tasks: projTasks()}
	tt := p.tester()
	src, _ := tt.Find("AVE-1 Write the overture")
	dst, _ := tt.Find("AVE-4 Rehearse act two")
	calDrag(tt, src.X+20, src.Y+20, dst.X+20, dst.Y+10)
	if kanbanOrder(p.tasks, "review") != "AVE-1,AVE-4" || kanbanOrder(p.tasks, "todo") != "" || len(p.moves) != 1 ||
		p.moves[0] != (KanbanMove{ID: "AVE-1", From: "todo", To: "review", Index: 0}) {
		t.Fatalf("drag onto a card: %v %+v", kanbanOrder(p.tasks, "review"), p.moves)
	}
	src, _ = tt.Find("AVE-2 Tune the organ")
	col, _ := tt.Find("In Progress")
	calDrag(tt, src.X+20, src.Y+20, col.X+20, col.Y+400)
	if kanbanOrder(p.tasks, "doing") != "AVE-3,AVE-2" {
		t.Fatalf("drag to the column end: %v", kanbanOrder(p.tasks, "doing"))
	}
}

func TestKanbanBoardKeyboardMoves(t *testing.T) {
	p := &kanbanProbe{tasks: projTasks()}
	tt := p.tester()
	_ = tt.Click("AVE-2 Tune the organ")
	tt.Key(ui.Alt, ui.KeyDown)
	if kanbanOrder(p.tasks, "doing") != "AVE-3,AVE-2" || !tt.Focused("AVE-2 Tune the organ") {
		t.Fatalf("Alt+Down %v", kanbanOrder(p.tasks, "doing"))
	}
	tt.Key(ui.Alt, ui.KeyRight)
	if kanbanOrder(p.tasks, "review") != "AVE-4,AVE-2" || !tt.Focused("AVE-2 Tune the organ") {
		t.Fatalf("Alt+Right %v", kanbanOrder(p.tasks, "review"))
	}
	tt.Key(ui.Alt, ui.KeyUp)
	if kanbanOrder(p.tasks, "review") != "AVE-2,AVE-4" {
		t.Fatalf("Alt+Up %v", kanbanOrder(p.tasks, "review"))
	}
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 2 || p.opened[1] != "AVE-2" {
		t.Fatalf("opened %v", p.opened)
	}
}

func TestKanbanBoardWIPLimit(t *testing.T) {
	p := &kanbanProbe{tasks: projTasks(), cols: []KanbanColumn{{Status: "todo"}, {Status: "doing", WIP: 2}, {Status: "review"}, {Status: "done"}}}
	tt := p.tester()
	if !tt.HasText("2 of 2 cards") {
		t.Fatalf("WIP count %v", tt.Texts())
	}
	_ = tt.Click("AVE-1 Write the overture")
	tt.Key(ui.Alt, ui.KeyRight)
	if kanbanOrder(p.tasks, "todo") != "AVE-1" || !tt.HasText("In Progress is at its limit of 2 cards.") {
		t.Fatalf("WIP not enforced: %v", tt.Texts())
	}
	p.cols[1].WIP = 1
	tt.Frame()
	if !tt.HasText("1 over the limit") {
		t.Fatal("over-limit warning missing")
	}
}

func TestKanbanBoardDisabledAndInvalid(t *testing.T) {
	p := &kanbanProbe{tasks: projTasks(), disabled: true}
	tt := p.tester()
	src, _ := tt.Find("AVE-1 Write the overture")
	dst, _ := tt.Find("AVE-4 Rehearse act two")
	calDrag(tt, src.X+20, src.Y+20, dst.X+20, dst.Y+10)
	_ = tt.Click("AVE-1 Write the overture")
	tt.Key(ui.Alt, ui.KeyRight)
	if len(p.moves) != 0 || len(p.opened) != 0 || kanbanOrder(p.tasks, "todo") != "AVE-1" {
		t.Fatal("disabled board changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for a repeated column")
		}
	}()
	q := &kanbanProbe{tasks: projTasks(), cols: []KanbanColumn{{Status: "todo"}, {Status: "todo"}}}
	q.tester()
}

func TestKanbanBoardLargeData(t *testing.T) {
	p := &kanbanProbe{}
	for i := range 100000 {
		p.tasks = append(p.tasks, Task{ID: fmt.Sprintf("BIG-%d", i+1), Title: "Card", Status: kanbanColumns3[i%4].Status})
	}
	p.cols = []KanbanColumn{{Status: "todo"}, {Status: "doing"}, {Status: "review"}, {Status: "done"}}
	tt := p.tester()
	if n := len(tt.Texts()); n > 600 {
		t.Fatalf("built %d labels", n)
	}
}

func TestKanbanBoardSnapshot(t *testing.T) {
	tasks := projTasks()
	snapshot(t, "kanban_board", 1240, 460, func(c *ui.Context) {
		ui.Column(c).Height(460).Padding(16).Children(func() {
			KanbanBoard(c, &tasks, KanbanBoardOptions{Columns: kanbanColumns3, Statuses: projStatuses, Members: projMembers, Labels: projLabels})
		})
	})
}

func TestKanbanBoardColumnChanges(t *testing.T) {
	p := &kanbanProbe{tasks: projTasks(), cols: []KanbanColumn{{Status: "todo"}, {Status: "doing"}}}
	tt := p.tester()
	p.cols = []KanbanColumn{{Status: "done"}, {Status: "review"}, {Status: "doing"}}
	tt.Frame()
	if !tt.HasText("AVE-5 Book the hall") || !tt.HasText("AVE-4 Rehearse act two") {
		t.Fatal("columns retain old indexes")
	}
}

func TestKanbanBoardForeignDrop(t *testing.T) {
	left := []Task{{ID: "A-1", Title: "First", Status: "todo"}}
	right := []Task{{ID: "B-1", Title: "Second", Status: "todo"}}
	opts := KanbanBoardOptions{Columns: []KanbanColumn{{Status: "todo"}}, Statuses: projStatuses}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Height(400).Children(func() {
			ui.Box(c).Key("left").Width(300).Height(400).Children(func() { KanbanBoard(c, &left, opts) })
			ui.Box(c).Key("right").Width(300).Height(400).Children(func() { KanbanBoard(c, &right, opts) })
		})
	}, 640, 420)
	src, _ := tt.Find("A-1 First")
	dst, _ := tt.Find("B-1 Second")
	calDrag(tt, src.X+20, src.Y+20, dst.X+20, dst.Y+10)
	if len(left) != 1 || left[0].ID != "A-1" || len(right) != 1 || right[0].ID != "B-1" {
		t.Fatal("foreign drop changed another board")
	}
}
