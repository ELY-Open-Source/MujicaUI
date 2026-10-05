package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var selectTestOptions = []SelectOption[string]{
	{Value: "apple", Label: "Apple"},
	{Value: "banana", Label: "Banana"},
	{Value: "blueberry", Label: "Blueberry"},
	{Value: "cherry", Label: "Cherry", Disabled: true},
	{Value: "date", Label: "Date"},
}

func selectTester(value *string, options []SelectOption[string], opts SelectOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if Select(c, value, options, opts).Changed() {
				*changes++
			}
		})
	}, 400, 400)
}

func TestSelectClickChooses(t *testing.T) {
	v, changes := "", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{Label: "Fruit"}, &changes)
	if !tt.HasText("Select…") {
		t.Fatal("placeholder missing")
	}
	if err := tt.Click("Select…"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Banana"); err != nil {
		t.Fatal(err)
	}
	if v != "banana" || changes != 1 {
		t.Fatalf("value %q changes %d", v, changes)
	}
	if tt.HasText("Apple") {
		t.Fatal("popup stayed open")
	}
}

func TestSelectDisabledOptionNotChosen(t *testing.T) {
	v, changes := "apple", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{}, &changes)
	_ = tt.Click("Apple")
	_ = tt.Click("Cherry")
	if v != "apple" {
		t.Fatalf("disabled option chosen: %q", v)
	}
	// The arrows skip it: Banana, Blueberry, then Date.
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	for range 3 {
		tt.Key(0, ui.KeyDown)
	}
	tt.Key(0, ui.KeyEnter)
	if v != "date" {
		t.Fatalf("arrows chose %q, want date", v)
	}
}

func TestSelectKeyboard(t *testing.T) {
	v, changes := "apple", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.HasText("Banana") {
		t.Fatal("Down did not open the popup")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if v != "banana" {
		t.Fatalf("value %q", v)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Date") || v != "banana" {
		t.Fatalf("Escape: open=%v value %q", tt.HasText("Date"), v)
	}
}

func TestSelectTypeahead(t *testing.T) {
	v, changes := "apple", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.TypeKey(0, ui.KeyB, "b")
	if v != "banana" {
		t.Fatalf("b chose %q", v)
	}
	tt.TypeKey(0, ui.KeyB, "b")
	if v != "blueberry" {
		t.Fatalf("second b chose %q", v)
	}
	tt.TypeKey(0, ui.KeyC, "c")
	if v != "blueberry" {
		t.Fatalf("typed bbc chose %q", v)
	}
}

func TestSelectClear(t *testing.T) {
	v, changes := "date", 0
	tt := selectTester(&v, selectTestOptions, SelectOptions{Clearable: true}, &changes)
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if v != "" || tt.HasText("Banana") {
		t.Fatalf("value %q after clear, popup open %v", v, tt.HasText("Banana"))
	}
	v = "apple"
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyBackspace)
	if v != "" || changes != 2 {
		t.Fatalf("Backspace left %q, changes %d", v, changes)
	}
}

func TestSelectEmptyAndDisabled(t *testing.T) {
	v, changes := "", 0
	tt := selectTester(&v, nil, SelectOptions{}, &changes)
	_ = tt.Click("Select…")
	if !tt.HasText("No options") {
		t.Fatal("empty state missing")
	}
	v2 := "apple"
	tt = selectTester(&v2, selectTestOptions, SelectOptions{Disabled: true}, &changes)
	_ = tt.Click("Apple")
	if tt.HasText("Banana") {
		t.Fatal("disabled select opened")
	}
}

func TestSelectSnapshot(t *testing.T) {
	snapshot(t, "select", 660, 90, func(c *ui.Context) {
		a, b, d := "", "banana", "apple"
		ui.Row(c).Padding(20).Gap(16).Wrap().Children(func() {
			Select(c, &a, selectTestOptions, SelectOptions{})
			Select(c, &b, selectTestOptions, SelectOptions{Clearable: true})
			Select(c, &d, selectTestOptions, SelectOptions{Disabled: true})
		})
	})
}
