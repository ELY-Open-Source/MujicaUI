package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var multiSelectTestOptions = []SelectOption[int]{
	{Value: 1, Label: "Rose"},
	{Value: 2, Label: "Lily"},
	{Value: 3, Label: "Iris"},
	{Value: 4, Label: "Violet", Disabled: true},
	{Value: 5, Label: "Lilac"},
}

func multiSelectTester(values *[]int, opts MultiSelectOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if MultiSelect(c, values, multiSelectTestOptions, opts).Changed() {
				*changes++
			}
		})
	}, 500, 500)
}

func TestMultiSelectCheckAndRemove(t *testing.T) {
	var values []int
	changes := 0
	tt := multiSelectTester(&values, MultiSelectOptions{Label: "Flowers"}, &changes)
	if err := tt.Click("Select…"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Iris")
	_ = tt.Click("Rose")
	_ = tt.Click("Violet")
	if !slices.Equal(values, []int{3, 1}) || changes != 2 {
		t.Fatalf("values %v changes %d", values, changes)
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Lily") {
		t.Fatal("Escape did not close")
	}
	if err := tt.Click("Remove Iris"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, []int{1}) {
		t.Fatalf("after remove %v", values)
	}
}

func TestMultiSelectSearchAndSelectAll(t *testing.T) {
	var values []int
	changes := 0
	tt := multiSelectTester(&values, MultiSelectOptions{}, &changes)
	_ = tt.Click("Select…")
	tt.Type("li")
	if tt.HasText("Rose") || !tt.HasText("Lily") || !tt.HasText("Lilac") {
		t.Fatalf("filter shows %v", tt.Texts())
	}
	_ = tt.Click("Select all")
	if !slices.Equal(values, []int{2, 5}) {
		t.Fatalf("select all of matches: %v", values)
	}
	_ = tt.Click("Select all")
	if len(values) != 0 {
		t.Fatalf("second select all left %v", values)
	}
	_ = tt.Click("Search")
	tt.Type("zz")
	if !tt.HasText("No results") {
		t.Fatal("no-results state missing")
	}
}

func TestMultiSelectLimit(t *testing.T) {
	values := []int{1}
	changes := 0
	tt := multiSelectTester(&values, MultiSelectOptions{Max: 2}, &changes)
	_ = tt.Click("Rose")
	_ = tt.Click("Select all")
	if !slices.Equal(values, []int{1, 2}) {
		t.Fatalf("select all past limit: %v", values)
	}
	if !tt.HasText("2 of 2 chosen") {
		t.Fatal("limit note missing")
	}
	_ = tt.Click("Lilac")
	if len(values) != 2 {
		t.Fatalf("limit broken: %v", values)
	}
}

func TestMultiSelectKeyboard(t *testing.T) {
	var values []int
	changes := 0
	tt := multiSelectTester(&values, MultiSelectOptions{Label: "Flowers"}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.Focused("Search") {
		t.Fatal("search not focused on open")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(values, []int{1, 2}) {
		t.Fatalf("values %v", values)
	}
	tt.Key(0, ui.KeyEscape)
	if !tt.Focused("Flowers") {
		t.Fatal("focus did not return to the field")
	}
	tt.Key(0, ui.KeyBackspace)
	if !slices.Equal(values, []int{1}) {
		t.Fatalf("Backspace left %v", values)
	}
}

func TestMultiSelectCollapsesTags(t *testing.T) {
	values := []int{1, 2, 3, 5}
	changes := 0
	tt := multiSelectTester(&values, MultiSelectOptions{MaxTags: 2}, &changes)
	if !tt.HasText("+2") || tt.HasText("Iris") {
		t.Fatalf("texts %v", tt.Texts())
	}
}

func TestMultiSelectSnapshot(t *testing.T) {
	snapshot(t, "multi_select", 460, 110, func(c *ui.Context) {
		var none []int
		some := []int{1, 2, 3, 5}
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			MultiSelect(c, &none, multiSelectTestOptions, MultiSelectOptions{})
			MultiSelect(c, &some, multiSelectTestOptions, MultiSelectOptions{MaxTags: 2})
		})
	})
}
