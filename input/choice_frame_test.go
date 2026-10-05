package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// Keys delivered in the frame a control or option turns disabled must change nothing.

func TestMultiSelectBackspaceAsItDisables(t *testing.T) {
	values, disabled := []int{1, 2}, false
	tt := choiceBetween(func(c *ui.Context) {
		MultiSelect(c, &values, multiSelectTestOptions, MultiSelectOptions{Disabled: disabled, Label: "Flowers"})
	})
	_ = tt.Click("Before")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Flowers") {
		t.Fatal("field not focused")
	}
	disabled = true
	tt.Key(0, ui.KeyBackspace)
	if !slices.Equal(values, []int{1, 2}) {
		t.Fatalf("Backspace removed a value: %v", values)
	}
}

func TestComboboxEnterAsItDisables(t *testing.T) {
	value, disabled := "", false
	tt := choiceBetween(func(c *ui.Context) {
		Combobox(c, &value, comboboxTestFonts, ComboboxOptions{AllowCustom: true, Disabled: disabled, Label: "Font"})
	})
	_ = tt.Click("Before")
	tt.Key(0, ui.KeyTab)
	tt.Type("Caslon")
	tt.Key(0, ui.KeyEscape)
	disabled = true
	tt.Key(0, ui.KeyEnter)
	if value != "" {
		t.Fatalf("Enter submitted %q", value)
	}
}

func TestCascaderRightAsRowDisables(t *testing.T) {
	var path []string
	f := newCascaderFake()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Cascader(c, &path, CascaderOptions{Level: f.level, Load: f.load, AnyLevel: true, Label: "Place"})
	}, 700, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Europe") {
		t.Fatal("Europe not focused")
	}
	roots := f.levels[""]
	roots.Nodes = slices.Clone(roots.Nodes)
	roots.Nodes[0].Disabled = true
	f.levels[""] = roots
	tt.Key(0, ui.KeyRight)
	if len(path) != 0 || len(f.loads) != 0 {
		t.Fatalf("Right on a disabled row chose %v, loaded %v", path, f.loads)
	}
}

func TestSelectEnterOnOptionAsItDisables(t *testing.T) {
	v := "apple"
	options := slices.Clone(selectTestOptions)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Select(c, &v, options, SelectOptions{Label: "Fruit"})
	}, 400, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	options[1].Disabled = true
	tt.Key(0, ui.KeyEnter)
	if v != "apple" {
		t.Fatalf("Enter chose the disabled option: %q", v)
	}
}

func TestEntryFocusSkipsDisabledFirstRow(t *testing.T) {
	var path []string
	f := &cascaderFake{levels: map[string]CascaderLevel{
		"": {Nodes: []CascaderNode{{Value: "aq", Label: "Antarctica", Disabled: true}, {Value: "eu", Label: "Europe", Leaf: true}}},
	}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Cascader(c, &path, CascaderOptions{Level: f.level, Load: f.load, Label: "Place"})
	}, 700, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Europe") {
		t.Fatal("cascader entry did not reach the first enabled row")
	}
	nodes := []TreeSelectNode{{ID: "a", Label: "Archived", Disabled: true}, {ID: "b", Label: "Brass"}, {ID: "c", Label: "Copper"}}
	var sel []string
	tt = ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TreeSelect(c, &sel, nodes, TreeSelectOptions{Label: "Metal"})
	}, 500, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Brass") {
		t.Fatal("tree entry did not reach the first enabled row")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Copper") {
		t.Fatal("arrows stopped working after entry")
	}
}
