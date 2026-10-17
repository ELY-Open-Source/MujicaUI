package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestFormFocusFollowsFieldAfterReorder(t *testing.T) {
	var st FormState
	first, second := "", ""
	swapped := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Form(c, &st, FormOptions{}, func(f *FormScope) {
			fields := []func(){
				func() {
					FormItem(c, f, "first", &first, FormItemOptions[string]{Label: "First"}, func() ui.Element {
						return TextInput(c, &first, TextInputOptions{Label: "First"}).Input
					})
				},
				func() {
					FormItem(c, f, "second", &second, FormItemOptions[string]{Label: "Second"}, func() ui.Element {
						return TextInput(c, &second, TextInputOptions{Label: "Second"}).Input
					})
				},
			}
			if swapped {
				slices.Reverse(fields)
			}
			for _, fn := range fields {
				fn()
			}
		})
	}, 500, 400)
	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	tt.Type("a")
	swapped = true
	tt.Frame()
	tt.Type("b")
	if first != "ab" || second != "" {
		t.Fatalf("first %q second %q, want typing to stay in First", first, second)
	}
}

func TestToggleGroupFocusFollowsValueAfterReorder(t *testing.T) {
	var on []string
	items := []ToggleItem[string]{{Value: "left", Label: "Left"}, {Value: "right", Label: "Right"}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToggleGroup(c, &on, items, ToggleGroupOptions{Multiple: true, Label: "Align"})
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	slices.Reverse(items)
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(on, []string{"left"}) {
		t.Fatalf("pressed %v, want [left]", on)
	}
}
