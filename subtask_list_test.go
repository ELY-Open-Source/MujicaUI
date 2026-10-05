package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func subtaskSample() []Subtask {
	return []Subtask{{ID: "a", Title: "Sketch", Done: true}, {ID: "b", Title: "Score"}, {ID: "c", Title: "Rehearse"}}
}

func subtaskTester(subs *[]Subtask, opts *SubtaskListOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if SubtaskList(c, subs, *opts).Changed() {
				*n++
			}
		})
	}, 480, 480)
}

func subtaskTitles(s []Subtask) string {
	out := ""
	for _, x := range s {
		out += x.ID
	}
	return out
}

func TestSubtaskListProgressCheckAdd(t *testing.T) {
	subs, n := subtaskSample(), 0
	tt := subtaskTester(&subs, &SubtaskListOptions{}, &n)
	if !tt.HasText("1 of 3 done") {
		t.Fatalf("progress %v", tt.Texts())
	}
	_ = tt.Click("Score")
	if !subs[1].Done || !tt.HasText("2 of 3 done") {
		t.Fatal("check")
	}
	_ = tt.Click("Add a subtask")
	tt.Type("Bow")
	tt.Key(0, ui.KeyEnter)
	if len(subs) != 4 || subs[3].Title != "Bow" || subs[3].ID != "sub-4" || !tt.HasText("2 of 4 done") {
		t.Fatalf("add %+v", subs)
	}
	_ = tt.Click("Remove Sketch")
	if subtaskTitles(subs) != "bcsub-4" || n != 3 {
		t.Fatalf("remove %v (%d)", subtaskTitles(subs), n)
	}
}

func TestSubtaskListReorder(t *testing.T) {
	subs, n := subtaskSample(), 0
	tt := subtaskTester(&subs, &SubtaskListOptions{}, &n)
	_ = tt.Click("Rehearse")
	tt.Key(ui.Alt, ui.KeyUp)
	if subtaskTitles(subs) != "acb" {
		t.Fatalf("Alt+Up %v", subtaskTitles(subs))
	}
	g, _ := tt.Find("Drag Sketch")
	target, _ := tt.Find("Drag Score")
	calDrag(tt, g.X+g.W/2, g.Y+g.H/2, target.X+40, target.Y+target.H/2)
	if subtaskTitles(subs) != "cba" {
		t.Fatalf("drag %v", subtaskTitles(subs))
	}
}

func TestSubtaskListDisabledAndInvalid(t *testing.T) {
	subs, n := subtaskSample(), 0
	tt := subtaskTester(&subs, &SubtaskListOptions{Disabled: true}, &n)
	_ = tt.Click("Score")
	_ = tt.Click("Remove Sketch")
	g, _ := tt.Find("Drag Sketch")
	calDrag(tt, g.X+g.W/2, g.Y+g.H/2, g.X+g.W/2, g.Y+120)
	if n != 0 || subtaskTitles(subs) != "abc" || subs[1].Done {
		t.Fatal("disabled list changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for repeated IDs")
		}
	}()
	bad := []Subtask{{ID: "x"}, {ID: "x"}}
	subtaskTester(&bad, &SubtaskListOptions{}, &n)
}

func TestSubtaskListSnapshot(t *testing.T) {
	subs := subtaskSample()
	snapshot(t, "subtask_list", 440, 240, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { SubtaskList(c, &subs, SubtaskListOptions{}) })
	})
}

func TestSubtaskListLargeData(t *testing.T) {
	subs := make([]Subtask, 100000)
	for i := range subs {
		subs[i] = Subtask{ID: fmt.Sprint(i), Title: "Step"}
	}
	n := 0
	tt := subtaskTester(&subs, &SubtaskListOptions{}, &n)
	if built := len(tt.Texts()); built > 200 {
		t.Fatalf("built %d labels", built)
	}
}
