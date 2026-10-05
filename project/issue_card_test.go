package project

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func issueCardTester(task Task, opts IssueCardOptions, n *int) *ui.Tester {
	opts.Statuses, opts.Members, opts.Labels = projStatuses, projMembers, projLabels
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Width(300).Children(func() {
			if IssueCard(c, task, opts).Clicked() {
				*n++
			}
		})
	}, 360, 260)
}

func TestIssueCardShowsFieldsAndActivates(t *testing.T) {
	n := 0
	tt := issueCardTester(projTasks()[1], IssueCardOptions{}, &n)
	for _, s := range []string{"AVE-2", "Tune the organ", "Bug", "Urgent", "Mutsumi Wakaba", "Uika Misumi"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	_ = tt.Click("AVE-2 Tune the organ")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 2 {
		t.Fatalf("activations %d", n)
	}
}

func TestIssueCardDisabledAndInvalid(t *testing.T) {
	n := 0
	tt := issueCardTester(projTasks()[0], IssueCardOptions{Disabled: true}, &n)
	_ = tt.Click("AVE-1 Write the overture")
	if n != 0 {
		t.Fatal("disabled card clicked")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an unknown status")
		}
	}()
	bad := projTasks()[0]
	bad.Status = "lost"
	issueCardTester(bad, IssueCardOptions{}, &n)
}

func TestIssueCardSnapshot(t *testing.T) {
	tasks := projTasks()
	snapshot(t, "issue_card", 640, 200, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).AlignItems(ui.Start).Children(func() {
			for i, task := range tasks[:3] {
				ui.Box(c).Width(190).Children(func() {
					IssueCard(c, task, IssueCardOptions{Statuses: projStatuses, Members: projMembers, Labels: projLabels, Selected: i == 1})
				})
			}
		})
	})
}

func TestIssueCardFocusRing(t *testing.T) {
	n := 0
	tt := issueCardTester(projTasks()[0], IssueCardOptions{}, &n)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	r, ok := tt.Find("AVE-1 Write the overture")
	if !ok {
		t.Fatal("card missing")
	}
	img, found := tt.Image(), false
	for x := int(r.X) - 5; x < int(r.X)-1; x++ {
		for y := int(r.Y) + 8; y < int(r.Y+r.H)-8; y++ {
			if datetimeNear(img.RGBAAt(x, y), theme.Light().Focus) {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("card keyboard focus has no outer ring")
	}
}
