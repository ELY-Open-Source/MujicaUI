package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func branchSample() []Branch {
	return []Branch{
		{Name: "main", Current: true, Upstream: "origin/main", Ahead: 2, Behind: 1},
		{Name: "feature/graph", Upstream: "origin/feature/graph"},
		{Name: "fix/blame"},
		{Name: "origin/main", Remote: true},
		{Name: "origin/release", Remote: true},
	}
}

func TestBranchListRowsAndActions(t *testing.T) {
	var s ListState[string]
	var acted [][2]string
	actions := []GitAction{{ID: "checkout", Label: "Check out", Icon: "git-branch"}, {ID: "delete", Label: "Delete", Icon: "trash-2"}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := BranchList(c, &s, branchSample(), BranchListOptions{Actions: actions})
		r.Element.Height(300)
		if a, item, ok := r.Action(); ok {
			acted = append(acted, [2]string{a, item})
		}
	}, 700, 340)
	if !tt.HasText("current") || !tt.HasText("origin/feature/graph") {
		t.Fatalf("current, ahead, behind: %q", tt.Texts())
	}
	for _, name := range []string{"2 ahead of upstream", "1 behind upstream"} {
		if _, ok := tt.Find(name); !ok {
			t.Fatalf("%q not shown", name)
		}
	}
	r, _ := tt.Find("fix/blame")
	if err := tt.Click("fix/blame"); err != nil {
		t.Fatal(err)
	}
	if key, _ := s.Selected(); key != "fix/blame" {
		t.Fatalf("selected %q", key)
	}
	// The third row's Delete button sits at the row's right end.
	tt.ClickAt(700-60, r.Y+r.H/2)
	if len(acted) != 1 || acted[0] != [2]string{"delete", "fix/blame"} {
		t.Fatalf("actions %q", acted)
	}
}

func TestBranchListInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		b []Branch
		a []GitAction
	}{
		"name":   {[]Branch{{}}, nil},
		"ahead":  {[]Branch{{Name: "x", Ahead: -1}}, nil},
		"action": {nil, []GitAction{{ID: "x"}}},
		"dup":    {nil, []GitAction{{ID: "x", Label: "X", Icon: "x"}, {ID: "x", Label: "Y", Icon: "x"}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				var s ListState[string]
				BranchList(c, &s, tc.b, BranchListOptions{Actions: tc.a}).Element.Height(100)
			}, 300, 120, 1)
		})
	}
}

func TestBranchListSnapshot(t *testing.T) {
	snapshot(t, "branch_list", 600, 260, func(c *ui.Context) {
		s := ui.Local(c.Root(), "list", func() *ListState[string] { return &ListState[string]{} })
		ui.Column(c).Padding(16).Children(func() {
			BranchList(c, *s, branchSample(), BranchListOptions{Actions: []GitAction{{ID: "checkout", Label: "Check out", Icon: "git-branch"}}}).Element.Height(220)
		})
	})
}
