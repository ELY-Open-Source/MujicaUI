package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func tagSample() []GitTag {
	return []GitTag{
		{Name: "v1.1.0", Target: gitTestHash(42), Annotation: "Five hundred components"},
		{Name: "v1.0.0", Target: gitTestHash(7)},
	}
}

func TestTagListJumpAndSelect(t *testing.T) {
	var s data.ListState[string]
	jump := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := TagList(c, &s, tagSample(), TagListOptions{})
		r.Element.Height(200)
		if a, item, ok := r.Action(); ok && a == "jump" {
			jump = item
		}
	}, 600, 240)
	if !tt.HasText("Five hundred components") || !tt.HasText(gitShortHash(gitTestHash(42))) {
		t.Fatalf("row: %q", tt.Texts())
	}
	if err := tt.Click("v1.0.0"); err != nil {
		t.Fatal(err)
	}
	if key, _ := s.Selected(); key != "v1.0.0" {
		t.Fatalf("selected %q", key)
	}
	if err := tt.Click("Go to commit"); err != nil {
		t.Fatal(err)
	}
	if jump != gitTestHash(42) {
		t.Fatalf("jump %q", jump)
	}
}

func TestTagListInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var s data.ListState[string]
		TagList(c, &s, []GitTag{{Name: "v1"}}, TagListOptions{}).Element.Height(100)
	}, 300, 120, 1)
}

func TestTagListSnapshot(t *testing.T) {
	snapshot(t, "tag_list", 600, 140, func(c *ui.Context) {
		s := ui.Local(c.Root(), "list", func() *data.ListState[string] { return &data.ListState[string]{} })
		ui.Column(c).Padding(16).Children(func() {
			TagList(c, *s, tagSample(), TagListOptions{}).Element.Height(100)
		})
	})
}
