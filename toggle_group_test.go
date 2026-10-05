package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

var toggleGroupTestItems = []ToggleItem[string]{
	{Value: "left", Label: "Left"},
	{Value: "center", Label: "Center"},
	{Value: "right", Label: "Right"},
	{Value: "justify", Label: "Justify", Disabled: true},
}

func toggleGroupTester(value *[]string, opts ToggleGroupOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if ToggleGroup(c, value, toggleGroupTestItems, opts).Changed() {
			*changes++
		}
	}, 500, 120)
}

func TestToggleGroupSingle(t *testing.T) {
	value, changes := []string{"left"}, 0
	tt := toggleGroupTester(&value, ToggleGroupOptions{}, &changes)
	if err := tt.Click("Center"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(value, []string{"center"}) || changes != 1 {
		t.Fatalf("value %v changes %d", value, changes)
	}
	// The pressed item stays without AllowEmpty.
	_ = tt.Click("Center")
	if !slices.Equal(value, []string{"center"}) || changes != 1 {
		t.Fatalf("released without AllowEmpty: %v", value)
	}
}

func TestToggleGroupAllowEmpty(t *testing.T) {
	value, changes := []string{"left"}, 0
	tt := toggleGroupTester(&value, ToggleGroupOptions{AllowEmpty: true}, &changes)
	_ = tt.Click("Left")
	if len(value) != 0 || changes != 1 {
		t.Fatalf("value %v changes %d", value, changes)
	}
}

func TestToggleGroupMultiple(t *testing.T) {
	value, changes := []string{"left"}, 0
	tt := toggleGroupTester(&value, ToggleGroupOptions{Multiple: true}, &changes)
	_ = tt.Click("Right")
	if !slices.Equal(value, []string{"left", "right"}) {
		t.Fatalf("value %v", value)
	}
	_ = tt.Click("Left")
	_ = tt.Click("Right")
	if !slices.Equal(value, []string{"right"}) || changes != 2 {
		t.Fatalf("last item released: %v, changes %d", value, changes)
	}
	_ = tt.Click("Justify")
	if slices.Contains(value, "justify") {
		t.Fatal("disabled item pressed")
	}
}

func TestToggleGroupKeyboard(t *testing.T) {
	value, changes := []string{"left"}, 0
	tt := toggleGroupTester(&value, ToggleGroupOptions{}, &changes)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Left") {
		t.Fatal("Tab did not enter at the pressed item")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(value, []string{"right"}) {
		t.Fatalf("value %v after arrows and Space", value)
	}
	// The disabled item is skipped: Right wraps to the first.
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Left") {
		t.Fatal("arrow did not skip the disabled item")
	}
}

func TestToggleGroupRejectsTwoValuesInSingleMode(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	value := []string{"left", "right"}
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		ToggleGroup(c, &value, toggleGroupTestItems, ToggleGroupOptions{})
	}, 300, 100, 1)
}

func TestToggleGroupSnapshot(t *testing.T) {
	snapshot(t, "toggle_group", 520, 150, func(c *ui.Context) {
		single, multi := []string{"center"}, []string{"left", "right"}
		ui.Column(c).Padding(20).Gap(16).Children(func() {
			ToggleGroup(c, &single, toggleGroupTestItems, ToggleGroupOptions{Label: "Align"})
			ToggleGroup(c, &multi, toggleGroupTestItems, ToggleGroupOptions{Multiple: true})
		})
	})
}
