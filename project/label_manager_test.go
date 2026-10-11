package project

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/datetime"
)

func labelTester(labels *[]TaskLabel, sel *[]string, opts *LabelManagerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if LabelManager(c, labels, sel, *opts).Changed() {
				*n++
			}
		})
	}, 480, 480)
}

func TestLabelManagerSelectCreateDelete(t *testing.T) {
	labels, sel, n := slices.Clone(projLabels), []string{}, 0
	tt := labelTester(&labels, &sel, &LabelManagerOptions{}, &n)
	_ = tt.Click("UI")
	if !slices.Equal(sel, []string{"ui"}) {
		t.Fatalf("select %v", sel)
	}
	_ = tt.Click("New label")
	tt.Type("Docs")
	tt.Key(0, ui.KeyEnter)
	if len(labels) != 3 || !tt.HasText(`A label named "Docs" exists.`) {
		t.Fatalf("duplicate accepted: %v", labels)
	}
	tt.Command("selectAll")
	tt.Type("Performance")
	tt.Key(0, ui.KeyEnter)
	if len(labels) != 4 || labels[3].Name != "Performance" || labels[3].ID != "label-4" {
		t.Fatalf("create %+v", labels)
	}
	_ = tt.Click("Delete UI")
	if len(labels) != 3 || len(sel) != 0 || n != 3 {
		t.Fatalf("delete left %+v sel %v (%d)", labels, sel, n)
	}
}

func TestLabelManagerRenameAndColor(t *testing.T) {
	labels, sel, n := slices.Clone(projLabels), []string{}, 0
	tt := labelTester(&labels, &sel, &LabelManagerOptions{}, &n)
	_ = tt.Click("Edit Bug")
	_ = tt.Click("Green")
	_ = tt.Click("Label name")
	tt.Command("selectAll")
	tt.Type("UI")
	tt.Key(0, ui.KeyEnter)
	if labels[0].Color != datetime.EventSuccess || labels[0].Name != "Bug" || !tt.HasText(`A label named "UI" exists.`) {
		t.Fatalf("color/duplicate %+v", labels[0])
	}
	tt.Command("selectAll")
	tt.Type("Defect")
	tt.Key(0, ui.KeyEnter)
	if labels[0].Name != "Defect" || tt.HasText("Label name") {
		t.Fatalf("rename %+v", labels[0])
	}
}

func TestLabelManagerDisabledAndInvalid(t *testing.T) {
	labels, sel, n := slices.Clone(projLabels), []string{}, 0
	tt := labelTester(&labels, &sel, &LabelManagerOptions{Disabled: true}, &n)
	_ = tt.Click("UI")
	_ = tt.Click("Delete Bug")
	_ = tt.Click("Edit Bug")
	if n != 0 || len(labels) != 3 || len(sel) != 0 || tt.HasText("Label name") {
		t.Fatal("disabled manager changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an unknown selected label")
		}
	}()
	bad := []string{"ghost"}
	labelTester(&labels, &bad, &LabelManagerOptions{}, &n)
}

func TestLabelManagerSnapshot(t *testing.T) {
	labels, sel := slices.Clone(projLabels), []string{"bug"}
	snapshot(t, "label_manager", 400, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { LabelManager(c, &labels, &sel, LabelManagerOptions{}) })
	})
}

func TestLabelManagerLargeData(t *testing.T) {
	labels := make([]TaskLabel, 100000)
	for i := range labels {
		labels[i] = TaskLabel{ID: fmt.Sprint(i), Name: "Tag"}
	}
	selected, n := []string{}, 0
	tt := labelTester(&labels, &selected, &LabelManagerOptions{}, &n)
	if built := len(tt.Texts()); built > 200 {
		t.Fatalf("built %d labels", built)
	}
}
