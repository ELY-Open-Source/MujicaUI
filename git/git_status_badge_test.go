package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestGitStatusBadgeLettersAndNames(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(8).Children(func() {
			for s := GitModified; s <= GitCopied; s++ {
				GitStatusBadge(c, s, GitStatusBadgeOptions{})
			}
		})
	}, 400, 80)
	for i, name := range []string{"Modified", "Added", "Deleted", "Untracked", "Renamed", "Copied"} {
		if _, ok := tt.Find(name); !ok {
			t.Errorf("no badge named %q", name)
		}
		if got := GitStatus(i).Letter(); got != string("MADURC"[i]) {
			t.Errorf("letter %d = %q", i, got)
		}
	}
	if !tt.HasText("R") || !tt.HasText("U") {
		t.Fatalf("letters not shown: %q", tt.Texts())
	}
}

func TestGitStatusBadgeLabelAndInvalid(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		GitStatusBadge(c, GitDeleted, GitStatusBadgeOptions{ShowLabel: true})
	}, 300, 80)
	if !tt.HasText("已删除") {
		t.Fatalf("label missing: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown status did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		GitStatusBadge(c, GitStatus(9), GitStatusBadgeOptions{})
	}, 100, 50, 1)
}

func TestGitStatusBadgeSnapshot(t *testing.T) {
	snapshot(t, "git_status_badge", 520, 90, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			ui.Row(c).Gap(10).Children(func() {
				for s := GitModified; s <= GitCopied; s++ {
					GitStatusBadge(c, s, GitStatusBadgeOptions{})
				}
			})
			ui.Row(c).Gap(14).Children(func() {
				for _, s := range []GitStatus{GitModified, GitAdded, GitDeleted} {
					GitStatusBadge(c, s, GitStatusBadgeOptions{ShowLabel: true})
				}
			})
		})
	})
}
