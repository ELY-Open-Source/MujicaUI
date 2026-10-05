package git

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// commitGraphSample is main with a feature branch merged back.
func commitGraphSample() []Commit {
	cs := []Commit{
		{Hash: "m3", Parents: []string{"m2", "f2"}, Subject: "Merge feature", Refs: []string{"main"}},
		{Hash: "f2", Parents: []string{"f1"}, Subject: "Feature work"},
		{Hash: "m2", Parents: []string{"m1"}, Subject: "Fix on main"},
		{Hash: "f1", Parents: []string{"m1"}, Subject: "Start feature"},
		{Hash: "m1", Subject: "Initial commit"},
	}
	for i := range cs {
		cs[i].Author = []string{"Sakiko", "Mutsumi", "Uika"}[i%3]
		cs[i].When = gitTestTime.Add(-time.Duration(i) * 3 * time.Hour)
	}
	return cs
}

func TestCommitGraphLayoutLanes(t *testing.T) {
	rows, lanes := commitGraphLayout(commitGraphSample())
	if lanes != 2 {
		t.Fatalf("lanes %d, want 2", lanes)
	}
	want := []struct {
		node  int
		merge bool
	}{{0, true}, {1, false}, {0, false}, {1, false}, {0, false}}
	for i, w := range want {
		if rows[i].node != w.node || rows[i].merge != w.merge {
			t.Errorf("row %d: node %d merge %v, want %+v", i, rows[i].node, rows[i].merge, w)
		}
	}
	// The merge sends a line to the feature lane; f1 joins main again at m1.
	if l := rows[0].bottom; len(l) != 2 || l[1].from != 0 || l[1].to != 1 {
		t.Errorf("merge links %+v", l)
	}
	if l := rows[4].top; len(l) != 2 || l[1].from != 1 || l[1].to != 0 {
		t.Errorf("fork point links %+v", l)
	}
}

func TestCommitGraphVirtualized(t *testing.T) {
	commits := gitHistory(100000)
	var s CommitGraphState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CommitGraph(c, &s, commits, CommitGraphOptions{}).Element.Height(400)
	}, 700, 440)
	if n := len(tt.Texts()); n > 400 {
		t.Fatalf("%d texts built for 100,000 commits", n)
	}
	rows := &s.rows[0]
	tt.Scroll(300, 200, 0, 50000)
	tt.Frame()
	if &s.rows[0] != rows {
		t.Fatal("layout recomputed without a data change")
	}
	if tt.HasText("Commit number 100000") {
		t.Fatal("scroll did not move the window of rows")
	}
}

func TestCommitGraphChooseAndKeys(t *testing.T) {
	var s CommitGraphState
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := CommitGraph(c, &s, commitGraphSample(), CommitGraphOptions{})
		r.Element.Height(300)
		if r.Changed() {
			changes++
		}
	}, 700, 400)
	if err := tt.Click("Fix on main"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if key, _ := s.List.Selected(); key != "f1" || changes != 2 {
		t.Fatalf("selected %q changes %d", key, changes)
	}
}

func TestCommitGraphSnapshot(t *testing.T) {
	commits := commitGraphSample()
	snapshot(t, "commit_graph", 640, 240, func(c *ui.Context) {
		s := ui.Local(c.Root(), "graph", func() *CommitGraphState { return &CommitGraphState{} })
		ui.Column(c).Padding(16).Children(func() {
			CommitGraph(c, *s, commits, CommitGraphOptions{}).Element.Height(200)
		})
	})
}
