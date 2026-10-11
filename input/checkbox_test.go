package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCheckboxLabelClickAndSpace(t *testing.T) {
	on, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Checkbox(c, &on, "Remember me", CheckboxOptions{}).Changed() {
			changes++
		}
	}, 300, 200)
	_ = tt.Click("Remember me")
	if !on || changes != 1 {
		t.Fatalf("label click: %v changes %d", on, changes)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if on || changes != 2 {
		t.Fatalf("Space: %v changes %d", on, changes)
	}
}

func TestCheckboxIndeterminateChecks(t *testing.T) {
	for _, start := range []bool{false, true} {
		on := start
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			Checkbox(c, &on, "Some", CheckboxOptions{Indeterminate: true})
		}, 300, 200)
		_ = tt.Click("Some")
		if !on {
			t.Fatalf("from mixed with %v, click left the box clear", start)
		}
	}
}

func TestCheckboxDisabled(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Checkbox(c, &on, "Off", CheckboxOptions{Disabled: true})
	}, 300, 200)
	_ = tt.Click("Off")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if on || tt.Focused("Off") {
		t.Fatalf("disabled box changed %v or took focus", on)
	}
}

func TestCheckboxGroupSelectAll(t *testing.T) {
	mail, cal, locked := true, false, false
	changes := 0
	items := []CheckboxItem{{Label: "Mail", Value: &mail}, {Label: "Calendar", Value: &cal}, {Label: "Admin", Value: &locked, Disabled: true}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if CheckboxGroup(c, "Notifications", items, CheckboxGroupOptions{}).Changed() {
			changes++
		}
	}, 300, 300)
	_ = tt.Click("Notifications")
	if !mail || !cal || locked || changes != 1 {
		t.Fatalf("select all: mail %v cal %v admin %v changes %d", mail, cal, locked, changes)
	}
	_ = tt.Click("Notifications")
	if mail || cal || locked {
		t.Fatalf("clear all: mail %v cal %v admin %v", mail, cal, locked)
	}
	_ = tt.Click("Calendar")
	if !cal || changes != 3 {
		t.Fatalf("item click: cal %v changes %d", cal, changes)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Notifications") {
		t.Fatal("Tab did not reach the select-all box")
	}
	tt.Key(0, ui.KeySpace)
	if !mail || !cal {
		t.Fatalf("Space on select-all: mail %v cal %v", mail, cal)
	}
}

func checkboxSnapshotView(c *ui.Context) {
	a, b, m, x, y := false, true, false, true, false
	ui.Row(c).Padding(24).Gap(32).AlignItems(ui.Start).Children(func() {
		ui.Column(c).Children(func() {
			Checkbox(c, &a, "Unchecked", CheckboxOptions{})
			Checkbox(c, &b, "Checked", CheckboxOptions{})
			Checkbox(c, &m, "Indeterminate", CheckboxOptions{Indeterminate: true})
			Checkbox(c, &b, "Disabled", CheckboxOptions{Disabled: true})
		})
		CheckboxGroup(c, "Notifications", []CheckboxItem{{Label: "Mail", Value: &x}, {Label: "Calendar", Value: &y}}, CheckboxGroupOptions{})
	})
}

func TestCheckboxSnapshot(t *testing.T) {
	snapshot(t, "checkbox", 420, 220, checkboxSnapshotView)
	buttonDensitySnapshots(t, "checkbox", 420, 260, checkboxSnapshotView)
}
