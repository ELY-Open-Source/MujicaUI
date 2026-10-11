package project

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func statusTester(v *string, opts *StatusSelectOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if StatusSelect(c, v, projStatuses, *opts).Changed() {
				*n++
			}
		})
	}, 400, 520)
}

func TestStatusSelectGroupsAndPick(t *testing.T) {
	v, n := "todo", 0
	tt := statusTester(&v, &StatusSelectOptions{}, &n)
	_ = tt.Click("Status")
	for _, g := range []string{"Backlog", "To do", "In progress", "Done", "Canceled"} {
		if !tt.HasText(g) {
			t.Fatalf("group %q missing: %v", g, tt.Texts())
		}
	}
	_ = tt.Click("In Review")
	if v != "review" || n != 1 {
		t.Fatalf("picked %q", v)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if v != "done" || n != 2 {
		t.Fatalf("keyboard picked %q", v)
	}
}

func TestStatusSelectDisabledAndInvalid(t *testing.T) {
	v, n := "todo", 0
	opts := &StatusSelectOptions{Disabled: true}
	tt := statusTester(&v, opts, &n)
	_ = tt.Click("Status")
	tt.Key(0, ui.KeyDown)
	if tt.HasText("In Review") || v != "todo" || n != 0 {
		t.Fatal("disabled select opened")
	}
	for _, bad := range []string{"nope", ""} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%q: expected panic", bad)
				}
			}()
			w := bad
			statusTester(&w, &StatusSelectOptions{}, &n)
		}()
	}
}

func TestStatusSelectSnapshot(t *testing.T) {
	values := []string{"backlog", "todo", "doing", "done", "dropped"}
	snapshot(t, "status_select", 440, 176, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(8).Wrap().Children(func() {
			for i := range values {
				StatusSelect(c, &values[i], projStatuses, StatusSelectOptions{})
			}
		})
	})
}

func TestStatusSelectAncestorDisableCloses(t *testing.T) {
	value, disabled := "todo", false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { StatusSelect(c, &value, projStatuses, StatusSelectOptions{}) })
	}, 400, 520)
	if err := tt.Click("Status"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	if tt.HasText("In Review") {
		t.Fatal("ancestor-disabled popup remains visible")
	}
	disabled = false
	tt.Frame()
	if tt.HasText("In Review") {
		t.Fatal("popup reopens after enabling")
	}
}
